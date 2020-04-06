package db

import (
	"database/sql"
	"fmt"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/jackc/pgx/v4"
	"strings"
)

type pointsRepo struct {
	Db *DbConnection
}

func NewPointsRepository(db *DbConnection) core.PointsRepository {
	return &pointsRepo{Db: db}
}

func (r *pointsRepo) Add(ctx core.ReqContext, entityType domain.EntityType, entityId int64, userId string, value int) bool {
	tr := ctx.StartTrace("PointsRepository.Add")
	defer ctx.StopTrace(tr)
	query := fmt.Sprintf("INSERT INTO points_%ss (entityid,userid,date,value) VALUES ($1,$2,now(),$3 )", domain.EntityTypeToString(entityType))
	_, err := r.Db.Conn.Exec(ctx, query, entityId, userId, value)
	if err != nil {
		r.Db.LogError(err, query)
		return false
	}
	return true
}

func (r *pointsRepo) Get(ctx core.ReqContext, userid string, entityType domain.EntityType, entityId int64) *domain.Points {
	tr := ctx.StartTrace("PointsRepository.Get")
	defer ctx.StopTrace(tr)
	var query string
	var row pgx.Row
	entityName := domain.EntityTypeToString(entityType)
	if userid == "" {
		query = fmt.Sprintf("SELECT entityid, updatedate, count, value, avg, false FROM points_aggregated_%ss WHERE entityid=$1", entityName)
		row = r.Db.Conn.QueryRow(ctx, query, entityId)
	} else {
		query = fmt.Sprintf("SELECT entityid, updatedate, count, value, avg, EXISTS (SELECT entityid FROM points_%ss WHERE entityid=ps.entityid AND userid=$2) FROM points_aggregated_%ss ps WHERE entityid=$1", entityName, entityName)
		row = r.Db.Conn.QueryRow(ctx, query, entityId, userid)
	}

	dbo, err := r.scanRow(row)
	if err == sql.ErrNoRows {
		return nil
	}

	if err != nil {
		r.Db.LogError(err, query)
		return nil
	}
	return dbo.ToPoints()
}

func (r *pointsRepo) GetList(ctx core.ReqContext, userid string, entityType domain.EntityType, entityId []int64) []domain.Points {
	tr := ctx.StartTrace("PointsRepository.GetList")
	defer ctx.StopTrace(tr)
	idList := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(entityId)), ","), "[]")
	var query string
	var rows pgx.Rows
	var err error
	entityName := domain.EntityTypeToString(entityType)
	if userid == "" {
		query = fmt.Sprintf("SELECT entityid, updatedate, count, value, avg, false FROM points_aggregated_%ss WHERE entityid IN (%s)", entityName, idList)
		rows, err = r.Db.Conn.Query(ctx, query)
	} else {
		query = fmt.Sprintf("SELECT entityid, updatedate, count, value, avg, EXISTS (SELECT entityid FROM points_%ss WHERE entityid=ps.entityid AND userid=$1) FROM points_aggregated_%ss ps WHERE entityid IN (%s)", entityName, entityName, idList)
		rows, err = r.Db.Conn.Query(ctx, query, userid)
	}

	if err != nil {
		r.Db.LogError(err, query)
		return []domain.Points{}
	}
	defer rows.Close()
	return r.scanRows(rows)
}

func (r *pointsRepo) scanRows(rows pgx.Rows) []domain.Points {
	plans := make([]domain.Points, 0)
	for rows.Next() {
		dbo, err := r.scanRow(rows)
		if err != nil {
			return []domain.Points{}
		}
		plans = append(plans, *dbo.ToPoints())
	}
	return plans
}

func (r *pointsRepo) scanRow(row pgx.Row) (*PointsDBO, error) {
	dbo := PointsDBO{}
	err := row.Scan(&dbo.Id, &dbo.Update, &dbo.Count, &dbo.Value, &dbo.Avg, &dbo.Voted)
	if err != nil && err.Error() == "no rows in result set" {
		return &dbo, sql.ErrNoRows
	}
	return &dbo, err
}
