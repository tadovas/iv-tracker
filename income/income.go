package income

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/tadovas/iv-tracker/log"
)

type ID int

type CountryCode string

type Year int

type Income struct {
	ID      ID
	Amount  Money
	Date    time.Time
	Origin  CountryCode
	Comment string
}

type IncomeList []Income

func (il IncomeList) Total() Money {
	var res Money
	for _, income := range il {
		res = res.Add(income.Amount)
	}
	return res
}

type Repository struct {
	DB *sql.DB
}

func (repo Repository) AddNewIncome(income Income) (ID, error) {
	stmt, err := repo.DB.Prepare("INSERT INTO incomes (amount,earned,year,origin,comment) VALUES(?,?,YEAR(?),?,?)")
	if err != nil {
		return ID(0), fmt.Errorf("new income statement prepare error: %w", err)
	}
	defer log.IfError("new income statement close", stmt.Close)

	res, err := stmt.Exec(income.Amount, income.Date, income.Date, income.Origin, income.Comment)
	if err != nil {
		return ID(0), fmt.Errorf("new income statement execution error: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return ID(0), fmt.Errorf("new income last insert id error: %w", err)
	}
	return ID(id), nil
}

func (repo Repository) ListIncomesByYear(year Year) (IncomeList, error) {

	rows, err := repo.DB.Query("SELECT id, amount, earned,origin, comment FROM incomes WHERE year=? ORDER BY earned DESC", year)
	if err != nil {
		return nil, fmt.Errorf("list income query error: %w", err)
	}
	defer log.IfError("list incomes rows close", rows.Close)

	var incomeList IncomeList
	for rows.Next() {
		var income Income
		if err := rows.Scan(&income.ID, &income.Amount, &income.Date, &income.Origin, &income.Comment); err != nil {
			return nil, fmt.Errorf("list income row scan error: %w", err)
		}
		incomeList = append(incomeList, income)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("list income row iterator error: %w", err)
	}
	return incomeList, nil
}

type YearSummary struct {
	Year         Year
	TotalAmount  Money
	TotalIncomes int
}

func (repo Repository) ListAllYears() ([]YearSummary, error) {
	rows, err := repo.DB.Query("select year, sum(amount), count(id) from incomes group by year")
	if err != nil {
		return nil, fmt.Errorf("list all years query error: %w", err)
	}
	defer log.IfError("list all years sql rows close", rows.Close)

	var years []YearSummary
	for rows.Next() {
		var yearSummary YearSummary
		if err := rows.Scan(&yearSummary.Year, &yearSummary.TotalAmount, &yearSummary.TotalIncomes); err != nil {
			return nil, fmt.Errorf("list all years row scan error: %w", err)
		}
		years = append(years, yearSummary)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("list all years query iterator error: %w", rows.Err())
	}

	return years, nil
}
