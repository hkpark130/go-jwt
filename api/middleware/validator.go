package middleware

import (
	"golang/jwt/api/handlers/auth"
	"net/http"
	"net/mail"
	"strings"

	"github.com/gin-gonic/gin"
)

func LoginFormValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Request.FormValue("email")
		password := c.Request.FormValue("password")
		_, email_err := mail.ParseAddress(email)

		if len(email) < 1 || len(password) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Please enter your email or password.",
			})
			c.Abort()
			return
		}
		if email_err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Please check your email again.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("Authorization")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No Authorization cookie.",
			})
			c.Abort()
			return
		}

		if !auth.IsTokenVerified(strings.Split(cookie.Value, " ")[1]) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Fail to verify.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
