package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

type ExportController struct {
	service *service.ExportService
}

func NewExportController(service *service.ExportService) *ExportController {
	return &ExportController{service: service}
}

func (c *ExportController) ExportTransactions(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	var req request.ExportReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	data, filename, err := c.service.ExportTransactions(userID, req.StartDate, req.EndDate, req.Format)
	if err != nil {
		handleError(ctx, err)
		return
	}

	contentType := "text/csv"
	if req.Format == "json" {
		contentType = "application/json"
	}

	ctx.Header("Content-Type", contentType)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Data(http.StatusOK, contentType, data)
}

func (c *ExportController) ExportAccounts(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	var req request.ExportReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	data, filename, err := c.service.ExportAccounts(userID, req.Format)
	if err != nil {
		handleError(ctx, err)
		return
	}

	contentType := "text/csv"
	if req.Format == "json" {
		contentType = "application/json"
	}

	ctx.Header("Content-Type", contentType)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Data(http.StatusOK, contentType, data)
}

func (c *ExportController) ExportBudgets(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	var req request.ExportReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	data, filename, err := c.service.ExportBudgets(userID, req.Format)
	if err != nil {
		handleError(ctx, err)
		return
	}

	contentType := "text/csv"
	if req.Format == "json" {
		contentType = "application/json"
	}

	ctx.Header("Content-Type", contentType)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Data(http.StatusOK, contentType, data)
}

func (c *ExportController) ExportCategories(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	var req request.ExportReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	data, filename, err := c.service.ExportCategories(userID, req.Format)
	if err != nil {
		handleError(ctx, err)
		return
	}

	contentType := "text/csv"
	if req.Format == "json" {
		contentType = "application/json"
	}

	ctx.Header("Content-Type", contentType)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Data(http.StatusOK, contentType, data)
}

func (c *ExportController) ExportTags(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")

	var req request.ExportReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	data, filename, err := c.service.ExportTags(userID, req.Format)
	if err != nil {
		handleError(ctx, err)
		return
	}

	contentType := "text/csv"
	if req.Format == "json" {
		contentType = "application/json"
	}

	ctx.Header("Content-Type", contentType)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Data(http.StatusOK, contentType, data)
}
