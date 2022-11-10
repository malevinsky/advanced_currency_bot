package messages

import (
	"errors"
	"fmt"
	"gitlab.ozon.dev/amalevinskaya/teodora-malevinskaia/internal/storage"
	"strconv"
	"strings"
)

func formatLimit(expenses []*storage.Expense, rates []*storage.Rates, limit float64, id int64) (*storage.Expense, error) {
	/**
	Цель — достать значения и красиво их вывести. Достаё
	*/

	expensesByCategory := make(map[string]int)
	for _, expense := range expenses {
		result := revert(expense)
		expensesByCategory[expense.Category] += int(result)
	}

	var formattedResult strings.Builder

	var s int
	for amount := range expensesByCategory {
		intVar, err := strconv.Atoi(amount)
		if err != nil {
			errors.New("Ошибка: сумма лимита на месяц должна быть цифрой")
		}
		s += intVar
		fmt.Println(formattedResult)
		//formattedResult.WriteString(fmt.Sprintf("\n%s: %d", category, amount))

		formattedResult.WriteString(fmt.Sprintf(string(s)))
	}

	if s < int(limit) {
		return nil, fmt.Errorf("Вы превысили лимит")

	} else {
		return nil, nil
	}
}

func formatExpenses(expenses []*storage.Expense, rates []*storage.Rates) string {

	if len(expenses) == 0 {
		return "Вы пока не добавили трату."
	}
	/**
	Цель — достать значения и красиво их вывести. Достаё
	*/

	expensesByCategory := make(map[string]int)
	for _, expense := range expenses {
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
		fmt.Println(expense.Amount)
		fmt.Println(cny)
		fmt.Print(difference)
		return finalAmount
	}
	return expense.Amount
}
