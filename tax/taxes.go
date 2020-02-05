package tax

import "github.com/tadovas/iv-tracker/income"

type Type string

const (
	GPM = Type("GPM")
	PSD = Type("PSD")
	VSD = Type("VSD")
)

type Tax struct {
	Type   Type
	Amount income.Money
}

type Provider interface {
	Calculate(total income.Money) (Tax, error)
}

type List []Tax

func (tl List) TotalTaxAmount() income.Money {
	var total income.Money
	for _, taxVal := range tl {
		total = total.Add(taxVal.Amount)
	}
	return total
}
