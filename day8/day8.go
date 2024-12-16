package day8

import (
	"maps"

	"github.com/mikeramage/aoc2024/position"
	"github.com/mikeramage/aoc2024/utils"
)

func Day8() (int, int) {
	lines := utils.Lines("./input/day8.txt")
	antennaMap := buildAntennaMap(lines)
	rows := len(lines)
	cols := len(lines[0])

	antinodeMap1 := make(map[position.Position]bool)
	antinodeMap2 := make(map[position.Position]bool)
	var part1, part2 int
	for antenna := range maps.Keys(antennaMap) {
		count1, count2 := collectAntinodes(antennaMap[antenna], antinodeMap1, antinodeMap2, rows, cols)
		part1 += count1
		part2 += count2
	}

	return part1, part2
}

func buildAntennaMap(lines []string) map[byte][]position.Position {
	antennaMap := make(map[byte][]position.Position)
	for row, line := range lines {
		for col, antenna := range line {
			if antenna != '.' {
				antennaMap[byte(antenna)] = append(antennaMap[byte(antenna)], position.Position{Row: row, Col: col})
			}
		}
	}
	return antennaMap
}

func collectAntinodes(antennaPositions []position.Position, antinodeMap1 map[position.Position]bool, antinodeMap2 map[position.Position]bool, rows, cols int) (int, int) {
	count1, count2 := 0, 0
	for i, posA := range antennaPositions {
		for j := i + 1; j < len(antennaPositions); j++ {
			posB := antennaPositions[j]
			dr, dc := posB.Row-posA.Row, posB.Col-posA.Col
			count1, count2 = traverseDirections(posA, posB, dr, dc, antinodeMap1, antinodeMap2, rows, cols, count1, count2)
		}
	}
	return count1, count2
}

func traverseDirections(posA, posB position.Position, dr, dc int, antinodeMap1, antinodeMap2 map[position.Position]bool, rows, cols, count1, count2 int) (int, int) {
	for multiplier := 0; ; multiplier++ {
		posPlus := position.Position{Row: posB.Row + multiplier*dr, Col: posB.Col + multiplier*dc}
		posMinus := position.Position{Row: posA.Row - multiplier*dr, Col: posA.Col - multiplier*dc}
		validPlus := position.WithinBounds(posPlus.Row, posPlus.Col, rows, cols)
		validMinus := position.WithinBounds(posMinus.Row, posMinus.Col, rows, cols)
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

func processAntinode(pos position.Position, multiplier int, antinodeMap1, antinodeMap2 map[position.Position]bool, count1, count2 int) (int, int) {
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
