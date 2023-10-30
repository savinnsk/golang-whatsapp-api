package usecase

import (
	"context"

	"github.com/go-redis/redis/v8"
	infra "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/whatsmeow"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func Init(client *whatsmeow.Client, evt *events.Message, redisClient *redis.Client) {

	msg := `*Olá! Por favor, escolha uma das seguintes opções:*

	1. Deseja Agendar Horário? 📅
	2. Ver Preços? 💲
	3. Verificar seu Horário? ⏰
	4. Ligar para atendente? 📞
	5. Sair. 🚪
 
	
	Responda com o número correspondente à sua escolha.`

	chatNumber := infra.WhatsmeowSendResponse(client, evt, msg)
	result, _ := redisClient.HSet(context.Background(), chatNumber, "currentChatId").Result()
	// currentChatId, _ := redisClient.HGet(context.Background(), evt.Info.Chat.String(), "currentChatId").Result()
	// if result != "" && currentChatId != "Initial" {
	//     fmt.Println("Has a active conversation ")
	//     return

	// }
	switch res {
	case "1":

		infra.WhatsmeowSendResponse(client, evt, `*Você escolheu a opção 1: Deseja Agendar Horário.*
        
        os seguintes horários estão disponíveis:

        1 - 12:00am
        2 - 15:00pm
        
        `)

		infra.WhatsmeowSendResponse(client, evt, "*Por favor, responda com o número correspondente ao horário desejado.*")

		evt.Message.Reset()
	case "2":

		infra.WhatsmeowSendResponse(client, evt, "Você escolheu a opção 2: Ver Preços.")
	case "3":

		infra.WhatsmeowSendResponse(client, evt, "Você escolheu a opção 3: Verificar seu Horário.")
	case "4":

		infra.WhatsmeowSendResponse(client, evt, "Você escolheu a opção 4: Ligar para atendente.")
	case "5":

		infra.WhatsmeowSendResponse(client, evt, "Você escolheu a opção 5: Sair.")
	}
}
