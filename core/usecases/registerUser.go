package usecases

import (
	"github.com/google/uuid"
	"roadmaps/core"
	"roadmaps/domain"
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

	err := r.validate(ctx, name, email, password)
	if err != nil {
		return nil, err
	}

	hash, salt := r.Hash.HashPassword(password)
	user := &domain.User{
		Id:     uuid.New().String(),
		Name:   name,
		Email:  email,
		Rights: domain.U}

	if r.UserRepo.Create(user, hash, salt) {
		return user, nil
	} else {
		return nil, core.NewError(core.InternalError)
	}
}

func (r *registerUser) validate(ctx core.ReqContext, name string, email string, password string) error {

	if ok, c := core.IsValidUserName(name); !ok {
		r.Log.Infow("Username is not valid",
			"ReqId", ctx.ReqId(),
			"Email", email,
			"Name", name)
		return core.NewError(c)
	}

	if ok := core.IsValidEmail(email); !ok {
		r.Log.Infow("Username is not valid",
			"ReqId", ctx.ReqId(),
			"Email", email,
			"Name", name)
		return core.NewError(core.BadEmail)
	}

	if ok, c := core.IsValidPassword(password); !ok {
		r.Log.Infow("Password is not valid",
			"ReqId", ctx.ReqId(),
			"Email", email,
			"Name", name)
		return core.NewError(c)
	}

	if r.UserRepo.ExistsName(name) {
		r.Log.Infow("User with same name already registered",
			"ReqId", ctx.ReqId(),
			"Email", email,
			"Name", name)
		return core.NewError(core.NameAlreadyExists)
	}

	if r.UserRepo.ExistsEmail(email) {
		r.Log.Infow("User with same email already registered",
			"ReqId", ctx.ReqId(),
			"Email", email,
			"Name", name)
		return core.NewError(core.EmailAlreadyExists)
	}

	return nil
}
