package db

import (
	"context"
	"fmt"
	"github.com/NeekUP/roadmaps/core"
	"os"

	"github.com/jackc/pgx/v4"
)

type DbConnection struct {
	Conn *pgx.Conn
	Log  core.AppLogger
}

func NewDbConnection(connString string, log core.AppLogger) *DbConnection {
	pgdb, err := pgx.ConnectConfig(context.Background(), configure(connString, log))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}

	err = pgdb.Ping(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}

	return &DbConnection{
		Conn: pgdb,
		Log:  log,
	}
}

func configure(cfg string, log core.AppLogger) *pgx.ConnConfig {
	config, _ := pgx.ParseConfig(cfg)
	config.Logger = &DbLogger{Logger: log}
	return config
}

type DbLogger struct {
	Logger core.AppLogger
}

func (driver *DbLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	arr := make([]string, len(data)*2)
	for k, v := range data {
		arr = append(arr, k, fmt.Sprintf("%v", v))
	}

	switch level {
	case pgx.LogLevelNone:
		return
	case pgx.LogLevelError:
		driver.Logger.Errorw(msg, arr)
	case pgx.LogLevelWarn:
		driver.Logger.Warnw(msg, arr)
	case pgx.LogLevelInfo:
		driver.Logger.Infow(msg, arr)
	case pgx.LogLevelDebug:
		driver.Logger.Debugw(msg, arr)
	case pgx.LogLevelTrace:
		driver.Logger.Debugw(msg, arr)
	}
}
