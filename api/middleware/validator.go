package middleware

import (
	"fmt"
	"golang/jwt/api/handlers"
	"golang/jwt/api/handlers/auth"
	"golang/jwt/api/repository"
	"net/http"
	"net/mail"
	"os"
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

func Authorization(jwtUserRepository *repository.JwtUserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("Authorization")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No Authorization cookie.",
			})
			c.Abort()
			return
		}

		if isVerified, err := auth.IsTokenVerified(strings.Split(cookie.Value, " ")[1]); !isVerified {
			if err.Error() == "Expired JWT Token" {
				refreshToken, err := c.Request.Cookie("Refresh")
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{
						"error": "No Refresh Token cookie.",
					})
					c.Abort()
					return
				}

				// Access Token 再発行
				if isVerified, _ := auth.IsTokenVerified(refreshToken.Value); isVerified {
					pldat, err := auth.Decode(strings.Split(cookie.Value, " ")[1])
					accessToken, refreshToken, err := auth.ReissueToken(pldat, jwtUserRepository)
					if err != nil {
						c.JSON(http.StatusUnauthorized, gin.H{
							"error": "Fail to Reissue.",
						})
						c.Abort()
						return
					}

					var domainName string
					if domainName = "localhost"; os.Getenv("SECRET_KEY") == "production.conf" {
						domainName = "hkpark130.p-e.kr"
					}
					http.SetCookie(c.Writer, &http.Cookie{
						Domain:   domainName,
						Name:     "Authorization",
						Value:    fmt.Sprintf("Bearer %s", accessToken),
						Expires:  handlers.ExpiresCookie,
						Path:     "/",
						Secure:   true,
						HttpOnly: true,
					})

					http.SetCookie(c.Writer, &http.Cookie{
						Domain:   domainName,
						Name:     "Refresh",
						Value:    fmt.Sprintf("%s", refreshToken),
						Expires:  handlers.ExpiresCookie,
						Path:     "/",
						Secure:   true,
						HttpOnly: true,
					})
				} else {
					c.JSON(http.StatusUnauthorized, gin.H{
						"error": "Fail to verify.",
					})
					c.Abort()
					return
				}
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Fail to verify.",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
