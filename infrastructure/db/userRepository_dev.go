// +build DEV

package db

import (
	"database/sql"
	"roadmaps/core"
	"roadmaps/domain"
	"sync"
)

var (
	Users = make([]domain.User, 0)
	mux   sync.Mutex
)

type userRepoInMemory struct {
	Conn *sql.DB
}

func NewUserRepository(conn *sql.DB) core.UserRepository {
	return &userRepoInMemory{
		Conn: conn}
}

func (this *userRepoInMemory) Get(id string) *domain.User {

	mux.Lock()
	defer mux.Unlock()
	for i := 0; i < len(Users); i++ {
		if Users[i].Id == id {
			return &Users[i]
		}
	}
	return nil
}

func (this *userRepoInMemory) Create(user *domain.User, passHash []byte, salt []byte) bool {
	user.Pass = passHash
	user.Salt = salt

	mux.Lock()
	defer mux.Unlock()

	for i := 0; i < len(Users); i++ {
		if Users[i].Email == user.Email || Users[i].Name == user.Name || Users[i].Id == user.Id {
			return false
		}
	}

	Users = append(Users, *user)
	return true
}

func (this *userRepoInMemory) Update(user *domain.User) bool {

	mux.Lock()
	defer mux.Unlock()

	for i := 0; i < len(Users); i++ {
		if Users[i].Id == user.Id {
			Users[i] = *user
			return true
		}
	}

	return false
}

func (this *userRepoInMemory) ExistsName(name string) bool {

	mux.Lock()
	defer mux.Unlock()

	for i := 0; i < len(Users); i++ {
		if Users[i].Name == name {
			return true
		}
	}

	return false
}

func (this *userRepoInMemory) ExistsEmail(email string) bool {
	// for general purposes
	return this.FindByEmail(email) != nil
}

func (this *userRepoInMemory) FindByEmail(email string) *domain.User {
	mux.Lock()
	defer mux.Unlock()

	for i := 0; i < len(Users); i++ {
		if Users[i].Email == email {
			return &Users[i]
		}
	}
	return nil
}
