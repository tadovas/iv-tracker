package tax

import (
	"database/sql"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/tadovas/iv-tracker/expense"
	"github.com/tadovas/iv-tracker/income"
)

type CalculatorDBLoader struct {
	DB *sql.DB
}

func (cl CalculatorDBLoader) LoadFor(year income.Year) (YearlyCalculator, error) {

	var sb SocialBase
	row := cl.DB.QueryRow("SELECT vdu, count, percentage FROM social_base WHERE year = ?", year)
	if err := row.Scan(&sb.VDU, &sb.Count, &sb.Percentage); err != nil {
		return YearlyCalculator{}, fmt.Errorf("social base query error: %w", err)
	}

	vsdProvider := VSDProvider{
		SB: sb,
	}
	row = cl.DB.QueryRow("SELECT percentage FROM vsd_tax WHERE year=?", year)
	if err := row.Scan(&vsdProvider.Percentage); err != nil {
		return YearlyCalculator{}, fmt.Errorf("vsd tax query error: %w", err)
	}

	psdProvider := PSDProvider{
		SB: sb,
	}
	row = cl.DB.QueryRow("SELECT percentage FROM psd_tax WHERE year=?", year)
	if err := row.Scan(&psdProvider.Percentage); err != nil {
		return YearlyCalculator{}, fmt.Errorf("vsd tax query error: %w", err)
	}

	var gpmProvider GPMProvider
	row = cl.DB.QueryRow("SELECT percentage FROM gpm_tax WHERE year = ?", year)
	if err := row.Scan(&gpmProvider.Percentage); err != nil {
		return YearlyCalculator{}, fmt.Errorf("gpm query error: %w", err)
	}

	return FromProviders(
		expense.FixedPercentage{
			Percentage: decimal.NewFromFloat(0.3),
		},
		gpmProvider, vsdProvider, psdProvider,
	), nil
}
