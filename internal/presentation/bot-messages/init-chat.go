package presentation

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	gorm "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/gorm"
	infra "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/whatsmeow"

	usecase "github.com/savinnsk/prototype_bot_whatsapp/internal/usecase"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func Init(client *whatsmeow.Client, evt *events.Message, redisClient *redis.Client) {
	if evt.Message.GetConversation() == "1" {
		handleShowUserSchedules(client, evt, redisClient)
		return
	} else if evt.Message.GetConversation() == "2" {
		handlerShowSchedules(client, evt, redisClient)
	} else if evt.Message.GetConversation() == "3" {
		handlerWithScheduleCancel(client, evt, redisClient)

	} else if evt.Message.GetConversation() == "4" {
		handleWithSendContact(client, evt, redisClient)
	} else {
		handleDefaultConversation(client, evt, redisClient)
	}
}

func handleDefaultConversation(client *whatsmeow.Client, evt *events.Message, redisClient *redis.Client) {
	user, _ := gorm.FindUserByPhone(evt.Info.Chat.String())

	userPhone := evt.Info.Chat.String()

	currentChatId, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String(), "current.chat.id").Result()

	if user == nil && currentChatId == "" {
		infra.WhatsmeowSendResponse(client, evt, GetMessage().Greetings)
	}

	redisClient.HSet(context.Background(), userPhone, "phone", userPhone)

	redisClient.HSet(context.Background(), userPhone, "current.chat.id", "chat.init").Err()

	infra.WhatsmeowSendResponse(client, evt, GetMessage().MenuInteractionText)

}

func handleShowUserSchedules(client *whatsmeow.Client, evt *events.Message, redisClient *redis.Client) {

	schedules, err := usecase.LoadAllUserSchedules(evt.Info.Chat.String())
	print()
	if err != nil {
		infra.WhatsmeowSendResponse(client, evt, GetMessage().SchedulesNotFound)
		time.AfterFunc(1*time.Second, func() {
			handleDefaultConversation(client, evt, redisClient)
		})
		return
	}

	msg := "*Seus Agendamentos Abaixo:* \n"

	for _, schedule := range schedules {
		msg += fmt.Sprintf("\n🕥 - *HORA*: %s , *DATA* : %s", schedule.Time, schedule.Date)
	}

	infra.WhatsmeowSendResponse(client, evt, msg)
	time.AfterFunc(5*time.Second, func() {
		handleDefaultConversation(client, evt, redisClient)
	})

}

// [DEAL] with error type at new_schedule
func handlerShowSchedules(client *whatsmeow.Client, evt *events.Message, redisClient *redis.Client) {

	schedulesFiltered := usecase.FilterSchedules()

	if len(schedulesFiltered) == 0 {
		msg := "*🙁 - Não há agendamentos para hoje.*"
		infra.WhatsmeowSendResponse(client, evt, msg)
		return
	}

	schedulesJSON, _ := json.Marshal(schedulesFiltered)

	redisClient.HSet(context.Background(), evt.Info.Chat.String(), "schedules", schedulesJSON).Result()
	msg := GetMessage().SchedulesAvailableTitle

	for i, schedule := range schedulesFiltered {
		msg += fmt.Sprintf("\n%d - *HORA*: %s Hoje 🕥", i+2, schedule)
	}

	msg += GetMessage().ScheduleOtherTime + GetMessage().BackButton + GetMessage().DefaultFooter

	redisClient.HSet(context.Background(), evt.Info.Chat.String(), "current.chat.id", "chat.show.schedules").Result()
	infra.WhatsmeowSendResponse(client, evt, msg)

}

// [DEAL] with error type at cancel_schedule
func handlerWithScheduleCancel(client *whatsmeow.Client, evt *events.Message, redisClient *redis.Client) {
	schedules, err := usecase.LoadAllUserSchedules(evt.Info.Chat.String())
	if err != nil {
		infra.WhatsmeowSendResponse(client, evt, GetMessage().SchedulesNotFound)
		time.AfterFunc(1*time.Second, func() {
			handleDefaultConversation(client, evt, redisClient)
		})
		return
	}

	var userScheduleArray []string

	for _, schedule := range schedules {
		userScheduleArray = append(userScheduleArray, schedule.Time)

	}

	schedulesJSON, _ := json.Marshal(userScheduleArray)
	redisClient.HSet(context.Background(), evt.Info.Chat.String(), "schedules.to.cancel", schedulesJSON).Result()
	msg := "*Seus Agendamentos Abaixo:* \n"

	for i, schedule := range schedules {
		msg += fmt.Sprintf("\n%d - *HORA*: %s , *DATA* : %s", i+2, schedule.Time, schedule.Date)
	}

	msg += GetMessage().BackButton
	msg += GetMessage().DefaultFooter

	redisClient.HSet(context.Background(), evt.Info.Chat.String(), "current.chat.id", "chat.cancel.schedule").Result()
	infra.WhatsmeowSendResponse(client, evt, msg)

}

func handleWithSendContact(client *whatsmeow.Client, evt *events.Message, redisClient *redis.Client) {

	infra.WhatsmeowSendResponse(client, evt, GetMessage().Contacts)
	time.AfterFunc(5*time.Second, func() {
		handleDefaultConversation(client, evt, redisClient)
	})

}
