package app

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/tamboto2000/otaqku-tasks/internal/config"
	"github.com/tamboto2000/otaqku-tasks/internal/modules/auth"
	authHttp "github.com/tamboto2000/otaqku-tasks/internal/modules/auth/http"
	"github.com/tamboto2000/otaqku-tasks/internal/modules/task"
	taskHttp "github.com/tamboto2000/otaqku-tasks/internal/modules/task/http"
	"github.com/vinovest/sqlx"
)

type repositories struct {
	accRepo  auth.AccountRepository
	taskRepo task.TaskRepository
}

func newRepositories(db *sqlx.DB, logger *slog.Logger) repositories {
	return repositories{
		accRepo:  auth.NewPostgreAccountRepository(db, logger),
		taskRepo: task.NewPostgreTaskRepository(db, logger),
	}
}

type services struct {
	authSvc auth.AuthService
	taskSvc task.TaskService
}

func newServices(cfg config.Config, repos repositories, logger *slog.Logger) services {
	return services{
		authSvc: auth.NewAuthService(cfg.JWT, repos.accRepo, logger),
		taskSvc: task.NewTaskService(repos.taskRepo),
	}
}

func registerHandlers(router *echo.Echo, logger *slog.Logger, svcs services) {
	authHandler := authHttp.NewAuthHandler(svcs.authSvc)
	authHttp.RegisterAuthHandler(authHandler, router)

	taskHandler := taskHttp.NewTaskHandler(svcs.taskSvc, logger, AuthMiddleware(svcs.authSvc))
	taskHttp.RegisterTaskHandler(taskHandler, router)
}
