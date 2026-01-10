package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/transport/http/middleware"
	usecaselocation "backend/internal/usecase/location"
)

// LocationHandler exposes location endpoints.
type LocationHandler struct {
	uc *usecaselocation.Service
}

func NewLocationHandler(uc *usecaselocation.Service) *LocationHandler {
	return &LocationHandler{uc: uc}
}

func (h *LocationHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/locations", middleware.RequirePermission("location.write"), h.create)
	rg.GET("/locations", middleware.RequirePermission("location.read"), h.list)
}

func (h *LocationHandler) create(c *gin.Context) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	location, err := h.uc.Create(c.Request.Context(), req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "created", "data": location})
}

func (h *LocationHandler) list(c *gin.Context) {
	locations, err := h.uc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "success", "data": locations})
}
