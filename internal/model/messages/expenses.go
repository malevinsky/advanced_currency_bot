package messages


import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.ozon.dev/amalevinskaya/teodora-malevinskaia/internal/storage"
	"io"
	"io/ioutil"
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

// Parse and store expense for given user
func AddExpense(id int64, message string) error {
	expense, err := parseExpense(message) //

	if err != nil {
		return err
	}

	storage.AddExpense(id, expense)
	return nil
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
	fmt.Print(parts)


	if len(parts) != 4 { // Проверяем, что передаётся верное количество аргументов — 3.
		_ = fmt.Sprint("Ошибка: введите три параметра.")
		return nil, nil
	}
	MainCurr = parts[0]
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
	var i int = int(currency)
	s2 := strconv.Itoa(i)
	fmt.Println(currency)

	if parts[0] == "RUB" {
		textgreting := "Трата записана:\n- Категория: " + parts[2] + "\n- Сумма: " + parts[1] + " " + parts[0] + "\n- Дата: " + parts[3] + "\n\nДополнительные команды: \nsummary — сумма всех трат"
		Greting = textgreting
	} else {
		textgreting := "Трата записана:\n- Категория: " + parts[2] + "\n- Сумма: " + parts[1] + " " + parts[0] + "\n- Сумма в рублях: " + s2 + "\n- Дата: " + parts[3] + "\n\nДополнительные команды: \nsummary — сумма всех трат"
		Greting = textgreting
	}

	return storage.NewExpense(amount, parts[2], ts, amountfl), nil
}

func ValidCurr(currency string, amountfl float64) float64 {

	switch currency {
	case "USD":
		rubles := parseapi(1)
		return amountfl * rubles

	case "CNY":
		rubles := parseapi(2)
		return amountfl * rubles

	case "EUR":
		rubles := parseapi(3)
		return amountfl * rubles

	case "RUB":
		return amountfl
	}
	return 0
}

func parseapi(num int) float64{

	/**
	Функция parseapi нужна, чтобы достать данные из API exchangeratesapi
	1. Выше объявлены две структуры: Currency и Rates, в них запишутся данные из API.
	2. HTTP-запрос к веб-ресурсу отправляется через http.Client.
	3. Для создания объекта используется http.NewRequest().
	4. Для отправки объекта используется Do().
	5. Содержимое буфера вычитывается, чтобы записать в файл.
	*/

	url :="http://api.exchangeratesapi.io/v1/latest?access_key=9c484230306ca3014e2eb4c8575de8df&symbols=USD,CNY,RUB&format=1"

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

			}
		}(res.Body)
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	Currency1 := Currency{}
	jsonErr := json.Unmarshal(body, &Currency1)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	Currency1.Rates.EUR = 1

	switch num {
	case 1:
		converted := Currency1.Rates.RUB / Currency1.Rates.USD
		return converted //1

	case 2:
		converted := Currency1.Rates.RUB / Currency1.Rates.CNY
		return converted //7

	case 3:
		return Currency1.Rates.RUB //60

	}
	return 0
}
