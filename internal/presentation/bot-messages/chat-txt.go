package presentation

type Message struct {
	Greetings               string
	MenuInteractionText     string
	SchedulesAvailableTitle string
	ErrorDefault            string
	SchedulesNotFound       string
	NotUnderstand           string
	BackButton              string
	ScheduleOtherTime       string
	DefaultFooter           string
	Contacts                string
}

func GetMessage() Message {
	// INIT
	var Greetings = "*👋 Bem-vindo esse é um bot de teste da IS Idea Software!*\n\n Somos apaixonados por transformar atendimentos e experiências com a ajuda da tecnologia. ✨\n\n 🤖 O foco desse bot é ajudar em processos de atendimento e automação de processos de agendamento, cancelamento, e notificar os interessados.\n\n 💼 Oferecemos soluções de automação de atendimento, criamos bots inteligentes para agilizar suas interações com os clientes e fornecemos serviços de marketing digital para impulsionar seus negócios.\n\n  🚀 Nossa equipe dedicada e experiente está pronta para ajudar sua empresa a atingir seus objetivos, seja através de automação, chatbots inovadores ou soluções de software personalizadas.\n\n💻 Descubra como podemos trabalhar juntos para criar experiências excepcionais e melhorar seus resultados. 💬\n\nEntre em contato conosco e leve sua empresa a um novo patamar com a IS Idea Software. Estamos ansiosos para ajudar você a alcançar o sucesso!💡🌟\n"
	var MenuInteractionText = "*Olá! Por favor, escolha uma das seguintes opções de 1 a 4:*\n\n1. VER SEUS AGENDAMENTO ? 👁️\n2. VER HORÁRIOS DISPONÍVEIS ? 👀\n3. CANCELAR UM AGENDAMENTO ? ❌\n4. ENTRAR EM CONTATO ? 📞\n\n_Responda com o número correspondente à sua escolha._"
	var SchedulesAvailableTitle = "*Horários disponíveis para agendamento:*\n\n"
	var Contacts = "*Contatos* ☎️\n\n Sávio Picanço *+5522996043721*"

	// errors
	var ErrorDefault = "*Desculpe, algo deu errado. Por favor, tente novamente mais tarde.*"
	var SchedulesNotFound = "*Não encontramos nenhum agendamento para você.*"
	var NotUnderstand = "*Desculpe, não entendi, pode repetir.*"

	// layout
	var BackButton = "\n_0 - VOLTAR_  ◀️"
	var ScheduleOtherTime = "\n\n_1 - AGENDAR OUTRA DATA 📅_"
	var DefaultFooter = "\n\n_Responda com o número correspondente à sua escolha_"

	return Message{
		Greetings:               Greetings,
		MenuInteractionText:     MenuInteractionText,
		SchedulesAvailableTitle: SchedulesAvailableTitle,
		ErrorDefault:            ErrorDefault,
		SchedulesNotFound:       SchedulesNotFound,
		NotUnderstand:           NotUnderstand,
		BackButton:              BackButton,
		ScheduleOtherTime:       ScheduleOtherTime,
		DefaultFooter:           DefaultFooter,
		Contacts:                Contacts,
	}
}
