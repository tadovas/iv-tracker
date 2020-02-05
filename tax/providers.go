package tax

import (
	"github.com/shopspring/decimal"
	"github.com/tadovas/iv-tracker/income"
)

type GPMProvider struct {
	Percentage decimal.Decimal
}

func (gp GPMProvider) Calculate(total income.Money) (Tax, error) {
	return Tax{
		Amount: total.Multiply(gp.Percentage),
		Type:   GPM,
	}, nil
}

var _ Provider = GPMProvider{}

type SocialBase struct {
	VDU        income.Money
	Count      decimal.Decimal
	Percentage decimal.Decimal
}

func (sb SocialBase) Calculate(amount income.Money) (income.Money, error) {
	socialCeil := sb.VDU.Multiply(sb.Count)
	incomeBase := amount.Multiply(sb.Percentage)
	return income.Min(socialCeil, incomeBase), nil
}

type PSDProvider struct {
	Percentage decimal.Decimal
	SB         SocialBase
}

func (pp PSDProvider) Calculate(total income.Money) (Tax, error) {
	baseAmount, err := pp.SB.Calculate(total)
	if err != nil {
		return Tax{}, err
	}
	return Tax{
		Amount: baseAmount.Multiply(pp.Percentage),
		Type:   PSD,
	}, nil
}

type VSDProvider struct {
	Percentage decimal.Decimal
	SB         SocialBase
}

func (vp VSDProvider) Calculate(total income.Money) (Tax, error) {
	baseAmount, err := vp.SB.Calculate(total)
	if err != nil {
		return Tax{}, err
	}
	return Tax{
		Amount: baseAmount.Multiply(vp.Percentage),
		Type:   VSD,
	}, nil
}

var _ Provider = VSDProvider{}
