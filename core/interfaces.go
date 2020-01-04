package core

import (
	"github.com/NeekUP/roadmaps/domain"
	"time"
)

type UserRepository interface {
	Get(id string) *domain.User
	Save(user *domain.User) (bool, *AppError)
	Update(user *domain.User) (bool, *AppError)
	ExistsName(name string) (exists bool, ok bool)
	ExistsEmail(email string) (exists bool, ok bool)
	FindByEmail(email string) *domain.User
	Count() (count int, ok bool)

	//dev
	All() []domain.User
}

type SourceRepository interface {
	Get(id int64) *domain.Source
	FindByIdentifier(identifier string) *domain.Source
	Save(source *domain.Source) (bool, *AppError)
	Update(source *domain.Source) (bool, *AppError)
	GetOrAddByIdentifier(source *domain.Source) *domain.Source

	//dev
	All() []domain.Source
}

type TopicRepository interface {
	Get(name string) *domain.Topic
	GetById(id int) *domain.Topic
	Save(source *domain.Topic) (bool, *AppError)
	Update(source *domain.Topic) (bool, *AppError)
	TitleLike(str string, count int) []domain.Topic
	AddTag(tagname, topicname string) bool
	DeleteTag(tagname, topicname string) bool
	GetTags(topicnames []string) []domain.TopicTag

	//dev
	All() []domain.Topic
}

type PlanRepository interface {
	SaveWithSteps(plan *domain.Plan) (bool, *AppError)
	// should includes steps
	Get(id int) *domain.Plan
	GetList(id []int) []domain.Plan
	GetPopularByTopic(topic string, count int) []domain.Plan
	Update(plan *domain.Plan) (bool, *AppError)
	Delete(planId int) (bool, *AppError)
	//dev
	All() []domain.Plan
}

type StepRepository interface {
	GetByPlan(planid int) []domain.Step
	//dev
	All() []domain.Step
}

type UsersPlanRepository interface {
	// true - if already exists and when added new row
	// if same topic exists for user, delete exists and add new
	// false only if db error
	Add(userId string, topicName string, planId int) (bool, *AppError)
	Remove(userId string, planId int) (bool, *AppError)
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
