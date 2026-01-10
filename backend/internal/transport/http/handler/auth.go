package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	domainuser "backend/internal/domain/user"
	usecaseauth "backend/internal/usecase/auth"
)

// AuthHandler exposes authentication endpoints.
type AuthHandler struct {
	authUC *usecaseauth.Service
	secure bool
}

func NewAuthHandler(authUC *usecaseauth.Service, secure bool) *AuthHandler {
	return &AuthHandler{authUC: authUC, secure: secure}
}

func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", h.register)
	rg.POST("/login", h.login)
	rg.POST("/refresh", h.refresh)
}

func (h *AuthHandler) register(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	session, err := h.authUC.Register(c.Request.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	setRefreshCookie(c, session.Tokens.RefreshToken, session.Tokens.RefreshExp, h.secure)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "registered",
		"data": gin.H{
			"access_token":      session.Tokens.AccessToken,
			"access_expires_at": session.Tokens.AccessExp,
			"user":              userResponse(session.User),
		},
	})
}

func (h *AuthHandler) login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	session, err := h.authUC.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		status := http.StatusUnauthorized
		c.JSON(status, gin.H{"success": false, "message": err.Error()})
		return
	}

	setRefreshCookie(c, session.Tokens.RefreshToken, session.Tokens.RefreshExp, h.secure)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "success",
		"data": gin.H{
			"access_token":      session.Tokens.AccessToken,
			"access_expires_at": session.Tokens.AccessExp,
			"user":              userResponse(session.User),
		},
	})
}

func (h *AuthHandler) refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "missing refresh token"})
		return
	}
	tokens, err := h.authUC.Refresh(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
		return
	}
	setRefreshCookie(c, tokens.RefreshToken, tokens.RefreshExp, h.secure)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "success", "data": gin.H{"access_token": tokens.AccessToken, "access_expires_at": tokens.AccessExp}})
}

func setRefreshCookie(c *gin.Context, token string, exp time.Time, secure bool) {
	maxAge := int(time.Until(exp).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}
	c.SetCookie("refresh_token", token, maxAge, "/", "", secure, true)
}

type authUser struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

func userResponse(u domainuser.User) authUser {
	return authUser{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Roles: u.Roles,
	}
}
