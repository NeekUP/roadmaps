package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type TreeNode struct {
	TopicName  string
	TopicTitle string
	PlanId     int
	PlanTitle  string
	Child      []TreeNode
}

type GetPlanTree interface {
	Do(ctx core.ReqContext, identifiers []int) ([]TreeNode, error)
	DoByTopic(ctx core.ReqContext, name string) ([]TreeNode, error)
}

type getPlanTree struct {
	PlanRepo   core.PlanRepository
	TopicRepo  core.TopicRepository
	UsersPlans core.UsersPlanRepository
	Log        core.AppLogger
}

func NewGetPlanTree(planRepo core.PlanRepository, topics core.TopicRepository, uplans core.UsersPlanRepository, log core.AppLogger) GetPlanTree {
	return &getPlanTree{PlanRepo: planRepo, TopicRepo: topics, UsersPlans: uplans, Log: log}
}

func (this *getPlanTree) Do(ctx core.ReqContext, ids []int) ([]TreeNode, error) {
	appErr := this.validate(ids)
	if appErr != nil {
		this.Log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", appErr.Error(),
		)
		return nil, appErr
	}

	// get plans
	plans := this.PlanRepo.GetList(ids)
	if len(plans) == 0 {
		return nil, core.NewError(core.InvalidRequest)
	}

	result := make([]TreeNode, 0, 0)

	// for every plan
	for i := 0; i < len(plans); i++ {
		// get topic
		t := this.TopicRepo.Get(plans[i].TopicName)
		if t == nil {
			return nil, core.NewError(core.InvalidRequest)
		}
		// create tree node
		parent := TreeNode{
			TopicTitle: t.Title,
			TopicName:  t.Name,
			PlanId:     plans[i].Id,
			PlanTitle:  plans[i].Title,
		}

		userId := ctx.UserId()
		userFavorits := this.getUserFavoritsPlans(userId)

		// for every plan step with topic
		for j := 0; j < len(plans[i].Steps); j++ {
			if plans[i].Steps[j].ReferenceType == domain.TopicReference {
				// get topic
				t := this.TopicRepo.GetById(int(plans[i].Steps[j].ReferenceId))
				if t != nil {
					if userFavorits[t.Name] != 0 {
						plan := this.PlanRepo.Get(userFavorits[t.Name])
						t.Plans = []domain.Plan{*plan}
					} else {
						t.Plans = this.PlanRepo.GetPopularByTopic(t.Name, 1)
					}
					chPlanId := -1
					chPlanTitle := ""

					if len(t.Plans) >= 1 {
						chPlanId = t.Plans[0].Id
						chPlanTitle = t.Plans[0].Title
					}

					// create tree node
					ch := TreeNode{
						TopicTitle: t.Title,
						TopicName:  t.Name,
						PlanId:     chPlanId,
						PlanTitle:  chPlanTitle,
					}
					parent.Child = append(parent.Child, ch)
				}
			}
		}
		result = append(result, parent)
	}

	return result, nil
}

func (this *getPlanTree) getUserFavoritsPlans(userid string) map[string]int {
	userFavorits := make(map[string]int)
	if userid == "" {
		return userFavorits
	}

	userId := userid
	if userId != "" {
		uplans := this.UsersPlans.GetByUser(userId)
		for i := 0; i < len(uplans); i++ {
			userFavorits[uplans[i].TopicName] = uplans[i].PlanId
		}
	}

	return userFavorits
}

func (this *getPlanTree) DoByTopic(ctx core.ReqContext, name string) ([]TreeNode, error) {
	topic := this.TopicRepo.Get(name)
	if topic == nil {
		return nil, core.NewError(core.InvalidRequest)
	}

	up := this.UsersPlans.GetByTopic(ctx.UserId(), topic.Name)
	if up != nil {
		p := this.PlanRepo.Get(up.PlanId)
		topic.Plans = []domain.Plan{*p}
	} else {
		topic.Plans = this.PlanRepo.GetPopularByTopic(topic.Name, 1)
	}

	if len(topic.Plans) == 0 {
		return []TreeNode{TreeNode{
			TopicName:  topic.Name,
			TopicTitle: topic.Title,
		}}, nil
	}

	return this.Do(ctx, []int{topic.Plans[0].Id})
}

func (this *getPlanTree) validate(identifiers []int) *core.AppError {
	errors := make(map[string]string)

	if identifiers == nil {
		errors["id"] = core.InvalidFormat.String()
		return core.ValidationError(errors)
	}

	l := len(identifiers)
	if l == 0 || l > 5 {
		errors["id"] = core.InvalidFormat.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}
	return nil
}
