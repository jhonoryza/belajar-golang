package helper

import (
	"encoding/json"
	"net/http"
)

type WebResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (web *WebResponse) ToJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err := encoder.Encode(web)
	PanicIfError(err)
}
