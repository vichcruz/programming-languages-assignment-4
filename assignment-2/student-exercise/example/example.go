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
