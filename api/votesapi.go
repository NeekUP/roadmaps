package api

import (
	"encoding/json"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/core/usecases"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/NeekUP/roadmaps/infrastructure"
	"net/http"
	"strconv"
)

type addVoteRequest struct {
	Id    string `json:"id"`
	Type  string `json:"type"`
	Value int    `json:"value"`
}

type addVoteResponse struct {
	Success bool `json:"success"`
}

func AddPoints(addSource usecases.AddVote, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		data := new(addVoteRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		var entityType domain.EntityType
		var isValidType bool
		if isValidType, entityType = domain.EntityTypeFromString(data.Type); !isValidType {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		entityId, err := strconv.ParseInt(data.Type, 10, 64)
		if err != nil || entityType == domain.PlanEntity {
			id, err := core.DecodeStringToNum(data.Id)
			if err != nil {
				errors := make(map[string]string)
				errors["id"] = core.InvalidValue.String()
				badRequest(w, core.ValidationError(errors))
				return
			}
			entityId = int64(id)
		}

		success, err := addSource.Do(infrastructure.NewContext(r.Context()), entityType, entityId, data.Value)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		valueResponse(w, &addVoteResponse{
			Success: success,
		})
	}
}
