package outboundproduct

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/adafatya/wms-backend/internal/models"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Customer
func (h *Handler) CreateCustomer(c *gin.Context) {
	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: err.Error()})
		return
	}

	res, err := h.service.CreateCustomer(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.StandardResponse{Data: res})
}

func (h *Handler) GetCustomer(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	res, err := h.service.GetCustomer(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.StandardResponse{Message: "customer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.StandardResponse{Data: res})
}

func (h *Handler) ListCustomers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	res, pagination, err := h.service.ListCustomers(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{Data: res, Pagination: pagination})
}

func (h *Handler) UpdateCustomer(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: err.Error()})
		return
	}

	res, err := h.service.UpdateCustomer(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{Data: res})
}

func (h *Handler) DeleteCustomer(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.DeleteCustomer(c.Request.Context(), id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.StandardResponse{Message: "customer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.StandardResponse{Message: "customer deleted"})
}

// Delivery Order
func (h *Handler) CreateDeliveryOrder(c *gin.Context) {
	var req CreateDeliveryOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: err.Error()})
		return
	}

	res, err := h.service.CreateDeliveryOrder(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.StandardResponse{Data: res})
}

func (h *Handler) GetDeliveryOrder(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	res, err := h.service.GetDeliveryOrder(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.StandardResponse{Message: "delivery order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.StandardResponse{Data: res})
}

func (h *Handler) ListDeliveryOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	res, pagination, err := h.service.ListDeliveryOrders(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{Data: res, Pagination: pagination})
}

func (h *Handler) UpdateDeliveryOrder(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req UpdateDeliveryOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: err.Error()})
		return
	}

	res, err := h.service.UpdateDeliveryOrder(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{Data: res})
}

func (h *Handler) DeleteDeliveryOrder(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.DeleteDeliveryOrder(c.Request.Context(), id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.StandardResponse{Message: "delivery order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.StandardResponse{Message: "delivery order deleted"})
}

// Delivery
func (h *Handler) CreateDelivery(c *gin.Context) {
	var req CreateDeliveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: err.Error()})
		return
	}

	res, err := h.service.CreateDelivery(c.Request.Context(), req)
	if err != nil {
		// Handle insufficient stock error specifically
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.StandardResponse{Data: res})
}

func (h *Handler) GetDelivery(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	res, err := h.service.GetDelivery(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.StandardResponse{Message: "delivery not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.StandardResponse{Data: res})
}

func (h *Handler) ListDeliveries(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	res, pagination, err := h.service.ListDeliveries(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{Data: res, Pagination: pagination})
}

func (h *Handler) UpdateDelivery(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req UpdateDeliveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: err.Error()})
		return
	}

	res, err := h.service.UpdateDelivery(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{Data: res})
}

func (h *Handler) DeleteDelivery(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.DeleteDelivery(c.Request.Context(), id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.StandardResponse{Message: "delivery not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.StandardResponse{Message: "delivery deleted"})
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	outbound := router.Group("/outbound-products")
	{
		// Customers
		outbound.POST("/customer", h.CreateCustomer)
		outbound.GET("/customer", h.ListCustomers)
		outbound.GET("/customer/:id", h.GetCustomer)
		outbound.PUT("/customer/:id", h.UpdateCustomer)
		outbound.DELETE("/customer/:id", h.DeleteCustomer)

		// Delivery Orders
		outbound.POST("/delivery-order", h.CreateDeliveryOrder)
		outbound.GET("/delivery-order", h.ListDeliveryOrders)
		outbound.GET("/delivery-order/:id", h.GetDeliveryOrder)
		outbound.PUT("/delivery-order/:id", h.UpdateDeliveryOrder)
		outbound.DELETE("/delivery-order/:id", h.DeleteDeliveryOrder)

		// Deliveries
		outbound.POST("/delivery", h.CreateDelivery)
		outbound.GET("/delivery", h.ListDeliveries)
		outbound.GET("/delivery/:id", h.GetDelivery)
		outbound.PUT("/delivery/:id", h.UpdateDelivery)
		outbound.DELETE("/delivery/:id", h.DeleteDelivery)
	}
}
