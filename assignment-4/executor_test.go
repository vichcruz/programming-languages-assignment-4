package main

import (
	"bytes"
	"context"
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
	err = RunTasks(context.Background(), cfg, 2, buf)

	if err == nil {
		t.Fatal("expected error for missing dependency, got nil")
	}

	if !strings.Contains(err.Error(), "DoesNotExist") {
		t.Fatalf("unexpected error: %v", err)
	}
}
