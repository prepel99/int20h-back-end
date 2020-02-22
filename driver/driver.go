package driver

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"int20h-back-end/logger"
	"time"
)

var logr logger.Logger

func ConnectDB(url string) *mongo.Client {
	for {
		client, err := connectDB(url)
		if err != nil {
			logr.LogErr(err)
			fmt.Println("Cannot connect to Mongo. Next try in 5 sec:")
			time.Sleep(5 * time.Second)
		} else {
			return client
		}
	}

}
func connectDB(url string) (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to MongoDB!")
	return client, nil
}
