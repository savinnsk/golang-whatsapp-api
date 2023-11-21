package presentation

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	gorm "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/gorm"
	infra "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/whatsmeow"

	usecase "github.com/savinnsk/prototype_bot_whatsapp/internal/usecase"
)

func Init(cs ChatSetup) {

	messageFromClient := cs.evt.Message.GetConversation()

	if messageFromClient == "1" {
		HandleShowUserSchedules(cs)
		return
	} else if messageFromClient == "2" {
		HandlerShowSchedules(cs)
		return
	} else if messageFromClient == "3" {
		HandlerWithScheduleCancel(cs)
		return
	} else if messageFromClient == "4" {
		HandleWithSendContact(cs)
		return
	} else {
		HandleDefaultConversation(cs)
		return
	}
}

func HandleDefaultConversation(cs ChatSetup) {
	user, _ := gorm.FindUserByPhone(cs.evt.Info.Chat.String())

	userPhone := cs.evt.Info.Chat.String()

	currentChatId, _ := cs.redisClient.HGet(context.Background(), cs.evt.Info.Chat.String(), "current.chat.id").Result()

	if user == nil && currentChatId == "" {
		infra.WhatsmeowSendResponse(cs.client, cs.evt, GetMessage().Greetings)
	}

	cs.redisClient.HSet(context.Background(), userPhone, "phone", userPhone)

	cs.redisClient.HSet(context.Background(), userPhone, "current.chat.id", "chat.init").Err()

	infra.WhatsmeowSendResponse(cs.client, cs.evt, GetMessage().MenuInteractionText)

}

func HandleShowUserSchedules(cs ChatSetup) {

	schedules, err := usecase.LoadAllUserSchedules(cs.evt.Info.Chat.String())
	print()
	if err != nil {
		infra.WhatsmeowSendResponse(cs.client, cs.evt, GetMessage().SchedulesNotFound)
		time.AfterFunc(1*time.Second, func() {
			HandleDefaultConversation(cs)
		})
		return
	}

	msg := "*Seus Agendamentos Abaixo:* \n"

	for _, schedule := range schedules {
		msg += fmt.Sprintf("\n🕥 - *HORA*: %s , *DATA* : %s", schedule.Time, schedule.Date)
	}

	infra.WhatsmeowSendResponse(cs.client, cs.evt, msg)
	time.AfterFunc(5*time.Second, func() {
		HandleDefaultConversation(cs)
	})

}

func HandlerShowSchedules(cs ChatSetup) {

	schedulesFiltered := usecase.FilterSchedules()

	if len(schedulesFiltered) == 0 {
		msg := "*🙁 - Não há agendamentos para hoje.*"
		infra.WhatsmeowSendResponse(cs.client, cs.evt, msg)
		return
	}

	schedulesJSON, _ := json.Marshal(schedulesFiltered)

	cs.redisClient.HSet(context.Background(), cs.evt.Info.Chat.String(), "schedules", schedulesJSON).Result()
	msg := GetMessage().SchedulesAvailableTitle

	for i, schedule := range schedulesFiltered {
		msg += fmt.Sprintf("\n%d - *HORA*: %s Hoje 🕥", i+2, schedule)
	}

	msg += GetMessage().ScheduleOtherTime + GetMessage().BackButton + GetMessage().DefaultFooter

	cs.redisClient.HSet(context.Background(), cs.evt.Info.Chat.String(), "current.chat.id", "chat.show.schedules").Result()
	infra.WhatsmeowSendResponse(cs.client, cs.evt, msg)

}

// [DEAL] with error type at cancel_schedule
func HandlerWithScheduleCancel(cs ChatSetup) {
	schedules, err := usecase.LoadAllUserSchedules(cs.evt.Info.Chat.String())
	if err != nil {
		infra.WhatsmeowSendResponse(cs.client, cs.evt, GetMessage().SchedulesNotFound)
		time.AfterFunc(1*time.Second, func() {
			HandleDefaultConversation(cs)
		})
		return
	}

	var userScheduleArray []string

	for _, schedule := range schedules {
		userScheduleArray = append(userScheduleArray, schedule.Time)

	}

	schedulesJSON, _ := json.Marshal(userScheduleArray)
	cs.redisClient.HSet(context.Background(), cs.evt.Info.Chat.String(), "schedules.to.cancel", schedulesJSON).Result()
	msg := "*Seus Agendamentos Abaixo:* \n"

	for i, schedule := range schedules {
		msg += fmt.Sprintf("\n%d - *HORA*: %s , *DATA* : %s", i+2, schedule.Time, schedule.Date)
	}

	msg += GetMessage().BackButton
	msg += GetMessage().DefaultFooter

	cs.redisClient.HSet(context.Background(), cs.evt.Info.Chat.String(), "current.chat.id", "chat.cancel.schedule").Result()
	infra.WhatsmeowSendResponse(cs.client, cs.evt, msg)

}

func HandleWithSendContact(cs ChatSetup) {

	infra.WhatsmeowSendResponse(cs.client, cs.evt, GetMessage().Contacts)
	time.AfterFunc(5*time.Second, func() {
		HandleDefaultConversation(cs)
	})

}
