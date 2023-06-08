package docker

import (
	"bytes"
	"constester-go/internal/repository"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"io"
	"regexp"
	"time"
)

// RunTestsCPP runs the c++ code tests in the task by its id
func (cl *ClientDocker) RunTestsCPP(ctx context.Context, tests []repository.Test, duration time.Duration, code []byte) (TestResult, error) {
	containerName := GenerateContainerName()
	imageName := "cxx-image"
	var emptyResult TestResult

	create, err := cl.Client.ContainerCreate(ctx, &container.Config{
		Image:        imageName,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
	}, nil, nil, nil, containerName)
	if err != nil {
		return emptyResult, err
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
		return emptyResult, err
	}

	// compile code
	execConfig := types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          []string{"g++", "-o", "code", "main.cpp"},
	}
	r, err := cl.Client.ContainerExecCreate(ctx, containerName, execConfig)
	if err != nil {
		return emptyResult, err
	}
	//err = cl.Client.ContainerExecStart(context.Background(), r.ID, types.ExecStartCheck{
	//	Tty: true,
	//})
	//if err != nil {
	//	return nil, err
	//}

	attach, err := cl.Client.ContainerExecAttach(ctx, r.ID, types.ExecStartCheck{})
	if err != nil {
		return emptyResult, err
	}
	defer attach.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, attach.Reader)
	if err != nil {
		return emptyResult, err
	}

	if buf.String() != "" {
		res := TestResult{
			Code: "Compilation Error",
			Info: buf.String(),
		}
		return res, nil
	}

	//testing code
	for _, test := range tests {
		output, err := cl.RunCodeInContainer(containerName, test.Input, duration)
		if err != nil {
			return emptyResult, err
		}

		if output != test.Expected {
			res := TestResult{
				Code: "Wrong Answer",
				Info: fmt.Sprintf("test %d", test.Name),
			}
			return res, nil
		}

	}

	res := TestResult{
		Code: "OK",
	}

	return res, nil
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
