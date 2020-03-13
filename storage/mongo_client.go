package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	client *mongo.Client
	Db     *mongo.Database
}

func NewMongoClient(connection string) (*MongoClient, error) {
	opts := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	err = client.Connect(context.TODO())
	if err != nil {
		return nil, err
	}
	db := client.Database(connection)
	return &MongoClient{client: client, Db: db}, nil
}

func (client *MongoClient) Disconnect() error {
	return client.client.Disconnect(context.TODO())
}
