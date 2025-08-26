package main 

import (
	"github.com/gin-gonic/gin"
)

func main(){
	router := gin.Default()

	router.POST("/original-urls", submitURL)
	// TODO: router.GET("/originals-urls:param")

	router.Run("localhost:8080")

}