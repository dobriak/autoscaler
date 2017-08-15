package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			log.Fatalln("Received an interrupt, stopping tickers")
			RepoRemoveAllApps()
			cleanupDone <- true
		}
	}()

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
	//<-cleanupDone
}
