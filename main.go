package main

import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/NeekUP/roadmaps/api"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/core/usecases"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/NeekUP/roadmaps/infrastructure"
	"github.com/NeekUP/roadmaps/infrastructure/db"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
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
	cors := cors.New(cors.Options{
		AllowOriginFunc:  AllowOriginFunc,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	/*
		Middlewares
	*/
	r.Use(cors.Handler)
	r.Use(infrastructure.RequestID)
	r.Use(middleware.RealIP)
	r.Use(httpLogger(newLogger("http")))
	r.Use(recoverer(newLogger("recoverer")))
	r.Use(contentTypeMiddleware)

	/*
		Infrastructure initialization
	*/
	dbConnection := db.NewDbConnection(Cfg.Db.ConnString, newLogger("database"))
	hashProvider := infrastructure.NewSha256HashProvider()
	userRepo := db.NewUserRepository(dbConnection)
	sourceRepo := db.NewSourceRepository(dbConnection)
	topicRepo := db.NewTopicRepository(dbConnection)
	planRepo := db.NewPlansRepository(dbConnection)
	usersPlanRepo := db.NewUsersPlanRepository(dbConnection)
	captcha := infrastructure.SuccessCaptcha{}
	tokenService := infrastructure.NewJwtTokenService(userRepo, JwtSecret)
	imageManager := infrastructure.NewImageManager(Cfg.ImgSaver.LocalFolder, Cfg.ImgSaver.UriPath)

	//stepRepo := db.NewStepsRepository(dbConnection.Conn)

	/*
		Usecases
	*/
	regUser := usecases.NewRegisterUser(userRepo, newLogger("registerUser"), hashProvider)
	loginUser := usecases.NewLoginUser(userRepo, newLogger("loginUser"), hashProvider, tokenService)
	refreshToken := usecases.NewRefreshToken(userRepo, newLogger("refreshToken"), tokenService, JwtSecret)
	addSource := usecases.NewAddSource(sourceRepo, newLogger("addSource"), imageManager)
	addTopic := usecases.NewAddTopic(topicRepo, newLogger("addTopic"))
	addPlan := usecases.NewAddPlan(planRepo, newLogger("addPlan"))
	getTopic := usecases.NewGetTopic(topicRepo, planRepo, usersPlanRepo, newLogger("getTopic"))
	getPlanTree := usecases.NewGetPlanTree(planRepo, topicRepo, usersPlanRepo, newLogger("getPlanTree"))
	addUserPlan := usecases.NewAddUserPlan(planRepo, usersPlanRepo, newLogger("addUserPlan"))
	removeUserPlan := usecases.NewRemoveUserPlan(usersPlanRepo, newLogger("removeUserPlan"))
	getPlan := usecases.NewGetPlan(planRepo, userRepo, newLogger("getPlan"))
	getPlanList := usecases.NewGetPlanList(planRepo, userRepo, newLogger("getPlanList"))
	getUsersPlans := usecases.NewGetUsersPlans(planRepo, usersPlanRepo, newLogger("getUsersPlans"))
	/*
		Api methods
	*/
	apiReqUser := api.RegUser(regUser, newLogger("apiReqUser"), captcha)
	apiLoginUser := api.Login(loginUser, newLogger("apiLogin"), captcha)
	apiRefreshToken := api.RefreshToken(refreshToken, newLogger("apiRefreshToken"), captcha)
	apiAddSource := api.AddSource(addSource, newLogger("apiAddSource"))
	apiAddTopic := api.AddTopic(addTopic, newLogger("apiAddTopic"))
	apiAddPlan := api.AddPlan(addPlan, newLogger("apiAddPlan"))
	apiGetPlanTree := api.GetPlanTree(getPlanTree, newLogger("apiGetPlanTree"))
	apiGetTopicTree := api.GetTopicTree(getPlanTree, newLogger("apiGetTopicTree"))
	apiAddUserPlan := api.AddUserPlan(addUserPlan, newLogger("apiAddUserPlanTree"))
	apiRemoveAddUserPlan := api.RemoveUserPlan(removeUserPlan, newLogger("apiRemoveUserPlanTree"))
	apiGetTopic := api.GetTopic(getTopic, newLogger("apiGetTopic"))
	apiGetPlan := api.GetPlan(getPlan, newLogger("apiGetPlan"))
	apiGetPlanList := api.GetPlanList(getPlanList, getUsersPlans, newLogger("getPlanList"))
	/*
		Database
	*/
	dbSeed := infrastructure.NewDbSeed(regUser, userRepo)
	dbSeed.Seed()

	/*
		Http server
	*/

	// for all
	r.Group(func(r chi.Router) {
		r.Use(api.Auth(domain.All, tokenService))
		r.Post("/api/user/registration", apiReqUser)
		r.Post("/api/user/login", apiLoginUser)
		r.Post("/api/user/refresh", apiRefreshToken)
		r.Post("/api/plan/get", apiGetPlan)
		r.Post("/api/plan/list", apiGetPlanList)
		r.Post("/api/plan/tree", apiGetPlanTree)
		r.Post("/api/topic/tree", apiGetTopicTree)
		r.Post("/api/topic/get", apiGetTopic)
	})

	// for users
	r.Group(func(r chi.Router) {
		r.Use(api.Auth(domain.U, tokenService))
		r.Post("/api/source/add", apiAddSource)
		r.Post("/api/topic/add", apiAddTopic)
		r.Post("/api/plan/add", apiAddPlan)
		r.Post("/api/user/plan/favorite", apiAddUserPlan)
		r.Post("/api/user/plan/unfavorite", apiRemoveAddUserPlan)
	})

	// // for development only
	// listTopicsDev := usecases.NewListTopicsDev(topicRepo)
	// listPlansDev := usecases.NewListPlansDev(planRepo)
	// listStepsDev := usecases.NewListStepsDev(stepRepo)
	// listSourcesDev := usecases.NewListSourcesDev(sourceRepo)
	// listUsersDev := usecases.NewListUsersDev(userRepo)

	// apiListTopicsDev := api.ListTopics(listTopicsDev)
	// apiListPlansDev := api.ListPlans(listPlansDev)
	// apiListStepsDev := api.ListSteps(listStepsDev)
	// apiListSourcesDev := api.ListSources(listSourcesDev)
	// apiListUsersDev := api.ListUsers(listUsersDev)

	// r.Group(func(r chi.Router) {
	// 	r.Use(api.Auth(domain.All, tokenService))
	// 	r.Post("/api/dev/list/topics", apiListTopicsDev)
	// 	r.Post("/api/dev/list/plans", apiListPlansDev)
	// 	r.Post("/api/dev/list/steps", apiListStepsDev)
	// 	r.Post("/api/dev/list/source", apiListSourcesDev)
	// 	r.Post("/api/dev/list/users", apiListUsersDev)
	// })

	log.Printf("Listening %s", Cfg.HTTPServer.Host+":"+Cfg.HTTPServer.Port)

	srv := &http.Server{
		Handler:           r,
		Addr:              Cfg.HTTPServer.Host + ":" + Cfg.HTTPServer.Port,
		WriteTimeout:      time.Duration(Cfg.HTTPServer.WriteTimeoutSec) * time.Second,
		ReadTimeout:       time.Duration(Cfg.HTTPServer.ReadTimeoutSec) * time.Second,
		ReadHeaderTimeout: time.Duration(Cfg.HTTPServer.ReadHeadersTimeoutSec) * time.Second,
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

func recoverer(l core.AppLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {

					if l != nil {
						l.Panic(rvr, debug.Stack())
					} else {
						fmt.Fprintf(os.Stderr, "Panic:%v \r\n%s", rvr, string(debug.Stack()))
						debug.PrintStack()
					}

					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
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
func AllowOriginFunc(r *http.Request, origin string) bool {
	if origin == Cfg.Client.Host {
		return true
	}
	return false
}
