package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/adafatya/wms-backend/db/sqlc"
	"go.uber.org/zap"
)

type Server struct {
	store  db.Store
	router *gin.Engine
	logger *zap.Logger
}

func NewServer(store db.Store, logger *zap.Logger) *Server {
	server := &Server{
		store:  store,
		logger: logger,
	}
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
