package presentation

import (
	"fmt"
	infra "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/whatsmeow"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func SendBotErrorDefault(reason string, client *whatsmeow.Client, evt *events.Message) {
	fmt.Println("Erro to save and init conversation:", reason)
	infra.WhatsmeowSendResponse(client, evt, GetMessage().ErrorDefault)
}
