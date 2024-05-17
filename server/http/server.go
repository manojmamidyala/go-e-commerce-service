package http

import (
	"fmt"

	"mami/e-commerce/commons/logger"
	"mami/e-commerce/config"
	healthHttp "mami/e-commerce/health/http"
	userHttp "mami/e-commerce/user/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	engine    *gin.Engine
	cfg       *config.EnvConfigs
	db        config.IDatabase
	validator *validator.Validate
}

func NewServer(db config.IDatabase, validator *validator.Validate) *Server {
	return &Server{
		engine:    gin.Default(),
		cfg:       config.GetConfig(),
		db:        db,
		validator: validator,
	}
}

func (s Server) GetEngine() *gin.Engine {
	return s.engine
}

func (s Server) Run() error {
	_ = s.engine.SetTrustedProxies(nil)
	if s.cfg.Environment == config.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	if err := s.MapRoutes(); err != nil {
		logger.Fatalf("MapRoutes Error: %v", err)
	}

	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logger.Infof("HTTP server is listening on PORT: %v", s.cfg.HttpPort)
	if err := s.engine.Run(fmt.Sprintf(":%d", s.cfg.HttpPort)); err != nil {
		logger.Fatalf("Running HTTP server: %v", err)
	}

	return nil
}

func (s Server) MapRoutes() error {
	v1 := s.engine.Group("/api/v1")
	healthHttp.Routes(v1, s.db)
	userHttp.Routes(v1, s.db, *s.validator)
	return nil
}
