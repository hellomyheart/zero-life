// Package controller 提供HTTP请求处理控制器
// AttachmentController 附件控制器
package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// AttachmentController 附件管理控制器
// 处理附件的上传、下载、在线预览、列表查询和删除等HTTP请求
// 附件可以关联到交易、账单等不同类型的对象（attachable）
type AttachmentController struct {
	service *service.AttachmentService // 附件业务服务
}

// NewAttachmentController 创建附件控制器实例
// 参数：
//   - service: 附件业务服务实例
// 返回：
//   - *AttachmentController: 附件控制器实例
func NewAttachmentController(service *service.AttachmentService) *AttachmentController {
	return &AttachmentController{service: service}
}

// Upload 上传附件
// 接收multipart/form-data格式的文件上传请求，将文件关联到指定对象
// 参数：
//   - ctx: Gin上下文，包含用户身份
// 表单字段：attachable_type（关联对象类型）、attachable_id（关联对象ID）、file（文件）
// 响应：上传成功的附件信息
// @Summary      Upload attachment
// @Description  Upload a file and attach it to an object
// @Tags         attachments
// @Accept       multipart/form-data
// @Produce      json
// @Param        attachable_type formData string true "Attachable object type"
// @Param        attachable_id formData uint64 true "Attachable object ID"
// @Param        file formData file true "File to upload"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/attachments/upload [post]
// @Security     BearerAuth
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

// Download 下载附件
// 以附件形式下载指定ID的文件
// 参数：
//   - ctx: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（附件ID）
// 响应：文件流（Content-Disposition: attachment）
// @Summary      Download attachment
// @Description  Download an attachment file by ID
// @Tags         attachments
// @Accept       json
// @Produce      octet-stream
// @Param        id path uint64 true "Attachment ID"
// @Success      200  {file} file
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/attachments/{id}/download [get]
// @Security     BearerAuth
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

// View 在线预览附件
// 以内联方式在浏览器中预览指定ID的文件
// 参数：
//   - ctx: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（附件ID）
// 响应：文件流（Content-Disposition: inline）
// @Summary      View attachment
// @Description  View an attachment file inline in browser
// @Tags         attachments
// @Accept       json
// @Produce      octet-stream
// @Param        id path uint64 true "Attachment ID"
// @Success      200  {file} file
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/attachments/{id}/view [get]
// @Security     BearerAuth
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

// List 获取附件列表
// 根据关联对象类型和ID查询附件列表
// 参数：
//   - ctx: Gin上下文，包含用户身份
// 查询参数：attachable_type（关联对象类型）、attachable_id（关联对象ID）
// 响应：附件列表
// @Summary      List attachments
// @Description  Get attachment list by attachable type and ID
// @Tags         attachments
// @Accept       json
// @Produce      json
// @Param        attachable_type query string false "Attachable object type"
// @Param        attachable_id query uint64 false "Attachable object ID"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/attachments [get]
// @Security     BearerAuth
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

// Delete 删除附件
// 参数：
//   - ctx: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（附件ID）
// 响应：删除成功返回nil
// @Summary      Delete attachment
// @Description  Delete an attachment by ID
// @Tags         attachments
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Attachment ID"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/attachments/{id} [delete]
// @Security     BearerAuth
func (c *AttachmentController) Delete(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	id := parseIDParam(ctx, "id")

	if err := c.service.Delete(userID, id); err != nil {
		handleError(ctx, err)
		return
	}
	Success(ctx, nil)
}
