package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	router := httprouter.New()

	router.GET("/ping", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("pong"))
	})
	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	reponse := recorder.Result()
	body, _ := io.ReadAll(reponse.Body)
	assert.Equal(t, "pong", string(body))
}

func TestParams(t *testing.T) {
	router := httprouter.New()

	router.GET("/products/:id", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id := params.ByName("id")
		fmt.Fprint(w, id)
	})
	request := httptest.NewRequest(http.MethodGet, "/products/2", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	reponse := recorder.Result()
	body, _ := io.ReadAll(reponse.Body)
	assert.Equal(t, "2", string(body))
}

func TestCatchAll(t *testing.T) {
	router := httprouter.New()

	router.GET("/products/*id", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id := params.ByName("id")
		fmt.Fprint(w, id)
	})
	request := httptest.NewRequest(http.MethodGet, "/products/gas/aha", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	reponse := recorder.Result()
	body, _ := io.ReadAll(reponse.Body)
	assert.Equal(t, "/gas/aha", string(body))
}

//go:embed resources
var resources embed.FS

func TestServeFiles(t *testing.T) {
	router := httprouter.New()

	dir, _ := fs.Sub(resources, "resources")
	router.ServeFiles("/files/*filepath", http.FS(dir))

	request := httptest.NewRequest(http.MethodGet, "/files/a.txt", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	reponse := recorder.Result()
	body, _ := io.ReadAll(reponse.Body)
	assert.Equal(t, "aaaa", string(body))
}

func TestPanic(t *testing.T) {
	router := httprouter.New()

	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, panic interface{}) {
		fmt.Fprint(w, "Panic: ", panic)
	}

	router.GET("/panic", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		panic("ups")
	})
	request := httptest.NewRequest(http.MethodGet, "/panic", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	reponse := recorder.Result()
	body, _ := io.ReadAll(reponse.Body)
	assert.Equal(t, "Panic: ups", string(body))
}

func TestNotFound(t *testing.T) {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "not found")
	})

	request := httptest.NewRequest(http.MethodGet, "/panic", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	reponse := recorder.Result()
	body, _ := io.ReadAll(reponse.Body)
	assert.Equal(t, "not found", string(body))
}

func TestNotAllowed(t *testing.T) {
	router := httprouter.New()

	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "not allowed")
	})

	router.POST("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "not allowed")
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	reponse := recorder.Result()
	body, _ := io.ReadAll(reponse.Body)
	assert.Equal(t, "not allowed", string(body))
}
