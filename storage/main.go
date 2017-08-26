package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func postImage(c *gin.Context) {
	image := Image{}
	log.Printf("[DEBUG] image: %+v\n", image)
	if err := c.Bind(&image); err != nil {
		log.Printf("[ERROR] %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := image.Build(); err != nil {
		log.Printf("[ERROR] %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, image)
}

func main() {
	r := gin.Default()

	images := r.Group("/images")
	images.POST("", postImage)

	port := os.Getenv("PORT")
	r.Run(":" + port)
}
