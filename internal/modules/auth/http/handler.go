package http

import (
	"github.com/labstack/echo/v4"
	"github.com/tamboto2000/otaqku-tasks/internal/common"
	"github.com/tamboto2000/otaqku-tasks/internal/dto"
	"github.com/tamboto2000/otaqku-tasks/internal/modules/auth"
)

type AuthHandler struct {
	authSvc auth.AuthService
}

func NewAuthHandler(authSvc auth.AuthService) AuthHandler {
	return AuthHandler{authSvc: authSvc}
}

func RegisterAuthHandler(h AuthHandler, router *echo.Echo) {
	group := router.Group("/auth")
	group.POST("/account", h.RegisterAccount)
	group.POST("/login", h.Login)
	group.POST("/refresh_token", h.ExchangeRefreshToken)
}

func (h AuthHandler) RegisterAccount(ectx echo.Context) error {
	ctx := ectx.Request().Context()

	var req dto.CreateAccountRequest
	if err := ectx.Bind(&req); err != nil {
		return common.InvalidReqBodyResponse(ectx, err)
	}

	if err := h.authSvc.RegisterAccount(ctx, req); err != nil {
		return common.ErrorResponse(ectx, err)
	}

	return common.OKResponse(ectx, "success", nil)
}

func (h AuthHandler) Login(ectx echo.Context) error {
	ctx := ectx.Request().Context()

	var req dto.LoginRequest
	if err := ectx.Bind(&req); err != nil {
		return common.InvalidReqBodyResponse(ectx, err)
	}

	tokens, err := h.authSvc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return common.ErrorResponse(ectx, err)
	}

	return common.OKResponse(ectx, "success", tokens)
}

func (h AuthHandler) ExchangeRefreshToken(ectx echo.Context) error {
	ctx := ectx.Request().Context()

	var req dto.ExchangeRefreshTokenRequest
	if err := ectx.Bind(&req); err != nil {
		return common.InvalidReqBodyResponse(ectx, err)
	}

	tokens, err := h.authSvc.ExchangeRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return common.ErrorResponse(ectx, err)
	}

	return common.OKResponse(ectx, "success", tokens)
}
