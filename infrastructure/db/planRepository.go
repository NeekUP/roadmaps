package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/jackc/pgx/v4"
	"strings"
)

type planRepo struct {
	Db *DbConnection
}

func NewPlansRepository(db *DbConnection) core.PlanRepository {
	return &planRepo{
		Db: db}
}

func (r *planRepo) SaveWithSteps(plan *domain.Plan) (bool, *core.AppError) {

	if len(plan.Steps) == 0 {
		return false, core.NewError(core.InvalidRequest)
	}
	tx, err := r.Db.Conn.BeginTx(context.Background(), pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})
	insertPlanQuery := "INSERT INTO plans(title, topic, owner, points) VALUES ($1, $2, $3, $4) RETURNING id;"
	err = r.Db.Conn.QueryRow(context.Background(), insertPlanQuery, plan.Title, plan.TopicName, plan.OwnerId, plan.Points).Scan(&plan.Id)
	if err != nil {
		//r.Db.Log.Errorw(err.)
		if e := tx.Rollback(context.Background()); e != nil {
			r.Db.Log.Errorw("Tx not rolled back", "err", e.Error())
		}
		return false, r.Db.LogError(err, insertPlanQuery)
	}

	for i := 0; i < len(plan.Steps); i++ {
		plan.Steps[i].PlanId = plan.Id
		query := "INSERT INTO steps( planid, referenceid, referencetype, position) VALUES ($1, $2, $3, $4) RETURNING id;"
		err := r.Db.Conn.QueryRow(context.Background(), query, plan.Steps[i].PlanId, plan.Steps[i].ReferenceId, plan.Steps[i].ReferenceType, plan.Steps[i].Position).Scan(&plan.Steps[i].Id)
		if err != nil {
			if e := tx.Rollback(context.Background()); e != nil {
				r.Db.Log.Errorw("Tx not rolled back", "err", e.Error())
			}
			return false, r.Db.LogError(err, query)
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return false, r.Db.LogError(err, "")
	}
	return true, nil
}

func (r *planRepo) Get(id int) *domain.Plan {
	query := `SELECT id, title, topic, owner, points FROM plans WHERE id=$1;`
	row := r.Db.Conn.QueryRow(context.Background(), query, id)
	p, err := r.scanRow(row)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		r.Db.LogError(err, query)
		return nil
	}
	return p.ToPlan()
}

func (r *planRepo) GetList(id []int) []domain.Plan {
	query := "SELECT id, title, topic, owner, points FROM plans WHERE id IN (%s);"
	query = fmt.Sprintf(query, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(id)), ","), "[]"))
	rows, err := r.Db.Conn.Query(context.Background(), query)
	if err != nil {
		r.Db.LogError(err, query)
		return []domain.Plan{}
	}
	defer rows.Close()
	return r.scanRows(rows)
}

func (r *planRepo) GetPopularByTopic(topic string, count int) []domain.Plan {
	query := "SELECT id, title, topic, owner, points FROM plans WHERE topic=$1 ORDER BY points DESC LIMIT $2"
	rows, err := r.Db.Conn.Query(context.Background(), query, topic, count)
	if err != nil {
		r.Db.LogError(err, query)
		return []domain.Plan{}
	}
	defer rows.Close()
	return r.scanRows(rows)
}

func (r *planRepo) scanRows(rows pgx.Rows) []domain.Plan {
	plans := make([]domain.Plan, 0)
	for rows.Next() {
		dbo, err := r.scanRow(rows)
		if err != nil {
			return []domain.Plan{}
		}
		plans = append(plans, *dbo.ToPlan())
	}
	return plans
}

func (r *planRepo) All() []domain.Plan {
	query := "SELECT id, title, topic, owner, points FROM plans"
	rows, err := r.Db.Conn.Query(context.Background(), query)
	if err != nil {
		r.Db.LogError(err, query)
		return []domain.Plan{}
	}
	defer rows.Close()
	return r.scanRows(rows)
}

func (r *planRepo) scanRow(row pgx.Row) (*PlanDBO, error) {
	dbo := PlanDBO{}
	err := row.Scan(&dbo.Id, &dbo.Title, &dbo.TopicName, &dbo.OwnerId, &dbo.Points)
	return &dbo, err
}
