package db

import (
	"context"
	"database/sql"

	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/jackc/pgx/v4"
)

type commentsRepo struct {
	Db *DbConnection
}

func NewCommentsRepository(db *DbConnection) core.CommentsRepository {
	return &commentsRepo{Db: db}
}

func (r *commentsRepo) Add(ctx core.ReqContext, comment *domain.Comment) (bool, error) {
	dbo := &CommentDBO{}
	dbo.FromComment(comment)
	query := `INSERT INTO comments (entitytype, entityid, date, parentid, threadid, userid, text, title, deleted)
		VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;`
	tr := ctx.StartTrace("CommentsRepository.Add")
	defer ctx.StopTrace(tr)
	row := r.Db.Conn.QueryRow(context.Background(), query, dbo.EntityType, dbo.EntityId, dbo.Date, dbo.ParentId, dbo.ThreadId, dbo.UserId, dbo.Text, dbo.Title, dbo.Deleted)
	err := row.Scan(&comment.Id)
	if err != nil {
		return false, r.Db.LogError(err, query)
	}
	return true, nil
}

func (r *commentsRepo) Update(ctx core.ReqContext, id int64, text, title string) (bool, error) {
	query := `UPDATE comments SET text=$1, title=$2, WHERE id=$3;`
	tr := ctx.StartTrace("CommentsRepository.Update")
	defer ctx.StopTrace(tr)
	tag, err := r.Db.Conn.Exec(context.Background(), query, text, ToNullString(title), id)
	if err != nil {
		return false, r.Db.LogError(err, query)
	}

	return tag.RowsAffected() > 0, nil
}

func (r *commentsRepo) Delete(ctx core.ReqContext, id int64) (bool, error) {
	query := `UPDATE comments SET deleted=$1 WHERE id=$2;`
	tr := ctx.StartTrace("CommentsRepository.Delete")
	defer ctx.StopTrace(tr)
	tag, err := r.Db.Conn.Exec(context.Background(), query, true, id)
	if err != nil {
		return false, r.Db.LogError(err, query)
	}

	return tag.RowsAffected() > 0, nil
}

func (r *commentsRepo) Get(ctx core.ReqContext, id int64) *domain.Comment {
	query := `SELECT id, entitytype, entityid, date, parentid, threadid, userid, text, title, deleted, points FROM comments WHERE id=$1;`
	tr := ctx.StartTrace("CommentsRepository.Get")
	defer ctx.StopTrace(tr)
	row := r.Db.Conn.QueryRow(context.Background(), query, id)
	p, err := r.scanRow(row)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		r.Db.LogError(err, query)
		return nil
	}
	return p.ToComment()
}

func (r *commentsRepo) GetThreadList(ctx core.ReqContext, entityType int, entityId int64, count int, page int) []domain.Comment {
	query := `SELECT id, entitytype, entityid, date, parentid, threadid, userid, text, title, deleted, points 
	FROM comments 
	WHERE entitytype=$1 
		AND entityid=$2 
		AND threadId is null
	ORDER BY id, points DESC
	LIMIT $3 OFFSET $4;`
	tr := ctx.StartTrace("CommentsRepository.GetThreadList")
	defer ctx.StopTrace(tr)
	rows, err := r.Db.Conn.Query(context.Background(), query, entityType, entityId, count, page*count)
	if err != nil {
		r.Db.LogError(err, query)
		return []domain.Comment{}
	}
	defer rows.Close()
	return r.scanRows(rows)
}

func (r *commentsRepo) GetThread(ctx core.ReqContext, entityType int, entityId int64, threadId int64) []domain.Comment {
	query := `SELECT id, entitytype, entityid, date, parentid, threadid, userid, text, title, deleted, points 
	FROM comments 
	WHERE entitytype = $1 
		AND entityid = $2 
		AND threadId = $3
	ORDER BY id;`
	tr := ctx.StartTrace("CommentsRepository.GetThread")
	defer ctx.StopTrace(tr)
	rows, err := r.Db.Conn.Query(context.Background(), query, entityType, entityId, threadId)
	if err != nil {
		r.Db.LogError(err, query)
		return []domain.Comment{}
	}
	defer rows.Close()
	return r.scanRows(rows)
}

func (r *commentsRepo) scanRows(rows pgx.Rows) []domain.Comment {
	comments := make([]domain.Comment, 0)
	for rows.Next() {
		dbo, err := r.scanRow(rows)
		if err != nil {
			return []domain.Comment{}
		}
		comments = append(comments, *dbo.ToComment())
	}
	return comments
}

func (r *commentsRepo) scanRow(row pgx.Row) (*CommentDBO, error) {
	dbo := CommentDBO{}
	err := row.Scan(&dbo.Id, &dbo.EntityType, &dbo.EntityId, &dbo.Date, &dbo.ParentId, &dbo.ThreadId, &dbo.UserId, &dbo.Text, &dbo.Title, &dbo.Deleted, &dbo.Points)
	if err != nil && err.Error() == "no rows in result set" {
		return &dbo, sql.ErrNoRows
	}
	return &dbo, err
}
