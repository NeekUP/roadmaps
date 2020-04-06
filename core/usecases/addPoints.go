package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type AddVote interface {
	Do(ctx core.ReqContext, entityType domain.EntityType, id int64, value int) (bool, error)
}

type addVote struct {
	pointsRepo core.PointsRepository
	log        core.AppLogger
}

func NewAddPoints(pointsRepo core.PointsRepository, log core.AppLogger) AddVote {
	return &addVote{pointsRepo: pointsRepo, log: log}
}

func (usecase addVote) Do(ctx core.ReqContext, entityType domain.EntityType, id int64, value int) (bool, error) {
	trace := ctx.StartTrace("addVote")
	defer ctx.StopTrace(trace)
	appErr := usecase.validate(entityType, id, value)
	if appErr != nil {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"error", appErr.Error(),
		)
		return false, appErr
	}

	result := usecase.pointsRepo.Add(ctx, entityType, id, ctx.UserId(), value)
	return result, nil
}

func (usecase addVote) validate(entityType domain.EntityType, id int64, value int) *core.AppError {
	errors := make(map[string]string)

	if !core.IsValidEntityType(entityType) {
		errors["type"] = core.InvalidValue.String()
	}

	if id <= 0 {
		errors["id"] = core.InvalidValue.String()
	}

	if value <= 0 || value > 10 {
		errors["value"] = core.InvalidValue.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}
