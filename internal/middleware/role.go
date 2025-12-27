package middleware

import (
	"net/http"

	"github.com/amirullazmi0/kratify-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// RequireRole middleware checks if user has required role
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			response.Error(c, http.StatusForbidden, "Role information not found", nil)
			c.Abort()
			return
		}

		role := userRole.(string)

		// Check if user role is in allowed roles
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, "You don't have permission to access this resource", nil)
		c.Abort()
	}
}

// RequireSuperAdmin middleware allows only superadmin
func RequireSuperAdmin() gin.HandlerFunc {
	return RequireRole("SUPERADMIN")
}

// RequireAdmin middleware allows superadmin and admin
func RequireAdmin() gin.HandlerFunc {
	return RequireRole("SUPERADMIN", "ADMIN")
}
