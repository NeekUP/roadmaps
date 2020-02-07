package domain

import "time"

type ChangeLogRecord struct {
	Id         int64
	Date       time.Time
	Action     ChangeType
	UserId     string
	EntityType EntityType
	EntityId   int64
	Diff       string
	Points     int
}
