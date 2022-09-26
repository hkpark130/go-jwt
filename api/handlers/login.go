package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang/jwt/api/domain"
	"golang/jwt/api/handlers/auth"
	"golang/jwt/api/repository"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ExpiresCookie = time.Now().Add(time.Hour * 24 * 1)
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func AdminHandler(c *gin.Context, jwtUserRepository *repository.JwtUserRepository) {
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

	if payload.Permission == "role_admin" {
		users, err := jwtUserRepository.GetUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"status": http.StatusInternalServerError,
					"error": "Failed to read data."})
			c.Abort()
			return
		}

		jsonUsers, err := json.Marshal(users)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"status": http.StatusInternalServerError,
					"error": "Failed to read data."})
			c.Abort()
			return
		}

		c.Data(http.StatusOK,
			"text/html; charset=utf-8",
			[]byte(jsonUsers))
	} else {
		c.JSON(http.StatusForbidden,
			gin.H{"status": http.StatusForbidden,
				"error": "Unauthorized user."})
		c.Abort()
		return
	}
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

	refresh, err := jwtUserRepository.GetRefreshToken(payload.Email)
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"status": http.StatusUnauthorized,
					"error": "The refresh token provided has expired."})
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

	result, err := json.Marshal([]string{strings.Split(cookie.Value, " ")[1], refresh})
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError,
				"error": "Failed to read data."})
		c.Abort()
		return
	}

	c.Data(http.StatusOK,
		"text/html; charset=utf-8",
		[]byte(result))
}

func LogoutHandler(c *gin.Context, jwtUserRepository *repository.JwtUserRepository) {
	_, err := c.Request.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusUnauthorized,
			gin.H{"status": http.StatusUnauthorized,
				"error": "Failed to get Authorization cookie."})
		c.Abort()
		return
	}

	var domainName string
	if domainName = "localhost"; os.Getenv("CONF_FILE") == "production.conf" {
		domainName = "hkpark130.p-e.kr"
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Domain:   domainName,
		Name:     "Authorization",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Domain:   domainName,
		Name:     "Refresh",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	c.JSON(http.StatusOK, "OK")
}

func IsRegisteredUser(c *gin.Context, payload *auth.Payload, password string, jwtUserRepository *repository.JwtUserRepository) bool {
	jwtUser := NewJwtUser(payload.Email, "", password)
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
				"error": "Please check your email address and password again."})
		c.Abort()
		return false
	}

	if (domain.JwtUser{}) != *user {
		return true
	}

	c.JSON(http.StatusUnauthorized,
		gin.H{"status": http.StatusUnauthorized,
			"error": "Please check your email address and password again."})
	c.Abort()
	return false
}

func SetPayloadFromDB(payload *auth.Payload, jwtUserRepository *repository.JwtUserRepository) bool {
	jwtUser := NewJwtUser(payload.Email, "", "")
	user, err := jwtUserRepository.LoginEmailPassword(jwtUser)
	if err != nil {
		return false
	}

	payload.Permission = user.Permission
	return true
}

func Login(c *gin.Context, jwtUserRepository *repository.JwtUserRepository) {
	email := c.Request.FormValue("email")
	payload := auth.CreatePayload(email, "")

	if !IsRegisteredUser(c, payload, c.Request.FormValue("password"), jwtUserRepository) {
		return
	}

	if isSet := SetPayloadFromDB(payload, jwtUserRepository); !isSet {
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError,
				"error": "Failed to set payload from DB."})
		c.Abort()
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

	var domainName string
	if domainName = "localhost"; os.Getenv("CONF_FILE") == "production.conf" {
		domainName = "hkpark130.p-e.kr"
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Domain:   domainName,
		Name:     "Authorization",
		Value:    fmt.Sprintf("Bearer %s", accessToken),
		Expires:  ExpiresCookie,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Domain:   domainName,
		Name:     "Refresh",
		Value:    fmt.Sprintf("%s", refreshToken),
		Expires:  ExpiresCookie,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	c.JSON(http.StatusOK, "OK")
}
