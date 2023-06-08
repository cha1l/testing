package service

import (
	"constester-go/internal/docker"
	"constester-go/internal/repository"
	"context"
	"errors"
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

func (t *TaskService) RunTests(ctx context.Context, code docker.Code) (docker.TestResult, error) {
	tests, dur, err := t.repo.GetTests(code.TaskName)
	if err != nil {
		return docker.TestResult{}, err
	}

	var res docker.TestResult

	switch code.Language {
	case "cpp":
		res, err = t.container.RunTestsCPP(ctx, tests, dur, []byte(code.Code))
	default:
		return docker.TestResult{}, errors.New("invalid or unsupported language")
	}

	return res, err
}

func (t *TaskService) GetAllTasks() (*[]repository.Task, error) {
	return t.repo.GetAllTasks()
}
