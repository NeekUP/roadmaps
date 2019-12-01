package db

import (
	"context"
	"database/sql"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"log"
)

type userRepository struct {
	Db *DbConnection
}

func NewUserRepository(db *DbConnection) core.UserRepository {
	return &userRepository{Db: db}
}

func (r *userRepository) Get(id string) *domain.User {
	row := r.Db.Conn.QueryRow(context.Background(), "SELECT id, name, normalizedname, email, emailconfirmed, emailconfirmation, img, tokens, rights, password, salt "+
		"FROM users where id=$1", id)
	dbo, err := r.scanRow(row)
	if err != nil {
		log.Fatal(err)
	}
	if err == sql.ErrNoRows {
		return nil
	}
	return dbo.ToUser()
}

func (r *userRepository) Save(user *domain.User) (bool, *core.AppError) {
	dbo := UserDBO{}
	dbo.FromUser(user)
	dbo.Id = uuid.New().String()
	query := "INSERT INTO users " +
		"	(id, name, normalizedname, email, emailconfirmed, emailconfirmation, img, tokens, rights, password, salt) " +
		"VALUES " +
		"	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) " +
		"RETURNING id;"

	row := r.Db.Conn.QueryRow(context.Background(), query, dbo.Id, dbo.Name, dbo.NormalizedName, dbo.Email, dbo.EmailConfirmed, dbo.EmailConfirmation, dbo.Img, dbo.Tokens, dbo.Rights, dbo.Pass, dbo.Salt)
	err := row.Scan(&user.Id)

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

func (r *userRepository) Update(user *domain.User) (bool, *core.AppError) {
	dbo := &UserDBO{}
	dbo.FromUser(user)
	query := "UPDATE users " +
		"SET name=$1, normalizedname=$2, email=$3, emailconfirmed=$4, emailconfirmation=$5, img=$6, tokens=$7, rights=$8, password=$9, salt=$10 " +
		"WHERE id = $11;"

	tag, err := r.Db.Conn.Exec(context.Background(), query, dbo.Name, dbo.NormalizedName, dbo.Email, dbo.EmailConfirmed, dbo.EmailConfirmation, dbo.Img, dbo.Tokens, dbo.Rights, dbo.Pass, dbo.Salt, dbo.Id)
	if err != nil {
		return false, core.NewError(core.InternalError)
	}

	return tag.RowsAffected() > 0, nil
}

func (r *userRepository) ExistsName(name string) (exists bool, ok bool) {
	query := "select exists(select 1 from users where normalizedname=$1)"
	err := r.Db.Conn.QueryRow(context.Background(), query, name).Scan(&exists)
	if err != nil {
		return false, false
	}
	return exists, true
}

func (r *userRepository) ExistsEmail(email string) (exists bool, ok bool) {
	query := "select exists(select 1 from users where email=$1)"
	err := r.Db.Conn.QueryRow(context.Background(), query, email).Scan(&exists)
	if err != nil {
		return false, false
	}
	return exists, true
}

func (r *userRepository) FindByEmail(email string) *domain.User {
	query := "SELECT id, name, normalizedname, email, emailconfirmed, emailconfirmation, img, tokens, rights, password, salt " +
		"FROM users where email=$1"

	row := r.Db.Conn.QueryRow(context.Background(), query, email)
	dbo, err := r.scanRow(row)
	if err != nil {
		log.Fatal(err)
	}
	if err == sql.ErrNoRows {
		return nil
	}
	return dbo.ToUser()
}

func (r *userRepository) Count() (count int, ok bool) {
	query := "select count(id) from users;"
	err := r.Db.Conn.QueryRow(context.Background(), query).Scan(&count)
	if err != nil {
		return 0, false
	}
	return count, true
}

func (r *userRepository) All() []domain.User {
	query := "select id, name, normalizedname, email, emailconfirmed, emailconfirmation, img, tokens, rights, password, salt" +
		"FROM users"
	rows, err := r.Db.Conn.Query(context.Background(), query)
	if err != nil {
		return []domain.User{}
	}
	defer rows.Close()
	users := make([]domain.User, 0)
	for rows.Next() {
		dbo, err := r.scanRow(rows)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, *dbo.ToUser())
	}

	return users
}

func (r *userRepository) scanRow(row pgx.Row) (*UserDBO, error) {
	dbo := UserDBO{}
	err := row.Scan(&dbo.Id, &dbo.Name, &dbo.NormalizedName, &dbo.Email, &dbo.EmailConfirmed, &dbo.EmailConfirmation, &dbo.Img, &dbo.Tokens, &dbo.Rights, &dbo.Pass, &dbo.Salt)
	return &dbo, err
}
