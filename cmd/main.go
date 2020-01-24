package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
)

func main() {

	server := http.Server{}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	go func() {
		<-stop
		if err := server.Close(); err != nil {
			fmt.Println("Service stop error:", err)
		}
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Println("Http service error:", err)
	}
	fmt.Println("Terminated")
}
