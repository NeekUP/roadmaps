package db

import (
	"database/sql"
	"encoding/json"
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
}

func (dbo *TopicDBO) ToTopic() *domain.Topic {
	t := &domain.Topic{
		Id:          dbo.Id,
		Name:        dbo.Name,
		Title:       dbo.Title,
		Description: dbo.Description.String,
		Creator:     dbo.Creator,
	}
	if dbo.Tags == nil {
		t.Tags = []string{}
	} else {
		t.Tags = dbo.Tags
	}
	return t
}

func (dbo *TopicDBO) FromTopic(d *domain.Topic) {
	dbo.Id = d.Id
	dbo.Name = d.Name
	dbo.Title = d.Title
	dbo.Description = ToNullString(d.Description)
	dbo.Creator = d.Creator
	if d.Tags == nil {
		dbo.Tags = []string{}
	} else {
		dbo.Tags = d.Tags
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
	Points    int
}

func (dbo *PlanDBO) ToPlan() *domain.Plan {
	return &domain.Plan{
		Id:        dbo.Id,
		Title:     dbo.Title,
		TopicName: dbo.TopicName,
		OwnerId:   dbo.OwnerId,
		Points:    dbo.Points,
	}
}

func (dbo *PlanDBO) FromPlan(plan *domain.Plan) {
	dbo.Id = plan.Id
	dbo.Title = plan.Title
	dbo.TopicName = plan.TopicName
	dbo.OwnerId = plan.OwnerId
	dbo.Points = plan.Points
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
}

func (dbo *StepDBO) FromStep(step *domain.Step) {
	dbo.Id = step.Id
	dbo.PlanId = step.PlanId
	dbo.Position = step.Position
	dbo.ReferenceId = step.ReferenceId
	dbo.ReferenceType = string(step.ReferenceType)
}

func (dbo *StepDBO) ToStep() *domain.Step {
	return &domain.Step{
		Id:            dbo.Id,
		PlanId:        dbo.PlanId,
		ReferenceId:   dbo.ReferenceId,
		ReferenceType: domain.ReferenceType(dbo.ReferenceType),
		Position:      dbo.Position,
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

func ToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}
