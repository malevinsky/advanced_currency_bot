package storage

import (
	"net/http"
	"time"
)

type Repository struct {
	client *http.Client
}

func NewRepository(client *http.Client) *Repository {
	return &Repository{
		client:   client,
	}
}

func GetRepository() []*Repository {
	var result []*Repository

	return result
}
var userStorage map[int64]User = make(map[int64]User)

type userStorage2 struct {
	Id	int64
	User	[]*User
}

type User struct {
	ID       int64
	Expenses []*Expense
}

type Expense struct {
	Amount   float64
	Category string
	Ts       time.Time
	Rubles   float64
}

func NewExpense(amount float64, category string, ts time.Time, rubles float64) *Expense {
	/**
	Добавляем передаваемые значения в структуру Expense.
	*/
	return &Expense{
		Amount:   amount,
		Category: category,
		Ts:       ts,
		Rubles:   rubles,
	}
}

type Currency struct {
	Success   bool
	Timestamp int
	Base      string
	Date      string
	Rates     Rates
}

type Rates struct {
	USD float64
	CNY float64
	RUB float64
	EUR float64
}

func GetRates() []*Rates {
	var result []*Rates

	return result
}

func CurrencyStorage2(usd float64, cny float64, rub float64, eur float64) *Rates {
	return &Rates{
		USD: usd,
		CNY: cny,
		RUB: rub,
		EUR: eur,
	}
}



func AddExpense(userID int64, expense *Expense) {
	user, ok := userStorage[userID]
	if !ok {
		user = User{ID: userID}
	}
	user.addExpense(expense)
	userStorage[user.ID] = user
}

func (u *User) addExpense(e *Expense) {
	u.Expenses = append(u.Expenses, e)
}

func GetExpenses(userID int64, start_period time.Time) []*Expense {
	user, ok := userStorage[userID]
	if !ok {
		return nil
	}

	var result []*Expense

	for _, expense := range user.Expenses {
		if expense.Ts.After(start_period) {
			result = append(result, expense)
		}
	}

	return result
}

func Get1(userID int64, start_period time.Time) []*Expense {
	user, ok := userStorage[userID]
	if !ok {
		return nil
	}

	var result []*Expense

	for _, expense := range user.Expenses {
		if expense.Ts.After(start_period) {
			result = append(result, expense)
		}
	}

	return result
}
