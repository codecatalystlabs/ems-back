package http

import "github.com/gin-gonic/gin"

type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

func (h *Handler) Login(c *gin.Context) {
	c.JSON(200, gin.H{"message": "implement login"})
}
