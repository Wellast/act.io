package gin

import (
	"controller/config"
	_ "net/http"

	"github.com/gin-gonic/gin"
)

var appConfig config.IConf

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/api/v1/deployment", getDeployment)
	r.POST("/api/v1/deployment", createDeployment)
	r.DELETE("/api/v1/deployment/:namespace", deleteDeployment)

	return r
}

func RunGin(conf config.IConf) error {
	appConfig = conf

	r := setupRouter()
	return r.Run(":8080")
}
