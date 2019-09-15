package main

import "C"
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"roadmaps/api"
	"roadmaps/core"
	"roadmaps/core/usecases"
	"roadmaps/domain"
	"roadmaps/infrastructure"
	"roadmaps/infrastructure/db"
	"runtime"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	AppLog    core.AppLogger
	Cfg       *infrastructure.Config
	JwtSecret = "ih7Cp1aB0exNXzsHjV9Z66qBczoG8g14_bBBW7iK1L-szDYVIbhWDZv6R-d_PD_TOjriomFr44UYMky2snKInO_7UL23uBmsH6hFlaqGJv12SQl4LC_1D7DW1iNLWSB22u1f3YowVH8YS_odqsUs5klaR7BlsvnQxucJcqSom6JuuZynz3j8p-8MevBDWTPAD7QeD4NUjTp55JftBEEg8J3Qf0ZrFOxkP2ULKvX-VbTwBN2U3YnNHJsdQ5aleUH-62NiG9EUiEDrLuEWw73oHaSCDPLVhIM1zCHW25Nmy8oxzW7rBVPwyLHC9v63QBSH7JXVhBOfDm-F55eOG0zlBw"
)

func init() {
	cfg, err := ioutil.ReadFile("conf.json")
	panicError(err)
	Cfg = initConfig(cfg)
	AppLog = newLogger("app")

	AppLog.Infow("Inited.", "time", time.Now().String())
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	r := chi.NewRouter()

	/*
		Middlewares
	*/
	r.Use(infrastructure.RequestID)
	r.Use(middleware.RealIP)
	r.Use(httpLogger(newLogger("http")))
	r.Use(middleware.Recoverer)
	r.Use(contentTypeMiddleware)

	/*
		Infrastructure initialization
	*/
	dbConnection := db.NewDbConnection(Cfg.Db.ConnString)
	hashProvider := infrastructure.NewSha256HashProvider()
	userRepo := db.NewUserRepository(dbConnection.Db)
	sourceRepo := db.NewSourceRepository(dbConnection.Db)
	captcha := infrastructure.SuccessCaptcha{}
	tokenService := infrastructure.NewJwtTokenService(userRepo, JwtSecret)
	imageManager := infrastructure.NewImageManager(Cfg.ImgSaver.LocalFolder, Cfg.ImgSaver.UriPath)

	/*
		Usecases
	*/
	regUser := usecases.NewRegisterUser(userRepo, newLogger("registerUser"), hashProvider)
	loginUser := usecases.NewLoginUser(userRepo, newLogger("loginUser"), hashProvider, tokenService)
	refreshToken := usecases.NewRefreshToken(userRepo, newLogger("refreshToken"), tokenService, JwtSecret)
	addSource := usecases.NewAddSource(sourceRepo, newLogger("addSource"), imageManager)
	/*
		Api methods
	*/
	apiReqUser := api.RegUser(regUser, newLogger("apiReqUser"), captcha)
	apiLoginUser := api.Login(loginUser, newLogger("apiLogin"), captcha)
	apiRefreshToken := api.RefreshToken(refreshToken, newLogger("apiRefreshToken"), captcha)
	apiAddSource := api.AddSource(addSource, newLogger("apiAddSource"))
	/*
		Database
	*/
	dbSeed := infrastructure.NewDbSeed(regUser, userRepo)
	dbSeed.Seed()

	/*
		Http server
	*/
	r.Group(func(r chi.Router) {
		// public
		r.Post("/user/reqistration", apiReqUser)
		r.Post("/user/login", apiLoginUser)
		r.Post("/user/refresh", apiRefreshToken)

		// not public
		r.Group(func(r chi.Router) {
			r.Use(api.Auth(domain.U, tokenService))
			r.Post("/source/add", apiAddSource)
		})
	})

	//port := os.Getenv("PORT")
	port := Cfg.HTTPServer.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)

	srv := &http.Server{
		Handler:      r,
		Addr:         Cfg.HTTPServer.Host + ":" + Cfg.HTTPServer.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func initConfig(dat []byte) *infrastructure.Config {
	var cfg infrastructure.Config
	err := json.Unmarshal(dat, &cfg)
	panicError(err)
	return &cfg
}

func newLogger(name string) core.AppLogger {

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

func httpLogger(l core.AppLogger) func(next http.Handler) http.Handler {
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

func contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}
