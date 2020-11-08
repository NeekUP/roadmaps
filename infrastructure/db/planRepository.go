package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/jackc/pgx/v4"
)

type planRepo struct {
	Db *DbConnection
}

func NewPlansRepository(db *DbConnection) core.PlanRepository {
	return &planRepo{
		Db: db}
}

func (r *planRepo) SaveWithSteps(ctx core.ReqContext, plan *domain.Plan) (bool, *core.AppError) {
	tr := ctx.StartTrace("PlanRepository.SaveWithSteps")
	defer ctx.StopTrace(tr)
	if len(plan.Steps) == 0 {
		return false, core.NewError(core.InvalidRequest)
	}
	tx, err := r.Db.Conn.BeginTx(context.Background(), pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})
	insertPlanQuery := "INSERT INTO plans(title, topic, owner, isdraft) VALUES ($1, $2, $3, $4) RETURNING id;"

	err = r.Db.Conn.QueryRow(context.Background(), insertPlanQuery, plan.Title, plan.TopicName, plan.OwnerId, plan.IsDraft).Scan(&plan.Id)
	tr.Point("insert plan")
	if err != nil {
		if e := tx.Rollback(context.Background()); e != nil {
			r.Db.Log.Errorw("Tx not rolled back", "err", e.Error())
		}
		return false, r.Db.LogError(err, insertPlanQuery)
	}

	for i := 0; i < len(plan.Steps); i++ {
		plan.Steps[i].PlanId = plan.Id
		query := "INSERT INTO steps( planid, referenceid, referencetype, position, title) VALUES ($1, $2, $3, $4, $5) RETURNING id;"
		err := r.Db.Conn.QueryRow(context.Background(),
			query,
			plan.Steps[i].PlanId,
			plan.Steps[i].ReferenceId,
			plan.Steps[i].ReferenceType,
			plan.Steps[i].Position,
			plan.Steps[i].Title).
			Scan(&plan.Steps[i].Id)

		tr.Point("insert steps")
		if err != nil {
			if e := tx.Rollback(context.Background()); e != nil {
				r.Db.Log.Errorw("Tx not rolled back", "err", e.Error())
			}
			return false, r.Db.LogError(err, query)
		}
	}

	err = tx.Commit(context.Background())
	// if err is serialization error, we should repeat transaction
	if err != nil {
		return false, r.Db.LogError(err, "")
	}
	return true, nil
}

func (r *planRepo) Update(ctx core.ReqContext, plan *domain.Plan) (bool, *core.AppError) {
	tr := ctx.StartTrace("PlanRepository.Update")
	defer ctx.StopTrace(tr)

	if len(plan.Steps) == 0 {
		return false, core.NewError(core.InvalidRequest)
	}

	updatePlanQuery := `UPDATE plans SET title = $1, topic = $2, isdraft = $4 WHERE id = $3`
	_, err := r.Db.Conn.Exec(context.Background(), updatePlanQuery, plan.Title, plan.TopicName, plan.Id, plan.IsDraft)
	tr.Point("update plans")
	if err != nil {
		return false, r.Db.LogError(err, updatePlanQuery)
	}

	tx, err := r.Db.Conn.BeginTx(context.Background(), pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})
	deleteStepsQuery := `DELETE FROM steps WHERE planid = $1`
	_, err = r.Db.Conn.Exec(context.Background(), deleteStepsQuery, plan.Id)
	tr.Point("delete steps")
	if err != nil {
		if e := tx.Rollback(context.Background()); e != nil {
			r.Db.Log.Errorw("Tx not rolled back", "err", e.Error())
		}
		return false, r.Db.LogError(err, updatePlanQuery)
	}

	for i := 0; i < len(plan.Steps); i++ {
		plan.Steps[i].PlanId = plan.Id
		query := "INSERT INTO steps( planid, referenceid, referencetype, position, title) VALUES ($1, $2, $3, $4, $5) RETURNING id;"
		err := r.Db.Conn.QueryRow(context.Background(),
			query,
			plan.Steps[i].PlanId,
			plan.Steps[i].ReferenceId,
			plan.Steps[i].ReferenceType,
			plan.Steps[i].Position,
			plan.Steps[i].Title).
			Scan(&plan.Steps[i].Id)

		tr.Point("insert steps")
		if err != nil {
			if e := tx.Rollback(context.Background()); e != nil {
				r.Db.Log.Errorw("Tx not rolled back", "err", e.Error())
			}
			return false, r.Db.LogError(err, query)
		}
	}

	err = tx.Commit(context.Background())
	// if err is serialization error, we should repeat transaction
	if err != nil {
		return false, r.Db.LogError(err, "")
	}

	return true, nil
}

func (r *planRepo) Delete(ctx core.ReqContext, planId int) (bool, *core.AppError) {
	tr := ctx.StartTrace("PlanRepository.Delete")
	defer ctx.StopTrace(tr)

	deletePlanQuery := `DELETE FROM plans WHERE id = $1`
	_, err := r.Db.Conn.Exec(context.Background(), deletePlanQuery, planId)
	if err != nil {
		return false, r.Db.LogError(err, deletePlanQuery)
	}
	deleteStepsQuery := `DELETE FROM steps WHERE planid = $1`
	_, err = r.Db.Conn.Exec(context.Background(), deleteStepsQuery, planId)
	if err != nil {
		return false, r.Db.LogError(err, deletePlanQuery)
	}
	return true, nil
}

func (r *planRepo) Get(ctx core.ReqContext, id int) *domain.Plan {
	tr := ctx.StartTrace("PlanRepository.Get")
	defer ctx.StopTrace(tr)

	query := `SELECT id, title, topic, owner, isdraft FROM plans WHERE id=$1 AND isdraft=false ;`
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

func (r *planRepo) GetWithDraft(ctx core.ReqContext, id int, userid string) *domain.Plan {
	tr := ctx.StartTrace("PlanRepository.Get")
	defer ctx.StopTrace(tr)

	query := `SELECT id, title, topic, owner, isdraft FROM plans WHERE id=$1 AND ( isdraft=false OR owner=$2 );`
	row := r.Db.Conn.QueryRow(context.Background(), query, id, userid)
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

func (r *planRepo) GetList(ctx core.ReqContext, id []int) []domain.Plan {
	tr := ctx.StartTrace("PlanRepository.GetList")
	defer ctx.StopTrace(tr)

	query := "SELECT id, title, topic, owner, isdraft FROM plans WHERE id IN (%s) AND isdraft=false;"
	query = fmt.Sprintf(query, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(id)), ","), "[]"))
	rows, err := r.Db.Conn.Query(context.Background(), query)
	if err != nil {
		r.Db.LogError(err, query)
		return []domain.Plan{}
	}
	defer rows.Close()
	return r.scanRows(rows)
}

func (r *planRepo) GetByUser(ctx core.ReqContext, userid string, count int, page int) []domain.Plan {
	tr := ctx.StartTrace("PlanRepository.GetByUser")
	defer ctx.StopTrace(tr)

	query := "SELECT id, title, topic, owner, isdraft " +
		"FROM plans " +
		"WHERE owner =$1 ORDER BY id DESC LIMIT $2 OFFSET $3;"
	rows, err := r.Db.Conn.Query(context.Background(), query, userid, count, page*count)
	if err != nil {
		r.Db.LogError(err, query)
		return []domain.Plan{}
	}
	defer rows.Close()
	return r.scanRows(rows)
}

func (r *planRepo) GetPopularByTopic(ctx core.ReqContext, topic string, count int) []domain.Plan {
	tr := ctx.StartTrace("PlanRepository.GetPopularByTopic")
	defer ctx.StopTrace(tr)

	query := "SELECT p.id, p.title, p.topic, p.owner, p.isdraft FROM plans p LEFT JOIN points_aggregated_plans ps ON p.id=ps.entityid WHERE p.topic=$1 AND p.isdraft=false ORDER BY ps.value DESC LIMIT $2"
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
	query := "SELECT id, title, topic, owner, isdraft FROM plans"
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
	err := row.Scan(&dbo.Id, &dbo.Title, &dbo.TopicName, &dbo.OwnerId, &dbo.IsDraft)
	if err != nil && err.Error() == "no rows in result set" {
		return &dbo, sql.ErrNoRows
	}
	return &dbo, err
}
