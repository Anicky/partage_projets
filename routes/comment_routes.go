package routes

import (
	"partage-projets/controllers"
	"partage-projets/middlewares"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(router *gin.Engine) {
	routesGroup := router.Group("/comments")

	routesGroup.Use(middlewares.Authentication())

	{
		routesGroup.POST("/", controllers.PostComment)
	}
}
