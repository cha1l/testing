package docker

import (
	"constester-go/internal/repository"
	"context"
	"time"
)

type Code struct {
	TaskName string `json:"task_name"`
	Language string `json:"language"`
	Code     string `json:"code"`
}

type TestResult struct {
	Code string `json:"code"`
	Info string `json:"info"`
}

type Container interface {
	RunTestsCPP(ctx context.Context, tests []repository.Test, duration time.Duration, code []byte) (TestResult, error)
}

type Image struct {
	Container
}

func NewImages(client *ClientDocker) *Image {
	return &Image{
		Container: client,
	}
}
