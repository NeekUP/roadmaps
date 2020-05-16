package api

import (
	"encoding/json"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/core/usecases"
	"github.com/NeekUP/roadmaps/infrastructure"
	"net/http"
)

type addTopicRequest struct {
	Title string   `json:"title"`
	Desc  string   `json:"desc"`
	Tags  []string `json:"tags"`
	IsTag bool     `json:"istag"`
}

func (req *addTopicRequest) Sanitize() {
	req.Title = StrictSanitize(req.Title)
	req.Desc = StrictSanitize(req.Desc)
	sanitizedTags := make([]string, len(req.Tags))
	for i, v := range req.Tags {
		sanitizedTags[i] = StrictSanitize(v)
	}
	req.Tags = sanitizedTags
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
		data.Sanitize()
		topic, err := addTopic.Do(context, data.Title, data.Desc, data.IsTag, data.Tags)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		valueResponse(w, NewTopicDto(topic))
	}
}

type editTopicRequest struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
	IsTag bool   `json:"istag"`
}

func (req *editTopicRequest) Sanitize() {
	req.Title = StrictSanitize(req.Title)
	req.Desc = StrictSanitize(req.Desc)
}

func EditTopic(editTopic usecases.EditTopic, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		data := new(editTopicRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		context := infrastructure.NewContext(r.Context())
		data.Sanitize()
		_, err = editTopic.Do(context, data.Id, data.Title, data.Desc, data.IsTag)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		statusResponse(w, &status{Code: 200})
	}
}

type getTopicTreeRequest struct {
	Name string `json:"name"`
}

func (req *getTopicTreeRequest) Sanitize() {
	req.Name = StrictSanitize(req.Name)
}

type getTopicTreeResponse struct {
	Nodes []treeNode `json:"nodes"`
}

func GetTopicTree(getTopicTree usecases.GetPlanTree, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(getTopicTreeRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}
		data.Sanitize()
		trees, err := getTopicTree.DoByTopic(infrastructure.NewContext(r.Context()), data.Name)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		roots := make([]treeNode, len(trees))
		if len(trees) > 0 {
			for i := 0; i < len(trees); i++ {
				newPlanTree(trees[i], &roots[i])
			}
		}

		valueResponse(w, &getPlanTreeResponse{Nodes: roots})
	}
}

type getTopicRequest struct {
	Name      string `json:"name"`
	PlanCount int    `json:"planCount"`
}

func (req *getTopicRequest) Sanitize() {
	req.Name = StrictSanitize(req.Name)
}

type getTopicResponse struct {
	Topic *topic `json:"topic"`
}

func GetTopic(getTopic usecases.GetTopic, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(getTopicTreeRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		t, err := getTopic.Do(infrastructure.NewContext(r.Context()), data.Name, 10)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		topicDto := NewTopicDto(t)

		valueResponse(w, &getTopicResponse{Topic: topicDto})
	}
}

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

func SearchTopic(searchTopic usecases.SearchTopic, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
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

type addTopicTagReq struct {
	TopicName string `json:"topicname"`
	TagName   string `json:"tagname"`
}

func (req *addTopicTagReq) Sanitize() {
	req.TopicName = StrictSanitize(req.TopicName)
	req.TagName = StrictSanitize(req.TagName)
}

type addTopicTagRes struct {
	Added bool `json:"added"`
}

func AddTopicTag(addTopicTag usecases.AddTopicTag, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(addTopicTagReq)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}
		data.Sanitize()
		added, err := addTopicTag.Do(infrastructure.NewContext(r.Context()), data.TagName, data.TopicName)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		valueResponse(w, &addTopicTagRes{Added: added})
	}
}

type removeTopicTagReq struct {
	TopicName string `json:"topicname"`
	TagName   string `json:"tagname"`
}

func (req *removeTopicTagReq) Sanitize() {
	req.TopicName = StrictSanitize(req.TopicName)
	req.TagName = StrictSanitize(req.TagName)
}

type removeTopicTagRes struct {
	Removed bool `json:"removed"`
}

func RemoveTopicTag(removeTopicTag usecases.RemoveTopicTag, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(removeTopicTagReq)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}
		data.Sanitize()
		removed, err := removeTopicTag.Do(infrastructure.NewContext(r.Context()), data.TagName, data.TopicName)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		valueResponse(w, &removeTopicTagRes{Removed: removed})
	}
}
