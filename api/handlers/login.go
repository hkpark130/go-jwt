package handlers

import (
	"fmt"
	"golang/jwt/api/handlers/auth"
	"golang/jwt/api/repository"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TODO: requset data -> auth package -> return data here

func GetTokenHandler(c *gin.Context) {

	c.Data(http.StatusOK,
		"text/html; charset=utf-8",
		[]byte("token"))
}

func IsRegisteredUser(payload *auth.Payload, jwtUserRepository *repository.JwtUserRepository) bool {
	err := jwtUserRepository.CheckUser(payload.Email, payload.Password)
	if err == nil {
		return true
	}
	if err.Error() != "record not found" {
		log.Fatal("Failed to read user form DB:", err)
	}

	return false
}

func Authentication(c *gin.Context, jwtUserRepository *repository.JwtUserRepository) {
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")
	payload := &auth.Payload{
		Exp:      time.Now().Add(time.Second * time.Duration(3600)),
		Iat:      time.Now(),
		Email:    email,
		Password: password}

	if !IsRegisteredUser(payload, jwtUserRepository) {
		// TODO: Redirect login page
		c.JSON(http.StatusUnauthorized,
			gin.H{"status": http.StatusUnauthorized,
				"error": "メールアドレスとパスワードをもう一度確認してください。"})
		c.Abort()
		return
	}

	token := auth.Hashing(payload)
	c.Header("Authorization", fmt.Sprintf("Bearer %s", token))

	c.JSON(http.StatusOK, "OK")
}
