package presentation

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	t "time"

	"github.com/go-redis/redis/v8"
	"github.com/savinnsk/prototype_bot_whatsapp/internal/configs"
	dto "github.com/savinnsk/prototype_bot_whatsapp/internal/domain/dto"
	gorm "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/gorm"
	redisC "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/redis"
	infra "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/whatsmeow"
	usecase "github.com/savinnsk/prototype_bot_whatsapp/internal/usecase"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func NewSchedule(cs ChatSetup) {
	clientIndex := cs.evt.Info.Chat.String()
	messageFromClient := cs.evt.Message.GetConversation()
	stageConversation, err := redisC.GetValue(cs.redisClient, context.Background(), clientIndex, "chat.stage") //cs.redisClient.HGet(context.Background(), cs.evt.Info.Chat.String(), "chat.stage").Result()
	if err != nil {
		configs.MakeError(err, "schedules.message.go / line : 24 ")
	}

	_, MessageNotNumeric := strconv.Atoi(messageFromClient)

	if messageFromClient == "1" && stageConversation == "" {
		SetOption1ScheduleAnotherDate(cs)
		return
	}

	if stageConversation == "schedule.another.date" {
		HandleScheduleAnotherDate(cs)
		return

	}

	if stageConversation == "chat.time" && messageFromClient != "0" {

		schedulesJson, _ := cs.redisClient.HGet(context.Background(), cs.evt.Info.Chat.String(), "schedules").Result()
		var schedules []string
		if err := json.Unmarshal([]byte(schedulesJson), &schedules); err != nil {
			fmt.Print("\nError unmarshaling JSON:", err, "\n")
			return
		}
		timeChose, err := usecase.VerifyScheduleBasedAtArray(cs.evt.Message.GetConversation(), schedules)

		if MessageNotNumeric != nil || err != nil || timeChose == "" {
			infra.WhatsmeowSendResponse(cs.client, cs.evt, GetMessage().VerifyValues)
			return
		}
		cs.redisClient.HSet(context.Background(), cs.evt.Info.Chat.String(), "chat.time.chose", timeChose)
		phone := cs.evt.Info.Chat.String()
		user, _ := gorm.FindUserByPhone(phone)

		if user == nil {
			msg := `Digite seu *Nome* e *Sobrenome*:.*`
			cs.redisClient.HSet(context.Background(), cs.evt.Info.Chat.String(), "chat.stage", "chat.get.name").Result()
			infra.WhatsmeowSendResponse(cs.client, cs.evt, msg)
			return
		}

		msg := "*Seu nome é :*" + " " + user.Name
		msg += "\n\n_1 - SIM_"
		msg += "\n_0 - NÂO_"
		cs.redisClient.HSet(context.Background(), cs.evt.Info.Chat.String(), "chat.stage", "chat.get.name").Result()
		infra.WhatsmeowSendResponse(cs.client, cs.evt, msg)
		return
	}

	if stageConversation == "chat.get.name" {
		date, _ := cs.redisClient.HGet(context.Background(), cs.evt.Info.Chat.String(), "target.date").Result()
		timeChose, _ := cs.redisClient.HGet(context.Background(), cs.evt.Info.Chat.String(), "chat.time.chose").Result()
		targetDate, _ := t.Parse("2006-01-02T15:04:05Z", date)
		targetDateFormatted := targetDate.Format("02/01/2006")

		user, err := gorm.FindUserByPhone(cs.evt.Info.Chat.String())

		if err != nil && user != nil {
			result := SaveScheduleWithAUserAlreadySignUp(user.Name, user.Phone, targetDateFormatted, timeChose)

			if result != "ok" {
				SendBotErrorDefault("Error to save user with account", cs.client, cs.evt)

			}
		}

		result := SaveScheduleWithAUserAlreadySignUp(cs.evt.Message.GetConversation(), cs.evt.Info.Chat.String(), targetDateFormatted, timeChose)
		if result != "ok" {
			SendBotErrorDefault("Error to save user with no account", cs.client, cs.evt)
		}
		//
		msg := fmt.Sprintf("_seu agendamento foi concluído_ para : *%s* 🕥", targetDate.Format("02/01/2006"))
		cs.redisClient.HSet(context.Background(), cs.evt.Info.Chat.String(), "chat.stage", "finish").Result()
		infra.WhatsmeowSendResponse(cs.client, cs.evt, msg)
		return
	}

	if cs.evt.Message.GetConversation() == "0" {
		BackToMenu(cs)
	}

}

func SetOption1ScheduleAnotherDate(cs ChatSetup) {
	err := redisC.AddValues(cs.redisClient, context.Background(), cs.evt.Info.Chat.String(), []string{"chat.stage"}, []any{"schedule.another.date"})
	if err != nil {
		configs.MakeError(err, "schedules.message.go / line : 110")
	}
	infra.WhatsmeowSendResponse(cs.client, cs.evt, GetMessage().ScheduleAnotherDate)
}

func HandleScheduleAnotherDate(cs ChatSetup) {
	messageFromClient := cs.evt.Message.GetConversation()

	if messageFromClient == "0" {
		redisC.DeleteValues(cs.redisClient, context.Background(), cs.evt.Info.Chat.String(), []string{"chat.stage", "schedules"})
		HandlerShowSchedules(cs)
		return
	}

	targetDate, err := t.Parse("02/01/2006", cs.evt.Message.GetConversation())
	if err != nil {
		configs.MakeError(err, "schedules.message.go / line : 126")
		infra.WhatsmeowSendResponse(cs.client, cs.evt, GetMessage().NotUnderstand)
		return
	}

	if targetDate.Before(t.Now()) {
		configs.MakeError(err, "schedules.message.go / line : 126")
		infra.WhatsmeowSendResponse(cs.client, cs.evt, GetMessage().VerifyValues)
		return
	}

	schedulesFiltered := usecase.FilterSchedules()
	schedulesJSON, err := json.Marshal(schedulesFiltered)

	if err != nil {
		configs.MakeError(err, "schedules.message.go / line : 139")
	}

	err = redisC.AddValues(cs.redisClient, context.Background(), cs.evt.Info.Chat.String(), []string{"target.date", "schedules", "chat.stage"}, []any{targetDate, schedulesJSON, "chat.time"})
	if err != nil {
		configs.MakeError(err, "schedules.message.go / line : 146")
	}

	msg := `Horários Disponíveis para :*` + targetDate.Format("02/01/2006") + "* Abaixo:\n\n"

	for i, schedule := range schedulesFiltered {
		msg += fmt.Sprintf("\n%d - *HORA*: %s  *%s* 🕥", i+2, schedule, targetDate.Format("02/01/2006"))
	}

	msg += GetMessage().BackButton + GetMessage().DefaultFooter

	infra.WhatsmeowSendResponse(cs.client, cs.evt, msg)
	return
}

// refactor

func ScheduleAnotherTimeWithUserAlreadyRegistered(client *whatsmeow.Client, evt *events.Message, redisClient *redis.Client) {
	user, _ := gorm.FindUserByPhone(evt.Info.Chat.String())
	redisClient.HSet(context.Background(), evt.Info.Chat.String(), "name.user", user.Name).Result()
	msg := `*Confirmar Agendamento ?*`
	msg += "\n\n_1 - SIM_"
	msg += "\n_0 - NÂO_"
	redisClient.HSet(context.Background(), evt.Info.Chat.String(), "currentChatId", "chat.save").Result()
	infra.WhatsmeowSendResponse(client, evt, msg)

}

func SaveScheduleWithAUserAlreadySignUp(nameClient string, phone string, targetDate string, targetTime string) string {

	userData := dto.CreateUserDto{
		Name:  nameClient,
		Phone: phone,
		Role:  "user",
	}

	data := dto.SaveNewUserAndSchedule{
		CreateUserDto: userData,
		ScheduleTime:  targetTime,
		ScheduleDate:  targetDate,
	}

	result := usecase.ProcessNewSchedule(data)

	return result

}

func BackToMenu(cs ChatSetup) {
	Init(cs)

}
