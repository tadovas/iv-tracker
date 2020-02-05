package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/tadovas/iv-tracker/income"
	"github.com/tadovas/iv-tracker/log"
	"github.com/tadovas/iv-tracker/saving"
)

func AddNewSaving(repo saving.Repository) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var newSaving saving.Saving
		if err := json.NewDecoder(request.Body).Decode(&newSaving); err != nil {
			http.Error(writer, fmt.Sprintf("new saving parse body error: %v", err), http.StatusBadRequest)
			return
		}
		if newSaving.CreatedAt.IsZero() {
			newSaving.CreatedAt = time.Now()
		}
		if err := repo.AddNewSaving(newSaving); err != nil {
			http.Error(writer, fmt.Sprintf("new saving store error: %v", err), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusCreated)
	}
}

type SavingsListResponse struct {
	Savings []saving.Saving
}

func ListSavings(repo saving.Repository) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		year, err := strconv.Atoi(chi.URLParam(request, "year"))
		if err != nil {
			http.Error(writer, fmt.Sprintf("URL param parse error: %v", err), http.StatusBadRequest)
			return
		}
		list, err := repo.SavingsByYear(income.Year(year))
		if err != nil {
			http.Error(writer, fmt.Sprintf("savings load error: %v", err), http.StatusInternalServerError)
			return
		}
		resp := SavingsListResponse{Savings: list}
		log.IfError("render saving list response", func() error {
			return json.NewEncoder(writer).Encode(&resp)
		})
	}
}
