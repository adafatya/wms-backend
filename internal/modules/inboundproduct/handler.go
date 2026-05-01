package inboundproduct

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

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/inbound-products")
	group.POST("/schedule", h.CreateSchedule)
	group.GET("/schedule", h.ListSchedules)
	group.GET("/schedule/:id", h.GetSchedule)
	group.PUT("/schedule/:id", h.UpdateSchedule)
	group.DELETE("/schedule/:id", h.DeleteSchedule)

	group.POST("/receipt", h.CreateReceipt)
	group.GET("/receipt", h.ListReceipts)
	group.GET("/receipt/:id", h.GetReceipt)
}

func (h *Handler) CreateSchedule(c *gin.Context) {
	var req CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: err.Error()})
		return
	}

	res, err := h.service.CreateSchedule(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.StandardResponse{
		Message: "Schedule created successfully",
		Data:    res,
	})
}

func (h *Handler) GetSchedule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: "invalid id"})
		return
	}

	res, err := h.service.GetSchedule(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.StandardResponse{Message: "schedule not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Message: "Schedule fetched successfully",
		Data:    res,
	})
}

func (h *Handler) ListSchedules(c *gin.Context) {
	page, limit, err := models.ParsePaginationQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: err.Error()})
		return
	}

	res, pagination, err := h.service.ListSchedules(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Message:    "Schedules fetched successfully",
		Data:       res,
		Pagination: pagination,
	})
}

func (h *Handler) UpdateSchedule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: "invalid id"})
		return
	}

	var req UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: err.Error()})
		return
	}

	res, err := h.service.UpdateSchedule(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Message: "Schedule updated successfully",
		Data:    res,
	})
}

func (h *Handler) DeleteSchedule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: "invalid id"})
		return
	}

	if err := h.service.DeleteSchedule(c.Request.Context(), id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.StandardResponse{Message: "schedule not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{Message: "Schedule deleted successfully"})
}

func (h *Handler) CreateReceipt(c *gin.Context) {
	var req CreateReceiptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: err.Error()})
		return
	}

	res, err := h.service.CreateReceipt(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.StandardResponse{
		Message: "Receipt created successfully",
		Data:    res,
	})
}

func (h *Handler) GetReceipt(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: "invalid id"})
		return
	}

	res, err := h.service.GetReceipt(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.StandardResponse{Message: "receipt not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Message: "Receipt fetched successfully",
		Data:    res,
	})
}

func (h *Handler) ListReceipts(c *gin.Context) {
	page, limit, err := models.ParsePaginationQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: err.Error()})
		return
	}

	res, pagination, err := h.service.ListReceipts(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Message:    "Receipts fetched successfully",
		Data:       res,
		Pagination: pagination,
	})
}
