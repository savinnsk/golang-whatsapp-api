package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	// Import the SQLite driver package
	_ "github.com/mattn/go-sqlite3"
	//"github.com/mdp/qrterminal/v3"
	//"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	wp "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
	//"image/png"
)

var client *whatsmeow.Client

// eventHandler handles incoming events and dispatches them to the appropriate functions
func eventHandler(evt interface{}) {
    // Check if event is a message
    if evt, ok := evt.(*events.Message); ok {
        // Handle image messages

        // Handle text messages
        if evt.Message.GetConversation() != "" {
            go handleConversation(evt)
        }
    }
}

// handleConversation handles incoming text messages by generating an AI-generated text response and sending it back
func handleConversation(evt *events.Message) {
    // Generate an AI-generated text response using the message text
    msg := "oi Sávio não estar deixe seu recado"

    // Create a message to send back containing the generated text
    client.SendMessage(context.Background(), evt.Info.Chat, &wp.Message{
        ExtendedTextMessage: &wp.ExtendedTextMessage{
            Text: proto.String(msg),
            ContextInfo: &wp.ContextInfo{
                QuotedMessage: &wp.Message{
                    Conversation: proto.String(evt.Message.GetConversation()),
                },
                StanzaId:    proto.String(evt.Info.ID),
                Participant: proto.String(evt.Info.Sender.ToNonAD().String()),
            },
        },
    })
}

func main() {
    configureLogging()

    container, err := initializeSQLStore()
    handleError(err)

  

	deviceStore, err := getFirstDeviceFromContainer(container)
	handleError(err)

    handleError(err)


	initializeWhatsMeowClient(deviceStore)

    // Check if the client is already logged in or perform the login process
    handleLogin()

    waitForInterruptSignal()
}

func configureLogging() {
    dbLog := waLog.Stdout("Database", "INFO", true)
    // Use dbLog if needed for database logging
    _ = dbLog
    // Configure other logging settings if needed
    // ...
}

func initializeSQLStore() (*sqlstore.Container, error) {
	dbLog := waLog.Stdout("Database", "INFO", true)
    return sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
}

func getFirstDeviceFromContainer(container *sqlstore.Container) (*store.Device, error)  {
 return  container.GetFirstDevice()  
}

func initializeWhatsMeowClient(deviceStore *store.Device) {
    clientLog := waLog.Stdout("Client", "DEBUG", true)
    client = whatsmeow.NewClient(deviceStore, clientLog)
    client.AddEventHandler(eventHandler)
}

func handleLogin() {
    if client.Store.ID == nil {
        qrChan, _ := client.GetQRChannel(context.Background())
        err := client.Connect()
        handleError(err)

        // Print the QR code to the console
        for evt := range qrChan {
            if evt.Event == "code" {
                // qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
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
