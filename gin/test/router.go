package test

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	return r
}
