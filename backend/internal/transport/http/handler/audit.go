package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/transport/http/middleware"
	usecaseaudit "backend/internal/usecase/audit"
)

// AuditHandler exposes audit log endpoints.
type AuditHandler struct {
	uc *usecaseaudit.Service
}

func NewAuditHandler(uc *usecaseaudit.Service) *AuditHandler {
	return &AuditHandler{uc: uc}
}

func (h *AuditHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/audit/logs", middleware.RequirePermission("audit.read"), h.list)
}

func (h *AuditHandler) list(c *gin.Context) {
	actorID := c.Query("actor_id")
	logs, err := h.uc.List(c.Request.Context(), actorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "success", "data": logs})
}
