package docker

import (
	"archive/tar"
	"bytes"
	"constester-go/internal/repository"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"io"
	"regexp"
	"strings"
	"time"
)

type ClientDocker struct {
	Client *client.Client
}

func NewClientDocker() (*ClientDocker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	return &ClientDocker{
		Client: cli,
	}, err
}

// RunTestsCPP runs the c++ code tests in the task by its id
func (cl *ClientDocker) RunTestsCPP(ctx context.Context, tests []repository.Test, duration time.Duration, code []byte) (any, error) {
	containerName := GenerateContainerName()
	imageName := "cxx-image"

	create, err := cl.Client.ContainerCreate(ctx, &container.Config{
		Image:        imageName,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
	}, nil, nil, nil, containerName)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := cl.Client.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			panic(err)
		}
	}()

	if err := cl.Client.ContainerStart(ctx, create.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	// inserting code
	err = cl.InsertFile(ctx, code, "main.cpp", containerName)
	if err != nil {
		return nil, err
	}

	// compile code
	execConfig := types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          []string{"g++", "-o", "code", "main.cpp"},
	}
	r, err := cl.Client.ContainerExecCreate(ctx, containerName, execConfig)
	if err != nil {
		return nil, err
	}
	err = cl.Client.ContainerExecStart(context.Background(), r.ID, types.ExecStartCheck{
		Tty: true,
	})
	if err != nil {
		return nil, err
	}

	//testing code
	for _, test := range tests {
		output, err := cl.RunCodeInContainer(containerName, test.Input, duration)
		if err != nil {
			return nil, err
		}

		//todo : add difference between error and wrong test

		if output != test.Expected {
			msg := fmt.Sprintf("Failed on test %d", test.Name)
			return msg, nil
		}

	}

	log.Info("successful testing")

	return "You got it", nil
}

// RunCodeInContainer runs the code in container with containerID
func (cl *ClientDocker) RunCodeInContainer(containerID, input string, duration time.Duration) (string, error) {
	execConfig := types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{"./code"},
	}

	newResp, err := cl.Client.ContainerExecCreate(context.Background(), containerID, execConfig)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	attach, err := cl.Client.ContainerExecAttach(ctx, newResp.ID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}
	defer attach.Close()

	inputData := []byte(input)
	_, err = attach.Conn.Write(append(inputData, '\n'))
	if err != nil {
		return "", err
	}

	// reading input
	var buf bytes.Buffer
	_, err = io.Copy(&buf, attach.Reader)
	if err != nil {
		return "", err
	}

	reg := regexp.MustCompile(`[[:cntrl:]]`)
	out := reg.ReplaceAllString(buf.String(), "")

	return out[GetTestInputLength(input):], nil
}

// InsertFile inserts file to app catalog of docker container
func (cl *ClientDocker) InsertFile(ctx context.Context, fileContents []byte, destPath string, containerName string) error {
	tarReader, err := cl.CreateTarArchive(fileContents, destPath)
	if err != nil {
		return err
	}

	err = cl.Client.CopyToContainer(ctx, containerName, "/app", tarReader, types.CopyToContainerOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (cl *ClientDocker) CreateTarArchive(fileContents []byte, destPath string) (*bytes.Reader, error) {
	tarBuf := new(bytes.Buffer)
	tw := tar.NewWriter(tarBuf)

	hdr := &tar.Header{
		Name: destPath,
		Mode: 0700,
		Size: int64(len(fileContents)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return nil, err
	}

	if _, err := tw.Write(fileContents); err != nil {
		return nil, err
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}

	return bytes.NewReader(tarBuf.Bytes()), nil
}

func GenerateContainerName() string {
	//TODO: just use any string generator
	return "id"
}

func GetTestInputLength(input string) int {
	a := strings.Count(input, "\n")
	b := len(input)
	return b - a
}
