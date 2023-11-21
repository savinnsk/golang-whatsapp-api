package infra

import (
	"context"

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

func DeleteValues(ctx context.Context, index string, values []string) error {

	for _, key := range values {
		err := redisClient.HDel(ctx, index, key).Err()
		if err != nil {
			return err
		}
	}

	return nil

}

func AddValues(ctx context.Context, index string, keysValues []string, values []string) error {
	for _, keyValue := range keysValues {
		for _, value := range values {
			err := redisClient.HSet(ctx, index, keyValue, value).Err()
			if err != nil {
				return err
			}
		}

		return nil
	}
	return nil
}

func GetValues(ctx context.Context, index string, keys []string) ([]string, error) {

	var values []string
	for _, key := range keys {
		value, err := redisClient.HGet(ctx, index, key).Result()
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	return values, nil

}
