package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	domaininventory "backend/internal/domain/inventory"
	usecaseinventory "backend/internal/usecase/inventory"
)

// InventoryHandler exposes inventory endpoints.
type InventoryHandler struct {
	uc *usecaseinventory.Service
}

func NewInventoryHandler(uc *usecaseinventory.Service) *InventoryHandler {
	return &InventoryHandler{uc: uc}
}

func (h *InventoryHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/inventory", h.create)
	rg.PATCH("/inventory/:id/adjust", h.adjust)
	rg.GET("/inventory/product/:product_id", h.listByProduct)
}

func (h *InventoryHandler) create(c *gin.Context) {
	var req struct {
		ProductID string `json:"product_id"`
		Quantity  int64  `json:"quantity"`
		Location  string `json:"location"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	inv := domaininventory.Inventory{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Location:  req.Location,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	created, err := h.uc.Create(c.Request.Context(), inv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "created", "data": created})
}

func (h *InventoryHandler) adjust(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Delta int64 `json:"delta"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := h.uc.Adjust(c.Request.Context(), id, req.Delta); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "updated"})
}

func (h *InventoryHandler) listByProduct(c *gin.Context) {
	productID := c.Param("product_id")
	res, err := h.uc.ListByProduct(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "success", "data": res})
}
