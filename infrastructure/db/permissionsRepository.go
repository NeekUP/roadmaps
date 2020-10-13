package db

import (
	"database/sql"
	"fmt"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/jackc/pgx/v4"
)

type permissionsRepo struct {
	Db *DbConnection
}

func NewPermissionsRepository(db *DbConnection) core.PermissionsRepository {
	return &permissionsRepo{Db: db}
}

func (r *permissionsRepo) Get(ctx core.ReqContext, userId string, entityType domain.EntityType, entityId int64) uint64 {
	tr := ctx.StartTrace("PermissionsRepository.Get")
	defer ctx.StopTrace(tr)

	entityName := domain.EntityTypeToString(entityType)
	query := fmt.Sprintf(`SELECT permissions FROM permissions_%ss WHERE entityid=$1 AND userid=$2;`, entityName)
	row := r.Db.Conn.QueryRow(ctx, query, entityId, userId)
	permissions, err := r.scanRow(row)
	if err == sql.ErrNoRows {
		return 0
	}
	if err != nil {
		r.Db.LogError(err, query)
		return 0
	}
	return permissions
}

func (r *permissionsRepo) GetGlobal(ctx core.ReqContext, userId string) uint64 {
	tr := ctx.StartTrace("PermissionsRepository.Get")
	defer ctx.StopTrace(tr)

	query := `SELECT permissions FROM permissions_base WHERE entityid=$1 AND userid=$2;`
	row := r.Db.Conn.QueryRow(ctx, query, 0, userId)
	permissions, err := r.scanRow(row)
	if err == sql.ErrNoRows {
		return 0
	}
	if err != nil {
		r.Db.LogError(err, query)
		return 0
	}
	return permissions
}

func (r *permissionsRepo) Set(ctx core.ReqContext, userId string, entityType domain.EntityType, entityId int64, permissions uint64) *core.AppError {
	tr := ctx.StartTrace("PermissionsRepository.SetGlobal")
	defer ctx.StopTrace(tr)

	entityName := domain.EntityTypeToString(entityType)
	query := fmt.Sprintf(`INSERT INTO permissions_%ss SET permissions=$1 WHERE entityid=$2 AND userid=$3;`, entityName)
	_, err := r.Db.Conn.Exec(ctx, query, permissions, entityId, userId)
	if err != nil {
		return r.Db.LogError(err, query)
	}
	return nil
}

func (r *permissionsRepo) SetGlobal(ctx core.ReqContext, userId string, permissions uint64) *core.AppError {
	tr := ctx.StartTrace("PermissionsRepository.SetGlobal")
	defer ctx.StopTrace(tr)

	query := `INSERT INTO permissions_base SET permissions=$1 WHERE entityid=$1 AND userid=$2;`
	_, err := r.Db.Conn.Exec(ctx, query, 0, userId)
	if err != nil {
		return r.Db.LogError(err, query)
	}
	return nil
}

func (r permissionsRepo) scanRow(row pgx.Row) (uint64, error) {
	var perm uint64
	err := row.Scan(&perm)
	if err != nil && err.Error() == "no rows in result set" {
		return 0, sql.ErrNoRows
	}
	return perm, err
}
