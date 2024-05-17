package http

import (
	"mami/e-commerce/config"
	"mami/e-commerce/pkg/middleware"
	"mami/e-commerce/user/repository"
	"mami/e-commerce/user/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Routes(r *gin.RouterGroup, sqlDB config.IDatabase, validator validator.Validate) {
	userRepo := repository.NewUserRepository(sqlDB)
	userService := service.NewUserService(userRepo, validator)
	userHandler := NewUserhandler(userService)

	authMiddleware := middleware.JWTAuth()
	refreshAuthMiddleware := middleware.JWTRefresh()
	authRoute := r.Group("/auth")
	{
		authRoute.POST("/register", userHandler.Register)
		authRoute.POST("/login", userHandler.Login)
		authRoute.POST("/refresh", refreshAuthMiddleware, userHandler.RefreshToken)
		authRoute.GET("/me", authMiddleware, userHandler.GetMe)
		authRoute.PUT("/change-password", authMiddleware, userHandler.ChangePassword)
	}
}
