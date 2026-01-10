package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/transport/http/middleware"
	usecaseproduct "backend/internal/usecase/product"
)

// ProductHandler exposes product endpoints.
type ProductHandler struct {
	uc *usecaseproduct.Service
}

func NewProductHandler(uc *usecaseproduct.Service) *ProductHandler {
	return &ProductHandler{uc: uc}
}

func (h *ProductHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/products", middleware.RequirePermission("product.write"), h.create)
	rg.GET("/products", middleware.RequirePermission("product.read"), h.list)
	rg.PUT("/products/:id", middleware.RequirePermission("product.write"), h.update)
	rg.DELETE("/products/:id", middleware.RequirePermission("product.write"), h.delete)
}

func (h *ProductHandler) create(c *gin.Context) {
	var req struct {
		SKU         string  `json:"sku"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	product, err := h.uc.Create(c.Request.Context(), req.SKU, req.Name, req.Description, req.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "created", "data": product})
}

func (h *ProductHandler) list(c *gin.Context) {
	products, err := h.uc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "success", "data": products})
}

func (h *ProductHandler) update(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		SKU         string  `json:"sku"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	product, err := h.uc.Update(c.Request.Context(), id, req.SKU, req.Name, req.Description, req.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "updated", "data": product})
}

func (h *ProductHandler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "deleted"})
}
