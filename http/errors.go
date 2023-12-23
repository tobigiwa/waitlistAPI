package http

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

const (
	invalidFormData = 900
)

func clientError(w http.ResponseWriter, errStatus int, err error) {
	http.Error(w, err.Error(), errStatus)
}

func serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	http.Error(w, fmt.Sprintf("error: %s\ntrace:\n%s", err.Error(), trace), http.StatusInternalServerError)
}
