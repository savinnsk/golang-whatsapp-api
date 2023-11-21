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
	redisC "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/redis"
	infra "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/whatsmeow"
	usecase "github.com/savinnsk/prototype_bot_whatsapp/internal/usecase"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func NewSchedule(cs *ChatSetup) {

	messageFromClient := cs.evt.Message.GetConversation()
	_, errNotNumeric := strconv.Atoi(messageFromClient)
	stageConversation, _ := redisC.GetValues(context.TODO(), cs.evt.Info.Chat.String(), []string{"chat.stage"}) //cs.redisClient.HGet(context.Background(), cs.evt.Info.Chat.String(), "chat.stage").Result()

	if messageFromClient == "1" && stageConversation[0] == "" {
		Option1ScheduleAnotherDate(cs)

	}

	if stageConversation[0] == "schedule.another.date" {

		if messageFromClient == "0" {
			cs.redisClient.HDel(context.Background(), cs.evt.Info.Chat.String(), "chat.stage")
			cs.redisClient.HDel(context.Background(), cs.evt.Info.Chat.String(), "schedules")
			Init(cs.client, cs.evt, cs.redisClient)
			return
		}

		targetDate, err := t.Parse("02/01/2006", cs.evt.Message.GetConversation())

		if err != nil {
			fmt.Println(err, "line 43 schedules.message")
			infra.WhatsmeowSendResponse(cs.client, cs.evt, GetMessage().NotUnderstand)
			return
		}

		if targetDate.Before(t.Now()) {
			infra.WhatsmeowSendResponse(cs.client, cs.evt, GetMessage().VerifyValues)
			return
		}

		cs.redisClient.HSet(context.Background(), cs.evt.Info.Chat.String(), "target.date", targetDate).Result()

		schedulesFiltered := usecase.FilterSchedules()
		schedulesJSON, _ := json.Marshal(schedulesFiltered)

		cs.redisClient.HSet(context.Background(), cs.evt.Info.Chat.String(), "schedules", schedulesJSON).Result()
		msg := `Horários Disponíveis para :*` + targetDate.Format("02/01/2006") + "* Abaixo:\n\n"

		for i, schedule := range schedulesFiltered {
			msg += fmt.Sprintf("\n%d - *HORA*: %s  *%s* 🕥", i+2, schedule, targetDate.Format("02/01/2006"))
		}

		msg += "\n\n_0 - VOLTAR  ◀️_"
		msg += "\n\n_Responda com o número correspondente à sua escolha. Para agendar_ 📅"

		cs.redisClient.HSet(context.Background(), cs.evt.Info.Chat.String(), "chat.stage", "chat.time").Result()
		infra.WhatsmeowSendResponse(cs.client, cs.evt, msg)
		return
	}

	if stageConversation[0] == "chat.time" && cs.evt.Message.GetConversation() != "0" {

		schedulesJson, _ := cs.redisClient.HGet(context.Background(), cs.evt.Info.Chat.String(), "schedules").Result()
		var schedules []string
		if err := json.Unmarshal([]byte(schedulesJson), &schedules); err != nil {
			fmt.Print("\nError unmarshaling JSON:", err, "\n")
			return
		}
		timeChose, err := usecase.VerifyScheduleBasedAtArray(cs.evt.Message.GetConversation(), schedules)

		if errNotNumeric != nil || err != nil || timeChose == "" {
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

	if stageConversation[0] == "chat.get.name" {
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
		BackToMenu(cs.client, cs.evt, cs.redisClient)
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

func Option1ScheduleAnotherDate(cs *ChatSetup) {
	err := redisC.AddValues(context.Background(), cs.evt.Info.Chat.String(), []string{"chat.stage"}, []string{"schedule.another.date"})
	if err != nil {
		print("💀 error at : schedules.message.go / line : 236")
	}
	infra.WhatsmeowSendResponse(cs.client, cs.evt, GetMessage().ScheduleAnotherDate)
}

func HandleScheduleAnotherDate(cs *ChatSetup) {
	messageFromClient := cs.evt.Message.GetConversation()

	if messageFromClient == "0" {
		cs.redisClient.HDel(context.Background(), cs.evt.Info.Chat.String(), "chat.stage")
		cs.redisClient.HDel(context.Background(), cs.evt.Info.Chat.String(), "schedules")
		Init(cs.client, cs.evt, cs.redisClient)
		return
	}

	targetDate, err := t.Parse("02/01/2006", cs.evt.Message.GetConversation())

	if err != nil {
		fmt.Println(err, "line 43 schedules.message")
		infra.WhatsmeowSendResponse(cs.client, cs.evt, GetMessage().NotUnderstand)
		return
	}

	if targetDate.Before(t.Now()) {
		infra.WhatsmeowSendResponse(cs.client, cs.evt, GetMessage().VerifyValues)
		return
	}

	cs.redisClient.HSet(context.Background(), cs.evt.Info.Chat.String(), "target.date", targetDate).Result()

	schedulesFiltered := usecase.FilterSchedules()
	schedulesJSON, _ := json.Marshal(schedulesFiltered)

	cs.redisClient.HSet(context.Background(), cs.evt.Info.Chat.String(), "schedules", schedulesJSON).Result()
	msg := `Horários Disponíveis para :*` + targetDate.Format("02/01/2006") + "* Abaixo:\n\n"

	for i, schedule := range schedulesFiltered {
		msg += fmt.Sprintf("\n%d - *HORA*: %s  *%s* 🕥", i+2, schedule, targetDate.Format("02/01/2006"))
	}

	msg += "\n\n_0 - VOLTAR  ◀️_"
	msg += "\n\n_Responda com o número correspondente à sua escolha. Para agendar_ 📅"

	cs.redisClient.HSet(context.Background(), cs.evt.Info.Chat.String(), "chat.stage", "chat.time").Result()
	infra.WhatsmeowSendResponse(cs.client, cs.evt, msg)
	return
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

func BackToMenu(client *whatsmeow.Client, evt *events.Message, redisClient *redis.Client) {
	Init(client, evt, redisClient)

}
