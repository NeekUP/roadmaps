package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type EditPlan interface {
	Do(ctx core.ReqContext, req EditPlanReq) (bool, error)
}

type editPlan struct {
	planRepo     core.PlanRepository
	sourceRepo   core.SourceRepository
	topicRepo    core.TopicRepository
	projectsRepo core.ProjectsRepository
	log          core.AppLogger
	changeLog    core.ChangeLog
}

type EditPlanReq struct {
	Id        int
	TopicName string
	Title     string
	Steps     []PlanStep
}

func NewEditPlan(planRepo core.PlanRepository, sourceRepo core.SourceRepository, topicRepo core.TopicRepository, projectsRepo core.ProjectsRepository, changeLog core.ChangeLog, log core.AppLogger) EditPlan {
	return &editPlan{planRepo: planRepo,
		sourceRepo:   sourceRepo,
		topicRepo:    topicRepo,
		projectsRepo: projectsRepo,
		changeLog:    changeLog,
		log:          log}
}

func (usecase *editPlan) Do(ctx core.ReqContext, req EditPlanReq) (bool, error) {
	trace := ctx.StartTrace("editPlan")
	defer ctx.StopTrace(trace)

	old := usecase.planRepo.Get(ctx, req.Id)
	userId := ctx.UserId()
	appErr := usecase.validate(ctx, req, userId, old)
	if appErr != nil {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"error", appErr.Error(),
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
			Title:         v.Title,
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
			usecase.log.Errorw("invalid request",
				"reqid", ctx.ReqId(),
				"error", err.Error(),
			)
		}
		return false, err
	}

	usecase.changeLog.Edited(domain.PlanEntity, int64(plan.Id), userId, old, plan)
	return true, nil
}

func (usecase *editPlan) validate(ctx core.ReqContext, req EditPlanReq, userId string, plan *domain.Plan) *core.AppError {
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
			errors["type"] = core.InvalidValue.String()
		}
		if !core.IsValidStepTitle(v.Title) {
			errors["title"] = core.InvalidValue.String()
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
