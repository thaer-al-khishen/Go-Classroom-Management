package Middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webapptrials/Classroom/Auth/Models"
	"webapptrials/Classroom/Auth/Utils"
	"webapptrials/Classroom/Shared"
)

// RoleCheckMiddleware returns a Gin middleware that checks for a specific role
func RoleCheckMiddleware(requiredRole Models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the user's role from the JWT token
		// This is simplified; you'll need to parse the JWT token and extract the Claims
		claims, err := Utils.ParseToken(c.GetHeader("Authorization"))
		if err != nil {
			// Handle parsing error (e.g., token is expired or invalid)
			Shared.SendGinGenericApiResponse(c, http.StatusUnauthorized, "Unauthorized", "", err.Error())
			c.Abort()
			return
		}

		userRole := claims.Role

		if userRole != requiredRole {
			Shared.SendGinGenericApiResponse(c, http.StatusForbidden, "You can't access this resource", "", "You can't access this resource")
			c.Abort()
			return
		}

		c.Next()
	}

}
