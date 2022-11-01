package messages

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.ozon.dev/amalevinskaya/teodora-malevinskaia/internal/storage"
	"io"
	"log"
	"net/http"
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
	expense, err := parseExpense(message) //

	if err != nil {
		return err
	}

	storage.AddExpense(id, expense)
	return nil
}

func AddExpense(id int64, message string) error {
	expense, err := parseExpense(message) //

	if err != nil {
		return err
	}
	storage.AddExpense(id, expense)
	return nil
}

func parseCurrency(message string) (string, error) {
	normalizedMessage := strings.TrimSpace(strings.TrimPrefix(message, "/currency"))
	parts := strings.Split(normalizedMessage, " ")

	if len(parts) != 1 {
		return "Напишите валюту правильно. Например: /currency EUR", nil
	}
	if parts[0] == "" {
		return "Напишите валюту правильно. Например: /currency EUR", nil
	} else {
		currencyupper := strings.ToUpper(normalizedMessage) //всё привожу к большим буквам, чтобы не вылезла ошибка, если отправят EuR или rUb
		MainCurr = currencyupper
		return "Успешно установлена валюта: ", nil
	}
}

func parseExpense(message string) (*storage.Expense, error) {
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

	if len(parts) != 4 { // Проверяем, что передаётся верное количество аргументов — 3.
		_ = fmt.Sprint("Ошибка: введите три параметра.")
		return nil, nil
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
		return nil, errors.New("Сумма должна быть цифрой")
	}

	//const layout = "2021-11-22"
	ts, err := time.Parse("2006-01-02", parts[3])

	if err != nil {
		return nil, fmt.Errorf("Напишите дату в формате день-месяц-год")
	}

	currency := ValidCurr(parts[0], amountfl)
	MainCurr = parts[0]
	fmt.Print("parts[0]")
	fmt.Print(parts[0])

	s2 := fmt.Sprintf("%f", currency)
	//s2 := strconv.Itoa(int(currency))

	if parts[0] == "RUB" {
		textgreting := "Трата записана:\n- Категория: " + parts[2] + "\n- Сумма: " + parts[1] + " " + parts[0] + "\n- Дата: " + parts[3] + "\n\nПолучить сумму всех трат по датам и категориям: \n/get + year | week | day." + "\n\nВалюта, которую вы сейчас используете: " + parts[0]
		Greting = textgreting
	} else {
		textgreting := "Трата записана:\n- Категория: " + parts[2] + "\n- Сумма: " + parts[1] + " " + parts[0] + "\n- Сумма в рублях: " + s2 + "\n- Дата: " + parts[3] + "\n\nПолучить сумму всех трат по датам и категориям: \n/get + year | week | day." + "\n\nВалюта, которую вы сейчас используете: " + parts[0]
		Greting = textgreting
	}

	fmt.Print("\n\n")
	fmt.Print(currency)
	fmt.Print(int(currency))
	fmt.Print("\n\n")
	return storage.NewExpense(currency, parts[2], ts, amountfl), nil
}

func ValidCurr(currency string, amountfl float64) float64 {
	/**
	Получаем сумму в рублях
	 */
	switch currency {
	case "USD":
		rubles := parseapi(1)
		return amountfl * rubles

	case "CNY":
		rubles := parseapi(2)
		//fmt.Print("\n\n")
		//fmt.Print(rubles)
		//fmt.Print(amountfl * rubles)
		//fmt.Print("\n\n")

		return amountfl * rubles

	case "EUR":
		rubles := parseapi(3, )
		return amountfl * rubles

	case "RUB":
		return amountfl
	}
	return 0
}

func Parseapibeginning()  {
	/**
	Функция parseapi нужна, чтобы достать данные из API exchangeratesapi
	1. Выше объявлены две структуры: Currency и Rates, в них запишутся данные из API.
	2. HTTP-запрос к веб-ресурсу отправляется через http.Client.
	3. Для создания объекта используется http.NewRequest().
	4. Для отправки объекта используется Do().
	5. Содержимое буфера вычитывается, чтобы записать в файл.
	*/
	url := "http://api.exchangeratesapi.io/v1/latest?access_key=9c484230306ca3014e2eb4c8575de8df&symbols=USD,CNY,RUB&format=1"

	spaceClient := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatal(getErr)
			}
		}(res.Body)
	}
	body, readErr := io.ReadAll(res.Body)

	if readErr != nil {
		log.Fatal(readErr)
	}
	Currency1 := Currency{}
	jsonErr := json.Unmarshal(body, &Currency1)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	usd = Currency1.Rates.USD
	cny = Currency1.Rates.CNY
	rub = Currency1.Rates.RUB
	eur = float64(1)

	storage.CurrencyStorage2(usd, cny, rub, eur)
}

func parseapi(num int) float64 {

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
	Currency1 := Currency{}
	switch num {
	case 1:
		//converted := Currency{}
		converted := Currency1.Rates.RUB / Currency1.Rates.USD
		return converted

	case 2:
		converted := Currency1.Rates.RUB / Currency1.Rates.CNY
		fmt.Print("\n currency diff")
		fmt.Print(Currency1.Rates.CNY)
		return converted

	case 3:
		return Currency1.Rates.RUB

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

func parsePeriod(message string) (*time.Time, error) {
	/**
	В этой функции происходят похожие процессы на обработку траты в expenses:
	1. Получаем команду в string, её нужно разбить на части.
	2. Обрабатываем части. В нашем случае она одна — week, month, year.
	3. Если пользователь записал что-то криво, вылезет ошибка с подсказкой, как исправить.
	*/
	normalizedMessage := strings.TrimSpace(strings.TrimPrefix(message, "/get"))
	parts := strings.Split(normalizedMessage, " ")

	if len(parts) != 1 {
		return nil, errors.New("Допишите период, за который нужно получить отчёт: week, month, year. Например, /get year.")
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

// Format expenses by category into table
func formatExpenses(expenses []*storage.Expense, rates []*storage.Rates) string {

	if len(expenses) == 0 {
		return "Вы пока не добавили трату."
	}
	/**
	Цель — достать значения и красиво их вывести. Достаё
	*/
	expensesByCategory := make(map[string]int)
	for _, expense := range expenses {
		//fmt.Println(reflect.TypeOf(expense))
		result := revert(expense)
		expensesByCategory[expense.Category] += int(result)
	}

	var formattedResult strings.Builder

	for category, amount := range expensesByCategory {
		formattedResult.WriteString(fmt.Sprintf("\n%s: %d", category, amount))
	}

	return formattedResult.String()
}

func revert(expense *storage.Expense) float64 {
	//Currency1 := Currency{}

	switch MainCurr {
	case "EUR":
		difference := rub / eur
		finalAmount := float64(expense.Amount) / difference
		return finalAmount

	case "USD":
		difference := rub / usd
		finalAmount := float64(expense.Amount) / difference
		return finalAmount

	case "CNY":
		difference := rub / cny
		finalAmount := expense.Amount / difference
		//fmt.Println(expense.Amount)
		//fmt.Println(cny)
		//fmt.Print(difference)
		return finalAmount
	}
	return expense.Amount
}
