package task_get_all

import (
	"log/slog"
	"net/http"
	"rest-arch-training/internal/models"
	httperrors "rest-arch-training/internal/server/utils/http/http.errors"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	Tasks []*models.TaskEntity
}

type TaskGetter interface {
	GetAllUserTasks(userId uint) ([]*models.TaskEntity, error)
}

func GetAllUserTasks(log *slog.Logger, taskGetter TaskGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.tasks.GetAllUserTasks"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userId := r.Context().Value("userId").(uint)
		tasks, err := taskGetter.GetAllUserTasks(userId)
		if err != nil {
			log.Error("Failed to get tasks", "Error", err.Error())
			httperrors.SetErrResponse(w, r, err)

			return
		}

		log.Info("Tasks getted", "Id", tasks)

		render.JSON(w, r, Response{
			Tasks: tasks,
		})
	}
}
