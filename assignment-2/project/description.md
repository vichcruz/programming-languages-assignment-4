# Project Proposal: Concurrent CLI Task Runner in Go

**Author:** Victor Cruz da Silva  
**Email:** victorhugo.cruzdasilva@uzh.ch  
**Language:** Go

---

## Project Title and Brief Description

**Title:** Concurrent CLI Task Runner

**Description:**  
The goal of this project is to build a command-line tool that can execute multiple user-defined tasks concurrently — similar in spirit to `make`, `npm run`, or `Taskfile`. The tool will read a configuration file (e.g., YAML or JSON) describing various commands, their dependencies, and execution constraints. It will then run them efficiently using Go’s goroutines and channels.

The project demonstrates how Go’s concurrency model can be leveraged to build a practical developer utility capable of coordinating multiple long-running processes while managing synchronization, logging, and error handling.

---

## Motivation: Why This Project Fits Go

Go was designed for building fast, concurrent systems with minimal complexity. A CLI task runner is a perfect example of a system that benefits from Go’s strengths:

- **Concurrency:** Many tasks can be executed in parallel.
- **Simplicity:** Go’s lightweight syntax and native tooling (like `go run` and `go build`) make it easy to write and distribute CLI applications.
- **Performance:** Go’s compiled binaries execute quickly and are cross-platform.
- **Reliability:** Static typing and strong error handling make the tool robust and maintainable.

This project embodies Go’s philosophy of simple, composable concurrency applied to a real-world problem.

---

## Key Language Features to Utilize

1. **Goroutines:**  
   Each task runs in its own goroutine to achieve true concurrent execution without explicit threading management.

2. **Channels:**  
   Channels will be used for communication between the scheduler, task workers, and logger — ensuring safe and synchronized message passing.

3. **WaitGroups:**  
   To coordinate the completion of multiple concurrent tasks and wait for them before program termination.

4. **Structs and JSON/YAML Decoding:**  
   Task configurations will be represented as structs, populated via Go’s standard `encoding/json` or external YAML package.

5. **Error Handling and Contexts:**  
   The `context` package will be used to handle task cancellations, timeouts, and dependency failures gracefully.

---

## Basic System Architecture

### Overview

config.yaml  
↓  
Config Parser → Task Manager → Worker Pool  
↓  
Goroutines + Channels  
↓  
Logger / CLI Output

### Components

- **Parser:** Reads and validates the configuration file, constructing a dependency graph.
- **Task Manager:** Schedules tasks based on dependencies and manages their lifecycle.
- **Worker Pool:** Executes tasks concurrently using goroutines, respecting a configurable concurrency limit.
- **Logger:** Streams task outputs and errors in real-time with color-coded terminal output.

---

## Scope

### Must-Haves

- Parse a configuration file (JSON or YAML) containing tasks and dependencies.
- Run tasks concurrently using goroutines and channels.
- Support task dependencies (i.e., run `build` before `deploy`).
- Limit the number of concurrent workers via CLI flag (e.g. `--max=4`).
- Display real-time logs and exit status for each task.

### Nice-to-Haves

- Colorized terminal output and progress indicators.
- Support for environment variables in task definitions.
- Persistent log files.
- Cross-platform builds using `go build`.

### Out of Scope

- Complex dependency resolution or DAG visualization.
- Web-based or GUI interfaces.
- Integration with remote task runners or Docker.

---

## Potential Challenges and Mitigation Strategies

| Challenge                    | Description                                                  | Strategy                                                                 |
| ---------------------------- | ------------------------------------------------------------ | ------------------------------------------------------------------------ |
| **Dependency Management**    | Tasks may depend on one another, requiring careful ordering. | Implement a topological sort or dependency check before execution.       |
| **Error Propagation**        | Failures in one goroutine could affect others.               | Use channels and `context.Context` to cancel dependent tasks on failure. |
| **Output Synchronization**   | Multiple goroutines may write to stdout simultaneously.      | Use a centralized logging channel and a mutex-protected writer.          |
| **Configuration Validation** | Invalid config could cause runtime errors.                   | Implement schema validation and clear error reporting.                   |
| **Scalability**              | Too many concurrent tasks could overwhelm system resources.  | Allow a user-defined concurrency limit and worker pool management.       |

---

## Summary

The **Concurrent CLI Task Runner** project combines Go’s hallmark features — simplicity, concurrency, and performance — into a practical developer tool. It challenges the developer to coordinate multiple concurrent processes, handle communication safely, and produce clean, reliable command-line software.

This project is well-scoped for individual work and offers both a learning opportunity and a tangible, useful outcome.

---
