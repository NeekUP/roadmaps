package infrastructure

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/r3labs/diff"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	Add    = 1
	Edit   = 2
	Delete = 3
)

type ChangesCollector struct {
	changeLogRepo core.ChangeLogRepository
	log           core.AppLogger
}

func NewChangesCollector(repo core.ChangeLogRepository, logger core.AppLogger) core.ChangeLog {
	return &ChangesCollector{changeLogRepo: repo, log: logger}
}

func (collector *ChangesCollector) Added(entityType domain.EntityType, entityId int64, userId string) {
	action, err := getActionType(entityType, Add)
	if err != nil {
		collector.log.Errorw("Changes not logged", "error", err, "entityType", entityType, "entityId", entityId, "userId", userId, "action", Add)
		return
	}

	collector.saveRecord(action, entityType, entityId, userId, "")
}

func (collector *ChangesCollector) Edited(entityType domain.EntityType, entityId int64, userId string, before interface{}, after interface{}) {
	action, err := getActionType(entityType, Edit)
	if err != nil {
		collector.log.Errorw("Changes not logged", "error", err, "entityType", entityType, "entityId", entityId, "userId", userId, "action", Add)
		return
	}

	difference, err := diffEntities(before, after)
	if err != nil {
		collector.log.Errorw("Fail to get deff between entities", "error", err, "entityType", entityType, "entityId", entityId, "userId", userId, "action", Add)
		return
	}

	collector.saveRecord(action, entityType, entityId, userId, string(difference))
}

func (collector *ChangesCollector) Deleted(entityType domain.EntityType, entityId int64, userId string) {
	action, err := getActionType(entityType, Add)
	if err != nil {
		collector.log.Errorw("Changes not logged", "error", err, "entityType", entityType, "entityId", entityId, "userId", userId, "action", Add)
		return
	}
	collector.saveRecord(action, entityType, entityId, userId, "")
}

func (collector *ChangesCollector) saveRecord(action domain.ChangeType, entityType domain.EntityType, entityId int64, userId string, diff string) {
	record := &domain.ChangeLogRecord{
		Action:     action,
		UserId:     userId,
		EntityType: entityType,
		EntityId:   entityId,
		Diff:       diff,
	}

	if !collector.changeLogRepo.Add(record) {
		collector.log.Errorw("Changes not saved into db", "entityType", entityType, "entityId", entityId, "userId", userId, "action", Add)
	}
}

func getActionType(eType domain.EntityType, action int) (domain.ChangeType, error) {
	switch eType {
	case domain.PlanEntity:
		switch action {
		case Add:
			return domain.AddPlan, nil
		case Edit:
			return domain.EditPlan, nil
		case Delete:
			return domain.DeletePlan, nil
		}
	case domain.TopicEntity:
		switch action {
		case Add:
			return domain.AddTopic, nil
		case Edit:
			return domain.EditTopic, nil
		case Delete:
			return domain.DeleteTopic, nil
		}
	case domain.ProjectEntity:
		switch action {
		case Add:
			return domain.AddProject, nil
		case Edit:
			return domain.EditProject, nil
		case Delete:
			return domain.DeleteProject, nil
		}
	case domain.ResourceEntity:
		switch action {
		case Add:
			return domain.AddResource, nil
		case Edit:
			return domain.EditResource, nil
		case Delete:
			return domain.DeleteResource, nil
		}
	case domain.CommentEntity:
		switch action {
		case Add:
			return domain.AddComment, nil
		case Edit:
			return domain.EditComment, nil
		case Delete:
			return domain.DeleteComment, nil
		}
	case domain.UserEntity:
		switch action {
		case Add:
			return domain.AddUser, nil
		case Edit:
			return domain.EditUser, nil
		case Delete:
			return domain.DeleteUser, nil
		}
	}
	return 0, fmt.Errorf("unknown entityType: %v or action: %v", eType, action)
}

func diffEntities(before interface{}, after interface{}) ([]byte, error) {
	var err error
	var difference []byte

	switch before.(type) {
	case *domain.Plan:
		difference, err = diffPlans(before.(*domain.Plan), after.(*domain.Plan))
	case *domain.Topic:
		difference, err = diffTopic(before.(*domain.Topic), after.(*domain.Topic))
	// case *domain.Project:
	// 	difference, err = diffProject(before.(*domain.Project), after.(*domain.Project))
	case *domain.Comment:
		difference, err = diffComment(before.(*domain.Comment), after.(*domain.Comment))
	case *domain.User:
		difference, err = diffUser(before.(*domain.User), after.(*domain.User))
	default:
		return nil, fmt.Errorf("Unexpected type. before: %#v after: %#v", before, after)
	}

	return difference, err
}

func diffPlans(before *domain.Plan, after *domain.Plan) ([]byte, error) {
	dmp := diffmatchpatch.New()

	sort.Slice(before.Steps, func(i, j int) bool { return before.Steps[i].Position < before.Steps[j].Position })
	sort.Slice(after.Steps, func(i, j int) bool { return after.Steps[i].Position < after.Steps[j].Position })

	stepsChanges, err := diff.Diff(before.Steps, after.Steps)
	if err != nil {
		return nil, err
	}

	stepsDiff := make([]sliceDiff, len(stepsChanges))
	for i, v := range stepsChanges {
		stepsDiff[i] = sliceDiff{
			Position: i,
			Type:     v.Type,
			Attr:     v.Path[1],
			From:     v.From,
			To:       v.To,
		}
	}

	d := &planDiff{
		Title:   diffStrings(dmp, before.Title, after.Title),
		OwnerId: diffStrings(dmp, before.OwnerId, after.OwnerId),
		Steps:   stepsDiff,
	}

	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func diffTopic(before *domain.Topic, after *domain.Topic) ([]byte, error) {
	dmp := diffmatchpatch.New()

	SortTopicTags(before.Tags)
	SortTopicTags(after.Tags)

	tagsChanges, err := diff.Diff(before.Tags, after.Tags)
	tagsDiff := make([]sliceDiff, len(tagsChanges))

	for i, v := range tagsChanges {
		tagsDiff[i] = sliceDiff{
			Position: i,
			Type:     v.Type,
			Attr:     v.Path[1],
			From:     v.From,
			To:       v.To,
		}
	}
	if err != nil {
		return nil, err
	}
	d := &topicDiff{
		Name:        diffStrings(dmp, before.Name, after.Name),
		Title:       diffStrings(dmp, before.Title, after.Title),
		Description: diffStrings(dmp, before.Description, after.Description),
		Creator:     diffStrings(dmp, before.Creator, after.Creator),
		IsTag:       diffStrings(dmp, strconv.FormatBool(before.IsTag), strconv.FormatBool(after.IsTag)),
		Tags:        tagsDiff,
	}

	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func SortTopicTags(tags []domain.TopicTag) {
	bt := make([]string, len(tags))
	bm := make(map[string]domain.TopicTag)

	for i, v := range tags {
		bm[v.Name] = v
		bt[i] = v.Name
	}
	sort.Strings(bt)
	for i, v := range bt {
		tags[i] = bm[v]
	}
}

func diffComment(before *domain.Comment, after *domain.Comment) ([]byte, error) {
	dmp := diffmatchpatch.New()

	d := &commentDiff{
		Text:    diffStrings(dmp, before.Text, after.Text),
		Title:   diffStrings(dmp, before.Title, after.Title),
		Deleted: diffStrings(dmp, strconv.FormatBool(before.Deleted), strconv.FormatBool(after.Deleted)),
	}

	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func diffUser(before *domain.User, after *domain.User) ([]byte, error) {
	dmp := diffmatchpatch.New()

	d := &userDiff{
		Name:   diffStrings(dmp, before.Name, after.Name),
		Email:  diffStrings(dmp, before.Email, after.Email),
		Img:    diffStrings(dmp, before.Img, after.Img),
		Rights: diffStrings(dmp, strconv.Itoa(int(before.Rights)), strconv.Itoa(int(after.Rights))),
	}

	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func diffStrings(dmp *diffmatchpatch.DiffMatchPatch, one, two string) string {
	diffs := dmp.DiffMain(one, two, false)
	return dmp.DiffPrettyText(diffs)
}

type planDiff struct {
	Title   string
	OwnerId string
	Steps   []sliceDiff
}

type sliceDiff struct {
	Position int
	Type     string
	Attr     string
	From     interface{}
	To       interface{}
}

type topicDiff struct {
	Name        string
	Title       string
	Description string
	Creator     string
	IsTag       string
	Tags        []sliceDiff
}

type commentDiff struct {
	Text    string
	Title   string
	Deleted string
}

type userDiff struct {
	Name   string
	Email  string
	Img    string
	Rights string
}
