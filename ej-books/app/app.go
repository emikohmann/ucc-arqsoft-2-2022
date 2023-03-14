package app

import "github.com/gin-gonic/gin"

func StartApp() {
	router := gin.Default()
	deps := BuildDependencies()
	MapUrls(router, deps)
	_ = router.Run(":8090")
}
