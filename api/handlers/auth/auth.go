package auth

import (
	"golang/jwt/api/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Jwt struct {
	AccessToken  string
	RefreshToken string
}

func GetTokenHandler(c *gin.Context) {

	c.Data(http.StatusOK,
		"text/html; charset=utf-8",
		[]byte("token"))
}

func Authentication(c *gin.Context) {
	email := c.DefaultPostForm("email", "")
	password := c.DefaultPostForm("password", "")
	user := &domain.JwtUser{Email: email, Password: password}

	c.JSON(http.StatusOK, user)
}
