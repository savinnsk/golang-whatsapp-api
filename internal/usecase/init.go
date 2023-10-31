package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	gorm "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/gorm"
	infra "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/whatsmeow"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func Init(client *whatsmeow.Client, evt *events.Message, redisClient *redis.Client, currentChatId string) {

	if evt.Message.GetConversation() == "1" {
		msg := `*Seus Agendamentos Abaixo:*

1 - *HORA* : ~11:00~ - *DATA* 12/12/2023 🕥
2 - *HORA* : ~13:00~ - *DATA* 12/12/2023 🕥

	  
_0 - VOLTAR ? ◀️_`

		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "currentChatId", "SHOW_USER_SCHEDULE").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)
		return

	}

	if evt.Message.GetConversation() == "2" {
		schedules, err := gorm.LoadAllSchedules()
		if err != nil {
			println("Error to load Schedules")
		}

		msg := `*Horários De Hoje Disponíveis Abaixo:*
		`

		for i, schedule := range schedules {
			if schedule.Available && !schedule.Disabled {
				msg += fmt.Sprintf("\n%d - *HORA*: %s  *HOJE* 🕥", i+2, schedule.Time)
			}
		}

		msg += "\n\n_1 - AGENDAR / OUTRA DATA 📅_"
		msg += "\n_0 - VOLTAR  ◀️_"

		redisClient.HSet(context.Background(), evt.Info.Chat.String(), "currentChatId", "NEW_SCHEDULE").Result()
		infra.WhatsmeowSendResponse(client, evt, msg)
		return

	}

	if evt.Message.GetConversation() == "3" {
		msg := `*Qual dos seus horários você deseja cancelar?:*

1 - *HORA* : ~11:00~ - *DATA* 12/12/25 🕥
2 - *HORA* : ~13:00~ - *DATA* 12/12/25 🕥


_0 - VOLTAR ◀️_

_Responda com o número correspondente à sua escolha._
		`

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
	expirationDuration := 10 * time.Minute
	redisClient.Expire(context.Background(), evt.Info.Chat.String(), expirationDuration).Result()

	msg := `*Olá! Por favor, escolha uma das seguintes opções de 0 a 4:*

1. VER SEU AGENDAMENTO ? 👁️
2. VER HORÁRIOS DISPONÍVEIS ? 👀
3. CANCELAR UM AGENDAMENTO ? ❌
4. ENTRAR EM CONTATO ? 📞

_Responda com o número correspondente à sua escolha._`

	infra.WhatsmeowSendResponse(client, evt, msg)

}
