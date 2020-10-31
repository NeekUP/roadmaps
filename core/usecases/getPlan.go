package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type GetPlan interface {
	Do(ctx core.ReqContext, id int) (*domain.Plan, error)
}

func NewGetPlan(plans core.PlanRepository,
	users core.UserRepository,
	steps core.StepRepository,
	sources core.SourceRepository,
	topics core.TopicRepository,
	//projectRepo core.ProjectsRepository,
	logger core.AppLogger) GetPlan {
	return &getPlan{
		planRepo:   plans,
		stepRepo:   steps,
		userRepo:   users,
		sourceRepo: sources,
		topicRepo:  topics,
		//projectRepo: projectRepo,
		log: logger,
	}
}

type getPlan struct {
	planRepo   core.PlanRepository
	stepRepo   core.StepRepository
	sourceRepo core.SourceRepository
	topicRepo  core.TopicRepository
	userRepo   core.UserRepository
	//projectRepo core.ProjectsRepository
	log core.AppLogger
}

func (usecase *getPlan) Do(ctx core.ReqContext, id int) (*domain.Plan, error) {
	trace := ctx.StartTrace("getPlan")
	defer ctx.StopTrace(trace)

	appErr := usecase.validate(id)
	if appErr != nil {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"error", appErr.Error(),
		)
		return nil, appErr
	}

	plan := usecase.planRepo.GetWithDraft(ctx, id, ctx.UserId())
	if plan != nil {
		plan.Steps = usecase.stepRepo.GetByPlan(ctx, plan.Id)
		plan.Owner = usecase.userRepo.Get(ctx, plan.OwnerId)
		usecase.fillSteps(ctx, plan)
		return plan, nil
	}

	return nil, core.NewError(core.NotExists)
}

func (usecase *getPlan) fillSteps(ctx core.ReqContext, plan *domain.Plan) {
	for i := 0; i < len(plan.Steps); i++ {
		if plan.Steps[i].ReferenceType == domain.TopicReference {
			t := usecase.topicRepo.GetById(ctx, int(plan.Steps[i].ReferenceId))
			plan.Steps[i].Source = t
		} else if plan.Steps[i].ReferenceType == domain.ProjectReference {
			//p := usecase.projectRepo.Get(ctx, int(plan.Steps[i].ReferenceId))
			//plan.Steps[i].Source = p
		} else {
			s := usecase.sourceRepo.Get(ctx, plan.Steps[i].ReferenceId)
			if s == nil {
				break
			}
			plan.Steps[i].Source = &domain.Source{
				Id:                   s.Id,
				Title:                s.Title,
				Identifier:           s.Identifier,
				NormalizedIdentifier: s.NormalizedIdentifier,
				Type:                 s.Type,
				Img:                  s.Img,
				Properties:           s.Properties,
				Desc:                 s.Desc,
			}
		}
	}
}

func (usecase *getPlan) validate(id int) *core.AppError {
	errors := make(map[string]string)
	if id < 0 {
		errors["id"] = core.InvalidValue.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}
