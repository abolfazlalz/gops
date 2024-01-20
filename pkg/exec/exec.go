package exec

import (
	"bufio"
	"command-server/pkg/db"
	"command-server/pkg/models"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func handleScannerStd(std io.ReadCloser, textCh chan string) {
	scanner := bufio.NewScanner(std)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		textCh <- m
	}
}

type Option struct {
	dir     string
	command string
	args    []string
}

type OptionCallback func(*Option)

func WithArgs(args ...string) OptionCallback {
	return func(option *Option) {
		option.args = args
	}
}

func WithDir(dir string) OptionCallback {
	return func(option *Option) {
		if dir != "" {
			option.dir = dir
		} else {
			option.dir = "~"
		}
	}
}

func Execute(ctx context.Context, w io.Writer, command string, callbacks ...OptionCallback) (err error) {
	option := Option{}

	for _, callback := range callbacks {
		callback(&option)
	}

	cmd := exec.Command(command, option.args...)
	cmd.Dir = option.dir

	startTime := time.Now()

	cmdLog := models.CommandLog{
		Command: command,
		Args:    strings.Join(option.args, " "),
	}

	database, err := db.Database()
	if err != nil {
		return err
	}

	if err := database.Create(&cmdLog).Error; err != nil {
		return err
	}

	defer func() {
		endTime := time.Now()
		cmdLog.Duration = endTime.Sub(startTime)
		cmdLog.Completed = true

		err = database.Updates(&cmdLog).Error
		return
	}()

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	runCh := make(chan error)

	go func() {
		runCh <- cmd.Run()
	}()

	flusher, ok := w.(http.Flusher)
	if !ok {
		return errors.New("streaming not supported")
	}

	textCh := make(chan string)
	go handleScannerStd(stdout, textCh)
	go handleScannerStd(stderr, textCh)

	isDone := true
	for isDone {
		select {
		case val := <-textCh:
			if _, err := fmt.Fprintf(w, "%s\n", val); err != nil {
				return err
			}
			flusher.Flush()
			break
		case <-runCh:
			isDone = false
			if _, err := fmt.Fprintf(w, "Done\n"); err != nil {
				return err
			}
			flusher.Flush()
			break
		case <-ctx.Done():
			fmt.Println("canceled")
			return nil
		}
	}
	return cmd.Wait()
}
