package handler

import (
	"planet/internal/service"
)

type TaskHandler interface {
}

type taskHandler struct {
	taskSvc service.TaskService
}

func NewTaskHandler(taskSvc service.TaskService) TaskHandler {
	return &taskHandler{taskSvc: taskSvc}
}
