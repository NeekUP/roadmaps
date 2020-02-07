package tests

import (
	"fmt"
	"github.com/NeekUP/roadmaps/infrastructure"
	"strings"
	"testing"

	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/core/usecases"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/NeekUP/roadmaps/infrastructure/db"
)

func TestAddFirstCommetSuccess(t *testing.T) {

	u := registerUser("TestAddFirstCommetSuccess", "TestAddFirstCommetSuccess@w.ww", "TestAddFirstCommetSuccess")
	if u != nil {
		defer DeleteUser(u.Id)
	}

	addCommentUsecase := usecases.NewAddComment(db.NewCommentsRepository(DB), db.NewPlansRepository(DB), infrastructure.NewChangesCollector(db.NewChangeLogRepository(DB), &appLoggerForTests{}), appLoggerForTests{})

	topic, err := createTopic(u)
	if err != nil {
		t.Errorf("Topic not created: %v", err)
		return
	}
	defer DeleteTopic(topic.Id)

	plan, err := createPlan(topic, u)
	if err != nil {
		t.Errorf("Plan not created: %v", err)
		return
	}
	defer DeletePlan(plan.Id)

	comment, err := addCommentUsecase.Do(newContext(u), domain.PlanEntity, int64(plan.Id), 0, "text", "title")
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	if comment == nil {
		t.Error("Comment is nil")
		return
	}

	if comment.Id == 0 {
		t.Error("Id == 0")
	}
}

func TestAddFirstComment_NotExistsTarget(t *testing.T) {

	u := registerUser("TestAddFirstComment_NotExistsTar", "TestAddFirstComment_NotExistsTarget@w.ww", "TestAddFirstComment_NotExistsTarget")
	if u != nil {
		defer DeleteUser(u.Id)
	}

	addCommentUsecase := usecases.NewAddComment(db.NewCommentsRepository(DB), db.NewPlansRepository(DB), infrastructure.NewChangesCollector(db.NewChangeLogRepository(DB), &appLoggerForTests{}), appLoggerForTests{})

	comment, err := addCommentUsecase.Do(newContext(u), domain.PlanEntity, 1000, 0, "text", "title")
	if err == nil {
		t.Error("Error is null")
		return
	}

	if !strings.Contains(err.Error(), core.InvalidRequest.String()) {
		t.Errorf("Not expexted error: %v", err.Error())
		return
	}

	if comment != nil {
		t.Error("comment created")
	}
}

func TestGetCommentThreads(t *testing.T) {

	u := registerUser("TestGetCommentThreads", "TestGetCommentThreads@w.ww", "TestGetCommentThreads")
	if u != nil {
		defer DeleteUser(u.Id)
	}

	addCommentUsecase := usecases.NewAddComment(db.NewCommentsRepository(DB), db.NewPlansRepository(DB), infrastructure.NewChangesCollector(db.NewChangeLogRepository(DB), &appLoggerForTests{}), appLoggerForTests{})

	topic, err := createTopic(u)
	if err != nil {
		t.Errorf("Topic not created: %v", err)
		return
	}
	defer DeleteTopic(topic.Id)

	plan, err := createPlan(topic, u)
	if err != nil {
		t.Errorf("Plan not created: %v", err)
		return
	}
	defer DeletePlan(plan.Id)

	comment1, err := addCommentUsecase.Do(newContext(u), domain.PlanEntity, int64(plan.Id), 0, "text1", "title1")
	comment2, err := addCommentUsecase.Do(newContext(u), domain.PlanEntity, int64(plan.Id), 0, "text2", "title2")

	subComment1, err := addCommentUsecase.Do(newContext(u), domain.PlanEntity, int64(plan.Id), comment1.Id, "text1", "")
	//subComment2, err := addCommentUsecase.Do(newContext(u), domain.PlanEntity, int64(plan.Id), comment2.Id, "text2", "")

	subComment12, err := addCommentUsecase.Do(newContext(u), domain.PlanEntity, int64(plan.Id), subComment1.Id, "text12", "")

	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	if subComment12 == nil {
		t.Error("Comment is nil")
		return
	}

	if subComment12.Id == 0 {
		t.Error("Id == 0")
	}

	if subComment12.ParentId != subComment1.Id {
		t.Error("Unexpected parentId")
	}

	if subComment12.ThreadId != comment1.Id {
		t.Error("Unexpected threadId")
	}

	getCommentsThreadaUsecase := usecases.NewGetCommentsThreads(db.NewCommentsRepository(DB), db.NewUserRepository(DB), appLoggerForTests{})
	threads, hasmore, err := getCommentsThreadaUsecase.Do(newContext(u), domain.PlanEntity, int64(plan.Id), 10, 0)

	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	if len(threads) != 2 {
		t.Error("Unexpected threads count")
	}

	if hasmore {
		t.Error("Unexpected hasMore value")
	}

	if threads[0].Id != comment1.Id && threads[0].Id != comment2.Id {
		t.Error("Comment1 not in list")
	}

	if threads[1].Id != comment1.Id && threads[1].Id != comment2.Id {
		t.Error("Comment1 not in list")
	}

	if threads[0].Id == threads[1].Id {
		t.Error("wft")
	}
}

func TestWwe(t *testing.T) {
	id := []string{"q", "A"}
	query := "select id, name, normalizedname, email, emailconfirmed, emailconfirmation, img, tokens, rights, password, salt " +
		"FROM users WHERE Id IN ('%s')"
	query = fmt.Sprintf(query, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(id)), "','"), "[]"))
	t.Error(query)
}

func createPlan(topic *domain.Topic, u *domain.User) (*domain.Plan, error) {
	addPlanReq := usecases.AddPlanReq{
		Title:     "TestAddFirstCommetSuccess Plan",
		TopicName: topic.Name,
		Steps: []usecases.PlanStep{
			usecases.PlanStep{
				ReferenceId:   int64(topic.Id),
				ReferenceType: domain.TopicReference,
			},
		},
	}
	addPlanUsecase := usecases.NewAddPlan(db.NewPlansRepository(DB), infrastructure.NewChangesCollector(db.NewChangeLogRepository(DB), &appLoggerForTests{}), &appLoggerForTests{})
	plan, err := addPlanUsecase.Do(newContext(u), addPlanReq)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func createTopic(u *domain.User) (*domain.Topic, error) {
	newTopicUsecase := usecases.NewAddTopic(db.NewTopicRepository(DB), infrastructure.NewChangesCollector(db.NewChangeLogRepository(DB), &appLoggerForTests{}), appLoggerForTests{})
	topic, err := newTopicUsecase.Do(newContext(u), "TestAddFirstCommetSuccess", "", true, []string{})
	if err != nil {
		return nil, err
	}
	return topic, nil
}
