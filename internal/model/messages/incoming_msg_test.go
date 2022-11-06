package messages

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "gitlab.ozon.dev/amalevinskaya/teodora-malevinskaia/internal/mocks/messages"
)

func Test_OnStartCommand_ShouldAnswerWithIntroMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	model := New(sender)

	sender.EXPECT().SendMessage("Привет! Это бот для учёта трат. Использование бота:\nДобавьте трату командой /add + валюта + сумма + категория + дата.", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/start",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnUnknownCommand_ShouldAnswerWithHelpMessage(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	sender.EXPECT().SendMessage("Что-то не то, отправьте правильную команду", int64(123))
	model := New(sender)

	err := model.IncomingMessage(Message{
		Text:   "some text",
		UserID: 123,
	})

	assert.NoError(t, err)
}



func Test_UnknownCurrency(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	model := New(sender)

	sender.EXPECT().SendMessage("Введите правильную валюту. Например, «/currency CNY»", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/currency приветики",
		UserID: 123,
	})
	assert.NoError(t, err)
}

func Test_CurrencyWithNoCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	model := New(sender)

	sender.EXPECT().SendMessage("Вы забыли ввести валюту. Бот поддерживает RUB, EUR, CNY и USD. Пример правильной команды: «/currency CNY»", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/currency",
		UserID: 123,
	})

	assert.NoError(t, err)
}


func Test_CurrencyWithTooManyCommands(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	model := New(sender)

	sender.EXPECT().SendMessage("Вы забыли ввести валюту. Бот поддерживает RUB, EUR, CNY и USD. Пример правильной команды: «/currency CNY»", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/currency привет бот как дела",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_AddNotFourElements(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	model := New(sender)

	sender.EXPECT().SendMessage("Ошибка: введите четыре параметра.", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/add EUR, 50, Еда",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_AmountIsNotANumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	model := New(sender)

	sender.EXPECT().SendMessage("Ошибка: сумма должна быть цифрой", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/add EUR, привет, Еда, 2022-10-10",
		UserID: 123,
	})
	assert.NoError(t, err)
}

func Test_WrongDataFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	model := New(sender)

	sender.EXPECT().SendMessage("Ошибка: напишите дату в формате день-месяц-год", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/add EUR, 50, Еда, первое апреля",
		UserID: 123,
	})
	assert.NoError(t, err)
}

func Test_NoDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	model := New(sender)

	sender.EXPECT().SendMessage("Ошибка: введите четыре параметра.", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/add",
		UserID: 123,
	})
	assert.NoError(t, err)
}



func Test_GetWithoutArguments(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	model := New(sender)

	sender.EXPECT().SendMessage("Допишите период, за который нужно получить отчёт: week, month, year. Например, /get year.", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/get",
		UserID: 123,
	})
	assert.NoError(t, err)
}

func Test_TooManyArguments(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	model := New(sender)

	sender.EXPECT().SendMessage("Слишком много аргументов. Допишите период, за который нужно получить отчёт: week, month, year. Например, /get year.", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/get много аргументов тут",
		UserID: 123,
	})
	assert.NoError(t, err)
}

func Test_UnknownParameter(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	model := New(sender)

	sender.EXPECT().SendMessage("Неправильная команда. Я использую только week, month, year. Например, /get year.", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/get привет",
		UserID: 123,
	})
	assert.NoError(t, err)
}

