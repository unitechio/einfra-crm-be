
package main

import (
	"context"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"mymodule/internal/adapter/http/handler"
	"mymodule/internal/adapter/http/router"
	"mymodule/internal/config"
	"mymodule/internal/domain"
	"mymodule/internal/infrastructure/email"
	oauthinfra "mymodule/internal/infrastructure/oauth"
	"mymodule/internal/infrastructure/repository"
	"mymodule/internal/logger"
	"mymodule/internal/realtime"
	"mymodule/internal/usecase"
)

func main() {
	// --- Configuration ---
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}

	// --- Logger ---
	zapLogger := logger.NewZapLogger(cfg.Log.Level, cfg.Log.File)

	// --- Database Connection & Migration ---
	dbpool, err := pgxpool.New(context.Background(), cfg.DB.PostgresURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	if err := runMigrations(cfg.DB.PostgresURL); err != nil {
		log.Fatalf("Could not run database migrations: %v\n", err)
	}

	// --- Dependency Injection ---
	userRepo := repository.NewPostgresUserRepository(dbpool)
	imageRepo := repository.NewPostgresImageRepository(dbpool)
	tokenRepo := repository.NewTokenRepository()
	googleProvider := oauthinfra.NewGoogleProvider(&cfg.Auth.Google)
	azureProvider := oauthinfra.NewAzureProvider(&cfg.Auth.Azure)
	oauthProviders := map[string]domain.OAuthProvider{
		"google": googleProvider,
		"azure":  azureProvider,
	}
	emailService, err := email.NewSmtpService(&cfg.SMTP)
	if err != nil {
		log.Fatalf("failed to create email service: %v", err)
	}
	authUsecase := usecase.NewAuthUsecase(userRepo, tokenRepo)
	oauthUsecase := usecase.NewOAuthUsecase(userRepo, tokenRepo, oauthProviders)
	emailUsecase := usecase.NewEmailUsecase(emailService)
	userUsecase := usecase.NewUserUsecase(userRepo)
	imageUsecase := usecase.NewImageUsecase(imageRepo, cfg.Storage.ImagePath)
	notificationHandler := handler.NewNotificationHandler(nil) // Placeholder
	auditService := &usecase.AuditUsecase{}                      // Placeholder
	hub := realtime.NewHub()
	healthHandler := handler.NewHealthHandler(nil)                                     // Placeholder
	systemSettingHandler := handler.NewSystemSettingHandler(nil)                     // Placeholder
	featureFlagHandler := handler.NewFeatureFlagHandler(nil)                         // Placeholder
	userSettingsHandler := handler.NewUserSettingsHandler(userUsecase)

	// --- HTTP Server Setup ---
	ginRouter := router.SetupRouter(cfg, zapLogger, notificationHandler, &handler.AuditHandler{}, auditService, hub, healthHandler, systemSettingHandler, featureFlagHandler, userSettingsHandler)

	// --- Start Server ---
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := ginRouter.Run(cfg.Server.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

// runMigrations executes the database migrations.
func runMigrations(databaseURL string) error {
	migrationPath := "file://internal/infrastructure/repository/migrations"
	m, err := migrate.New(migrationPath, databaseURL)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	log.Println("Database migrated successfully")
	return nil
}
