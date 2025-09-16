package app

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/tamboto2000/otaqku-tasks/internal/config"
	"github.com/tamboto2000/otaqku-tasks/internal/modules/auth"
	authHttp "github.com/tamboto2000/otaqku-tasks/internal/modules/auth/http"
	"github.com/vinovest/sqlx"
)

type repositories struct {
	accRepo auth.AccountRepository
}

func newRepositories(db *sqlx.DB, logger *slog.Logger) repositories {
	return repositories{
		accRepo: auth.NewPostgreAccountRepository(db, logger),
	}
}

type services struct {
	authSvc auth.AuthService
}

func newServices(cfg config.Config, repos repositories, logger *slog.Logger) services {
	return services{
		authSvc: auth.NewAuthService(cfg.JWT, repos.accRepo, logger),
	}
}

func registerHandlers(router *echo.Echo, svcs services) {
	authHandler := authHttp.NewAuthHandler(svcs.authSvc)
	authHttp.RegisterAuthHandler(authHandler, router)
}
