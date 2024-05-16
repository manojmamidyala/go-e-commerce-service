package http

import (
	"mami/e-commerce/config"
	"mami/e-commerce/user/repository"
	"mami/e-commerce/user/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Routes(r *gin.RouterGroup, sqlDB config.IDatabase, validator validator.Validate) {
	userRepo := repository.NewUserRepository(sqlDB)
	userService := service.NewUserService(userRepo, validator)
	userHandler := NewUserhandler(userService)

	authRoute := r.Group("/auth")
	{
		authRoute.POST("/login", userHandler.Login)
	}

}
