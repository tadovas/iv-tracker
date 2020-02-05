package saving

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/tadovas/iv-tracker/income"
)

type Saving struct {
	Amount    income.Money
	CreatedAt time.Time `json:"created_at"`
	Comment   string
}

type List []Saving

func (sl List) TotalSaved() income.Money {
	var saved income.Money
	for _, saving := range sl {
		saved = saved.Add(saving.Amount)
	}
	return saved
}

type Repository struct {
	DB *sql.DB
}

func (r Repository) AddNewSaving(saving Saving) error {
	query := `
		INSERT INTO tax_savings (year, created_at, amount, comment)
			VALUES(YEAR(?), ? , ? , ?)
	`
	if _, err := r.DB.Exec(query, saving.CreatedAt, saving.CreatedAt, saving.Amount, saving.Comment); err != nil {
		return fmt.Errorf("new saving insert error: %v", err)
	}
	return nil
}

func (r Repository) SavingsByYear(year income.Year) (List, error) {
	var list List
	rows, err := r.DB.Query("SELECT amount, comment, created_at FROM tax_savings WHERE year=?", year)
	if err != nil {
		return list, fmt.Errorf("tax savings query error: %v", err)
	}
	for rows.Next() {
		var saving Saving
		if err := rows.Scan(&saving.Amount, &saving.Comment, &saving.CreatedAt); err != nil {
			return list, fmt.Errorf("tax savings row iteration error: %v", err)
		}
		list = append(list, saving)
	}
	return list, nil
}
