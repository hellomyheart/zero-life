// Package controller 提供HTTP请求处理控制器
// RecurringTransactionController 循环交易控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

type RecurringTransactionController struct {
	service *service.RecurringTransactionService
}

func NewRecurringTransactionController(service *service.RecurringTransactionService) *RecurringTransactionController {
	return &RecurringTransactionController{service: service}
}

func (ctrl *RecurringTransactionController) Create(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.CreateRecurringTransactionReq
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

func (ctrl *RecurringTransactionController) Get(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	result, err := ctrl.service.Get(userID, id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *RecurringTransactionController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.RecurringTransactionListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.service.List(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	SuccessPage(c, result)
}

func (ctrl *RecurringTransactionController) Update(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	var req request.UpdateRecurringTransactionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.service.Update(userID, id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *RecurringTransactionController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.service.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

func (ctrl *RecurringTransactionController) ProcessDue(c *gin.Context) {
	userID := c.GetUint64("user_id")

	created, err := ctrl.service.ProcessDue(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, gin.H{"created": created})
}