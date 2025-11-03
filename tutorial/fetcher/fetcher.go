package fetcher

import (
	"io"
	"net/http"
	"strings"
)

func Fetch(url string) string {
	// Simple HTTP GET request to check if the URL is reachable
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}

	
	// Read actual response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	bodyString := string(bodyBytes)
	// fmt.Println("Response Body:", bodyString)

	// Schedule the closing of response body to the end of the crawl function (even error or panic cases)
	// This ensures resources are freed properly
	defer resp.Body.Close()

	// Get title from response body
	titleStart := strings.Index(bodyString, "<title>")
	titleEnd := strings.Index(bodyString, "</title>")
	if titleStart != -1 && titleEnd != -1 {
		title := bodyString[titleStart+len("<title>") : titleEnd]
		return title
	} else {
		return ""
	}
}