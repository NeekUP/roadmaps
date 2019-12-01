package db

import (
	"context"
	"database/sql"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"log"
)

type sourceRepo struct {
	Db *DbConnection
}

func NewSourceRepository(db *DbConnection) core.SourceRepository {
	return &sourceRepo{
		Db: db}
}

func (repo *sourceRepo) Get(id int64) *domain.Source {
	query := "SELECT id, title, identifier, normalizedidentifier, type, properties, img, description FROM sources WHERE id=$1;"
	row := repo.Db.Conn.QueryRow(context.Background(), query, id)
	dbo, err := repo.scanRow(row)
	if err != nil {
		log.Fatal(err)
	}
	if err == sql.ErrNoRows {
		return nil
	}
	return dbo.ToSource()
}

func (repo *sourceRepo) FindByIdentifier(nIdentifier string) *domain.Source {
	query := "SELECT id, title, identifier, normalizedidentifier, type, properties, img, description FROM sources WHERE normalizedidentifier=$1;"
	row := repo.Db.Conn.QueryRow(context.Background(), query, nIdentifier)
	dbo, err := repo.scanRow(row)

	if err != nil {
		return nil
	}

	return dbo.ToSource()
}

func (repo *sourceRepo) Save(source *domain.Source) (bool, *core.AppError) {
	dbo := &SourceDBO{}
	dbo.FromSource(source)
	query := `INSERT INTO sources(
		title, identifier, normalizedidentifier, type, properties, img, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
	row := repo.Db.Conn.QueryRow(context.Background(), query, dbo.Title, dbo.Identifier, dbo.NormalizedIdentifier, dbo.Type, dbo.Properties, dbo.Img, dbo.Desc)
	err := row.Scan(&source.Id)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.Code == "23505" {
				return false, core.NewError(core.AlreadyExists)
			}
		}
		return false, core.NewError(core.InternalError)
	}
	return true, nil
}

func (repo *sourceRepo) Update(source *domain.Source) (bool, *core.AppError) {
	dbo := &SourceDBO{}
	dbo.FromSource(source)
	query := `UPDATE sources
		SET title=$2, identifier=$3, normalizedidentifier=$4, type=$5, properties=$6, img=$7, description=$8
		WHERE id=$1;`
	tag, err := repo.Db.Conn.Exec(context.Background(), query, dbo.Id, dbo.Title, dbo.Identifier, dbo.NormalizedIdentifier, dbo.Type, dbo.Properties, dbo.Img, dbo.Desc)
	if err != nil {
		return false, core.NewError(core.InternalError)
	}
	return tag.RowsAffected() > 0, nil
}

func (repo *sourceRepo) GetOrAddByIdentifier(source *domain.Source) *domain.Source {
	dbo := &SourceDBO{}
	dbo.FromSource(source)
	query := `
INSERT INTO sources(
	title, identifier, normalizedidentifier, type, properties, img, description)
	VALUES ($1, $2, $3, $4, $5, $6, $7) 
ON CONFLICT (normalizedidentifier) DO NOTHING;`

	_, err := repo.Db.Conn.Exec(context.Background(), query, dbo.Title, dbo.Identifier, dbo.NormalizedIdentifier, dbo.Type, dbo.Properties, dbo.Img, dbo.Desc)
	if err != nil {
		return nil
	}
	return repo.FindByIdentifier(dbo.NormalizedIdentifier)
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
		if err != nil {
			log.Fatal(err)
		}
		sources = append(sources, *dbo.ToSource())
	}

	return sources
}

func (repo *sourceRepo) scanRow(row pgx.Row) (*SourceDBO, error) {
	dbo := SourceDBO{}
	err := row.Scan(&dbo.Id, &dbo.Title, &dbo.Identifier, &dbo.NormalizedIdentifier, &dbo.Type, &dbo.Properties, &dbo.Img, &dbo.Desc)
	return &dbo, err
}
