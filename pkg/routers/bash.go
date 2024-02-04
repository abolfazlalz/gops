package routers

import (
	"command-server/pkg/controllers"
	"github.com/gin-gonic/gin"
)

func AddBashRoutes(route *gin.RouterGroup) {
	exec := controllers.NewBash()
	route.POST("/upload", exec.Upload)
}
