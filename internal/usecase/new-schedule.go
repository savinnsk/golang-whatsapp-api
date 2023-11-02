package usecase

import (
	"context"

	//	"fmt"
	//"time"

	"github.com/go-redis/redis/v8"
	//"gorm.io/gorm"
	entity "github.com/savinnsk/prototype_bot_whatsapp/internal/entity"
	gorm "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/gorm"
	infra "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/whatsmeow"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func NewSchedule(client *whatsmeow.Client, evt *events.Message, redisClient *redis.Client, currentChatId string) {

	stage, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String(), "chatStage").Result()

	if stage == "SAVE" && evt.Message.GetConversation() != "0" {
		time, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "time").Result()
		date, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "date").Result()
		name, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "name").Result()

		schedule, err := gorm.FindScheduleByTime(time)
		if err != nil {
			infra.WhatsmeowSendResponse(client, evt, "ERRO")
			return
		}

		user := entity.NewUser(name, evt.Info.Chat.String(), "user")

		gorm.CreateUser(user)
		newUser, err := gorm.FindUserByPhone(evt.Info.Chat.String())
		if err != nil {
			infra.WhatsmeowSendResponse(client, evt, "ERRO")
			return
		}
		println(">>>>>>>>>>>>>>>>>>>", date)
		gorm.CreateUserSchedule(newUser.Id, schedule.Id, date, time)

		msg := `_SEU AGENDAMENTO FOI SALVO COM SUCESSO_`
		//redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chatStage", "NAME").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)
		return
	}

	if stage == "NAME" && evt.Message.GetConversation() != "0" {
		redisClient.HSet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "name", evt.Message.GetConversation()).Result()
		msg := `*Confirmar ?*`
		msg += "\n\n_1 - SIM_"
		msg += "\n_0 - NÂO_"
		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chatStage", "SAVE").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)
		return
	}

	if stage == "TIME" && evt.Message.GetConversation() != "0" {
		redisClient.HSet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "time", evt.Message.GetConversation()).Result()
		// time, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "time").Result()
		// date, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "date").Result()

		msg := `*Qual seu nome?*`
		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chatStage", "NAME").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)
		return
	}

	if stage == "DATA" && evt.Message.GetConversation() != "0" {
		redisClient.HSet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "date", evt.Message.GetConversation()).Result()

		msg := `*Digite qual horário você deseja*

		- *Nesse formato separado por dois pontos ex:*  13:30
	   
	   
		`
		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chatStage", "TIME").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)
		return
	}

	if evt.Message.GetConversation() == "1" {
		msg := `*Digite qual data você quer agendar:*

*Nesse formato separado por barras* 

ex : 25/06/2023 
	 *dia/mes/ano*

	  
_0 - VOLTAR ? ◀️_`

		//redisClient.HSet(context.Background(), evt.Info.Chat.String(), "currentChatId", "NEW_SCHEDULE").Result()
		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chatStage", "DATA").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)
		return

	}

}
