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
	GetOrbitSchedulesByMonth(c *gin.Context)
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

	if tasks == nil {
		tasks = []*dto.GetTasksByMonthResponse{} // null 대신 빈 배열
	}

	pkg.Success(c, 200, tasks)
}

func (h *taskHandler) GetOrbitSchedulesByMonth(c *gin.Context) {
	userID := c.Param("userid")
	requesterID := c.GetString("userID")

	// Orbit Schedule은 "내가 Orbit한 사람들의 일정"이라, 본인 것만 조회 가능하다.
	// 다른 사람의 Orbit 그래프를 들여다보는 API는 존재하지 않는다 (문서 8번 정책).
	if userID != requesterID {
		pkg.Fail(c, 403, "본인의 Orbit Schedule만 조회할 수 있습니다")
		return
	}

	var req dto.GetOrbitSchedulesByMonthRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pkg.Fail(c, 400, err.Error())
		return
	}
	req.OrbiterID = requesterID

	schedules, err := h.taskSvc.GetOrbitSchedulesByMonth(&req)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}

	if schedules == nil {
		schedules = []*dto.OrbitScheduleResponse{}
	}

	pkg.Success(c, 200, schedules)
}

func (h *taskHandler) ToggleTask(c *gin.Context) {
	taskID := c.Param("task_id")

	req := dto.ToggleTaskRequest{
		ID:     taskID,
		UserID: c.GetString("userID"),
	}

	task, err := h.taskSvc.ToggleTask(&req)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}

	pkg.Success(c, 201, task)
}
