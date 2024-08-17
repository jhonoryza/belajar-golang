package controller

import (
	"api/helper"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func HomeIndex(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	data := map[string]interface{}{
		"code": http.StatusOK,
		"msg":  "server ok",
	}
	jsonData, _ := json.Marshal(data)
	_, err := w.Write(jsonData)
	helper.PanicIfError(err)
}
