package tests

import (
	"context"
	"fmt"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/core/usecases"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/NeekUP/roadmaps/infrastructure"
	"github.com/NeekUP/roadmaps/infrastructure/db"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var DB *db.DbConnection

func init() {
	DB = db.NewDbConnection("user=postgres password=1004287 host=localhost port=5432 database=roadmaps_tests sslmode=disable", newLogger("database_tests"))
}

func newLogger(name string) core.AppLogger {

	mainLogger := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s.log", "./log/", name),
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})

	encoderConfig := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey: "message",

		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,

		TimeKey:    "time",
		EncodeTime: zapcore.ISO8601TimeEncoder,

		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,

		EncodeDuration: zapcore.StringDurationEncoder,
	})

	logger := zap.New(zapcore.NewCore(encoderConfig, mainLogger, zapcore.DebugLevel))
	defer logger.Sync()

	return logger.Sugar()
}

func registerUser(name, email, pass string) *domain.User {
	regUserAction := usecases.NewRegisterUser(db.NewUserRepository(DB), appLoggerForTests{}, infrastructure.NewSha256HashProvider())
	u, _ := regUserAction.Do(infrastructure.NewContext(nil), name, email, pass)
	return u
}

func newContext(user *domain.User) core.ReqContext {
	ctx := context.Background()
	if user == nil {
		return infrastructure.NewContext(ctx)
	}

	ctx = context.WithValue(ctx, infrastructure.ReqRights, user.Rights)
	ctx = context.WithValue(ctx, infrastructure.ReqUserId, user.Id)
	ctx = context.WithValue(ctx, infrastructure.ReqUserName, user.Name)

	return infrastructure.NewContext(ctx)
}
