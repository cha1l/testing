package service

import (
	"constester-go/internal/docker"
	"constester-go/internal/repository"
	"context"
)

type TaskService struct {
	repo      repository.Tasks
	container docker.Container
}

func NewTaskService(repo repository.Tasks, container docker.Container) *TaskService {
	return &TaskService{
		repo:      repo,
		container: container,
	}
}

func (t *TaskService) AddTask(task repository.Task) error {
	for i := range task.Tests {
		task.Tests[i].Name = i + 1
	}

	return t.repo.InsertTask(task)
}

func (t *TaskService) RunTestsCPP(ctx context.Context) error {
	return nil
}
