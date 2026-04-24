package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/controller"
	"github.com/hellomyheart/zero-life/server/internal/middleware"
	"github.com/hellomyheart/zero-life/server/internal/pkg/jwt"
)

type Router struct {
	engine *gin.Engine

	authCtrl              *controller.AuthController
	accountCtrl           *controller.AccountController
	txnCtrl               *controller.TransactionController
	categoryCtrl          *controller.CategoryController
	tagCtrl               *controller.TagController
	budgetCtrl            *controller.BudgetController
	billCtrl              *controller.BillController
	currencyCtrl          *controller.CurrencyController
	ruleCtrl              *controller.RuleController
	reportCtrl            *controller.ReportController
	dashboardCtrl         *controller.DashboardController
	importCtrl            *controller.ImportController
	piggyBankCtrl         *controller.PiggyBankController
	attachmentCtrl        *controller.AttachmentController
	autocompleteCtrl      *controller.AutocompleteController
	exportCtrl            *controller.ExportController
	recurringTxnCtrl      *controller.RecurringTransactionController
	webhookCtrl           *controller.WebhookController
	objectGroupCtrl       *controller.ObjectGroupController
	transactionLinkCtrl   *controller.TransactionLinkController
	preferenceCtrl        *controller.PreferenceController
	reconciliationCtrl    *controller.ReconciliationController
	chartCtrl             *controller.ChartController
	insightCtrl           *controller.InsightController
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
	recurringTxnCtrl *controller.RecurringTransactionController,
	webhookCtrl *controller.WebhookController,
	objectGroupCtrl *controller.ObjectGroupController,
	transactionLinkCtrl *controller.TransactionLinkController,
	preferenceCtrl *controller.PreferenceController,
	reconciliationCtrl *controller.ReconciliationController,
	chartCtrl *controller.ChartController,
	insightCtrl *controller.InsightController,
) *Router {
	return &Router{
		engine:               engine,
		authCtrl:             authCtrl,
		accountCtrl:          accountCtrl,
		txnCtrl:              txnCtrl,
		categoryCtrl:         categoryCtrl,
		tagCtrl:              tagCtrl,
		budgetCtrl:           budgetCtrl,
		billCtrl:             billCtrl,
		currencyCtrl:         currencyCtrl,
		ruleCtrl:             ruleCtrl,
		reportCtrl:           reportCtrl,
		dashboardCtrl:        dashboardCtrl,
		importCtrl:           importCtrl,
		piggyBankCtrl:        piggyBankCtrl,
		attachmentCtrl:       attachmentCtrl,
		autocompleteCtrl:     autocompleteCtrl,
		exportCtrl:           exportCtrl,
		recurringTxnCtrl:     recurringTxnCtrl,
		webhookCtrl:          webhookCtrl,
		objectGroupCtrl:      objectGroupCtrl,
		transactionLinkCtrl:  transactionLinkCtrl,
		preferenceCtrl:       preferenceCtrl,
		reconciliationCtrl:   reconciliationCtrl,
		chartCtrl:            chartCtrl,
		insightCtrl:          insightCtrl,
	}
}

func (r *Router) Setup(jwtService *jwt.Service) {
	r.engine.Use(middleware.CORS())

	v1 := r.engine.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/register", r.authCtrl.Register)
		auth.POST("/login", r.authCtrl.Login)
		auth.POST("/refresh", r.authCtrl.RefreshToken)
		auth.POST("/forgot-password", r.authCtrl.ForgotPassword)
		auth.POST("/reset-password", r.authCtrl.ResetPassword)
	}

	authenticated := v1.Group("")
	authenticated.Use(middleware.Auth(jwtService))
	{
		authAuth := authenticated.Group("/auth")
		{
			authAuth.GET("/profile", r.authCtrl.GetProfile)
			authAuth.PUT("/profile", r.authCtrl.UpdateProfile)
			authAuth.PUT("/password", r.authCtrl.ChangePassword)
		}

		accounts := authenticated.Group("/accounts")
		{
			accounts.POST("", r.accountCtrl.Create)
			accounts.GET("", r.accountCtrl.List)
			accounts.GET("/:id", r.accountCtrl.Get)
			accounts.PUT("/:id", r.accountCtrl.Update)
			accounts.DELETE("/:id", r.accountCtrl.Delete)
		}

		transactions := authenticated.Group("/transactions")
		{
			transactions.POST("", r.txnCtrl.Create)
			transactions.GET("", r.txnCtrl.List)
			transactions.GET("/search", r.txnCtrl.Search)
			transactions.GET("/:id", r.txnCtrl.Get)
			transactions.PUT("/:id", r.txnCtrl.Update)
			transactions.DELETE("/:id", r.txnCtrl.Delete)
			// 交易拆分相关路由
			transactions.POST("/:id/split", r.txnCtrl.Split)
			transactions.GET("/:id/splits", r.txnCtrl.GetSplits)
			transactions.POST("/:id/merge", r.txnCtrl.MergeSplits)
		}

		categories := authenticated.Group("/categories")
		{
			categories.POST("", r.categoryCtrl.Create)
			categories.GET("", r.categoryCtrl.List)
			categories.PUT("/:id", r.categoryCtrl.Update)
			categories.DELETE("/:id", r.categoryCtrl.Delete)
		}

		tags := authenticated.Group("/tags")
		{
			tags.POST("", r.tagCtrl.Create)
			tags.GET("", r.tagCtrl.List)
			tags.PUT("/:id", r.tagCtrl.Update)
			tags.DELETE("/:id", r.tagCtrl.Delete)
		}

		budgets := authenticated.Group("/budgets")
		{
			budgets.POST("", r.budgetCtrl.Create)
			budgets.GET("", r.budgetCtrl.List)
			budgets.GET("/:id", r.budgetCtrl.Get)
			budgets.PUT("/:id", r.budgetCtrl.Update)
			budgets.DELETE("/:id", r.budgetCtrl.Delete)
			budgets.GET("/:id/history", r.budgetCtrl.GetHistory)
		}

		bills := authenticated.Group("/bills")
		{
			bills.POST("", r.billCtrl.Create)
			bills.GET("", r.billCtrl.List)
			bills.GET("/:id", r.billCtrl.Get)
			bills.PUT("/:id", r.billCtrl.Update)
			bills.DELETE("/:id", r.billCtrl.Delete)
		}

		currencies := authenticated.Group("/currencies")
		{
			currencies.GET("", r.currencyCtrl.List)
			currencies.PUT("/:id/status", r.currencyCtrl.UpdateStatus)
			currencies.PUT("/:id/default", r.currencyCtrl.SetDefault)
			currencies.GET("/exchange-rates", r.currencyCtrl.GetExchangeRates)
			currencies.POST("/exchange-rates", r.currencyCtrl.SetExchangeRate)
		}

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

		reports := authenticated.Group("/reports")
		{
			reports.GET("/income-expense", r.reportCtrl.IncomeExpense)
			reports.GET("/category", r.reportCtrl.Category)
			reports.GET("/budget", r.reportCtrl.Budget)
			reports.GET("/net-worth", r.reportCtrl.NetWorth)
			reports.GET("/trend", r.reportCtrl.Trend)
			reports.GET("/tag", r.reportCtrl.Tag)
		}

		imports := authenticated.Group("/imports")
		{
			imports.POST("/upload", r.importCtrl.Upload)
			imports.POST("/parse", r.importCtrl.Parse)
			imports.POST("/execute", r.importCtrl.Execute)
		}

		dashboard := authenticated.Group("/dashboard")
		{
			dashboard.GET("", r.dashboardCtrl.Get)
		}

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

		attachments := authenticated.Group("/attachments")
		{
			attachments.POST("/upload", r.attachmentCtrl.Upload)
			attachments.GET("/:id/download", r.attachmentCtrl.Download)
			attachments.GET("/:id/view", r.attachmentCtrl.View)
			attachments.GET("", r.attachmentCtrl.List)
			attachments.DELETE("/:id", r.attachmentCtrl.Delete)
		}

		autocomplete := authenticated.Group("/autocomplete")
		{
			autocomplete.GET("/accounts", r.autocompleteCtrl.Accounts)
			autocomplete.GET("/categories", r.autocompleteCtrl.Categories)
			autocomplete.GET("/tags", r.autocompleteCtrl.Tags)
			autocomplete.GET("/currencies", r.autocompleteCtrl.Currencies)
			autocomplete.GET("/budgets", r.autocompleteCtrl.Budgets)
			autocomplete.GET("/bills", r.autocompleteCtrl.Bills)
		}

		exports := authenticated.Group("/exports")
		{
			exports.GET("/transactions", r.exportCtrl.ExportTransactions)
			exports.GET("/accounts", r.exportCtrl.ExportAccounts)
			exports.GET("/budgets", r.exportCtrl.ExportBudgets)
			exports.GET("/categories", r.exportCtrl.ExportCategories)
			exports.GET("/tags", r.exportCtrl.ExportTags)
		}

		recurringTxns := authenticated.Group("/recurring-transactions")
		{
			recurringTxns.POST("", r.recurringTxnCtrl.Create)
			recurringTxns.GET("", r.recurringTxnCtrl.List)
			recurringTxns.GET("/:id", r.recurringTxnCtrl.Get)
			recurringTxns.PUT("/:id", r.recurringTxnCtrl.Update)
			recurringTxns.DELETE("/:id", r.recurringTxnCtrl.Delete)
			recurringTxns.POST("/process-due", r.recurringTxnCtrl.ProcessDue)
		}

		webhooks := authenticated.Group("/webhooks")
		{
			webhooks.POST("", r.webhookCtrl.Create)
			webhooks.GET("", r.webhookCtrl.List)
			webhooks.GET("/:id", r.webhookCtrl.Get)
			webhooks.PUT("/:id", r.webhookCtrl.Update)
			webhooks.DELETE("/:id", r.webhookCtrl.Delete)
			webhooks.GET("/:id/deliveries", r.webhookCtrl.ListDeliveries)
		}

		objectGroups := authenticated.Group("/object-groups")
		{
			objectGroups.POST("", r.objectGroupCtrl.Create)
			objectGroups.GET("", r.objectGroupCtrl.List)
			objectGroups.GET("/:id", r.objectGroupCtrl.Get)
			objectGroups.PUT("/:id", r.objectGroupCtrl.Update)
			objectGroups.DELETE("/:id", r.objectGroupCtrl.Delete)
		}

		transactionLinks := authenticated.Group("/transaction-links")
		{
			transactionLinks.POST("", r.transactionLinkCtrl.Create)
			transactionLinks.GET("", r.transactionLinkCtrl.List)
			transactionLinks.DELETE("/:id", r.transactionLinkCtrl.Delete)
		}

		preferences := authenticated.Group("/preferences")
		{
			preferences.GET("", r.preferenceCtrl.List)
			preferences.GET("/:key", r.preferenceCtrl.Get)
			preferences.PUT("", r.preferenceCtrl.Set)
			preferences.DELETE("/:key", r.preferenceCtrl.Delete)
		}

		reconciliations := authenticated.Group("/reconciliations")
		{
			reconciliations.POST("", r.reconciliationCtrl.Create)
			reconciliations.GET("", r.reconciliationCtrl.List)
			reconciliations.GET("/:id", r.reconciliationCtrl.Get)
			reconciliations.PUT("/:id", r.reconciliationCtrl.Update)
			reconciliations.DELETE("/:id", r.reconciliationCtrl.Delete)
		}

		// Chart routes - 图表数据API
		chart := authenticated.Group("/chart")
		{
			chart.GET("/account/:id", r.chartCtrl.Account)
			chart.GET("/budget/:id", r.chartCtrl.Budget)
			chart.GET("/category", r.chartCtrl.Category)
			chart.GET("/tag", r.chartCtrl.Tag)
			chart.GET("/transaction", r.chartCtrl.Transaction)
		}

		// Insight routes - 数据洞察API
		insight := authenticated.Group("/insight")
		{
			insight.GET("/expense", r.insightCtrl.Expense)
			insight.GET("/income", r.insightCtrl.Income)
			insight.GET("/transfer", r.insightCtrl.Transfer)
		}

		// Health check (public within authenticated group)
		authenticated.GET("/health", r.healthCheck)
	}
}

func (r *Router) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "healthy"})
}