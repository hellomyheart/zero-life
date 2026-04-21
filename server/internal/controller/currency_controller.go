package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

type CurrencyController struct {
	currencyService *service.CurrencyService
}

func NewCurrencyController(currencyService *service.CurrencyService) *CurrencyController {
	return &CurrencyController{currencyService: currencyService}
}

func (ctrl *CurrencyController) List(c *gin.Context) {
	result, err := ctrl.currencyService.List()
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *CurrencyController) UpdateStatus(c *gin.Context) {
	id := parseIDParam(c, "id")

	var req request.SetCurrencyStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := ctrl.currencyService.UpdateStatus(id, req.IsEnabled); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

func (ctrl *CurrencyController) SetDefault(c *gin.Context) {
	id := parseIDParam(c, "id")

	if err := ctrl.currencyService.SetDefault(id); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

func (ctrl *CurrencyController) GetExchangeRates(c *gin.Context) {
	result, err := ctrl.currencyService.GetExchangeRates()
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *CurrencyController) SetExchangeRate(c *gin.Context) {
	var req request.SetExchangeRateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := ctrl.currencyService.SetExchangeRate(&req); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}
