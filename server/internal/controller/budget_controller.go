package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zero-life/server/internal/dto/request"
	"github.com/zero-life/server/internal/pkg/errcode"
	"github.com/zero-life/server/internal/service"
)

type BudgetController struct {
	budgetService *service.BudgetService
}

func NewBudgetController(budgetService *service.BudgetService) *BudgetController {
	return &BudgetController{budgetService: budgetService}
}

func (ctrl *BudgetController) Create(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.CreateBudgetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.budgetService.Create(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *BudgetController) Get(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	result, err := ctrl.budgetService.Get(userID, id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *BudgetController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.budgetService.List(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *BudgetController) Update(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	var req request.UpdateBudgetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.budgetService.Update(userID, id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *BudgetController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.budgetService.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

func (ctrl *BudgetController) GetHistory(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	result, err := ctrl.budgetService.GetHistory(userID, id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}
