package usecases

import (
	"time"

	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type AddComment interface {
	Do(ctx core.ReqContext, entityType domain.EntityType, entityId int64, parentId int64, text string, title string) (*domain.Comment, error)
}

type addComment struct {
	CommentsRepo core.CommentsRepository
	PlanRepo     core.PlanRepository
	Log          core.AppLogger
}

func NewAddComment(commentsRepo core.CommentsRepository, planRepo core.PlanRepository, log core.AppLogger) AddComment {
	return &addComment{CommentsRepo: commentsRepo, PlanRepo: planRepo, Log: log}
}

func (usecase *addComment) Do(ctx core.ReqContext, entityType domain.EntityType, entityId int64, parentId int64, text string, title string) (*domain.Comment, error) {
	appErr := usecase.validate(entityType, entityId, parentId, text, title)
	if appErr != nil {
		usecase.Log.Errorw("Invalid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return nil, appErr
	}

	var threadId int64 = 0
	if parentId > 0 {
		parent := usecase.CommentsRepo.Get(parentId)
		if parent == nil {
			errors := make(map[string]string)
			errors["parentId"] = core.InvalidValue.String()
			return nil, core.ValidationError(errors)
		}

		threadId = parent.ThreadId
		if threadId == 0 {
			threadId = parentId
		}
	}

	comment := &domain.Comment{
		EntityType: entityType,
		EntityId:   entityId,
		ThreadId:   threadId,
		ParentId:   parentId,
		UserId:     ctx.UserId(),
		Text:       string(core.SanitizeString(text)),
		Title:      title,
		Date:       time.Now().UTC(),
		Deleted:    false,
	}

	if ok, err := usecase.CommentsRepo.Add(comment); !ok {
		if err != nil {
			usecase.Log.Errorw("Invalid request",
				"ReqId", ctx.ReqId(),
				"Error", err.Error(),
			)
		}
		return nil, err
	}

	return comment, nil
}

func (ac *addComment) validate(entityType domain.EntityType, entityId int64, parentId int64, text string, title string) *core.AppError {
	errors := make(map[string]string)
	if !entityType.IsValid() {
		errors["entityType"] = core.InvalidValue.String()
	}

	if entityId < 0 {
		errors["entityId"] = core.InvalidValue.String()
	}

	switch entityType {
	case domain.PlanEntity:
		if ac.PlanRepo.Get(int(entityId)) == nil {
			errors["entityId"] = core.InvalidValue.String()
		}
	default:
		errors["entityId"] = core.InvalidValue.String()
	}

	if !core.IsValidCommentText(text) {
		errors["text"] = core.InvalidValue.String()
	}

	if parentId != 0 && len(title) > 0 {
		errors["title"] = core.InvalidValue.String()
	}

	if !core.IsValidCommentTitle(title) {
		errors["title"] = core.InvalidValue.String()
	}

	if parentId == 0 && len(title) == 0 {
		errors["title"] = core.InvalidValue.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}
