package api

import (
	"net/http"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger() {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

// HandleError logs the error and writes an HTTP error response
func HandleError(w http.ResponseWriter, msg string, err error, status int) {
	Logger.Error(msg, zap.Error(err))
	http.Error(w, msg, status)
}
