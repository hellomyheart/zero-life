package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

type AdminController struct {
	adminService *service.AdminService
}

func NewAdminController(adminService *service.AdminService) *AdminController {
	return &AdminController{adminService: adminService}
}

func (ctrl *AdminController) GetConfiguration(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.adminService.GetConfiguration(name)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *AdminController) ListConfigurations(c *gin.Context) {
	result, err := ctrl.adminService.ListConfigurations()
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *AdminController) UpdateConfiguration(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	var req request.AdminUpdateConfigurationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.adminService.UpdateConfiguration(name, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *AdminController) TestEmail(c *gin.Context) {
	var req request.AdminTestEmailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.adminService.TestEmail(req.Email)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}
