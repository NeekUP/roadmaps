package db

import (
	"context"
	"database/sql"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/jackc/pgx/v4"
)

type sourceRepo struct {
	Db *DbConnection
}

func NewSourceRepository(db *DbConnection) core.SourceRepository {
	return &sourceRepo{
		Db: db}
}

func (repo *sourceRepo) Get(ctx core.ReqContext, id int64) *domain.Source {
	tr := ctx.StartTrace("SourceRepository.Get")
	defer ctx.StopTrace(tr)
	query := "SELECT id, title, identifier, normalizedidentifier, type, properties, img, description FROM sources WHERE id=$1;"
	row := repo.Db.Conn.QueryRow(context.Background(), query, id)
	dbo, err := repo.scanRow(row)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		repo.Db.LogError(err, query)
		return nil
	}
	return dbo.ToSource()
}

func (repo *sourceRepo) FindByIdentifier(ctx core.ReqContext, nIdentifier string) *domain.Source {
	tr := ctx.StartTrace("SourceRepository.FindByIdentifier")
	defer ctx.StopTrace(tr)
	query := "SELECT id, title, identifier, normalizedidentifier, type, properties, img, description FROM sources WHERE normalizedidentifier=$1;"
	row := repo.Db.Conn.QueryRow(context.Background(), query, nIdentifier)
	dbo, err := repo.scanRow(row)

	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		repo.Db.LogError(err, query)
		return nil
	}
	return dbo.ToSource()
}

func (repo *sourceRepo) Save(ctx core.ReqContext, source *domain.Source) (bool, *core.AppError) {
	tr := ctx.StartTrace("SourceRepository.Save")
	defer ctx.StopTrace(tr)
	dbo := &SourceDBO{}
	dbo.FromSource(source)
	query := `INSERT INTO sources(
		title, identifier, normalizedidentifier, type, properties, img, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
	row := repo.Db.Conn.QueryRow(context.Background(), query, dbo.Title, dbo.Identifier, dbo.NormalizedIdentifier, dbo.Type, dbo.Properties, dbo.Img, dbo.Desc)
	err := row.Scan(&source.Id)
	if err != nil {
		return false, repo.Db.LogError(err, query)
	}
	return true, nil
}

func (repo *sourceRepo) Update(ctx core.ReqContext, source *domain.Source) (bool, *core.AppError) {
	tr := ctx.StartTrace("SourceRepository.Update")
	defer ctx.StopTrace(tr)
	dbo := &SourceDBO{}
	dbo.FromSource(source)
	query := `UPDATE sources
		SET title=$2, identifier=$3, normalizedidentifier=$4, type=$5, properties=$6, img=$7, description=$8
		WHERE id=$1;`
	tag, err := repo.Db.Conn.Exec(context.Background(), query, dbo.Id, dbo.Title, dbo.Identifier, dbo.NormalizedIdentifier, dbo.Type, dbo.Properties, dbo.Img, dbo.Desc)
	if err != nil {
		return false, repo.Db.LogError(err, query)
	}
	return tag.RowsAffected() > 0, nil
}

func (repo *sourceRepo) GetOrAddByIdentifier(ctx core.ReqContext, source *domain.Source) *domain.Source {
	tr := ctx.StartTrace("SourceRepository.GetOrAddByIdentifier")
	defer ctx.StopTrace(tr)
	dbo := &SourceDBO{}
	dbo.FromSource(source)
	query := `
INSERT INTO sources(
	title, identifier, normalizedidentifier, type, properties, img, description)
	VALUES ($1, $2, $3, $4, $5, $6, $7) 
ON CONFLICT (normalizedidentifier) DO NOTHING;`

	_, err := repo.Db.Conn.Exec(context.Background(), query, dbo.Title, dbo.Identifier, dbo.NormalizedIdentifier, dbo.Type, dbo.Properties, dbo.Img, dbo.Desc)
	if err != nil {
		repo.Db.LogError(err, query)
		return nil
	}
	return repo.FindByIdentifier(ctx, dbo.NormalizedIdentifier)
}

func (repo *sourceRepo) All() []domain.Source {
	query := "SELECT id, title, identifier, normalizedidentifier, type, properties, img, description FROM sources;"
	rows, err := repo.Db.Conn.Query(context.Background(), query)
	if err != nil {
		return []domain.Source{}
	}
	defer rows.Close()
	sources := make([]domain.Source, 0)
	for rows.Next() {
		dbo, err := repo.scanRow(rows)
		if err == sql.ErrNoRows {
			return nil
		}
		if err != nil {
			repo.Db.LogError(err, query)
			return nil
		}
		sources = append(sources, *dbo.ToSource())
	}

	return sources
}

func (repo *sourceRepo) scanRow(row pgx.Row) (*SourceDBO, error) {
	dbo := SourceDBO{}
	err := row.Scan(&dbo.Id, &dbo.Title, &dbo.Identifier, &dbo.NormalizedIdentifier, &dbo.Type, &dbo.Properties, &dbo.Img, &dbo.Desc)
	if err != nil && err.Error() == "no rows in result set" {
		return &dbo, sql.ErrNoRows
	}
	return &dbo, err
}
