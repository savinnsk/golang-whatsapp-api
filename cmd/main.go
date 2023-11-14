package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v8"
	_ "github.com/mattn/go-sqlite3"

	//"github.com/mdp/qrterminal/v3"
	//"github.com/skip2/go-qrcode"
	gr "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/gorm"
	redisInstance "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/redis"
	whatsmeowInstance "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/whatsmeow"
	evPresent "github.com/savinnsk/prototype_bot_whatsapp/internal/presentation/bot-messages"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var client *whatsmeow.Client
var redisClient *redis.Client

func eventHandler(evt interface{}) {
	if evt, ok := evt.(*events.Message); ok {
		evPresent.EventsMapper(client, evt, redisClient)
	}
}

func main() {
	configureLogging()

	redisClient = redisInstance.Init()
	gr.Init()
	container, err := initializeSQLStore()
	handleError(err)

	deviceStore, err := getFirstDeviceFromContainer(container)
	handleError(err)

	whatsmeow := whatsmeowInstance.Init(deviceStore)
	whatsmeow.AddEventHandler(eventHandler)
	client = whatsmeow
	handleLogin()

	waitForInterruptSignal()
}

func configureLogging() {
	dbLog := waLog.Stdout("Database", "INFO", true)
	_ = dbLog

}

func initializeSQLStore() (*sqlstore.Container, error) {
	dbLog := waLog.Stdout("Database", "INFO", true)
	return sqlstore.New("sqlite3", "file:whatsmeow.db?_foreign_keys=on", dbLog)
}

func getFirstDeviceFromContainer(container *sqlstore.Container) (*store.Device, error) {
	return container.GetFirstDevice()
}

func handleLogin() {
	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err := client.Connect()
		handleError(err)

		// Create a file to save the QR code as a PNG image
		qrFile, err := os.Create("qrcode.png")

		if err != nil {
			panic(err)
		}
		handleError(err)
		defer qrFile.Close()

		// Print the QR code to the file
		for evt := range qrChan {
			if evt.Event == "code" {
				// Print the QR code to the file

				_, err := qrFile.Write([]byte(evt.Code))
				handleError(err)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		err := client.Connect()
		handleError(err)
	}
}

func waitForInterruptSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
