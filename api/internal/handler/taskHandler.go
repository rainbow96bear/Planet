package handler

import (
	"planet/internal/dto"
	"planet/internal/pkg"
	"planet/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskHandler interface {
	CreateTask(c *gin.Context)
	DeleteTask(c *gin.Context)
	GetTasksByMonth(c *gin.Context)
	ToggleTask(c *gin.Context)
}

type taskHandler struct {
	taskSvc service.TaskService
}

func NewTaskHandler(taskSvc service.TaskService) TaskHandler {
	return &taskHandler{
		taskSvc: taskSvc,
	}
}

func (h *taskHandler) CreateTask(c *gin.Context) {
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, 400, err.Error())
		return
	}

	req.UserID = c.GetUint("userID")

	task, err := h.taskSvc.CreateTask(&req)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}
	pkg.Success(c, 201, task)
}

func (h *taskHandler) DeleteTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Fail(c, 400, "잘못된 요청입니다")
		return
	}

	req := dto.DeleteTaskRequest{
		ID:     uint(taskID),
		UserID: c.GetUint("userID"),
	}

	if err := h.taskSvc.DeleteTask(&req); err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}
	pkg.Success(c, 204, nil)
}

func (h *taskHandler) GetTasksByMonth(c *gin.Context) {
	username := c.Param("username")

	yearStr := c.Query("year")
	monthStr := c.Query("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		pkg.Fail(c, 400, "잘못된 year 값입니다")
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		pkg.Fail(c, 400, "잘못된 month 값입니다")
		return
	}

	req := dto.GetTasksByMonthRequest{
		Username:          username,
		RequesterUsername: c.GetString("username"),
		Year:              year,
		Month:             month,
	}

	tasks, err := h.taskSvc.GetTasksByMonth(&req)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}

	pkg.Success(c, 200, tasks)
}

func (h *taskHandler) ToggleTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Fail(c, 400, "잘못된 요청입니다")
		return
	}

	req := dto.ToggleTaskRequest{
		ID: uint(taskID),
	}

	task, err := h.taskSvc.ToggleTask(&req)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}

	pkg.Success(c, 201, task)
}
