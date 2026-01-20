package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"partage-projets/config"
	"partage-projets/controllers"
	"partage-projets/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal("Unable to setup database: ", err)
	}

	err = db.AutoMigrate(&models.Project{}, &models.User{}, &models.Comment{})
	if err != nil {
		log.Fatal("Unable to migrate database: ", err)
	}

	project := models.Project{
		Name:        "Test project",
		Description: "Test description",
	}
	db.Create(&project)

	comment := models.Comment{
		ProjectID: project.ID,
		Content:   "Test comment",
	}
	db.Create(&comment)

	return db
}

func TestGetProjects(testing *testing.T) {
	gin.SetMode(gin.TestMode)
	config.DB = setupTestDatabase()

	router := gin.Default()
	router.GET("/projects", controllers.GetProjects)

	request, err := http.NewRequest(http.MethodGet, "/projects", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "Test project")
	assert.Contains(testing, body, "Test comment")
}

func TestPostProject(testing *testing.T) {
	gin.SetMode(gin.TestMode)
	config.DB = setupTestDatabase()

	router := gin.Default()
	router.POST("/projects", controllers.PostProject)

	project := map[string]interface{}{
		"name":        "Test project",
		"description": "Test description",
		"skills":      []string{"Go", "Testing", "SQLite"},
	}

	data, err := json.Marshal(project)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/projects", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusCreated, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "Test project")
	assert.Contains(testing, body, "Test description")
	assert.Contains(testing, body, "Testing")
}
