package signupHandler

import (
	"log/slog"
	"net/http"
	"rest-arch-training/internal/models"
	utilhttp "rest-arch-training/internal/server/utils/http"
	httperrors "rest-arch-training/internal/server/utils/http/http.errors"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

type Request struct {
	User models.UserEntity `json:"user"`
}

type Response struct {
	Id uint `json:"id"`
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=UserSaver
type UserSaver interface {
	CreateUser(user *models.UserEntity) (uint, error)
}

func SignUp(log *slog.Logger, userSaver UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.Signup"

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

		password, err := bcrypt.GenerateFromPassword([]byte(data.User.Password), 10)
		if err != nil {
			log.Error("Password hash generate failed", "Error", err.Error())

			httperrors.SetErrResponse(w, r, err)

			return
		}

		id, err := userSaver.CreateUser(&models.UserEntity{Username: data.User.Username, Email: data.User.Email, Password: string(password)})
		if err != nil {
			log.Error("Failed to save user", "Error", err.Error())

			httperrors.SetErrResponse(w, r, err)

			return
		}

		log.Info("User saved", "Id", id)

		render.JSON(w, r, Response{
			Id: id,
		})
	}
}
