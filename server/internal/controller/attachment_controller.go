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
	userID := ctx.GetUint64("user_id")

	attachableType := ctx.PostForm("attachable_type")
	attachableIDStr := ctx.PostForm("attachable_id")
	attachableID, err := strconv.ParseUint(attachableIDStr, 10, 64)
	if err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		Error(ctx, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	resp, err := c.service.Upload(userID, attachableType, attachableID, fileHeader)
	if err != nil {
		handleError(ctx, err)
		return
	}
	Success(ctx, resp)
}

func (c *AttachmentController) Download(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	id := parseIDParam(ctx, "id")

	filePath, err := c.service.Download(userID, id)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.FileAttachment(filePath, "")
}

func (c *AttachmentController) View(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	id := parseIDParam(ctx, "id")

	result, err := c.service.View(userID, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.Header("Content-Type", result.Mime)
	ctx.Header("Content-Disposition", "inline; filename="+result.Filename)
	ctx.File(result.FilePath)
}

func (c *AttachmentController) List(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	attachableType := ctx.Query("attachable_type")
	attachableID, _ := strconv.ParseUint(ctx.Query("attachable_id"), 10, 64)

	resp, err := c.service.List(userID, attachableType, attachableID)
	if err != nil {
		handleError(ctx, err)
		return
	}
	Success(ctx, resp)
}

func (c *AttachmentController) Delete(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	id := parseIDParam(ctx, "id")

	if err := c.service.Delete(userID, id); err != nil {
		handleError(ctx, err)
		return
	}
	Success(ctx, nil)
}
