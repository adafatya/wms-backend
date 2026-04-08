package roles

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
	router.POST("/roles", h.CreateRole)
	router.GET("/roles", h.ListRoles)
	router.GET("/roles/:id", h.GetRole)
	router.PUT("/roles/:id", h.UpdateRole)
	router.DELETE("/roles/:id", h.DeleteRole)
}

func (h *Handler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{
			Message: err.Error(),
		})
		return
	}

	role, err := h.service.CreateRole(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.StandardResponse{
		Message: "Role created successfully",
		Data:    role,
	})
}

func (h *Handler) GetRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{
			Message: "invalid role id",
		})
		return
	}

	role, err := h.service.GetRole(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.StandardResponse{
			Message: "role not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Message: "Role fetched successfully",
		Data:    role,
	})
}

func (h *Handler) ListRoles(c *gin.Context) {
	page, limit, err := models.ParsePaginationQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: err.Error()})
		return
	}

	roles, pagination, err := h.service.ListRoles(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Message:    "Roles fetched successfully",
		Data:       roles,
		Pagination: pagination,
	})
}

func (h *Handler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{
			Message: "invalid role id",
		})
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{
			Message: err.Error(),
		})
		return
	}

	role, err := h.service.UpdateRole(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Message: "Role updated successfully",
		Data:    role,
	})
}

func (h *Handler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{
			Message: "invalid role id",
		})
		return
	}

	if err := h.service.DeleteRole(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Message: "Role deleted successfully",
	})
}
