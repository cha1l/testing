package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Tasks interface {
	GetTests() ([]Test, time.Duration, error)
	GetText()
	InsertTask(task Task) error
}

type Repository struct {
	Tasks
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Tasks: NewTaskRepository(db),
	}
}
