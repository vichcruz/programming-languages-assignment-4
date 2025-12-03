package main

import (
	"encoding/json"
	"fmt"
	"io"
)

// Task represents a single CLI task
type Task struct {
	DESC string   `json:"desc"`
	CMD  string   `json:"cmd"`
	CWD  string   `json:"cwd"`
	DEPS []string `json:"deps"`
}

// ConfigFile is a map of task names to Tasks
type ConfigFile map[string]Task

// LoadConfig reads JSON config and validates unknown fields
func LoadConfig(r io.Reader) (ConfigFile, error) {
	var cfg ConfigFile
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	return cfg, dec.Decode(&cfg)
}

// detectCycle returns true if a cyclic dependency exists
func detectCycle(name string, cfg ConfigFile, visiting, visited map[string]bool) bool {
	if visiting[name] {
		return true
	}
	if visited[name] {
		return false
	}

	visiting[name] = true
	task := cfg[name]
	for _, dep := range task.DEPS {
		if detectCycle(dep, cfg, visiting, visited) {
			return true
		}
	}
	visiting[name] = false
	visited[name] = true
	return false
}

// ValidateNoCycles checks all tasks for cyclic dependencies
func ValidateNoCycles(cfg ConfigFile) error {
	visiting := map[string]bool{}
	visited := map[string]bool{}
	for name := range cfg {
		if detectCycle(name, cfg, visiting, visited) {
			return fmt.Errorf("cyclic dependency detected at task %s", name)
		}
	}
	return nil
}
