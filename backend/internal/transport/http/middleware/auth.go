package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	usecaseauth "backend/internal/usecase/auth"
)

const (
	ContextUserID = "user_id"
	ContextRoles  = "roles"
	ContextPerms  = "perms"
)

// Auth middleware validates JWT access token and populates context.
func Auth(authUC *usecaseauth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.GetHeader("Authorization")
		if raw == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "missing token"})
			return
		}
		token := strings.TrimPrefix(raw, "Bearer ")
		claims, err := authUC.ParseAccess(token)
		if err != nil {
			status := http.StatusUnauthorized
			if usecaseauth.IsTokenExpired(err) {
				status = http.StatusUnauthorized
			}
			c.AbortWithStatusJSON(status, gin.H{"success": false, "message": "invalid token"})
			return
		}
		c.Set(ContextUserID, claims.Subject)
		c.Set(ContextRoles, claims.Roles)
		c.Set(ContextPerms, claims.Permissions)
		c.Next()
	}
}

// RequirePermission enforces RBAC at handler-level.
func RequirePermission(code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get(ContextPerms)
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "message": "missing permissions"})
			return
		}
		perms, ok := val.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "message": "invalid permissions context"})
			return
		}
		for _, p := range perms {
			if p == code {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "message": "forbidden"})
	}
}
