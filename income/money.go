package income

import (
	"database/sql/driver"

	"github.com/shopspring/decimal"
)

type Money decimal.Decimal

func FromFloat(f float64) Money {
	return Money(decimal.NewFromFloat(f))
}

func (m Money) Multiply(d decimal.Decimal) Money {
	val := decimal.Decimal(m)
	return Money(val.Mul(d))
}

func (m Money) Sub(other Money) Money {
	val := decimal.Decimal(m)
	otherVal := decimal.Decimal(other)
	return Money(val.Sub(otherVal))
}

func (m Money) Add(amount Money) Money {
	val := decimal.Decimal(m)
	other := decimal.Decimal(amount)
	return Money(val.Add(other))
}

func Min(a Money, b Money) Money {
	aVal := decimal.Decimal(a)
	bVal := decimal.Decimal(b)
	return Money(decimal.Min(aVal, bVal))
}

func (m *Money) Scan(val interface{}) error {
	d := (*decimal.Decimal)(m)
	return d.Scan(val)
}

func (m Money) MarshalJSON() ([]byte, error) {
	d := (decimal.Decimal)(m).Round(2)
	return d.MarshalJSON()
}

func (m *Money) UnmarshalJSON(bytes []byte) error {
	d := (*decimal.Decimal)(m)
	return d.UnmarshalJSON(bytes)
}

func (m Money) Value() (driver.Value, error) {
	d := (decimal.Decimal)(m)
	return d.Value()
}

func (m Money) String() string {
	d := (decimal.Decimal)(m)
	return d.Round(2).String()
}
