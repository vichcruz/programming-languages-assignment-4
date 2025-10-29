package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func Fetch(url string) {
	// Simple HTTP GET request to check if the URL is reachable
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error accessing URL:", err)
		return
	}

	
	// Read actual response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	bodyString := string(bodyBytes)
	// fmt.Println("Response Body:", bodyString)

	// Get title from response body
	titleStart := strings.Index(bodyString, "<title>")
	titleEnd := strings.Index(bodyString, "</title>")
	if titleStart != -1 && titleEnd != -1 {
		title := bodyString[titleStart+len("<title>") : titleEnd]
		fmt.Println("Page Title:", title)
		fmt.Println("Status Code:", resp.StatusCode)
		} else {
			fmt.Println("No title found")
			fmt.Println("Status Code:", resp.StatusCode)
	}

	// Schedule the closing of response body to the end of the crawl function (even error or panic cases)
	// This ensures resources are freed properly
	defer resp.Body.Close()

	// Print the status code
}