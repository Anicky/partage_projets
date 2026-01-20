package tests

import (
	"log"
	"net/http"
	"os"
	"partage-projets/config"
	"partage-projets/models"
	"partage-projets/routes"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	config.DB = setupTestDatabase()

	router := gin.Default()

	routes.ProjectRoutes(router)
	routes.UserRoutes(router)
	routes.CommentRoutes(router)

	return router
}

func AuthenticateUser(request *http.Request) {
	token := generateTestToken(1)

	request.Header.Set("Authorization", "Bearer "+token)
}

func setupTestDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal("Unable to setup database: ", err)
	}

	err = db.AutoMigrate(&models.Project{}, &models.User{}, &models.Comment{})
	if err != nil {
		log.Fatal("Unable to migrate database: ", err)
	}

	project1 := models.Project{
		Name:        "Test project 1",
		Description: "Test description 1",
	}
	db.Create(&project1)

	project2 := models.Project{
		Name:        "Test project 2",
		Description: "Test description 2",
	}
	db.Create(&project2)

	user := models.User{
		Email:    "user1@example.com",
		Password: "Password123!",
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Unable to hash password: ", err)
	}

	user.Password = string(hashedPassword)
	db.Create(&user)

	comment := models.Comment{
		ProjectID: project1.ID,
		Content:   "Test comment on project 1",
	}
	db.Create(&comment)

	return db
}

func generateTestToken(userID uint) string {
	err := os.Setenv("JWT_SECRET", "test_secret")
	if err != nil {
		log.Fatal("Unable to set JWT_SECRET environment variable: ", err)
	}

	claims := jwt.MapClaims{
		"UserID": float64(userID),
		"exp":    time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString([]byte("test_secret"))

	return tokenString
}
