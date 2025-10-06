package main

import (
	"cmp"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
)

// Co-Pilot was used for the entire quick sort
// Quicksort is a generic implementation of the Quicksort algorithm.
// It uses Go's generics to work with any type that implements the constraints.Ordered interface.
func Quicksort[T cmp.Ordered](arr []T) []T {
	// Base case: arrays with 0 or 1 elements are already sorted
	if len(arr) <= 1 {
		return arr
	}

	// Choose the pivot (we'll use the first element for simplicity)
	pivot := arr[0]

	// Partition the array into three parts: less, equal, and greater
	var less, equal, greater []T
	for _, v := range arr {
		switch {
		case v < pivot:
			less = append(less, v)
		case v == pivot:
			equal = append(equal, v)
		case v > pivot:
			greater = append(greater, v)
		}
	}

	// Recursively sort the less and greater partitions, then concatenate the results
	return append(append(Quicksort(less), equal...), Quicksort(greater)...)
}

func callExternalCmd(cmd string, args ...string) string {
	command := exec.Command(cmd, args...)
	out, err := command.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

// Co-Pilot was used for the xml downloading and parsing
type Book struct {
	Category string   `xml:"category,attr"`
	Cover    string   `xml:"cover,attr,omitempty"`
	Title    string   `xml:"title"`
	Author   []string `xml:"author"`
	Year     int      `xml:"year"`
	Price    float64  `xml:"price"`
}

type Bookstore struct {
	Books []Book `xml:"book"`
}

func calculateTotalPriceAndCount(url string) (int, float64, error) {
	// Download the XML file
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	// Parse the XML
	var bookstore Bookstore
	err = xml.Unmarshal(body, &bookstore)
	if err != nil {
		return 0, 0, err
	}

	// Calculate the total price and count the books
	totalPrice := 0.0
	for _, book := range bookstore.Books {
		totalPrice += book.Price
	}

	return len(bookstore.Books), totalPrice, nil
}

func main() {
	// Example usage of the Quicksort function with integers
	ints := []int{3, 6, 8, 10, 1, 2, 1}
	fmt.Println("Sorted integers:", Quicksort(ints))

	// Example usage of the Quicksort function with floating-point numbers
	floats := []float64{3.1, 6.4, 8.2, 10.0, 1.5, 2.3, 1.1}
	fmt.Println("Sorted floats:", Quicksort(floats))

	// Example usage of the Quicksort function with strings
	strings := []string{"banana", "apple", "cherry", "date", "fig", "grape"}
	fmt.Println("Sorted strings:", Quicksort(strings))

	// Example usage of external command function
	fmt.Println(callExternalCmd("uname", "-a"))

	// Example usage of calculateTotalPriceAndCount function
	url := "https://www.w3schools.com/xml/books.xml"
	count, totalPrice, err := calculateTotalPriceAndCount(url)
	if err != nil {
		log.Fatalf("Error calculating total price and count: %v", err)
	}
	fmt.Printf("Number of books: %d, Total price: %.2f\n", count, totalPrice)
}
