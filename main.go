package main

import (
	"log"
	"net/http"

	"github.com/ajthomp2/ascii_art_service/asciiimageservice"
	"github.com/gin-gonic/gin"
)

func main() {
	service := asciiimageservice.New()

	r := gin.Default()

	r.POST("/image", func(c *gin.Context) {
		if c.Request.Body == nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "Message": "Image not in request"})
			return
		}

		imageId, err := service.SaveAsAscii(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "Message": "An unexpected error occurred"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"imageId": imageId})
	})

	r.GET("/images", func(c *gin.Context) {
		imageIds, err := service.ListAllIds()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "Message": "An unexpected error occurred"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"imageIds": imageIds})
	})

	r.GET("/images/:imageId", func(c *gin.Context) {
		if c.Param("imageId") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "Message": "Id of image must be specified in path param"})
			return
		}
		asciiimage, err := service.GetById(c.Param("imageId"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "Message": "An unexpected error occurred"})
			return
		}

		c.String(http.StatusOK, asciiimage.Data)
	})

	r.Run(":8080")
}
