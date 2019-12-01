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

type topicRepo struct {
	Db *DbConnection
}

func NewTopicRepository(db *DbConnection) core.TopicRepository {
	return &topicRepo{Db: db}
}

func (repo *topicRepo) Get(name string) *domain.Topic {
	row := repo.Db.Conn.QueryRow(context.Background(), "SELECT id, name, title, description, creator FROM topics WHERE name=$1", name)
	dbo, err := repo.scanRow(row)
	if err != nil {
		log.Fatal(err)
	}
	if err == sql.ErrNoRows {
		return nil
	}
	return dbo.ToTopic()
}

func (repo *topicRepo) GetById(id int) *domain.Topic {
	row := repo.Db.Conn.QueryRow(context.Background(), "SELECT id, name, title, description, creator FROM topics WHERE id=$1", id)
	dbo, err := repo.scanRow(row)
	if err != nil {
		log.Fatal(err)
	}
	if err == sql.ErrNoRows {
		return nil
	}
	return dbo.ToTopic()
}

func (repo *topicRepo) Save(topic *domain.Topic) (bool, *core.AppError) {
	dbo := TopicDBO{}
	dbo.FromTopic(topic)
	query := "INSERT INTO topics( name, title, description, creator) VALUES ($1, $2, $3, $4) RETURNING id;"
	row := repo.Db.Conn.QueryRow(context.Background(), query, dbo.Name, dbo.Title, dbo.Description, dbo.Creator)
	err := row.Scan(&topic.Id)
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

func (repo *topicRepo) Update(topic *domain.Topic) (bool, *core.AppError) {
	dbo := TopicDBO{}
	dbo.FromTopic(topic)
	query := "UPDATE topics SET name=$2, title=$3, description=$4, creator=$5 WHERE id=$1;"
	tag, err := repo.Db.Conn.Exec(context.Background(), query, dbo.Id, dbo.Name, dbo.Title, dbo.Description, dbo.Creator)
	if err != nil {
		return false, core.NewError(core.InternalError)
	}

	return tag.RowsAffected() > 0, nil
}

func (repo *topicRepo) All() []domain.Topic {
	query := "SELECT id, name, title, description, creator FROM topics"
	rows, err := repo.Db.Conn.Query(context.Background(), query)
	if err != nil {
		return []domain.Topic{}
	}
	defer rows.Close()
	users := make([]domain.Topic, 0)
	for rows.Next() {
		dbo, err := repo.scanRow(rows)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, *dbo.ToTopic())
	}

	return users
}

func (repo *topicRepo) scanRow(row pgx.Row) (*TopicDBO, error) {
	dbo := TopicDBO{}
	err := row.Scan(&dbo.Id, &dbo.Name, &dbo.Title, &dbo.Description, &dbo.Creator)
	return &dbo, err
}
