package utils

import (
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

func UploadImage(context *gin.Context) (*string, error) {
	file, err := context.FormFile("image")

	if err == nil {
		path := "uploads/" + file.Filename
		if err := context.SaveUploadedFile(file, path); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to upload image."})

			return nil, err
		}

		img, _ := imaging.Open(path)

		resized := imaging.Resize(img, 800, 0, imaging.Lanczos)
		if err := imaging.Save(resized, path); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to resize image."})

			return nil, err
		}

		return &path, nil
	}

	return nil, nil
}
