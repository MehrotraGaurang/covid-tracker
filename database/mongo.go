package database

import (
	fileutil "app/fileUtil"
	"app/models"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
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

func UpdateDb(items []models.StateObject) {

	var userCollection = Db().Database(fileutil.AppConfigProperties["database"]).Collection(fileutil.AppConfigProperties["collection"])

	for _, item := range items {

		_, err := userCollection.InsertOne(context.TODO(), item)

		if err != nil {
			log.Fatal(err)
		}

	}

}

func GetCovidData(stateCode string) models.StateObject {

	fmt.Println("Starting to get data now for ", stateCode)

	var userCollection = Db().Database(fileutil.AppConfigProperties["database"]).Collection(fileutil.AppConfigProperties["collection"])

	filterCursor, err := userCollection.Find(context.TODO(), bson.M{"statecode": stateCode})

	fmt.Println("Data is here")

	if err != nil {
		log.Fatal(err)
	}

	var stateInfoBlocks []models.StateObject

	fmt.Println("Initializing StateObject list")

	if err = filterCursor.All(context.TODO(), &stateInfoBlocks); err != nil {
		log.Fatal(err)
	}

	if len(stateInfoBlocks) == 0 {
		var stateObj models.StateObject
		return stateObj
	}

	fmt.Println("Now filtering it!")

	var maxLastTime = stateInfoBlocks[0]

	for i := 1; i < len(stateInfoBlocks); i++ {

		item := stateInfoBlocks[i]

		if maxLastTime.LastUpdated.Before(item.LastUpdated) {
			maxLastTime = item
		}

	}

	return maxLastTime
}
