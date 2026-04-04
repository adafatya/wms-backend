package server

import (
	"github.com/adafatya/wms-backend/internal/db/sqlc"
	"github.com/adafatya/wms-backend/internal/modules/users"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	store  sqlc.Store
	router *gin.Engine
	logger *zap.Logger
}

func NewServer(store sqlc.Store, logger *zap.Logger) *Server {
	server := &Server{
		store:  store,
		logger: logger,
	}
	router := gin.Default()

	// Initial repository, service, and handler for users
	userRepo := users.NewRepository(store)
	userService := users.NewService(userRepo)
	userHandler := users.NewHandler(userService)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// User routes
	userHandler.RegisterRoutes(router)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
