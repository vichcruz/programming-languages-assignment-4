package concurrentcrawler

import (
	"fmt"
	"sync"

	"parent/fetcher"
)

func ConcCrawl (urls []string) {
	var wg sync.WaitGroup
	// Start of the concurrent web crawler
	fmt.Println("Concurrent crawling")
	for index, url := range urls {
		// Concurrenty crawler
		wg.Go(func() {
			fetcher.Fetch(url)
			fmt.Println(index)
		})
	}

	wg.Wait()
}