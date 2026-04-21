package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/zero-life/server/internal/config"
	"github.com/zero-life/server/internal/controller"
	"github.com/zero-life/server/internal/model"
	"github.com/zero-life/server/internal/pkg/jwt"
	"github.com/zero-life/server/internal/repository"
	"github.com/zero-life/server/internal/router"
	"github.com/zero-life/server/internal/service"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Load config
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
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

	// Initialize database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.C.DB.User,
		config.C.DB.Password,
		config.C.DB.Host,
		config.C.DB.Port,
		config.C.DB.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("Failed to get database instance", zap.Error(err))
	}
	sqlDB.SetMaxIdleConns(config.C.DB.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.C.DB.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.C.DB.MaxLifetime)

	// Auto migrate
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
		&model.Rule{},
		&model.RuleCondition{},
		&model.RuleAction{},
	); err != nil {
		logger.Fatal("Failed to auto migrate", zap.Error(err))
	}

	// Initialize Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.C.Redis.Host, config.C.Redis.Port),
		Password: config.C.Redis.Password,
		DB:       config.C.Redis.DB,
	})

	// Initialize JWT service
	jwtService := jwt.NewService()

	// Initialize repositories
	authRepo := repository.NewAuthRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	txnRepo := repository.NewTransactionRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	tagRepo := repository.NewTagRepository(db)
	budgetRepo := repository.NewBudgetRepository(db)
	billRepo := repository.NewBillRepository(db)
	currencyRepo := repository.NewCurrencyRepository(db)
	ruleRepo := repository.NewRuleRepository(db)

	// Initialize services
	authService := service.NewAuthService(authRepo, jwtService, rdb)
	accountService := service.NewAccountService(accountRepo)
	txnService := service.NewTransactionService(txnRepo, accountRepo, db)
	categoryService := service.NewCategoryService(categoryRepo, db)
	tagService := service.NewTagService(tagRepo)
	budgetService := service.NewBudgetService(budgetRepo, txnRepo)
	billService := service.NewBillService(billRepo)
	currencyService := service.NewCurrencyService(currencyRepo, accountRepo)
	ruleService := service.NewRuleService(ruleRepo)
	reportService := service.NewReportService(txnRepo, accountRepo, budgetRepo, categoryRepo)
	dashboardService := service.NewDashboardService(txnRepo, accountRepo, budgetRepo, billRepo)
	importService := service.NewImportService(txnRepo, accountRepo, db)

	// Initialize default currencies
	if err := currencyService.InitDefaultCurrencies(); err != nil {
		logger.Warn("Failed to initialize default currencies", zap.Error(err))
	}

	// Initialize controllers
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

	// Initialize Gin engine
	if config.C.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())

	// Setup router
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
	)
	r.Setup(jwtService)

	// Start server
	addr := ":" + config.C.App.Port
	logger.Info("Server starting", zap.String("addr", addr))
	if err := engine.Run(addr); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
