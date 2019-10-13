package usecases

import (
	"roadmaps/core"
	"roadmaps/domain"
)

type AddPlan interface {
	Do(ctx core.ReqContext, req AddPlanReq) (*domain.Plan, error)
}

type addPlan struct {
	PlanRepo core.PlanRepository
	Log      core.AppLogger
}

type AddPlanReq struct {
	TopicName string
	Title     string
	Steps     []PlanStep
}

type PlanStep struct {
	ReferenceId   int
	ReferenceType domain.ReferenceType
}

func NewAddPlan(planRepo core.PlanRepository, log core.AppLogger) AddPlan {
	return &addPlan{PlanRepo: planRepo, Log: log}
}

func (this *addPlan) Do(ctx core.ReqContext, req AddPlanReq) (*domain.Plan, error) {

	appErr := this.validate(req)
	if appErr != nil {
		this.Log.Errorw("Invalid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return nil, appErr
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
		TopicName: req.TopicName,
		Title:     req.Title,
		OwnerId:   userId,
		Steps:     steps,
	}

	/*
		Сохранить в транзакции в базу все
		В базе должны быть проверки уникальности и вторичные ключи
	*/
	err := this.PlanRepo.SaveWithSteps(plan)
	if err != nil {
		return nil, core.NewError(core.InvalidRequest)
	}
	return plan, nil
}

func (this *addPlan) validate(req AddPlanReq) *core.AppError {
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

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}
