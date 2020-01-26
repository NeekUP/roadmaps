package domain

import "time"

type Comment struct {
	Id         int64
	EntityType EntityType
	EntityId   int64
	ThreadId   int64 // id родительского комментария 0 уровня
	ParentId   int64
	Date       time.Time
	UserId     string
	User       *User
	Text       string
	Title      string
	Deleted    bool
	Points     int
	Childs     []Comment
}
