package service

import (
	"constester-go/internal/docker"
	"constester-go/internal/repository"
	"context"
)

type Tasks interface {
	AddTask(task repository.Task) error
	RunTestsCPP(ctx context.Context, code docker.Code) (any, error)
}

type Service struct {
	Tasks
}

func NewService(image *docker.Image, repo *repository.Repository) *Service {
	return &Service{
		Tasks: NewTaskService(repo.Tasks, image.Container),
	}
}
