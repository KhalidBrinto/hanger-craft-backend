package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddBannerImages(c *gin.Context) {
	var payload struct {
		Position string
		Image    []string
	}

	// Bind the incoming JSON to the Category struct
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var contents []models.ContentImage
	for _, image := range payload.Image {
		contents = append(contents, models.ContentImage{
			Position: payload.Position,
			Image:    image,
		})

	}

	// Insert the category into the database
	if err := config.DB.Create(&contents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the newly created category
	c.JSON(http.StatusCreated, gin.H{"message": "content added successfully"})
}

// api to fetch banner contents
func GetBannerImages(c *gin.Context) {
	response := map[string]interface{}{
		"left_banner":    []string{},
		"right_banner_1": []string{},
		"right_banner_2": []string{},
	}

	var content []*models.ContentImage

	config.DB.Model(&content).Find(&content).Order("id DESC")

	if len(content) != 0 {
		for _, image := range content {
			switch image.Position {
			case "left_banner":
				response["left_banner"] = append(response["left_banner"].([]string), image.Image)
			case "right_banner_1":
				response["right_banner_1"] = append(response["right_banner_1"].([]string), image.Image)
			case "right_banner_2":
				response["right_banner_2"] = append(response["right_banner_2"].([]string), image.Image)
			}

		}
	}
	c.JSON(http.StatusOK, response)
}
