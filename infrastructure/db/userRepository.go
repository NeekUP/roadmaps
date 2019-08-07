// +build PROD

package db

import (
	"database/sql"
	"roadmaps/domain"
)

type userRepository struct {
	Conn *sql.DB
}

//func NewUserRepository(conn *sql.DB) core.UserRepository{
//	return &userRepository{Conn: conn}
//}

func (r *userRepository) Get(id string) *domain.User {
	return nil
}

func (r *userRepository) Create(user *domain.User, passHash []byte, salt []byte) bool {
	r.Conn.Exec("")
	return false
}

func (r *userRepository) Update(user *domain.User) bool {
	r.Conn.Exec("")
	return false
}

func (r *userRepository) CheckPass(id string, pass string) bool {
	r.Conn.Exec("")
	return false
}

func (r *userRepository) ExistsName(name string) bool {
	r.Conn.Exec("")
	return false
}

func (r *userRepository) ExistsEmail(email string) bool {
	r.Conn.Exec("")
	return false
}

func (r *userRepository) FindByEmail(email string) *domain.User {
	panic("implement me")
}
