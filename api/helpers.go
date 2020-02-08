package api

import (
	"encoding/json"
	"github.com/microcosm-cc/bluemonday"
	"net/http"
	"regexp"
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

func SanitizeText(text string) string {
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	data := []byte(text)
	return string(p.SanitizeBytes(data))
}

func StrictSanitize(text string) string {
	p := bluemonday.StrictPolicy()
	return string(p.SanitizeBytes([]byte(text)))
}

func sanitizeInput(data []byte) []byte {
	p := bluemonday.UGCPolicy()
	bluemonday.StrictPolicy()
	p.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	return p.SanitizeBytes(data)
}
