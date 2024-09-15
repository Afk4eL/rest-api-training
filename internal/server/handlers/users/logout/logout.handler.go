package logout

import (
	"log/slog"
	"net/http"
	"rest-arch-training/internal/server/utils/csrf"
	"rest-arch-training/internal/server/utils/jwt"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func Logout(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.Logout"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Logout", "User", r.Context().Value("user_id"))

		cookie := &http.Cookie{
			Name:     jwt.JwtCookieName,
			Value:    "",
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)

		cookie = &http.Cookie{
			Name:     csrf.CSRFHeader,
			Value:    "",
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)

		render.JSON(w, r, nil)
	}
}
