package handler

import (
	"encoding/json"
	"net/http"

	"video-server/module/entity"
)

func BuildErrorResponse(w http.ResponseWriter, err error) {
	e, ok := err.(entity.RequestError)
	if !ok {
		WriteHTTPResponse(w, map[string]string{"message": err.Error()}, http.StatusInternalServerError)
		return
	}

	WriteHTTPResponse(w, nil, e.StatusCode)
}

func WriteHTTPResponse(w http.ResponseWriter, body interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}
