package presentation

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	infra "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/whatsmeow"
	usecase "github.com/savinnsk/prototype_bot_whatsapp/internal/usecase"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func Init(client *whatsmeow.Client, evt *events.Message, redisClient *redis.Client, currentChatId string) {
	if evt.Message.GetConversation() == "1" {

		schedules, err := usecase.LoadAllUserSchedules(evt.Info.Chat.String())
		if err != nil {
			msg := "_Não há agendamentos registrados._"
			redisClient.HSet(context.Background(), evt.Info.Chat.String(), "currentChatId", "INIT").Result()
			infra.WhatsmeowSendResponse(client, evt, msg)
			return
		}

		msg := "*Seus Agendamentos Abaixo:* \n"

		for _, schedule := range schedules {
			msg += fmt.Sprintf("\n🕥 - *HORA*: %s , *DATA* : %s", schedule.Time, schedule.Date)
		}

		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "currentChatId", "INIT").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)
		return

	}

	if evt.Message.GetConversation() == "2" {
		schedulesFiltered := usecase.FilterSchedules()

		if len(schedulesFiltered) == 0 {
			msg := "_🙁 - Não há agendamentos para hoje."
			msg += "\n\n_1 - AGENDAR OUTRA DATA 📅_"
			msg += "\n_0 - VOLTAR  ◀️_"
			msg += "\n\n_Responda com o número correspondente à sua escolha. Para agendar_ 📅"
			redisClient.HSet(context.Background(), evt.Info.Chat.String(), "currentChatId", "NEW_SCHEDULE").Result()
			redisClient.HSet(context.Background(), evt.Info.Chat.String(), "currentChatId", "INIT").Result()
			infra.WhatsmeowSendResponse(client, evt, msg)
			return
		}

		schedulesJSON, _ := json.Marshal(schedulesFiltered)

		redisClient.HSet(context.Background(), evt.Info.Chat.String()+"NEW_SCHEDULE", "schedules", schedulesJSON).Result()
		msg := "Horários Disponíveis de Hoje: \n\n"

		for i, schedule := range schedulesFiltered {
			msg += fmt.Sprintf("\n%d - *HORA*: %s Hoje 🕥", i+2, schedule)
		}

		msg += "\n\n_1 - AGENDAR OUTRA DATA 📅_"
		msg += "\n_0 - VOLTAR  ◀️_"

		msg += "\n\n_Responda com o número correspondente à sua escolha. Para agendar_ 📅"

		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "currentChatId", "NEW_SCHEDULE").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)
		return

	}

	if evt.Message.GetConversation() == "3" {
		schedules, err := usecase.LoadAllUserSchedules(evt.Info.Chat.String())
		var userScheduleArray []string

		for _, schedule := range schedules {
			userScheduleArray = append(userScheduleArray, schedule.Time)

		}

		schedulesJSON, _ := json.Marshal(userScheduleArray)
		redisClient.HSet(context.Background(), evt.Info.Chat.String()+"CANCEL_SCHEDULE", "schedules", schedulesJSON).Result()
		msg := "*Seus Agendamentos Abaixo:* \n"

		if err != nil {
			msg := "_🙁 - Não há agendamentos para cancelar._"
			redisClient.HSet(context.Background(), evt.Info.Chat.String(), "currentChatId", "INIT").Result()
			infra.WhatsmeowSendResponse(client, evt, msg)
			return
		}
		for i, schedule := range schedules {
			msg += fmt.Sprintf("\n%d - *HORA*: %s , *DATA* : %s", i+2, schedule.Time, schedule.Date)
		}

		msg += "\n\n_0 - VOLTAR  ◀️_"
		msg += "\n\n_Responda com o número correspondente à sua escolha. Para *cancelar ❌*_ 📅"

		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "currentChatId", "CANCEL_SCHEDULE").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)

		return
	}

	if evt.Message.GetConversation() == "4" {
		msg := `*Contato*
		☎️ - Sávio +5522996043721
		`

		infra.WhatsmeowSendResponse(client, evt, msg)

		return
	}

	fields := map[string]interface{}{
		"phone":         evt.Info.Chat.String(),
		"currentChatId": "INIT",
	}

	err := redisClient.HMSet(context.Background(), evt.Info.Chat.String(), fields).Err()
	if err != nil {
		fmt.Println("Erro to save and init conversation:", err)
	}
	// expirationDuration := 10 * time.Minute
	// redisClient.Expire(context.Background(), evt.Info.Chat.String(), expirationDuration).Result()

	msg := `*Olá! Por favor, escolha uma das seguintes opções de 1 a 4:*

1. VER SEUS AGENDAMENTO ? 👁️
2. VER HORÁRIOS DISPONÍVEIS ? 👀
3. CANCELAR UM AGENDAMENTO ? ❌
4. ENTRAR EM CONTATO ? 📞

_Responda com o número correspondente à sua escolha._`

	infra.WhatsmeowSendResponse(client, evt, msg)

}
