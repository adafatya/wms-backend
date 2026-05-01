package server

import (
	"github.com/adafatya/wms-backend/internal/db/sqlc"
	"github.com/adafatya/wms-backend/internal/models"
	"github.com/adafatya/wms-backend/internal/modules/roles"
	"github.com/adafatya/wms-backend/internal/modules/users"
	"github.com/adafatya/wms-backend/internal/modules/inventory"
	"github.com/adafatya/wms-backend/internal/modules/inboundproduct"
	"github.com/adafatya/wms-backend/internal/modules/outboundproduct"
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

	// Initial repository, service, and handler for roles
	roleRepo := roles.NewRepository(store)
	roleService := roles.NewService(roleRepo)
	roleHandler := roles.NewHandler(roleService)

	// Initial repository, service, and handler for users
	userRepo := users.NewRepository(store)
	userService := users.NewService(userRepo, roleRepo)
	userHandler := users.NewHandler(userService)

	// Initial repository, service, and handler for inventory
	inventoryRepo := inventory.NewRepository(store)
	inventoryService := inventory.NewService(inventoryRepo, store)
	inventoryHandler := inventory.NewHandler(inventoryService)
	
	// Initial repository, service, and handler for inbound product
	inboundRepo := inboundproduct.NewRepository(store)
	inboundService := inboundproduct.NewService(inboundRepo, store)
	inboundHandler := inboundproduct.NewHandler(inboundService)

	// Initial repository, service, and handler for outbound product
	outboundRepo := outboundproduct.NewRepository(store)
	outboundService := outboundproduct.NewService(outboundRepo, store)
	outboundHandler := outboundproduct.NewHandler(outboundService)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, models.StandardResponse{
			Message: "pong",
		})
	})

	// Register routes
	roleHandler.RegisterRoutes(router)
	userHandler.RegisterRoutes(router)
	inventoryHandler.RegisterRoutes(router)
	inboundHandler.RegisterRoutes(router)
	outboundHandler.RegisterRoutes(router)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
