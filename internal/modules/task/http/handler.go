package http

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/tamboto2000/otaqku-tasks/internal/common"
	"github.com/tamboto2000/otaqku-tasks/internal/dto"
	"github.com/tamboto2000/otaqku-tasks/internal/modules/task"
)

type TaskHandler struct {
	taskSvc  task.TaskService
	logger   *slog.Logger
	authMddl echo.MiddlewareFunc
}

func NewTaskHandler(taskSvc task.TaskService, logger *slog.Logger, authMddl echo.MiddlewareFunc) TaskHandler {
	return TaskHandler{taskSvc: taskSvc, logger: logger, authMddl: authMddl}
}

func RegisterTaskHandler(h TaskHandler, router *echo.Echo) {
	group := router.Group("tasks", h.authMddl)
	group.POST("", h.CreateTask)
}

func (h TaskHandler) CreateTask(ectx echo.Context) error {
	ctx := ectx.Request().Context()

	var req dto.CreateTaskRequest
	if err := ectx.Bind(&req); err != nil {
		return common.InvalidReqBodyResponse(ectx, err)
	}

	accId, err := common.AccountIDFromEchoCtx(ectx)
	if err != nil {
		return common.InternalServerErrorResponse(ectx, h.logger, err)
	}

	if err := h.taskSvc.CreateTask(ctx, accId, req); err != nil {
		return common.ErrorResponse(ectx, err)
	}

	return common.OKResponse(ectx, "success", nil)
}
