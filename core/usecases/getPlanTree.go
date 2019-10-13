package usecases

import (
	"roadmaps/core"
	"roadmaps/domain"
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
}

type getPlanTree struct {
	PlanRepo core.PlanRepository
	GetTopic GetTopic
	Log      core.AppLogger
}

func NewGetPlanTree(planRepo core.PlanRepository, getTopic GetTopic, log core.AppLogger) GetPlanTree {
	return &getPlanTree{PlanRepo: planRepo, GetTopic: getTopic, Log: log}
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
		t, err := this.GetTopic.Do(ctx, plans[i].TopicName, 1)
		if err != nil {
			return nil, core.NewError(core.InvalidRequest)
		}
		// create tree node
		parent := TreeNode{
			TopicTitle: t.Title,
			TopicName:  t.Name,
			PlanId:     plans[i].Id,
			PlanTitle:  plans[i].Title,
		}

		// for every plan step with topic
		for j := 0; j < len(plans[i].Steps); j++ {
			if plans[i].Steps[j].ReferenceType == domain.TopicReference {
				// get topic
				t, _ := this.GetTopic.DoById(ctx, plans[i].Steps[j].ReferenceId, 1)
				if t != nil {
					chPlanId := -1
					chPlanTitle := ""

					if len(t.Plans) > 0 {
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
