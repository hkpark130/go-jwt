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

	http.Handle("/token", auth.GetTokenHandler)
	http.Handle("/login", auth.RenderLoginView)

	http.Handle("/api/login", auth.Authentication)

	r.Run(":3000")
}
