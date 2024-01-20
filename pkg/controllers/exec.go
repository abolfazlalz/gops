package controllers

import (
	"command-server/pkg/exec"
	"context"
	"github.com/gin-gonic/gin"
	"io"
)

type Exec struct {
}

func NewExec() *Exec {
	return &Exec{}
}

type RunExec struct {
	Dir     string   `json:"directory"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

func (e Exec) Run(c *gin.Context) {
	var data RunExec
	if err := c.BindJSON(&data); err != nil {
		panic(err)
	}

	c.Stream(func(w io.Writer) bool {
		ctx, cancel := context.WithCancel(c)
		defer cancel()
		if err := exec.Execute(
			ctx,
			w,
			data.Command,
			exec.WithArgs(data.Args...),
			exec.WithDir(data.Dir),
		); err != nil {
			return false
		}
		return false
	})
}
