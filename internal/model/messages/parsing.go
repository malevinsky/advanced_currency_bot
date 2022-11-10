package messages

import (
	"encoding/json"
	"fmt"
	"gitlab.ozon.dev/amalevinskaya/teodora-malevinskaia/internal/storage"
	"io"
	"log"
	"net/http"
	"time"
)

func Parseapibeginning() {
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
		return
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
