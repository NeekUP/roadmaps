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
	EntityType string `json:"entityType"`
	EntityId   string `json:"entityId"`
	ParentId   int64  `json:"parentId"`
	Text       string `json:"text"`
	Title      string `json:"title"`
}

func (req *addCommentRequest) Sanitize() {
	req.Title = StrictSanitize(req.Title)
	req.Text = SanitizeText(req.Text)
	req.EntityId = StrictSanitize(req.EntityId)
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

		data.Sanitize()
		var entityType domain.EntityType
		var isValidType bool
		if isValidType, entityType = domain.EntityTypeFromString(data.EntityType); !isValidType {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		entityId, err := strconv.ParseInt(data.EntityId, 10, 64)
		if err != nil || entityType == domain.PlanEntity {
			id, err := core.DecodeStringToNum(data.EntityId)
			if err != nil {
				errors := make(map[string]string)
				errors["entityId"] = core.InvalidValue.String()
				badRequest(w, core.ValidationError(errors))
				return
			}
			entityId = int64(id)
		}

		comment, err := addComment.Do(infrastructure.NewContext(r.Context()), entityType, entityId, data.ParentId, data.Text, data.Title)
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

func (req *editCommentRequest) Sanitize() {
	req.Title = StrictSanitize(req.Title)
	req.Text = SanitizeText(req.Text)
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
		data.Sanitize()
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
	EntityType string `json:"entityType"`
	EntityId   string `json:"entityId"`
	Count      int    `json:"count"`
	Page       int    `json:"page"`
}

func (req *getCommentThreadsRequest) Sanitize() {
	req.EntityId = StrictSanitize(req.EntityId)
}

type getCommentThreadsResponse struct {
	HasMore  bool      `json:"hasMore"`
	Page     int       `json:"page"`
	Comments []comment `json:"comments"`
}

func GetThreads(getThreads usecases.GetCommentsThreads, getPointsList usecases.GetPointsList, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(getCommentThreadsRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}
		data.Sanitize()
		var entityType domain.EntityType
		var isValidType bool
		if isValidType, entityType = domain.EntityTypeFromString(data.EntityType); !isValidType {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		entityId, err := strconv.ParseInt(data.EntityId, 10, 64)
		if err != nil || entityType == domain.PlanEntity {
			id, err := core.DecodeStringToNum(data.EntityId)
			if err != nil {
				errors := make(map[string]string)
				errors["entityId"] = core.InvalidValue.String()
				badRequest(w, core.ValidationError(errors))
				return
			}
			entityId = int64(id)
		}

		ctx := infrastructure.NewContext(r.Context())
		comments, hasMore, err := getThreads.Do(ctx, entityType, entityId, data.Count, data.Page)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		attachePoints(comments, getPointsList, ctx, log)

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
	EntityType string `json:"entityType"`
	EntityId   string `json:"entityId"`
	ThreadId   int64  `json:"threadId"`
}

func (req *getCommentThreadRequest) Sanitize() {
	req.EntityId = StrictSanitize(req.EntityId)
}

type getCommentThreadResponse struct {
	Comments []comment `json:"comments"`
}

func GetThread(getThread usecases.GetCommentsThread, getPointsList usecases.GetPointsList, log core.AppLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		data := new(getCommentThreadRequest)
		err := decoder.Decode(data)
		defer r.Body.Close()

		if err != nil {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}
		data.Sanitize()
		var entityType domain.EntityType
		var isValidType bool
		if isValidType, entityType = domain.EntityTypeFromString(data.EntityType); !isValidType {
			statusResponse(w, &status{Code: http.StatusBadRequest})
			return
		}

		entityId, err := strconv.ParseInt(data.EntityId, 10, 64)
		if err != nil || entityType == domain.PlanEntity {
			id, err := core.DecodeStringToNum(data.EntityId)
			if err != nil {
				errors := make(map[string]string)
				errors["entityId"] = core.InvalidValue.String()
				badRequest(w, core.ValidationError(errors))
				return
			}
			entityId = int64(id)
		}

		ctx := infrastructure.NewContext(r.Context())
		comments, err := getThread.Do(ctx, entityType, entityId, data.ThreadId)
		if err != nil {
			if err.Error() != core.InternalError.String() {
				badRequest(w, err)
			} else {
				statusResponse(w, &status{Code: 500})
			}
			return
		}

		attachePoints(comments, getPointsList, ctx, log)

		c := make([]comment, len(comments))
		for i := 0; i < len(comments); i++ {
			c[i] = *NewCommentDto(&comments[i])
		}

		valueResponse(w, c)
	}
}

func attachePoints(comments []domain.Comment, getPointsList usecases.GetPointsList, ctx core.ReqContext, log core.AppLogger) {
	idList := make([]int64, len(comments))
	for i := 0; i < len(comments); i++ {
		idList[i] = int64(comments[i].Id)
	}

	points, err := getPointsList.Do(ctx, domain.CommentEntity, idList)
	if err != nil {
		log.Errorw("fail to retrieve points for comments",
			"reqid", ctx.ReqId(),
			"error", "see db log")
	} else {
		for i := 0; i < len(comments); i++ {
			for j := 0; j < len(comments); j++ {
				if int64(comments[j].Id) == points[i].Id {
					comments[j].Points = &points[i]
					break
				}
			}
		}
	}
}
