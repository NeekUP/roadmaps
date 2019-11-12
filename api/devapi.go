// +build DEV

package api

import (
	"net/http"
	"roadmaps/core/usecases"
)

func ListTopics(usecase usecases.ListTopicsDev) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		valueResponse(w, usecase.Do())
	}
}

func ListPlans(usecase usecases.ListPlansDev) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		valueResponse(w, usecase.Do())
	}
}

func ListSteps(usecase usecases.ListStepsDev) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		valueResponse(w, usecase.Do())
	}
}

func ListSources(usecase usecases.ListSourcesDev) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		valueResponse(w, usecase.Do())
	}
}

func ListUsers(usecase usecases.ListUsersDev) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		valueResponse(w, usecase.Do())
	}
}
