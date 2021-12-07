package handlers

import (
	"errors"
	"fmt"
	"golang/jwt/api/domain"
	"golang/jwt/api/handlers/auth"
	"golang/jwt/api/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTokenHandler(c *gin.Context) {
	cookie, err := c.Request.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest,
				"error": "Failed to get Authorization cookie."})
		c.Abort()
		return
	}

	c.Data(http.StatusOK,
		"text/html; charset=utf-8",
		[]byte(cookie.Value))
}

func IsRegisteredUser(c *gin.Context, payload *auth.Payload, jwtUserRepository *repository.JwtUserRepository) bool {
	jwtUser := &domain.JwtUser{Email: payload.Email, Password: payload.Password}
	user, err := jwtUserRepository.LoginEmailPassword(jwtUser)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError,
				"error": fmt.Sprintf("Failed to read user form DB: %s", err)})
		c.Abort()
		return false
	}

	if (domain.JwtUser{}) != *user {
		return true
	}

	c.JSON(http.StatusUnauthorized,
		gin.H{"status": http.StatusUnauthorized,
			"error": "メールアドレスとパスワードをもう一度確認してください。"})
	c.Abort()
	return false
}

func Login(c *gin.Context, jwtUserRepository *repository.JwtUserRepository) {
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")
	payload := &auth.Payload{
		Exp:      time.Now().Add(time.Second * time.Duration(3600)),
		Iat:      time.Now(),
		Email:    email,
		Password: password}

	if !IsRegisteredUser(c, payload, jwtUserRepository) {
		return
	}

	token, err := auth.Hashing(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError,
				"error": "Failed to hashing."})
		c.Abort()
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "Authorization",
		Value:    fmt.Sprintf("Bearer %s", token),
		Expires:  payload.Exp,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	c.JSON(http.StatusOK, "OK")
}
