package getUserHandler

import (
	"clean-rest-arch/internal/models"
	utilhttp "clean-rest-arch/internal/server/utils/http"
	httperrors "clean-rest-arch/internal/server/utils/http/http.errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"gorm.io/gorm"
)

type Request struct {
	Id uint `json:"id" validate:"required"`
}

type Response struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=UserSaver
type UserGetter interface {
	GetUserById(id uint) (*models.UserEntity, error)
}

func GetUserById(log *slog.Logger, userGetter UserGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.GetuserById"

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

		user, err := userGetter.GetUserById(data.Id)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Info("User not found")
				httperrors.SetErrResponse(w, r, err)

				return
			}

			log.Error("Failed to get user", "Error", err.Error())
			httperrors.SetErrResponse(w, r, err)

			return
		}

		log.Info("Get user", "User", user)

		render.JSON(w, r, Response{
			Username: user.Username,
			Email:    user.Email,
		})
	}
}
