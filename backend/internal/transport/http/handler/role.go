package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/transport/http/middleware"
	usecaserole "backend/internal/usecase/role"
)

// RoleHandler exposes role endpoints.
type RoleHandler struct {
	uc *usecaserole.Service
}

func NewRoleHandler(uc *usecaserole.Service) *RoleHandler {
	return &RoleHandler{uc: uc}
}

func (h *RoleHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/roles", middleware.RequirePermission("role.write"), h.create)
	rg.GET("/roles", middleware.RequirePermission("role.read"), h.list)
}

func (h *RoleHandler) create(c *gin.Context) {
	var req struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Permissions []string `json:"permissions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	role, err := h.uc.Create(c.Request.Context(), req.Name, req.Description, req.Permissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "created", "data": role})
}

func (h *RoleHandler) list(c *gin.Context) {
	roles, err := h.uc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "success", "data": roles})
}
