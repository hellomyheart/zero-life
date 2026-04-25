// Package controller 提供HTTP请求处理控制器
// LinkTypeController 关联类型控制器
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// LinkTypeController 关联类型管理控制器
// 处理关联类型的创建、查询、更新和删除等HTTP请求
// 关联类型定义了交易之间关联关系的语义（如冲正、对账、关联等）
// 关联类型是系统级资源，不区分用户
type LinkTypeController struct {
	linkTypeService *service.LinkTypeService // 关联类型业务服务
}

// NewLinkTypeController 创建关联类型控制器实例
// 参数：
//   - linkTypeService: 关联类型业务服务实例
// 返回：
//   - *LinkTypeController: 关联类型控制器实例
func NewLinkTypeController(linkTypeService *service.LinkTypeService) *LinkTypeController {
	return &LinkTypeController{linkTypeService: linkTypeService}
}

// Create 创建关联类型
// 参数：
//   - c: Gin上下文
// 请求体：CreateLinkTypeReq（包含名称、正向描述、反向描述、是否有方向性）
// 响应：创建成功的关联类型信息
// @Summary      Create link type
// @Description  Create a new link type defining transaction relationship semantics
// @Tags         link-types
// @Accept       json
// @Produce      json
// @Param        body body request.CreateLinkTypeReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/link-types [post]
// @Security     BearerAuth
func (ctrl *LinkTypeController) Create(c *gin.Context) {
	var req request.CreateLinkTypeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.linkTypeService.Create(&req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Get 获取关联类型详情
// 参数：
//   - c: Gin上下文，包含URL路径参数id
// 路径参数：id（关联类型ID）
// 响应：关联类型详细信息
// @Summary      Get link type
// @Description  Get link type details by ID
// @Tags         link-types
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Link type ID"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/link-types/{id} [get]
// @Security     BearerAuth
func (ctrl *LinkTypeController) Get(c *gin.Context) {
	id := parseIDParam(c, "id")

	result, err := ctrl.linkTypeService.Get(id)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// List 获取关联类型列表（分页）
// 参数：
//   - c: Gin上下文
// 查询参数：LinkTypeListReq（包含page、page_size）
// 响应：关联类型分页列表
// @Summary      List link types
// @Description  Get paginated list of link types
// @Tags         link-types
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/link-types [get]
// @Security     BearerAuth
func (ctrl *LinkTypeController) List(c *gin.Context) {
	var req request.LinkTypeListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.linkTypeService.List(&req)
	if err != nil {
		handleError(c, err)
		return
	}

	SuccessPage(c, result)
}

// Update 更新关联类型
// 参数：
//   - c: Gin上下文，包含URL路径参数id
// 路径参数：id（关联类型ID）
// 请求体：UpdateLinkTypeReq（包含需要更新的字段）
// 响应：更新后的关联类型信息
// @Summary      Update link type
// @Description  Update link type by ID
// @Tags         link-types
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Link type ID"
// @Param        body body request.UpdateLinkTypeReq true "request body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/link-types/{id} [put]
// @Security     BearerAuth
func (ctrl *LinkTypeController) Update(c *gin.Context) {
	id := parseIDParam(c, "id")

	var req request.UpdateLinkTypeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.linkTypeService.Update(id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

// Delete 删除关联类型
// 参数：
//   - c: Gin上下文，包含URL路径参数id
// 路径参数：id（关联类型ID）
// 响应：删除成功返回nil
// @Summary      Delete link type
// @Description  Delete link type by ID
// @Tags         link-types
// @Accept       json
// @Produce      json
// @Param        id path uint64 true "Link type ID"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/link-types/{id} [delete]
// @Security     BearerAuth
func (ctrl *LinkTypeController) Delete(c *gin.Context) {
	id := parseIDParam(c, "id")

	if err := ctrl.linkTypeService.Delete(id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}
