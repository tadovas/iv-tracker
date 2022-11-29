package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/tadovas/iv-tracker/db"
	"github.com/tadovas/iv-tracker/income"
	"github.com/tadovas/iv-tracker/log"
	"github.com/tadovas/iv-tracker/reports"
	"github.com/tadovas/iv-tracker/rest"
	"github.com/tadovas/iv-tracker/saving"
	"github.com/tadovas/iv-tracker/server"
	"github.com/tadovas/iv-tracker/tax"
)

func main() {
	var dbFlags db.Flags
	db.RegisterFlags(&dbFlags)

	var serverFlags server.Flags
	server.RegisterFlags(&serverFlags)

	var credentials rest.Credentials
	rest.RegisterFlags(&credentials)

	flag.Parse()

	DB, err := db.Setup(dbFlags)
	failOnError(err)
	failOnError(DB.Ping())
	failOnError(db.Migrate(DB, filepath.Join("db", "migrations")))

	incomeRepository := income.Repository{
		DB: DB,
	}

	taxCalcLoader := tax.CalculatorDBLoader{DB: DB}

	savingsRepository := saving.Repository{DB: DB}

	reportsGenerator := reports.JournalGenerator{IncomeRepository: incomeRepository}

	globalRouter := chi.NewRouter()
	globalRouter.Use(middleware.SetHeader("Content-type", "application/json"))
	if credentials.IsSet() {
		log.Info("Basic auth enabled for user: ", credentials.Username)
		globalRouter.Use(rest.BasicAuth(credentials))
	}

	incomesRouter := globalRouter.Route("/incomes", nil)
	incomesRouter.Post("/", rest.AddIncome(incomeRepository))
	incomesRouter.Get("/years", rest.ListIncomeYears(incomeRepository))
	incomesRouter.Get("/{year}", rest.ListIncomesByYear(incomeRepository))
	incomesRouter.Get("/{year}/declaration", rest.YearlyIncomeDeclaration(incomeRepository))

	taxesRouter := globalRouter.Route("/taxes", nil)
	taxesRouter.Get("/{year}", rest.TaxSummaryView(incomeRepository, savingsRepository, taxCalcLoader))

	savingsRouter := globalRouter.Route("/savings", nil)
	savingsRouter.Post("/", rest.AddNewSaving(savingsRepository))
	savingsRouter.Get("/{year}", rest.ListSavings(savingsRepository))

	reportsRouter := globalRouter.Route("/reports", nil)
	reportsRouter.Get("/journal", rest.GenerateIncomeExpenseJournal(reportsGenerator))

	httpServer, err := server.Setup(serverFlags, globalRouter)
	failOnError(err)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	go func() {
		<-stop
		if err := httpServer.Close(); err != nil {
			fmt.Println("Service stop error:", err)
		}
	}()
	log.Info("Serving at", serverFlags.Address)
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Error("Http service error:", err)
	}
	log.Info("Terminated")
}

func failOnError(err error) {
	if err != nil {
		log.Error("Init error:", err)
		os.Exit(1)
	}
}
