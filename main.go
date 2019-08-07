package main

import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"roadmaps/api"
	"roadmaps/core/usecases"
	"roadmaps/infrastructure"
	"roadmaps/infrastructure/db"
	"runtime"
	"time"
)

var (
	AppLog    infrastructure.AppLogger
	Cfg       *infrastructure.Config
	JwtSecret = "ih7Cp1aB0exNXzsHjV9Z66qBczoG8g14_bBBW7iK1L-szDYVIbhWDZv6R-d_PD_TOjriomFr44UYMky2snKInO_7UL23uBmsH6hFlaqGJv12SQl4LC_1D7DW1iNLWSB22u1f3YowVH8YS_odqsUs5klaR7BlsvnQxucJcqSom6JuuZynz3j8p-8MevBDWTPAD7QeD4NUjTp55JftBEEg8J3Qf0ZrFOxkP2ULKvX-VbTwBN2U3YnNHJsdQ5aleUH-62NiG9EUiEDrLuEWw73oHaSCDPLVhIM1zCHW25Nmy8oxzW7rBVPwyLHC9v63QBSH7JXVhBOfDm-F55eOG0zlBw"
)

func init() {
	cfg, err := ioutil.ReadFile("conf.json")
	panicError(err)
	Cfg = initConfig(cfg)
	AppLog = initLogger("app")

	AppLog.Infow("Inited.", "time", time.Now().String())
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	r := chi.NewRouter()

	r.Use(infrastructure.RequestID)
	r.Use(middleware.RealIP)
	r.Use(httpLogger(initLogger("http")))
	r.Use(middleware.Recoverer)

	dbConnection := db.NewDbConnection(Cfg.Db.ConnString)
	hashProvider := infrastructure.NewSha256HashProvider()
	userRepo := db.NewUserRepository(dbConnection.Db)
	captcha := infrastructure.SuccessCaptcha{}
	regUser := usecases.NewRegisterUser(userRepo, initLogger("registerUser"), hashProvider)
	apiReqUser := api.RegUser(regUser, initLogger("apiReqUser"), captcha)
	r.Post("/", apiReqUser)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	_, err := fmt.Fprint(w, "Hello, World!")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func initConfig(dat []byte) *infrastructure.Config {
	var cfg infrastructure.Config
	err := json.Unmarshal(dat, &cfg)
	panicError(err)
	return &cfg
}

func initLogger(name string) infrastructure.AppLogger {

	mainLogger := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s.log", Cfg.Logger.Path, name),
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

func httpLogger(l infrastructure.AppLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				l.Info("",
					zap.Int("status", ww.Status()),
					zap.String("reqId", middleware.GetReqID(r.Context())),
					zap.String("proto", r.Proto),
					zap.String("path", r.URL.Path),
					zap.Duration("elapsed", time.Since(t1)),
					zap.Int("size", ww.BytesWritten()),
					zap.String("ip", r.RemoteAddr))
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}
