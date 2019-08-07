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

	json.NewEncoder(w).Encode(status)
}
