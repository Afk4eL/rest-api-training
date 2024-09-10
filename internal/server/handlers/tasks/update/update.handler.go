package task_update

import (
	"clean-rest-arch/internal/models"
	utilhttp "clean-rest-arch/internal/server/utils/http"
	httperrors "clean-rest-arch/internal/server/utils/http/http.errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Task models.TaskEntity
}

type Response struct {
	Status string `json:"status"`
}

type TaskUpdater interface {
	UpdateTask(task *models.TaskEntity) error
}

func UpdateTask(log *slog.Logger, taskUpdater TaskUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.tasks.UpdateTask"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var data Request

		if err := utilhttp.ReadRequest(r.Body, &data); err != nil {
			log.Error("Failed to read request", "Error", err.Error())
			httperrors.SetErrResponse(w, r, err)

			return
		}

		data.Task.Id = r.Context().Value("userId").(uint)
		err := taskUpdater.UpdateTask(&data.Task)
		if err != nil {
			log.Error("Failed to update task", "Error", err.Error())
			httperrors.SetErrResponse(w, r, err)

			return
		}

		log.Info("Task updated", "Id", data.Task.Id)

		render.JSON(w, r, Response{
			Status: "Task updated",
		})
	}
}
