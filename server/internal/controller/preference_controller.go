package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

type PreferenceController struct {
	prefService *service.PreferenceService
}

func NewPreferenceController(prefService *service.PreferenceService) *PreferenceController {
	return &PreferenceController{prefService: prefService}
}

func (ctrl *PreferenceController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.prefService.List(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *PreferenceController) Get(c *gin.Context) {
	userID := c.GetUint64("user_id")
	name := c.Param("name")

	result, err := ctrl.prefService.Get(userID, name)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *PreferenceController) Update(c *gin.Context) {
	userID := c.GetUint64("user_id")
	name := c.Param("name")

	var req request.UpdatePreferenceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.prefService.Set(userID, name, req.Value)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}
