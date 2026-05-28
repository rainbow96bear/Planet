package handler

import (
	"planet/internal/dto"
	"planet/internal/pkg"
	"planet/internal/service"

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

	req.UserID = c.GetString("userID")

	task, err := h.taskSvc.CreateTask(&req)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}
	pkg.Success(c, 201, task)
}

func (h *taskHandler) DeleteTask(c *gin.Context) {
	taskID := c.Param("task_id")

	req := dto.DeleteTaskRequest{
		ID:     taskID,
		UserID: c.GetString("userID"),
	}

	if err := h.taskSvc.DeleteTask(&req); err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}
	pkg.Success(c, 204, nil)
}

func (h *taskHandler) GetTasksByMonth(c *gin.Context) {
	userID := c.Param("userid")

	var req dto.GetTasksByMonthRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pkg.Fail(c, 400, err.Error())
		return
	}

	req.UserID = userID
	req.RequesterUserId = c.GetString("userID")

	tasks, err := h.taskSvc.GetTasksByMonth(&req)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}

	pkg.Success(c, 200, tasks)
}

func (h *taskHandler) ToggleTask(c *gin.Context) {
	taskID := c.Param("task_id")

	req := dto.ToggleTaskRequest{
		ID: taskID,
	}

	task, err := h.taskSvc.ToggleTask(&req)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}

	pkg.Success(c, 201, task)
}
