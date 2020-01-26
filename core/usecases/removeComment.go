package usecases

import (
	"fmt"

	"github.com/NeekUP/roadmaps/core"
)

type RemoveComment interface {
	Do(ctx core.ReqContext, id int64) (bool, error)
}

type removeComment struct {
	CommentsRepo core.CommentsRepository
	Log          core.AppLogger
}

func NewRemoveComments(commentsRepo core.CommentsRepository, log core.AppLogger) RemoveComment {
	return &removeComment{CommentsRepo: commentsRepo, Log: log}
}

func (usecase removeComment) Do(ctx core.ReqContext, id int64) (bool, error) {

	comment := usecase.CommentsRepo.Get(id)
	if comment == nil || comment.UserId != ctx.UserId() {
		usecase.Log.Errorw("Invalid request",
			"ReqId", ctx.ReqId(),
			"UserId", ctx.UserId(),
			"Error", fmt.Sprintf("access denied or not comment existed. id: %v", id),
		)
		return false, core.NewError(core.AccessDenied)
	}

	ok, err := usecase.CommentsRepo.Delete(id)
	return ok, err
}
