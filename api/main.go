package main

import (
	"api/controller"
	"api/exception"
	"api/helper"
	"api/middleware"
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/belajar_go?parseTime=true")
	helper.PanicIfError(err)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Second)
	db.SetConnMaxLifetime(60 * time.Second)
	defer db.Close()

	router := httprouter.New()
	router.GET("/", controller.HomeIndex)

	catController := controller.NewCategoryController(db)
	router.GET("/api/categories", catController.Index)
	router.GET("/api/categories/:categoryId", catController.Show)
	router.POST("/api/categories", catController.Store)
	router.PUT("/api/categories/:categoryId", catController.Update)
	router.DELETE("/api/categories/:categoryId", catController.Destroy)

	router.PanicHandler = exception.ErrorHandler

	logMiddleware := middleware.LogMiddleware{Handler: router}

	log.Println("listen and serve http://localhost:8080")
	err = http.ListenAndServe("localhost:8080", &logMiddleware)
	helper.PanicIfError(err)
}
