package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type GetCommentsThreads interface {
	Do(ctx core.ReqContext, entityType domain.EntityType, entityId int64, count int, page int) (comments []domain.Comment, hasMore bool, err error)
}

type getCommentsThreads struct {
	commentsRepo core.CommentsRepository
	usersRepo    core.UserRepository
	log          core.AppLogger
}

func NewGetCommentsThreads(commentsRepo core.CommentsRepository, usersRepo core.UserRepository, log core.AppLogger) GetCommentsThreads {
	return &getCommentsThreads{commentsRepo: commentsRepo, usersRepo: usersRepo, log: log}
}

func (usecase *getCommentsThreads) Do(ctx core.ReqContext, entityType domain.EntityType, entityId int64, count int, page int) (comments []domain.Comment, hasMore bool, err error) {
	appErr := usecase.validate(entityType, entityId, count, page)
	if appErr != nil {
		usecase.log.Errorw("Invalid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return nil, false, appErr
	}

	list := usecase.commentsRepo.GetThreadList(int(entityType), entityId, count, page)
	userIds := make(map[string]*domain.User)
	for i := 0; i < len(list); i++ {
		if _, ok := userIds[list[i].UserId]; !ok {
			userIds[list[i].UserId] = nil
		}
	}

	idList := make([]string, 0, len(userIds))
	for k := range userIds {
		idList = append(idList, k)
	}

	for _, v := range usecase.usersRepo.GetList(idList) {
		userIds[v.Id] = &v
	}

	for i := 0; i < len(list); i++ {
		list[i].User = userIds[list[i].UserId]
	}

	return list, len(list) == count, nil
}

func (usecase *getCommentsThreads) validate(entityType domain.EntityType, entityId int64, count int, page int) *core.AppError {
	errors := make(map[string]string)
	if !entityType.IsValid() {
		errors["entityType"] = core.InvalidValue.String()
	}

	if entityId < 0 {
		errors["entityId"] = core.InvalidValue.String()
	}

	if count <= 0 {
		errors["count"] = core.InvalidValue.String()
	}

	if page < 0 {
		errors["page"] = core.InvalidValue.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}
