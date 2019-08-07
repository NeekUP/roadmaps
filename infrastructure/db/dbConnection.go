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
	//dbConnection := new(DbConnection)
	//dbConnection.Db = conn
	//return dbConnection
	return nil
}
