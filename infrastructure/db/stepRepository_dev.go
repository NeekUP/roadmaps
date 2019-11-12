// +build DEV

package db

import (
	"database/sql"
	"roadmaps/core"
	"roadmaps/domain"
	"sync"
)

var (
	Steps        = make([]domain.Step, 0)
	StepsMux     sync.Mutex
	StepsCounter int
)

type stepRepoInMemory struct {
	Conn *sql.DB
}

func NewStepsRepository(conn *sql.DB) core.StepRepository {
	return &stepRepoInMemory{
		Conn: conn}
}

func (this stepRepoInMemory) All() []domain.Step {
	return Steps
}
