package synchronouscrawler

import (
	"fmt"

	"parent/fetcher"
)

func SyncCrawl (urls []string)  {
	// Start of the synchronous web crawler
	fmt.Println("Synchronous crawling")
	for _, url := range urls {
		// This will result in main goroutine exiting before crawl completes
		// go crawl(url)

		// Synchronous call to crawl
		fetcher.Fetch(url)
	}
}