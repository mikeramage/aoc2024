package utils

import (
	"bufio"
	"log"
	"os"
)

func Lines(fileName string) []string {

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("Could not open file for reading:", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println("Warning: failed to close file:", err)
		}
	}()

	scanner := bufio.NewScanner(f)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
