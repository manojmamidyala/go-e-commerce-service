package config

import (
	"fmt"
	"log"
	"net/http"

	responsedto "mami/e-commerce/responseDto"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	engine *gin.Engine
	cfg    *envConfigs
	db     IDatabase
}

func NewServer(db IDatabase) *Server {
	return &Server{
		engine: gin.Default(),
		cfg:    GetConfig(),
		db:     db,
	}
}

func (s Server) GetEngine() *gin.Engine {
	return s.engine
}

func (s Server) Run() error {
	_ = s.engine.SetTrustedProxies(nil)
	if s.cfg.Environment == ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	if err := s.MapRoutes(); err != nil {
		log.Fatalf("MapRoutes Error: %v", err)
	}

	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.engine.GET("/health", func(c *gin.Context) {
		responsedto.JSON(c, http.StatusOK, nil)
		return
	})

	log.Printf("HTTP server is listening on PORT: %v", s.cfg.HttpPort)
	if err := s.engine.Run(fmt.Sprintf(":%d", s.cfg.HttpPort)); err != nil {
		log.Fatalf("Running HTTP server: %v", err)
	}

	return nil
}

func (s Server) MapRoutes() error {
	v1 := s.engine.Group("/api/v1")
	return nil
}
