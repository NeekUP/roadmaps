package usecases

import (
	"strings"

	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"

	"github.com/google/uuid"
)

type RegisterUser interface {
	Do(ctx core.ReqContext, name, email, password string) (*domain.User, error)
}

type registerUser struct {
	userRepo     core.UserRepository
	log          core.AppLogger
	hash         core.HashProvider
	emailService core.EmailSender
	imgManage    core.ImageManager
}

func NewRegisterUser(userRepo core.UserRepository, emailService core.EmailSender, hash core.HashProvider, imgManager core.ImageManager, log core.AppLogger) RegisterUser {
	return &registerUser{
		userRepo:     userRepo,
		emailService: emailService,
		hash:         hash,
		imgManage:    imgManager,
		log:          log,
	}
}

func (usecase *registerUser) Do(ctx core.ReqContext, name string, email string, password string) (*domain.User, error) {
	trace := ctx.StartTrace("registerUser")
	defer ctx.StopTrace(trace)

	appErr := usecase.validate(ctx, name, email, password)
	if appErr != nil {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"email", email,
			"error", appErr.Error(),
		)
		return nil, appErr
	}

	if ok := core.IsValidEmailHost(email); !ok {
		err := core.NewError(core.BadEmail)
		usecase.log.Errorw("Not valid email host",
			"reqid", ctx.ReqId(),
			"email", email,
			"error", err.Error(),
		)
		return nil, appErr
	}

	var avatarName string
	avatar, err := usecase.imgManage.GenerateAvatar(name)
	if err != nil {
		usecase.log.Errorw("Fail to generate avatar",
			"reqid", ctx.ReqId(),
			"email", email,
			"error", err.Error(),
			"name", name,
		)
	} else {
		avatarName = uuid.New().String() + ".png"
	}

	hash, salt := usecase.hash.HashPassword(password)
	user := &domain.User{
		Id:                name,
		Name:              name,
		NormalizedName:    strings.ToUpper(name),
		Email:             email,
		Rights:            domain.U,
		Pass:              hash,
		Salt:              salt,
		Img:               avatarName,
		EmailConfirmation: uuid.New().String(),
	}

	if _, err := usecase.userRepo.Save(ctx, user); err != nil {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"error", err.Error(),
		)
		return nil, err
	}

	if avatarName != "" {
		err = usecase.imgManage.SaveAvatar(avatar, avatarName)
		if err != nil {
			usecase.log.Errorw("Fail to generate avatar",
				"reqid", ctx.ReqId(),
				"email", email,
				"error", err.Error(),
				"name", name,
			)
		}
	}

	go usecase.emailService.Registration(email, user.Id, user.EmailConfirmation)

	return user, nil
}

func (usecase *registerUser) validate(ctx core.ReqContext, name string, email string, password string) *core.AppError {

	errors := make(map[string]string)

	if !core.IsValidUserName(name) {
		errors["name"] = core.InvalidFormat.String()
	}

	if !core.IsValidEmail(email) {
		errors["email"] = core.InvalidFormat.String()
	}

	if !core.IsValidPassword(password) {
		errors["pass"] = core.InvalidFormat.String()
	}

	if exists, ok := usecase.userRepo.ExistsName(ctx, name); ok && exists {
		errors["name"] = core.AlreadyExists.String()
	}

	if exists, ok := usecase.userRepo.ExistsEmail(ctx, email); ok && exists {
		errors["email"] = core.AlreadyExists.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}

	return nil
}
