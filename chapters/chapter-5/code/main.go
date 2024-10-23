package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	totalBytes, _ := fileLen("go.mod")
	fmt.Println(totalBytes)

	s := make([]int, 0, 5)
	s = append(s, 1, 2, 3)
	s2 := s[:]
	s2 = append(s2, 88)
	// modSlice(s)
	fmt.Println(s)
	fmt.Println(s2)

	m := map[int]int{}
	m[11] = 1
	m[22] = 2
	modMap(m)
	fmt.Println(m)
}

// Returns number of bytes in file
func fileLen(filename string) (int, error) {
	// Handle open, error, and close
	f, err := os.Open(filename)
	if err != nil {
		return -1, errors.New("error opening file")
	}
	defer f.Close()

	// Get number of bytes
	data := make([]byte, 2048)
	var totalBytes int

	for {
		count, err := f.Read(data)
		os.Stdout.Write(data[:count])
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
		totalBytes += count
	}

	// for i := 0; i < totalBytes; i++ {
	// 	fmt.Println(i, string(data[i]))
	// }

	return totalBytes, nil
}

func modMap(m map[int]int) {
	m[22] = 23
	m[99] = 9
	// clear(m)
	delete(m, 22)
}
func modSlice(s []int) {
	s[0] = 99
	s = append(s, 45)
}
