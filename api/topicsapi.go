package api

import (
	"encoding/json"
	"net/http"
	"roadmaps/core"
	"roadmaps/core/usecases"
	"roadmaps/infrastructure"
)

type addTopicRequest struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type addTopicResponse struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

func AddTopic(addTopic usecases.AddTopic, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		data := new(addTopicRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		context := infrastructure.NewContext(r.Context())
		topic, err := addTopic.Do(context, data.Title, data.Desc)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		valueResponse(w, &addTopicResponse{
			Title: topic.Title,
			Name:  topic.Name,
			Desc:  topic.Description,
		})
	}
}
