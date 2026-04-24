// Package controller 提供HTTP请求处理控制器
package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

// PopupController 弹窗数据控制器
// 提供各种弹窗选择器所需的数据接口
type PopupController struct {
	accountService  *service.AccountService
	categoryService *service.CategoryService
	tagService      *service.TagService
	currencyService *service.CurrencyService
}

// NewPopupController 创建弹窗控制器实例
// 参数：
//   - accountService: 账户业务服务
//   - categoryService: 分类业务服务
//   - tagService: 标签业务服务
//   - currencyService: 货币业务服务
// 返回：
//   - *PopupController: 弹窗控制器实例
func NewPopupController(
	accountService *service.AccountService,
	categoryService *service.CategoryService,
	tagService *service.TagService,
	currencyService *service.CurrencyService,
) *PopupController {
	return &PopupController{
		accountService:  accountService,
		categoryService: categoryService,
		tagService:      tagService,
		currencyService: currencyService,
	}
}

// GetAccounts 获取账户选择弹窗数据
// 参数：
//   - c: Gin上下文
func (ctrl *PopupController) GetAccounts(c *gin.Context) {
	userID := c.GetUint64("user_id")

	// 获取所有账户（不分页）
	accounts, err := ctrl.accountService.List(userID, &request.AccountListReq{
		Page:     1,
		PageSize: 1000,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, accounts)
}

// GetCategories 获取分类选择弹窗数据
// 参数：
//   - c: Gin上下文
func (ctrl *PopupController) GetCategories(c *gin.Context) {
	userID := c.GetUint64("user_id")

	// 获取所有分类
	categories, err := ctrl.categoryService.List(userID, &request.CategoryListReq{
		Page:     1,
		PageSize: 1000,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, categories)
}

// GetTags 获取标签选择弹窗数据
// 参数：
//   - c: Gin上下文
func (ctrl *PopupController) GetTags(c *gin.Context) {
	userID := c.GetUint64("user_id")

	// 获取所有标签
	tags, err := ctrl.tagService.List(userID, &request.TagListReq{
		Page:     1,
		PageSize: 1000,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, tags)
}

// GetCurrencies 获取货币选择弹窗数据
// 参数：
//   - c: Gin上下文
func (ctrl *PopupController) GetCurrencies(c *gin.Context) {
	// 获取所有启用的货币
	currencies, err := ctrl.currencyService.List(&request.CurrencyListReq{
		Page:     1,
		PageSize: 1000,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, currencies)
}
