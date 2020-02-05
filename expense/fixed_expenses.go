package expense

import (
	"github.com/shopspring/decimal"
	"github.com/tadovas/iv-tracker/income"
)

type FixedPercentage struct {
	Percentage decimal.Decimal
}

func (fp FixedPercentage) ExpenseFor(totalIncome income.Money) (income.Money, error) {
	return totalIncome.Multiply(fp.Percentage), nil
}
