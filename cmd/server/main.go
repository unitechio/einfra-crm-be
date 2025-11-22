package main

import (
	"context"
	"log"
	"os"

	"github.com/unitechio/einfra-be/internal/auth"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/unitechio/einfra-be/internal/adapter/http/handler"
	"github.com/unitechio/einfra-be/internal/adapter/http/router"
	"github.com/unitechio/einfra-be/internal/auth"
	"github.com/unitechio/einfra-be/internal/config"
	"github.com/unitechio/einfra-be/internal/domain"
	"github.com/unitechio/einfra-be/internal/infrastructure/email"
	oauthinfra "github.com/unitechio/einfra-be/internal/infrastructure/oauth"
	_repository "github.com/unitechio/einfra-be/internal/infrastructure/repository"
	"github.com/unitechio/einfra-be/internal/logger"
	"github.com/unitechio/einfra-be/internal/realtime"
	"github.com/unitechio/einfra-be/internal/usecase"
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
	db, err := pgxpool.New(context.Background(), cfg.DB.PostgresURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	if err := runMigrations(cfg.DB.PostgresURL); err != nil {
		log.Fatalf("Could not run database migrations: %v\n", err)
	}

	// --- Dependency Injection ---

	jwtService := auth.NewJWTService(cfg.Auth)

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

	userRepo := _repository.NewUserRepository(db)
	roleRepo := _repository.NewRoleRepository(db)
	permissionRepo := _repository.NewPermissionRepository(db)
	sessionRepo := _repository.NewSessionRepository(db)
	auditRepo := _repository.NewAuditRepository(db)
	authRepo := _repository.NewAuthRepository(db)
	loginAttemptRepo := _repository.NewPostgresLoginAttemptRepository(db)
	imageRepo := _repository.NewPostgresImageRepository(db)
	tokenRepo := _repository.NewTokenRepository()
	notifRepo := _repository.NewNotificationRepository(db)
	templateRepo := _repository.NewNotificationTemplateRepository(db)
	prefRepo := _repository.NewNotificationPreferenceRepository(db)

	// Messaging and Realtime
	// producer := messaging.NewKafkaProducer(...) // Placeholder
	hub := realtime.NewHub()

	notificationUsecase := usecase.NewNotificationUsecase(
		notifRepo,
		templateRepo,
		prefRepo,
		userRepo,
		nil, // producer placeholder
		hub,
		zapLogger,
	)

	authUsecase := usecase.NewAuthUsecase(
		authRepo,
		userRepo,
		sessionRepo,
		loginAttemptRepo,
		notificationUsecase,
		jwtService,
		cfg.Auth,
	)
	oauthUsecase := usecase.NewOAuthUsecase(userRepo, authRepo, jwtService, oauthProviders)
	emailUsecase := usecase.NewEmailUsecase(emailService)
	userUsecase := usecase.NewUserUsecase(userRepo, authRepo)
	imageUsecase := usecase.NewImageUsecase(imageRepo, cfg.Storage.ImagePath)

	notificationHandler := handler.NewNotificationHandler(notificationUsecase)
	auditService := &usecase.AuditUsecase{}                      // Placeholder
	healthHandler := handler.NewHealthHandler(nil)               // Placeholder
	systemSettingHandler := handler.NewSystemSettingHandler(nil) // Placeholder
	featureFlagHandler := handler.NewFeatureFlagHandler(nil)     // Placeholder
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
