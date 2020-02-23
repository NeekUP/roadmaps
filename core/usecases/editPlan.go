package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type EditPlan interface {
	Do(ctx core.ReqContext, req EditPlanReq) (bool, error)
}

type editPlan struct {
	planRepo  core.PlanRepository
	log       core.AppLogger
	changeLog core.ChangeLog
}

type EditPlanReq struct {
	Id        int
	TopicName string
	Title     string
	Steps     []PlanStep
}

func NewEditPlan(planRepo core.PlanRepository, changeLog core.ChangeLog, log core.AppLogger) EditPlan {
	return &editPlan{planRepo: planRepo, changeLog: changeLog, log: log}
}

func (usecase *editPlan) Do(ctx core.ReqContext, req EditPlanReq) (bool, error) {
	trace := ctx.StartTrace("editPlan")
	defer ctx.StopTrace(trace)

	old := usecase.planRepo.Get(ctx, req.Id)
	userId := ctx.UserId()
	appErr := usecase.validate(req, userId, old)
	if appErr != nil {
		usecase.log.Errorw("Invalid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return false, appErr
	}

	stepsCount := len(req.Steps)
	steps := make([]domain.Step, 0, stepsCount)

	for i, v := range req.Steps {
		step := domain.Step{
			ReferenceId:   v.ReferenceId,
			ReferenceType: v.ReferenceType,
			Position:      i,
		}
		steps = append(steps, step)
	}

	plan := &domain.Plan{
		Id:        req.Id,
		TopicName: req.TopicName,
		Title:     req.Title,
		OwnerId:   userId,
		Steps:     steps,
	}

	if ok, err := usecase.planRepo.Update(ctx, plan); !ok {
		if err != nil {
			usecase.log.Errorw("Invalid request",
				"ReqId", ctx.ReqId(),
				"Error", err.Error(),
			)
		}
		return false, err
	}

	usecase.changeLog.Edited(domain.PlanEntity, int64(plan.Id), userId, old, plan)
	return true, nil
}

func (usecase *editPlan) validate(req EditPlanReq, userId string, plan *domain.Plan) *core.AppError {
	errors := make(map[string]string)
	if !core.IsValidTopicName(req.TopicName) {
		errors["topic"] = core.InvalidFormat.String()
	}

	if !core.IsValidPlanTitle(req.Title) {
		errors["title"] = core.InvalidFormat.String()
	}

	if len(req.Steps) == 0 {
		errors["steps"] = core.InvalidCount.String()
	}

	if plan == nil {
		errors["id"] = core.NotExists.String()
	}

	if plan.OwnerId != userId {
		errors["id"] = core.AccessDenied.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}

	return nil
}
