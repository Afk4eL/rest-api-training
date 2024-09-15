package task_delete

import (
	"log/slog"
	"net/http"
	utilhttp "rest-arch-training/internal/server/utils/http"
	httperrors "rest-arch-training/internal/server/utils/http/http.errors"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"gorm.io/gorm"
)

type Request struct {
	Id uint `json:"id" validate:"required"`
}

type Response struct {
	Id uint `json:"id"`
}

type TaskDeleter interface {
	DeleteTask(userId uint, taskId uint) error
}

func DeleteTask(log *slog.Logger, taskDeleter TaskDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.tasks.DeleteTask"

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

		userId := r.Context().Value("userId").(uint)
		err = taskDeleter.DeleteTask(userId, data.Id)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Error("Task not found", "Error", err.Error())
				httperrors.SetErrResponse(w, r, err)

				return
			}
			log.Error("Failed to delete task", "Error", err.Error())

			httperrors.SetErrResponse(w, r, err)

			return
		}

		log.Info("Task deleted", "Id", data.Id)

		render.JSON(w, r, Response{
			Id: data.Id,
		})
	}
}
