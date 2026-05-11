package handler

import (
	"planet/internal/dto"
	"planet/internal/pkg"
	"planet/internal/service"

	"github.com/gin-gonic/gin"
)

type SearchHandler interface {
	SearchUsers(c *gin.Context)
}

type searchHandler struct {
	searchSvc service.SearchService
}

func NewSearchHandler(searchSvc service.SearchService) SearchHandler {
	return &searchHandler{
		searchSvc: searchSvc,
	}
}

func (h *searchHandler) SearchUsers(c *gin.Context) {
	var req dto.SearchUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pkg.Fail(c, 400, err.Error())
		return
	}

	req.RequesterUserId = c.GetString("userID")
	users, err := h.searchSvc.SearchUsers(&req)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}

	pkg.Success(c, 200, users)
}
