package api

import (
	"context"
	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/NeekUP/roadmaps/infrastructure"
	"net/http"
)

func Auth(rights domain.Rights, ts core.TokenService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			ctx := r.Context()
			if len(authHeader) > 0 {
				userId, userName, userRights, err := ts.Validate(authHeader[7:])
				if userId == "" || err != nil {
					statusResponse(w, &status{Code: http.StatusUnauthorized})
					return
				}

				if !rights.HasFlag(domain.God) && !domain.Rights(userRights).HasFlag(rights) {
					statusResponse(w, &status{Code: http.StatusForbidden})
					return
				}

				ctx = context.WithValue(ctx, infrastructure.ReqRights, userRights)
				ctx = context.WithValue(ctx, infrastructure.ReqUserId, userId)
				ctx = context.WithValue(ctx, infrastructure.ReqUserName, userName)

			} else if !rights.HasFlag(domain.God) {
				statusResponse(w, &status{Code: http.StatusUnauthorized})
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
