package income

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/tadovas/iv-tracker/log"
)

type ID int

type Money float64

type CountryCode string

type Year int

type Income struct {
	ID      ID
	Amount  Money
	Date    time.Time
	Origin  CountryCode
	Comment string
}

type Repository struct {
	DB sql.DB
}

func (repo Repository) AddNewIncome(income Income) (ID, error) {
	stmt, err := repo.DB.Prepare("INSERT INTO income () VALUES()")
	if err != nil {
		return ID(0), fmt.Errorf("new income statement prepare error: %w", err)
	}
	defer log.IfError("new income statement close", stmt.Close)

	res, err := stmt.Exec(income.Amount, income.Date, income.Origin, income.Comment)
	if err != nil {
		return ID(0), fmt.Errorf("new income statement execution error: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return ID(0), fmt.Errorf("new income last insert id error: %w", err)
	}
	return ID(id), nil
}

func (repo Repository) ListIncomesByYear(year Year) ([]Income, error) {
	stmt, err := repo.DB.Prepare("")
	if err != nil {
		return nil, fmt.Errorf("list income statement prepare error: %w", err)
	}
	defer log.IfError("list income statement close", stmt.Close)

	rows, err := stmt.Query(year)
	if err != nil {
		return nil, fmt.Errorf("list income query error: %w", err)
	}
	var incomes []Income
	for rows.Next() {
		var income Income
		if err := rows.Scan(income.ID, income.Amount, income.Date, income.Origin, income.Comment); err != nil {
			return nil, fmt.Errorf("list income row scan error: %w", err)
		}
		incomes = append(incomes, income)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("list income row iterator error: %w", err)
	}
	return incomes, nil
}
