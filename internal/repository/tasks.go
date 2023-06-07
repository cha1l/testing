package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Test struct {
	Name     int    `bson:"name"`
	Input    string `json:"input" bson:"input"`
	Expected string `json:"expected" bson:"expected"`
}

type Task struct {
	Name     string        `bson:"name" json:"name"`
	Text     string        `bson:"text" json:"text"`
	Duration time.Duration `bson:"duration" json:"duration"`
	Tests    []Test        `bson:"tests" json:"tests"`
}

type TaskRepository struct {
	coll *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) *TaskRepository {
	return &TaskRepository{
		coll: db.Collection("tasks"),
	}
}

func (t *TaskRepository) GetTests() ([]Test, time.Duration, error) {

	// todo : get tests from data base

	return []Test{
		{
			Name:     1,
			Input:    "1\n2 3",
			Expected: "6",
		},
		{
			Name:     2,
			Input:    "11\n12 5",
			Expected: "28",
		},
		{
			Name:     3,
			Input:    "6 7\n9",
			Expected: "22",
		},
	}, 1 * time.Second, nil
}

func (t *TaskRepository) GetText() {

	// todo : get task text from data base

}

func (t *TaskRepository) InsertTask(task Task) error {
	ctx := context.TODO()

	mod := mongo.IndexModel{
		Keys: bson.M{
			"name": 1,
		},
		Options: &options.IndexOptions{
			Unique: SetBool(true),
		},
	}
	_, err := t.coll.Indexes().CreateOne(ctx, mod)
	if err != nil {
		return err
	}

	_, err = t.coll.InsertOne(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

func SetBool(value bool) *bool {
	b := value
	return &b
}
