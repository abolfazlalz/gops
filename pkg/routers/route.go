package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func StartRoutes(engine *gin.Engine) {
	engine.Use(func(ctx *gin.Context) {
		defer func() {
			re := recover()
			if err, ok := re.(error); ok {
				ctx.String(http.StatusOK, err.Error())
			}
		}()
		ctx.Next()
	})
	execRoute := engine.Group("/exec")
	AddExecRoutes(execRoute)

	bashRoute := engine.Group("/bash")
	AddBashRoutes(bashRoute)
}
