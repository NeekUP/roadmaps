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
}

type SourceRepository interface {
	Get(id int) *domain.Source
	FindByIdentifier(identifier string) *domain.Source
	Save(source *domain.Source) bool
	Update(source *domain.Source) bool
	GetOrAddByIdentifier(source *domain.Source) *domain.Source
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
	Validate(authToken string) (userID string, rights int, err error)
}

type ReqContext interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
	ReqId() string
}

type ImageManager interface {
	Save(data []byte, name string) error
	GetUrl(name string) string
}
