package router

import (
	"log/slog"
	task_create "rest-arch-training/internal/server/handlers/tasks/create"
	task_delete "rest-arch-training/internal/server/handlers/tasks/delete"
	task_get "rest-arch-training/internal/server/handlers/tasks/get"
	task_get_all "rest-arch-training/internal/server/handlers/tasks/get-all"
	task_update "rest-arch-training/internal/server/handlers/tasks/update"
	getUserHandler "rest-arch-training/internal/server/handlers/users/get"
	"rest-arch-training/internal/server/handlers/users/logout"
	signinHandler "rest-arch-training/internal/server/handlers/users/signin"
	signupHandler "rest-arch-training/internal/server/handlers/users/signup"
	"rest-arch-training/internal/server/middlewares"
	"rest-arch-training/storage/repos"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(log *slog.Logger, userRepo repos.UserRepository, taskRepo repos.TaskRepository) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middlewares.MiddlewareMetrics)

	router.Handle("/metrics", promhttp.Handler())

	router.Route("/auth", func(r chi.Router) {
		r.Use(middlewares.CSRFValidate)
		r.Use(middlewares.MiddlewareAuth)

		r.Get("/get-user-id", getUserHandler.GetUserById(log, userRepo))
		r.Post("/create-task", task_create.CreateTask(log, taskRepo))
		r.Delete("/delete-task", task_delete.DeleteTask(log, taskRepo))
		r.Get("/get-task-id", task_get.GetTask(log, taskRepo))
		r.Get("/get-all-tasks", task_get_all.GetAllUserTasks(log, taskRepo))
		r.Patch("/update-task", task_update.UpdateTask(log, taskRepo))
	})

	router.Post("/signup", signupHandler.SignUp(log, userRepo))
	router.Post("/signin", signinHandler.Signin(log, userRepo))
	router.Post("/logout", logout.Logout(log))

	return router
}
