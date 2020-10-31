package db

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/NeekUP/roadmaps/domain"
)

/*
	Topic
 ******************/
type TopicDBO struct {
	Id          int
	Name        string
	Title       string
	Description sql.NullString
	Creator     string
	Tags        []string
	IsTag       bool
}

func (dbo *TopicDBO) ToTopic(tags []domain.TopicTag) *domain.Topic {
	t := &domain.Topic{
		Id:          dbo.Id,
		Name:        dbo.Name,
		Title:       dbo.Title,
		Description: dbo.Description.String,
		Creator:     dbo.Creator,
		Tags:        tags,
		IsTag:       dbo.IsTag,
	}
	return t
}

func (dbo *TopicDBO) FromTopic(d *domain.Topic) {
	dbo.Id = d.Id
	dbo.Name = d.Name
	dbo.Title = d.Title
	dbo.Description = ToNullString(d.Description)
	dbo.Creator = d.Creator
	dbo.IsTag = d.IsTag
	dbo.Tags = make([]string, len(d.Tags))

	for i, tag := range d.Tags {
		dbo.Tags[i] = tag.Name
	}
}

/*
	Plan
 ******************/
type PlanDBO struct {
	Id        int
	Title     string
	TopicName string
	OwnerId   string
	IsDraft   bool
}

func (dbo *PlanDBO) ToPlan() *domain.Plan {
	return &domain.Plan{
		Id:        dbo.Id,
		Title:     dbo.Title,
		TopicName: dbo.TopicName,
		OwnerId:   dbo.OwnerId,
		IsDraft:   dbo.IsDraft,
	}
}

func (dbo *PlanDBO) FromPlan(plan *domain.Plan) {
	dbo.Id = plan.Id
	dbo.Title = plan.Title
	dbo.TopicName = plan.TopicName
	dbo.OwnerId = plan.OwnerId
	dbo.IsDraft = plan.IsDraft
}

/*
	Source
 ******************/
type SourceDBO struct {
	Id                   int64
	Title                string
	Identifier           string
	NormalizedIdentifier string
	Type                 string
	Properties           sql.NullString
	Img                  sql.NullString
	Desc                 sql.NullString
}

func (dbo *SourceDBO) ToSource() *domain.Source {
	return &domain.Source{
		Id:                   dbo.Id,
		Title:                dbo.Title,
		Identifier:           dbo.Identifier,
		NormalizedIdentifier: dbo.NormalizedIdentifier,
		Type:                 domain.SourceType(dbo.Type),
		Properties:           dbo.Properties.String,
		Img:                  dbo.Img.String,
		Desc:                 dbo.Desc.String,
	}
}

func (dbo *SourceDBO) FromSource(s *domain.Source) {
	dbo.Id = s.Id
	dbo.Title = s.Title
	dbo.Type = string(s.Type)
	dbo.Identifier = s.Identifier
	dbo.NormalizedIdentifier = s.NormalizedIdentifier
	dbo.Desc = ToNullString(s.Desc)
	dbo.Img = ToNullString(s.Img)
	dbo.Properties = ToNullString(s.Properties)
}

/*
	Step
 ******************/
type StepDBO struct {
	Id            int64
	PlanId        int
	ReferenceId   int64
	ReferenceType string
	Position      int
	Title         string
}

func (dbo *StepDBO) FromStep(step *domain.Step) {
	dbo.Id = step.Id
	dbo.PlanId = step.PlanId
	dbo.Position = step.Position
	dbo.ReferenceId = step.ReferenceId
	dbo.ReferenceType = string(step.ReferenceType)
	dbo.Title = step.Title
}

func (dbo *StepDBO) ToStep() *domain.Step {
	return &domain.Step{
		Id:            dbo.Id,
		PlanId:        dbo.PlanId,
		ReferenceId:   dbo.ReferenceId,
		ReferenceType: domain.ReferenceType(dbo.ReferenceType),
		Position:      dbo.Position,
		Title:         dbo.Title,
	}
}

/*
	User
 ******************/
type UserDBO struct {
	Id                string
	Name              string
	NormalizedName    string
	Email             string
	EmailConfirmed    bool
	EmailConfirmation string
	Img               sql.NullString
	Tokens            sql.NullString
	Rights            int
	Pass              []byte
	Salt              []byte
}

func (dbo *UserDBO) ToUser() *domain.User {
	tokens := make([]domain.UserToken, 0)
	json.Unmarshal([]byte(dbo.Tokens.String), &tokens)

	return &domain.User{
		Id:                dbo.Id,
		Name:              dbo.Name,
		NormalizedName:    dbo.NormalizedName,
		Email:             dbo.Email,
		EmailConfirmed:    dbo.EmailConfirmed,
		EmailConfirmation: dbo.EmailConfirmation,
		Img:               dbo.Img.String,
		Tokens:            tokens,
		Rights:            domain.Rights(dbo.Rights),
		Pass:              dbo.Pass,
		Salt:              dbo.Salt,
	}
}

func (dbo *UserDBO) FromUser(u *domain.User) {
	tokensStr, _ := json.Marshal(u.Tokens)
	dbo.Id = u.Id
	dbo.Name = u.Name
	dbo.NormalizedName = u.NormalizedName
	dbo.Email = u.Email
	dbo.EmailConfirmed = u.EmailConfirmed
	dbo.EmailConfirmation = u.EmailConfirmation
	dbo.Img = ToNullString(u.Img)
	dbo.Tokens = ToNullString(string(tokensStr))
	dbo.Rights = int(u.Rights)
	dbo.Pass = u.Pass
	dbo.Salt = u.Salt
}

/*
	Users Plan
 ******************/
type UsersPlanDBO struct {
	UserId    string
	TopicName string
	PlanId    int
}

func (dbo *UsersPlanDBO) ToUsersPlan() *domain.UsersPlan {
	return &domain.UsersPlan{
		UserId:    dbo.UserId,
		TopicName: dbo.TopicName,
		PlanId:    dbo.PlanId,
	}
}

func (dbo *UsersPlanDBO) FromUsersPlan(up *domain.UsersPlan) {
	dbo.PlanId = up.PlanId
	dbo.UserId = up.UserId
	dbo.TopicName = up.TopicName
}

/*
	Topic Tag
 ******************/
type TopicTagDBO struct {
	Name  string
	Title string
}

func (dbo *TopicTagDBO) ToTopicTag() *domain.TopicTag {
	return &domain.TopicTag{
		Name:  dbo.Name,
		Title: dbo.Title,
	}
}

func (dbo *TopicTagDBO) FromTopicTag(tag *domain.TopicTag) {
	dbo.Name = tag.Name
	dbo.Title = tag.Title
}

/*
	Comments
*******************/
type CommentDBO struct {
	Id         int64
	EntityType int
	EntityId   int64
	ThreadId   sql.NullInt64 // id родительского комментария 0 уровня
	ParentId   sql.NullInt64
	Date       time.Time
	UserId     string
	Text       string
	Title      sql.NullString
	Deleted    bool
}

func (dbo *CommentDBO) ToComment() *domain.Comment {
	return &domain.Comment{
		Id:         dbo.Id,
		EntityType: domain.EntityType(dbo.EntityType),
		EntityId:   dbo.EntityId,
		ThreadId:   dbo.ThreadId.Int64,
		ParentId:   dbo.ParentId.Int64,
		Date:       dbo.Date,
		UserId:     dbo.UserId,
		Text:       dbo.Text,
		Title:      dbo.Title.String,
		Deleted:    dbo.Deleted,
		Childs:     []domain.Comment{},
	}
}

func (dbo *CommentDBO) FromComment(c *domain.Comment) {
	dbo.Id = c.Id
	dbo.EntityType = int(c.EntityType)
	dbo.EntityId = c.EntityId
	dbo.ThreadId = ToNullInt64(c.ThreadId)
	dbo.ParentId = ToNullInt64(c.ParentId)
	dbo.Date = c.Date
	dbo.UserId = c.UserId
	dbo.Text = c.Text
	dbo.Title = ToNullString(c.Title)
	dbo.Deleted = c.Deleted
}

/*
	ChangeLogRecord
*******************/
type ChangeLogRecordDBO struct {
	Id         int64
	Date       time.Time
	Action     int
	UserId     string
	EntityType int
	EntityId   int64
	Diff       sql.NullString
	Points     int
}

func (dbo *ChangeLogRecordDBO) ToChangeLogRecord() *domain.ChangeLogRecord {
	return &domain.ChangeLogRecord{
		Id:         dbo.Id,
		Date:       dbo.Date,
		Action:     domain.ChangeType(dbo.Action),
		UserId:     dbo.UserId,
		EntityType: domain.EntityType(dbo.EntityType),
		EntityId:   dbo.EntityId,
		Diff:       dbo.Diff.String,
		Points:     dbo.Points,
	}
}

func (dbo *ChangeLogRecordDBO) FromChangeLogRecord(c *domain.ChangeLogRecord) {
	dbo.Id = c.Id
	dbo.Date = c.Date
	dbo.Action = int(c.Action)
	dbo.UserId = c.UserId
	dbo.EntityType = int(c.EntityType)
	dbo.EntityId = c.EntityId
	dbo.Diff = ToNullString(c.Diff)
	dbo.Points = c.Points
}

/*
	Project
*******************/

type ProjectDBO struct {
	Id      int
	Title   string
	Text    string
	Tags    []string
	OwnerId string
}

func (dbo *ProjectDBO) ToProject(tags []domain.TopicTag) *domain.Project {
	return &domain.Project{
		Id:      dbo.Id,
		Title:   dbo.Title,
		Text:    dbo.Text,
		Tags:    tags,
		OwnerId: dbo.OwnerId,
	}
}

func (dbo *ProjectDBO) FromProject(p *domain.Project) {
	dbo.Id = p.Id
	dbo.Title = p.Title
	dbo.Text = p.Text
	dbo.Tags = make([]string, len(p.Tags))
	dbo.OwnerId = p.OwnerId

	for i, tag := range p.Tags {
		dbo.Tags[i] = tag.Name
	}
}

type PointsDBO struct {
	Id     int64
	Update time.Time
	Count  int
	Avg    float32
	Value  float32
	Voted  bool
}

func (dbo *PointsDBO) ToPoints() *domain.Points {
	return &domain.Points{
		Id:    dbo.Id,
		Count: dbo.Count,
		Avg:   dbo.Avg,
		Value: dbo.Value,
		Voted: dbo.Voted,
	}
}

func (dbo *PointsDBO) FromPoints(p *domain.Points) {
	dbo.Id = p.Id
	dbo.Count = p.Count
	dbo.Avg = p.Avg
	dbo.Value = p.Value
	dbo.Voted = p.Voted
}

func ToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func ToNullInt64(s int64) sql.NullInt64 {
	return sql.NullInt64{Int64: s, Valid: s != 0}
}
