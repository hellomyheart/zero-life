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
	accountRepo  *repository.AccountRepository
	categoryRepo *repository.CategoryRepository
	tagRepo      *repository.TagRepository
	currencyRepo *repository.CurrencyRepository
	budgetRepo   *repository.BudgetRepository
	rtRepo       *repository.RecurringTransactionRepository
}

func NewAutocompleteController(
	accountRepo *repository.AccountRepository,
	categoryRepo *repository.CategoryRepository,
	tagRepo *repository.TagRepository,
	currencyRepo *repository.CurrencyRepository,
	budgetRepo *repository.BudgetRepository,
	rtRepo *repository.RecurringTransactionRepository,
) *AutocompleteController {
	return &AutocompleteController{
		accountRepo:  accountRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
		currencyRepo: currencyRepo,
		budgetRepo:   budgetRepo,
		rtRepo:       rtRepo,
	}
}

// Accounts 账户自动补全
// 根据查询关键词返回匹配的账户列表
// 参数：
//   - ctx: Gin上下文
// @Summary      Autocomplete accounts
// @Description  Search accounts by keyword for autocomplete
// @Tags         autocomplete
// @Accept       json
// @Produce      json
// @Param        q query string false "Search keyword"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/autocomplete/accounts [get]
// @Security     BearerAuth
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
// @Summary      Autocomplete categories
// @Description  Search categories by keyword for autocomplete
// @Tags         autocomplete
// @Accept       json
// @Produce      json
// @Param        q query string false "Search keyword"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/autocomplete/categories [get]
// @Security     BearerAuth
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
// @Summary      Autocomplete tags
// @Description  Search tags by keyword for autocomplete
// @Tags         autocomplete
// @Accept       json
// @Produce      json
// @Param        q query string false "Search keyword"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/autocomplete/tags [get]
// @Security     BearerAuth
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
// @Summary      Autocomplete currencies
// @Description  Search currencies by keyword for autocomplete
// @Tags         autocomplete
// @Accept       json
// @Produce      json
// @Param        q query string false "Search keyword"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/autocomplete/currencies [get]
// @Security     BearerAuth
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
// @Summary      Autocomplete budgets
// @Description  Search budgets by keyword for autocomplete
// @Tags         autocomplete
// @Accept       json
// @Produce      json
// @Param        q query string false "Search keyword"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/autocomplete/budgets [get]
// @Security     BearerAuth
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
// @Summary      Autocomplete bills
// @Description  Search bills by keyword for autocomplete
// @Tags         autocomplete
// @Accept       json
// @Produce      json
// @Param        q query string false "Search keyword"
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/v1/autocomplete/bills [get]
// @Security     BearerAuth
func (c *AutocompleteController) Bills(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	query := strings.ToLower(ctx.Query("q"))

	rts, err := c.rtRepo.List(userID, nil, 0, 10000)
	if err != nil {
		handleError(ctx, err)
		return
	}

	items := make([]response.AutocompleteItemResp, 0)
	for _, rt := range rts {
		if query != "" && !strings.Contains(strings.ToLower(rt.Description), query) {
			continue
		}
		items = append(items, response.AutocompleteItemResp{
			ID:   rt.ID,
			Name: rt.Description,
		})
		if len(items) >= autocompleteLimit {
			break
		}
	}
	Success(ctx, items)
}
