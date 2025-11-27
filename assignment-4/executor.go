package main

import (
	"context"
	"encoding/json"
	"fmt"
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

func main() {
	f, err := os.Open("cli-tasks.json")
	if err != nil {
		fmt.Printf("Reading file error %s\n", err)
		return
	}

	// Defer file close to ensure file is closed after reading regardless of success or failure
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}()

	// Decode JSON config according to our types
	var cfg ConfigFile
	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields() // helps catch typos in JSON keys
	if err := dec.Decode(&cfg); err != nil {
		fmt.Println("json decode error:", err)
		return
	}

	// Create visiting and visited maps for cycle detection
	visiting := map[string]bool{}
	visited := map[string]bool{}

	for name := range cfg {
		if detectCycle(name, cfg, visiting, visited) {
			fmt.Println("Cyclic dependency detected!")
			return
		}
	}

	// TODO: check errGroup number of goroutines limitations
	// This creates a new errorGroup with a cancelable context derived from its parent (here: context.Background())
	eg, ctx := errgroup.WithContext(context.Background())

	// Mapping of task name to channel
	taskChans := make(map[string]chan struct{})

	// Make sure channels exist before running the task goroutines
	for name := range cfg {
		// Create a channel for each task and map it to its name
		ch := make(chan struct{})
		taskChans[name] = ch
	}

	for configName, configTask := range cfg {
		// Create local copies to avoid loop variable capture issue
		name := configName
		task := configTask
		// Start goroutine for each task inside an errgroup
		eg.Go(func() error {
			fmt.Printf("Task Name: %s\nDescription: %s\n", name, task.DESC)

			// Wait for dependencies to complete
			for _, dep := range task.DEPS {
				depChan, exists := taskChans[dep]
				if !exists {
					return fmt.Errorf("dependency %s not found for task %s", dep, configName)
				}
				// Wait for dependency to finish
				// This blocks further execution until we get a value or depChan is closed
				// In case it was closed, the zero value of the type inside the channel is returned and execution can continue
				// Reads on a closed channel proceed immediately but a send would panic
				// If context was cancelled due to some error during execution, the task is aborted
				select {
				case <-depChan:
					// dependency finished normally
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			// Split CMD into command + args
			parts := strings.Fields(task.CMD)
			// Create command to execute
			cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
			// Load environment variables and set working directory
			cmd.Env = os.Environ()
			cmd.Dir = task.CWD

			// Execute command and capture combined output (stdout + stderr)
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("command failed: %v\noutput: \n%s\n", err, string(out))
				return err
			}
			// Print command output
			fmt.Printf("Output from task %s:\n%s\n", name, string(out))

			// Signal task completion by closing channel
			defer close(taskChans[name])

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
