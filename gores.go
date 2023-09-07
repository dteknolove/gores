package gores

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"net/http"
)

const (
	errEncode       = "error to encode: "
	errUnauthorized = "error you are not authorized"
	keyData         = "data"
	keyMessage      = "message"
	keyStatusCode   = "statusCode"
	keySuccess      = "success"
	contentType     = "Content-Type"
	appJson         = "application/json"
	msgCreateUpdate = "create or update"
	msgSuccess      = "success"
)

func Success(data interface{}, message string, statusCode int, w http.ResponseWriter) {
	slog.Info(msgSuccess, slog.Any(keyData, data))
	response := map[string]interface{}{
		keyData:       data,
		keyMessage:    message,
		keyStatusCode: statusCode,
		keySuccess:    true,
	}
	w.Header().Set(contentType, appJson)
	w.WriteHeader(statusCode)
	errJson(response, w)
}

func SuccessCreateOrUpdate(data interface{}, message string, w http.ResponseWriter) {
	slog.Info(msgCreateUpdate, slog.Any(keyData, data))
	w.Header().Set(contentType, appJson)
	w.WriteHeader(http.StatusCreated)
	errJson(successResponse(data, message, http.StatusCreated), w)
}

func UnAuthorized(w http.ResponseWriter, err error) {
	slog.Error(err.Error())
	w.Header().Set(contentType, appJson)
	w.WriteHeader(http.StatusUnauthorized)
	errJson(errResponse(errUnauthorized, err, http.StatusUnauthorized), w)
}

func Error(err error, message string, statusCode int, w http.ResponseWriter) {
	defer slog.Error(err.Error())
	errResponse(message, err, statusCode)
	w.Header().Set(contentType, appJson)
	w.WriteHeader(statusCode)
	errJson(errResponse(message, err, statusCode), w)
}

func ErrorBool(err error, errName string, status int, w http.ResponseWriter) bool {
	if err != nil {
		Error(err, fmt.Sprintf("%s: %s", errName, err.Error()), status, w)
		return true
	}
	return false
}

func ErrorTx(tx pgx.Tx, ctx context.Context, errName string, w http.ResponseWriter) {
	errRoll := tx.Rollback(ctx)
	if errRoll != nil {
		Error(errRoll, fmt.Sprintf("ERROR rollback: %s, ON: %s", errRoll.Error(), errName), http.StatusInternalServerError, w)
		return
	}
	Error(nil, fmt.Sprintf("%s, execute ROLLBACK", errName), http.StatusInternalServerError, w)
}

func ErrorBoolTx(tx pgx.Tx, ctx context.Context, err error, errName string, w http.ResponseWriter) bool {
	if err != nil {
		errRoll := tx.Rollback(ctx)
		if errRoll != nil {
			slog.Error(errName, errName)
			Error(errRoll, fmt.Sprintf("%s: %s, ON: %s", "error ROLLBACK", errRoll.Error(), errName), http.StatusInternalServerError, w)
			return true
		}
		slog.Error(errName, err)
		Error(err, fmt.Sprintf("%s | execute ROLLBACK, %s", errName, err.Error()), http.StatusInternalServerError, w)
		return true
	}
	return false
}

func errJson(response map[string]interface{}, w http.ResponseWriter) {
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error(errEncode+err.Error(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func errResponse(msg string, err error, statusCode int) map[string]interface{} {
	return map[string]interface{}{
		keyData: nil,
		keyMessage: map[string]any{
			"errorMessage": msg,
			"error":        err,
		},
		"statusCode": statusCode,
		"success":    false,
	}
}
func successResponse(data interface{}, message string, statusCode int) map[string]interface{} {
	return map[string]interface{}{
		keyData:       data,
		keyMessage:    message,
		keyStatusCode: statusCode,
		keySuccess:    true,
	}
}
