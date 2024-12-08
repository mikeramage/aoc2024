package internal

import (
	"maps"
)

func Day8() (int, int) {
	lines := Lines("./input/day8.txt")
	antennaMap := buildAntennaMap(lines)
	rows := len(lines)
	cols := len(lines[0])

	antinodeMap1 := make(map[Position]bool)
	antinodeMap2 := make(map[Position]bool)
	var part1, part2 int
	for antenna := range maps.Keys(antennaMap) {
		count1, count2 := collectAntinodes(antennaMap[antenna], antinodeMap1, antinodeMap2, rows, cols)
		part1 += count1
		part2 += count2
	}

	return part1, part2
}

func buildAntennaMap(lines []string) map[byte][]Position {
	antennaMap := make(map[byte][]Position)
	for row, line := range lines {
		for col, antenna := range line {
			if antenna != '.' {
				antennaMap[byte(antenna)] = append(antennaMap[byte(antenna)], Position{row, col})
			}
		}
	}
	return antennaMap
}

func collectAntinodes(antennaPositions []Position, antinodeMap1 map[Position]bool, antinodeMap2 map[Position]bool, rows, cols int) (int, int) {
	count1, count2 := 0, 0
	for i, posA := range antennaPositions {
		for j := i + 1; j < len(antennaPositions); j++ {
			posB := antennaPositions[j]
			dr, dc := posB.row-posA.row, posB.col-posA.col
			count1, count2 = traverseDirections(posA, posB, dr, dc, antinodeMap1, antinodeMap2, rows, cols, count1, count2)
		}
	}
	return count1, count2
}

func traverseDirections(posA, posB Position, dr, dc int, antinodeMap1, antinodeMap2 map[Position]bool, rows, cols, count1, count2 int) (int, int) {
	for multiplier := 0; ; multiplier++ {
		posPlus := Position{row: posB.row + multiplier*dr, col: posB.col + multiplier*dc}
		posMinus := Position{row: posA.row - multiplier*dr, col: posA.col - multiplier*dc}
		validPlus := withinBounds(posPlus.row, posPlus.col, rows, cols)
		validMinus := withinBounds(posMinus.row, posMinus.col, rows, cols)
		if !validPlus && !validMinus {
			break
		}
		if validPlus {
			count1, count2 = processAntinode(posPlus, multiplier, antinodeMap1, antinodeMap2, count1, count2)
		}
		if validMinus {
			count1, count2 = processAntinode(posMinus, multiplier, antinodeMap1, antinodeMap2, count1, count2)
		}
	}
	return count1, count2
}

func processAntinode(pos Position, multiplier int, antinodeMap1, antinodeMap2 map[Position]bool, count1, count2 int) (int, int) {
	if multiplier == 1 {
		if !antinodeMap1[pos] {
			count1++
		}
		antinodeMap1[pos] = true
	}
	if !antinodeMap2[pos] {
		count2++
	}
	antinodeMap2[pos] = true
	return count1, count2
}
