package handlers

import (
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

func IsRegistered(payload *auth.Payload, jwtUserRepository *repository.JwtUserRepository) bool {
	_, err := jwtUserRepository.GetUserByEmail(payload.Email)
	if err == nil {
		return true
	} else if err.Error() == "record not found" {
		return false
	} else {
		log.Fatal("Failed to read user form DB:", err)
		return false
	}
}

func Authentication(c *gin.Context, jwtUserRepository *repository.JwtUserRepository) {
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")
	payload := &auth.Payload{
		Exp:      time.Now().Add(time.Second * time.Duration(3600)),
		Iat:      time.Now(),
		Email:    email,
		Password: password}

	if !IsRegistered(payload, jwtUserRepository) {
		// TODO: Redirect login page
		c.JSON(http.StatusUnauthorized,
			gin.H{"status": http.StatusUnauthorized,
				"error": "Unregistered user"})
		c.Abort()
		return
	}

	token := auth.Hashing(payload)

	claim := auth.Decode(token)

	c.JSON(http.StatusOK, claim)
}
