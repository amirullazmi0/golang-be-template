package handler

import (
	"github.com/amirullazmi0/kratify-backend/config"
	"github.com/amirullazmi0/kratify-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, userHandler *UserHandler, addressHandler *AddressHandler, cfg *config.Config) {
	// API routes
	api := router.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.POST("/refresh", userHandler.RefreshToken)
			auth.POST("/logout", middleware.JWTAuth(&cfg.JWT), userHandler.Logout)
		}

		// User routes (protected)
		users := api.Group("/users")
		users.Use(middleware.JWTAuth(&cfg.JWT))
		{
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
			users.PUT("/change-password", userHandler.ChangePassword)

			// Admin only routes
			users.GET("", middleware.RequireRole("ADMIN"), userHandler.GetAllUsers)
			users.DELETE("/:id", middleware.RequireSuperAdmin(), userHandler.DeleteUser)
		}

		// Address routes (protected)
		addresses := api.Group("/addresses")
		addresses.Use(middleware.JWTAuth(&cfg.JWT))
		{
			addresses.GET("", addressHandler.GetAddressByAuth)
			addresses.POST("", addressHandler.CreateAddress)
			addresses.GET("/:id", addressHandler.GetAddressByID)
			addresses.PUT("/:id", addressHandler.UpdateAddress)
			addresses.DELETE("/:id", addressHandler.DeleteAddress)
		}
	}
}
