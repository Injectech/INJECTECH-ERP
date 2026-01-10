package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	domaininventory "backend/internal/domain/inventory"
	"backend/internal/transport/http/middleware"
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
	rg.POST("/inventory", middleware.RequirePermission("inventory.write"), h.create)
	rg.PATCH("/inventory/:id/adjust", middleware.RequirePermission("inventory.write"), h.adjust)
	rg.PATCH("/inventory/:id/location", middleware.RequirePermission("inventory.write"), h.updateLocation)
	rg.GET("/inventory/product/:product_id", middleware.RequirePermission("inventory.read"), h.listByProduct)
	rg.GET("/inventory", middleware.RequirePermission("inventory.read"), h.list)
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
		ID:        uuid.NewString(),
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

func (h *InventoryHandler) updateLocation(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Location string `json:"location"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	if err := h.uc.UpdateLocation(c.Request.Context(), id, req.Location); err != nil {
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

func (h *InventoryHandler) list(c *gin.Context) {
	res, err := h.uc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "success", "data": res})
}
