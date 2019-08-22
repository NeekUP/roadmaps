package usecases

import (
	"roadmaps/core"
	"roadmaps/domain"
)

type LoginUser interface {
	Do(tx core.ReqContext, email, password, fingerprint, useragent string) (*domain.User, string, string, error)
}

func NewLoginUser(ur core.UserRepository, log core.AppLogger, hash core.HashProvider, ts core.TokenService) LoginUser {
	return &loginUser{UserRepo: ur, Log: log, Hash: hash, TokenService: ts}
}

type loginUser struct {
	UserRepo     core.UserRepository
	Log          core.AppLogger
	Hash         core.HashProvider
	TokenService core.TokenService
}

func (this *loginUser) Do(ctx core.ReqContext, email, password, fingerprint, useragent string) (*domain.User, string, string, error) {
	err := this.validate(ctx, email, password)
	if err != nil {
		return nil, "", "", core.NewError(core.AuthenticationError)
	}

	useragent = core.UserAgentFingerprint(useragent)

	user := this.UserRepo.FindByEmail(email)
	if user == nil {
		this.Log.Infow("User not found",
			"ReqId", ctx.ReqId(),
			"Email", email)
		return nil, "", "", core.NewError(core.AuthenticationError)
	}

	if !this.Hash.CheckPassword(password, user.Pass, user.Salt) {
		this.Log.Infow("Password is wrong",
			"ReqId", ctx.ReqId(),
			"Email", email)
		return nil, "", "", core.NewError(core.AuthenticationError)
	}

	aToken, rToken, err := this.TokenService.Create(user, fingerprint, useragent)
	if err != nil {
		this.Log.Errorw("Fail to create token pair",
			"ReqId", ctx.ReqId(),
			"Email", email,
			"Error", err.Error())
		return nil, "", "", core.NewError(core.AuthenticationError)
	}

	return user, aToken, rToken, nil
}

func (r *loginUser) validate(ctx core.ReqContext, email string, password string) error {

	if ok, c := core.IsValidEmail(email); !ok {
		r.Log.Infow("Username is not valid",
			"ReqId", ctx.ReqId(),
			"Email", email)
		return core.NewError(c)
	}

	if ok, c := core.IsValidPassword(password); !ok {
		r.Log.Infow("Password is not valid",
			"ReqId", ctx.ReqId(),
			"Email", email)
		return core.NewError(c)
	}

	return nil
}
