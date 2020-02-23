package db

import (
	"context"
	"database/sql"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/jackc/pgx/v4"
)

type topicRepo struct {
	Db *DbConnection
}

func NewTopicRepository(db *DbConnection) core.TopicRepository {
	return &topicRepo{Db: db}
}

func (repo *topicRepo) Get(ctx core.ReqContext, name string) *domain.Topic {
	tr := ctx.StartTrace("TopicRepository.Get")
	defer ctx.StopTrace(tr)
	row := repo.Db.Conn.QueryRow(context.Background(), "SELECT id, name, title, description, creator, tags, istag FROM topics WHERE name=$1", name)
	dbo, err := repo.scanRow(row)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		repo.Db.Log.Errorw("", "error", err.Error())
	}

	topic := dbo.ToTopic(repo.GetTags(ctx, dbo.Tags))
	return topic
}

func (repo *topicRepo) GetById(ctx core.ReqContext, id int) *domain.Topic {
	tr := ctx.StartTrace("TopicRepository.GetById")
	defer ctx.StopTrace(tr)
	row := repo.Db.Conn.QueryRow(context.Background(), "SELECT id, name, title, description, creator, tags, istag FROM topics WHERE id=$1", id)
	dbo, err := repo.scanRow(row)

	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		repo.Db.Log.Errorw("", "error", err.Error())
		return nil
	}

	topic := dbo.ToTopic(repo.GetTags(ctx, dbo.Tags))
	return topic
}

func (repo *topicRepo) Save(ctx core.ReqContext, topic *domain.Topic) (bool, *core.AppError) {
	tr := ctx.StartTrace("TopicRepository.Save")
	defer ctx.StopTrace(tr)
	dbo := TopicDBO{}
	dbo.FromTopic(topic)
	query := "INSERT INTO topics( name, title, description, creator, tags, istag) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;"
	row := repo.Db.Conn.QueryRow(context.Background(), query, dbo.Name, dbo.Title, dbo.Description, dbo.Creator, dbo.Tags, dbo.IsTag)
	err := row.Scan(&topic.Id)
	if err != nil {
		return false, repo.Db.LogError(err, query)
	}

	if len(topic.Tags) > 0 {
		for _, tag := range topic.Tags {
			repo.AddTag(ctx, tag.Name, topic.Name)
		}
	}
	return true, nil
}

func (repo *topicRepo) Update(ctx core.ReqContext, topic *domain.Topic) (bool, *core.AppError) {
	tr := ctx.StartTrace("TopicRepository.Update")
	defer ctx.StopTrace(tr)
	dbo := TopicDBO{}
	dbo.FromTopic(topic)
	oldTopic := repo.GetById(ctx, topic.Id)
	if oldTopic == nil {
		return false, core.NewError(core.InternalError)
	}

	tx, err := repo.Db.Conn.BeginTx(context.Background(), pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})
	query := "UPDATE topics SET name=$2, title=$3, description=$4, istag=$5 WHERE id=$1;"
	tag, err := repo.Db.Conn.Exec(context.Background(), query, dbo.Id, dbo.Name, dbo.Title, dbo.Description, dbo.IsTag)
	if err != nil {
		if e := tx.Rollback(context.Background()); e != nil {
			repo.Db.Log.Errorw("Tx not rolled back", "err", e.Error())
		}
		return false, repo.Db.LogError(err, query)
	}

	if tag.RowsAffected() > 0 {
		_, err := repo.updateTags(oldTopic.Name, dbo.Name)
		if err != nil {
			if e := tx.Rollback(context.Background()); e != nil {
				repo.Db.Log.Errorw("Tx not rolled back", "err", e.Error())
			}
			return false, repo.Db.LogError(err, "")
		}
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return false, core.NewError(core.InternalError)
	}
	return true, nil
}

func (repo *topicRepo) updateTags(oldname, newname string) (bool, error) {
	query := `UPDATE topics SET tags = array_replace(tags,$1,$2)`
	_, err := repo.Db.Conn.Exec(context.Background(), query, oldname, newname)
	if err != nil {
		return false, repo.Db.LogError(err, query)
	}
	return true, nil
}

func (repo *topicRepo) All() []domain.Topic {
	query := "SELECT id, name, title, description, creator, tags, istag FROM topics"
	rows, err := repo.Db.Conn.Query(context.Background(), query)
	if err != nil {
		return []domain.Topic{}
	}
	defer rows.Close()
	topics := make([]domain.Topic, 0)
	for rows.Next() {
		dbo, err := repo.scanRow(rows)
		if err != nil {
			repo.Db.LogError(err, query)
			return []domain.Topic{}
		}
		topics = append(topics, *dbo.ToTopic([]domain.TopicTag{}))
	}

	return topics
}

func (repo *topicRepo) TitleLike(ctx core.ReqContext, str string, count int) []domain.Topic {
	tr := ctx.StartTrace("TopicRepository.TitleLike")
	defer ctx.StopTrace(tr)
	query := "SELECT id, name, title, description, creator, tags , istag FROM topics WHERE title ILIKE $1 LIMIT $2"
	rows, err := repo.Db.Conn.Query(context.Background(), query, "%"+str+"%", count)
	if err != nil {
		return []domain.Topic{}
	}
	defer rows.Close()
	topics := make([]domain.Topic, 0)
	for rows.Next() {
		dbo, err := repo.scanRow(rows)
		if err != nil {
			repo.Db.LogError(err, query)
			return []domain.Topic{}
		}
		topic := *dbo.ToTopic(repo.GetTags(ctx, dbo.Tags))
		topics = append(topics, topic)
	}

	return topics
}

func (repo *topicRepo) GetTags(ctx core.ReqContext, topicNames []string) []domain.TopicTag {
	tr := ctx.StartTrace("TopicRepository.GetTags")
	defer ctx.StopTrace(tr)

	if topicNames == nil || len(topicNames) == 0 {
		return []domain.TopicTag{}
	}

	query := `SELECT name,title FROM topics WHERE name=ANY($1) AND istag = true`

	rows, err := repo.Db.Conn.Query(context.Background(), query, topicNames)
	if err != nil {
		return []domain.TopicTag{}
	}
	defer rows.Close()
	tags := make([]domain.TopicTag, 0)
	for rows.Next() {
		dbo := TopicTagDBO{}
		err := rows.Scan(&dbo.Name, &dbo.Title)
		if err != nil {
			repo.Db.LogError(err, query)
			return []domain.TopicTag{}
		}
		tags = append(tags, *dbo.ToTopicTag())
	}

	return tags
}

func (repo *topicRepo) AddTag(ctx core.ReqContext, tagname, topicname string) bool {
	tr := ctx.StartTrace("TopicRepository.AddTag")
	defer ctx.StopTrace(tr)

	query := `UPDATE topics 
SET tags = array_cat(tags, $1) 
WHERE name=$2 
	AND array_position(tags, $3) IS NULL
	AND EXISTS( SELECT id FROM topics WHERE name=$3 AND istag = true );`
	t, err := repo.Db.Conn.Exec(context.Background(), query, []string{tagname}, topicname, tagname)
	if err != nil {
		repo.Db.LogError(err, query)
		return false
	}
	return t.RowsAffected() > 0
}

func (repo *topicRepo) DeleteTag(ctx core.ReqContext, tagname, topicname string) bool {
	tr := ctx.StartTrace("TopicRepository.DeleteTag")
	defer ctx.StopTrace(tr)

	query := `UPDATE topics SET tags = array_remove(tags, $1) WHERE name=$2;`
	_, err := repo.Db.Conn.Exec(context.Background(), query, tagname, topicname)
	if err != nil {
		repo.Db.LogError(err, query)
		return false
	}
	return true
}

func (repo *topicRepo) scanRow(row pgx.Row) (*TopicDBO, error) {
	dbo := TopicDBO{}
	err := row.Scan(&dbo.Id, &dbo.Name, &dbo.Title, &dbo.Description, &dbo.Creator, &dbo.Tags, &dbo.IsTag)
	if err != nil && err.Error() == "no rows in result set" {
		return &dbo, sql.ErrNoRows
	}
	return &dbo, err
}
