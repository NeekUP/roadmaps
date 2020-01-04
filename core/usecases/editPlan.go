package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type EditPlan interface {
	Do(ctx core.ReqContext, req EditPlanReq) (bool, error)
}

type editPlan struct {
	PlanRepo core.PlanRepository
	Log      core.AppLogger
}

type EditPlanReq struct {
	Id        int
	TopicName string
	Title     string
	Steps     []PlanStep
}

func NewEditPlan(planRepo core.PlanRepository, log core.AppLogger) EditPlan {
	return &editPlan{PlanRepo: planRepo, Log: log}
}

func (this *editPlan) Do(ctx core.ReqContext, req EditPlanReq) (bool, error) {
	appErr := this.validate(req, ctx.UserId())
	if appErr != nil {
		this.Log.Errorw("Invalid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return false, appErr
	}

	userId := ctx.UserId()
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

	if ok, err := this.PlanRepo.Update(plan); !ok {
		if err != nil {
			this.Log.Errorw("Invalid request",
				"ReqId", ctx.ReqId(),
				"Error", err.Error(),
			)
		}
		return false, err
	}

	return true, nil
}

func (this *editPlan) validate(req EditPlanReq, userId string) *core.AppError {
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

	plan := this.PlanRepo.Get(req.Id)
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
