package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type Body struct {
	Name string `json:"name"`
}

func main() {
	router := gin.Default()
	router.GET("/test", func(context *gin.Context) {
		body := Body{
			Name: os.Getenv("HOSTNAME"),
		}
		context.JSON(http.StatusAccepted, &body)
	})
	_ = router.Run(":3000")
}
