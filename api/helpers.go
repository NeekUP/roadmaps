package api

import (
	"encoding/json"
	"net/http"
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
