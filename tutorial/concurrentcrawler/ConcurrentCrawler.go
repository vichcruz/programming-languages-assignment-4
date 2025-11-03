package concurrentcrawler

import (
	"fmt"
	"sync"

	"parent/fetcher"
)

func ConcCrawl (urls []string, result chan string) {
	var wg sync.WaitGroup
	// Start of the concurrent web crawler
	fmt.Println("Concurrent crawling")
	for _, url := range urls {
		// Concurrenty crawler
		wg.Go(func() {
			title := fetcher.Fetch(url)
			result <- title
		})
	}

	wg.Wait()
}