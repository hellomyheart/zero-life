package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

type ReconciliationController struct {
	reconService *service.ReconciliationService
}

func NewReconciliationController(reconService *service.ReconciliationService) *ReconciliationController {
	return &ReconciliationController{reconService: reconService}
}

func (ctrl *ReconciliationController) GetReconciliation(c *gin.Context) {
	userID := c.GetUint64("user_id")
	accountID := parseIDParam(c, "id")

	var req request.GetReconciliationReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.reconService.GetReconciliationData(userID, accountID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *ReconciliationController) SubmitReconciliation(c *gin.Context) {
	userID := c.GetUint64("user_id")
	accountID := parseIDParam(c, "id")

	var req request.SubmitReconciliationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.reconService.SubmitReconciliation(userID, accountID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}
