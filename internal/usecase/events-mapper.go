package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func EventsMapper(client *whatsmeow.Client, evt interface{}, redisClient *redis.Client) {

	if evt, ok := evt.(*events.Message); ok {

		result, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String(), "phone").Result()
		if result != "" {
			fmt.Println("Has a active conversation ")
			return

		}
		fields := map[string]interface{}{
			"lastMessage":   evt.Message.GetConversation(),
			"phone":         evt.Info.Chat.String(),
			"currentIdChat": "Initial",
		}

		err := redisClient.HMSet(context.Background(), evt.Info.Chat.String(), fields).Err()
		if err != nil {
			fmt.Println("Erro to save and init conversation:", err)
		}

		expirationDuration := 10 * time.Minute
		redisClient.Expire(context.Background(), evt.Info.Chat.String(), expirationDuration).Result()

		Init(client, evt, redisClient)

	}
}
