package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Task struct {
	DESC string   `json:"desc"`
	CMD  string   `json:"cmd"`
	CWD  string   `json:"cwd"`
	DEPS []string `json:"deps"`
}

type ConfigFile map[string]Task

func main()  {
	// TODO: add a dag analysis to check for cyclic dependenciees first

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

	// TODO: use errGroup.Group or errorChannel to handle errors and cancel tasks if any dependecy failed
	var wg sync.WaitGroup

	// Mapping of task name to channel
	taskChans := make(map[string]chan struct{})

	// Make sure channels exist before running the task goroutines
	for name := range cfg {
		// Create a channel for each task and map it to its name
		ch := make(chan struct{})
		taskChans[name] = ch
	}

	for configName, configTask := range cfg {
		// Start goroutine for each task inside a waitGroup
		wg.Add(1)
		go func(name string, task Task) {
			defer wg.Done()
			fmt.Printf("Task Name: %s\nDescription: %s\n", name, task.DESC)

			// Wait for dependencies to complete
			for _, dep := range task.DEPS {
				depChan, exists := taskChans[dep]
				if !exists {
					fmt.Printf("Dependency %s not found for task %s\n", dep, name)
					return
				}
				// Wait for dependency to finish
				<-depChan
			}
			// Split CMD into command + args
			parts := strings.Fields(task.CMD)
			// Create command to execute
			cmd := exec.Command(parts[0], parts[1:]...)
			// Load environment variables and set working directory
			cmd.Env = os.Environ()
			cmd.Dir = task.CWD
		
			// Execute command and capture combined output (stdout + stderr)
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("command failed: %v\noutput: \n%s\n", err, string(out))
				return
			}
			// Print command output
			fmt.Printf("Output from task %s:\n%s\n",name, string(out))

			// Signal task completion by closing channel
			close(taskChans[name])
		}(configName, configTask)
	}

	wg.Wait()
}