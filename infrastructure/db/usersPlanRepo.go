package db

import (
	"context"
	"database/sql"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/jackc/pgx/v4"
)

type usersPlanRepo struct {
	Db *DbConnection
}

func NewUsersPlanRepository(db *DbConnection) core.UsersPlanRepository {
	return &usersPlanRepo{
		Db: db}
}

func (repo *usersPlanRepo) Add(ctx core.ReqContext, userId string, topicName string, planId int) (bool, *core.AppError) {
	tr := ctx.StartTrace("UsersPlanRepository.Add")
	defer ctx.StopTrace(tr)

	dbo := &UsersPlanDBO{}
	dbo.FromUsersPlan(&domain.UsersPlan{UserId: userId, TopicName: topicName, PlanId: planId})
	query := `INSERT INTO usersplans (userid, topic, planid) VALUES ($1, $2, $3);`
	tag, err := repo.Db.Conn.Exec(context.Background(), query, dbo.UserId, dbo.TopicName, dbo.PlanId)
	if err != nil {
		return false, repo.Db.LogError(err, query)
	}
	return tag.RowsAffected() > 0, nil
}

func (repo *usersPlanRepo) Remove(ctx core.ReqContext, userId string, planId int) (bool, *core.AppError) {
	tr := ctx.StartTrace("UsersPlanRepository.Remove")
	defer ctx.StopTrace(tr)
	query := `DELETE FROM usersplans WHERE userid=$1 AND planid=$2;`
	t, err := repo.Db.Conn.Exec(context.Background(), query, userId, planId)
	if err != nil {
		return false, repo.Db.LogError(err, query)
	}
	return t.RowsAffected() > 0, nil
}

func (repo *usersPlanRepo) GetByTopic(ctx core.ReqContext, userId, topicName string) *domain.UsersPlan {
	tr := ctx.StartTrace("UsersPlanRepository.GetByTopic")
	defer ctx.StopTrace(tr)
	query := `SELECT userid, topic, planid FROM usersplans WHERE userid=$1 AND topic=$2`
	row := repo.Db.Conn.QueryRow(context.Background(), query, userId, topicName)
	dbo, err := repo.scanRow(row)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		repo.Db.LogError(err, query)
		return nil
	}
	return dbo.ToUsersPlan()
}

func (repo *usersPlanRepo) GetByUser(ctx core.ReqContext, userId string) []domain.UsersPlan {
	tr := ctx.StartTrace("UsersPlanRepository.GetByUser")
	defer ctx.StopTrace(tr)

	query := `SELECT userid, topic, planid FROM usersplans WHERE userid=$1`
	rows, err := repo.Db.Conn.Query(context.Background(), query, userId)
	if err == sql.ErrNoRows {
		return []domain.UsersPlan{}
	}
	if err != nil {
		repo.Db.LogError(err, query)
		return []domain.UsersPlan{}
	}
	defer rows.Close()
	usersPlans := make([]domain.UsersPlan, 0)
	for rows.Next() {
		dbo, err := repo.scanRow(rows)
		if err != nil {
			return []domain.UsersPlan{}
		}
		usersPlans = append(usersPlans, *dbo.ToUsersPlan())
	}
	return usersPlans
}

func (repo *usersPlanRepo) scanRow(row pgx.Row) (*UsersPlanDBO, error) {
	dbo := UsersPlanDBO{}
	err := row.Scan(&dbo.UserId, &dbo.TopicName, &dbo.PlanId)
	if err != nil && err.Error() == "no rows in result set" {
		return &dbo, sql.ErrNoRows
	}
	return &dbo, err
}
