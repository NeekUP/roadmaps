package tests

import (
	"github.com/NeekUP/roadmaps/core/usecases"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/NeekUP/roadmaps/infrastructure/db"
	"testing"
)

func TestGetPlanTreeSuccess(t *testing.T) {
	u := registerUser("TestAddTopicSuccess", "TestAddTopicSuccess@w.ww", "TestAddTopicSuccess")
	if u != nil {
		defer DeleteUser(u.Id)
	}
	newTopicUsecase := usecases.NewAddTopic(db.NewTopicRepository(DB), log)
	topic1, _ := newTopicUsecase.Do(newContext(u), "Topic1", "")
	if topic1 != nil {
		defer DeleteTopic(topic1.Id)
	}

	addPlansReq := []usecases.AddPlanReq{
		usecases.AddPlanReq{
			Title:     "Plan #1 !!!",
			TopicName: topic1.Name,
			Steps: []usecases.PlanStep{
				usecases.PlanStep{
					ReferenceId:   int64(topic1.Id),
					ReferenceType: domain.TopicReference,
				},
				usecases.PlanStep{
					ReferenceId:   int64(topic1.Id),
					ReferenceType: domain.TopicReference,
				},
			},
		},
	}

	plans := make([]domain.Plan, len(addPlansReq), len(addPlansReq))
	usecase := usecases.NewAddPlan(db.NewPlansRepository(DB), &appLoggerForTests{})
	for i, v := range addPlansReq {
		plan, err := usecase.Do(newContext(u), v)
		if err != nil {
			t.Error("Fail to create plan")
			return
		} else {
			defer DeletePlan(plan.Id)
		}
		plans[i] = *plan
	}

	getPlanTree := usecases.NewGetPlanTree(db.NewPlansRepository(DB), db.NewTopicRepository(DB), db.NewUsersPlanRepository(DB), log)
	result, err := getPlanTree.Do(newContext(u), []int{plans[0].Id})
	if err != nil {
		t.Errorf("Error while getting plan tree: %s", err.Error())
		return
	}

	if result == nil {
		t.Error("Results is nil")
		return
	}

	for i, v := range result {
		if v.PlanId == 0 {
			t.Errorf("plan id is 0")
		}

		if v.TopicTitle != topic1.Title {
			t.Errorf("TopicName title not expected: %s, expected: %s", v.TopicTitle, plans[i].Title)
		}

		if v.PlanTitle != plans[i].Title {
			t.Errorf("Plan title not expected: %s, expected: %s", v.PlanTitle, plans[i].Title)
		}

		if v.TopicName != topic1.Name {
			t.Errorf("TopicName name for main plan not excpected: %s", v.TopicName)
		}

		if len(v.Child) == 0 {
			t.Error("Child plans is empty")
		}

		for j, s := range v.Child {
			if s.PlanId == 0 {
				t.Errorf("plan id is 0")
			}

			if s.TopicTitle != topic1.Title {
				t.Errorf("TopicName title not expected: %s, expected: %s", s.TopicTitle, plans[j].Title)
			}

			if s.PlanTitle != plans[i].Title {
				t.Errorf("Plan title not expected: %s, expected: %s", s.PlanTitle, plans[j].Title)
			}

			if s.TopicName != topic1.Name {
				t.Errorf("TopicName name for main plan not excpected: %s", s.TopicName)
			}
		}
	}
}
