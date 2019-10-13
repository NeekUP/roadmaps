package tests

import (
	"roadmaps/core/usecases"
	"roadmaps/domain"
	"roadmaps/infrastructure"
	"roadmaps/infrastructure/db"
	"testing"
)

func TestGetPlanTreeSuccess(t *testing.T) {
	newTopicUsecase := usecases.NewAddTopic(db.NewTopicRepository(nil), log)
	topic1, _ := newTopicUsecase.Do(infrastructure.NewContext(nil), "Topic1", "")

	addPlansReq := []usecases.AddPlanReq{
		usecases.AddPlanReq{
			Title:     "Plan #1 !!!",
			TopicName: topic1.Name,
			Steps: []usecases.PlanStep{
				usecases.PlanStep{
					ReferenceId:   topic1.Id,
					ReferenceType: domain.TopicReference,
				},
				usecases.PlanStep{
					ReferenceId:   topic1.Id,
					ReferenceType: domain.TopicReference,
				},
			},
		},
	}

	plans := make([]domain.Plan, len(addPlansReq), len(addPlansReq))
	usecase := usecases.NewAddPlan(db.NewPlansRepository(nil), &appLoggerForTests{})
	for i, v := range addPlansReq {
		plan, err := usecase.Do(infrastructure.NewContext(nil), v)
		if err != nil {
			t.Error("Fail to create plan")
			return
		}
		plans[i] = *plan
	}

	getTopic := usecases.NewGetTopic(db.NewTopicRepository(nil), db.NewPlansRepository(nil), log)
	getPlanTree := usecases.NewGetPlanTree(db.NewPlansRepository(nil), getTopic, log)
	result, err := getPlanTree.Do(infrastructure.NewContext(nil), []int{plans[0].Id})
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
			t.Errorf("Topic title not expected: %s, expected: %s", v.TopicTitle, plans[i].Title)
		}

		if v.PlanTitle != plans[i].Title {
			t.Errorf("Plan title not expected: %s, expected: %s", v.PlanTitle, plans[i].Title)
		}

		if v.TopicName != topic1.Name {
			t.Errorf("Topic name for main plan not excpected: %s", v.TopicName)
		}

		if len(v.Child) == 0 {
			t.Error("Child plans is empty")
		}

		for j, s := range v.Child {
			if s.PlanId == 0 {
				t.Errorf("plan id is 0")
			}

			if s.TopicTitle != topic1.Title {
				t.Errorf("Topic title not expected: %s, expected: %s", s.TopicTitle, plans[j].Title)
			}

			if s.PlanTitle != plans[i].Title {
				t.Errorf("Plan title not expected: %s, expected: %s", s.PlanTitle, plans[j].Title)
			}

			if s.TopicName != topic1.Name {
				t.Errorf("Topic name for main plan not excpected: %s", s.TopicName)
			}
		}
	}
}
