package handler

import (
	"planet/internal/dto"
	"planet/internal/pkg"
	"planet/internal/service"

	"github.com/gin-gonic/gin"
)

type TaskHandler interface {
	CreateTask(c *gin.Context)
}

type taskHandler struct {
	taskSvc service.TaskService
}

func NewTaskHandler(taskSvc service.TaskService) TaskHandler {
	return &taskHandler{taskSvc: taskSvc}
}

func (h *taskHandler) CreateTask(c *gin.Context) {
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, 400, err.Error())
		return
	}

	task, err := h.taskSvc.CreateTask(&req)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}
	pkg.Success(c, 201, task)
}
