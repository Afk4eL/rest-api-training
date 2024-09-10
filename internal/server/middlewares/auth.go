package middlewares

import (
	httperrors "clean-rest-arch/internal/server/utils/http/http.errors"
	"clean-rest-arch/internal/server/utils/jwt"
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
)

// type contextKey string

// const (
// 	UserIdContextKey contextKey = "id"
// )

func MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie(jwt.JwtCookieName)
		if err != nil {
			httperrors.SetErrResponse(w, r, err)

			return
		}

		claims, err := jwt.ParseJWT(cookie)
		if err != nil {
			if err == jwt.ErrIncorrectCookieName {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, httperrors.ErrNoCookie)

				return
			}
			httperrors.SetErrResponse(w, r, err)

			return
		}

		idToCtx, _ := strconv.Atoi(claims.Subject)
		ctx := context.WithValue(r.Context(), "userId", uint(idToCtx))

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
