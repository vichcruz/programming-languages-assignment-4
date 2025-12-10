package main

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMissingDeps(t *testing.T) {
	json := `{
        "clean": {
            "desc": "Clean test",
            "cmd": "echo hi",
            "cwd": ".",
            "deps": ["DoesNotExist"]
        }
    }`

	cfg, err := LoadConfig(strings.NewReader(json))
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	buf := &bytes.Buffer{}
	ctx, cancel := context.WithCancel(context.Background())
	runErr := RunTasks(ctx, cancel, cfg, 1, buf)

	if runErr == nil {
		t.Fatal("expected error for missing dependency, got nil")
	}

	if !strings.Contains(runErr.Error(), "dependency DoesNotExist not found") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCyclicDeps(t *testing.T) {
	json := `{
        "clean": {
            "desc": "Clean test",
            "cmd": "echo hi",
            "cwd": ".",
            "deps": ["build"]
        },
        "build": {
            "desc": "Clean test",
            "cmd": "echo hi",
            "cwd": ".",
            "deps": ["test"]
        },
        "test": {
            "desc": "Clean test",
            "cmd": "echo hi",
            "cwd": ".",
            "deps": ["clean"]
        }
    }`

	cfg, err := LoadConfig(strings.NewReader(json))
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	buf := &bytes.Buffer{}
	ctx, cancel := context.WithCancel(context.Background())
	runErr := RunTasks(ctx, cancel, cfg, 1, buf)

	if runErr == nil {
		t.Fatal("expected error for cyclic dependencies, got nil")
	}

	if !strings.Contains(runErr.Error(), "cyclic dependency detected") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	_, err := LoadConfig(strings.NewReader("{invalid-json}"))
	if err == nil {
		t.Fatal("expected JSON parse error, got nil")
	}
}

func TestRunTasksSuccess(t *testing.T) {
	json := `{
        "a": {
            "desc": "Task A",
            "cmd": "echo A",
            "cwd": ".",
            "deps": []
        },
        "b": {
            "desc": "Task B",
            "cmd": "echo B",
            "cwd": ".",
            "deps": ["a"]
        }
    }`

	cfg, err := LoadConfig(strings.NewReader(json))
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	buf := &bytes.Buffer{}
	ctx, cancel := context.WithCancel(context.Background())
	if err := RunTasks(ctx, cancel, cfg, 1, buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()

	// Check headers
	if !strings.Contains(out, "---- a started ----") {
		t.Fatalf("expected task header for 'a', got: %s", out)
	}

	if !strings.Contains(out, "---- b started ----") {
		t.Fatalf("expected task header for 'b', got: %s", out)
	}

	// Check actual command output
	if !strings.Contains(out, "A") {
		t.Fatalf("expected output from task a, got: %s", out)
	}

	if !strings.Contains(out, "B") {
		t.Fatalf("expected output from task b, got: %s", out)
	}
}

func TestRunTasksOrder(t *testing.T) {
	json := `{
        "first": {
            "desc": "First",
            "cmd": "echo first",
            "cwd": ".",
            "deps": []
        },
        "second": {
            "desc": "Second",
            "cmd": "echo second",
            "cwd": ".",
            "deps": ["first"]
        }
    }`

	cfg, _ := LoadConfig(strings.NewReader(json))
	buf := &bytes.Buffer{}
	ctx, cancel := context.WithCancel(context.Background())
	if err := RunTasks(ctx, cancel, cfg, 1, buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	firstIndex := strings.Index(out, "---- first finished successfully ----")
	secondIndex := strings.Index(out, "---- second finished successfully ----")

	if secondIndex < firstIndex {
		t.Fatalf("dependency order violated:\n%s", out)
	}
}

func TestRunMultiTasksOrder(t *testing.T) {
	json := `{
        "first": {
            "desc": "First",
            "cmd": "echo first",
            "cwd": ".",
            "deps": []
        },
        "second": {
            "desc": "Second",
            "cmd": "echo second",
            "cwd": ".",
            "deps": ["first"]
        },
        "third": {
            "desc": "Second",
            "cmd": "echo third",
            "cwd": ".",
            "deps": ["first", "second"]
        }
    }`

	cfg, _ := LoadConfig(strings.NewReader(json))
	buf := &bytes.Buffer{}
	ctx, cancel := context.WithCancel(context.Background())
	if err := RunTasks(ctx, cancel, cfg, 1, buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()

	firstIndex := strings.Index(out, "---- first finished successfully ----")
	secondIndex := strings.Index(out, "---- second finished successfully ----")
	thirdIndex := strings.Index(out, "---- third finished successfully ----")

	if (secondIndex < firstIndex) || (thirdIndex < firstIndex) {
		t.Fatalf("dependency order violated:\n%s", out)
	}
}

func TestRunTasksWithCwd(t *testing.T) {
	dir := t.TempDir()

	// create a file in the temp directory
	path := filepath.Join(dir, "hello.txt")
	if err := os.WriteFile(path, []byte("hi"), fs.ModeTemporary); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	json := fmt.Sprintf(`{
        "test": {
            "desc": "test",
            "cmd": "ls",
            "cwd": "%s",
            "deps": []
        }
    }`, dir)

	cfg, _ := LoadConfig(strings.NewReader(json))
	buf := &bytes.Buffer{}
	ctx, cancel := context.WithCancel(context.Background())
	if err := RunTasks(ctx, cancel, cfg, 1, buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()

	if !strings.Contains(out, "hello.txt") {
		t.Fatalf("expected ls to list hello.txt, got: %s", out)
	}
}
