package main

import (
	"context"
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		limit int
		file  string
	)

	flag.IntVar(&limit, "max", 4, "Max concurrent tasks")
	flag.StringVar(&file, "file", "cli-tasks.json", "Config file")
	flag.Parse()

	f, err := os.Open(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open config file: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close config file: %v\n", err)
		}
	}()

	cfg, err := LoadConfig(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse config: %v\n", err)
		os.Exit(1)
	}

	mainCtx := context.Background()
	ctx, cancel := context.WithCancel(mainCtx)
	if err := RunTasks(ctx, cancel, cfg, limit, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error running tasks: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("All tasks completed successfully.")
}
