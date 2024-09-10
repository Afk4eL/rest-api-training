package task_get

import (
	"clean-rest-arch/internal/models"
	utilhttp "clean-rest-arch/internal/server/utils/http"
	httperrors "clean-rest-arch/internal/server/utils/http/http.errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	Task models.TaskEntity
}

type Request struct {
	TaskId uint `json:"id" validate:"required"`
}
type TaskGetter interface {
	GetTask(userId uint, taskId uint) (*models.TaskEntity, error)
}

func GetTask(log *slog.Logger, taskGetter TaskGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.tasks.GetTask"

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
		task, err := taskGetter.GetTask(userId, data.TaskId)
		if err != nil {
			log.Error("Failed to get task", "Error", err.Error())
			httperrors.SetErrResponse(w, r, err)

			return
		}

		log.Info("Task getted", "Id", task)

		render.JSON(w, r, Response{
			Task: *task,
		})
	}
}
