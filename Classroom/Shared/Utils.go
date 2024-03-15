package Shared

import "github.com/gin-gonic/gin"

func GetQueryParam(c *gin.Context, keys ...string) string {
	for _, key := range keys {
		if value, exists := c.GetQuery(key); exists && value != "" {
			return value
		}
	}
	return ""
}
