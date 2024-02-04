package controllers

import (
	"command-server/pkg/db"
	"command-server/pkg/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type Bash struct {
}

func NewBash() *Bash {
	return &Bash{}
}

type uploadData struct {
	Title string `form:"title" binding:"required"`
}

func (ctrl Bash) Upload(c *gin.Context) {
	file, _ := c.FormFile("file")

	var data uploadData
	if err := c.Bind(&data); err != nil {
		panic(err)
	}

	ext := ""
	fileParts := strings.Split(file.Filename, ".")

	if len(fileParts) > 1 {
		ext = "." + fileParts[1]
	}

	dst := fmt.Sprintf("uploads/%d%s", time.Now().Unix(), ext)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		panic(err)
	}

	bash := models.Bash{
		Title: data.Title,
		File:  dst,
	}

	database, _ := db.Database()
	database.Create(&bash)

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded\n", file.Filename))
}
