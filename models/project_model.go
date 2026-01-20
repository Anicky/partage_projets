package models

import (
	"errors"
	"net/http"
	"partage-projets/config"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Project struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string `binding:"required"`
	Description string `binding:"required"`
	Image       string
	Skills      datatypes.JSONSlice[string] `gorm:"type:json" swaggertype:"array,string"`
	Comments    []Comment                   `gorm:"foreignKey:ProjectID"`
	Likes       []User                      `gorm:"many2many:project_likes"`
}

type ProjectUpdateInput struct {
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	Skills      *[]string `json:"skills"`
}

func FindProjectById(context *gin.Context) (project *Project, err error) {
	idParam := context.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID."})

		return nil, err
	}

	if err = config.DB.Preload("Likes").Preload("Comments").First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"error": "Project not found."})

			return nil, err
		}

		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch project."})

		return nil, err
	}

	return project, nil
}
