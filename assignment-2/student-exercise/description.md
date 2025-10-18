# Concurrent File Downloader in Go

**Author:** Victor Cruz da Silva  
**Email:** victorhugo.cruzdasilva@uzh.ch

## Objective

Create a Go program that reads a list of URLs from a text file and downloads the files concurrently. This exercise will help you practice Go's concurrency features, including goroutines and channels.

## Why Go?

Go (Golang) is particularly well-suited for concurrent programming. It was designed from the ground up with concurrency as a core feature, making it perfect for tasks like downloading multiple files simultaneously. The language provides simple yet powerful primitives (goroutines and channels) that make concurrent programming much easier than in traditional languages.

## Learning Goals

- Use **goroutines** to perform tasks concurrently.
- Use **channels** to communicate between goroutines.
- Handle errors gracefully in a concurrent environment.
- Work with Go's standard library (`net/http` for requests, `os` for file operations).
- Understand synchronization and waiting for concurrent tasks to complete.

- Understand synchronization and waiting for concurrent tasks to complete.

## Introduction to Concurrency in Go

### What are Goroutines?

Goroutines are lightweight threads managed by the Go runtime. They allow you to run functions concurrently without the overhead of traditional operating system threads. Starting a goroutine is as simple as adding the `go` keyword before a function call:

```go
go myFunction()  // Runs myFunction concurrently
```

### What are Channels?

Channels are Go's way of allowing goroutines to communicate with each other safely. Think of a channel as a pipe through which you can send and receive values between goroutines.

```go
// Create a channel
messages := make(chan string)

// Send a value into the channel (from one goroutine)
messages <- "Hello"

// Receive a value from the channel (in another goroutine)
msg := <-messages
```

### Simple Example

Here's a basic example that demonstrates goroutines and channels (you can find this code in the `example/` folder):

```go
package main

import (
	"fmt"
	"time"
)

// Function that sends a message after a delay
func greet(name string, ch chan string) {
    time.Sleep(1 * time.Second)
    ch <- "Hello, " + name
}

func main() {
    messages := make(chan string)

    // Start two goroutines
    go greet("Alice", messages)
    go greet("Bob", messages)

    // Receive messages from the channel
    fmt.Println(<-messages)
    fmt.Println(<-messages)
}
```

In this example:

- Two goroutines run concurrently, each calling the `greet` function
- Each goroutine sends a message through the `messages` channel
- The main function waits to receive both messages before completing
- The program demonstrates how goroutines can work simultaneously and communicate safely through channels

## Instructions

### Task Requirements

1. Create a text file named `urls.txt` containing a list of URLs, one per line. You can use image URLs, text files, or any other downloadable content.

2. Write a Go program (`main.go`) that:

   - Reads the URLs from the `urls.txt` file.
   - Creates a folder named `downloads` to store the downloaded files.
   - Starts a goroutine to download each URL concurrently.
   - Uses channels to track when downloads complete.
   - Waits until all downloads are finished before exiting.
   - Saves each file in the `downloads` folder with a unique name (e.g., `file1.jpg`, `file2.txt`).
   - Handles errors gracefully (e.g., if a URL is invalid or unreachable, print an error message but continue with other downloads).

3. Print progress messages for each download:
   - "Starting download: [URL]"
   - "Download completed: [filename]" or "Download failed: [URL] - [error]"
   - "All downloads completed!"

### Hints

- Use `os.ReadFile()` or `bufio.Scanner` to read URLs from the file.
- Use `http.Get()` to download files from URLs.
- Use `os.Create()` to create files locally.
- Use `io.Copy()` to write the HTTP response body to a file.
- Consider using a `sync.WaitGroup` or counting channel receives to wait for all goroutines.

### Example Code Structure

```go
go downloadFile(url, filename, doneChannel)
```

## Passing Criteria

Your solution will be considered a **pass** if it meets ALL of the following criteria:

1. ✅ **Concurrency**: The program downloads files concurrently using goroutines (not sequentially).
2. ✅ **Channel Communication**: Channels are used to signal completion or communicate between goroutines.
3. ✅ **File I/O**: The program correctly reads URLs from `urls.txt` and saves downloaded files to a `downloads` folder.
4. ✅ **Error Handling**: The program handles errors gracefully (e.g., invalid URLs, network failures) without crashing, and prints meaningful error messages.
5. ✅ **Progress Tracking**: The program prints clear progress messages showing when downloads start, complete, or fail.
6. ✅ **Synchronization**: The program waits for all downloads to complete before exiting (no premature termination).

### Additional Example

You can find a small working example demonstrating goroutines and channels in the `example/` folder. This example shows the basic concepts you'll need for the main task.

## Estimated Time

This task is designed to take approximately **12-16 hours** for a beginner, including:

- Setting up Go and understanding the basics (2-3 hours)
- Learning about goroutines and channels (2-3 hours)
- Implementing the file downloader (6-8 hours)
- Testing and debugging (2-3 hours)

## Optional Bonus Challenges

Want to go further? Try these optional enhancements:

- **Concurrency Limit**: Limit the number of concurrent downloads using buffered channels (e.g., max 3 simultaneous downloads).
- **Progress Bar**: Display a progress bar for each download showing bytes downloaded.
- **Retry Logic**: Automatically retry failed downloads up to 3 times.
- **Download Statistics**: Show total download time and average speed for each file.

## Submission Instructions

### Directory Structure

Your submission should have the following structure:

```
exercise/
├── main.go          # Your main program
├── urls.txt         # Text file with URLs (at least 5 URLs for testing)
├── downloads/       # Folder with downloaded files (created by your program)
└── README.md        # (Optional) Any additional notes or instructions
```

### File Naming Conventions

- Main program: `main.go`
- URL list: `urls.txt`
- Downloaded files: Store in `downloads/` folder with descriptive names

### Submission Format

Create a **ZIP file** named `exercise.zip` containing the `exercise/` folder and all its contents.

**Important**: Make sure your program compiles and runs without errors before submitting!
