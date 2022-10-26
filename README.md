
## Структура проекта

```plain
├── cmd/
│   └── bot/
│       └── main.go					# Входная точка в проект
│
├── internal/
│   ├── clients/             #
│   ├── config/              # Обработка конфиг-файла
│   ├── mocks/     			 # Файл из шаблона, подключающий фреймворк Gomock для тестирования, но в этом проекте тестов нет.
│   ├── model/  
│ 	│ 	└── messages/
│ 	│		└── incoming_msg.go    #     
│ 	│		└── expenses.go 		
│	│		└── report.go
│   │
│   └── storage          			#package, в котором хранятся данные по тратам
│       └── user.go		
│ 
```

## Как пользоваться ботом

Бот доступен по логину в телеграме: @Ozon_theo_bot.

#### Методы

| Команда  | Примеры  | Описание |
|------------- |---------------| -------------|
| `/start`      | — | Подсказка и приветствие, начало работы с ботом. </br></br> ![image info](./img/start.jpg) |
| `/currency` | `/currency EUR` | Добавление основной валюты, в которой рассчитываются траты. </br></br> ![image info](./img/currency.jpg) </br></br> Примечание: в памяти траты хранятся в рублях, но перед выдачей я перевожу его в выбранную пользователем валюту.
| `/add`      | `/add EUR, 50, Еда, 2022-10-10` | Добавление основной валюты и траты. Строгий формат: `/add ВАЛЮТА(USD, CNY, EUR, RUB), СУММА, КАТЕГОРИЯ, ДАТА-В-ФОРМАТЕ-2022-10-10` </br></br> ![image info](./img/output.jpg) |
| `/get` | `/get year` | Получение статистики по тратам за год, месяц, неделю. Формат: `/get СРОК(year, month, week)` </br></br> ![image info](./img/get.jpg) |

--------
## Откуда берём валюты

Бот работает с 4 валютами: USD, CNY, EUR, RUB.

API для получения валют онлайн: `http://api.exchangeratesapi.io/v1/latest?access_key=9c484230306ca3014e2eb4c8575de8df&symbols=USD,CNY,RUB&format=1`

Всё, что с этим связано, лежит в `model/messages/expenses`. Примечание к иллюстрации: в структуре есть в EUR, он равен единице.
![image info](./img/img.png)

### Как высчитывается валюта
**Base** — базовая валюта, от которой высчитываются значения остальных. В апи, который я использую, это EUR.

Например, мы знаем, что: 
1 USD = 0.87 EUR, 
1 USD = 0.73 GBP, 
Чтобы узнать, сколько EUR равняется GBP, используется формула:

```
a / b = c
```

```
EURGBP = (USDGBP / USDEUR) = (0.73 / 0.87) = 0.84
```

-------

## Описание заданий
### Домашняя работа 1
Нужно добавить функционал:
| Условие  | Как я это сделала  |
|------------- |---------------| 
|✅ Команда добавления новой финансовой "траты". В трате должна присутствовать сумма, категория и дата. Но можете добавить еще поля, если считаете нужным. Придумайте, как оформить команду так, чтобы пользователю было удобно ее использовать. | Команда реализована — Строгий формат: `/add ВАЛЮТА(USD, CNY, EUR, RUB), СУММА, КАТЕГОРИЯ, ДАТА-В-ФОРМАТЕ-2022-10-10. |
|✅ Хранение трат в памяти, базы данных пока не используем. | Реализовано, я храню данные в структурах и переменных. Большая часть из них в storage |
|✅ Команда запроса отчета за последнюю неделю/месяц/год. В отчете должны быть суммы трат по категориям. | Команда реализована — это `/get ПЕРИОД`. Вывод категорий, трат и периодов тоже работает как описано в тз. |

### Домашняя работа 2
| Условие  | Как я это сделала  | 
|------------- |---------------| 
|✅ Команда переключения бота на конкретную валюту - "выбрать валюту" | Команда реализована — это `/currency ВАЛЮТА` |
|✅ После ввода команды бот предлагает выбрать интересующую валюту из четырех: USD, CNY, EUR, RUB | Все 4 валюты поддерживаются, траты можно записывать в любой из них, также можно между ними переключаться. |
|✅ При нажатии на нужную валюту переключаем бота на нее - результат получение трат конвертируется в выбранную валюту. | Реализовано. Чтобы переключиться на другую валюту, используйте `/currency ВАЛЮТА` |
|✅ Храним траты всегда в рублях, конвертацию используем только для отображения, ввода и отчетов. | Я записываю траты в package `storage.Expense.Amount` в рублях. Дополнительно вывожу в оповещении о трате суммы и в выбранной валюте, и в рублях. |
|✅ При запуске сервиса мы в отдельном потоке запрашиваем курсы валют. | Реализовано, я использую горутину в функции парсинга валют из открытого апи. |
|✅ Запрос курса валют происходит из любого из открытых источников. | Я использую открытый апи http://api.exchangeratesapi.io/v1/latest?access_key=9c484230306ca3014e2eb4c8575de8df&symbols=USD,CNY,RUB&format=1 |
|✅ Сервис должен завершаться gracefully. | Если я правильно поняла описания в интернете, у меня оно именно так работает.  |
