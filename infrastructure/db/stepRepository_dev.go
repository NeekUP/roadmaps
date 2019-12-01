//
//
package db

//
//import (
//	"sync"
//
//	"github.com/NeekUP/roadmaps/core"
//	"github.com/NeekUP/roadmaps/domain"
//	"github.com/jackc/pgx/v4"
//)
//
//var (
//	Steps        = make([]domain.Step, 0)
//	StepsMux     sync.Mutex
//	StepsCounter int
//)
//
//type stepRepoInMemory struct {
//	Conn *pgx.Conn
//}
//
//func NewStepsRepository(conn *pgx.Conn) core.StepRepository {
//	return &stepRepoInMemory{
//		Conn: conn}
//}
//
//func (this stepRepoInMemory) All() []domain.Step {
//	return Steps
//}
