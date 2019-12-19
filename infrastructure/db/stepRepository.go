package db

import (
	"context"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/jackc/pgx/v4"
)

type stepRepo struct {
	Db *DbConnection
}

func NewStepsRepository(db *DbConnection) core.StepRepository {
	return &stepRepo{
		Db: db}
}

func (r stepRepo) All() []domain.Step {
	query := "SELECT id, planid, referenceid, referencetype, position FROM steps;"
	rows, err := r.Db.Conn.Query(context.Background(), query)
	if err != nil {
		return []domain.Step{}
	}
	defer rows.Close()
	steps := make([]domain.Step, 0)
	for rows.Next() {
		dbo, err := r.scanRow(rows)
		if err != nil {
			return []domain.Step{}
		}
		steps = append(steps, *dbo.ToStep())
	}

	return steps
}

func (r stepRepo) GetByPlan(planid int) []domain.Step {
	query := "SELECT id, planid, referenceid, referencetype, position FROM steps WHERE planid=$1;"
	rows, err := r.Db.Conn.Query(context.Background(), query, planid)
	if err != nil {
		r.Db.LogError(err, query)
		return []domain.Step{}
	}
	defer rows.Close()
	steps := make([]domain.Step, 0)
	for rows.Next() {
		dbo, err := r.scanRow(rows)
		if err != nil {
			return []domain.Step{}
		}
		steps = append(steps, *dbo.ToStep())
	}

	return steps
}

func (r *stepRepo) scanRow(row pgx.Row) (*StepDBO, error) {
	st := StepDBO{}
	err := row.Scan(&st.Id, &st.PlanId, &st.ReferenceId, &st.ReferenceType, &st.Position)
	return &st, err
}
