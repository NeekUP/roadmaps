package usecases

import (
	"context"
	"roadmaps/core"
	"roadmaps/domain"
	"roadmaps/infrastructure"
)

type RegisterUser interface {
	Do(ctx context.Context, name, email, password string) (*domain.User, error)
}

type registerUser struct {
	UserRepo core.UserRepository
	Log      infrastructure.AppLogger
	Hash     core.HashProvider
}

func NewRegisterUser(userRepo core.UserRepository, log infrastructure.AppLogger, hash core.HashProvider) RegisterUser {
	return &registerUser{
		UserRepo: userRepo,
		Log:      log,
		Hash:     hash}
}

func (r *registerUser) Do(ctx context.Context, name string, email string, password string) (*domain.User, error) {

	err := r.validate(ctx, name, email, password)
	if err != nil {
		return nil, err
	}

	hash, salt := r.Hash.HashPassword(password)
	user := &domain.User{
		Name:   name,
		Email:  email,
		Rights: domain.U}

	if r.UserRepo.Create(user, hash, salt) {
		return user, nil
	} else {
		return nil, core.NewError(core.InternalError)
	}
}

func (r *registerUser) validate(ctx context.Context, name string, email string, password string) error {

	if ok, c := core.IsValidUserName(name); !ok {
		r.Log.Infow("Username is not valid",
			"ReqId", infrastructure.GetReqID(ctx),
			"Email", email,
			"Name", name)
		return core.NewError(c)
	}

	if ok := core.IsValidEmail(email); !ok {
		r.Log.Infow("Username is not valid",
			"ReqId", infrastructure.GetReqID(ctx),
			"Email", email,
			"Name", name)
		return core.NewError(core.BadEmail)
	}

	if ok, c := core.IsValidPassword(password); !ok {
		r.Log.Infow("Password is not valid",
			"ReqId", infrastructure.GetReqID(ctx),
			"Email", email,
			"Name", name)
		return core.NewError(c)
	}

	if r.UserRepo.ExistsName(name) {
		r.Log.Infow("User with same name already registered",
			"ReqId", infrastructure.GetReqID(ctx),
			"Email", email,
			"Name", name)
		return core.NewError(core.NameAlreadyExists)
	}

	if r.UserRepo.ExistsEmail(email) {
		r.Log.Infow("User with same email already registered",
			"ReqId", infrastructure.GetReqID(ctx),
			"Email", email,
			"Name", name)
		return core.NewError(core.EmailAlreadyExists)
	}

	return nil
}
