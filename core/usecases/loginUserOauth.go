package usecases

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
)

type LoginUserOauth interface {
	Do(tx core.ReqContext, provider, openid, fingerprint, useragent string) (*domain.User, string, string, error)
}

func NewLoginUserOauth(ur core.UserRepository, ts core.TokenService, log core.AppLogger) LoginUserOauth {
	return &loginUserOauth{userRepo: ur, log: log, tokenService: ts}
}

type loginUserOauth struct {
	userRepo     core.UserRepository
	log          core.AppLogger
	tokenService core.TokenService
}

func (usecase *loginUserOauth) Do(ctx core.ReqContext, provider, openid, fingerprint, useragent string) (*domain.User, string, string, error) {
	trace := ctx.StartTrace("loginUserOauth")
	defer ctx.StopTrace(trace)

	useragent = core.UserAgentFingerprint(useragent)

	user := usecase.userRepo.FindByOauth(ctx, provider, openid)
	if user == nil {
		usecase.log.Infow("User not found",
			"reqid", ctx.ReqId(),
			"provider", provider,
			"openid", openid)
		return nil, "", "", core.NewError(core.AuthenticationError)
	}

	aToken, rToken, err := usecase.tokenService.Create(ctx, user, fingerprint, useragent)
	if err != nil {
		usecase.log.Errorw("Fail to create token pair",
			"reqid", ctx.ReqId(),
			"provider", provider,
			"openid", openid)
		return nil, "", "", core.NewError(core.AuthenticationError)
	}

	return user, aToken, rToken, nil
}

func (r *loginUserOauth) validate(ctx core.ReqContext, email string, password string) *core.AppError {

	errors := make(map[string]string)
	if !core.IsValidEmail(email) {
		errors["email"] = core.InvalidFormat.String()
	}

	if !core.IsValidPassword(password) {
		errors["pass"] = core.InvalidFormat.String()
	}

	if len(errors) > 0 {
		return core.ValidationError(errors)
	}

	return nil
}
