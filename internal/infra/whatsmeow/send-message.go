package infra

import (
	"context"
	"go.mau.fi/whatsmeow"
	wp "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func WhatsmeowSendResponse(client *whatsmeow.Client, evt *events.Message, responseText string) string {
	client.SendMessage(context.Background(), evt.Info.Chat, &wp.Message{
		ExtendedTextMessage: &wp.ExtendedTextMessage{
			Text: proto.String(responseText),
			ContextInfo: &wp.ContextInfo{
				QuotedMessage: &wp.Message{
					Conversation: proto.String(evt.Message.GetConversation()),
				},
				StanzaId:    proto.String(evt.Info.ID),
				Participant: proto.String(evt.Info.Sender.ToNonAD().String()),
			},
		},
	})

	return evt.Info.Chat.String()
}
