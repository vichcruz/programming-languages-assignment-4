package main

import (
	"fmt"
	"time"

	"parent/concurrentcrawler"
	"parent/synchronouscrawler"
)



func main() {
	websiteURLs := []string{
		"https://crawler-test.com/content/custom_text",
		"https://crawler-test.com/content/no_h1",
		"https://crawler-test.com/content/error_page",
		"https://www.wikipedia.org/",
		"https://www.bbc.com/",
		"https://www.nytimes.com/",
		"https://edition.cnn.com/",
		"https://www.reddit.com/",
		"https://www.stackoverflow.com/",
		"https://news.ycombinator.com/",
		"https://go.dev/",
		"https://developer.mozilla.org/",
		"https://www.w3.org/",
		"https://httpstat.us/200",
		"https://httpstat.us/404",
		"https://httpstat.us/500",
		"https://example.org/",
		"https://www.ecma-international.org/",
		"https://example.com",
		"https://httpbin.org/html",
		"https://httpbingo.org/html",
		"https://reqres.in/",
		"https://www.iana.org/domains/reserved",
	}
	fmt.Println("Number of websites to crawl:", len(websiteURLs))

	// Buffered channel to hold results from concurrent crawler
	// concResult := make(chan string, len(websiteURLs))

	// Unbuffered channel to hold results from concurrent crawler
	// Only one send can proceed when a receive is ready
	concResult := make(chan string)

	// Concurrent part
	// How it works
	/* 
		1. Main goroutine starts the concurrent crawler as a goroutine
		2. Concurrent crawler launches 23 goroutines (one for each URL) - all start fetching simultaneously
		3. Main goroutine enters the for loop and waits at title := <-concResult
		4. First goroutine to finish sends its title to the channel
		5. Main goroutine immediately receives it, prints it, and loops back to wait for the next one
		6. Second goroutine to finish sends its title → main receives and prints
		7. Third goroutine to finish sends → main receives and prints
		... and so on until all 23 titles are received
	*/
	
	startConc := time.Now()
	go concurrentcrawler.ConcCrawl(websiteURLs, concResult)
	// Collect results from the channel
	for range websiteURLs {
		title := <-concResult
		// This leads to deadlock
		/*
		if <-result == "" {
			fmt.Println("Concurrently fetched title: No title found")
			} else {
				fmt.Println("Concurrently fetched title:", <-result)
				}		
				*/
				// The result is already consumed in the if statement
				// If it's not true another receive operation is attempted leading to deadlock once all sends are done
				
				// Correct way
				if title == "" {
					fmt.Println("Concurrently fetched title: No title found")
					} else {
						fmt.Println("Concurrently fetched title:", title)
					}
	}
	
				
	durationConc := time.Since(startConc)
	fmt.Println("Concurrent crawling duration:", durationConc)

	fmt.Println("---------------------------------------------------")

	// Synchronous part
	startSync := time.Now()
	syncTitles := synchronouscrawler.SyncCrawl(websiteURLs)
	
	for i := range websiteURLs {
		title := syncTitles[i]
		if title == "" {
			fmt.Println("Synchronously fetched title: No title found")
			} else {
				fmt.Println("Synchronously fetched title:", title)
			}
		}
		
	durationSync := time.Since(startSync)
	fmt.Println("Synchronous crawling duration:", durationSync)

}