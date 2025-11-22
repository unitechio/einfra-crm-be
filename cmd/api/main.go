// @title Ebetool API
// @version 1.0
// @description This is a sample server for a content management system.
// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"net/http"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/unitechio/einfra-be/internal/auth"
	"github.com/unitechio/einfra-be/internal/config"
	"github.com/unitechio/einfra-be/internal/http/router"
	"github.com/unitechio/einfra-be/internal/infrastructure/database"
)

func main() {
	env := os.Getenv("APP_ENV")
	configFile := ".env.dev"
	if env == "production" {
		configFile = ".env.production"
	}
	log.Printf("📦 Loading configuration from %s", configFile)

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	dbConnections, err := database.InitDatabases(cfg.Databases)
	if err != nil {
		log.Fatalf("❌ Failed to initialize databases: %v", err)
	}
	defer database.CloseAll(dbConnections)

	dbProvider := provider.NewDatabaseProvider(dbConnections)

	ocsClient := client.NewOCSClient(string(constants.OCSUrlBase))
	crmClient := client.NewCRMClient(string(constants.CRMUrlBase))

	limiter := usecase.NewPlateChangeLimiter(3, 5*time.Minute)

	// Repositories
	userRepo := repository.NewUserRepository(dbProvider)
	roleRepo := repository.NewRoleRepository(dbProvider)
	permissionRepo := repository.NewPermissionRepository(dbProvider)
	userRoleRepo := repository.NewUserRoleRepository(dbProvider)
	rolePermissionRepo := repository.NewRolePermissionRepository(dbProvider)
	auditRepo := repository.NewAuditLogRepository(dbProvider)
	authRepo := repository.NewAuthRepository(dbProvider)
	parkingRepo := repository.NewParkingRepository(dbConnections)
	sqlRepo := repository.NewSqlRepository(dbConnections)
	imRepo := repository.NewIMRepository(dbProvider)
	customerRepo := repository.NewCustomerRepository(dbConnections, sqlRepo)
	vehicleRepo := repository.NewVehicleRepository(dbConnections, imRepo, customerRepo)
	systemSettingRepo := repository.NewSystemSettingRepository(dbProvider)
	featureFlagRepo := repository.NewFeatureFlagRepository(dbProvider)

	jwtService := auth.NewJWTService(&cfg.JWT)

	// Usecases
	authUsecase := usecase.NewAuthUsecase(dbProvider, userRepo, roleRepo, userRoleRepo, rolePermissionRepo, authRepo, jwtService, &cfg.JWT)
	userUsecase := usecase.NewUserUsecase(userRepo, roleRepo, userRoleRepo, jwtService, dbProvider)
	roleUsecase := usecase.NewRoleUsecase(roleRepo, rolePermissionRepo, dbProvider)
	userRoleUsecase := usecase.NewUserRoleUseCase(dbProvider, userRoleRepo)
	permissionUsecase := usecase.NewPermissionUsecase(permissionRepo)
	rolePermissionUsecase := usecase.NewRolePermissionUseCase(dbProvider, rolePermissionRepo)
	auditUsecase := usecase.NewAuditLogUsecase(auditRepo, dbProvider)
	parkingUsecase := usecase.NewParkingUsecase(parkingRepo)
	sqlUsecase := usecase.NewSqlUsecase(sqlRepo)
	imUsecase := usecase.NewIMUsecase(imRepo, dbProvider)
	ocsUsecase := usecase.NewOCSUsecase(ocsClient, customerRepo)
	vehicleUsecase := usecase.NewVehicleUsecase(vehicleRepo, imRepo, ocsUsecase, dbProvider)
	customerUsecase := usecase.NewCustomerUsecase(customerRepo, crmClient)
	systemSettingUsecase := usecase.NewSystemSettingUsecase(systemSettingRepo)
	featureFlagUsecase := usecase.NewFeatureFlagUsecase(featureFlagRepo)
	operatorUsecase := usecase.NewOperatorUsecase(sqlRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authUsecase)
	userHandler := handler.NewUserHandler(userUsecase, roleUsecase, userRoleUsecase)
	roleHandler := handler.NewRoleHandler(roleUsecase)
	permissionHandler := handler.NewPermissionHandler(permissionUsecase)
	userRoleHandler := handler.NewUserRoleHandler(userRoleUsecase)
	rolePermissionHandler := handler.NewRolePermissionHandler(rolePermissionUsecase)
	auditHandler := handler.NewAuditHandler(auditUsecase)
	customerHandler := handler.NewCustomerHandler(parkingUsecase, customerUsecase)
	imHandler := handler.NewIMHandler(imUsecase)
	vehicleHandler := handler.NewVehicleHandler(vehicleUsecase, limiter)
	ocsHandler := handler.NewOCSHandler(ocsUsecase)
	operatorHandler := handler.NewOperatorHandler(operatorUsecase)
	systemSettingHandler := handler.NewSystemSettingHandler(systemSettingUsecase)
	featureFlagHandler := handler.NewFeatureFlagHandler(featureFlagUsecase)
	sqlHandler := handler.NewSqlHandler(sqlUsecase)

	r := router.InitRouter(cfg,
		customerHandler,
		ocsHandler,
		authHandler,
		userHandler,
		roleHandler,
		permissionHandler,
		userRoleHandler,
		rolePermissionHandler,
		auditHandler,
		sqlHandler,
		auditUsecase,
		imHandler,
		vehicleHandler,
		operatorHandler,
		systemSettingHandler,
		featureFlagHandler,
	)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Custom HTTP server with timeout
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 600 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("🚀 Server starting on %s", srv.Addr)
	log.Printf("📖 Swagger: http://localhost:%d/swagger/index.html", cfg.Server.Port)
	log.Println("✅ Server is ready!")

	// Start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	// Limiter cleanup
	go func() {
		for range time.Tick(time.Minute * 10) {
			limiter.Cleanup()
		}
	}()

	// Graceful shutdown on context cancel (optional)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	<-ctx.Done()
}
