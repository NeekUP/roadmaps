package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/google/uuid"
	"strings"
)

type RegisterUserOauth interface {
	Do(ctx core.ReqContext, name string, email string, providerName string, id string) (*domain.User, error)
}

type registerUserOauth struct {
	userRepo  core.UserRepository
	hash      core.HashProvider
	log       core.AppLogger
	imgManage core.ImageManager
}

func NewRegisterUserOauth(userRepo core.UserRepository, hash core.HashProvider, imgManager core.ImageManager, log core.AppLogger) RegisterUserOauth {
	return &registerUserOauth{
		userRepo:  userRepo,
		hash:      hash,
		imgManage: imgManager,
		log:       log,
	}
}

func (usecase *registerUserOauth) Do(ctx core.ReqContext, name, email, providerName, id string) (*domain.User, error) {
	trace := ctx.StartTrace("registerUserOauth")
	defer ctx.StopTrace(trace)

	appErr := usecase.validate(ctx, name, email)
	if appErr != nil {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"email", email,
			"error", appErr.Error(),
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

	hash, salt := usecase.hash.HashPassword(uuid.New().String())
	user := &domain.User{
		Id:                uuid.New().String(),
		Name:              name,
		NormalizedName:    strings.ToUpper(name),
		Email:             email,
		Rights:            domain.U,
		Pass:              hash,
		Salt:              salt,
		Img:               avatarName,
		EmailConfirmation: uuid.New().String(),
		EmailConfirmed:    true,
	}

	if _, err := usecase.userRepo.Save(ctx, user); err != nil {
		usecase.log.Errorw("invalid request",
			"reqid", ctx.ReqId(),
			"error", err.Error(),
		)
		return nil, err
	}

	if ok, err := usecase.userRepo.AddOauth(ctx, user.Id, providerName, id); !ok {
		if err != nil {
			usecase.log.Errorw("Oauth not added",
				"reqid", ctx.ReqId(),
				"error", err.Error(),
			)
		}
		usecase.userRepo.Delete(ctx, user.Id)
		return nil, core.NewError(core.InvalidRequest)
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

	user.OAuth = true
	return user, nil
}

func (usecase *registerUserOauth) validate(ctx core.ReqContext, name, email string) *core.AppError {

	errors := make(map[string]string)

	if !core.IsValidUserName(name) {
		errors["name"] = core.InvalidFormat.String()
	}

	if !core.IsValidEmail(email) {
		errors["email"] = core.InvalidFormat.String()
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
