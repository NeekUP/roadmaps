package usecases

import (
	"sort"

	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type GetCommentsThread interface {
	Do(ctx core.ReqContext, entityType domain.EntityType, entityId int64, threadId int64) ([]domain.Comment, error)
}

type getCommentsThread struct {
	commentsRepo core.CommentsRepository
	usersRepo    core.UserRepository
	log          core.AppLogger
}

func NewGetCommentsThread(commentsRepo core.CommentsRepository, usersRepo core.UserRepository, log core.AppLogger) GetCommentsThread {
	return &getCommentsThread{commentsRepo: commentsRepo, usersRepo: usersRepo, log: log}
}

func (usecase *getCommentsThread) Do(ctx core.ReqContext, entityType domain.EntityType, entityId int64, threadId int64) ([]domain.Comment, error) {
	trace := ctx.StartTrace("getCommentsThread")
	defer ctx.StopTrace(trace)

	appErr := usecase.validate(entityType, entityId, threadId)
	if appErr != nil {
		usecase.log.Errorw("Invalid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return nil, appErr
	}

	allcomments := usecase.commentsRepo.GetThread(ctx, int(entityType), entityId, threadId)
	if len(allcomments) == 0 {
		return []domain.Comment{}, nil
	}

	// attache users
	userIds := make(map[string]*domain.User)
	for i := 0; i < len(allcomments); i++ {
		if _, ok := userIds[allcomments[i].UserId]; !ok {
			userIds[allcomments[i].UserId] = nil
		}
	}

	idList := make([]string, 0, len(userIds))
	for k := range userIds {
		idList = append(idList, k)
	}

	for _, v := range usecase.usersRepo.GetList(ctx, idList) {
		userIds[v.Id] = &v
	}

	for i := 0; i < len(allcomments); i++ {
		allcomments[i].User = userIds[allcomments[i].UserId]
	}

	// sort
	m := make(map[int64]domain.Comment)
	for _, v := range allcomments {
		m[v.Id] = v
	}

	for key, val := range m {
		if v, ok := m[val.ParentId]; ok {
			if v.Childs == nil {
				v.Childs = []domain.Comment{val}
			} else {
				v.Childs = append(v.Childs, val)
				sort.Slice(v.Childs, func(i, j int) bool { return v.Childs[i].Id < v.Childs[j].Id })
			}
			m[val.ParentId] = v
			delete(m, key)
		}
	}

	comments := make([]domain.Comment, len(m))
	var i int = 0
	for _, val := range m {
		comments[i] = val
		i++
	}
	return comments, nil
}

func (usecase *getCommentsThread) validate(entityType domain.EntityType, entityId int64, threadId int64) *core.AppError {
	errors := make(map[string]string)
	if !entityType.IsValid() {
		errors["entityType"] = core.InvalidValue.String()
	}

	if entityId < 0 {
		errors["entityId"] = core.InvalidValue.String()
	}

	if threadId <= 0 {
		errors["threadId"] = core.InvalidValue.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}
