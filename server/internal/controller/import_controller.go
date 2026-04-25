// Package controller 提供HTTP请求处理控制器
// ImportController 数据导入控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// ImportController 数据导入控制器
// 处理CSV等格式文件的上传、解析预览和正式导入等HTTP请求
// 导入流程：上传文件 → 解析预览（确认字段映射） → 执行导入
type ImportController struct {
	importService *service.ImportService // 导入业务服务
}

// NewImportController 创建导入控制器实例
// 参数：
//   - importService: 导入业务服务实例
// 返回：
//   - *ImportController: 导入控制器实例
func NewImportController(importService *service.ImportService) *ImportController {
	return &ImportController{importService: importService}
}

// Upload 上传导入文件
// 接收CSV等格式的文件上传，返回文件ID用于后续解析和导入
// 参数：
//   - c: Gin上下文
// 表单字段：file（导入文件）
// 响应：ImportUploadResp（包含file_id）
// @Summary      Upload import file
// @Description  Upload a CSV file for data import, returns file ID for subsequent parsing
// @Tags         imports
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "File to import"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/imports/upload [post]
// @Security     BearerAuth
func (ctrl *ImportController) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrImportFileInvalid)
		return
	}
	defer file.Close()

	result, err := ctrl.importService.Upload(header.Filename, file)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Parse 解析导入文件（预览）
// 根据文件ID和字段映射规则，解析上传的文件并返回预览数据
// 参数：
//   - c: Gin上下文，包含用户身份
// 请求体：ImportParseReq（包含file_id和mapping字段映射）
// 响应：ImportPreviewResp（包含总行数、有效行数、无效行数和每行预览数据）
// @Summary      Parse import file
// @Description  Parse uploaded file with field mapping and return preview data
// @Tags         imports
// @Accept       json
// @Produce      json
// @Param        body body request.ImportParseReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/imports/parse [post]
// @Security     BearerAuth
func (ctrl *ImportController) Parse(c *gin.Context) {
	// 获取用户ID，用于后续权限校验（确保用户只能解析自己上传的文件）
	userID := c.GetUint64("user_id")

	var req request.ImportParseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.importService.Parse(req.FileID, req.Mapping)
	if err != nil {
		handleError(c, err)
		return
	}

	// userID已获取，后续可根据需要添加文件归属校验
	_ = userID

	Success(c, result)
}

// Execute 执行导入
// 根据文件ID和字段映射规则，正式执行数据导入，将文件中的数据创建为交易记录
// 参数：
//   - c: Gin上下文，包含用户身份
// 请求体：ImportExecuteReq（包含file_id和mapping字段映射）
// 响应：ImportResultResp（包含总数、成功数、失败数、跳过数）
// @Summary      Execute import
// @Description  Execute data import with field mapping, creating transaction records from file
// @Tags         imports
// @Accept       json
// @Produce      json
// @Param        body body request.ImportExecuteReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/imports/execute [post]
// @Security     BearerAuth
func (ctrl *ImportController) Execute(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.ImportExecuteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.importService.Execute(userID, req.FileID, req.Mapping)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}
