package core

import (
	"github.com/NeekUP/nptrace"
	"github.com/NeekUP/roadmaps/domain"
	"time"
)

type UserRepository interface {
	Get(ctx ReqContext, id string) *domain.User
	GetList(ctx ReqContext, id []string) []domain.User
	Save(ctx ReqContext, user *domain.User) (bool, *AppError)
	Update(ctx ReqContext, user *domain.User) (bool, *AppError)
	ExistsName(ctx ReqContext, name string) (exists bool, ok bool)
	ExistsEmail(ctx ReqContext, email string) (exists bool, ok bool)
	FindByEmail(ctx ReqContext, email string) *domain.User
	Count(ctx ReqContext) (count int, ok bool)

	//dev
	All() []domain.User
}

type SourceRepository interface {
	Get(ctx ReqContext, id int64) *domain.Source
	FindByIdentifier(ctx ReqContext, identifier string) *domain.Source
	Save(ctx ReqContext, source *domain.Source) (bool, *AppError)
	Update(ctx ReqContext, source *domain.Source) (bool, *AppError)
	GetOrAddByIdentifier(ctx ReqContext, source *domain.Source) *domain.Source

	//dev
	All() []domain.Source
}

type TopicRepository interface {
	Get(ctx ReqContext, name string) *domain.Topic
	GetById(ctx ReqContext, id int) *domain.Topic
	Save(ctx ReqContext, source *domain.Topic) (bool, *AppError)
	Update(ctx ReqContext, source *domain.Topic) (bool, *AppError)
	TitleLike(ctx ReqContext, str string, count int) []domain.Topic
	AddTag(ctx ReqContext, tagname, topicname string) bool
	DeleteTag(ctx ReqContext, tagname, topicname string) bool
	GetTags(ctx ReqContext, topicnames []string) []domain.TopicTag

	//dev
	All() []domain.Topic
}

type PlanRepository interface {
	SaveWithSteps(ctx ReqContext, plan *domain.Plan) (bool, *AppError)
	// should includes steps
	Get(ctx ReqContext, id int) *domain.Plan
	GetList(ctx ReqContext, id []int) []domain.Plan
	GetPopularByTopic(ctx ReqContext, topic string, count int) []domain.Plan
	Update(ctx ReqContext, plan *domain.Plan) (bool, *AppError)
	Delete(ctx ReqContext, planId int) (bool, *AppError)
	//dev
	All() []domain.Plan
}

type StepRepository interface {
	GetByPlan(ctx ReqContext, planid int) []domain.Step
	//dev
	All() []domain.Step
}

type UsersPlanRepository interface {
	// true - if already exists and when added new row
	// if same topic exists for user, delete exists and add new
	// false only if db error
	Add(ctx ReqContext, userId string, topicName string, planId int) (bool, *AppError)
	Remove(ctx ReqContext, userId string, planId int) (bool, *AppError)
	GetByTopic(ctx ReqContext, userId, topicName string) *domain.UsersPlan
	GetByUser(ctx ReqContext, userId string) []domain.UsersPlan
}

type CommentsRepository interface {
	Add(ctx ReqContext, comment *domain.Comment) (bool, error)
	Update(ctx ReqContext, id int64, text, title string) (bool, error)
	Delete(ctx ReqContext, id int64) (bool, error)
	Get(ctx ReqContext, id int64) *domain.Comment
	GetThreadList(ctx ReqContext, entityType int, entityId int64, count int, page int) []domain.Comment
	GetThread(ctx ReqContext, entityType int, entityId int64, threadId int64) []domain.Comment
}

type ChangeLogRepository interface {
	Add(record *domain.ChangeLogRecord) bool
}

type ProjectsRepository interface {
	Add(ctx ReqContext, project *domain.Project) (bool, error)
	Update(ctx ReqContext, project *domain.Project) (bool, error)
	Get(ctx ReqContext, id int) *domain.Project
}

type PointsRepository interface {
	Add(ctx ReqContext, entityType domain.EntityType, entityId int64, userId string, value int) bool
	Get(ctx ReqContext, userid string, entityType domain.EntityType, entityId int64) *domain.Points
	GetList(ctx ReqContext, userid string, entityType domain.EntityType, entityId []int64) []domain.Points
}

type ChangeLog interface {
	Added(entityType domain.EntityType, entityId int64, userId string)
	Edited(entityType domain.EntityType, entityId int64, userId string, before interface{}, after interface{})
	Deleted(entityType domain.EntityType, entityId int64, userId string)
}

type HashProvider interface {
	HashPassword(pass string) (hash []byte, salt []byte)
	CheckPassword(pass string, hash []byte, salt []byte) bool
}

type EmailSender interface {
	Send(recipient string, subject string, body string) (bool, error)
	Registration(recipient, userId, secret string) (bool, error)
}

type TokenService interface {
	Create(ctx ReqContext, user *domain.User, fingerprint, useragent string) (auth string, refresh string, err error)
	Refresh(ctx ReqContext, authToken, refreshToken, fingerprint, useragent string) (aToken string, rToken string, err error)
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
	StartTrace(name string, args ...interface{}) *nptrace.Trace
	StopTrace(t *nptrace.Trace)
}

type ImageManager interface {
	Save(data []byte, name string) error
	GetUrl(name string) string
}
