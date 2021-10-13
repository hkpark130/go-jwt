package handlers

import (
	"golang/jwt/api/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: requset data -> auth package -> return data here

func GetTokenHandler(c *gin.Context) {

	c.Data(http.StatusOK,
		"text/html; charset=utf-8",
		[]byte("token"))
}

func Authentication(c *gin.Context) {
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")
	user := &domain.JwtUser{Email: email, Password: password}

	c.JSON(http.StatusOK, user)
}
