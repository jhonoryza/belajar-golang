package exception

import (
	"api/helper"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	if notFoundError(w, r, err) {
		return
	}
	if validationError(w, r, err) {
		return
	}
	internalServerError(w, r, err)
}

func internalServerError(w http.ResponseWriter, r *http.Request, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	resp := helper.WebResponse{
		Code:    http.StatusInternalServerError,
		Message: "INTERNAL_SERVER_ERROR",
		Data:    err,
	}
	resp.ToJson(w)
}

func notFoundError(w http.ResponseWriter, r *http.Request, err interface{}) bool {
	exception, ok := err.(NotFoundError)
	if !ok {
		return false
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		resp := helper.WebResponse{
			Code:    http.StatusNotFound,
			Message: "NOT FOUND",
			Data:    exception.Error,
		}
		resp.ToJson(w)
		return true
	}
}

func validationError(w http.ResponseWriter, r *http.Request, err interface{}) bool {
	exception, ok := err.(validator.ValidationErrors)
	if !ok {
		return false
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := helper.WebResponse{
			Code:    http.StatusBadRequest,
			Message: "BAD REQUEST",
			Data:    exception.Error(),
		}
		resp.ToJson(w)
		return true
	}
}
