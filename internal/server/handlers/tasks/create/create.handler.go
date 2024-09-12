package task_create

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
	Task models.TaskEntity `json:"task"`
}

type Response struct {
	Id uint `json:"id"`
}

type TaskSaver interface {
	CreateTask(task *models.TaskEntity) (uint, error)
}

func CreateTask(log *slog.Logger, taskSaver TaskSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.tasks.CreateTask"

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

		data.Task.UserId = r.Context().Value("userId").(uint)
		id, err := taskSaver.CreateTask(&data.Task)
		if err != nil {
			log.Error("Failed to create task", "Error", err.Error())
			httperrors.SetErrResponse(w, r, err)

			return
		}

		log.Info("Task created", "Id", id)

		render.JSON(w, r, Response{
			Id: id,
		})
	}
}
