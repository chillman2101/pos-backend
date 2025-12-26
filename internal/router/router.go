package router

import (
	"pos-backend/internal/config"
	"pos-backend/internal/handler"
	"pos-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, authHandler *handler.AuthHandler, userHandler *handler.UserHandler, categoryHandler *handler.CategoryHandler, productHandler *handler.ProductHandler, transactionHandler *handler.TransactionHandler) *gin.Engine {
	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Global middlewares
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CORSMiddleware())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"message": "POS Backend API is running",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Public routes (no auth required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes (auth required)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// Example protected route
			protected.GET("/profile", func(c *gin.Context) {
				userID := c.GetString("user_id")
				username := c.GetString("username")
				role := c.GetString("role")

				c.JSON(200, gin.H{
					"user_id":  userID,
					"username": username,
					"role":     role,
				})
			})

			// Users routes
			users := protected.Group("/users")
			{
				users.GET("", userHandler.GetAll)
				users.GET("/:id", userHandler.GetByID)
				users.PUT("/:id", userHandler.Update)
				users.DELETE("/:id", middleware.RoleMiddleware("admin"), userHandler.Delete)
			}

			// Categories routes
			categories := protected.Group("/categories")
			{
				categories.GET("", categoryHandler.GetAll)
				categories.GET("/:id", categoryHandler.GetByID)
				categories.POST("", middleware.RoleMiddleware("admin", "manager"), categoryHandler.Create)
				categories.PUT("/:id", middleware.RoleMiddleware("admin", "manager"), categoryHandler.Update)
				categories.DELETE("/:id", middleware.RoleMiddleware("admin"), categoryHandler.Delete)
			}

			// Products routes
			products := protected.Group("/products")
			{
				products.GET("", productHandler.GetAll)
				products.GET("/category/:category_id", productHandler.GetByCategory)
				products.GET("/sku/:sku", productHandler.GetBySKU)
				products.GET("/:id", productHandler.GetByID)
				products.POST("", middleware.RoleMiddleware("admin", "manager"), productHandler.Create)
				products.PUT("/:id", middleware.RoleMiddleware("admin", "manager"), productHandler.Update)
				products.DELETE("/:id", middleware.RoleMiddleware("admin"), productHandler.Delete)
			}

			// Transactions routes
			transactions := protected.Group("/transactions")
			{
				transactions.GET("", transactionHandler.GetAll)
				transactions.GET("/:id", transactionHandler.GetByID)
				transactions.POST("", transactionHandler.Create)
				transactions.POST("/bulk-sync", transactionHandler.BulkSync)
				transactions.PATCH("/:id/cancel", middleware.RoleMiddleware("admin", "manager"), transactionHandler.Cancel)
			}

			// Inventory routes
			// inventory := protected.Group("/inventory")
			// {
			// 	inventory.GET("/movements", inventoryHandler.GetMovements)
			// 	inventory.POST("/adjustment", middleware.RoleMiddleware("admin", "manager"), inventoryHandler.Adjustment)
			// }
		}
	}

	return r
}
