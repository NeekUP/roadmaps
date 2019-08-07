package core

import (
	"roadmaps/domain"
)

type UserRepository interface {
	Get(id string) *domain.User
	// Should be transaction with check name and email
	Create(user *domain.User, passHash []byte, salt []byte) bool
	Update(user *domain.User) bool
	ExistsName(name string) bool
	ExistsEmail(email string) bool
	FindByEmail(email string) *domain.User
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
	Validate(authToken string) (bool, error)
}
