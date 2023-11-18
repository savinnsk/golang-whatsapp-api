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

	messageFromClient := evt.Message.GetConversation()
	_, errNotNumeric := strconv.Atoi(messageFromClient)

	stage, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String(), "chat.stage").Result()

	if evt.Message.GetConversation() == "1" && stage == "" {
		msg := "*Digite qual data você quer agendar:*"

		msg += "\n\n*Nesse formato separado por barras*"
		msg += "\n_👉 *dia/mes/ano* 👈_"
		msg += "\n_👉 Exemplo: 01/01/2000_"
		msg += "\n\n_0 - VOLTAR ? ◀️_"

		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chat.stage", "schedule.another.date").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)
		return

	}

	if stage == "schedule.another.date" {

		if messageFromClient == "0" {
			redisClient.HDel(context.Background(), evt.Info.Chat.String(), "chat.stage")
			redisClient.HDel(context.Background(), evt.Info.Chat.String(), "schedules")
			Init(client, evt, redisClient)
			return
		}

		targetDate, err := t.Parse("02/01/2006", evt.Message.GetConversation())

		if err != nil {
			fmt.Println(err, "line 43 schedules.message")
			infra.WhatsmeowSendResponse(client, evt, GetMessage().NotUnderstand)
			return
		}

		if targetDate.Before(t.Now()) {
			infra.WhatsmeowSendResponse(client, evt, GetMessage().VerifyValues)
			return
		}

		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "target.date", targetDate).Result()

		schedulesFiltered := usecase.FilterSchedules()
		schedulesJSON, _ := json.Marshal(schedulesFiltered)

		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "schedules", schedulesJSON).Result()
		msg := `Horários Disponíveis para :*` + targetDate.Format("02/01/2006") + "* Abaixo:\n\n"

		for i, schedule := range schedulesFiltered {
			msg += fmt.Sprintf("\n%d - *HORA*: %s  *%s* 🕥", i+2, schedule, targetDate.Format("02/01/2006"))
		}

		msg += "\n\n_0 - VOLTAR  ◀️_"
		msg += "\n\n_Responda com o número correspondente à sua escolha. Para agendar_ 📅"

		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chat.stage", "chat.time").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)
		return
	}

	if stage == "chat.time" && evt.Message.GetConversation() != "0" {

		schedulesJson, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String(), "schedules").Result()
		var schedules []string
		if err := json.Unmarshal([]byte(schedulesJson), &schedules); err != nil {
			fmt.Print("\nError unmarshaling JSON:", err, "\n")
			return
		}
		timeChose, err := usecase.VerifyScheduleBasedAtArray(evt.Message.GetConversation(), schedules)

		if errNotNumeric != nil || err != nil || timeChose == "" {
			infra.WhatsmeowSendResponse(client, evt, GetMessage().VerifyValues)
			return
		}
		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chat.time.chose", timeChose)
		phone := evt.Info.Chat.String()
		user, _ := gorm.FindUserByPhone(phone)

		if user == nil {
			msg := `Digite seu *Nome* e *Sobrenome*:.*`
			redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chat.stage", "chat.get.name").Result()
			infra.WhatsmeowSendResponse(client, evt, msg)
			return
		}

		msg := "*Seu nome é :*" + " " + user.Name
		msg += "\n\n_1 - SIM_"
		msg += "\n_0 - NÂO_"
		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chat.stage", "chat.get.name").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)
		return
	}

	if stage == "chat.get.name" {
		date, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String(), "target.date").Result()
		timeChose, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String(), "chat.time.chose").Result()
		targetDate, _ := t.Parse("2006-01-02T15:04:05Z", date)
		targetDateFormatted := targetDate.Format("02/01/2006")

		user, err := gorm.FindUserByPhone(evt.Info.Chat.String())

		if err != nil && user != nil {
			result := SaveScheduleWithAUserAlreadySignUp(user.Name, user.Phone, targetDateFormatted, timeChose)

			if result != "ok" {
				SendBotErrorDefault("Error to save user with account", client, evt)

			}
		}

		result := SaveScheduleWithAUserAlreadySignUp(evt.Message.GetConversation(), evt.Info.Chat.String(), targetDateFormatted, timeChose)
		if result != "ok" {
			SendBotErrorDefault("Error to save user with no account", client, evt)
		}
		//
		msg := fmt.Sprintf("_seu agendamento foi concluído_ para : *%s* 🕥", targetDate.Format("02/01/2006"))
		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chat.stage", "finish").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)
		return
	}

	if evt.Message.GetConversation() == "0" {
		BackToMenu(client, evt, redisClient)
	}

	// if stage == "SAVE" && evt.Message.GetConversation() != "0" {
	// 	time, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String(), "time").Result()
	// 	date, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String(), "date").Result()
	// 	name, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String(), "name").Result()

	// 	userData := dto.CreateUserDto{
	// 		Name:  name,
	// 		Phone: evt.Info.Chat.String(),
	// 		Role:  "user",
	// 	}

	// 	if date == "" {
	// 		currentDate := t.Now().Format("02/01/2006")
	// 		date = currentDate
	// 	}
	// 	data := dto.SaveNewUserAndSchedule{
	// 		CreateUserDto: userData,
	// 		ScheduleTime:  time,
	// 		ScheduleDate:  date,
	// 	}

	// 	result := usecase.ProcessNewSchedule(data)
	// 	if result != "ok" {
	// 		redisClient.Del(context.Background(), evt.Info.Chat.String())
	// 		redisClient.Del(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE")
	// 		infra.WhatsmeowSendResponse(client, evt, "\n\n_"+result+"_")
	// 		return
	// 	}

	// 	msg := `_😀 SEU AGENDAMENTO FOI SALVO COM SUCESSO_`

	// 	redisClient.Del(context.Background(), evt.Info.Chat.String())
	// 	redisClient.Del(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE")
	// 	infra.WhatsmeowSendResponse(client, evt, msg)

	// 	return
	// }

	// if stage == "NAME" && evt.Message.GetConversation() != "0" {
	// 	redisClient.HSet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "name", evt.Message.GetConversation()).Result()
	// 	msg := `*Confirmar Agendamento ?*`
	// 	msg += "\n\n_1 - SIM_"
	// 	msg += "\n_0 - NÂO_"
	// 	redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chatStage", "SAVE").Result()
	// 	infra.WhatsmeowSendResponse(client, evt, msg)
	// 	return
	// }

	// if stage == "DATA" && evt.Message.GetConversation() != "0" {
	// 	redisClient.HSet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "date", evt.Message.GetConversation()).Result()
	// 	date, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "date").Result()
	// 	schedulesFiltered := usecase.FilterSchedules()
	// 	schedulesJSON, _ := json.Marshal(schedulesFiltered)

	// 	redisClient.HSet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "schedules", schedulesJSON).Result()
	// 	msg := `Horários Disponíveis para :*` + date + "* Abaixo:\n\n"

	// 	for i, schedule := range schedulesFiltered {
	// 		msg += fmt.Sprintf("\n%d - *HORA*: %s  *%s* 🕥", i+2, schedule, date)
	// 	}

	// 	msg += "\n\n_0 - VOLTAR  ◀️_"
	// 	msg += "\n\n_Responda com o número correspondente à sua escolha. Para agendar_ 📅"

	// 	redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chatStage", "TIME").Result()
	// 	infra.WhatsmeowSendResponse(client, evt, msg)
	// 	return
	// }

	// schedulesJson, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "schedules").Result()
	// var schedules []string
	// if err := json.Unmarshal([]byte(schedulesJson), &schedules); err != nil {
	// 	fmt.Println("Error unmarshaling JSON:", err)
	// 	return
	// }

	// timeChose, _ := usecase.VerifyScheduleBasedAtArray(evt.Message.GetConversation(), schedules)
	// redisClient.HSet(context.Background(), evt.Info.Chat.String(), "time.chose", timeChose).Result()

	// user, _ := gorm.FindUserByPhone(evt.Info.Chat.String())

	// if user == nil {
	// 	msg := `Digite seu *Nome* e *Sobrenome*:*`
	// 	redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chatStage", "NAME").Result()
	// 	infra.WhatsmeowSendResponse(client, evt, msg)
	// }

	// msg := "*Seu nome é :*" + " " + user.Name + "?"
	// msg += "\n\n_1 - SIM_"
	// msg += "\n_0 - NÂO_"
	// redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chatStage", "NAME").Result()
	// infra.WhatsmeowSendResponse(client, evt, msg)

}

func BackToMenu(client *whatsmeow.Client, evt *events.Message, redisClient *redis.Client) {
	Init(client, evt, redisClient)

}

func ScheduleAnotherDate(client *whatsmeow.Client, evt *events.Message, redisClient *redis.Client) {
	msg := "*Digite qual data você quer agendar:*"

	msg += "\n\n*Nesse formato separado por barras*"
	msg += "\n_👉 *dia/mes/ano* 👈_"
	msg += "\n_👉 Exemplo: 01/01/2000_"
	msg += "\n\n_0 - VOLTAR ? ◀️_"

	redisClient.HSet(context.Background(), evt.Info.Chat.String(), "chat.stage", "schedule.another.date").Result()
	infra.WhatsmeowSendResponse(client, evt, msg)
}

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
