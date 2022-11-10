package messages

import (
	"errors"
	"fmt"
	"gitlab.ozon.dev/amalevinskaya/teodora-malevinskaia/internal/storage"
	//"reflect"
	"strconv"
	"strings"
	"time"
)

/**
MainCurr — валюта, на которую переключён бот. Для удобства я сделала так, чтобы её указывали прям при записи траты.
Greeting — string, в который я записываю данные из команды, чтобы наглядно выводить отчэт о записанной траты.
*/
var MainCurr = ""
var Greting = ""

var usd = 6.6
var cny = 6.6
var rub = 6.6
var eur = 6.6

var limit = 6.6

/**
Currency и Rates — это структуры, куда записываются данные из API. Дополнительно EUR — его значение «1»,
потому что это базовая валюта, от которой высчитываются значения остальных.
*/
type Currency struct {
	Success   bool   `json:"success"`
	Timestamp int    `json:"timestamp"`
	Base      string `json:"base"`
	Date      string `json:"date"`
	Rates     Rates  `json:"rates"`
}

type Rates struct {
	USD float64 `json:"USD"`
	CNY float64 `json:"CNY"`
	RUB float64 `json:"RUB"`
	EUR float64
}

const ExpensesPrefix = "/add"

func AddCurrency(id int64, message string) error {
	expense, err := parseExpense(message, id) //

	if err != nil {
		return err
	}

	storage.AddExpense(id, expense)
	return nil
}

func AddExpense(id int64, message string) error {
	expense, err := parseExpense(message, id) //

	if err != nil {
		return err
	}
	storage.AddExpense(id, expense)
	return nil
}


func parseExpense(message string, id int64) (*storage.Expense, error) {
	/**
	У нас есть строка с параметрами, которую отправил пользователь.
	1. Чтобы обрабатывать каждый аргумент отдельно, их нужно поделить с помощью метода Split. Разграничитель — запятая.
	2. Проверяем, что передаётся верное количество аргументов — 4.

	Например:
	- parts[0] — RUB
	- parts[1] — 500
	- parts[2] — еда
	- parts[0] — 2022-01-01
	*/
	normalizedMessage := strings.TrimSpace(strings.TrimPrefix(message, ExpensesPrefix))
	parts := strings.Split(normalizedMessage, ", ")

	if len(parts) != 5 { // Проверяем, что передаётся верное количество аргументов — 4.
		return nil, errors.New("Ошибка: введите четыре параметра.")
	}

	if parts[0] != "RUB" && parts[0] != "EUR" && parts[0] != "CNY" && parts[0] != "USD" {
		//проверка, чтобы не записать в валюту тарабарщину
		return nil, fmt.Errorf("Ошибка: неправильная валюта. Бот поддерживает RUB, EUR, CNY и USD.")
	}
	/**
	Сейчас все элементы в строке — String. Чтобы траты можно было складывать, нужно перевести её в Int.
	1. Преобразуем сумму траты из String в Int с помощью Atoi.
	2. Преобразуем время траты из String в Time с помощью time.Parse.
	3. Отправляем данные в нужном формате в функцию NewExpense, которая возвращает значения в структурку.
	*/
	amount, err := strconv.Atoi(parts[1])
	amountfl := float64(amount)

	if err != nil {
		return nil, errors.New("Ошибка: сумма должна быть цифрой")
	}

	amountLimit, err := strconv.Atoi(parts[1])
	amountLimitFloat := float64(amountLimit)

	if err != nil {
		return nil, errors.New("Ошибка: сумма лимита на месяц должна быть цифрой")
	}

	//const layout = "2021-11-22"
	ts, err := time.Parse("2006-01-02", parts[3])

	if err != nil {
		return nil, fmt.Errorf("Ошибка: напишите дату в формате день-месяц-год")
	}

	fmt.Print(ts)
	limit = ValidCurr(parts[0], amountLimitFloat)
	currency := ValidCurr(parts[0], amountfl)

	MainCurr = parts[0]
	//fmt.Print("parts[0]")
	//fmt.Print(parts[0])

	s2 := fmt.Sprintf("%f", currency)
	//s2 := strconv.Itoa(int(currency))

	if parts[0] == "RUB" {
		textgreting := "Трата записана:\n- Категория: " + parts[2] + "\n- Сумма: " + parts[1] + " " + parts[0] + "\n- Дата: " + parts[3] + "\n\nПолучить сумму всех трат по датам и категориям: \n/get + year | week | day." + "\n\nВалюта, которую вы сейчас используете: " + parts[0] + "\nТакже мы учли, что вы установили лимит трат на месяц: " + parts[4]
		Greting = textgreting
	} else {
		textgreting := "Трата записана:\n- Категория: " + parts[2] + "\n- Сумма: " + parts[1] + " " + parts[0] + "\n- Сумма в рублях: " + s2 + "\n- Дата: " + parts[3] + "\n\nПолучить сумму всех трат по датам и категориям: \n/get + year | week | day." + "\n\nВалюта, которую вы сейчас используете: " + parts[0] + "\nТакже мы учли, что вы установили лимит трат на месяц: " + parts[4]
		Greting = textgreting
	}


	start_period := time.Now()
	start_period = start_period.AddDate(0, -1, 0)

	expenses := storage.GetExpenses(id, start_period)
	rates := storage.GetRates()
	value, err4 := formatLimit(expenses, rates, limit, id)
	//fmt.Print(value)
	return value, err4
}


func ValidCurr(currency string, amountfl float64) float64 {
	/**
	Получаем сумму в рублях
	*/

	switch currency {
	case "USD":
		rubles := parseapi("USD")
		return amountfl * rubles

	case "CNY":
		rubles := parseapi("CNY")
		return amountfl * rubles

	case "EUR":
		rubles := parseapi("EUR")
		return amountfl * rubles

	case "RUB":
		return amountfl
	}
	return 0
}

func parseapi(num string) float64 {

	/**
	Что происходит ниже:
	Например, мы знаем, что:
	1 USD = 0.87 EUR,
	1 USD = 0.73 GBP,
	Чтобы узнать, сколько EUR равняется GBP, используется формула: a / b = c
	Тогда: EURGBP = (USDGBP / USDEUR) = (0.73 / 0.87) = 0.84

	Ниже я высчитываю курс валюты к валюте с помощью формулы выше. Лучший вариант — найти апи,
	где базовая валюьа — рубль, но я взяла такой.
	*/

	//expenses := storage.GetRates()
	//fmt.Print(expenses)
	Rates4 := storage.Rates{}
	fmt.Print(Rates4)
	switch num {
	case "USD":

		converted := rub / usd
		fmt.Print("attention!")
		//fmt.Print(expenses.RUB)
		return converted

	case "CNY":
		converted := rub / cny
		fmt.Print("\n currency diff")
		fmt.Print(cny)
		return converted

	case "EUR":
		return rub
	}
	return 0
}

func GetReport(userID int64, message string) (string, error) {

	if strings.HasPrefix(message, "/currency") {
		answer, err := parseCurrency(message)

		if err != nil {
			return "", err
		}
		return answer, nil

	} else {
		start_period, err := parsePeriod(message)

		if err != nil {
			return "", err
		}
		expenses := storage.GetExpenses(userID, *start_period)
		rates := storage.GetRates()
		fmt.Print(rates)
		return formatExpenses(expenses, rates), nil
	}
}
