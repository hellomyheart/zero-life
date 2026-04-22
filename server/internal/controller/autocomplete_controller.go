package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/response"
	"github.com/hellomyheart/zero-life/server/internal/repository"
)

const autocompleteLimit = 20

type AutocompleteController struct {
	accountRepo  *repository.AccountRepository
	categoryRepo *repository.CategoryRepository
	tagRepo      *repository.TagRepository
	currencyRepo *repository.CurrencyRepository
}

func NewAutocompleteController(
	accountRepo *repository.AccountRepository,
	categoryRepo *repository.CategoryRepository,
	tagRepo *repository.TagRepository,
	currencyRepo *repository.CurrencyRepository,
) *AutocompleteController {
	return &AutocompleteController{
		accountRepo:  accountRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
		currencyRepo: currencyRepo,
	}
}

func (c *AutocompleteController) Accounts(ctx *gin.Context) {
	userID := ctx.GetUint64("userID")
	query := strings.ToLower(ctx.Query("q"))

	accounts, err := c.accountRepo.List(userID, "", query, "name", 0, autocompleteLimit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	ctx.JSON(http.StatusOK, items)
}

func (c *AutocompleteController) Categories(ctx *gin.Context) {
	userID := ctx.GetUint64("userID")
	query := strings.ToLower(ctx.Query("q"))

	categories, err := c.categoryRepo.List(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	ctx.JSON(http.StatusOK, items)
}

func (c *AutocompleteController) Tags(ctx *gin.Context) {
	userID := ctx.GetUint64("userID")
	query := strings.ToLower(ctx.Query("q"))

	tags, err := c.tagRepo.List(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	ctx.JSON(http.StatusOK, items)
}

func (c *AutocompleteController) Currencies(ctx *gin.Context) {
	query := strings.ToLower(ctx.Query("q"))

	currencies, err := c.currencyRepo.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	ctx.JSON(http.StatusOK, items)
}
