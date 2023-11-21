package presentation

import (
	"context"

	"github.com/go-redis/redis/v8"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func EventsMapper(client *whatsmeow.Client, evt interface{}, redisClient *redis.Client) {

	if evt, ok := evt.(*events.Message); ok {
		currentChatId, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String(), "current.chat.id").Result()
		chatSetup := &ChatSetup{
			client:        client,
			evt:           evt,
			redisClient:   redisClient,
			currentChatId: currentChatId,
		}

		if currentChatId == "chat.show.schedules" {
			NewSchedule(*chatSetup)
			return
		}

		if currentChatId == "CANCEL_SCHEDULE" && evt.Message.GetConversation() == "0" {
			redisClient.HSet(context.Background(), evt.Info.Chat.String(), "currentChatId", "INIT").Result()
			Init(*chatSetup)
			return
		}
		if currentChatId == "CANCEL_SCHEDULE" {
			Cancel(client, evt, redisClient)
		}

		if currentChatId != "chat.back" {
			Init(*chatSetup)
		}

	}
}
