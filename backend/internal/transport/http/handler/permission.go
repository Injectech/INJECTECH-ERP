package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/transport/http/middleware"
	usecasepermission "backend/internal/usecase/permission"
)

// PermissionHandler exposes permission endpoints.
type PermissionHandler struct {
	uc *usecasepermission.Service
}

func NewPermissionHandler(uc *usecasepermission.Service) *PermissionHandler {
	return &PermissionHandler{uc: uc}
}

func (h *PermissionHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/permissions", middleware.RequirePermission("permission.write"), h.create)
	rg.GET("/permissions", middleware.RequirePermission("permission.read"), h.list)
}

func (h *PermissionHandler) create(c *gin.Context) {
	var req struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	p, err := h.uc.Create(c.Request.Context(), req.Code, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "created", "data": p})
}

func (h *PermissionHandler) list(c *gin.Context) {
	items, err := h.uc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "success", "data": items})
}
