package internal

import (
	"bufio"
	"log"
	"os"
)

type Queue[T any] interface {
	Append(v T)
	PopFront() T
	Len() int
}

// Queue implementation
type QueueImpl[T any] struct {
	elements []T
}

func (q *QueueImpl[T]) Append(v T) {
	q.elements = append(q.elements, v)
}

func (q *QueueImpl[T]) Len() int {
	return len(q.elements)
}

func (q *QueueImpl[T]) PopFront() T {
	if len(q.elements) == 0 {
		panic("Cannot pop from empty queue!")
	}
	var el T
	el, q.elements = q.elements[0], q.elements[1:]

	return el
}

func newQ[T any]() Queue[T] {
	return &QueueImpl[T]{elements: make([]T, 0)}
}

type Position struct {
	row, col int
}

func comparePositions(a, b Position) int {
	if a.row > b.row {
		return 1
	}

	if a.row < b.row {
		return -1
	}

	if a.col > b.col {
		return 1
	}

	if a.col < b.col {
		return -1
	}

	return 0
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

func withinBoundsPos(p Position, rows, cols int) bool {
	return p.row < rows && p.row >= 0 && p.col < cols && p.col >= 0
}
