package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	_ "github.com/mattn/go-sqlite3"
	//"github.com/mdp/qrterminal/v3"
	//"github.com/skip2/go-qrcode"
	infra "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/whatsmeow"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var client *whatsmeow.Client


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
    msg:= `*Olá! Por favor, escolha uma das seguintes opções:*

	1. Deseja Agendar Horário? 📅
	2. Ver Preços? 💲
	3. Verificar seu Horário? ⏰
	4. Ligar para atendente? 📞
	5. Sair. 🚪
 
	
Responda com o número correspondente à sua escolha.`
    // Verifique se a mensagem do usuário é um número válido (1 a 5)
    userInput := evt.Message.GetConversation()
    switch userInput {
    case "1", "2", "3", "4", "5":
        // O usuário escolheu uma opção válida, você pode implementar a lógica para cada opção aqui
        handleUserChoice(evt, userInput)
    default:
        // O usuário não digitou um número válido, envie uma mensagem de erro
        sendErrorMessage(evt, msg)
    }
}
func handleUserChoice(evt *events.Message, choice string) {
    // Implemente a lógica para cada opção escolhida pelo usuário
    switch choice {
    case "1":
        // O usuário escolheu a opção 1: Deseja Agendar Horário
        // Implemente a lógica para esta opção aqui
        infra.WhatsmeowSendResponse(client ,evt, `*Você escolheu a opção 1: Deseja Agendar Horário.*
        
        os seguintes horários estão disponíveis:

        1 - 12:00am
        2 - 15:00pm
        
        `)


        infra.WhatsmeowSendResponse(client,evt, "*Por favor, responda com o número correspondente ao horário desejado.*")
    case "2":
        // O usuário escolheu a opção 2: Ver Preços
        // Implemente a lógica para esta opção aqui
        infra.WhatsmeowSendResponse(client,evt, "Você escolheu a opção 2: Ver Preços.")
    case "3":
        // O usuário escolheu a opção 3: Verificar seu Horário
        // Implemente a lógica para esta opção aqui
        infra.WhatsmeowSendResponse(client, evt, "Você escolheu a opção 3: Verificar seu Horário.")
    case "4":
        // O usuário escolheu a opção 4: Ligar para atendente
        // Implemente a lógica para esta opção aqui
        infra.WhatsmeowSendResponse(client,evt, "Você escolheu a opção 4: Ligar para atendente.")
    case "5":
        // O usuário escolheu a opção 5: Sair
        // Implemente a lógica para esta opção aqui
        infra.WhatsmeowSendResponse(client,evt, "Você escolheu a opção 5: Sair.")
    }
}




// Função para enviar uma mensagem de erro ao usuário
func sendErrorMessage(evt *events.Message, errorMessage string) {
    infra.WhatsmeowSendResponse(client, evt, errorMessage)
}

func main() {
    configureLogging()

    container, err := initializeSQLStore()
    handleError(err)
  

	deviceStore, err := getFirstDeviceFromContainer(container)
	handleError(err)


	initializeWhatsMeowClient(deviceStore)

    handleLogin()

    waitForInterruptSignal()
}

func configureLogging() {
    dbLog := waLog.Stdout("Database", "INFO", true)
    _ = dbLog
  
}

func initializeSQLStore() (*sqlstore.Container, error) {
	dbLog := waLog.Stdout("Database", "INFO", true)
    return sqlstore.New("sqlite3", "file:../examplestore.db?_foreign_keys=on", dbLog)
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

		// Create a file to save the QR code as a PNG image
		qrFile, err := os.Create("../qrcode.png")
	
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
