package storage

import "time"


var userStorage map[int64]User = make(map[int64]User)

type User struct {
	ID       int64
	Expenses []*Expense
}

type Expense struct {
	Amount		int
	Category	string
	Ts			time.Time
	Rubles		float64
}

func NewExpense(amount int, category string, ts time.Time, rubles float64) *Expense {
	/**
	Добавляем передаваемые значения в структуру Expense.
	 */
	return &Expense{
		Amount:		amount,
		Category:	category,
		Ts:			ts,
		Rubles:		rubles,
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
