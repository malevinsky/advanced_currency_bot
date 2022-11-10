package messages

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func parsePeriod(message string) (*time.Time, error) {
	/**
	В этой функции происходят похожие процессы на обработку траты в expenses:
	1. Получаем команду в string, её нужно разбить на части.
	2. Обрабатываем части. В нашем случае она одна — week, month, year.
	3. Если пользователь записал что-то криво, вылезет ошибка с подсказкой, как исправить.
	*/
	normalizedMessage := strings.TrimSpace(strings.TrimPrefix(message, "/get"))
	parts := strings.Split(normalizedMessage, " ")

	if len(parts) > 2 {
		return nil, fmt.Errorf("Слишком много аргументов. Допишите период, за который нужно получить отчёт: week, month, year. Например, /get year.")
	}

	if parts[0] == "" {
		return nil, fmt.Errorf("Допишите период, за который нужно получить отчёт: week, month, year. Например, /get year.")
	}

	period := strings.ToLower(parts[0]) //на всякий случай, если напишет YeAr, а то может быть ошибка

	/**
	1. Узнаём время на момент отправки сообщения.
	2. Нужно получить период Time, за который мы выводим результаты.
	О цифрах в Addtime:
	- -7 — это нынешняя дата минус 7 дней;
	- -1 — это нынешняя дата минус один месяц;
	- -1 — это нынешняя дата минус год.
	*/

	now := time.Now()
	switch period {
	case "week":
		now = now.AddDate(0, 0, -7)
	case "month":
		now = now.AddDate(0, -1, 0)
	case "year":
		now = now.AddDate(-1, 0, 0)
	default:
		return nil, errors.New("Неправильная команда. Я использую только week, month, year. Например, /get year.")
	}

	return &now, nil
}
