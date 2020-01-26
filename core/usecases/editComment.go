package usecases

import (
	"fmt"

	"github.com/NeekUP/roadmaps/core"
)

type EditComment interface {
	Do(ctx core.ReqContext, id int64, text string, title string) (bool, error)
}

type editComment struct {
	CommentsRepo core.CommentsRepository
	Log          core.AppLogger
}

func NewEditComment(commentsRepo core.CommentsRepository, log core.AppLogger) EditComment {
	return &editComment{CommentsRepo: commentsRepo, Log: log}
}

func (ac *editComment) Do(ctx core.ReqContext, id int64, text string, title string) (bool, error) {
	appErr := ac.validate(id, text, title)
	if appErr != nil {
		ac.Log.Errorw("Invalid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return false, appErr
	}

	comment := ac.CommentsRepo.Get(id)
	if comment == nil || comment.Deleted {
		ac.Log.Errorw("Invalid request",
			"ReqId", ctx.ReqId(),
			"UserId", ctx.UserId(),
			"Error", fmt.Sprintf("comment deleted or not existed. id: %v", id),
		)
		return false, core.NewError(core.NotExists)
	}

	if comment.UserId != ctx.UserId() {
		ac.Log.Errorw("Invalid request",
			"ReqId", ctx.ReqId(),
			"UserId", ctx.UserId(),
			"Error", fmt.Sprintf("access denied. id: %v", id),
		)
		return false, core.NewError(core.AccessDenied)
	}

	if ok, err := ac.CommentsRepo.Update(id, string(core.SanitizeString(text)), core.SanitizeString(title)); !ok {
		if err != nil {
			ac.Log.Errorw("Invalid request",
				"ReqId", ctx.ReqId(),
				"Error", err.Error(),
			)
		}
		return false, err
	}

	return true, nil
}

func (ac *editComment) validate(id int64, text string, title string) *core.AppError {
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
