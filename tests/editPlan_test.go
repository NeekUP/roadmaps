package tests

import (
	"github.com/NeekUP/roadmaps/core/usecases"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/NeekUP/roadmaps/infrastructure/db"
	"testing"
)

func TestEditPlanSuccess(t *testing.T) {
	user := registerUser("TestEditPlanSuccess", "TestEditPlanSuccess@w.ww", "TestEditPlanSuccess")
	if user != nil {
		defer DeleteUser(user.Id)
	}
	newTopicUsecase := usecases.NewAddTopic(db.NewTopicRepository(DB), log)
	topic1, err := newTopicUsecase.Do(newContext(user), "Add Plan", "", true, []string{})
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
		plan, err := usecase.Do(newContext(user), v)

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

		updatePlan(user, plan, t)

		newSteps := db.NewStepsRepository(DB).GetByPlan(plan.Id)
		if len(newSteps) != 1 {
			t.Errorf("unexpected steps count after plan update")
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

func updatePlan(u *domain.User, plan *domain.Plan, t *testing.T) {
	editPlanUsecase := usecases.NewEditPlan(db.NewPlansRepository(DB), log)
	modified, err := editPlanUsecase.Do(newContext(u), usecases.EditPlanReq{
		Id:        plan.Id,
		TopicName: plan.TopicName,
		Title:     "new title",
		Steps: []usecases.PlanStep{
			usecases.PlanStep{
				ReferenceId:   plan.Steps[0].Id,
				ReferenceType: plan.Steps[0].ReferenceType,
			},
		},
	})

	if !modified {
		t.Errorf("Plan not modified")
	}

	if err != nil {
		t.Errorf("Edit plan ends with error: %s", err.Error())
	}
}
