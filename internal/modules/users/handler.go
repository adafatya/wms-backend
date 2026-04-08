package users

import (
	"net/http"
	"strconv"

	"github.com/adafatya/wms-backend/internal/models"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/users", h.CreateUser)
	router.GET("/users", h.ListUsers)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{
			Message: err.Error(),
		})
		return
	}

	res, err := h.service.CreateUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.StandardResponse{
		Message: "User created successfully",
		Data:    res,
	})
}

func (h *Handler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, pagination, err := h.service.ListUsers(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Message:    "Users fetched successfully",
		Data:       users,
		Pagination: pagination,
	})
}
