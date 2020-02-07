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
	//usersPlans core.UsersPlanRepository,
	logger core.AppLogger) GetPlan {
	return &getPlan{
		planRepo:   plans,
		stepRepo:   steps,
		userRepo:   users,
		sourceRepo: sources,
		topicRepo:  topics,
		//UsersPlansRepo: usersPlans,
		log: logger,
	}
}

type getPlan struct {
	planRepo   core.PlanRepository
	stepRepo   core.StepRepository
	sourceRepo core.SourceRepository
	topicRepo  core.TopicRepository
	userRepo   core.UserRepository
	log        core.AppLogger
}

func (usecase *getPlan) Do(ctx core.ReqContext, id int) (*domain.Plan, error) {
	appErr := usecase.validate(id)
	if appErr != nil {
		usecase.log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return nil, appErr
	}

	plan := usecase.planRepo.Get(id)
	if plan != nil {
		plan.Steps = usecase.stepRepo.GetByPlan(plan.Id)
		plan.Owner = usecase.userRepo.Get(plan.OwnerId)
		usecase.fillSteps(plan)
	}
	return plan, nil
}

func (usecase *getPlan) fillSteps(plan *domain.Plan) {
	for i := 0; i < len(plan.Steps); i++ {
		if plan.Steps[i].ReferenceType == domain.TopicReference {
			t := usecase.topicRepo.GetById(int(plan.Steps[i].ReferenceId))
			if t == nil {
				break
			}
			plan.Steps[i].Source = t
		} else if plan.Steps[i].ReferenceType == domain.TestReference {
			plan.Steps[i].Source = &domain.Source{
				Title:      "Test",
				Identifier: "Test",
				Desc:       "Not implementer yet",
			}
		} else {
			s := usecase.sourceRepo.Get(plan.Steps[i].ReferenceId)
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
