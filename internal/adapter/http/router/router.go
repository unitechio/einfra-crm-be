package router

import (
	"mymodule/internal/adapter/http/handler"
	"mymodule/internal/adapter/http/middleware"
	"mymodule/internal/config"
	"mymodule/internal/domain"
	"mymodule/internal/logger"
	"mymodule/internal/monitoring"
	"mymodule/internal/realtime"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// SetupRouter sets up the Gin router and registers all routes.
func SetupRouter(cfg *config.Config, log logger.Logger, notificationHandler *handler.NotificationHandler, auditHandler *handler.AuditHandler, auditService domain.AuditService, hub *realtime.Hub, healthHandler *handler.HealthHandler, systemSettingHandler *handler.SystemSettingHandler, featureFlagHandler *handler.FeatureFlagHandler, userSettingsHandler *handler.UserSettingsHandler) *gin.Engine {
	router := gin.Default()

	// --- Prometheus Metrics Setup ---
	reg := prometheus.NewRegistry()
	metrics := monitoring.NewMetrics(reg)
	reg.MustRegister(prometheus.NewGoCollector())

	// Middlewares
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.TimeoutMiddleware(30 * time.Second))
	router.Use(middleware.PrometheusMiddleware(metrics))

	// Health check endpoint
	router.GET("/health", healthHandler.HealthCheck)

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.HandlerFor(reg, promhttp.HandlerOpts{})))

	// Pprof endpoints
	pprof.Register(router)

	// Realtime handler
	handler.NewRealtimeHandler(router, hub)

	// API v1 group
	apiV1 := router.Group("/api/v1")
	apiV1.Use(middleware.AuthMiddleware(cfg.Auth.APIKey))
	apiV1.Use(middleware.AuditLog(auditService))
	{
		// Notification routes
		notificationRoutes := apiV1.Group("/notifications")
		{
			notificationRoutes.GET("", notificationHandler.GetAll)
			// ... other notification routes
		}

		// Audit routes
		auditRoutes := apiV1.Group("/audits")
		{
			// Mark the GET /audits endpoint as deprecated
			auditRoutes.GET("",
				middleware.WithEndpoint("/api/v1/audits"), // For consistent logging/metrics
				middleware.Deprecated(log, middleware.DeprecationOptions{
					Since:       "v1.1.0",
					NewEndpoint: "/api/v2/audits",
				}),
				auditHandler.GetAll,
			)

			auditRoutes.GET("/:id", auditHandler.GetByID)
			auditRoutes.POST("", auditHandler.Log)
			auditRoutes.PUT("/:id", auditHandler.Update)
			auditRoutes.DELETE("/:id", auditHandler.Delete)
		}

		// System Setting routes
		systemSettingRoutes := apiV1.Group("/system-settings")
		{
			systemSettingRoutes.POST("", systemSettingHandler.CreateSystemSetting)
			systemSettingRoutes.GET("", systemSettingHandler.GetAllSystemSettings)
			systemSettingRoutes.GET("/key/:key", systemSettingHandler.GetSystemSettingByKey)
			systemSettingRoutes.GET("/category/:category", systemSettingHandler.GetSystemSettingsByCategory)
			systemSettingRoutes.PUT("", systemSettingHandler.UpdateSystemSetting)
			systemSettingRoutes.DELETE("/:id", systemSettingHandler.DeleteSystemSetting)
		}

		// Feature Flag routes
		featureFlagRoutes := apiV1.Group("/feature-flags")
		{
			featureFlagRoutes.POST("", featureFlagHandler.CreateFeatureFlag)
			featureFlagRoutes.GET("", featureFlagHandler.GetAllFeatureFlags)
			featureFlagRoutes.GET("/name/:name", featureFlagHandler.GetFeatureFlagByName)
			featureFlagRoutes.GET("/category/:category", featureFlagHandler.GetFeatureFlagsByCategory)
			featureFlagRoutes.PUT("", featureFlagHandler.UpdateFeatureFlag)
			featureFlagRoutes.DELETE("/:id", featureFlagHandler.DeleteFeatureFlag)
		}

		// User Settings routes
		userSettingsRoutes := apiV1.Group("/users/:id/settings")
		{
			userSettingsRoutes.PUT("", userSettingsHandler.UpdateSettings)
		}
	}

	return router
}
