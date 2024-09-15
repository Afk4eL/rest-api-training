package signinHandler

import (
	"errors"
	"log/slog"
	"net/http"
	"rest-arch-training/internal/models"
	"rest-arch-training/internal/server/utils/csrf"
	utilhttp "rest-arch-training/internal/server/utils/http"
	httperrors "rest-arch-training/internal/server/utils/http/http.errors"
	"rest-arch-training/internal/server/utils/jwt"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Request struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	Id uint `json:"id"`
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=UserSaver
type FindUser interface {
	GetUserByUsername(username string) (*models.UserEntity, error)
}

func Signin(log *slog.Logger, userFind FindUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.Signin"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var data Request

		err := utilhttp.ReadRequest(r.Body, &data)
		if err != nil {
			log.Error("Failed to read request body", "Error", err.Error())
			httperrors.SetErrResponse(w, r, err)

			return
		}

		user, err := userFind.GetUserByUsername(data.Username)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Info("User not found", "Error", err.Error())

				httperrors.SetErrResponse(w, r, errors.Join(err, httperrors.ErrUserNotFound))

				return
			}

			log.Error("Failed to get username", "Error", err.Error())

			httperrors.SetErrResponse(w, r, err)

			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
			log.Info("Wrong password", "Info", err.Error())
			httperrors.SetErrResponse(w, r, err)

			return
		}

		log.Info("User signin", "User", user)

		token, err := jwt.GenerateJWT(int(user.Id))
		if err != nil {
			log.Error("Failed to generate JWT", "Error", err.Error())

			httperrors.SetErrResponse(w, r, err)

			return
		}

		csrfToken, err := csrf.MakeToken()
		if err != nil {
			log.Error("Failed to generate CSRF token", "Error", err.Error())

			httperrors.SetErrResponse(w, r, err)

			return
		}

		cookie := &http.Cookie{
			Name:     jwt.JwtCookieName,
			Value:    token,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, cookie)

		cookie = &http.Cookie{
			Name:     csrf.CSRFHeader,
			Value:    csrfToken,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, cookie)

		w.Header().Add(csrf.CSRFHeader, csrfToken)

		render.JSON(w, r, Response{
			Id: user.Id,
		})
	}
}
