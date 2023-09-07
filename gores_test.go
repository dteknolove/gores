package gores

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func dummyResult() map[string]interface{} {
	return map[string]interface{}{
		"keyString": "valueString",
		"keyInt":    999,
		"keyFloat":  20.6,
		"keyFunc":   []string{"one", "two", "three"},
	}
}

const (
	msgDumErr = "dummy error something"
)

func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	func(w http.ResponseWriter, r *http.Request) {
		dataMap := dummyResult()
		Success(dataMap, "success do something", http.StatusOK, w)
	}(w, r)
}
func TestSuccessCreateOrUpdate(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	func(w http.ResponseWriter, r *http.Request) {
		dataMap := dummyResult()
		SuccessCreateOrUpdate(dataMap, "success create", w)
	}(w, r)
}
func TestUnAuthorized(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	func(w http.ResponseWriter, r *http.Request) {
		err := errors.New(msgDumErr)
		UnAuthorized(w, err)
	}(w, r)
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	func(w http.ResponseWriter, r *http.Request) {
		err := errors.New(msgDumErr)
		Error(err, "error something", http.StatusInternalServerError, w)
	}(w, r)
}

func TestErrorBool(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	func(w http.ResponseWriter, r *http.Request) {
		err := errors.New(msgDumErr)
		if ErrorBool(err, "error name", http.StatusInternalServerError, w) {
			return
		}
	}(w, r)
}
