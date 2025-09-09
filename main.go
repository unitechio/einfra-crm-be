// @title My API
// @version 1.0
// @description This is a sample server for a Go-based application.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https

package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"mymodule/internal/adapter/http/handler"
	"mymodule/internal/adapter/http/middleware"
	"mymodule/internal/config"
	"mymodule/internal/domain"
	"mymodule/internal/infrastructure"
	"mymodule/internal/infrastructure/drivers"
	"mymodule/internal/usecase"
)

func main() {
	// Load configuration.
	cfg := config.LoadConfig()

	// Initialize database connections.
	db, err := drivers.NewMySQLConnection(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	redisClient, err := drivers.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	minioClient, err := drivers.NewMinioClient(cfg.Minio)
	if err != nil {
		log.Fatalf("Failed to connect to Minio: %v", err)
	}

	// Set up the router.
	r := gin.New()

	// Add middleware.
	r.Use(gin.Logger()) // Use Gin's default logger.
	r.Use(middleware.RecoveryMiddleware()) // Use our custom recovery middleware.
	r.Use(middleware.CorsMiddleware())     // Use our CORS middleware.

	// Set up rate limiting.
	rateLimiter := middleware.NewIPRateLimiter(cfg.RateLimit)
	r.Use(rateLimiter.RateLimitMiddleware())

	// Set up repositories.
	var authRepo domain.AuthRepository
	switch cfg.AuthProvider {
	case "keycloak":
		authRepo = infrastructure.NewKeycloakAuthRepository(cfg.Keycloak)
		log.Println("Using Keycloak for authentication")
	default:
		authRepo = infrastructure.NewJWTAuthRepository(cfg.JWT)
		log.Println("Using JWT for authentication")
	}
	healthRepo := infrastructure.NewHealthRepository(db, redisClient, minioClient)

	// Set up use cases.
	authUsecase := usecase.NewAuthUsecase(authRepo)
	healthUsecase := usecase.NewHealthUsecase(healthRepo)

	// Set up handlers.
	authHandler := handler.NewAuthHandler(authUsecase)
	healthHandler := handler.NewHealthHandler(healthUsecase)

	// Set up routes.
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", healthHandler.Check)
	api := r.Group("/api")
	{
		api.POST("/login", authHandler.Login)

		// Protected route with permission middleware.
		authRequired := api.Group("/")
		authRequired.Use(middleware.AuthMiddleware(authRepo))
		{
			authRequired.GET("/protected", authHandler.ProtectedData)

			// Admin-only route.
			adminRequired := authRequired.Group("/admin")
			adminRequired.Use(middleware.PermissionMiddleware("admin"))
			{
				adminRequired.GET("/data", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "This is admin data"})
				})
			}
		}
	}

	// Start the server.
	log.Printf("Server listening on %s", cfg.Server.Addr)
	if err := r.Run(cfg.Server.Addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
