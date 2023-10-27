package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	wp "go.mau.fi/whatsmeow/binary/proto"

	// Import the SQLite driver package
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

var client *whatsmeow.Client

// Event handler handles incoming events and dispatches them to the appropriate functions
func eventHandler(evt interface{}) {
	// Check if event is a message
	if evt, ok := evt.(*events.Message); ok {
		// Handle text messages
		if evt.Message.GetConversation() != "" {
			go handleConversation(evt)
		}
	}
}

// HandleConversation handles incoming text messages by generating an AI-generated text response and sending it back
func handleConversation(evt *events.Message) {
	msg := `*Olá! Por favor, escolha uma das seguintes opções:*

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
	switch choice {
	case "1":
		// O usuário escolheu a opção 1: Deseja Agendar Horário
		// Implemente a lógica para esta opção aqui
		//oisendResponse(evt, "*Você escolheu a opção 1: Deseja Agendar Horário.*")
		msg1 := &wp.ListMessage{
			Title:       proto.String("welcome"),
			Description: proto.String("test 2"),
			ButtonText:  proto.String("kk"),
			ListType:    wp.ListMessage_SINGLE_SELECT.Enum(),
			Sections: []*wp.ListMessage_Section{
				{
					Title: proto.String("test"),
					Rows: []*wp.ListMessage_Row{
						{
							RowId: proto.String("1"),
							Title: proto.String("testanto 1"),
							//	Description: proto.String("عمادة القبول والتسجيل"),
						},
						{
							RowId: proto.String("2"),
							Title: proto.String("testando 2"),
						},
						{
							RowId: proto.String("3"),
							Title: proto.String("testando 3"),
						},
					},
				},
			}}


	msg2 := &wp.ButtonsMessage{
		ContentText: proto.String("test"),
		HeaderType:  wp.ButtonsMessage_EMPTY.Enum(),
		Buttons: []*wp.ButtonsMessage_Button{
			{
				ButtonId:       proto.String("222"),
				ButtonText:     &wp.ButtonsMessage_Button_ButtonText{DisplayText: proto.String("222")},
				Type:           wp.ButtonsMessage_Button_RESPONSE.Enum(),
				NativeFlowInfo: &wp.ButtonsMessage_Button_NativeFlowInfo{},
			},
			{
				ButtonId:       proto.String("2222"),
				ButtonText:     &wp.ButtonsMessage_Button_ButtonText{DisplayText: proto.String("s")},
				Type:           wp.ButtonsMessage_Button_RESPONSE.Enum(), //proto.ButtonsMessage_Button_Type.Enum,
				NativeFlowInfo: &wp.ButtonsMessage_Button_NativeFlowInfo{},
			},
		},
	}
			//	ProductListInfo: &waProto.ListMessage_ProductListInfo{},
			//	FooterText:      new(string),
			//	ContextInfo:     &waProto.ContextInfo{},
		
		// Adicione os botões de agendamento
		// buttons := []*wp.HydratedTemplateButton{
		// 	{
		// 		Index: proto.Uint32(uint32(1)),
		// 		HydratedButton: &wp.HydratedTemplateButton_QuickReplyButton{
		// 			QuickReplyButton: &wp.HydratedTemplateButton_HydratedQuickReplyButton{
		// 				DisplayText: proto.String("Agendar para amanhã"),
		// 				Id:          proto.String("ScheduleTomorrow"),
		// 			},
		// 		},
		// 	},
		// 	{
		// 		Index: proto.Uint32(uint32(2)),
		// 		HydratedButton: &wp.HydratedTemplateButton_QuickReplyButton{
		// 			QuickReplyButton: &wp.HydratedTemplateButton_HydratedQuickReplyButton{
		// 				DisplayText: proto.String("Agendar para depois de amanhã"),
		// 				Id:          proto.String("ScheduleDayAfterTomorrow"),
		// 			},
		// 		},
		// 	},
		// 	{
		// 		Index: proto.Uint32(uint32(3)),
		// 		HydratedButton: &wp.HydratedTemplateButton_QuickReplyButton{
		// 			QuickReplyButton: &wp.HydratedTemplateButton_HydratedQuickReplyButton{
		// 				DisplayText: proto.String("Cancelar"),
		// 				Id:          proto.String("Cancel"),
		// 			},
		// 		},
		// 	},
		// }

		sendButtonResponseOnce(evt, msg2)
		sendListResponseOnce(evt, msg1)
		sendButtonResponse(evt, msg2)
		sendListResponse(evt, msg1)
		sendResponse(evt , "done")
	case "2":
		// O usuário escolheu a opção 2: Ver Preços
		// Implemente a lógica para esta opção aqui
		sendResponse(evt, "Você escolheu a opção 2: Ver Preços.")
	case "3":
		// O usuário escolheu a opção 3: Verificar seu Horário
		// Implemente a lógica para esta opção aqui
		sendResponse(evt, "Você escolheu a opção 3: Verificar seu Horário.")
	case "4":
		// O usuário escolheu a opção 4: Ligar para atendente
		// Implemente a lógica para esta opção aqui
		sendResponse(evt, "Você escolheu a opção 4: Ligar para atendente.")
	case "5":
		// O usuário escolheu a opção 5: Sair
		// Implemente a lógica para esta opção aqui
		sendResponse(evt, "Você escolheu a opção 5: Sair.")
	}
}

func sendResponse(evt *events.Message, responseText string) string {
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

// Função para enviar uma mensagem de erro ao usuário
func sendErrorMessage(evt *events.Message, errorMessage string) {
	sendResponse(evt, errorMessage)
}

// Função para enviar uma mensagem com botões
func sendButtonResponseOnce(evt *events.Message, buttons *wp.ButtonsMessage) any{
	client.SendMessage(context.Background(), evt.Info.Chat, &wp.Message{
		ViewOnceMessage: &wp.FutureProofMessage{
			Message: &wp.Message{
				ButtonsMessage: buttons,
				Conversation: proto.String(evt.Message.GetConversation()),
				// ListMessage: buttons,
			},
		},
	})

	return evt.Message.GetViewOnceMessage()
}

func sendListResponseOnce(evt *events.Message, buttons *wp.ListMessage) any{
	client.SendMessage(context.Background(), evt.Info.Chat, &wp.Message{
		ViewOnceMessage: &wp.FutureProofMessage{
			Message: &wp.Message{
				ListMessage: buttons,
				Conversation: proto.String(evt.Message.GetConversation()),
				// ListMessage: buttons,
			},
		},
	})

	return evt.Message.GetViewOnceMessage()
}


func sendButtonResponse(evt *events.Message, buttons *wp.ButtonsMessage) any{
	client.SendMessage(context.Background(), evt.Info.Chat, &wp.Message{
				ButtonsMessage:  &wp.ButtonsMessage{
					ContentText: proto.String("test"),
					HeaderType:  wp.ButtonsMessage_EMPTY.Enum(),
					Buttons: []*wp.ButtonsMessage_Button{
						{
							ButtonId:       proto.String("222"),
							ButtonText:     &wp.ButtonsMessage_Button_ButtonText{DisplayText: proto.String("222")},
							Type:           wp.ButtonsMessage_Button_RESPONSE.Enum(),
							NativeFlowInfo: &wp.ButtonsMessage_Button_NativeFlowInfo{},
						},
						{
							ButtonId:       proto.String("2222"),
							ButtonText:     &wp.ButtonsMessage_Button_ButtonText{DisplayText: proto.String("s")},
							Type:           wp.ButtonsMessage_Button_RESPONSE.Enum(), //proto.ButtonsMessage_Button_Type.Enum,
							NativeFlowInfo: &wp.ButtonsMessage_Button_NativeFlowInfo{},
						},
					},
				},
				Conversation: proto.String(evt.Message.GetConversation()),
				// ListMessage: buttons,
			
		
	})

	return evt.Message.GetConversation()
}

func sendListResponse(evt *events.Message, buttons *wp.ListMessage) any{
	client.SendMessage(context.Background(), evt.Info.Chat, &wp.Message{
		
		
				ListMessage: buttons,
				Conversation: proto.String(evt.Message.GetConversation()),
				// ListMessage: buttons,
			
		
	})

	return evt.Message.GetListMessage()
}

// func sendListResponseNew(evt *events.Message, buttons *wp.ListMessage) any{
// 	client.SendMessage(context.Background(), evt.Info.Chat, &wp.Message{
// 		ExtendedTextMessage: &wp.ExtendedTextMessage{
// 			ContextInfo: &wp.ContextInfo{
// 				QuotedMessage: &wp.Message{
// 					Conversation: proto.String(evt.Message.GetConversation()),
// 				},
// 				StanzaId:    proto.String(evt.Info.ID),
// 				Participant: proto.String(evt.Info.Sender.ToNonAD().String()),
// 			},
// 		},
// 	})

// 	return evt.Message.GetListMessage()
// }


func main() {
	configureLogging()

	container, err := initializeSQLStore()
	handleError(err)

	deviceStore, err := getFirstDeviceFromContainer(container)
	handleError(err)

	initializeWhatsMeowClient(deviceStore)

	// Check if the client is already logged in or perform the login process
	handleLogin()

	waitForInterruptSignal()
}

func configureLogging() {
	dbLog := waLog.Stdout("Database", "INFO", true)
	_ = dbLog
}

func initializeSQLStore() (*sqlstore.Container, error) {
	dbLog := waLog.Stdout("Database", "INFO", true)
	return sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
}

func getFirstDeviceFromContainer(container *sqlstore.Container) (*store.Device, error) {
	return container.GetFirstDevice()
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

		qrFile, err := os.Create("qrcode.png")
		if err != nil {
			panic(err)
		}
		handleError(err)
		defer qrFile.Close()

		for evt := range qrChan {
			if evt.Event == "code" {
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