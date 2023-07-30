package config

import (
	"github.com/julienschmidt/httprouter"

	"video-server/module/internal/handler"
)

func RegisterHandler(router *httprouter.Router, usecase *Usecase) {

	healthHandler := handler.NewHealthHandler()
	fileHandler := handler.NewFileHandler(usecase.FileUsecase)

	healthHandler.Register(router)
	fileHandler.Register(router)
}
