package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeController(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"token": "token",
	})
}
