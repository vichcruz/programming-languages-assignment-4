package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Please enter four positive integers (m n x y) or type 'quit': ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Quit command
		if strings.EqualFold(input, "quit") {
			fmt.Println("Goodbye!")
			break
		}

		// Split and validate count
		parts := strings.Fields(input)
		if len(parts) != 4 {
			fmt.Println("Please enter exactly 4 values.")
			continue
		}

		// Try to parse integers
		nums := make([]int, 4)
		valid := true
		for i, p := range parts {
			n, err := strconv.Atoi(p)
			if err != nil || n <= 0 {
				fmt.Printf("'%s' is not a positive integer.\n", p)
				valid = false
				break
			}
			nums[i] = n
		}
		if !valid {
			continue
		}

		m, n, x, y := nums[0], nums[1], nums[2], nums[3]

		// Sanity checks
		if m > n {
			fmt.Println("The first number must be smaller than the second number.")
			continue
		}
		if x == y {
			fmt.Println("The third and fourth number must be different.")
			continue
		}

		// FizzBuzz logic
		for i := m; i < n + 1; i++ {
			switch {
			case i%(x*y) == 0:
				fmt.Print("fizzbuzz ")
			case i%x == 0:
				fmt.Print("fizz ")
			case i%y == 0:
				fmt.Print("buzz ")
			default:
				fmt.Print(i, " ")
			}
		}
		fmt.Println("\n")
	}
}