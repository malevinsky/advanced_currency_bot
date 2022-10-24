package messages

import (
	"strings"
)

type MessageSender interface {
	SendMessage(text string, userID int64) error
}

type MessageSender2 interface {
	SendMessage(text string, userID int64, Summary int64, Category string) error
}
type Model struct {
	tgClient MessageSender
}

func New(tgClient MessageSender) *Model {
	return &Model{
		tgClient: tgClient,
	}
}

type Message struct {
	Text      string
	UserID    int64
}


func (s *Model) IncomingMessage(msg Message) error {
	response := "Unknown command"


	greetings := "Привет! Это бот для учёта трат. Использование бота:\nДобавьте трату командой /add + валюта + сумма + категория + дата."
	/**
	При нажатии start отображается приветствие
	 */
	if msg.Text == "/start" {
		//greetings := "Variable string %d content",

		return s.tgClient.SendMessage(greetings, msg.UserID)
	}

	/**
	При нажатии /currency
	 */
	if msg.Text == "/currency" {
		//добавить отдельную функцию для добавления валюты
		//return s.tgClient.SendMessage("hello", msg.UserID)
	}

	if strings.HasPrefix(msg.Text, "/add_expense") {
		err := AddExpense(msg.UserID, msg.Text)
		if err != nil {
			response = err.Error()
		} else {
			response = "Expense added successfully"
		}
	} else if strings.HasPrefix(msg.Text, "/get") {
		report, err := GetReport(msg.UserID, msg.Text)
		if err != nil {
			response = err.Error()
		} else {
			response = report
		}
	}
	response = Greting

	return s.tgClient.SendMessage(response, msg.UserID)
}
