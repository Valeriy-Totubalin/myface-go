package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "test ok",
	})
}

func (h *Handler) signIn(c *gin.Context) {

}
