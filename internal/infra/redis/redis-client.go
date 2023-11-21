package infra

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func Init() *redis.Client {

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	return redisClient

}

func DeleteValues(rc *redis.Client, ctx context.Context, index string, values []string) error {

	for _, key := range values {
		err := rc.HDel(ctx, index, key).Err()
		if err != nil {
			return err
		}
	}

	return nil

}

func AddValues(rc *redis.Client, ctx context.Context, index string, keysValues []string, values []any) error {

	for i, value := range values {
		fmt.Print("aaaaaaaaaaaaaaaaaaaaaaaaaaa  \n", keysValues[i], value)
		err := redisClient.HSet(ctx, index, keysValues[i], value).Err()
		if err != nil {
			return err
		}
	}

	return nil

}

func GetValue(rc *redis.Client, ctx context.Context, index string, key string) (string, error) {
	value, err := rc.HGet(ctx, index, key).Result()
	if err != nil {
		return "", err
	}

	return value, nil

}
