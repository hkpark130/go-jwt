package handlers

import (
	"errors"
	"fmt"
	"golang/jwt/api/domain"
	"golang/jwt/api/handlers/auth"
	"golang/jwt/api/repository"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ExpiresCookie = time.Now().Add(time.Hour * 24 * 7)
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetTokenHandler(c *gin.Context, jwtUserRepository *repository.JwtUserRepository) {
	cookie, err := c.Request.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusUnauthorized,
			gin.H{"status": http.StatusUnauthorized,
				"error": "Failed to get Authorization cookie."})
		c.Abort()
		return
	}

	payload, err := auth.Decode(strings.Split(cookie.Value, " ")[1])
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError,
				"error": "Failed to decode token."})
		c.Abort()
		return
	}

	refres, err := jwtUserRepository.GetRefreshToken(payload.Email)
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"status": http.StatusUnauthorized,
					"error": "トークンの有効期間が切りました。"})
			c.Abort()
			return
		} else {
			c.JSON(http.StatusInternalServerError,
				gin.H{"status": http.StatusInternalServerError,
					"error": "Failed to read refresh token."})
			c.Abort()
			return
		}
	}

	c.Data(http.StatusOK,
		"text/html; charset=utf-8",
		[]byte(cookie.Value+refres))
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
	payload := auth.CreatePayload(email)

	if !IsRegisteredUser(c, payload, c.Request.FormValue("password"), jwtUserRepository) {
		return
	}

	accessToken, err := auth.IssueToken(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError,
				"error": "Failed to issue access token."})
		c.Abort()
		return
	}

	modifiedPayload := auth.ModifyForRefreshToken(payload)
	refreshToken, err := auth.IssueToken(modifiedPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError,
				"error": "Failed to issue refresh token."})
		c.Abort()
		return
	}

	err = jwtUserRepository.SetRefreshToken(email, refreshToken, modifiedPayload.Exp)
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
		Expires:  ExpiresCookie,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "Refresh",
		Value:    fmt.Sprintf("%s", refreshToken),
		Expires:  ExpiresCookie,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	c.JSON(http.StatusOK, "OK")
}
