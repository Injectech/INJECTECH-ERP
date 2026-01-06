package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	usecaseuser "backend/internal/usecase/user"
)

// UserHandler exposes user management endpoints.
type UserHandler struct {
	uc *usecaseuser.Service
}

func NewUserHandler(uc *usecaseuser.Service) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/users", h.create)
	rg.GET("/users/:id", h.get)
}

func (h *UserHandler) create(c *gin.Context) {
	var req struct {
		Email    string   `json:"email"`
		Password string   `json:"password"`
		Name     string   `json:"name"`
		Roles    []string `json:"roles"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	user, err := h.uc.Create(c.Request.Context(), req.Email, req.Password, req.Name, req.Roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "created", "data": user})
}

func (h *UserHandler) get(c *gin.Context) {
	id := c.Param("id")
	user, err := h.uc.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "success", "data": user})
}
