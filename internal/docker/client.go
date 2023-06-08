package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"strings"
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
