package controllers

import (
	"net/http"
	"os"
	"partage-projets/config"
	"partage-projets/middlewares"
	"partage-projets/models"
	"partage-projets/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

// GetProjects godoc
// @Description Récupérer tous les projets
// @Tags Projects
// @Produce json
// @Success 200 {array} models.Project
// @Security BearerAuth
// @Router /projects [get]
func GetProjects(context *gin.Context) {
	var projects []models.Project

	if err := config.DB.Preload("Likes").Preload("Comments").Find(&projects).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch projects."})
		return
	}

	context.JSON(http.StatusOK, projects)
}

// GetProject godoc
// @Description Récupérer un projet par son ID
// @Tags Projects
// @Produce json
// @Param id path int true "ID du projet"
// @Success 200 {object} models.Project
// @Failure 400 {object} map[string]string "ID invalide"
// @Failure 404 {object} map[string]string "Projet non trouvé"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /projects/{id} [get]
func GetProject(context *gin.Context) {
	project, err := models.FindProjectById(context)

	if err == nil {
		context.JSON(http.StatusOK, project)
	}
}

// PostProject godoc
// @Description Créer un nouveau projet
// @Tags Projects
// @Accept json
// @Produce json
// @Param project body models.Project true "Données du projet"
// @Success 201 {object} models.Project
// @Failure 400 {object} map[string]string "Données invalides"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /projects [post]
func PostProject(context *gin.Context) {
	var project models.Project

	if err := context.ShouldBindJSON(&project); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data."})

		return
	}

	path, err := utils.UploadImage(context)
	if err != nil {
		return
	}

	if path != nil {
		project.Image = *path
	}

	if err := config.DB.Create(&project).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create project."})

		return
	}

	context.JSON(http.StatusCreated, project)
}

// PutProject godoc
// @Description Mettre à jour un projet existant
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path int true "ID du projet"
// @Param input body models.ProjectUpdateInput true "Données de mise à jour"
// @Success 200 {object} models.Project
// @Failure 400 {object} map[string]string "Données invalides"
// @Failure 404 {object} map[string]string "Projet non trouvé"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /projects/{id} [put]
func PutProject(context *gin.Context) {
	project, err := models.FindProjectById(context)

	if err == nil {
		var input models.ProjectUpdateInput
		if err = context.ShouldBindJSON(&input); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data."})

			return
		}

		updates := make(map[string]interface{})

		if input.Name != nil {
			updates["name"] = *input.Name
		}

		if input.Description != nil {
			updates["description"] = *input.Description
		}

		path, err := utils.UploadImage(context)
		if err != nil {
			return
		}

		if path != nil {
			if project.Image != "" {
				err = os.Remove(project.Image)

				if err != nil {
					context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete old image."})

					return
				}
			}

			updates["image"] = *path
		}

		if input.Skills != nil {
			updates["skills"] = datatypes.JSONSlice[string](*input.Skills)
		}

		if len(updates) == 0 {
			context.JSON(http.StatusBadRequest, gin.H{"error": "No data to update."})

			return
		}

		if err := config.DB.Model(&project).Updates(updates).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update project."})

			return
		}

		context.JSON(http.StatusOK, project)
	}
}

// DeleteProject godoc
// @Description Supprimer un projet
// @Tags Projects
// @Produce json
// @Param id path int true "ID du projet"
// @Success 200 {object} map[string]string "Message de succès"
// @Failure 400 {object} map[string]string "ID invalide"
// @Failure 404 {object} map[string]string "Projet non trouvé"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /projects/{id} [delete]
func DeleteProject(context *gin.Context) {
	project, err := models.FindProjectById(context)

	if err == nil {
		if err = config.DB.Delete(&project).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete project."})

			return
		}

		context.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully."})
	}
}

// LikeProject godoc
// @Description Liker ou déliker un projet
// @Tags Projects
// @Produce json
// @Param id path int true "ID du projet"
// @Success 200 {object} map[string]string "Message de succès"
// @Failure 400 {object} map[string]string "ID invalide"
// @Failure 404 {object} map[string]string "Projet non trouvé"
// @Failure 500 {object} map[string]string "Erreur interne"
// @Security BearerAuth
// @Router /projects/{id}/like [put]
func LikeProject(context *gin.Context) {
	var user models.User

	project, err := models.FindProjectById(context)

	if err == nil {
		userId := middlewares.GetUserId(context)

		if err := config.DB.First(&user, userId).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch user."})

			return
		}

		liked := false
		for _, u := range project.Likes {
			if u.ID == user.ID {
				liked = true
				break
			}
		}

		if liked {
			if err := config.DB.Model(&project).Association("Likes").Delete(&user); err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to unlike project."})

				return
			}

			context.JSON(http.StatusOK, gin.H{"message": "Project unliked successfully."})
		} else {
			if err := config.DB.Model(&project).Association("Likes").Append(&user); err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to like project."})

				return
			}

			context.JSON(http.StatusOK, gin.H{"message": "Project liked successfully."})
		}
	}
}
