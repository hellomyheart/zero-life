package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/controller"
	"github.com/hellomyheart/zero-life/server/internal/middleware"
	"github.com/hellomyheart/zero-life/server/internal/pkg/jwt"
)

type Router struct {
	engine *gin.Engine

	authCtrl        *controller.AuthController
	accountCtrl     *controller.AccountController
	txnCtrl         *controller.TransactionController
	categoryCtrl    *controller.CategoryController
	tagCtrl         *controller.TagController
	budgetCtrl      *controller.BudgetController
	billCtrl        *controller.BillController
	currencyCtrl    *controller.CurrencyController
	ruleCtrl        *controller.RuleController
	reportCtrl      *controller.ReportController
	dashboardCtrl   *controller.DashboardController
	importCtrl      *controller.ImportController
	piggyBankCtrl   *controller.PiggyBankController
	attachmentCtrl  *controller.AttachmentController
	autocompleteCtrl *controller.AutocompleteController
	exportCtrl      *controller.ExportController
}

func NewRouter(
	engine *gin.Engine,
	authCtrl *controller.AuthController,
	accountCtrl *controller.AccountController,
	txnCtrl *controller.TransactionController,
	categoryCtrl *controller.CategoryController,
	tagCtrl *controller.TagController,
	budgetCtrl *controller.BudgetController,
	billCtrl *controller.BillController,
	currencyCtrl *controller.CurrencyController,
	ruleCtrl *controller.RuleController,
	reportCtrl *controller.ReportController,
	dashboardCtrl *controller.DashboardController,
	importCtrl *controller.ImportController,
	piggyBankCtrl *controller.PiggyBankController,
	attachmentCtrl *controller.AttachmentController,
	autocompleteCtrl *controller.AutocompleteController,
	exportCtrl *controller.ExportController,
) *Router {
	return &Router{
		engine:           engine,
		authCtrl:         authCtrl,
		accountCtrl:      accountCtrl,
		txnCtrl:          txnCtrl,
		categoryCtrl:     categoryCtrl,
		tagCtrl:          tagCtrl,
		budgetCtrl:       budgetCtrl,
		billCtrl:         billCtrl,
		currencyCtrl:     currencyCtrl,
		ruleCtrl:         ruleCtrl,
		reportCtrl:       reportCtrl,
		dashboardCtrl:    dashboardCtrl,
		importCtrl:       importCtrl,
		piggyBankCtrl:    piggyBankCtrl,
		attachmentCtrl:   attachmentCtrl,
		autocompleteCtrl: autocompleteCtrl,
		exportCtrl:       exportCtrl,
	}
}

func (r *Router) Setup(jwtService *jwt.Service) {
	// Global middleware
	r.engine.Use(middleware.CORS())

	// API v1 group
	v1 := r.engine.Group("/api/v1")

	// Public routes (no auth required)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", r.authCtrl.Register)
		auth.POST("/login", r.authCtrl.Login)
		auth.POST("/refresh", r.authCtrl.RefreshToken)
		auth.POST("/forgot-password", r.authCtrl.ForgotPassword)
		auth.POST("/reset-password", r.authCtrl.ResetPassword)
	}

	// Authenticated routes
	authenticated := v1.Group("")
	authenticated.Use(middleware.Auth(jwtService))
	{
		// Auth profile
		authAuth := authenticated.Group("/auth")
		{
			authAuth.GET("/profile", r.authCtrl.GetProfile)
			authAuth.PUT("/profile", r.authCtrl.UpdateProfile)
			authAuth.PUT("/password", r.authCtrl.ChangePassword)
		}

		// Accounts
		accounts := authenticated.Group("/accounts")
		{
			accounts.POST("", r.accountCtrl.Create)
			accounts.GET("", r.accountCtrl.List)
			accounts.GET("/:id", r.accountCtrl.Get)
			accounts.PUT("/:id", r.accountCtrl.Update)
			accounts.DELETE("/:id", r.accountCtrl.Delete)
		}

		// Transactions
		transactions := authenticated.Group("/transactions")
		{
			transactions.POST("", r.txnCtrl.Create)
			transactions.GET("", r.txnCtrl.List)
			transactions.GET("/search", r.txnCtrl.Search)
			transactions.GET("/:id", r.txnCtrl.Get)
			transactions.PUT("/:id", r.txnCtrl.Update)
			transactions.DELETE("/:id", r.txnCtrl.Delete)
		}

		// Categories
		categories := authenticated.Group("/categories")
		{
			categories.POST("", r.categoryCtrl.Create)
			categories.GET("", r.categoryCtrl.List)
			categories.PUT("/:id", r.categoryCtrl.Update)
			categories.DELETE("/:id", r.categoryCtrl.Delete)
		}

		// Tags
		tags := authenticated.Group("/tags")
		{
			tags.POST("", r.tagCtrl.Create)
			tags.GET("", r.tagCtrl.List)
			tags.PUT("/:id", r.tagCtrl.Update)
			tags.DELETE("/:id", r.tagCtrl.Delete)
		}

		// Budgets
		budgets := authenticated.Group("/budgets")
		{
			budgets.POST("", r.budgetCtrl.Create)
			budgets.GET("", r.budgetCtrl.List)
			budgets.GET("/:id", r.budgetCtrl.Get)
			budgets.PUT("/:id", r.budgetCtrl.Update)
			budgets.DELETE("/:id", r.budgetCtrl.Delete)
			budgets.GET("/:id/history", r.budgetCtrl.GetHistory)
		}

		// Bills
		bills := authenticated.Group("/bills")
		{
			bills.POST("", r.billCtrl.Create)
			bills.GET("", r.billCtrl.List)
			bills.GET("/:id", r.billCtrl.Get)
			bills.PUT("/:id", r.billCtrl.Update)
			bills.DELETE("/:id", r.billCtrl.Delete)
		}

		// Currencies
		currencies := authenticated.Group("/currencies")
		{
			currencies.GET("", r.currencyCtrl.List)
			currencies.PUT("/:id/status", r.currencyCtrl.UpdateStatus)
			currencies.PUT("/:id/default", r.currencyCtrl.SetDefault)
			currencies.GET("/exchange-rates", r.currencyCtrl.GetExchangeRates)
			currencies.POST("/exchange-rates", r.currencyCtrl.SetExchangeRate)
		}

		// Rules
		rules := authenticated.Group("/rules")
		{
			rules.POST("", r.ruleCtrl.Create)
			rules.GET("", r.ruleCtrl.List)
			rules.GET("/:id", r.ruleCtrl.Get)
			rules.PUT("/:id", r.ruleCtrl.Update)
			rules.DELETE("/:id", r.ruleCtrl.Delete)
			rules.PUT("/:id/toggle", r.ruleCtrl.ToggleStatus)
			rules.POST("/:id/execute", r.ruleCtrl.Execute)
		}

		// Reports
		reports := authenticated.Group("/reports")
		{
			reports.GET("/income-expense", r.reportCtrl.IncomeExpense)
			reports.GET("/category", r.reportCtrl.Category)
			reports.GET("/budget", r.reportCtrl.Budget)
			reports.GET("/net-worth", r.reportCtrl.NetWorth)
			reports.GET("/trend", r.reportCtrl.Trend)
			reports.GET("/tag", r.reportCtrl.Tag)
		}

		// Import
		imports := authenticated.Group("/imports")
		{
			imports.POST("/upload", r.importCtrl.Upload)
			imports.POST("/parse", r.importCtrl.Parse)
			imports.POST("/execute", r.importCtrl.Execute)
		}

		// Dashboard
		dashboard := authenticated.Group("/dashboard")
		{
			dashboard.GET("", r.dashboardCtrl.Get)
		}

		// Piggy Banks
		piggyBanks := authenticated.Group("/piggy-banks")
		{
			piggyBanks.POST("", r.piggyBankCtrl.Create)
			piggyBanks.GET("", r.piggyBankCtrl.List)
			piggyBanks.GET("/:id", r.piggyBankCtrl.Get)
			piggyBanks.PUT("/:id", r.piggyBankCtrl.Update)
			piggyBanks.DELETE("/:id", r.piggyBankCtrl.Delete)
			piggyBanks.POST("/:id/add", r.piggyBankCtrl.AddAmount)
			piggyBanks.POST("/:id/remove", r.piggyBankCtrl.RemoveAmount)
			piggyBanks.GET("/:id/events", r.piggyBankCtrl.GetEvents)
		}

		// Attachments
		attachments := authenticated.Group("/attachments")
		{
			attachments.POST("/upload", r.attachmentCtrl.Upload)
			attachments.GET("/:id/download", r.attachmentCtrl.Download)
			attachments.GET("", r.attachmentCtrl.List)
			attachments.DELETE("/:id", r.attachmentCtrl.Delete)
		}

		// Autocomplete
		autocomplete := authenticated.Group("/autocomplete")
		{
			autocomplete.GET("/accounts", r.autocompleteCtrl.Accounts)
			autocomplete.GET("/categories", r.autocompleteCtrl.Categories)
			autocomplete.GET("/tags", r.autocompleteCtrl.Tags)
			autocomplete.GET("/currencies", r.autocompleteCtrl.Currencies)
		}

		// Exports
		exports := authenticated.Group("/exports")
		{
			exports.GET("/transactions", r.exportCtrl.ExportTransactions)
			exports.GET("/accounts", r.exportCtrl.ExportAccounts)
		}
	}
}
