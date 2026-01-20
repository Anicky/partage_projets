package controllers

import (
	"net/http"
	"partage-projets/config"
	"partage-projets/middlewares"
	"partage-projets/models"

	"github.com/gin-gonic/gin"
)

// PostComment godoc
// @Description Ajouter un commentaire à un projet
// @Tags Comments
// @Accept json
// @Produce json
// @Param comment body models.Comment true "Données du commentaire"
// @Success 201 {object} models.Comment
// @Failure 400 {object} map[string]string "Données invalides"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /comments [post]
func PostComment(context *gin.Context) {
	var comment models.Comment

	if err := context.ShouldBindJSON(&comment); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data."})

		return
	}

	comment.UserID = *middlewares.GetUserId(context)

	if err := config.DB.Create(&comment).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create comment."})

		return
	}

	context.JSON(http.StatusCreated, comment)
}
