package db

import (
	"database/sql"
	_ "github.com/jackc/pgx"
)

type DbConnection struct {
	Db *sql.DB
}

func NewDbConnection(connString string) *DbConnection {
	//conn, err := sql.Open("postgre", connString)
	//if err != nil {
	//	panic(err)
	//}
	//
	dbConnection := new(DbConnection)
	dbConnection.Db = nil
	return dbConnection
	//return nil
}
