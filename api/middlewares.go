package api

import (
	"context"
	"net/http"
	"roadmaps/core"
	"roadmaps/domain"
	"roadmaps/infrastructure"
)

func Auth(rights domain.Rights, ts core.TokenService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if len(authHeader) > 0 {
				userId, userName, userRights, err := ts.Validate(authHeader[7:])
				if userId == "" || err != nil {
					statusResponse(w, &status{Code: http.StatusUnauthorized})
					return
				}

				if !rights.HasFlag(domain.All) && !domain.Rights(userRights).HasFlag(rights) {
					statusResponse(w, &status{Code: http.StatusForbidden})
					return
				}

				ctx := r.Context()
				ctx = context.WithValue(ctx, infrastructure.ReqRights, userRights)
				ctx = context.WithValue(ctx, infrastructure.ReqUserId, userId)
				ctx = context.WithValue(ctx, infrastructure.ReqUserName, userName)

			} else if !rights.HasFlag(domain.All) {
				statusResponse(w, &status{Code: http.StatusUnauthorized})
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

