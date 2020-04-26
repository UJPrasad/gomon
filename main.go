package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/fsnotify/fsnotify"
)

var ctx context.Context
var cancel context.CancelFunc

func runAndPrint(gorun string) {
	ctx, cancel = context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "bash", "-c", gorun)
	mwriter := io.MultiWriter(os.Stdout)
	cmd.Stdout = mwriter
	cmd.Stderr = mwriter
	cmd.Run()
}

// Test to check..
func Test() {
	go func() {
		time.Sleep(500 * time.Millisecond)
		l, err := ioutil.ReadFile("test/main.go")
		if err != nil {
			fmt.Println(err, "read")
		}
		err = ioutil.WriteFile("test/main.go", l, 0644)
		if err != nil {
			fmt.Println(err, "write")
		}
	}()
	runInfinite("./test/", "3000", "main.go", true)
}

func runInfinite(path, port, main string, t bool) {
	gorun := "go run " + path + main
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()
	done := make(chan bool)
	quit := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					if cancel != nil {
						cancel()
					}
					s := fmt.Sprintf("lsof -i tcp:%s | awk 'NR!=1 {print $2}' | xargs kill", port)
					x := exec.Command("bash", "-c", s)
					x.Run()
					go runAndPrint(gorun)
					fmt.Println("Restarting project...")
					time.Sleep(500 * time.Millisecond)
				}

			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			case <-quit:
				if cancel != nil {
					cancel()
				}
				s := fmt.Sprintf("lsof -i tcp:%s | awk 'NR!=1 {print $2}' | xargs kill", port)
				x := exec.Command("bash", "-c", s)
				x.Run()
				done <- true
				return
			}
		}
	}()
	if err := watcher.Add(path); err != nil {
		fmt.Println("ERROR", err)
	}
	go runAndPrint(gorun)

	go func() {
		time.Sleep(600 * time.Millisecond)
		if t {
			quit <- true
		}
	}()
	<-done
}

func main() {
	path := flag.String("path", "", "--path \"/<abs-path>/<to>/<project>\"")
	port := flag.String("port", "", "--port \"<port>\"")
	main := flag.String("main", "app.go", "--main \"<file-name.go>\"")
	flag.Parse()
	if *path == "" || *port == "" {
		panic("path and port cannot be emtpy")
	}
	runInfinite(*path, *port, *main, false)
}
