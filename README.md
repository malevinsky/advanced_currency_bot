
## Структура проекта
### cmd/bot
`main.go` — входная точка в проект

- `internal/clients`
- `internal/config` — обработка конфиг-файла
- `internal/mocks` —
- `internal/model` — файл incoming_msg.go, в котором обрабатываются
- `internal/storage` — структуры, в которых хранятся данные

-----------

## Домашняя работа 1
Нужно добавить функционал:
- Команда добавления новой финансовой "траты". В трате должна присутствовать сумма, категория и дата. Но можете добавить еще поля, если считаете нужным. Придумайте, как оформить команду так, чтобы пользователю было удобно ее использовать.
- Хранение трат в памяти, базы данных пока не используем.
- Команда запроса отчета за последнюю неделю/месяц/год. В отчете должны быть суммы трат по категориям.

## Домашняя работа 2
- Команда переключения бота на конкретную валюту - "выбрать валюту"
- После ввода команды бот предлагает выбрать интересующую валюту из четырех: USD, CNY, EUR, RUB
- При нажатии на нужную валюту переключаем бота на нее - результат получение трат конвертируется в выбранную валюту.
- Храним траты всегда в рублях, конвертацию используем только для отображения, ввода и отчетов
                     
Особенности:
- При запуске сервиса мы в отдельном потоке запрашиваем курсы валют.
- Запрос курса валют происходит из любого из открытых источников.
- Сервис должен завершаться gracefully.

--------

API для получения валют онлайн: `http://api.exchangeratesapi.io/v1/latest?access_key=9c484230306ca3014e2eb4c8575de8df&symbols=USD,CNY,RUB&format=1`

### Как высчитывается валюта
**Base** — базовая валюта, от которой высчитываются значения остальных. 

Например, мы знаем, что: 
1 USD = 0.87 EUR, 
1 USD = 0.73 GBP, 
Чтобы узнать, сколько EUR равняется GBP?

Формула: 
```
a / b = c
```

```
EURGBP = (USDGBP / USDEUR) = (0.73 / 0.87) = 0.84
```

-------


![image info](./img/img.png)

