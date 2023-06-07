package repository

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

func ConnectToDB() (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info("connecting to mongo db ...")

	password := os.Getenv("MONGO_PASSWORD")
	username := os.Getenv("MONGO_USER")
	dbName := os.Getenv("MONGO_DB_NAME")

	connectURL := fmt.Sprintf("mongodb+srv://%s:%s@constester.z5qopzs.mongodb.net/?retryWrites=true&w=majority", username, password)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(connectURL).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	if err := client.Database("admin").RunCommand(ctx, bson.D{{"ping", 1}}).Err(); err != nil {
		return nil, err
	}

	log.Info("connected to mongo")

	db := client.Database(dbName)

	return db, nil
}
