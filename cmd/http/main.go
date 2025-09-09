
package main

import (
	"context"
	"fmt"
	"mymodule/internal/adapter/http/handler"
	"mymodule/internal/adapter/http/router"
	"mymodule/internal/config"
	"mymodule/internal/domain"
	"mymodule/internal/infrastructure"
	"mymodule/internal/logger"
	"mymodule/internal/realtime"
	"mymodule/internal/trace"
	"mymodule/internal/usecase"

	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// @title Notification API
// @version 1.0
// @description This is a sample server for a notification API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("could not load config: %v", err))
	}

	// Initialize logger
	log := logger.NewZapLogger(logger.LoggerConfig{
		Level: logger.LogLevel(cfg.Logger.Level),
		DevMode: true, // Assuming dev mode for now
		ServiceName: "notification-service",
		ServiceVersion: "1.0.0",
	})

	ctx := context.Background()

	cleanup := trace.InitTracer()
	defer cleanup()

	hub := realtime.NewHub()
	go hub.Run()

	// Database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(ctx, "failed to connect database", logger.LogField{Key: "error", Value: err})
	}

	// Migrate the schema
	db.AutoMigrate(&domain.Notification{}, &domain.AuditLog{})

	// Create Kafka Producer
	kafkaProducer, err := infrastructure.NewKafkaProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	if err != nil {
		log.Fatal(ctx, "Failed to create kafka producer", logger.LogField{Key: "error", Value: err})
	}
	defer kafkaProducer.Close()

	// Create repositories, usecases and handlers
	notificationRepo := infrastructure.NewNotificationRepository(db)
	notificationUsecase := usecase.NewNotificationUsecase(notificationRepo, hub, kafkaProducer, log)
	notificationHandler := handler.NewNotificationHandler(notificationUsecase, log)

	auditRepo := infrastructure.NewAuditRepository(db)
	auditService := usecase.NewAuditService(auditRepo)
	auditHandler := handler.NewAuditHandler(auditService, log)

	healthHandler := handler.NewHealthHandler()

	// Setup router
	router := router.SetupRouter(cfg, notificationHandler, auditHandler, auditService, hub, healthHandler)

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Info(ctx, fmt.Sprintf("Server starting on port %s", cfg.Server.Port))
	if err := router.Run(cfg.Server.Port); err != nil {
		log.Fatal(ctx, "failed to start server", logger.LogField{Key: "error", Value: err})
	}
}
