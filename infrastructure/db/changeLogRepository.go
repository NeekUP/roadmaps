package db

import (
	"context"

	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type changeLogRepo struct {
	Db *DbConnection
}

func NewChangeLogRepository(db *DbConnection) core.ChangeLogRepository {
	return &changeLogRepo{Db: db}
}

// TODO: add indexes int db!
func (r *changeLogRepo) Add(record *domain.ChangeLogRecord) bool {
	dbo := &ChangeLogRecordDTO{}
	dbo.FromChangeLogRecord(record)
	query := `INSERT INTO changelog(date, action, userid, entitytype, entityid, diff, points)
		VALUES (now(), $1, $2, $3, $4, $5, 0);`
	tag, err := r.Db.Conn.Exec(context.Background(), query, dbo.Action, dbo.UserId, dbo.EntityType, dbo.EntityId, dbo.Diff)
	if err != nil {
		r.Db.LogError(err, query)
		return false
	}
	return tag.RowsAffected() > 0
}
