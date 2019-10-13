package usecases

import (
	"roadmaps/core"
	"roadmaps/domain"

	"github.com/google/uuid"
)

type RegisterUser interface {
	Do(ctx core.ReqContext, name, email, password string) (*domain.User, error)
}

type registerUser struct {
	UserRepo core.UserRepository
	Log      core.AppLogger
	Hash     core.HashProvider
}

func NewRegisterUser(userRepo core.UserRepository, log core.AppLogger, hash core.HashProvider) RegisterUser {
	return &registerUser{
		UserRepo: userRepo,
		Log:      log,
		Hash:     hash}
}

func (r *registerUser) Do(ctx core.ReqContext, name string, email string, password string) (*domain.User, error) {

	appErr := r.validate(ctx, name, email, password)
	if appErr != nil {
		r.Log.Errorw("Not valid request",
			"reqId", ctx.ReqId(),
			"email", email,
			"error", appErr.Error(),
		)
		return nil, appErr
	}

	hash, salt := r.Hash.HashPassword(password)
	user := &domain.User{
		Id:     uuid.New().String(),
		Name:   name,
		Email:  email,
		Rights: domain.U}

	if r.UserRepo.Save(user, hash, salt) {
		return user, nil
	} else {
		return nil, core.NewError(core.InternalError)
	}
}

func (r *registerUser) validate(ctx core.ReqContext, name string, email string, password string) *core.AppError {

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

	if r.UserRepo.ExistsName(name) {
		errors["name"] = core.AlreadyExists.String()
	}

	if r.UserRepo.ExistsEmail(email) {
		errors["email"] = core.AlreadyExists.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}

	return nil
}
