package initialize

import (
	"mxshop-api/user-web/router"

	"github.com/gin-gonic/gin"
)

func Routes() *gin.Engine {
	rout := gin.Default()
	ApiGroup := rout.Group("/v1")
	router.InitUserRouter(ApiGroup)
	return rout
}
