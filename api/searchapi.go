package api

import (
	"encoding/json"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/core/usecases"
	"github.com/NeekUP/roadmaps/infrastructure"
	"net/http"
)

type searchTopicReq struct {
	Tags  []string `json:"tags"`
	Query string   `json:"query"`
	Count int      `json:"count"`
}

func (req *searchTopicReq) Sanitize() {
	req.Query = StrictSanitize(req.Query)
	for i, tag := range req.Tags {
		req.Tags[i] = StrictSanitize(tag)
	}
}

type searchTopicRes struct {
	Str    string  `json:"query"`
	Result []topic `json:"result"`
}

func SearchTopic(searchTopic usecases.Search, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(searchTopicReq)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}
		if data.Tags == nil {
			data.Tags = []string{}
		}

		data.Sanitize()
		topics := searchTopic.Do(infrastructure.NewContext(r.Context()), data.Query, data.Tags, data.Count)
		dtos := make([]topic, len(topics))
		for i, t := range topics {
			dtos[i] = *NewTopicDto(&t)
		}
		valueResponse(w, &searchTopicRes{Result: dtos, Str: data.Query})
	}
}
