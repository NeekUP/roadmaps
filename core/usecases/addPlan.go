package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type AddPlan interface {
	Do(ctx core.ReqContext, req AddPlanReq) (*domain.Plan, error)
}

type addPlan struct {
	planRepo  core.PlanRepository
	log       core.AppLogger
	changeLog core.ChangeLog
}

type AddPlanReq struct {
	TopicName string
	Title     string
	Steps     []PlanStep
}

type PlanStep struct {
	ReferenceId   int64
	ReferenceType domain.ReferenceType
}

func NewAddPlan(planRepo core.PlanRepository, changeLog core.ChangeLog, log core.AppLogger) AddPlan {
	return &addPlan{planRepo: planRepo, changeLog: changeLog, log: log}
}

func (usecase *addPlan) Do(ctx core.ReqContext, req AddPlanReq) (*domain.Plan, error) {

	appErr := usecase.validate(req)
	if appErr != nil {
		usecase.log.Errorw("Invalid request",
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

	if ok, err := usecase.planRepo.SaveWithSteps(plan); !ok {
		if err != nil {
			usecase.log.Errorw("Invalid request",
				"ReqId", ctx.ReqId(),
				"Error", err.Error(),
			)
		}
		return nil, err
	}

	usecase.changeLog.Added(domain.PlanEntity, int64(plan.Id), userId)
	return plan, nil
}

func (usecase *addPlan) validate(req AddPlanReq) *core.AppError {
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
