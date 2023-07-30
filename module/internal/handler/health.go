package handler

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Register(router *httprouter.Router) {
	router.GET("/v1/health", h.GetHealth)
}

func (h *HealthHandler) GetHealth(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}
