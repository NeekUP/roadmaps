package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type RemovePlan interface {
	Do(ctx core.ReqContext, id int) (bool, error)
}

type removePlan struct {
	repo      core.PlanRepository
	log       core.AppLogger
	changeLog core.ChangeLog
}

func NewRemovePlan(planRepo core.PlanRepository, changeLog core.ChangeLog, log core.AppLogger) RemovePlan {
	return &removePlan{repo: planRepo, changeLog: changeLog, log: log}
}

func (usecase *removePlan) Do(ctx core.ReqContext, id int) (bool, error) {
	plan := usecase.repo.Get(id)
	userId := ctx.UserId()
	appErr := usecase.validate(id, userId, plan)
	if appErr != nil {
		usecase.log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return false, appErr
	}

	deleted, err := usecase.repo.Delete(id)
	if err != nil {
		return false, err
	}

	if deleted {
		usecase.changeLog.Deleted(domain.PlanEntity, int64(id), userId)
	}
	return true, nil
}

func (usecase *removePlan) validate(id int, userid string, plan *domain.Plan) *core.AppError {
	errors := make(map[string]string)
	if id <= 0 {
		errors["id"] = core.InvalidFormat.String()
	}

	if plan == nil {
		errors["id"] = core.NotExists.String()
	}

	if plan.OwnerId != userid {
		errors["id"] = core.AccessDenied.String()
	}

	return nil
}
