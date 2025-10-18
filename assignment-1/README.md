# Programming Languages Assignment 1

This repository contains solutions for Assignment 1 with 4 tasks implemented in Go.

## Prerequisites

- Go 1.18 or higher
- Download [installer](https://go.dev/dl/) corresponding to your operating system
- Follow the instructions in the installer

## Project Structure

```
assignment-1/
├── task-1/          # Hello World program
├── task-2/          # Interactive FizzBuzz program
├── task-3/          # ROT cipher encoding/decoding
├── task-4/          # Code snippets (quicksort, XML parsing, external commands)
└── README.md        # This file
```

## How to Compile and Run

### Running Tasks

To run any task, navigate to its folder and execute:

```bash
cd task-X  # Replace X with task number (1, 2, 3, or 4)
go run .
```

### Compiling to Binary

To compile a task to a binary executable:

```bash
cd task-X
go build .
```

This creates an executable file that you can run directly.

### Examples

```bash
# Run task 1
cd task-1
go run .

# Compile task 3 to binary
cd task-3
go build
./rot  # Run the compiled binary
```

## Task Descriptions

- **Task 1**: Basic Hello World program
- **Task 2**: Interactive FizzBuzz program
- **Task 3**: ROT cipher implementation with encoding and decoding functions
- **Task 4**: Various Go code snippets including quicksort, XML parsing, and external command execution

## Submission Compliance

✅ All source code files compile and run without errors  
✅ README with compilation and execution instructions provided  
✅ Compiled binaries can be generated using `go build`  
✅ Screenshots of compilation and execution included (see screenshots/ folder)

## Screenshots

Screenshots of each task compiling and running are included in the `screenshots/` folder:

- `hello-world.png`
- `fizzbuzz.png`
- `rot.png`
- `snippets.png`
