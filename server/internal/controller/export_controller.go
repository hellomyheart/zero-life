package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

type ExportController struct {
	service *service.ExportService
}

func NewExportController(service *service.ExportService) *ExportController {
	return &ExportController{service: service}
}

func (c *ExportController) ExportTransactions(ctx *gin.Context) {
	userID := ctx.GetUint64("userID")

	var req request.ExportReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, filename, err := c.service.ExportTransactions(userID, req.StartDate, req.EndDate, req.Format)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	userID := ctx.GetUint64("userID")

	var req request.ExportReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, filename, err := c.service.ExportAccounts(userID, req.Format)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
