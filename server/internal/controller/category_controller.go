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

type CategoryController struct {
	categoryService *service.CategoryService
}

func NewCategoryController(categoryService *service.CategoryService) *CategoryController {
	return &CategoryController{categoryService: categoryService}
}

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

func (ctrl *CategoryController) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.categoryService.List(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

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

func (ctrl *CategoryController) Delete(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id := parseIDParam(c, "id")

	if err := ctrl.categoryService.Delete(userID, id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}
