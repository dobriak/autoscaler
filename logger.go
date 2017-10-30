package main

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

//Logger provides a decorator for logging
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)
		log.Infoln(
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
