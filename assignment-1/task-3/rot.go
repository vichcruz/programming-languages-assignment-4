package main

import (
	"fmt"
)

// Used Co-Pilot for this
func getAsciiString(encodedArr []int) string {
	result := make([]byte, len(encodedArr))
	for i, val := range encodedArr {
		result[i] = byte(val)
	}
	return string(result)
}

// Encoding function
func encode(s1 string, s2 string, offset int) string {

	encodedArr := make([]int, len(s1)+len(s1)-1)
	for i := range encodedArr {
		if i%2 == 0 {
			// Used Co-Pilot for normalisation formula
			encodedArr[i] = 97 + (int(s1[i/2 % len(s1)]) + offset - 97) % 26
		} else {
			encodedArr[i] = 97 + (int(s2[i % len(s2)]) + offset - 97) % 26
		}
	}

	return getAsciiString(encodedArr)
}

// Decoding function
func decode(encodedStr string, offset int) string {
	cleanedUpArr := []int{}
	for i := range encodedStr {
		if i % 2 == 0 {
			// Used Co-Pilot for normalisation formula
			cleanedUpArr = append(cleanedUpArr, 97 + ((int(encodedStr[i]) - 97 - offset + 26) % 26))
		}
	}
	return getAsciiString(cleanedUpArr)
}

func main() {
	s1 := "thishouse"
	s2 := "jjk"
	offset := 5
	
	// Example encode
	encodedStr := encode(s1, s2, offset)
	fmt.Println("Encoded string:", encodedStr)

	// And corresponding decode
	decodedStr := decode(encodedStr, offset)
	fmt.Println("Decoded string:", decodedStr)
}