package httpx

import "github.com/gin-gonic/gin"

func UserID(c *gin.Context) string {
	return c.GetString("user_id")
}
