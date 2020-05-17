package db

import (
	"context"
	"fmt"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/infrastructure"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
)

type DbConnection struct {
	Conn *pgxpool.Pool
	Log  core.AppLogger
}

func NewDbConnection(config infrastructure.DbConf, log core.AppLogger) *DbConnection {
	pgdb, err := pgxpool.ConnectConfig(context.Background(), configure(config, log))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}

	return &DbConnection{
		Conn: pgdb,
		Log:  log,
	}
}

func configure(conf infrastructure.DbConf, log core.AppLogger) *pgxpool.Config {
	config, _ := pgxpool.ParseConfig(conf.ConnString)
	config.MaxConns = conf.PoolSettings.MaxConn
	config.HealthCheckPeriod = time.Second * time.Duration(conf.PoolSettings.HealthCheckSeconds)
	logLevel, err := pgx.LogLevelFromString(conf.LogLevel)
	if err != nil {
		log.Errorw("LogLevel value is invalid, default level: error")
		logLevel = pgx.LogLevelError
	}
	config.ConnConfig.Logger = &logger{Logger: log, Level: logLevel}
	return config
}

type logger struct {
	Logger core.AppLogger
	Level  pgx.LogLevel
}

func (driver *logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	arr := make([]string, len(data)*2)
	for k, v := range data {
		arr = append(arr, k, fmt.Sprintf("%v", v))
	}

	if driver.Level == 0 {
		return
	}

	switch level {
	case pgx.LogLevelNone:
		return
	case pgx.LogLevelError:
		if driver.Level >= 2 {
			driver.Logger.Errorw(msg, arr)
		}
	case pgx.LogLevelWarn:
		if driver.Level >= 3 {
			driver.Logger.Warnw(msg, arr)
		}
	case pgx.LogLevelInfo:
		if driver.Level >= 4 {
			driver.Logger.Infow(msg, arr)
		}
	case pgx.LogLevelDebug:
		if driver.Level >= 5 {
			driver.Logger.Debugw(msg, arr)
		}
	case pgx.LogLevelTrace:
		if driver.Level >= 6 {
			driver.Logger.Debugw(msg, arr)
		}
	}
}

func (db *DbConnection) LogError(err error, query string) *core.AppError {
	if pgerr, ok := err.(*pgconn.PgError); ok {
		db.Log.Errorw(pgerr.Message, "code", pgerr.Code, "hint", pgerr.Hint, "position", pgerr.Position, "query", query)
		if pgerr.Code == "23505" {
			return core.NewError(core.AlreadyExists)
		}
	} else {
		db.Log.Errorw(err.Error(), "query", query)
	}
	return core.NewError(core.InternalError)

}
