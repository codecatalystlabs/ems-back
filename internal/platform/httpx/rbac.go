package httpx

import "github.com/gin-gonic/gin"

func ScopeType(c *gin.Context) string {
	return c.GetString("scope_type")
}

func ScopeID(c *gin.Context) *string {
	v := c.GetString("scope_id")
	if v == "" {
		return nil
	}
	return &v
}
