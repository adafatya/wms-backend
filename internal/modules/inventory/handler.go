package inventory

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
	// Locations
	router.POST("/locations", h.CreateLocation)
	router.GET("/locations", h.ListLocations)
	router.GET("/locations/:id", h.GetLocation)
	router.PUT("/locations/:id", h.UpdateLocation)
	router.DELETE("/locations/:id", h.DeleteLocation)

	// Products
	router.POST("/products", h.CreateProduct)
	router.GET("/products", h.ListProducts)
	router.GET("/products/:id", h.GetProduct)
	router.PUT("/products/:id", h.UpdateProduct)
	router.DELETE("/products/:id", h.DeleteProduct)

	// Inventories
	router.POST("/inventories", h.UpsertInventories)
	router.GET("/inventories", h.ListInventories)
	router.GET("/inventories/:location_id/:product_id", h.GetInventoryStock)
}

// Location handlers
func (h *Handler) CreateLocation(c *gin.Context) {
	var req CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: err.Error()})
		return
	}

	res, err := h.service.CreateLocation(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.StandardResponse{
		Message: "Location created successfully",
		Data:    res,
	})
}

func (h *Handler) ListLocations(c *gin.Context) {
	page, limit, err := models.ParsePaginationQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: err.Error()})
		return
	}
	search := c.Query("q")

	res, pagination, err := h.service.ListLocations(c.Request.Context(), search, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Message:    "Locations fetched successfully",
		Data:       res,
		Pagination: pagination,
	})
}

func (h *Handler) GetLocation(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: "invalid location id"})
		return
	}

	res, err := h.service.GetLocation(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.StandardResponse{Message: "location not found"})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Message: "Location fetched successfully",
		Data:    res,
	})
}

func (h *Handler) UpdateLocation(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: "invalid location id"})
		return
	}

	var req UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: err.Error()})
		return
	}

	res, err := h.service.UpdateLocation(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Message: "Location updated successfully",
		Data:    res,
	})
}

func (h *Handler) DeleteLocation(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: "invalid location id"})
		return
	}

	if err := h.service.DeleteLocation(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{Message: "Location deleted successfully"})
}

// Product handlers
func (h *Handler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: err.Error()})
		return
	}

	res, err := h.service.CreateProduct(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.StandardResponse{
		Message: "Product created successfully",
		Data:    res,
	})
}

func (h *Handler) ListProducts(c *gin.Context) {
	page, limit, err := models.ParsePaginationQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: err.Error()})
		return
	}
	search := c.Query("q")

	res, pagination, err := h.service.ListProducts(c.Request.Context(), search, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Message:    "Products fetched successfully",
		Data:       res,
		Pagination: pagination,
	})
}

func (h *Handler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: "invalid product id"})
		return
	}

	res, err := h.service.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.StandardResponse{Message: "product not found"})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Message: "Product fetched successfully",
		Data:    res,
	})
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: "invalid product id"})
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: err.Error()})
		return
	}

	res, err := h.service.UpdateProduct(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Message: "Product updated successfully",
		Data:    res,
	})
}

func (h *Handler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: "invalid product id"})
		return
	}

	if err := h.service.DeleteProduct(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{Message: "Product deleted successfully"})
}

// Inventory handlers
func (h *Handler) UpsertInventories(c *gin.Context) {
	var req []InventoryInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: err.Error()})
		return
	}

	if err := h.service.UpsertInventories(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{Message: "Inventories updated successfully"})
}

func (h *Handler) ListInventories(c *gin.Context) {
	page, limit, err := models.ParsePaginationQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: err.Error()})
		return
	}

	res, pagination, err := h.service.ListInventories(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Message:    "Inventories fetched successfully",
		Data:       res,
		Pagination: pagination,
	})
}

func (h *Handler) GetInventoryStock(c *gin.Context) {
	locationID, err := strconv.ParseInt(c.Param("location_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: "invalid location id"})
		return
	}
	productID, err := strconv.ParseInt(c.Param("product_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{Message: "invalid product id"})
		return
	}

	res, err := h.service.GetInventoryStock(c.Request.Context(), productID, locationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
