package api

import (
	"net/http"
	"roadmaps/core"
	"roadmaps/core/usecases"
)

func AddSource(addSource usecases.AddSource, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		statusResponse(w, &status{Code: 204})
	}
}
