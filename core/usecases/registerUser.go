package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
	"strings"

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
		Id:             uuid.New().String(),
		Name:           name,
		NormalizedName: strings.ToUpper(name),
		Email:          email,
		Rights:         domain.U}

	user.Pass = hash
	user.Salt = salt
	if _, err := r.UserRepo.Save(user); err != nil {
		r.Log.Errorw("Not valid request",
			"ReqId", ctx.ReqId(),
			"Error", err.Error(),
		)
		return nil, err
	}
	return user, nil
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

	if exists, ok := r.UserRepo.ExistsName(name); ok && exists {
		errors["name"] = core.AlreadyExists.String()
	}

	if exists, ok := r.UserRepo.ExistsEmail(email); ok && exists {
		errors["email"] = core.AlreadyExists.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}

	return nil
}
