package test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

type ErrorMiddleware struct {
	Handler http.Handler
}

func (middleware *ErrorMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, http.StatusText(500))
		}
	}()

	middleware.Handler.ServeHTTP(w, r)
}

func TestMiddleware(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		panic(errors.New("sample error"))
	})

	errMiddleware := ErrorMiddleware{
		Handler: mux,
	}

	err := http.ListenAndServe("localhost:8080", &errMiddleware)
	if err != nil {
		t.Error(err)
	}
}
