package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

type PiggyBankController struct {
	service *service.PiggyBankService
}

func NewPiggyBankController(service *service.PiggyBankService) *PiggyBankController {
	return &PiggyBankController{service: service}
}

func (c *PiggyBankController) Create(ctx *gin.Context) {
	userID := ctx.GetUint64("userID")
	var req request.CreatePiggyBankReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := c.service.Create(userID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, resp)
}

func (c *PiggyBankController) Get(ctx *gin.Context) {
	userID := ctx.GetUint64("userID")
	id := ctx.GetUint64("id")
	resp, err := c.service.Get(userID, id)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (c *PiggyBankController) List(ctx *gin.Context) {
	userID := ctx.GetUint64("userID")
	resp, err := c.service.List(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (c *PiggyBankController) Update(ctx *gin.Context) {
	userID := ctx.GetUint64("userID")
	id := ctx.GetUint64("id")
	var req request.UpdatePiggyBankReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := c.service.Update(userID, id, &req)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (c *PiggyBankController) Delete(ctx *gin.Context) {
	userID := ctx.GetUint64("userID")
	id := ctx.GetUint64("id")
	if err := c.service.Delete(userID, id); err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

func (c *PiggyBankController) AddAmount(ctx *gin.Context) {
	userID := ctx.GetUint64("userID")
	id := ctx.GetUint64("id")
	var req request.AddAmountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := c.service.AddAmount(userID, id, &req)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (c *PiggyBankController) RemoveAmount(ctx *gin.Context) {
	userID := ctx.GetUint64("userID")
	id := ctx.GetUint64("id")
	var req request.RemoveAmountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := c.service.RemoveAmount(userID, id, &req)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (c *PiggyBankController) GetEvents(ctx *gin.Context) {
	userID := ctx.GetUint64("userID")
	id := ctx.GetUint64("id")
	resp, err := c.service.GetEvents(userID, id)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (c *PiggyBankController) Reorder(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	var req struct {
		Orders map[uint64]int `json:"orders" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := c.service.Reorder(userID, req.Orders); err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, nil)
}

func (c *PiggyBankController) ResetHistory(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	id := parseIDParam(ctx, "id")

	resp, err := c.service.ResetHistory(userID, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	Success(ctx, resp)
}
