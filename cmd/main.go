package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/tadovas/iv-tracker/tax"

	"github.com/go-chi/chi/middleware"

	"github.com/tadovas/iv-tracker/income"
	"github.com/tadovas/iv-tracker/rest"

	"github.com/go-chi/chi"

	"github.com/tadovas/iv-tracker/server"

	"github.com/tadovas/iv-tracker/db"
	"github.com/tadovas/iv-tracker/log"
)

func main() {
	var dbFlags db.Flags
	db.RegisterFlags(&dbFlags)

	var serverFlags server.Flags
	server.RegisterFlags(&serverFlags)

	flag.Parse()

	DB, err := db.Setup(dbFlags)
	failOnError(err)
	failOnError(DB.Ping())
	failOnError(db.Migrate(DB, filepath.Join("db", "migrations")))

	incomeRepository := income.Repository{
		DB: DB,
	}

	taxCalcLoader := tax.CalculatorDBLoader{DB: DB}

	globalRouter := chi.NewRouter()
	globalRouter.Use(middleware.SetHeader("Content-type", "application/json"))

	incomesRouter := globalRouter.Route("/incomes", nil)
	incomesRouter.Post("/", rest.AddIncome(incomeRepository))
	incomesRouter.Get("/years", rest.ListIncomeYears(incomeRepository))
	incomesRouter.Get("/{year}", rest.ListIncomesByYear(incomeRepository))

	taxesRouter := globalRouter.Route("/taxes", nil)
	taxesRouter.Get("/{year}", rest.TaxSummaryView(incomeRepository, taxCalcLoader))

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
