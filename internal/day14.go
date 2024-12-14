package internal

import (
	"fmt"
	"maps"
	"math"
	"regexp"
	"strconv"
)

type Robot struct {
	pos, vel Position
}

func (robot *Robot) move(rows, cols int) {
	robot.pos.row += robot.vel.row
	if robot.pos.row >= rows {
		robot.pos.row %= rows
	}
	if robot.pos.row < 0 {
		robot.pos.row = rows - Abs(robot.pos.row)
	}
	robot.pos.col += robot.vel.col
	if robot.pos.col >= cols {
		robot.pos.col %= cols
	}
	if robot.pos.col < 0 {
		robot.pos.col = cols - Abs(robot.pos.col)
	}
}

func (robot Robot) String() string {
	return fmt.Sprintf("p(%v, %v), v(%v, %v)", robot.pos.row, robot.pos.col, robot.vel.row, robot.vel.col)
}

func Day14() (int, int) {
	var part1, part2 int

	robots, counts := parseRobots()
	rows, cols := 103, 101
	// rows, cols := 7, 11

	min_variance := float64(-1)
	tree_iteration := 0
	treeCounts := make(map[Position]int)
	for i := 0; i < 10_000; i++ {
		mean_row, mean_col := float64(0), float64(0)
		for j, robot := range robots {
			counts[robot.pos]-- //No need to check existence; know robot is there
			robot.move(rows, cols)
			counts[robot.pos]++
			mean_row = (float64(j)*mean_row + float64(robot.pos.row)) / float64(j+1)
			mean_col = (float64(j)*mean_col + float64(robot.pos.col)) / float64(j+1)
		}
		variance := float64(0)
		for j, robot := range robots {
			distanceSq := math.Pow((float64(robot.pos.row)-mean_row), 2) + math.Pow((float64(robot.pos.col)-mean_col), 2)
			variance = (float64(j)*variance + float64(distanceSq)) / float64(j+1)
		}

		if min_variance == -1 {
			min_variance = variance
			tree_iteration = i + 1
			maps.Copy(treeCounts, counts)
		} else if variance < min_variance {
			tree_iteration = i + 1
			min_variance = variance
			maps.Copy(treeCounts, counts)
		}

		if i == 99 {
			part1 = safetyFactor(robots, counts, rows, cols)
		}
	}

	visualizeCounts(treeCounts, rows, cols, fmt.Sprintf("Tree: iteration %v", tree_iteration))

	part2 = tree_iteration

	return part1, part2
}

func safetyFactor(robots []*Robot, counts map[Position]int, rows, cols int) int {
	safetyFactor := 0

	safetyFactor = quadrantCount(robots, counts, 0, rows/2, 0, cols/2)            //1st quadrant
	safetyFactor *= quadrantCount(robots, counts, 0, rows/2, cols/2+1, cols)      //2nd quadrant +1 to avoid middle column index
	safetyFactor *= quadrantCount(robots, counts, rows/2+1, rows, 0, cols/2)      //3rd quadrant
	safetyFactor *= quadrantCount(robots, counts, rows/2+1, rows, cols/2+1, cols) //4th quadrant

	return safetyFactor
}

// minRow inclusive, maxRow exclusive, same for cols
func quadrantCount(robots []*Robot, counts map[Position]int, minRow, maxRow, minCol, maxCol int) int {
	count := 0
	counted := make(map[Position]bool)
	for _, robot := range robots {
		if robot.pos.row >= minRow && robot.pos.row < maxRow && robot.pos.col >= minCol && robot.pos.col < maxCol && !counted[robot.pos] {
			count += counts[robot.pos]
			counted[robot.pos] = true
		}
	}
	return count
}

func visualizeCounts(counts map[Position]int, rows, cols int, label string) {
	fmt.Println(label)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			p := Position{r, c}
			count := counts[p]
			if count < 0 {
				panic("Can't have negative count!")
			}
			if count > 0 {
				fmt.Printf("%v", count)
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}

	fmt.Println()
}

var reRobot = regexp.MustCompile(`p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)

func parseRobots() ([]*Robot, map[Position]int) {
	counts := make(map[Position]int)
	var robots []*Robot
	lines := Lines("./input/day14.txt")
	for _, line := range lines {
		if match := reRobot.FindStringSubmatch(line); match != nil {
			r, _ := strconv.Atoi(match[2])
			c, _ := strconv.Atoi(match[1])
			vr, _ := strconv.Atoi(match[4])
			vc, _ := strconv.Atoi(match[3])
			position := Position{r, c}
			velocity := Position{vr, vc}
			robot := Robot{position, velocity}
			robots = append(robots, &robot)
			counts[robot.pos]++
		} else {
			panic("No Robot match!")
		}
	}
	return robots, counts
}
