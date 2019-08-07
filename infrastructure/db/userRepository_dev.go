// +build DEV

package db

import (
	"database/sql"
	"roadmaps/core"
	"roadmaps/domain"
	"strings"
)

type userRepoInMemory struct {
	Conn  *sql.DB
	Users []domain.User
}

func NewUserRepository(conn *sql.DB) core.UserRepository {
	return &userRepoInMemory{
		Conn:  conn,
		Users: SeedUsers()}
}

func (this *userRepoInMemory) Get(id string) *domain.User {
	for i := 0; i < len(this.Users); i++ {
		if this.Users[i].Id == id {
			return &this.Users[i]
		}
	}
	return nil
}

func (this *userRepoInMemory) Create(user *domain.User, passHash []byte, salt []byte) bool {
	user.Pass = passHash
	user.Salt = salt
	this.Users = append(this.Users, *user)
	return true
}

func (this *userRepoInMemory) Update(user *domain.User) bool {
	for i := 0; i < len(this.Users); i++ {
		if this.Users[i].Id == user.Id {
			this.Users[i] = *user
			return true
		}
	}

	return false
}

func (this *userRepoInMemory) ExistsName(name string) bool {

	name = strings.ToLower(name)
	for i := 0; i < len(this.Users); i++ {
		if this.Users[i].Name == name {
			return true
		}
	}

	return false
}

func (this *userRepoInMemory) ExistsEmail(email string) bool {
	return this.FindByEmail(email) != nil
}

func (this *userRepoInMemory) FindByEmail(email string) *domain.User {
	email = strings.ToLower(email)
	for i := 0; i < len(this.Users); i++ {
		if this.Users[i].Email == email {
			return &this.Users[i]
		}
	}

	return nil
}
