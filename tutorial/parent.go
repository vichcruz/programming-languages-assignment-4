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
	
	startConc := time.Now()
	concurrentcrawler.ConcCrawl(websiteURLs)
	durationConc := time.Since(startConc)
	fmt.Println("Concurrent crawling duration:", durationConc)

	// Synchronous part
	startSync := time.Now()
	synchronouscrawler.SyncCrawl(websiteURLs)
	durationSync := time.Since(startSync)
	fmt.Println("Synchronous crawling duration:", durationSync)

}