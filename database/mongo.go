package database

import (
	fileutil "app/fileUtil"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Db() *mongo.Client {
	clientOptions := options.Client().ApplyURI(fileutil.AppConfigProperties["url"])

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}
