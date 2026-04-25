// Package controller 提供HTTP请求处理控制器
// CategoryController 分类控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// CategoryController 分类管理控制器
// 处理分类的创建、列表查询、更新和删除等HTTP请求
// 分类支持层级结构（最多2级），用于对交易进行归类
type CategoryController struct {
	categoryService *service.CategoryService // 分类业务服务
}

// NewCategoryController 创建分类控制器实例
// 参数：
//   - categoryService: 分类业务服务实例
// 返回：
//   - *CategoryController: 分类控制器实例
func NewCategoryController(categoryService *service.CategoryService) *CategoryController {
	return &CategoryController{categoryService: categoryService}
}

// Create 创建分类
// 接收创建分类请求，支持设置父分类、图标和备注
// 参数：
//   - c: Gin上下文，包含用户身份和请求体
// 请求体：CreateCategoryReq（包含名称、父分类ID、图标、备注）
// 响应：创建成功的分类信息
// 业务规则：分类名称在同一用户下不能重复，层级不能超过2级
func (ctrl *CategoryController) Create(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.CreateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.categoryService.Create(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// List 获取分类列表
// 返回当前用户的所有分类，支持树形结构（包含子分类）
// 参数：
//   - c: Gin上下文，包含用户身份
// 响应：分类列表（含子分类）
func (ctrl *CategoryController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.categoryService.List(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Update 更新分类
// 根据URL路径中的分类ID和请求体中的更新字段，修改分类信息
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（分类ID）
// 请求体：UpdateCategoryReq（包含需要更新的字段）
// 响应：更新后的分类信息
func (ctrl *CategoryController) Update(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	var req request.UpdateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.categoryService.Update(userID, id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Delete 删除分类
// 根据URL路径中的分类ID删除分类
// 参数：
//   - c: Gin上下文，包含用户身份和URL路径参数id
// 路径参数：id（分类ID）
// 响应：删除成功返回nil
// 业务规则：有子分类或关联交易时可能无法删除
func (ctrl *CategoryController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.categoryService.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}
