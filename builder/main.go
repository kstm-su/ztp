package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func postImage(c *gin.Context) {
	image := Image{}
	if err := c.Bind(&image); err != nil {
		log.Printf("[ERROR] %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Printf("[DEBUG] form: %+v\n", image)
	image.Append()
	c.JSON(http.StatusOK, image)
}

func main() {
	r := gin.Default()

	images := r.Group("/images")
	images.POST("", postImage)

	port := os.Getenv("PORT")
	r.Run(":" + port)
}
