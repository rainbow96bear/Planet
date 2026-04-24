package handler

import (
	"planet/internal/dto"
	"planet/internal/pkg"
	"planet/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskHandler interface {
	CreateTask(c *gin.Context)
	DeleteTask(c *gin.Context)
}

type taskHandler struct {
	taskSvc service.TaskService
}

func NewTaskHandler(db *gorm.DB, taskSvc service.TaskService) TaskHandler {
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
		TaskID: uint(taskID),
		UserID: c.GetUint("userID"),
	}

	if err := h.taskSvc.DeleteTask(&req); err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}
	pkg.Success(c, 204, nil)
}
