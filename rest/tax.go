package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tadovas/iv-tracker/log"

	"github.com/go-chi/chi"
	"github.com/tadovas/iv-tracker/income"
	"github.com/tadovas/iv-tracker/tax"
)

type TaxView struct {
	Total income.Money
	Taxes []tax.Tax
}

type TaxSummary struct {
	Income income.Money
	Taxes  TaxView
}

func TaxSummaryView(incomeRepo income.Repository, taxCalcLoader tax.CalculatorDBLoader) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		yearVal, err := strconv.Atoi(chi.URLParam(request, "year"))
		if err != nil {
			http.Error(writer, fmt.Sprintf("URL param parse error: %v", err), http.StatusBadRequest)
			return
		}
		year := income.Year(yearVal)

		calc, err := taxCalcLoader.LoadFor(year)
		if err != nil {
			http.Error(writer, fmt.Sprintf("calc load error: %v", err), http.StatusInternalServerError)
			return
		}

		incomeList, err := incomeRepo.ListIncomesByYear(year)
		if err != nil {
			http.Error(writer, fmt.Sprintf("income load error: %v", err), http.StatusInternalServerError)
			return
		}

		taxList, err := calc.Taxes(incomeList.Total())
		if err != nil {
			http.Error(writer, fmt.Sprintf("tax calc error: %v", err), http.StatusInternalServerError)
			return
		}

		taxSummary := TaxSummary{
			Income: incomeList.Total(),
			Taxes: TaxView{
				Taxes: taxList,
				Total: taxList.TotalTaxAmount(),
			},
		}
		log.IfError("tax summary render", func() error {
			return json.NewEncoder(writer).Encode(&taxSummary)
		})
	}
}
