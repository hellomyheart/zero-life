package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

type AttachmentController struct {
	service *service.AttachmentService
}

func NewAttachmentController(service *service.AttachmentService) *AttachmentController {
	return &AttachmentController{service: service}
}

func (c *AttachmentController) Upload(ctx *gin.Context) {
	userID := ctx.GetUint64("userID")

	attachableType := ctx.PostForm("attachable_type")
	attachableIDStr := ctx.PostForm("attachable_id")
	attachableID, err := strconv.ParseUint(attachableIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid attachable_id"})
		return
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}

	resp, err := c.service.Upload(userID, attachableType, attachableID, fileHeader)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, resp)
}

func (c *AttachmentController) Download(ctx *gin.Context) {
	userID := ctx.GetUint64("userID")
	id := ctx.GetUint64("id")

	filePath, err := c.service.Download(userID, id)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.File(filePath)
}

func (c *AttachmentController) List(ctx *gin.Context) {
	userID := ctx.GetUint64("userID")
	attachableType := ctx.Query("attachable_type")
	attachableID, _ := strconv.ParseUint(ctx.Query("attachable_id"), 10, 64)

	resp, err := c.service.List(userID, attachableType, attachableID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (c *AttachmentController) Delete(ctx *gin.Context) {
	userID := ctx.GetUint64("userID")
	id := ctx.GetUint64("id")

	if err := c.service.Delete(userID, id); err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
