package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestRestart(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	Test()
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	req := "Restarting project..."
	if !strings.Contains(string(out), req) {
		t.Fatalf("Expecting %s got %s", req, string(out))
	}
}
