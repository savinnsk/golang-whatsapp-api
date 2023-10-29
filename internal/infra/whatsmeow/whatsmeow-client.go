package infra

import (
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var client *whatsmeow.Client

type EventHandler interface {
	HandleEvent(evt interface{})
}

func InitializeWhatsMeowClient(deviceStore *store.Device, e EventHandler) {
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(e.HandleEvent)
}
