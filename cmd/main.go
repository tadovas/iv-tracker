package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

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

	httpServer, err := server.Setup(serverFlags, chi.NewRouter())
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
