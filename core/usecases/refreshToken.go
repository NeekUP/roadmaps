package usecases

import (
	"context"
	"roadmaps/core"
	"roadmaps/domain"
	"roadmaps/infrastructure"
)

type RefreshToken interface {
	Do(ctx context.Context, authToken, refreshToken string) (*domain.User, error)
}

type refreshToken struct {
	UserRepo     core.UserRepository
	Log          infrastructure.AppLogger
	EmailChecker core.EmailChecker
	JwtSecret    string
	TokenService core.TokenService
}

func (this *refreshToken) Do(ctx context.Context, authToken, refreshToken, fingerprint, useragent string) (aToken string, rToken string, err error) {

	if ok, err := this.validate(authToken, refreshToken, fingerprint, useragent); !ok {
		this.Log.Infow("Not valid data",
			"ReqId", infrastructure.GetReqID(ctx),
			"authToken", authToken,
			"refreshToken", refreshToken,
			"fingerprint", fingerprint,
			"useragent", useragent,
			"error", err)
		return "", "", core.NewError(core.InvalidRequest)
	}
	useragent = infrastructure.UserAgentFingerprint(useragent)

	auth, refresh, err := this.TokenService.Refresh(authToken, refreshToken, fingerprint, useragent)

	if err != nil {
		this.Log.Infow("Fail to refresh token",
			"ReqId", infrastructure.GetReqID(ctx),
			"authToken", authToken,
			"refreshToken", refreshToken,
			"fingerprint", fingerprint,
			"useragent", useragent,
			"error", err)
		return "", "", core.NewError(core.InternalError)
	}
	return auth, refresh, err
}

func (this *refreshToken) validate(aToken, rToken, fingerprint, useragent string) (bool, error) {

	if ok, c := core.IsValidTokenFormat(aToken); !ok {
		return false, c
	}

	if ok, c := core.IsValidTokenFormat(rToken); !ok {
		return false, c
	}

	if ok, c := core.IsValidFingerprint(fingerprint); !ok {
		return false, c
	}

	if ok, c := core.IsValidUserAgent(useragent); !ok {
		return false, c
	}

	return true, nil
}
