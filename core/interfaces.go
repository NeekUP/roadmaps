package core

import (
	"roadmaps/domain"
	"time"
)

type UserRepository interface {
	Get(id string) *domain.User
	// Should be transaction with check name and email
	Save(user *domain.User, passHash []byte, salt []byte) bool
	Update(user *domain.User) bool
	ExistsName(name string) bool
	ExistsEmail(email string) bool
	FindByEmail(email string) *domain.User
	Count() int

	//dev
	All() []domain.User
}

type SourceRepository interface {
	Get(id int) *domain.Source
	FindByIdentifier(identifier string) *domain.Source
	Save(source *domain.Source) bool
	Update(source *domain.Source) bool
	GetOrAddByIdentifier(source *domain.Source) *domain.Source

	//dev
	All() []domain.Source
}

type TopicRepository interface {
	Get(name string) *domain.Topic
	GetById(id int) *domain.Topic
	Save(source *domain.Topic) bool
	Update(source *domain.Topic) bool

	//dev
	All() []domain.Topic
}

type PlanRepository interface {
	SaveWithSteps(plan *domain.Plan) error
	// should includes steps
	Get(id int) *domain.Plan
	GetList(id []int) []domain.Plan
	GetTopByTopicName(topic string, count int) []domain.Plan

	//dev
	All() []domain.Plan
}

type StepRepository interface {
	//dev
	All() []domain.Step
}

type UsersPlanRepository interface {
	// true - if already exists and when added new row
	// if same topic exists for user, delete exists and add new
	// false only if db error
	Add(userId string, topicName string, planId int) bool
	Remove(userId string, planId int) bool
	GetByTopic(userId, topicName string) *domain.UsersPlan
	GetByUser(userId string) []domain.UsersPlan
}

type HashProvider interface {
	HashPassword(pass string) (hash []byte, salt []byte)
	CheckPassword(pass string, hash []byte, salt []byte) bool
}

type EmailChecker interface {
	IsValid(email string) bool
	IsExists(email string) (exists bool, errCode string, errMeg string)
}

type TokenService interface {
	Create(user *domain.User, fingerprint, useragent string) (auth string, refresh string, err error)
	Refresh(authToken, refreshToken, fingerprint, useragent string) (aToken string, rToken string, err error)
	Validate(authToken string) (userID string, userName string, rights int, err error)
}

type ReqContext interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
	ReqId() string
	UserId() string
	UserName() string
}

type ImageManager interface {
	Save(data []byte, name string) error
	GetUrl(name string) string
}
