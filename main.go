package main

import (
	"fmt"
	"log"
	"net/http"
	"partage-projets/config"
	"partage-projets/models"
	"partage-projets/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "partage-projets/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Partage de projets
// @version 1.0
// @description Description du projet de partage de projets
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
func main() {
	router := gin.Default()

	err := router.SetTrustedProxies(nil)
	if err != nil {
		log.Fatal("Unable to set trusted proxies: ", err)
	}

	router.Use(config.SecurityMiddleware())
	router.Use(config.CORSMiddleware())
	router.Use(config.RateLimit(100))

	err = godotenv.Load()
	if err != nil {
		// If .env file is not found, it is not necessarily an error.
		// With Render, environment variables are injected; there is no need for .env file.
		log.Print("Unable to find .env file: ", err)
	}

	router.GET("/status", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	routes.ProjectRoutes(router)
	routes.UserRoutes(router)
	routes.CommentRoutes(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	config.ConnectDB()

	err = config.DB.AutoMigrate(&models.Project{}, &models.User{}, &models.Comment{})
	if err != nil {
		log.Fatal("Unable to auto migrate: ", err)
	}

	err = router.Run(":8080")
	if err != nil {
		log.Fatal("Unable to start server: ", err)
	}

	fmt.Println("Server started on http://localhost:8080.")
}
