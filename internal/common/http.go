package common

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HTTPResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func InvalidReqBodyResponse(ectx echo.Context, err error) error {
	resp := HTTPResponse{
		Message: "Invalid request body",
		Error:   Error{Message: err.Error()},
	}

	return ectx.JSON(http.StatusBadRequest, resp)
}

func InvalidQueryParamResponse(ectx echo.Context, err error) error {
	resp := HTTPResponse{
		Message: "Invalid filter",
		Error:   Error{Message: err.Error()},
	}

	return ectx.JSON(http.StatusBadRequest, resp)
}

func ErrorResponse(ectx echo.Context, err error) error {
	httpCode := http.StatusInternalServerError
	resp := HTTPResponse{Message: "Internal server error"}

	var xErr Error
	if errors.As(err, &xErr) {
		switch xErr.Code {
		case ErrCodeInputValidation:
			httpCode = http.StatusBadRequest

		case ErrCodeNotFound:
			httpCode = http.StatusNotFound

		case ErrCodeAlreadyExists:
			httpCode = http.StatusConflict

		case ErrCodeUnauthorized:
			httpCode = http.StatusUnauthorized
		}

		resp.Message = xErr.Message
		resp.Error = xErr

		return ectx.JSON(httpCode, resp)
	}

	return ectx.JSON(http.StatusInternalServerError, resp)
}

func OKResponse(ectx echo.Context, msg string, data any) error {
	resp := HTTPResponse{
		Message: msg,
		Data:    data,
	}

	return ectx.JSON(http.StatusOK, resp)
}

func InternalServerErrorResponse(ectx echo.Context, logger *slog.Logger, err error) error {
	resp := HTTPResponse{
		Message: "Internal server error",
	}

	logger.Error(fmt.Sprintf("internal server error: %v", err))

	return ectx.JSON(http.StatusInternalServerError, resp)
}
