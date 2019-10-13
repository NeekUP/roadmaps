package tests

import (
	"roadmaps/core/usecases"
	"roadmaps/domain"
	"roadmaps/infrastructure"
	"roadmaps/infrastructure/db"
	"testing"
)

func TestAddPlanSuccess(t *testing.T) {
	newTopicUsecase := usecases.NewAddTopic(db.NewTopicRepository(nil), log)
	topic1, _ := newTopicUsecase.Do(infrastructure.NewContext(nil), "Add Plan", "")

	plans := []usecases.AddPlanReq{
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

	usecase := usecases.NewAddPlan(db.NewPlansRepository(nil), &appLoggerForTests{})
	for _, v := range plans {
		plan, err := usecase.Do(infrastructure.NewContext(nil), v)

		if err != nil {
			t.Errorf("Plan not saved: %s", err.Error())
			return
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
			if step.Position != pos {
				t.Errorf("Step has wrong position: %d", step.Position)
			}

			if step.PlanId != plan.Id {
				t.Errorf("Step.PlanId not equals Plan.Id: %d", step.PlanId)
			}
		}
	}
}
