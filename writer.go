package main

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

// writer wraps an io.Writer to prefix output with task name
type writer struct {
	w    io.Writer
	mu   *sync.Mutex
	name string
}

func (w *writer) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	lines := strings.Split(string(p), "\n")
	for i, line := range lines {
		if line != "" {
			if _, err := fmt.Fprintf(w.w, "[%s] %s\n", w.name, line); err != nil {
				return 0, err
			}
		} else if i < len(lines)-1 {
			if _, err := fmt.Fprintln(w.w, ""); err != nil {
				return 0, err
			}
		}
	}
	return len(p), nil
}
