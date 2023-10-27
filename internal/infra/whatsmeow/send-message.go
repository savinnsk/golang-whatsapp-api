package infra

import (
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
	wp "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"
	"context"
)

func WhatsmeowSendResponse(client *whatsmeow.Client , evt *events.Message, responseText string) string {
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

    return evt.Message.GetConversation()
}