package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type GetPointsList interface {
	Do(ctx core.ReqContext, entityType domain.EntityType, id []int64) ([]domain.Points, error)
}

type getPointsList struct {
	pointsRepo core.PointsRepository
	log        core.AppLogger
}

func NewGetPointsList(pointsRepo core.PointsRepository, log core.AppLogger) GetPointsList {
	return &getPointsList{pointsRepo: pointsRepo, log: log}
}

func (usecase *getPointsList) Do(ctx core.ReqContext, entityType domain.EntityType, idList []int64) ([]domain.Points, error) {
	trace := ctx.StartTrace("getPointsList")
	defer ctx.StopTrace(trace)

	appErr := usecase.validate(entityType, idList)
	if appErr != nil {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"error", appErr.Error(),
		)
		return []domain.Points{}, appErr
	}

	points := usecase.pointsRepo.GetList(ctx, ctx.UserId(), entityType, idList)
	if points != nil {
		return points, nil
	}
	return []domain.Points{}, core.NewError(core.NotExists)
}

func (usecase *getPointsList) validate(entityType domain.EntityType, idList []int64) *core.AppError {
	errors := make(map[string]string)

	if !core.IsValidEntityType(entityType) {
		errors["type"] = core.InvalidValue.String()
	}

	for _, id := range idList {
		if id <= 0 {
			errors["id"] = core.InvalidValue.String()
			break
		}
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}
