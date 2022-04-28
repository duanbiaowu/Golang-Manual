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

func setupFormData() *gin.Engine {
	r := gin.Default()
	r.POST("user", func(c *gin.Context) {
		id := c.Request.PostFormValue("id")
		name := c.Request.PostFormValue("name")
		attrs := c.Request.PostFormValue("attrs")
		c.String(http.StatusOK, "id = %s, name = %s, attrs = %v", id, name, attrs)
	})
	return r
}
