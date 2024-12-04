package internal

func Day4() (int, int) {
	lines := Lines("./input/day4.txt")

	var chars [][]byte
	for _, line := range lines {
		var row []byte
		for i := 0; i < len(line); i++ {
			row = append(row, line[i])
		}
		chars = append(chars, row)
	}

	rows := len(chars)
	cols := len(chars[0])

	var part1, part2 int

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if chars[r][c] == byte('X') {
				part1 += countXmas(chars, r, c, rows, cols)
			}
			if chars[r][c] == byte('A') && isMasX(chars, r, c, rows, cols) {
				part2++
			}
		}
	}

	return part1, part2
}

func countXmas(chars [][]byte, r, c, rows, cols int) int {
	count := 8 //Assume XMAS in all directions
	directions := [][]int{{0, 1}, {-1, 1}, {-1, 0}, {-1, -1}, {0, -1}, {1, -1}, {1, 0}, {1, 1}}
	sequence := []byte{'M', 'A', 'S'}
	for _, direction := range directions {
		rr := r
		cc := c
		for _, letter := range sequence {
			rr += direction[0]
			cc += direction[1]
			if !withinBounds(rr, cc, rows, cols) || chars[rr][cc] != letter {
				count--
				break //Direction doesn't work
			}
		}
	}

	return count
}

func isMasX(chars [][]byte, r, c, rows, cols int) bool {
	return withinBounds(r-1, c-1, rows, cols) &&
		withinBounds(r+1, c+1, rows, cols) &&
		withinBounds(r-1, c+1, rows, cols) &&
		withinBounds(r+1, c-1, rows, cols) &&
		((chars[r-1][c-1] == 'M') && (chars[r+1][c+1] == 'S') ||
			(chars[r-1][c-1] == 'S') && (chars[r+1][c+1] == 'M')) &&
		((chars[r-1][c+1] == 'M') && (chars[r+1][c-1] == 'S') ||
			(chars[r-1][c+1] == 'S') && (chars[r+1][c-1] == 'M'))
}

func withinBounds(r, c, rows, cols int) bool {
	return r < rows && r >= 0 && c < cols && c >= 0
}
