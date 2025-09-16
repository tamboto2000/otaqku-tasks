package http

import (
	"errors"
	"log/slog"
	"strconv"

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
	group.GET("", h.GetTaskList)
	group.GET("/:id", h.GetByID)
	group.PUT("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
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

func (h TaskHandler) GetTaskList(ectx echo.Context) error {
	ctx := ectx.Request().Context()
	var req dto.Pagination
	if err := ectx.Bind(&req); err != nil {
		return common.InvalidQueryParamResponse(ectx, err)
	}

	accId, err := common.AccountIDFromEchoCtx(ectx)
	if err != nil {
		return common.InternalServerErrorResponse(ectx, h.logger, err)
	}

	taskList, err := h.taskSvc.GetTaskList(ctx, accId, req)
	if err != nil {
		return common.ErrorResponse(ectx, err)
	}

	return common.OKResponse(ectx, "success", taskList)
}

func (h TaskHandler) GetByID(ectx echo.Context) error {
	ctx := ectx.Request().Context()

	accId, err := common.AccountIDFromEchoCtx(ectx)
	if err != nil {
		return common.InternalServerErrorResponse(ectx, h.logger, err)
	}

	idStr := ectx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return common.InvalidQueryParamResponse(ectx, errors.New("invalid task number"))
	}

	task, err := h.taskSvc.GetByID(ctx, accId, id)
	if err != nil {
		return common.ErrorResponse(ectx, err)
	}

	return common.OKResponse(ectx, "success", task)
}

func (h TaskHandler) Delete(ectx echo.Context) error {
	ctx := ectx.Request().Context()

	accId, err := common.AccountIDFromEchoCtx(ectx)
	if err != nil {
		return common.InternalServerErrorResponse(ectx, h.logger, err)
	}

	idStr := ectx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return common.InvalidQueryParamResponse(ectx, errors.New("invalid task number"))
	}

	if err := h.taskSvc.Delete(ctx, accId, id); err != nil {
		return common.ErrorResponse(ectx, err)
	}

	return common.OKResponse(ectx, "success", nil)
}

func (h TaskHandler) Update(ectx echo.Context) error {
	ctx := ectx.Request().Context()

	var req dto.Task
	if err := ectx.Bind(&req); err != nil {
		return common.InvalidReqBodyResponse(ectx, err)
	}

	accId, err := common.AccountIDFromEchoCtx(ectx)
	if err != nil {
		return common.InternalServerErrorResponse(ectx, h.logger, err)
	}

	idStr := ectx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return common.InvalidQueryParamResponse(ectx, errors.New("invalid task number"))
	}

	req.ID = id

	if err := h.taskSvc.Update(ctx, accId, req); err != nil {
		return common.ErrorResponse(ectx, err)
	}

	return common.OKResponse(ectx, "success", nil)
}
