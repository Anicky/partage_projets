package controllers

import (
	"net/http"
	"os"
	"partage-projets/config"
	"partage-projets/models"
	"partage-projets/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaim struct {
	UserID uint
	jwt.RegisteredClaims
}

func Login(context *gin.Context) {
	var user models.User

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data."})

		return
	}

	var existingUser models.User
	if err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password."})

		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password."})

		return
	}

	claim := &CustomClaim{
		UserID: existingUser.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate token."})

		return
	}

	context.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Register(context *gin.Context) {
	var user models.User

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data."})

		return
	}

	var count int64
	config.DB.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)

	if count > 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Email already used."})

		return
	}

	if err := utils.ValidatePassword(user.Password); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to hash password."})

		return
	}

	user.Password = string(hashedPassword)

	if err := config.DB.Create(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create user."})

		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User created successfully."})
}
