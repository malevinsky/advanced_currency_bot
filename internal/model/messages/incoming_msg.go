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

/**
В сообщении от пользователя хранится две сущности:
- сам текст,
- идентификатор отправителя.
*/
type Message struct {
	Text      string
	UserID    int64
}


func (s *Model) IncomingMessage(msg Message) error {
	/**
	Пользователь что-нибудь отправляет, далее проверяем команду.
	Сценария четыре:
	- /start — приветствие и подсказка о работе бота;
	- /add — добавление траты и смена основной валюты;
	- /get — получение статистики по тратам.

	Хотелось добавить кнопки, но почему-то не вышло.
	 */
	response := "Что-то не то, отправьте правильную команду"

	if msg.Text == "/start" {
		greetings := "Привет! Это бот для учёта трат. Использование бота:\nДобавьте трату командой /add + валюта + сумма + категория + дата."
		return s.tgClient.SendMessage(greetings, msg.UserID)
	}

	if strings.HasPrefix(msg.Text, "/add") {
		err := AddExpense(msg.UserID, msg.Text) // -> файл expenses.go
		if err != nil {
			response = err.Error()
		} else {
			response = Greting 	// В expenses.go есть переменная, в которую я записываю данные траты.
								// Это позволяет посылать подтверждения с данными по каждой трате.
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
