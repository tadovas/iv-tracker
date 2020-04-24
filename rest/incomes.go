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
)

func AddIncome(repository income.Repository) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var newIncome income.Income
		if err := json.NewDecoder(request.Body).Decode(&newIncome); err != nil {
			http.Error(writer, fmt.Sprintf("json parse error: %v", err), http.StatusBadRequest)
			return
		}
		if newIncome.Date.IsZero() {
			newIncome.Date = time.Now()
		}
		id, err := repository.AddNewIncome(newIncome)
		if err != nil {
			http.Error(writer, fmt.Sprintf("income insert error: %v", err), http.StatusInternalServerError)
			return
		}
		var respWithID struct {
			ID int `json:"id"`
		}
		respWithID.ID = int(id)
		writer.WriteHeader(http.StatusCreated)
		log.IfError("serialize json response", func() error {
			return json.NewEncoder(writer).Encode(&respWithID)
		})
	}
}

type yearObject struct {
	Year   string       `json:"year"`
	Amount income.Money `json:"total"`
	Count  int          `json:"count"`
}
type yearsResponse struct {
	Years []yearObject `json:"years"`
}

func ListIncomeYears(repository income.Repository) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		yearSummaries, err := repository.ListAllYears()
		if err != nil {
			http.Error(writer, fmt.Sprintf("list years: %v", err), http.StatusInternalServerError)
			return
		}
		var yearsResp yearsResponse
		for _, yearSummary := range yearSummaries {
			yearsResp.Years = append(yearsResp.Years, yearObject{
				Year:   fmt.Sprint(yearSummary.Year),
				Amount: yearSummary.TotalAmount,
				Count:  yearSummary.TotalIncomes,
			})
		}
		log.IfError("serializing json", func() error {
			return json.NewEncoder(writer).Encode(&yearsResp)
		})
	}
}

func ListIncomesByYear(repository income.Repository) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		year, err := strconv.Atoi(chi.URLParam(request, "year"))
		if err != nil {
			http.Error(writer, fmt.Sprintf("URL param parse error: %v", err), http.StatusBadRequest)
			return
		}
		incomes, err := repository.ListIncomesByYear(income.Year(year))
		if err != nil {
			http.Error(writer, fmt.Sprintf("Incomes query error: %v", err), http.StatusInternalServerError)
			return
		}
		var incomesResponse struct {
			Incomes []income.Income `json:"incomes"`
		}
		incomesResponse.Incomes = incomes
		log.IfError("render incomes", func() error {
			return json.NewEncoder(writer).Encode(&incomesResponse)
		})
	}
}

func YearlyIncomeDeclaration(repository income.Repository) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		year, err := strconv.Atoi(chi.URLParam(request, "year"))
		if err != nil {
			http.Error(writer, fmt.Sprintf("URL param parse error: %v", err), http.StatusBadRequest)
			return
		}
		countryIncomes, err := repository.YearlyIncomesByCountry(income.Year(year))
		if err != nil {
			http.Error(writer, fmt.Sprintf("Incomes query error: %v", err), http.StatusInternalServerError)
			return
		}
		var incomesResponse struct {
			Declaration []income.CountryIncome `json:"declaration"`
		}
		incomesResponse.Declaration = countryIncomes
		log.IfError("render incomes", func() error {
			return json.NewEncoder(writer).Encode(&incomesResponse)
		})
	}
}
