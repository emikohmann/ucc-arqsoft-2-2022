package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func MapUrls(router *gin.Engine, dependencies *Dependencies) {
	// Products Mapping

	router.GET("/books/:id", dependencies.BookController.Get)
	router.POST("/books", dependencies.BookController.Insert)

	fmt.Println("Finishing mappings configurations")
}
