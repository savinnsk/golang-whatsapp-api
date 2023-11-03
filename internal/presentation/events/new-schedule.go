package presentation

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	dto "github.com/savinnsk/prototype_bot_whatsapp/internal/domain/dto"
	infra "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/whatsmeow"
	usecase "github.com/savinnsk/prototype_bot_whatsapp/internal/usecase"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func NewSchedule(client *whatsmeow.Client, evt *events.Message, redisClient *redis.Client, currentChatId string) {

	stage, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String(), "chatStage").Result()

	if stage == "SAVE" && evt.Message.GetConversation() != "0" {
		time, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "time").Result()
		date, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "date").Result()
		name, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "name").Result()

		userData := dto.CreateUserDto{
			Name:  name,
			Phone: evt.Info.Chat.String(),
			Role:  "user",
		}

		data := dto.SaveNewUserAndSchedule{
			CreateUserDto: userData,
			ScheduleTime:  time,
			ScheduleDate:  date,
		}

		result := usecase.ProcessNewSchedule(data)
		if result != "ok" {
			redisClient.HDel(context.Background(), evt.Info.Chat.String(), "chatStage")
			redisClient.HDel(context.Background(), evt.Info.Chat.String(), "currentChatId")
			infra.WhatsmeowSendResponse(client, evt, "\n\n_"+result+"_")
			return
		}

		msg := `_😀 SEU AGENDAMENTO FOI SALVO COM SUCESSO_`

		infra.WhatsmeowSendResponse(client, evt, msg)
		redisClient.HDel(context.Background(), evt.Info.Chat.String(), "chatStage")
		redisClient.HDel(context.Background(), evt.Info.Chat.String(), "currentChatId")
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

		msg := `*Qual seu nome?*`
		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chatStage", "NAME").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)
		return
	}

	if stage == "DATA" && evt.Message.GetConversation() != "0" {
		redisClient.HSet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "date", evt.Message.GetConversation()).Result()
		date, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "date").Result()
		schedules := usecase.LoadAllValidSchedulesDates()
		schedulesJSON, _ := json.Marshal(schedules)

		redisClient.HSet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "schedules", schedulesJSON).Result()
		msg := `Horários Disponíveis para :*` + date + `* Abaixo:`

		for i, schedule := range schedules {
			if schedule.Available && !schedule.Disabled {
				msg += fmt.Sprintf("\n\n%d - *HORA*: %s  *%s* 🕥", i+2, schedule.Time, date)
			}
		}

		msg += "\n\n_0 - VOLTAR  ◀️_"
		msg += "\n\n_Responda com o número correspondente à sua escolha. Para agendar_ 📅"

		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chatStage", "TIME").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)
		return
	}

	if evt.Message.GetConversation() == "1" {
		msg := `*Digite qual data você quer agendar:*

*Nesse formato separado por barras* 

*exemplo:* 

25/06/2023 
*dia/mes/ano*
	  
_0 - VOLTAR ? ◀️_`

		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chatStage", "DATA").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)
		return

	}

}
