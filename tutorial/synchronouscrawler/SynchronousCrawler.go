package synchronouscrawler

import (
	"fmt"

	"parent/fetcher"
)

func SyncCrawl (urls []string) []string {
	// Start of the synchronous web crawler
	fmt.Println("Synchronous crawling")
	var results []string
	for _, url := range urls {

		// Synchronous call to crawl
		title := fetcher.Fetch(url)
		results = append(results, title)
	}
	return results
}