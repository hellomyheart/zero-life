// @title           Zero-Life API
// @version         1.0
// @description     Zero-Life 财务管理系统后端API，提供账户、交易、分类、标签、预算、账单、货币、规则、报表等功能
// @termsOfService  http://swagger.io/terms/

// @contact.name   Zero-Life
// @contact.url    https://github.com/hellomyheart/zero-life

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/hellomyheart/zero-life/server/internal/config"
	"github.com/hellomyheart/zero-life/server/internal/controller"
	"github.com/hellomyheart/zero-life/server/internal/model"
	"github.com/hellomyheart/zero-life/server/internal/pkg/jwt"
	"github.com/hellomyheart/zero-life/server/internal/repository"
	"github.com/hellomyheart/zero-life/server/internal/router"
	"github.com/hellomyheart/zero-life/server/internal/service"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// main 是程序入口，按顺序完成所有依赖初始化后启动HTTP服务
func main() {
	// 加载配置文件（config.yaml）
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 根据环境初始化日志：生产环境用JSON格式，开发环境用控制台格式
	var logger *zap.Logger
	var err error
	if config.C.App.Env == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// 初始化SQLite数据库，确保数据库目录存在
	dbPath := config.C.DB.Path
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		logger.Fatal("Failed to create database directory", zap.Error(err))
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// 设置数据库连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("Failed to get database instance", zap.Error(err))
	}
	sqlDB.SetMaxIdleConns(config.C.DB.MaxIdleConns) // 最大空闲连接数
	sqlDB.SetMaxOpenConns(config.C.DB.MaxOpenConns) // 最大打开连接数
	sqlDB.SetConnMaxLifetime(config.C.DB.MaxLifetime) // 连接最大存活时间

	// SQLite优化设置：WAL模式提升并发读写性能，开启外键约束
	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA foreign_keys=ON")

	// 自动迁移：根据模型定义自动创建/更新数据库表结构
	if err := db.AutoMigrate(
		&model.User{},
		&model.Account{},
		&model.Transaction{},
		&model.TransactionTag{},
		&model.Category{},
		&model.Tag{},
		&model.Budget{},
		&model.BudgetCategory{},
		&model.BudgetHistory{},
		&model.Bill{},
		&model.Currency{},
		&model.ExchangeRate{},
		&model.RuleGroup{},
		&model.Rule{},
		&model.RuleCondition{},
		&model.RuleAction{},
		&model.PiggyBank{},
		&model.PiggyEvent{},
		&model.Attachment{},
		&model.Recurrence{},
		&model.RecurringTransaction{},
		&model.RecurringTransactionLog{},
		&model.Webhook{},
		&model.WebhookMessage{},
		&model.WebhookDelivery{},
		&model.Reconciliation{},
		&model.TransactionReconciliation{},
		&model.ReconciliationEntry{},
		&model.LinkType{},
		&model.TransactionJournalLink{},
		&model.Preference{},
		&model.ObjectGroup{},
		&model.Configuration{},
		&model.BackupCode{},
		&model.KVStore{},
	); err != nil {
		logger.Fatal("Failed to auto migrate", zap.Error(err))
	}

	// 初始化附件存储目录
	attachPath := config.C.Attach.Path
	if attachPath == "" {
		attachPath = "./data/attachments"
	}
	if err := os.MkdirAll(attachPath, 0755); err != nil {
		logger.Fatal("Failed to create attachment directory", zap.Error(err))
	}

	// 初始化键值存储仓库（替代 Redis，用于登录限制、密码重置令牌、请求限流）
	kvRepo := repository.NewKVRepository(db)

	// 初始化JWT服务，用于生成和验证访问令牌
	jwtService := jwt.NewService()

	// 初始化数据访问层（Repository），每个Repository对应一个数据库表的操作
	authRepo := repository.NewAuthRepository(db)
	userRepo := repository.NewUserRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	txnRepo := repository.NewTransactionRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	tagRepo := repository.NewTagRepository(db)
	budgetRepo := repository.NewBudgetRepository(db)
	billRepo := repository.NewBillRepository(db)
	currencyRepo := repository.NewCurrencyRepository(db)
	ruleRepo := repository.NewRuleRepository(db)
	ruleGroupRepo := repository.NewRuleGroupRepository(db)
	piggyBankRepo := repository.NewPiggyBankRepository(db)
	attachmentRepo := repository.NewAttachmentRepository(db)
	recurrenceRepo := repository.NewRecurrenceRepository(db)
	webhookRepo := repository.NewWebhookRepository(db)
	reconRepo := repository.NewReconciliationRepository(db)
	linkTypeRepo := repository.NewLinkTypeRepository(db)
	txnLinkRepo := repository.NewTransactionLinkRepository(db)
	prefRepo := repository.NewPreferenceRepository(db)
	rtRepo := repository.NewRecurringTransactionRepository(db)
	ogRepo := repository.NewObjectGroupRepository(db)
	configRepo := repository.NewConfigurationRepository(db)

	// 新增MFA服务（需在authService之前初始化，因为authService依赖mfaService）
	mfaService := service.NewMFAService(userRepo, db)
	// 初始化业务逻辑层（Service），Service组合Repository实现业务逻辑
	// 依赖注入：Service通过构造函数接收所需的Repository和其他Service
	authService := service.NewAuthService(authRepo, jwtService, kvRepo, mfaService)
	accountService := service.NewAccountService(accountRepo)
	ruleService := service.NewRuleService(ruleRepo, txnRepo, categoryRepo, budgetRepo, tagRepo)
	webhookService := service.NewWebhookService(webhookRepo)
	txnService := service.NewTransactionService(txnRepo, accountRepo, categoryRepo, tagRepo, db, ruleService, webhookService)
	categoryService := service.NewCategoryService(categoryRepo, db)
	tagService := service.NewTagService(tagRepo, db)
	budgetService := service.NewBudgetService(budgetRepo, txnRepo, categoryRepo)
	billService := service.NewBillService(billRepo, txnService)
	currencyService := service.NewCurrencyService(currencyRepo, accountRepo)
	ruleGroupService := service.NewRuleGroupService(ruleGroupRepo, ruleRepo, txnRepo, categoryRepo, tagRepo, budgetRepo)
	reportService := service.NewReportService(txnRepo, accountRepo, budgetRepo, categoryRepo, tagRepo)
	dashboardService := service.NewDashboardService(txnRepo, accountRepo, budgetRepo, billRepo, categoryRepo)
	importService := service.NewImportService(txnService, accountRepo, db)
	piggyBankService := service.NewPiggyBankService(piggyBankRepo, accountRepo)
	attachmentService := service.NewAttachmentService(attachmentRepo, attachPath)
	exportService := service.NewExportService(txnRepo, accountRepo, billRepo, budgetRepo, categoryRepo, tagRepo, piggyBankRepo, ruleRepo)
	recurrenceService := service.NewRecurrenceService(recurrenceRepo, txnService, accountRepo, tagRepo)
	cronService := service.NewCronService(recurrenceService, billService, billRepo, recurrenceRepo, txnService, budgetService, logger)
	reconService := service.NewReconciliationService(reconRepo, accountRepo, txnRepo)
	txnBulkService := service.NewTransactionBulkService(txnRepo, accountRepo, db)
	linkTypeService := service.NewLinkTypeService(linkTypeRepo)
	txnLinkService := service.NewTransactionLinkService(txnLinkRepo, txnRepo)
	prefService := service.NewPreferenceService(prefRepo)
	rtService := service.NewRecurringTransactionService(rtRepo, txnRepo, txnService, accountRepo, db)
	ogService := service.NewObjectGroupService(ogRepo)
	// 新增图表和洞察服务
	chartService := service.NewChartService(txnRepo, accountRepo, budgetRepo, categoryRepo, tagRepo)
	insightService := service.NewInsightService(txnRepo, accountRepo, categoryRepo)
	// 新增用户管理服务
	userService := service.NewUserService(userRepo, db)

	// 初始化默认货币数据（如CNY、USD等）
	if err := currencyService.InitDefaultCurrencies(); err != nil {
		logger.Warn("Failed to initialize default currencies", zap.Error(err))
	}

	// 初始化控制器层（Controller），Controller接收HTTP请求并调用Service处理
	authCtrl := controller.NewAuthController(authService)
	accountCtrl := controller.NewAccountController(accountService)
	txnCtrl := controller.NewTransactionController(txnService)
	categoryCtrl := controller.NewCategoryController(categoryService)
	tagCtrl := controller.NewTagController(tagService)
	budgetCtrl := controller.NewBudgetController(budgetService)
	billCtrl := controller.NewBillController(billService)
	currencyCtrl := controller.NewCurrencyController(currencyService)
	ruleCtrl := controller.NewRuleController(ruleService)
	reportCtrl := controller.NewReportController(reportService)
	dashboardCtrl := controller.NewDashboardController(dashboardService)
	importCtrl := controller.NewImportController(importService)
	piggyBankCtrl := controller.NewPiggyBankController(piggyBankService)
	attachmentCtrl := controller.NewAttachmentController(attachmentService)
	autocompleteCtrl := controller.NewAutocompleteController(accountRepo, categoryRepo, tagRepo, currencyRepo, budgetRepo, billRepo)
	exportCtrl := controller.NewExportController(exportService)
	webhookCtrl := controller.NewWebhookController(webhookService)
	reconCtrl := controller.NewReconciliationController(reconService)
	txnLinkCtrl := controller.NewTransactionLinkController(txnLinkService)
	prefCtrl := controller.NewPreferenceController(prefService)
	rtCtrl := controller.NewRecurringTransactionController(rtService)
	ogCtrl := controller.NewObjectGroupController(ogService)
	chartCtrl := controller.NewChartController(chartService)
	insightCtrl := controller.NewInsightController(insightService)
	mfaCtrl := controller.NewMFAController(mfaService)
	userCtrl := controller.NewUserController(userService)
	linkTypeCtrl := controller.NewLinkTypeController(linkTypeService)
	adminService := service.NewAdminService(userRepo, configRepo)
	adminCtrl := controller.NewAdminController(adminService)
	adminUserCtrl := controller.NewAdminUserController(adminService)
	txnBulkCtrl := controller.NewTransactionBulkController(txnBulkService)
	recurrenceCtrl := controller.NewRecurrenceController(recurrenceService)
	ruleGroupCtrl := controller.NewRuleGroupController(ruleGroupService)
	cronCtrl := controller.NewCronController(cronService)

	// 初始化Gin引擎，生产环境使用Release模式减少日志输出
	if config.C.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())

	// 设置路由：注册所有API端点和中间件
	r := router.NewRouter(
		engine,
		authCtrl,
		accountCtrl,
		txnCtrl,
		categoryCtrl,
		tagCtrl,
		budgetCtrl,
		billCtrl,
		currencyCtrl,
		ruleCtrl,
		reportCtrl,
		dashboardCtrl,
		importCtrl,
		piggyBankCtrl,
		attachmentCtrl,
		autocompleteCtrl,
		exportCtrl,
		rtCtrl,
		webhookCtrl,
		ogCtrl,
		txnLinkCtrl,
		prefCtrl,
		reconCtrl,
		chartCtrl,
		insightCtrl,
		mfaCtrl,
		userCtrl,
		linkTypeCtrl,
		adminCtrl,
		adminUserCtrl,
		txnBulkCtrl,
		recurrenceCtrl,
		ruleGroupCtrl,
		cronCtrl,
	)
	r.Setup(jwtService, db)

	cronService.StartScheduler()

	addr := ":" + config.C.App.Port
	logger.Info("Server starting", zap.String("addr", addr))
	if err := engine.Run(addr); err != nil {
		cronService.StopScheduler()
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
