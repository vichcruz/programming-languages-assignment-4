package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// RunTasks executes all tasks concurrently, respecting dependencies
func RunTasks(ctx context.Context, cancel context.CancelFunc, cfg ConfigFile, maxWorkers int, w io.Writer) error {
	// Check for cyclic dependencies
	if err := ValidateNoCycles(cfg); err != nil {
		return err
	}

	taskDone := make(map[string]chan struct{})
	for name := range cfg {
		taskDone[name] = make(chan struct{})
	}

	// Semaphore for max concurrency
	sem := make(chan struct{}, maxWorkers)
	var mu sync.Mutex
	var wg sync.WaitGroup

	runTask := func(name string) error {
		task := cfg[name]

		// Print task header first
		mu.Lock()
		if _, err := fmt.Fprintf(w, "---- %s started ----\n", name); err != nil {
			mu.Unlock()
			return err
		}
		mu.Unlock()

		// Wait for dependencies
		for _, dep := range task.DEPS {
			if dep == "" {
				continue
			}
			depChan, ok := taskDone[dep]
			if !ok {
				close(taskDone[name])
				return fmt.Errorf("dependency %s not found for task %s", dep, name)
			}

			mu.Lock()
			if _, err := fmt.Fprintf(w, "[%s] waiting for dependency: %s\n", name, dep); err != nil {
				mu.Unlock()
				return err
			}
			mu.Unlock()

			select {
			case <-depChan:
			case <-ctx.Done():
				fmt.Println("ctx canceled inside runTask BEFORE command finished:", name)
				close(taskDone[name])
				return ctx.Err()
			}
		}

		// Acquire concurrency slot
		sem <- struct{}{}
		defer func() { <-sem }()

		// Prepare command
		parts := strings.Fields(task.CMD)
		if len(parts) == 0 {
			close(taskDone[name])
			return fmt.Errorf("task %s has empty command", name)
		}
		cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
		cmd.Dir = task.CWD
		cmd.Env = os.Environ()

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			close(taskDone[name])
			return fmt.Errorf("stdout pipe error for %s: %v", name, err)
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			close(taskDone[name])
			return fmt.Errorf("stderr pipe error for %s: %v", name, err)
		}

		if err := cmd.Start(); err != nil {
			close(taskDone[name])
			return fmt.Errorf("start error for %s: %v", name, err)
		}

		// Real-time output
		var outWg sync.WaitGroup
		outWg.Add(2)
		go func() {
			defer outWg.Done()
			if _, err := io.Copy(&writer{w, &mu, name}, stdout); err != nil {
				fmt.Fprintf(os.Stderr, "[%s] stdout copy error: %v\n", name, err)
			}
		}()
		go func() {
			defer outWg.Done()
			if _, err := io.Copy(&writer{w, &mu, name}, stderr); err != nil {
				fmt.Fprintf(os.Stderr, "[%s] stderr copy error: %v\n", name, err)
			}
		}()

		// Wait for command and output to finish
		err = cmd.Wait()
		outWg.Wait()

		mu.Lock()
		if err != nil {
			if _, err2 := fmt.Fprintf(w, "---- %s finished with error ----\n %v\n", name, err); err2 != nil {
				fmt.Fprintf(os.Stderr, "[%s] log write error: : %v\n", name, err2)
			}
			return err
		} else {
			if _, err2 := fmt.Fprintf(w, "---- %s finished successfully ----\n", name); err2 != nil {
				fmt.Fprintf(os.Stderr, "[%s] log write error: : %v\n", name, err2)
			}
		}
		mu.Unlock()

		close(taskDone[name])

		return nil
	}

	// Launch all tasks
	wg.Add(len(cfg))
	errCh := make(chan error, len(cfg))

	for name := range cfg {
		go func(n string) {
			defer wg.Done()
			if err := runTask(n); err != nil {
				cancel()
				errCh <- err
			}
		}(name)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}
