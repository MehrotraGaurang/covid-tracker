package cache

import (
	"app/constants"
	fileutil "app/fileUtil"
	"app/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.TODO()

func getRdb() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     fileutil.AppConfigProperties["redis_url"],
		Password: "HsTF1TISXxtfKJNLCPNjm5xbKx3JJ2RF",
		DB:       0,
	})

	fmt.Println("Connected to redis, now pinging")

	_, err := rdb.Ping(ctx).Result()

	if err != nil {
		fmt.Println("Not able to connect to Redis!")
		panic(err)
	}

	fmt.Println("Able to ping redis")

	return rdb
}

func StoreInRedis(stateObj models.StateObject, duration time.Duration) {

	fmt.Println("Caching statecode for 30min", stateObj.StateCode)

	rdb := getRdb()

	stateObjStr, err := json.Marshal(stateObj)

	if err != nil {
		panic(err)
	}

	err = rdb.Set(ctx, stateObj.StateCode, stateObjStr, duration).Err()
	if err != nil {
		panic(err)
	}

}

func GetFromRedis(stateCode string) models.StateObject {

	stateObject := models.StateObject{}

	fmt.Println("Trying to get statecode", stateCode)

	rdb := getRdb()

	stateObject.StateCode = "Not_Found"

	value, err := rdb.Get(ctx, stateCode).Result()
	if err == redis.Nil {
		return stateObject
	} else if err != nil {
		panic(err)
	}

	valueBytes := []byte(value)

	json.Unmarshal(valueBytes, &stateObject)

	return stateObject

}

func StoreTs(ts time.Time) {
	fmt.Println("Caching timestamp", ts)

	rdb := getRdb()

	err := rdb.Set(ctx, "timestamp", ts, 0).Err()

	if err != nil {
		panic(err)
	}
}

func GetTs() time.Time {

	fmt.Println("Getting last saved timestamp")

	rdb := getRdb()

	value, err := rdb.Get(ctx, "timestamp").Result()
	if err != nil {
		panic(err)
	}

	timeStamp, error := time.Parse(constants.Layout, value)
	if error != nil {
		panic(err)
	}

	fmt.Print("Cached TS is", timeStamp)

	return timeStamp
}
