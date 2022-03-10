package handlers

import (
	"golang/jwt/api/domain"
	"golang/jwt/api/repository"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context, jwtUserRepository *repository.JwtUserRepository) {
	var jwtUser domain.JwtUser

	if err := c.Bind(&jwtUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
		log.Printf("Wrong parameter. %s ", err)
		c.Abort()
		return
	}

	hashedPassword, err := HashPassword(c.Request.FormValue("password"))
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError,
				"error": "Failed to hash password."})
		c.Abort()
		return
	}
	jwtUser.Password = hashedPassword

	err = jwtUserRepository.CreateUser(&jwtUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
		log.Printf("Failed to sign up. %s ", err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, "success")
}

func GetUserByIDHandler(c *gin.Context, jwtUserRepository *repository.JwtUserRepository) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
		log.Printf("Wrong parameter. %s ", err)
		c.Abort()
		return
	}

	user, err := jwtUserRepository.GetUserByID(userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
		log.Printf("Failed to read data. %s ", err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetUsersHandler(c *gin.Context, jwtUserRepository *repository.JwtUserRepository) {
	user, err := jwtUserRepository.GetUsers()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
		log.Printf("Failed to read data. %s ", err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateHandler(c *gin.Context, jwtUserRepository *repository.JwtUserRepository) {
	userID, err := strconv.ParseUint(c.Request.FormValue("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
		log.Printf("Wrong parameter. %s ", err)
		c.Abort()
		return
	}

	hashedPassword, err := HashPassword(c.Request.FormValue("password"))
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError,
				"error": "Failed to hash password."})
		c.Abort()
		return
	}
	user := &domain.JwtUser{Email: c.Request.FormValue("email"), Password: hashedPassword}

	result, err := jwtUserRepository.UpdateUser(userID, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
		log.Printf("Failed to update. %s ", err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, result)
}

func DeleteHandler(c *gin.Context, jwtUserRepository *repository.JwtUserRepository) {
	userID, err := strconv.ParseUint(c.Request.FormValue("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
		log.Printf("Wrong parameter. %s ", err)
		c.Abort()
		return
	}

	err = jwtUserRepository.DeleteUserByID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
		log.Printf("Failed to update. %s ", err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, "success")
}
