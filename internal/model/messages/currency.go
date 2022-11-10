package messages

import (
	"strings"
)

func parseCurrency(message string) (string, error) {
	//normalizedMessage := strings.TrimSpace(strings.TrimPrefix(message, "/currency"))
	parts := strings.Split(message, " ")

	if len(parts) != 2 {
		return "Вы забыли ввести валюту. Бот поддерживает RUB, EUR, CNY и USD. Пример правильной команды: «/currency CNY»", nil
	} else {
		currencyupper := strings.ToUpper(parts[1]) //всё привожу к большим буквам, чтобы не вылезла ошибка, если отправят EuR или rUb

		if currencyupper == "RUB" || currencyupper == "EUR" || currencyupper == "CNY" || currencyupper == "USD" {
			//проверка, чтобы не записать в валюту тарабарщину
			MainCurr = currencyupper
			return "Успешно установлена валюта: ", nil
		} else {
			return "Введите правильную валюту. Например, «/currency CNY»", nil
		}
	}
}
