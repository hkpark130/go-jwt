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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

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

func IsRegisteredUser(c *gin.Context, payload *auth.Payload, password string, jwtUserRepository *repository.JwtUserRepository) bool {
	jwtUser := &domain.JwtUser{Email: payload.Email, Password: password}
	user, err := jwtUserRepository.LoginEmailPassword(jwtUser)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError,
				"error": fmt.Sprintf("Failed to read user form DB: %s", err)})
		c.Abort()
		return false
	}

	if !CheckPasswordHash(jwtUser.Password, user.Password) {
		c.JSON(http.StatusUnauthorized,
			gin.H{"status": http.StatusUnauthorized,
				"error": "メールアドレスとパスワードをもう一度確認してください。"})
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
	payload := &auth.Payload{
		Exp:   time.Now().Add(time.Second * time.Duration(3600)),
		Iat:   time.Now(),
		Email: email}

	if !IsRegisteredUser(c, payload, c.Request.FormValue("password"), jwtUserRepository) {
		return
	}

	accessToken, err := auth.IssueAccessToken(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError,
				"error": "Failed to issue access token."})
		c.Abort()
		return
	}

	refreshToken, err := auth.IssueRefreshToken(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError,
				"error": "Failed to issue refresh token."})
		c.Abort()
		return
	}

	err = jwtUserRepository.SetRefreshToken(email, accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError,
				"error": "Failed to set token in redis."})
		c.Abort()
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "Authorization",
		Value:    fmt.Sprintf("Bearer %s", accessToken),
		Expires:  payload.Exp,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "Refresh",
		Value:    fmt.Sprintf("%s", refreshToken),
		Expires:  payload.Exp,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	c.JSON(http.StatusOK, "OK")
}
