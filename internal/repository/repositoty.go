package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Tasks interface {
	GetTests(taskName string) ([]Test, time.Duration, error)
	GetText()
	InsertTask(task Task) error
	GetAllTasks() (*[]Task, error)
}

type Repository struct {
	Tasks
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Tasks: NewTaskRepository(db),
	}
}
