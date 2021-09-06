package app

import (
	"net/http"
	"golang/jwt/auth"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, 
			"text/html; charset=utf-8", 
			[]byte("index"))	
	})

	r.GET("/token", func(c *gin.Context) { auth.GetTokenHandler(c) })
	r.GET("/login", func(c *gin.Context) { auth.RenderLoginView(c) })

	r.POST("/api/login", func(c *gin.Context) { auth.Authentication(c) })

	return r
}