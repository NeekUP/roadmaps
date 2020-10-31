package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type AddPlan interface {
	Do(ctx core.ReqContext, req AddPlanReq) (*domain.Plan, error)
}

type addPlan struct {
	planRepo     core.PlanRepository
	sourceRepo   core.SourceRepository
	topicRepo    core.TopicRepository
	projectsRepo core.ProjectsRepository
	log          core.AppLogger
	changeLog    core.ChangeLog
}

type AddPlanReq struct {
	TopicName string
	Title     string
	Steps     []PlanStep
}

type PlanStep struct {
	ReferenceId   int64
	ReferenceType domain.ReferenceType
	Title         string
}

func NewAddPlan(planRepo core.PlanRepository, sourceRepo core.SourceRepository, topicRepo core.TopicRepository, projectsRepo core.ProjectsRepository, changeLog core.ChangeLog, log core.AppLogger) AddPlan {
	return &addPlan{planRepo: planRepo,
		sourceRepo:   sourceRepo,
		topicRepo:    topicRepo,
		projectsRepo: projectsRepo,
		changeLog:    changeLog,
		log:          log}
}

func (usecase *addPlan) Do(ctx core.ReqContext, req AddPlanReq) (*domain.Plan, error) {
	trace := ctx.StartTrace("addPlan")
	defer ctx.StopTrace(trace)
	appErr := usecase.validate(ctx, req)
	if appErr != nil {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"error", appErr.Error(),
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
			Title:         v.Title,
		}
		steps = append(steps, step)
	}

	plan := &domain.Plan{
		TopicName: req.TopicName,
		Title:     req.Title,
		OwnerId:   userId,
		Steps:     steps,
	}

	if ok, err := usecase.planRepo.SaveWithSteps(ctx, plan); !ok {
		if err != nil {
			usecase.log.Errorw("invalid request",
				"reqid", ctx.ReqId(),
				"error", err.Error(),
			)
		}
		return nil, err
	}

	usecase.changeLog.Added(domain.PlanEntity, int64(plan.Id), userId)
	return plan, nil
}

func (usecase *addPlan) validate(ctx core.ReqContext, req AddPlanReq) *core.AppError {
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

	for _, v := range req.Steps {
		if v.ReferenceId == 0 {
			errors["source.id"] = core.InvalidValue.String()
		}
		if !core.IsValidReferenceType(v.ReferenceType) {
			errors["source.type"] = core.InvalidValue.String()
		}
		if !core.IsValidStepTitle(v.Title) {
			errors["source.title"] = core.InvalidValue.String()
		}
		switch v.ReferenceType {
		case domain.ResourceReference:
			if usecase.sourceRepo.Get(ctx, v.ReferenceId) == nil {
				errors["source.id"] = core.NotExists.String()
			}
		case domain.ProjectReference:
			if usecase.projectsRepo.Get(ctx, int(v.ReferenceId)) == nil {
				errors["source.id"] = core.NotExists.String()
			}
		case domain.TopicReference:
			if usecase.topicRepo.GetById(ctx, int(v.ReferenceId)) == nil {
				errors["source.id"] = core.NotExists.String()
			}
		}
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}
