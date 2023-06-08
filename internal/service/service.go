package service

import (
	"constester-go/internal/docker"
	"constester-go/internal/repository"
	"context"
)

type Tasks interface {
	AddTask(task repository.Task) error
	RunTests(ctx context.Context, code docker.Code) (docker.TestResult, error)
	GetAllTasks() (*[]repository.Task, error)
}

type Service struct {
	Tasks
}

func NewService(image *docker.Image, repo *repository.Repository) *Service {
	return &Service{
		Tasks: NewTaskService(repo.Tasks, image.Container),
	}
}
