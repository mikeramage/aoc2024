package internal

import (
	"bufio"
	"log"
	"os"
)

type Position struct {
	row, col int
}

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

func withinBounds(r, c, rows, cols int) bool {
	return r < rows && r >= 0 && c < cols && c >= 0
}
