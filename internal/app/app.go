package app

import (
	"clean-rest-arch/internal/config"
	"clean-rest-arch/internal/server/router"
	"clean-rest-arch/internal/server/utils/slogpretty"
	"clean-rest-arch/storage/postgres"
	"clean-rest-arch/storage/repos"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type App struct {
	cfg      *config.Config
	logger   *slog.Logger
	storage  *postgres.Database
	userRepo repos.UserRepository
	taskRepo repos.TaskRepository
	router   *chi.Mux
	server   *http.Server
}

func (app *App) readConfig() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Usage go run <path to main.go> [arguments] \n Required arguments: \n - Path to config file")
		os.Exit(1)
	}

	app.cfg = config.MustLoad(args[0])
}

func New() *App {
	return &App{}
}

func (app *App) Config() {
	app.readConfig()

	app.logger = SetupLogger(app.cfg.Env)
	app.logger.Debug("Debug messages are enabled")
	app.logger.Info("Starting server")

	storage, err := postgres.New(*app.cfg)
	if err != nil {
		app.logger.Error("Failed to init storage", "Error", err.Error())
		os.Exit(1)
	}
	app.storage = storage

	app.userRepo = repos.NewUserRepository(app.storage.Database)
	app.taskRepo = repos.NewTaskRepository(app.storage.Database)

	app.router = router.NewRouter(app.logger, app.userRepo, app.taskRepo)

	app.logger.Info("Starting server", slog.String("Address", app.cfg.Address))
	app.server = &http.Server{
		Addr:        app.cfg.Address,
		Handler:     app.router,
		ReadTimeout: app.cfg.Timeout,
		IdleTimeout: app.cfg.IdleTimeout,
	}
}

func (app *App) Run() {
	const op = "app.Run"

	if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		app.logger.Error("Fatal error", op, err.Error())
		return
	}
}

func (app *App) Stop() {
	const op = "app.Stop"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.storage.Stop()

	if err := app.server.Shutdown(ctx); err != nil {
		app.logger.Error("Server forced to shutdown", op, err.Error())
		return
	}

	app.logger.Info("Server stopped")
}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
