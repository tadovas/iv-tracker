package tax

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/tadovas/iv-tracker/expense"
	"github.com/tadovas/iv-tracker/income"
)

func Test2019TaxValues(t *testing.T) {
	sb := SocialBase{
		VDU:        income.FromFloat(1136.2),
		Count:      decimal.NewFromInt(43),
		Percentage: decimal.NewFromFloat(0.9),
	}
	taxes2019 := FromProviders(
		expense.FixedPercentage{Percentage: decimal.NewFromFloat(0.3)},
		GPMProvider{
			Percentage: decimal.NewFromFloat(0.15),
		},
		VSDProvider{
			SB:         sb,
			Percentage: decimal.NewFromFloat(0.1252),
		},
		PSDProvider{
			SB:         sb,
			Percentage: decimal.NewFromFloat(0.0698),
		},
	)

	totalIncome := income.FromFloat(200 * 1000)
	taxList, err := taxes2019.Taxes(totalIncome)
	assert.NoError(t, err)

	assert.Equal(t, List{
		Tax{
			Type:   GPM,
			Amount: income.Money(decimal.New(210, 2)),
		},
		Tax{
			Type:   VSD,
			Amount: income.FromFloat(6116.84632),
		},
		Tax{
			Type:   PSD,
			Amount: income.FromFloat(3410.19068),
		},
	}, taxList)

	taxList, err = taxes2019.Taxes(income.FromFloat(65223))
	assert.NoError(t, err)
	assert.Equal(t, List{
		Tax{
			Type:   GPM,
			Amount: income.FromFloat(6848.415),
		},
		Tax{
			Type:   VSD,
			Amount: income.FromFloat(5144.529348),
		},
		Tax{
			Type:   PSD,
			Amount: income.FromFloat(2868.116202),
		},
	}, taxList)
}

func Test2020TaxValues(t *testing.T) {
	sb := SocialBase{
		VDU:        income.FromFloat(1241.4),
		Count:      decimal.NewFromInt(43),
		Percentage: decimal.NewFromFloat(0.9),
	}
	taxes2019 := FromProviders(
		expense.FixedPercentage{Percentage: decimal.NewFromFloat(0.3)},
		GPMProvider{
			Percentage: decimal.NewFromFloat(0.15),
		},
		VSDProvider{
			SB:         sb,
			Percentage: decimal.NewFromFloat(0.1252),
		},
		PSDProvider{
			SB:         sb,
			Percentage: decimal.NewFromFloat(0.0698),
		},
	)

	totalIncome := income.FromFloat(200 * 1000)
	taxList, err := taxes2019.Taxes(totalIncome)
	assert.NoError(t, err)

	assert.Equal(t, List{
		Tax{
			Type:   GPM,
			Amount: income.Money(decimal.New(210, 2)),
		},
		Tax{
			Type:   VSD,
			Amount: income.FromFloat(6683.20104),
		},
		Tax{
			Type:   PSD,
			Amount: income.FromFloat(3725.93796),
		},
	}, taxList)

	taxList, err = taxes2019.Taxes(income.FromFloat(65223))
	assert.NoError(t, err)
	assert.Equal(t, List{
		Tax{
			Type:   GPM,
			Amount: income.FromFloat(6848.415),
		},
		Tax{
			Type:   VSD,
			Amount: income.FromFloat(5144.529348),
		},
		Tax{
			Type:   PSD,
			Amount: income.FromFloat(2868.116202),
		},
	}, taxList)

}
