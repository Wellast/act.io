package gin

import (
	_ "net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/api/v1/deployment", getDeployment)
	r.POST("/api/v1/deployment", createDeployment)
	r.DELETE("/api/v1/deployment/:name", deleteDeployment)

	return r
}

func RunGin() error {
	r := setupRouter()
	return r.Run(":8080")
}
