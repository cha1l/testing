package docker

import (
	"constester-go/internal/repository"
	"context"
	"time"
)

type Code struct {
	TaskName string `json:"task_name"`
	Code     string `json:"code"`
}

type Container interface {
	RunTestsCPP(ctx context.Context, tests []repository.Test, duration time.Duration, code []byte) (any, error)
}

type Image struct {
	Container
}

func NewImages(client *ClientDocker) *Image {
	return &Image{
		Container: client,
	}
}
