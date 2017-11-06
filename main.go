package main

import (
	"net/http"
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
)

func init() {
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "2006-01-02 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	log.Infoln("Autoscaler starting.")

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
}
