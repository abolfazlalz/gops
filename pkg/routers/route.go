package routers

import "github.com/gin-gonic/gin"

func StartRoutes(engine *gin.Engine) {
	execRoute := engine.Group("/exec")
	AddExecRoutes(execRoute)
}
