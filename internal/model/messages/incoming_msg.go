/**
Продолжаем работу над ботом, каркас которого создали на воркшопе.
Нужно добавить функционал:
- Команда добавления новой финансовой "траты". В трате должна присутствовать сумма, категория и дата. Но можете добавить еще поля, если считаете нужным. Придумайте, как оформить команду так, чтобы пользователю было удобно ее использовать.
- Хранение трат в памяти, базы данных пока не используем.
- Команда запроса отчета за последнюю неделю/месяц/год. В отчете должны быть суммы трат по категориям.
*/

package messages

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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
	category  string
	timecount time.Time
}

type wallet map[string]float64

var db = map[int64]wallet{}
var sum float64 = 0

func (s *Model) IncomingMessage(msg Message) error {
	command := strings.Split(msg.Text, ", ")

	current := time.Now()
	timecurrent := current.Format("01-02-2006")

	//println(len(command))
	if len(command) == 4 {
		amount, err := strconv.ParseFloat(command[2], 64)

		if err != nil {
			return s.tgClient.SendMessage("Ошибка", msg.UserID)
		}

		if _, ok := db[msg.UserID]; !ok {
			db[msg.UserID] = wallet{}
		}

		db[msg.UserID][command[1]] += amount
		//balanceText := fmt.Sprintf("%f\n", db[msg.UserID][command[1]])
		sumstr := strconv.FormatFloat(amount, 'f', -1, 64)
		//donestr := fmt.Sprintf("%f", amount)

		switch command[0] {
		case "/add":

			msg.category = command[3]
			sum = db[msg.UserID][command[1]]
			done := "Трата записана.\n- Категория: " + msg.category + "\n- Сумма: " + sumstr + " рублей\n- Дата: " + timecurrent + "\n\nДополнительные команды: \n summary — сумма всех трат"
			return s.tgClient.SendMessage(done, msg.UserID)

		}

	} else {
		switch command[0] {
		case "/start":
			user1 := msg.UserID
			fmt.Println(user1)
			return s.tgClient.SendMessage("Привет! Добавьте трату в формате «/add, название траты, сумма, категория».", msg.UserID)

		case "/summary":

			sumstr := strconv.FormatFloat(sum, 'f', -1, 64)

			return s.tgClient.SendMessage(" Сумма всех покупок: "+sumstr+" рублей ", msg.UserID)
		}
	}

	return s.tgClient.SendMessage("не знаю эту команду", msg.UserID)
}
