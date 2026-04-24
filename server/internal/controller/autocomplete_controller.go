// Package controller 提供HTTP请求处理控制器
// AutocompleteController 自动补全控制器，提供自动补全接口
package controller

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/repository"
)

// autocompleteLimit 自动补全结果数量限制
const autocompleteLimit = 20

// AutocompleteController 自动补全控制器
// 提供账户、分类、标签、货币、预算、账单等数据的自动补全功能
type AutocompleteController struct {
	accountRepo  *repository.AccountRepository  // 账户数据访问对象
	categoryRepo *repository.CategoryRepository // 分类数据访问对象
	tagRepo      *repository.TagRepository      // 标签数据访问对象
	currencyRepo *repository.CurrencyRepository // 货币数据访问对象
	budgetRepo   *repository.BudgetRepository   // 预算数据访问对象
	billRepo     *repository.BillRepository     // 账单数据访问对象
}

// NewAutocompleteController 创建自动补全控制器实例
// 参数：
//   - accountRepo: 账户数据访问对象
//   - categoryRepo: 分类数据访问对象
//   - tagRepo: 标签数据访问对象
//   - currencyRepo: 货币数据访问对象
//   - budgetRepo: 预算数据访问对象
//   - billRepo: 账单数据访问对象
// 返回：
//   - *AutocompleteController: 自动补全控制器实例
func NewAutocompleteController(
	accountRepo *repository.AccountRepository,
	categoryRepo *repository.CategoryRepository,
	tagRepo *repository.TagRepository,
	currencyRepo *repository.CurrencyRepository,
	budgetRepo *repository.BudgetRepository,
	billRepo *repository.BillRepository,
) *AutocompleteController {
	return &AutocompleteController{
		accountRepo:  accountRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
		currencyRepo: currencyRepo,
		budgetRepo:   budgetRepo,
		billRepo:     billRepo,
	}
}

// Accounts 账户自动补全
// 根据查询关键词返回匹配的账户列表
// 参数：
//   - ctx: Gin上下文
func (c *AutocompleteController) Accounts(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	query := strings.ToLower(ctx.Query("q"))

	accounts, err := c.accountRepo.List(userID, "", query, "name", 0, autocompleteLimit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	items := make([]response.AutocompleteItemResp, 0, len(accounts))
	for _, a := range accounts {
		items = append(items, response.AutocompleteItemResp{
			ID:   a.ID,
			Name: a.Name,
			Type: string(a.Type),
		})
	}
	Success(ctx, items)
}

// Categories 分类自动补全
// 根据查询关键词返回匹配的分类列表
// 参数：
//   - ctx: Gin上下文
func (c *AutocompleteController) Categories(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	query := strings.ToLower(ctx.Query("q"))

	categories, err := c.categoryRepo.List(userID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	items := make([]response.AutocompleteItemResp, 0)
	for _, cat := range categories {
		if query != "" && !strings.Contains(strings.ToLower(cat.Name), query) {
			continue
		}
		items = append(items, response.AutocompleteItemResp{
			ID:   cat.ID,
			Name: cat.Name,
		})
		if len(items) >= autocompleteLimit {
			break
		}
	}
	Success(ctx, items)
}

// Tags 标签自动补全
// 根据查询关键词返回匹配的标签列表
// 参数：
//   - ctx: Gin上下文
func (c *AutocompleteController) Tags(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	query := strings.ToLower(ctx.Query("q"))

	tags, err := c.tagRepo.List(userID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	items := make([]response.AutocompleteItemResp, 0)
	for _, t := range tags {
		if query != "" && !strings.Contains(strings.ToLower(t.Name), query) {
			continue
		}
		items = append(items, response.AutocompleteItemResp{
			ID:    t.ID,
			Name:  t.Name,
			Color: t.Color,
		})
		if len(items) >= autocompleteLimit {
			break
		}
	}
	Success(ctx, items)
}

// Currencies 货币自动补全
// 根据查询关键词返回匹配的货币列表
// 参数：
//   - ctx: Gin上下文
func (c *AutocompleteController) Currencies(ctx *gin.Context) {
	query := strings.ToLower(ctx.Query("q"))

	currencies, err := c.currencyRepo.List()
	if err != nil {
		handleError(ctx, err)
		return
	}

	items := make([]response.AutocompleteItemResp, 0)
	for _, cur := range currencies {
		if query != "" &&
			!strings.Contains(strings.ToLower(cur.Code), query) &&
			!strings.Contains(strings.ToLower(cur.Name), query) {
			continue
		}
		items = append(items, response.AutocompleteItemResp{
			ID:     cur.ID,
			Name:   cur.Name,
			Symbol: cur.Symbol,
			Code:   cur.Code,
		})
		if len(items) >= autocompleteLimit {
			break
		}
	}
	Success(ctx, items)
}

// Budgets 预算自动补全
// 根据查询关键词返回匹配的预算列表
// 参数：
//   - ctx: Gin上下文
func (c *AutocompleteController) Budgets(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	query := strings.ToLower(ctx.Query("q"))

	budgets, err := c.budgetRepo.List(userID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	items := make([]response.AutocompleteItemResp, 0)
	for _, b := range budgets {
		if query != "" && !strings.Contains(strings.ToLower(b.Name), query) {
			continue
		}
		items = append(items, response.AutocompleteItemResp{
			ID:   b.ID,
			Name: b.Name,
		})
		if len(items) >= autocompleteLimit {
			break
		}
	}
	Success(ctx, items)
}

// Bills 账单自动补全
// 根据查询关键词返回匹配的账单列表
// 参数：
//   - ctx: Gin上下文
func (c *AutocompleteController) Bills(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	query := strings.ToLower(ctx.Query("q"))

	bills, err := c.billRepo.List(userID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	items := make([]response.AutocompleteItemResp, 0)
	for _, b := range bills {
		if query != "" && !strings.Contains(strings.ToLower(b.Name), query) {
			continue
		}
		items = append(items, response.AutocompleteItemResp{
			ID:   b.ID,
			Name: b.Name,
		})
		if len(items) >= autocompleteLimit {
			break
		}
	}
	Success(ctx, items)
}
