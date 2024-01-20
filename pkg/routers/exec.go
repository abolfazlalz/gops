package routers

import (
	"command-server/pkg/controllers"
	"github.com/gin-gonic/gin"
)

func AddExecRoutes(route *gin.RouterGroup) {
	exec := controllers.NewExec()
	route.POST("/", exec.Run)
}
