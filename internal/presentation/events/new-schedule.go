package presentation

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	t "time"

	"github.com/go-redis/redis/v8"
	dto "github.com/savinnsk/prototype_bot_whatsapp/internal/domain/dto"
	gorm "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/gorm"
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

		if date == "" {
			currentDate := t.Now().Format("02/01/2006")
			date = currentDate
		}
		data := dto.SaveNewUserAndSchedule{
			CreateUserDto: userData,
			ScheduleTime:  time,
			ScheduleDate:  date,
		}

		result := usecase.ProcessNewSchedule(data)
		if result != "ok" {
			redisClient.Del(context.Background(), evt.Info.Chat.String())
			redisClient.Del(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE")
			infra.WhatsmeowSendResponse(client, evt, "\n\n_"+result+"_")
			return
		}

		msg := `_😀 SEU AGENDAMENTO FOI SALVO COM SUCESSO_`

		redisClient.Del(context.Background(), evt.Info.Chat.String())
		redisClient.Del(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE")
		infra.WhatsmeowSendResponse(client, evt, msg)

		return
	}

	if stage == "NAME" && evt.Message.GetConversation() == "1" {
		user, _ := gorm.FindUserByPhone(evt.Info.Chat.String())
		redisClient.HSet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "name", user.Name).Result()
		msg := `*Confirmar Agendamento ?*`
		msg += "\n\n_1 - SIM_"
		msg += "\n_0 - NÂO_"
		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chatStage", "SAVE").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)
		return
	}

	if stage == "NAME" && evt.Message.GetConversation() != "0" {
		redisClient.HSet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "name", evt.Message.GetConversation()).Result()
		msg := `*Confirmar Agendamento ?*`
		msg += "\n\n_1 - SIM_"
		msg += "\n_0 - NÂO_"
		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chatStage", "SAVE").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)
		return
	}

	if stage == "TIME" && evt.Message.GetConversation() != "0" {
		schedulesJson, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "schedules").Result()
		var schedules []string
		if err := json.Unmarshal([]byte(schedulesJson), &schedules); err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return
		}
		timeChose := usecase.VerifyScheduleBasedAtArray(evt.Message.GetConversation(), schedules)
		redisClient.HSet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "time", timeChose).Result()
		phone := evt.Info.Chat.String()
		user, _ := gorm.FindUserByPhone(phone)

		if user == nil {
			msg := `Digite seu *Nome* e *Sobrenome*:.*`
			redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chatStage", "NAME").Result()
			infra.WhatsmeowSendResponse(client, evt, msg)
			return
		}

		msg := "*Seu nome é :*" + " " + user.Name
		msg += "\n\n_1 - SIM_"
		msg += "\n_0 - NÂO_"
		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chatStage", "NAME").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)
		return
	}

	if stage == "DATA" && evt.Message.GetConversation() != "0" {
		redisClient.HSet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "date", evt.Message.GetConversation()).Result()
		date, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "date").Result()
		schedulesFiltered := usecase.FilterSchedules()
		schedulesJSON, _ := json.Marshal(schedulesFiltered)

		redisClient.HSet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "schedules", schedulesJSON).Result()
		msg := `Horários Disponíveis para :*` + date + "* Abaixo:\n\n"

		for i, schedule := range schedulesFiltered {
			msg += fmt.Sprintf("\n%d - *HORA*: %s  *%s* 🕥", i+2, schedule, date)
		}

		msg += "\n\n_0 - VOLTAR  ◀️_"
		msg += "\n\n_Responda com o número correspondente à sua escolha. Para agendar_ 📅"

		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chatStage", "TIME").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)
		return
	}

	if evt.Message.GetConversation() == "1" {
		msg := "*Digite qual data você quer agendar:*"

		msg += "\n\n*Nesse formato separado por barras*"
		msg += "\n_👉 *dia/mes/ano* 👈_"
		msg += "\n_👉 Exemplo: 01/01/2000_"
		msg += "\n\n_0 - VOLTAR ? ◀️_"

		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chatStage", "DATA").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)
		return

	}
	_, err := strconv.Atoi(evt.Message.GetConversation())
	if err != nil {
		infra.WhatsmeowSendResponse(client, evt, "_Desculpe não entendi, pode repetir sua resposta?_")
		return
	}
	schedulesJson, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "schedules").Result()
	var schedules []string
	if err := json.Unmarshal([]byte(schedulesJson), &schedules); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	timeChose := usecase.VerifyScheduleBasedAtArray(evt.Message.GetConversation(), schedules)
	redisClient.HSet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "time", timeChose).Result()

	user, _ := gorm.FindUserByPhone(evt.Info.Chat.String())

	if user == nil {
		msg := `Digite seu *Nome* e *Sobrenome*:*`
		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chatStage", "NAME").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)
	}

	msg := "*Seu nome é :*" + " " + user.Name + "?"
	msg += "\n\n_1 - SIM_"
	msg += "\n_0 - NÂO_"
	redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chatStage", "NAME").Result()
	infra.WhatsmeowSendResponse(client, evt, msg)

}
