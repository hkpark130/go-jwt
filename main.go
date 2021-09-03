package main

import (
	"net/http"
	"golang/jwt/auth"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, 
			"text/html; charset=utf-8", 
			[]byte("index"))	
	})

	r.GET("/token", func(c *gin.Context) {auth.GetTokenHandler(c)})
	http.Handle("/login", auth.RenderLoginView)

	http.Handle("/api/login", auth.Authentication)

	r.Run(":3000")
}
