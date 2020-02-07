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
	userRepo core.UserRepository
	log      core.AppLogger
	hash     core.HashProvider
}

func NewRegisterUser(userRepo core.UserRepository, log core.AppLogger, hash core.HashProvider) RegisterUser {
	return &registerUser{
		userRepo: userRepo,
		log:      log,
		hash:     hash,
	}
}

func (usecase *registerUser) Do(ctx core.ReqContext, name string, email string, password string) (*domain.User, error) {

	appErr := usecase.validate(ctx, name, email, password)
	if appErr != nil {
		usecase.log.Errorw("Not valid request",
			"reqId", ctx.ReqId(),
			"email", email,
			"error", appErr.Error(),
		)
		return nil, appErr
	}

	hash, salt := usecase.hash.HashPassword(password)
	user := &domain.User{
		Id:             uuid.New().String(),
		Name:           name,
		NormalizedName: strings.ToUpper(name),
		Email:          email,
		Rights:         domain.U}

	user.Pass = hash
	user.Salt = salt
	if _, err := usecase.userRepo.Save(user); err != nil {
		usecase.log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", err.Error(),
		)
		return nil, err
	}

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

	if exists, ok := usecase.userRepo.ExistsName(name); ok && exists {
		errors["name"] = core.AlreadyExists.String()
	}

	if exists, ok := usecase.userRepo.ExistsEmail(email); ok && exists {
		errors["email"] = core.AlreadyExists.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}

	return nil
}
