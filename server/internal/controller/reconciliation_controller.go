// Package controller 提供HTTP请求处理控制器
// ReconciliationController 对账控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

type ReconciliationController struct {
	service *service.ReconciliationService
}

func NewReconciliationController(service *service.ReconciliationService) *ReconciliationController {
	return &ReconciliationController{service: service}
}

func (ctrl *ReconciliationController) Create(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.CreateReconciliationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.service.Create(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *ReconciliationController) Get(c *gin.Context) {
	id := parseIDParam(c, "id")

	result, err := ctrl.service.Get(id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *ReconciliationController) List(c *gin.Context) {
	var req request.ReconciliationListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.service.List(&req)
	if err != nil {
		handleError(c, err)
		return
	}

	SuccessPage(c, result)
}

func (ctrl *ReconciliationController) Update(c *gin.Context) {
	id := parseIDParam(c, "id")

	var req request.UpdateReconciliationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.service.Update(id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *ReconciliationController) Delete(c *gin.Context) {
	id := parseIDParam(c, "id")

	if err := ctrl.service.Delete(id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}