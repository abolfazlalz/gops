package main

import (
	"command-server/pkg/db"
	"command-server/pkg/routers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	log.Println("CLI Streaming http")

	if err := db.Migrate(); err != nil {
		panic(err)
	}

	app := gin.Default()

	routers.StartRoutes(app)

	if err := app.Run(":8080"); err != nil {
		panic(err)
	}
}
