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
	commentsRepo core.CommentsRepository
	planRepo     core.PlanRepository
	log          core.AppLogger
	changeLog    core.ChangeLog
}

func NewAddComment(commentsRepo core.CommentsRepository, planRepo core.PlanRepository, changeLog core.ChangeLog, log core.AppLogger) AddComment {
	return &addComment{commentsRepo: commentsRepo, planRepo: planRepo, changeLog: changeLog, log: log}
}

func (usecase *addComment) Do(ctx core.ReqContext, entityType domain.EntityType, entityId int64, parentId int64, text string, title string) (*domain.Comment, error) {
	trace := ctx.StartTrace("addComment")
	defer ctx.StopTrace(trace)
	appErr := usecase.validate(ctx, entityType, entityId, parentId, text, title)
	if appErr != nil {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"error", appErr.Error(),
		)
		return nil, appErr
	}

	var threadId int64 = 0
	if parentId > 0 {
		parent := usecase.commentsRepo.Get(ctx, parentId)
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
	userId := ctx.UserId()
	comment := &domain.Comment{
		EntityType: entityType,
		EntityId:   entityId,
		ThreadId:   threadId,
		ParentId:   parentId,
		UserId:     userId,
		Text:       text,
		Title:      title,
		Date:       time.Now().UTC(),
		Deleted:    false,
	}

	if ok, err := usecase.commentsRepo.Add(ctx, comment); !ok {
		if err != nil {
			usecase.log.Errorw("invalid request",
				"reqid", ctx.ReqId(),
				"error", err.Error(),
			)
		}
		return nil, err
	}

	usecase.changeLog.Added(domain.CommentEntity, comment.Id, userId)
	return comment, nil
}

func (ac *addComment) validate(ctx core.ReqContext, entityType domain.EntityType, entityId int64, parentId int64, text string, title string) *core.AppError {
	errors := make(map[string]string)
	if !entityType.IsValid() {
		errors["entityType"] = core.InvalidValue.String()
	}

	if entityId < 0 {
		errors["entityId"] = core.InvalidValue.String()
	}

	switch entityType {
	case domain.PlanEntity:
		if ac.planRepo.Get(ctx, int(entityId)) == nil {
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
