package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	rbacapp "dispatch/internal/modules/rbac/application"
)

func ScopeContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		scopeType := strings.TrimSpace(c.GetHeader("X-Scope-Type"))
		scopeID := strings.TrimSpace(c.GetHeader("X-Scope-ID"))
		if scopeType == "" {
			scopeType = strings.TrimSpace(c.Query("scope_type"))
		}
		if scopeID == "" {
			scopeID = strings.TrimSpace(c.Query("scope_id"))
		}
		if scopeType != "" {
			c.Set("scope_type", strings.ToUpper(scopeType))
		}
		if scopeID != "" {
			c.Set("scope_id", scopeID)
		}
		c.Next()
	}
}

func RequirePermission(rbacSvc *rbacapp.Service, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthenticated"})
			return
		}

		scopeType := c.GetString("scope_type")
		var scopeID *string
		if v := c.GetString("scope_id"); v != "" {
			scopeID = &v
		}

		ok, err := rbacSvc.HasPermission(c.Request.Context(), userID, permission, scopeType, scopeID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to evaluate permission"})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message":    "forbidden",
				"permission": permission,
				"scope_type": scopeType,
				"scope_id":   scopeID,
			})
			return
		}
		c.Next()
	}
}
