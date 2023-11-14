package presentation

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	gorm "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/gorm"
	infra "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/whatsmeow"
	usecase "github.com/savinnsk/prototype_bot_whatsapp/internal/usecase"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func Cancel(client *whatsmeow.Client, evt *events.Message, redisClient *redis.Client) {

	schedulesJson, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String()+"CANCEL_SCHEDULE", "schedules").Result()
	var schedules []string
	if err := json.Unmarshal([]byte(schedulesJson), &schedules); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
	timeChose := usecase.VerifyScheduleBasedAtArray(evt.Message.GetConversation(), schedules)
	gorm.DeleUserScheduleByTime(timeChose)

	msg := "\n\n_Horário Excluído com sucesso ❌*_ 📅"

	redisClient.HDel(context.Background(), evt.Info.Chat.String(), "currentChatId").Result()
	redisClient.Del(context.Background(), evt.Info.Chat.String()+"CANCEL_SCHEDULE")
	infra.WhatsmeowSendResponse(client, evt, msg)

}
