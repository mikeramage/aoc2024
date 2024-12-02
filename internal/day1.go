package internal

import (
	"bufio"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func Day1() (int, int) {
	f, err := os.Open("./input/day1.txt")
	if err != nil {
		log.Fatalln("Could not open file for reading:", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println("Warning: failed to close file:", err)
		}
	}()

	scanner := bufio.NewScanner(f)

	var left, right []int

	for scanner.Scan() {
		input := strings.Fields(scanner.Text())
		l, _ := strconv.Atoi(input[0])
		r, _ := strconv.Atoi(input[1])
		left = append(left, l)
		right = append(right, r)
	}

	slices.Sort(left)
	slices.Sort(right)

	part1 := 0
	m := make(map[int]int)
	for i := 0; i < len(left); i++ {
		diff := abs(left[i] - right[i])
		part1 += diff
		m[right[i]]++
	}

	part2 := 0
	for _, l := range left {
		part2 += m[l] * l
	}

	return part1, part2
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
