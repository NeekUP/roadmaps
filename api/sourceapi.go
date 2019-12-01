package api

import (
	"encoding/json"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/core/usecases"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/NeekUP/roadmaps/infrastructure"
	"net/http"
)

type addSourceRequest struct {
	Identifier string            `json:"identifier"`
	Type       domain.SourceType `json:"type"`
	Props      map[string]string `json:"props"`
}

type addSourceResponse struct {
	Id         int64             `json:"id"`
	Title      string            `json:"title"`
	Identifier string            `json:"identifier"`
	Type       domain.SourceType `json:"type"`
	Img        string            `json:"img"`
	Desc       string            `json:"desc"`
}

func AddSource(addSource usecases.AddSource, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		data := new(addSourceRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		source, err := addSource.Do(infrastructure.NewContext(r.Context()), data.Identifier, data.Props, data.Type)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		valueResponse(w, &addSourceResponse{
			Id:         source.Id,
			Title:      source.Title,
			Identifier: source.Identifier,
			Type:       source.Type,
			Img:        source.Img,
			Desc:       source.Desc,
		})
	}
}
