package usecases

import (
	"fmt"
	"github.com/NeekUP/roadmaps/domain"

	"github.com/NeekUP/roadmaps/core"
)

type RemoveComment interface {
	Do(ctx core.ReqContext, id int64) (bool, error)
}

type removeComment struct {
	commentsRepo core.CommentsRepository
	log          core.AppLogger
	changeLog    core.ChangeLog
}

func NewRemoveComments(commentsRepo core.CommentsRepository, changeLog core.ChangeLog, log core.AppLogger) RemoveComment {
	return &removeComment{commentsRepo: commentsRepo, changeLog: changeLog, log: log}
}

func (usecase removeComment) Do(ctx core.ReqContext, id int64) (bool, error) {
	trace := ctx.StartTrace("removeComment")
	defer ctx.StopTrace(trace)

	userId := ctx.UserId()
	comment := usecase.commentsRepo.Get(ctx, id)
	if comment == nil || comment.UserId != userId {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"UserId", userId,
			"error", fmt.Sprintf("access denied or not comment existed. id: %v", id),
		)
		return false, core.NewError(core.AccessDenied)
	}

	ok, err := usecase.commentsRepo.Delete(ctx, id)
	if ok {
		changedComment := *comment
		changedComment.Deleted = true
		usecase.changeLog.Edited(domain.CommentEntity, comment.Id, userId, comment, &changedComment)
	}
	return ok, err
}
