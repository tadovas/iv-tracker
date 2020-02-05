package tax

import (
	"github.com/tadovas/iv-tracker/income"
)

type ExpenseProvider interface {
	ExpenseFor(totalAmount income.Money) (income.Money, error)
}

type YearlyCalculator struct {
	expenses  ExpenseProvider
	providers []Provider
}

func (yc YearlyCalculator) Taxes(total income.Money) (List, error) {
	expenses, err := yc.expenses.ExpenseFor(total)
	if err != nil {
		return nil, err
	}

	profit := total.Sub(expenses)

	var taxes List
	for _, provider := range yc.providers {
		taxVal, err := provider.Calculate(profit)
		if err != nil {
			return nil, err
		}
		taxes = append(taxes, taxVal)
	}
	return taxes, nil
}

func FromProviders(expenses ExpenseProvider, providers ...Provider) YearlyCalculator {
	return YearlyCalculator{
		expenses:  expenses,
		providers: providers,
	}
}
