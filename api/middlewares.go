package api

import (
	"context"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/NeekUP/roadmaps/infrastructure"
	"net/http"
)

func Auth(rights domain.Rights, ts core.TokenService, log core.AppLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			ctx := r.Context()
			if len(authHeader) > 0 {
				userId, userName, userRights, err := ts.Validate(authHeader[7:])
				if err != nil {
					log.Errorw("Unauthorized. Error", "path", r.URL.Path, "requiredRights", rights, "error", err.Error())
					statusResponse(w, &status{Code: http.StatusUnauthorized})
					return
				}

				if userId == "" {
					log.Errorw("Unauthorized", "path", r.URL.Path, "requiredRights", rights)
					statusResponse(w, &status{Code: http.StatusUnauthorized})
					return
				}

				if rights != domain.All && !domain.Rights(userRights).HasFlag(rights) {
					log.Infow("Forbidden", "path", r.URL.Path, "requiredRights", rights, "userId", userId, "userName", userName, "userRights", userRights)
					statusResponse(w, &status{Code: http.StatusForbidden})
					return
				}

				log.Infow("Authorized", "path", r.URL.Path, "requiredRights", rights, "userId", userId, "userName", userName, "userRights", userRights)

				ctx = context.WithValue(ctx, infrastructure.ReqRights, userRights)
				ctx = context.WithValue(ctx, infrastructure.ReqUserId, userId)
				ctx = context.WithValue(ctx, infrastructure.ReqUserName, userName)

			} else if rights != domain.All {
				log.Infow("Unauthorized. no auth", "path", r.URL.Path, "requiredRights", rights)
				statusResponse(w, &status{Code: http.StatusUnauthorized})
				return
			} else {
				log.Infow("NoAuth", "path", r.URL.Path, "requiredRights", rights)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
