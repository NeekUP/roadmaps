package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type GetPoints interface {
	Do(ctx core.ReqContext, entityType domain.EntityType, id int64) (*domain.Points, error)
}

type getPoints struct {
	pointsRepo core.PointsRepository
	log        core.AppLogger
}

func NewGetPoints(pointsRepo core.PointsRepository, log core.AppLogger) GetPoints {
	return &getPoints{pointsRepo: pointsRepo, log: log}
}

func (usecase *getPoints) Do(ctx core.ReqContext, entityType domain.EntityType, id int64) (*domain.Points, error) {
	trace := ctx.StartTrace("getPoints")
	defer ctx.StopTrace(trace)

	appErr := usecase.validate(entityType, id)
	if appErr != nil {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"error", appErr.Error(),
		)
		return nil, appErr
	}

	points := usecase.pointsRepo.Get(ctx, ctx.UserId(), entityType, id)
	if points != nil {
		return points, nil
	}
	return nil, core.NewError(core.NotExists)
}

func (usecase *getPoints) validate(entityType domain.EntityType, id int64) *core.AppError {
	errors := make(map[string]string)

	if !core.IsValidEntityType(entityType) {
		errors["type"] = core.InvalidValue.String()
	}

	if id <= 0 {
		errors["id"] = core.InvalidValue.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}
