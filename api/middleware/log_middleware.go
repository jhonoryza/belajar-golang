package middleware

import (
	"log"
	"net/http"
)

type LogMiddleware struct {
	Handler http.Handler
}

func (logMiddleware *LogMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("URL %v\n", r.URL)
	log.Printf("Method %v\n", r.Method)
	log.Printf("BODY %v\n", r.Body)
	logMiddleware.Handler.ServeHTTP(w, r)
}
