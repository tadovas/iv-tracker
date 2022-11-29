package rest

import (
	"fmt"
	"net/http"

	"github.com/tadovas/iv-tracker/log"

	"github.com/tadovas/iv-tracker/reports"
)

func GenerateIncomeExpenseJournal(reportsGenerator reports.JournalGenerator) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		report, err := reportsGenerator.GenerateIncomeExpenseJournalAsExcel()
		if err != nil {
			http.Error(writer, fmt.Sprintf("report error: %v", err), http.StatusInternalServerError)
		}
		writer.Header().Set("Content-Disposition", "attachment; filename=IV_pajamu_islaidu_zurnalas.xlsx")
		writer.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		writer.WriteHeader(http.StatusOK)
		if _, err := writer.Write(report); err != nil {
			log.Errorf("Error sending report to client: %v", err)
		}
	}
}
