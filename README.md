# CLI Task Runner

A concurrent task execution tool written in Go that manages dependencies between tasks and executes them efficiently.

## Overview

This CLI task runner reads a JSON configuration file defining tasks with their commands, working directories, and dependencies. It executes tasks concurrently while respecting dependency ordering, with configurable parallelism limits.

### Features

- **Concurrent Execution**: Run multiple independent tasks simultaneously
- **Dependency Management**: Automatically handles task dependencies and execution order
- **Cyclic Dependency Detection**: Validates configuration to prevent infinite loops
- **Real-time Output**: Streams task output with prefixed task names for easy identification
- **Configurable Concurrency**: Control maximum parallel tasks with the `-max` flag
- **Error Handling**: Gracefully handles missing dependencies and command failures

## Setup

### Prerequisites

- Go 1.16 or higher
- `golangci-lint` (optional, for running the lint task in the example config)

### Installation

1. Clone or download this repository:

```bash
cd /path/to/programming-languages-assignment-4
```

2. Build the project:

```bash
go build -o cli-task-runner
```

## Configuration

Create a JSON configuration file (default: `cli-tasks.json`) with the following structure:

```json
{
  "task-name": {
    "desc": "Task description",
    "cmd": "command to execute",
    "cwd": "./working/directory",
    "deps": ["dependency1", "dependency2"]
  }
}
```

### Configuration Fields

- `desc`: Human-readable description of the task
- `cmd`: Shell command to execute
- `cwd`: Working directory for the command (relative or absolute path)
- `deps`: Array of task names that must complete before this task runs (use empty array `[]` for no dependencies)

### Example Configuration

See `cli-tasks.json` for a complete example with build, test, lint, clean, and print tasks.

## Usage

### Basic Usage

Run tasks with default settings (max 4 concurrent tasks, using `cli-tasks.json`):

```bash
./cli-task-runner
```

### Command-Line Options

```bash
./cli-task-runner [options]
```

**Options:**

- `-max <number>`: Maximum number of concurrent tasks (default: 4)
- `-file <path>`: Path to configuration file (default: "cli-tasks.json")

### Examples

Run with custom concurrency limit:

```bash
./cli-task-runner -max 2
```

Use a different configuration file:

```bash
./cli-task-runner -file my-tasks.json
```

Combine options:

```bash
./cli-task-runner -max 8 -file build-tasks.json
```

## Running Tests

Execute the test suite:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

## Output Format

Tasks display output in the following format:

```
---- task-name started ----
[task-name] waiting for dependency: dependency-name
[task-name] command output line 1
[task-name] command output line 2
---- task-name finished successfully ----
```

Errors are reported as:

```
---- task-name finished with error ----
 error details
```

## Error Handling

The runner will detect and report:

- Missing task dependencies
- Cyclic dependencies between tasks
- Command execution failures
- Invalid JSON configuration
- Unknown fields in configuration

## Project Structure

```
.
├── main.go           # Entry point and CLI argument parsing
├── config.go         # Configuration loading and cycle detection
├── runner.go         # Task execution and dependency management
├── writer.go         # Custom writer for prefixed output
├── executor_test.go  # Test suite
├── cli-tasks.json    # Example configuration file
└── README.md         # This file
```
