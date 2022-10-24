package messages

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gitlab.ozon.dev/amalevinskaya/teodora-malevinskaia/internal/storage"
)

const ReportPrefix = "/get"

func GetReport(userID int64, message string) (string, error) {
	start_period, err := parsePeriod(message)

	if err != nil {
		return "", err
	}

	expenses := storage.GetExpenses(userID, *start_period)
	return formatExpenses(expenses), nil
}

func parsePeriod(message string) (*time.Time, error) {
	normalizedMessage := strings.TrimSpace(strings.TrimPrefix(message, ReportPrefix))
	parts := strings.Split(normalizedMessage, " ")

	if len(parts) != 1 {
		return nil, errors.New("Report must consist of one part: period")
	}

	period := strings.ToLower(parts[0])
	now := time.Now()

	switch period {
	case "week":
		now = now.AddDate(0, 0, -7)
	case "month":
		now = now.AddDate(0, -1, 0)
	case "year":
		now = now.AddDate(-1, 0, 0)
	default:
		return nil, errors.New("Report must be one of: week, month, year")
	}

	return &now, nil
}

// Format expenses by category into table
func formatExpenses(expenses []*storage.Expense) string {

	if len(expenses) == 0 {
		return "No expenses are found"
	}

	expensesByCategory := make(map[string]int)
	for _, expense := range expenses {
		expensesByCategory[expense.Category] += expense.Amount
	}

	var formattedResult strings.Builder

	for category, amount := range expensesByCategory {
		formattedResult.WriteString(fmt.Sprintf("%s: %d\n", category, amount))
	}

	return formattedResult.String()
}

