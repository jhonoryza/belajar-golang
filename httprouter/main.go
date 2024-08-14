package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	router.GET("/ping", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("pong"))
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}
