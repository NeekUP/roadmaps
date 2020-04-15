package main

import "C"
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NeekUP/nptrace"
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
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	AppLog core.AppLogger
	Cfg    *infrastructure.Config
	// TODO: remove
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
	**************************************/

	r.Use(cors.Handler)
	r.Use(infrastructure.RequestID)

	tracer := initTracer()
	if tracer != nil {
		r.Use(traceMiddleware(tracer))
	}

	r.Use(middleware.RealIP)
	r.Use(httpLogger(newLogger("http")))
	r.Use(recoverer(newLogger("recoverer")))

	/*
		Infrastructure initialization
	**************************************/

	dbConnection := db.NewDbConnection(Cfg.Db, newLogger("database"))
	cache := infrastructure.NewInMemoryCache()
	openAuthenticator := infrastructure.NewOpenAuthenticator(cache, Cfg.OAuth.ReturnUrl)
	hashProvider := infrastructure.NewSha256HashProvider()
	userRepo := db.NewUserRepository(dbConnection)
	sourceRepo := db.NewSourceRepository(dbConnection)
	topicRepo := db.NewTopicRepository(dbConnection)
	planRepo := db.NewPlansRepository(dbConnection)
	usersPlanRepo := db.NewUsersPlanRepository(dbConnection)
	captcha := infrastructure.SuccessCaptcha{}
	tokenService := infrastructure.NewJwtTokenService(userRepo, JwtSecret)
	imageManager := infrastructure.NewImageManager(Cfg.ImgSaver.LocalFolder, Cfg.ImgSaver.UriPath)
	stepRepo := db.NewStepsRepository(dbConnection)
	commentsRepo := db.NewCommentsRepository(dbConnection)
	pointsRepo := db.NewPointsRepository(dbConnection)
	changesRepository := db.NewChangeLogRepository(dbConnection)
	changeLog := infrastructure.NewChangesCollector(changesRepository, newLogger("changeLog"))
	emailService := infrastructure.NewEmailSender(Cfg.SiteHost, Cfg.SMTP.SenderEmail, Cfg.SMTP.SenderName, Cfg.SMTP.Host, Cfg.SMTP.Pass, Cfg.SMTP.Port, newLogger("emails"))

	api.ImgManager = imageManager

	for _, v := range Cfg.OAuth.Providers {
		openAuthenticator.AddProvider(v.Name, v.ClientId, v.Secret, v.Scope)
	}
	/*
		Usecases
	**************************************/

	// Users
	regUser := usecases.NewRegisterUser(userRepo, emailService, hashProvider, imageManager, newLogger("registerUser"))
	loginUser := usecases.NewLoginUser(userRepo, newLogger("loginUser"), hashProvider, tokenService)
	refreshToken := usecases.NewRefreshToken(userRepo, newLogger("refreshToken"), tokenService, JwtSecret)
	emailConfirmation := usecases.NewEmailConfirmation(userRepo, newLogger("emailConfirmation"))
	checkUser := usecases.NewCheckUser(userRepo, newLogger("checkUser"))
	registerUserOauth := usecases.NewRegisterUserOauth(userRepo, hashProvider, imageManager, newLogger("registerUserOauth"))
	loginUserOauth := usecases.NewLoginUserOauth(userRepo, tokenService, newLogger("loginUserOauth"))
	// Sources
	addSource := usecases.NewAddSource(sourceRepo, newLogger("addSource"), imageManager, changeLog)

	// Topics
	addTopic := usecases.NewAddTopic(topicRepo, changeLog, newLogger("addTopic"))
	getTopic := usecases.NewGetTopic(topicRepo, planRepo, usersPlanRepo, newLogger("getTopic"))
	searchTopic := usecases.NewSearchTopic(topicRepo, newLogger("getUsersPlans"))
	editTopic := usecases.NewEditTopic(topicRepo, changeLog, newLogger("editTopic"))

	// Plans
	addPlan := usecases.NewAddPlan(planRepo, changeLog, newLogger("addPlan"))
	getPlanTree := usecases.NewGetPlanTree(planRepo, topicRepo, stepRepo, usersPlanRepo, newLogger("getPlanTree"))
	getPlan := usecases.NewGetPlan(planRepo, userRepo, stepRepo, sourceRepo, topicRepo, newLogger("getPlan"))
	getPlanList := usecases.NewGetPlanList(planRepo, userRepo, newLogger("getPlanList"))
	editPlan := usecases.NewEditPlan(planRepo, changeLog, newLogger("editPlan"))
	removePlan := usecases.NewRemovePlan(planRepo, changeLog, newLogger("removePlan"))

	// Users Plans
	addUserPlan := usecases.NewAddUserPlan(planRepo, usersPlanRepo, newLogger("addUserPlan"))
	removeUserPlan := usecases.NewRemoveUserPlan(usersPlanRepo, newLogger("removeUserPlan"))
	getUsersPlans := usecases.NewGetUsersPlans(planRepo, usersPlanRepo, newLogger("getUsersPlans"))

	// Topic Tags
	addTopicTag := usecases.NewAddTopicTag(topicRepo, changeLog, newLogger("addTopicTag"))
	removeTopicTag := usecases.NewRemoveTopicTag(topicRepo, changeLog, newLogger("removeTopicTag"))

	// Comments
	addComment := usecases.NewAddComment(commentsRepo, planRepo, changeLog, newLogger("addComment"))
	editComment := usecases.NewEditComment(commentsRepo, changeLog, newLogger("editComment"))
	removeComment := usecases.NewRemoveComments(commentsRepo, changeLog, newLogger("removeComment"))
	getCommentsThreads := usecases.NewGetCommentsThreads(commentsRepo, userRepo, newLogger("getCommentsThreads"))
	getCommentsThread := usecases.NewGetCommentsThread(commentsRepo, userRepo, newLogger("getCommentsThread"))

	// Vote
	addPoints := usecases.NewAddPoints(pointsRepo, newLogger("addPoints"))
	getPoints := usecases.NewGetPoints(pointsRepo, newLogger("getPoints"))
	getPointsList := usecases.NewGetPointsList(pointsRepo, newLogger("getPointsList"))
	/*
		Api methods
	**************************************/

	// Users
	apiReqUser := api.RegUser(regUser, newLogger("registerUser"), captcha)
	apiLoginUser := api.Login(loginUser, newLogger("loginUser"), captcha)
	apiRefreshToken := api.RefreshToken(refreshToken, newLogger("refreshToken"), captcha)
	apiEmailConfirmation := api.EmailConfirmation(emailConfirmation, newLogger("emailConfirmation"))
	apiCheckUser := api.CheckUser(checkUser, newLogger("checkUser"))
	apiRegisterUserOauthLink := api.RegisterOAuthLink(checkUser, openAuthenticator, newLogger("registerUserOauth"))
	apiRegisterUserOauth := api.RegisterOAuth(registerUserOauth, loginUserOauth, openAuthenticator, newLogger("registerUserOauth"))
	apiLoginUserOauthLink := api.LoginOAuthLink(openAuthenticator, newLogger("loginUserOauth"))
	apiLoginUserOauth := api.LoginOauth(loginUserOauth, openAuthenticator, newLogger("loginUserOauth"))

	// Sources
	apiAddSource := api.AddSource(addSource, newLogger("addSource"))

	// Topics
	apiAddTopic := api.AddTopic(addTopic, newLogger("addTopic"))
	apiGetTopicTree := api.GetTopicTree(getPlanTree, newLogger("getTopicTree"))
	apiGetTopic := api.GetTopic(getTopic, newLogger("getTopic"))
	apiSearchTopic := api.SearchTopic(searchTopic, newLogger("searchTopic"))
	apiEditTopic := api.EditTopic(editTopic, newLogger("editTopic"))

	// Plans
	apiAddPlan := api.AddPlan(addPlan, newLogger("addPlan"))
	apiGetPlanTree := api.GetPlanTree(getPlanTree, newLogger("getPlanTree"))
	apiGetPlan := api.GetPlan(getPlan, getUsersPlans, getPoints, newLogger("getPlan"))
	apiGetPlanList := api.GetPlanList(getPlanList, getUsersPlans, getPointsList, newLogger("getPlanList"))
	apiEditPlan := api.EditPlan(editPlan, newLogger("editPlan"))
	apiRemovePlan := api.RemovePlan(removePlan, newLogger("removePlan"))

	// Users Plans
	apiAddUserPlan := api.AddUserPlan(addUserPlan, newLogger("addUserPlan"))
	apiRemoveAddUserPlan := api.RemoveUserPlan(removeUserPlan, newLogger("removeUserPlan"))

	// Topic Tag
	apiAddTopicTag := api.AddTopicTag(addTopicTag, newLogger("addTopicTag"))
	apiRemoveTopicTag := api.RemoveTopicTag(removeTopicTag, newLogger("removeTopicTag"))

	// Comments
	apiAddComment := api.AddComment(addComment, newLogger("addComment"))
	apiEditComment := api.EditComment(editComment, newLogger("editComment"))
	apiRemoveComment := api.DeleteComment(removeComment, newLogger("removeComment"))
	apiGetCommentsThreads := api.GetThreads(getCommentsThreads, getPointsList, newLogger("getCommentsThreads"))
	apiGetCommentsThread := api.GetThread(getCommentsThread, getPointsList, newLogger("getCommentsThread"))

	// Vote
	apiAddPoints := api.AddPoints(addPoints, newLogger("addPoints"))

	/*
		Database
	**************************************/

	dbSeed := infrastructure.NewDbSeed(regUser, userRepo)
	dbSeed.Seed()

	/*
		Http server
	**************************************/

	// for all
	r.Group(func(r chi.Router) {
		r.Use(api.Auth(domain.All, tokenService, newLogger("auth")))
		r.Post("/api/topic/tree", apiGetTopicTree)
		r.Post("/api/topic/get", apiGetTopic)
		r.Post("/api/topic/search", apiSearchTopic)

		r.Post("/api/plan/get", apiGetPlan)
		r.Post("/api/plan/list", apiGetPlanList)
		r.Post("/api/plan/tree", apiGetPlanTree)

		r.Post("/api/user/registration", apiReqUser)
		r.Post("/api/user/login", apiLoginUser)
		r.Post("/api/user/refresh", apiRefreshToken)
		r.Post("/api/user/check", apiCheckUser)
		r.Post("/api/user/oauth/registrationStart", apiRegisterUserOauthLink)
		r.Post("/api/user/oauth/registrationEnd", apiRegisterUserOauth)
		r.Post("/api/user/oauth/loginStart", apiLoginUserOauthLink)
		r.Post("/api/user/oauth/loginEnd", apiLoginUserOauth)
		r.Post("/api/comment/threads", apiGetCommentsThreads)
		r.Post("/api/comment/thread", apiGetCommentsThread)
		r.Get("/s/confirm", apiEmailConfirmation)
	})

	// for users
	r.Group(func(r chi.Router) {
		r.Use(api.Auth(domain.U, tokenService, newLogger("auth")))
		r.Post("/api/source/add", apiAddSource)
		r.Post("/api/topic/add", apiAddTopic)
		r.Post("/api/plan/add", apiAddPlan)
		r.Post("/api/plan/edit", apiEditPlan)
		r.Post("/api/plan/remove", apiRemovePlan)
		r.Post("/api/user/plan/favorite", apiAddUserPlan)
		r.Post("/api/user/plan/unfavorite", apiRemoveAddUserPlan)
		r.Post("/api/comment/add", apiAddComment)
		r.Post("/api/comment/edit", apiEditComment)
		r.Post("/api/comment/delete", apiRemoveComment)
		r.Post("/api/points/add", apiAddPoints)
	})

	// for moderators
	r.Group(func(r chi.Router) {
		r.Use(api.Auth(domain.M, tokenService, newLogger("auth")))
		r.Post("/api/topic/tag/add", apiAddTopicTag)
		r.Post("/api/topic/tag/remove", apiRemoveTopicTag)
		r.Post("/api/topic/edit", apiEditTopic)
	})

	// for development only
	listTopicsDev := usecases.NewListTopicsDev(topicRepo)
	listPlansDev := usecases.NewListPlansDev(planRepo)
	listStepsDev := usecases.NewListStepsDev(stepRepo)
	listSourcesDev := usecases.NewListSourcesDev(sourceRepo)
	listUsersDev := usecases.NewListUsersDev(userRepo)

	apiListTopicsDev := api.ListTopics(listTopicsDev)
	apiListPlansDev := api.ListPlans(listPlansDev)
	apiListStepsDev := api.ListSteps(listStepsDev)
	apiListSourcesDev := api.ListSources(listSourcesDev)
	apiListUsersDev := api.ListUsers(listUsersDev)

	r.Group(func(r chi.Router) {
		r.Use(api.Auth(domain.A, tokenService, newLogger("auth")))
		r.Post("/api/dev/list/topics", apiListTopicsDev)
		r.Post("/api/dev/list/plans", apiListPlansDev)
		r.Post("/api/dev/list/steps", apiListStepsDev)
		r.Post("/api/dev/list/source", apiListSourcesDev)
		r.Post("/api/dev/list/users", apiListUsersDev)
	})

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

func initTracer() *nptrace.NPTrace {
	filename := fmt.Sprintf("%s/trace.log", Cfg.Logger.Path)
	outCfg := nptrace.NewJsonEncoderConfig(time.StampMicro, func(d time.Duration) []byte {
		return []byte(strconv.FormatInt(d.Nanoseconds()/1000, 10))
	})
	cfg := nptrace.NewJsonEncoder(outCfg)
	var traceFile *os.File
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		traceFile, err = os.Create(filename)
		if err != nil {
			AppLog.Errorw("Fail to create file for performance tracing.", "file", filename, "err", err.Error())
		}
	} else {
		traceFile, err = os.OpenFile(filename, os.O_WRONLY, 0666)
		if err != nil {
			AppLog.Errorw("Fail to open file for performance tracing.", "file", filename, "err", err.Error())
		}
	}
	if traceFile != nil {
		return nptrace.NewTracer(cfg, traceFile)
	}

	return nil
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
					zap.String("reqid", middleware.GetReqID(r.Context())),
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
						l.Error(rvr, string(debug.Stack()))
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

func traceMiddleware(npTrace *nptrace.NPTrace) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			tracer := npTrace.New(ctx.Value(infrastructure.ReqId).(string), strings.Trim(r.URL.Path, "/"))
			defer npTrace.Close(tracer)

			ctx = context.WithValue(ctx, infrastructure.Tracer, tracer)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
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
