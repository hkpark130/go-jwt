package handlers

import (
	"golang/jwt/api/handlers/auth"
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

func Authentication(c *gin.Context) {
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")
	payload := &auth.Payload{
		Exp:      time.Now().Add(time.Second * time.Duration(3600)),
		Iat:      time.Now(),
		Email:    email,
		Password: password}

	token := auth.Hashing(payload)

	decode_payload := auth.Decode(token)

	c.JSON(http.StatusOK, decode_payload)
}
