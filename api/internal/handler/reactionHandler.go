package handler

import (
	"planet/internal/dto"
	"planet/internal/pkg"
	"planet/internal/service"

	"github.com/gin-gonic/gin"
)

type ReactionHandler interface {
	AddReaction(c *gin.Context)
	RemoveReaction(c *gin.Context)
}

type reactionHandler struct {
	reactionSvc service.ReactionService
}

func NewReactionHandler(reactionSvc service.ReactionService) ReactionHandler {
	return &reactionHandler{
		reactionSvc: reactionSvc,
	}
}

func (h *reactionHandler) AddReaction(c *gin.Context) {
	var req dto.AddReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, 400, err.Error())
		return
	}

	req.TaskID = c.Param("task_id")
	req.UserID = c.GetString("userID")

	reaction, err := h.reactionSvc.AddReaction(&req)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}

	pkg.Success(c, 201, reaction)
}

func (h *reactionHandler) RemoveReaction(c *gin.Context) {
	req := dto.RemoveReactionRequest{
		TaskID: c.Param("task_id"),
		UserID: c.GetString("userID"),
		Type:   c.Param("type"),
	}

	if err := h.reactionSvc.RemoveReaction(&req); err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}

	pkg.Success(c, 204, nil)
}
