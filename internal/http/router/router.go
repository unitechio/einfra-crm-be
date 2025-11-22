package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/unitechio/einfra-be/internal/adapter/http/handler"
	"github.com/unitechio/einfra-be/internal/adapter/http/middleware"
	"github.com/unitechio/einfra-be/internal/config"
	"github.com/unitechio/einfra-be/internal/domain"
)

// InitRouter initializes the Gin router with all routes and middleware
func InitRouter(
	cfg *config.Config,
	customerHandler *handler.CustomerHandler,
	ocsHandler *handler.OCSHandler,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	roleHandler *handler.RoleHandler,
	permissionHandler *handler.PermissionHandler,
	userRoleHandler *handler.UserRoleHandler,
	rolePermissionHandler *handler.RolePermissionHandler,
	auditHandler *handler.AuditHandler,
	sqlHandler *handler.SqlHandler,
	auditUsecase domain.AuditUsecase,
	imHandler *handler.IMHandler,
	vehicleHandler *handler.VehicleHandler,
	operatorHandler *handler.OperatorHandler,
	systemSettingHandler *handler.SystemSettingHandler,
	featureFlagHandler *handler.FeatureFlagHandler,
) *gin.Engine {
	if !cfg.Server.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Global Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.TimeoutMiddleware(30 * time.Second))

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Configure appropriately for production
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": cfg.Server.Name,
			"env":     cfg.Server.Env,
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// API v1 Group
	v1 := r.Group("/api/v1")

	// Audit Log Middleware for all v1 routes (except auth maybe?)
	// v1.Use(middleware.AuditLog(auditUsecase))

	// Auth Routes
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.POST("/verify-email", authHandler.VerifyEmail)
	}

	// Protected Routes
	protected := v1.Group("")
	// protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret)) // Uncomment when middleware is ready
	{
		// User Routes
		users := protected.Group("/users")
		{
			users.POST("", userHandler.Create)
			users.GET("", userHandler.List)
			users.GET("/:id", userHandler.Get)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
			users.PUT("/:id/password", userHandler.ChangePassword)
			users.PUT("/:id/settings", userHandler.UpdateSettings)
		}

		// Role Routes
		roles := protected.Group("/roles")
		{
			roles.POST("", roleHandler.Create)
			roles.GET("", roleHandler.List)
			roles.GET("/:id", roleHandler.Get)
			roles.PUT("/:id", roleHandler.Update)
			roles.DELETE("/:id", roleHandler.Delete)
		}

		// Permission Routes
		permissions := protected.Group("/permissions")
		{
			permissions.POST("", permissionHandler.Create)
			permissions.GET("", permissionHandler.List)
			permissions.GET("/:id", permissionHandler.Get)
			permissions.PUT("/:id", permissionHandler.Update)
			permissions.DELETE("/:id", permissionHandler.Delete)
		}

		// Audit Routes
		audits := protected.Group("/audits")
		{
			audits.GET("", auditHandler.List)
			audits.GET("/:id", auditHandler.Get)
			audits.GET("/stats", auditHandler.GetStatistics)
		}

		// System Settings
		settings := protected.Group("/system-settings")
		{
			settings.GET("", systemSettingHandler.List)
			settings.POST("", systemSettingHandler.Create)
			settings.PUT("/:id", systemSettingHandler.Update)
			settings.DELETE("/:id", systemSettingHandler.Delete)
		}

		// Feature Flags
		flags := protected.Group("/feature-flags")
		{
			flags.GET("", featureFlagHandler.List)
			flags.POST("", featureFlagHandler.Create)
			flags.PUT("/:id", featureFlagHandler.Update)
			flags.DELETE("/:id", featureFlagHandler.Delete)
		}

		// Other routes (placeholders as I don't have handler methods)
		// customer, ocs, im, vehicle, operator, sql
	}

	return r
}
