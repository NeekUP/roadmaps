package api

import (
	"encoding/json"
	"net/http"
	"text/template"
)

type status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func statusResponse(w http.ResponseWriter, status *status) {
	w.WriteHeader(status.Code)

	if status.Message == "" {
		status.Message = http.StatusText(status.Code)
	}
}

func valueResponse(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

func badRequest(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(payload)
}

func HtmlEscape(value string) string {
	return template.HTMLEscapeString(value)
}
