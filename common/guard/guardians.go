package guards

import (
	"lb/common/log"
	"net/http"
)

var logger = log.Log

func HttpThrowServerError(w http.ResponseWriter, err error, message string, v ...interface{}) {
	HttpThrowError(w, http.StatusInternalServerError, message, v)
	logger.Err(err).Msgf(message, v)
}

func HttpThrowError(w http.ResponseWriter, httpCode int, message string, v ...interface{}) {
	http.Error(w, message, httpCode)
	logger.Error().Msgf(message, v)
}

func FailOnError(err error, message string, v ...interface{}) {
	if err != nil {
		logger.Panic().Err(err).Msgf(message, v)
	}
}
