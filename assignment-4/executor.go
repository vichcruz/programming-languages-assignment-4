package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sync/errgroup"
)

type Task struct {
	DESC string   `json:"desc"`
	CMD  string   `json:"cmd"`
	CWD  string   `json:"cwd"`
	DEPS []string `json:"deps"`
}

type ConfigFile map[string]Task

// Cycle detection method that uses DFS
func detectCycle(name string, config ConfigFile, visiting, visited map[string]bool) bool {
	if visiting[name] {
		return true // cycle found
	}
	if visited[name] {
		return false // already processed, no cycle
	}

	visiting[name] = true

	task := config[name]
	for _, dep := range task.DEPS {
		if detectCycle(dep, config, visiting, visited) {
			return true
		}
	}

	visiting[name] = false
	visited[name] = true
	return false
}

func LoadConfig(r io.Reader) (ConfigFile, error) {
	var cfg ConfigFile
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	return cfg, dec.Decode(&cfg)
}

func RunTasks(ctx context.Context, cfg ConfigFile, limit int, w io.Writer) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(limit)

	visiting := map[string]bool{}
	visited := map[string]bool{}

	for name := range cfg {
		if detectCycle(name, cfg, visiting, visited) {
			return fmt.Errorf("cyclic dependency detected")
		}
	}

	taskChans := make(map[string]chan struct{})
	for name := range cfg {
		taskChans[name] = make(chan struct{})
	}

	for name, task := range cfg {
		taskName := name
		configTask := task

		eg.Go(func() error {
			defer close(taskChans[taskName])

			// --- Wait on dependencies ---
			for _, dep := range configTask.DEPS {
				depChan, ok := taskChans[dep]
				if !ok {
					return fmt.Errorf("dependency %s not found for task %s", dep, taskName)
				}

				select {
				case <-depChan:
				case <-ctx.Done():
					return ctx.Err()
				}
			}

			parts := strings.Fields(configTask.CMD)
			cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
			cmd.Dir = configTask.CWD
			cmd.Env = os.Environ()

			out, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("task %s failed: %w", taskName, err)
			}

			if _, err := fmt.Fprintf(w, "---- %s ----\n", taskName); err != nil {
				return err
			}

			if _, err := fmt.Fprintf(w, "Output:\n%s\n", string(out)); err != nil {
				return err
			}

			return nil
		})
	}

	return eg.Wait()
}

func main() {
	var (
		limit int
		file  string
	)

	flag.IntVar(&limit, "max", 4, "Max goroutines")
	flag.StringVar(&file, "file", "cli-tasks.json", "Config file")
	flag.Parse()

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close file: %v\n", err)
		}
	}()

	cfg, err := LoadConfig(f)
	if err != nil {
		panic(err)
	}

	if err := RunTasks(context.Background(), cfg, limit, os.Stdout); err != nil {
		fmt.Println("Error:", err)
	}
}
