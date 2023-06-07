package docker

import (
	"constester-go/internal/repository"
	"context"
	"time"
)

type Container interface {
	RunTestsCPP(ctx context.Context, tests []repository.Test, duration time.Duration) (any, error)
}

type Image struct {
	Container
}

func NewImages(client *ClientDocker) *Image {
	return &Image{
		Container: client,
	}
}
