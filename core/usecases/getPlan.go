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
		PlanRepo:   plans,
		StepRepo:   steps,
		UserRepo:   users,
		SourceRepo: sources,
		TopicRepo:  topics,
		//UsersPlansRepo: usersPlans,
		Log: logger,
	}
}

type getPlan struct {
	PlanRepo   core.PlanRepository
	StepRepo   core.StepRepository
	SourceRepo core.SourceRepository
	TopicRepo  core.TopicRepository
	UserRepo   core.UserRepository
	//UsersPlansRepo core.UsersPlanRepository
	Log core.AppLogger
}

func (this *getPlan) Do(ctx core.ReqContext, id int) (*domain.Plan, error) {
	appErr := this.validate(id)
	if appErr != nil {
		this.Log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return nil, appErr
	}

	plan := this.PlanRepo.Get(id)
	if plan != nil {
		plan.Steps = this.StepRepo.GetByPlan(plan.Id)
		plan.Owner = this.UserRepo.Get(plan.OwnerId)
		this.fillSteps(plan)
	}
	return plan, nil
}

func (this *getPlan) fillSteps(plan *domain.Plan) {
	for i := 0; i < len(plan.Steps); i++ {
		if plan.Steps[i].ReferenceType == domain.TopicReference {
			t := this.TopicRepo.GetById(int(plan.Steps[i].ReferenceId))
			if t == nil {
				break
			}
			plan.Steps[i].Source = &domain.Source{
				Id:         int64(t.Id),
				Title:      t.Title,
				Identifier: t.Name,
				Desc:       t.Description,
			}
		} else if plan.Steps[i].ReferenceType == domain.TestReference {
			plan.Steps[i].Source = &domain.Source{
				Title:      "Test",
				Identifier: "Test",
				Desc:       "Not implementer yet",
			}
		} else {
			s := this.SourceRepo.Get(plan.Steps[i].ReferenceId)
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

func (this *getPlan) validate(id int) *core.AppError {
	errors := make(map[string]string)
	if id < 0 {
		errors["id"] = core.InvalidValue.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}
