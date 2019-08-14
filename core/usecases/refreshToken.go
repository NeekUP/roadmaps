package usecases

import (
	"roadmaps/core"
)

type RefreshToken interface {
	Do(ctx core.ReqContext, authToken, refreshToken, fingerprint, useragent string) (aToken string, rToken string, err error)
}

func NewRefreshToken(ur core.UserRepository, log core.AppLogger, ts core.TokenService, secret string) RefreshToken {
	return &refreshToken{
		UserRepo:     ur,
		Log:          log,
		Secret:       secret,
		TokenService: ts,
	}
}

type refreshToken struct {
	UserRepo     core.UserRepository
	Log          core.AppLogger
	Secret       string
	TokenService core.TokenService
}

func (this *refreshToken) Do(ctx core.ReqContext, authToken, refreshToken, fingerprint, useragent string) (aToken string, rToken string, err error) {

	if ok, err := this.validate(authToken, refreshToken, fingerprint, useragent); !ok {
		this.Log.Infow("Not valid data",
			"reqId", ctx.ReqId(),
			"authToken", authToken,
			"refreshToken", refreshToken,
			"fingerprint", fingerprint,
			"useragent", useragent,
			"error", err)
		return "", "", core.NewError(core.InvalidRequest)
	}
	useragent = core.UserAgentFingerprint(useragent)

	auth, refresh, err := this.TokenService.Refresh(authToken, refreshToken, fingerprint, useragent)

	if err != nil {
		this.Log.Infow("Fail to refresh token",
			"reqId", ctx.ReqId(),
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
