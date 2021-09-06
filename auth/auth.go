package auth

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type Jwt struct {
	AccessToken string
	RefreshToken string
}

type User struct {
	Email string
	Password string
}

func GetTokenHandler(c *gin.Context) {
	c.Data(http.StatusOK, 
		"text/html; charset=utf-8", 
		[]byte("token"))
}

func RenderLoginView(c *gin.Context) {
	c.HTML(http.StatusOK, 
		"login.html", 
		gin.H{})
}

func Authentication(c *gin.Context) {
	email := c.DefaultPostForm("email", "")
	password := c.DefaultPostForm("password", "")
	user := &User{Email: email, Password: password}

	c.JSON(http.StatusOK, user)
}
