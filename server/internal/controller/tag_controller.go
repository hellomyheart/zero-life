// Package controller 提供HTTP请求处理控制器
// TagController 标签控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// TagController 标签管理控制器
// 处理标签的创建、列表查询、详情查询、更新和删除等HTTP请求
// 标签用于对交易进行标记和分类，支持自定义颜色和两级树形结构
type TagController struct {
	tagService *service.TagService // 标签业务服务
}

// NewTagController 创建标签控制器实例
// 参数：
//   - tagService: 标签业务服务实例
// 返回：
//   - *TagController: 标签控制器实例
func NewTagController(tagService *service.TagService) *TagController {
	return &TagController{tagService: tagService}
}

// Create 创建标签
// 接收创建标签请求，设置标签名称、颜色和父标签ID
// 参数：
//   - c: Gin上下文，包含用户身份和请求体
// 请求体：CreateTagReq（包含名称、颜色、父标签ID）
// 响应：创建成功的标签信息
// 业务规则：标签名称在同一用户下不能重复；最多两级树形结构
// @Summary      Create tag
// @Description  Create a new tag for the current user, supports parent_id for tree structure
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        body body request.CreateTagReq true "create tag request"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/tags [post]
// @Security     BearerAuth
func (ctrl *TagController) Create(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.CreateTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.tagService.Create(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// List 获取标签列表
// 返回当前用户的所有标签（树形结构），含关联交易数量
// 参数：
//   - c: Gin上下文，包含用户身份
// 响应：标签树形列表
// @Summary      List tags
// @Description  Get all tags for the current user in tree structure
// @Tags         tags
// @Accept       json
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/tags [get]
// @Security     BearerAuth
func (ctrl *TagController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.tagService.List(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Get 获取标签详情
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（标签ID）
// 响应：标签详情信息（含关联交易数）
// @Summary      Get tag
// @Description  Get tag by ID
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Tag ID"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/tags/{id} [get]
// @Security     BearerAuth
func (ctrl *TagController) Get(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	result, err := ctrl.tagService.Get(userID, id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Update 更新标签
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（标签ID）
// 请求体：UpdateTagReq（包含名称、颜色、父标签ID）
// 响应：更新后的标签信息
// @Summary      Update tag
// @Description  Update tag by ID, supports changing parent_id
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Tag ID"
// @Param        body body request.UpdateTagReq true "update tag request"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/tags/{id} [put]
// @Security     BearerAuth
func (ctrl *TagController) Update(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	var req request.UpdateTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.tagService.Update(userID, id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Delete 删除标签
// 同时删除所有子标签及关联的transaction_tags记录
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（标签ID）
// 响应：删除成功返回nil
// @Summary      Delete tag
// @Description  Delete tag by ID, also deletes child tags
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Tag ID"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/tags/{id} [delete]
// @Security     BearerAuth
func (ctrl *TagController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.tagService.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}
