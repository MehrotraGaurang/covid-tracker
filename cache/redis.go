package cache

import (
	fileutil "app/fileUtil"
	"app/models"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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

	stateObjStr, err := json.Marshal(stateObj)

	if err != nil {
		panic(err)
	}

	err = getRdb().Set(ctx, stateObj.StateCode, stateObjStr, duration).Err()
	if err != nil {
		panic(err)
	}

}

func GetFromRedis(stateCode string) models.StateObject {

	stateObject := models.StateObject{}

	fmt.Println("Trying to get statecode", stateCode)

	stateObject.StateCode = "Not_Found"

	value, err := getRdb().Get(ctx, stateCode).Result()
	if err == redis.Nil {
		return stateObject
	} else if err != nil {
		panic(err)
	}

	valueBytes := []byte(value)

	json.Unmarshal(valueBytes, &stateObject)

	return stateObject

}

func StoreTs(ts int64) {
	fmt.Println("Caching timestamp", ts)

	err := getRdb().Set(ctx, "timestamp-final", ts, 0).Err()

	if err != nil {
		panic(err)
	}
}

func GetTs() int64 {

	fmt.Println("Getting last saved timestamp")

	value, err := getRdb().Get(ctx, "timestamp-final").Result()
	if err != nil {
		fmt.Println(err)
		return 0
	}

	fmt.Println("Timestamp from redis", value)

	timeStamp, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		panic(err)
	}

	fmt.Println("Cached TS is", timeStamp)

	return timeStamp
}
