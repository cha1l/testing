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

func (t *TaskRepository) GetTests(taskName string) ([]Test, time.Duration, error) {
	ctx := context.TODO()

	var (
		opt  options.FindOneOptions
		task Task
	)

	filter := bson.M{"name": taskName}
	opt.SetProjection(bson.M{"tests": 1, "duration": 1})

	if err := t.coll.FindOne(ctx, filter, &opt).Decode(&task); err != nil {
		return nil, 0, err
	}

	return task.Tests, task.Duration, nil
}

func (t *TaskRepository) GetText() {

	// todo : get task text from data base

}

func (t *TaskRepository) InsertTask(task Task) error {
	ctx := context.Background()

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

func (t *TaskRepository) GetAllTasks() (*[]Task, error) {
	cxt := context.Background()

	var (
		tasks []Task
		opts  options.FindOptions
	)

	opts.SetProjection(bson.M{"name": 1, "text": 1, "duration": 1})

	cur, err := t.coll.Find(cxt, bson.D{{}}, &opts)
	if err != nil {
		return nil, err
	}

	for cur.Next(cxt) {
		var task Task

		if err := cur.Decode(&task); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	err = cur.Err()

	return &tasks, err
}

func SetBool(value bool) *bool {
	b := value
	return &b
}
