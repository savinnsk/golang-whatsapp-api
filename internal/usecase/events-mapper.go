package usecase

import (
	"context"

	"github.com/go-redis/redis/v8"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func EventsMapper(client *whatsmeow.Client, evt interface{}, redisClient *redis.Client) {

	if evt, ok := evt.(*events.Message); ok {

		currentChatId, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String(), "currentChatId").Result()

		if currentChatId == "SHOW_USER_SCHEDULE" && evt.Message.GetConversation() == "0" {
			redisClient.HSet(context.Background(), evt.Info.Chat.String(), "currentChatId", "INIT").Result()
			Init(client, evt, redisClient, currentChatId)
			return

		}

		if currentChatId == "NEW_SCHEDULE" && evt.Message.GetConversation() == "0" {
			redisClient.HSet(context.Background(), evt.Info.Chat.String(), "currentChatId", "INIT").Result()
			Init(client, evt, redisClient, currentChatId)
			return
		}

		if currentChatId == "NEW_SCHEDULE" {

			NewSchedule(client, evt, redisClient, currentChatId)
			return
		}

		if currentChatId == "CANCEL_SCHEDULE" && evt.Message.GetConversation() == "0" {
			redisClient.HSet(context.Background(), evt.Info.Chat.String(), "currentChatId", "INIT").Result()
			Init(client, evt, redisClient, currentChatId)
			return
		}

		if currentChatId == "INIT" || currentChatId == "" {
			Init(client, evt, redisClient, currentChatId)
		}

	}
}
