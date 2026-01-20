package routes

import (
	"partage-projets/controllers"
	"partage-projets/middlewares"

	"github.com/gin-gonic/gin"
)

func ProjectRoutes(router *gin.Engine) {
	routesGroup := router.Group("/projects")

	routesGroup.Use(middlewares.Authentication())

	{
		routesGroup.GET("/", controllers.GetProjects)
		routesGroup.GET("/:id", controllers.GetProject)
		routesGroup.POST("/", controllers.PostProject)
		routesGroup.PUT("/:id/like", controllers.LikeProject)
		routesGroup.PUT("/:id", controllers.PutProject)
		routesGroup.DELETE("/:id", controllers.DeleteProject)
	}
}
