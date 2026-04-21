package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zero-life/server/internal/dto/request"
	"github.com/zero-life/server/internal/pkg/errcode"
	"github.com/zero-life/server/internal/service"
)

type BillController struct {
	billService *service.BillService
}

func NewBillController(billService *service.BillService) *BillController {
	return &BillController{billService: billService}
}

func (ctrl *BillController) Create(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.CreateBillReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.billService.Create(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *BillController) Get(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	result, err := ctrl.billService.Get(userID, id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *BillController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.billService.List(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *BillController) Update(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	var req request.UpdateBillReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.billService.Update(userID, id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *BillController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.billService.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}
