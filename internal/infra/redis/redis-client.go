package infra

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func Init() *redis.Client {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Endereço do servidor Redis
		Password: "",               // Senha, se necessário
		DB:       0,                // Número do banco de dados Redis
	})

	// Teste a conexão com o Redis
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	return redisClient

}
