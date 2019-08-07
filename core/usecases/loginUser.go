package usecases

import (
	"context"
	"roadmaps/core"
	"roadmaps/domain"
	"roadmaps/infrastructure"
)

type LoginUser interface {
	Do(ctx context.Context, email, password string) (*domain.User, error)
}

type loginUser struct {
	UserRepo     core.UserRepository
	Log          infrastructure.AppLogger
	Hash         core.HashProvider
	EmailChecker core.EmailChecker
	TokenService infrastructure.JwtTokenService
}

func (this *loginUser) Do(ctx context.Context, email, password, fingerprint, useragent string) (*domain.User, string, string, error) {
	err := this.validate(ctx, email, password)
	if err != nil {
		return nil, "", "", err
	}

	useragent = infrastructure.UserAgentFingerprint(useragent)

	user := this.UserRepo.FindByEmail(email)
	if user == nil {
		this.Log.Infow("User not found",
			"ReqId", infrastructure.GetReqID(ctx),
			"Email", email)
		return nil, "", "", core.NewError(core.AuthenticationError)
	}

	if !this.Hash.CheckPassword(password, user.Pass, user.Salt) {
		this.Log.Infow("Password is wrong",
			"ReqId", infrastructure.GetReqID(ctx),
			"Email", email)
		return nil, "", "", core.NewError(core.AuthenticationError)
	}

	aToken, rToken, err := this.TokenService.Create(user, fingerprint, useragent)
	if err != nil {
		this.Log.Errorw("Fail to create token pair",
			"ReqId", infrastructure.GetReqID(ctx),
			"Email", email,
			"Error", err.Error())
		return nil, "", "", core.NewError(core.AuthenticationError)
	}

	return user, aToken, rToken, nil
}

func (r *loginUser) validate(ctx context.Context, email string, password string) error {

	if ok := r.EmailChecker.IsValid(email); !ok {
		r.Log.Infow("Username is not valid",
			"ReqId", infrastructure.GetReqID(ctx),
			"Email", email)
		return core.NewError(core.BadEmail)
	}

	if ok, c := core.IsValidPassword(password); !ok {
		r.Log.Infow("Password is not valid",
			"ReqId", infrastructure.GetReqID(ctx),
			"Email", email)
		return core.NewError(c)
	}

	return nil
}
