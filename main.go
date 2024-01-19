package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func handleScannerStd(std io.ReadCloser, textCh chan string) {
	scanner := bufio.NewScanner(std)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		textCh <- m
	}
}

func handleCmd(ctx context.Context, w io.Writer, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	runCh := make(chan error)

	go func() {
		fmt.Println("ok")
		runCh <- cmd.Run()
		fmt.Println("not ok")
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
			fmt.Fprintf(w, "Data chunk: %s\n", val)
			flusher.Flush()
			break
		case <-ctx.Done():
			isDone = false
			fmt.Fprintf(w, "Done")
			flusher.Flush()
			break
		case <-runCh:
			isDone = false
			fmt.Fprintf(w, "Done")
			flusher.Flush()
			break
		}
	}
	return nil
}

func executeCommandResponse(response http.ResponseWriter, req *http.Request, args ...string) {
	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	go handleCmd(ctx, response, args[0], args[1:]...)

	select {
	case <-ctx.Done():
		fmt.Fprintln(response, "Request canceled")
	}
}

func hello(response http.ResponseWriter, req *http.Request) {
	text, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	body := fmt.Sprintf("%s", text)
	args := strings.Split(body, " ")
	executeCommandResponse(response, req, args...)
}

func main() {
	log.Println("CLI Streaming http")
	http.HandleFunc("/", hello)
	log.Println("Start server")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
