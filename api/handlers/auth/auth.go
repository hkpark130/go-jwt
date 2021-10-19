package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTokenHandler(c *gin.Context) {

	c.Data(http.StatusOK,
		"text/html; charset=utf-8",
		[]byte("token"))
}
