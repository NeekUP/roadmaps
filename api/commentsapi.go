package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/core/usecases"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/NeekUP/roadmaps/infrastructure"
)

type addCommentRequest struct {
	EntityType int    `json:"entityType"`
	EntityId   string `json:"entityId"`
	ParentId   int64  `json:"parentId"`
	Text       string `json:"text"`
	Title      string `json:"title"`
}

type addCommentResponse struct {
	Id int64 `json:"id"`
}

func AddComment(addComment usecases.AddComment, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(addCommentRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		entityId, err := strconv.ParseInt(data.EntityId, 10, 64)
		if err != nil || data.EntityType == int(domain.PlanEntity) {
			id, err := core.DecodeStringToNum(data.EntityId)
			if err != nil {
				errors := make(map[string]string)
				errors["entityId"] = core.InvalidValue.String()
				badRequest(w, core.ValidationError(errors))
				return
			}
			entityId = int64(id)
		}

		comment, err := addComment.Do(infrastructure.NewContext(r.Context()), domain.EntityType(data.EntityType), entityId, data.ParentId, data.Text, data.Title)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		valueResponse(w, &addCommentResponse{
			Id: comment.Id,
		})
	}
}

type editCommentRequest struct {
	Id    int64  `json:"id"`
	Text  string `json:"text"`
	Title string `json:"title"`
}

func EditComment(editComment usecases.EditComment, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(editCommentRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		_, err = editComment.Do(infrastructure.NewContext(r.Context()), data.Id, data.Text, data.Title)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		statusResponse(w, &status{Code: http.StatusOK})
	}
}

type deleteCommentRequest struct {
	Id int64 `json:"id"`
}

func DeleteComment(removeComment usecases.RemoveComment, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(deleteCommentRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		_, err = removeComment.Do(infrastructure.NewContext(r.Context()), data.Id)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		statusResponse(w, &status{Code: http.StatusOK})
	}
}

type getCommentThreadsRequest struct {
	EntityType int    `json:"entityType"`
	EntityId   string `json:"entityId"`
	Count      int    `json:"count"`
	Page       int    `json:"page"`
}

type getCommentThreadsResponse struct {
	HasMore  bool      `json:"hasMore"`
	Page     int       `json:"page"`
	Comments []comment `json:"comments"`
}

func GetThreads(getThreads usecases.GetCommentsThreads, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(getCommentThreadsRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		entityId, err := strconv.ParseInt(data.EntityId, 10, 64)
		if err != nil || data.EntityType == int(domain.PlanEntity) {
			id, err := core.DecodeStringToNum(data.EntityId)
			if err != nil {
				errors := make(map[string]string)
				errors["entityId"] = core.InvalidValue.String()
				badRequest(w, core.ValidationError(errors))
				return
			}
			entityId = int64(id)
		}

		comments, hasMore, err := getThreads.Do(infrastructure.NewContext(r.Context()), domain.EntityType(data.EntityType), entityId, data.Count, data.Page)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		c := make([]comment, len(comments))
		for i := 0; i < len(comments); i++ {
			c[i] = *NewCommentDto(&comments[i])
		}

		valueResponse(w, &getCommentThreadsResponse{
			HasMore:  hasMore,
			Page:     data.Page,
			Comments: c,
		})
	}
}

type getCommentThreadRequest struct {
	EntityType int    `json:"entityType"`
	EntityId   string `json:"entityId"`
	ThreadId   int64  `json:"threadId"`
}

type getCommentThreadResponse struct {
	Comments []comment `json:"comments"`
}

func GetThread(getThread usecases.GetCommentsThread, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(getCommentThreadRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		entityId, err := strconv.ParseInt(data.EntityId, 10, 64)
		if err != nil || data.EntityType == int(domain.PlanEntity) {
			id, err := core.DecodeStringToNum(data.EntityId)
			if err != nil {
				errors := make(map[string]string)
				errors["entityId"] = core.InvalidValue.String()
				badRequest(w, core.ValidationError(errors))
				return
			}
			entityId = int64(id)
		}

		comments, err := getThread.Do(infrastructure.NewContext(r.Context()), domain.EntityType(data.EntityType), entityId, data.ThreadId)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		c := make([]comment, len(comments))
		for i := 0; i < len(comments); i++ {
			c[i] = *NewCommentDto(&comments[i])
		}

		valueResponse(w, c)
	}
}
