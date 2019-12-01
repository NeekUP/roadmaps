package tests

import (
	"github.com/NeekUP/roadmaps/core/usecases"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/NeekUP/roadmaps/infrastructure/db"
	"testing"
)

func TestAddPlanSuccess(t *testing.T) {
	u := registerUser("TestAddTopicSuccess", "TestAddTopicSuccess@w.ww", "TestAddTopicSuccess")
	if u != nil {
		defer DeleteUser(u.Id)
	}
	newTopicUsecase := usecases.NewAddTopic(db.NewTopicRepository(DB), log)
	topic1, err := newTopicUsecase.Do(newContext(u), "Add Plan", "")
	if err != nil {
		t.Errorf("Topic not created: %s", err.Error())
		return
	}
	defer DeleteTopic(topic1.Id)

	plans := []usecases.AddPlanReq{
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

	usecase := usecases.NewAddPlan(db.NewPlansRepository(DB), &appLoggerForTests{})
	for _, v := range plans {
		plan, err := usecase.Do(newContext(u), v)

		if err != nil {
			t.Errorf("Plan not saved: %s", err.Error())
			return
		} else {
			defer DeletePlan(plan.Id)
		}

		if plan.Title != v.Title {
			t.Error("Plan title has missing")
		}

		if plan.Id == 0 {
			t.Error("Plan id not defined")
		}

		if len(plan.Steps) != len(v.Steps) {
			t.Errorf("Steps count not expected: %d", len(plan.Steps))
		}

		for pos, step := range plan.Steps {
			defer DeleteStep(step.Id)
			if step.Position != pos {
				t.Errorf("Step has wrong position: %d", step.Position)
			}

			if step.PlanId != plan.Id {
				t.Errorf("Step.PlanId not equals Plan.Id: %d", step.PlanId)
			}
		}
	}
}
