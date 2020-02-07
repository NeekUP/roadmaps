package usecases

import (
	"fmt"
	"github.com/NeekUP/roadmaps/domain"

	"github.com/NeekUP/roadmaps/core"
)

type EditComment interface {
	Do(ctx core.ReqContext, id int64, text string, title string) (bool, error)
}

type editComment struct {
	commentsRepo core.CommentsRepository
	log          core.AppLogger
	changeLog    core.ChangeLog
}

func NewEditComment(commentsRepo core.CommentsRepository, changeLog core.ChangeLog, log core.AppLogger) EditComment {
	return &editComment{commentsRepo: commentsRepo, changeLog: changeLog, log: log}
}

func (usecase *editComment) Do(ctx core.ReqContext, id int64, text string, title string) (bool, error) {
	appErr := usecase.validate(id, text, title)
	if appErr != nil {
		usecase.log.Errorw("Invalid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return false, appErr
	}

	userId := ctx.UserId()
	comment := usecase.commentsRepo.Get(id)
	if comment == nil || comment.Deleted {
		usecase.log.Errorw("Invalid request",
			"ReqId", ctx.ReqId(),
			"UserId", userId,
			"Error", fmt.Sprintf("comment deleted or not existed. id: %v", id),
		)
		return false, core.NewError(core.NotExists)
	}

	if comment.UserId != userId {
		usecase.log.Errorw("Invalid request",
			"ReqId", ctx.ReqId(),
			"UserId", userId,
			"Error", fmt.Sprintf("access denied. id: %v", id),
		)
		return false, core.NewError(core.AccessDenied)
	}

	if ok, err := usecase.commentsRepo.Update(id, text, title); !ok {
		if err != nil {
			usecase.log.Errorw("Invalid request",
				"ReqId", ctx.ReqId(),
				"Error", err.Error(),
			)
		}
		return false, err
	}

	changedComment := *comment
	changedComment.Text = text
	changedComment.Title = title
	usecase.changeLog.Edited(domain.CommentEntity, comment.Id, userId, comment, &changedComment)
	return true, nil
}

func (usecase *editComment) validate(id int64, text string, title string) *core.AppError {
	errors := make(map[string]string)

	if id <= 0 {
		errors["id"] = core.InvalidValue.String()
	}

	if !core.IsValidCommentText(text) {
		errors["text"] = core.InvalidValue.String()
	}

	if !core.IsValidCommentTitle(title) {
		errors["title"] = core.InvalidValue.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}
