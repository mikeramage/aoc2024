package internal

import (
	"maps"
	"slices"
)

func Day12() (int, int) {
	input := Lines("./input/day12.txt")
	allRegions := parseRegions(input)
	part1, part2 := calculatePrice(allRegions)
	return part1, part2
}

func calculatePrice(allRegions map[string][]map[Position]bool) (int, int) {
	var part1, part2 int
	for _, regions := range allRegions {
		for _, region := range regions {
			rowEdgeSegments := make(map[int][]int)
			colEdgeSegments := make(map[int][]int)
			area := len(region)
			var perimeter int
			for _, pos := range slices.SortedFunc(maps.Keys(region), comparePositions) {
				directions := []Position{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
				for _, direction := range directions {
					p := Position{pos.row + direction.row, pos.col + direction.col}
					if !region[p] {
						perimeter++
						if direction.row == 0 && direction.col == 1 {
							//right neighbour not in region - this is column edge segment
							colEdgeSegments[pos.col+1] = append(colEdgeSegments[pos.col+1], pos.row)
						} else if direction.row == 1 && direction.col == 0 {
							//down neighbour not in region - row edge segment
							rowEdgeSegments[pos.row+1] = append(rowEdgeSegments[pos.row+1], pos.col)
						} else if direction.row == 0 && direction.col == -1 {
							//left neighbour not in region - this is column edge segment
							colEdgeSegments[pos.col] = append(colEdgeSegments[pos.col], pos.row)
						} else { //drc.row == -1 && drc.col == 0
							//up neighbour not in region - row edge segment
							rowEdgeSegments[pos.row] = append(rowEdgeSegments[pos.row], pos.col)
						}
					}
				}
			}
			part1 += area * perimeter
			part2 += area * countEdges(rowEdgeSegments, colEdgeSegments)
		}
	}
	return part1, part2
}

func parseRegions(input []string) map[string][]map[Position]bool {
	allRegions := make(map[string][]map[Position]bool)
	rows, cols := len(input), len(input[0])
	var plotGrid [][]string
	for _, line := range input {
		var row []string
		for _, plot := range line {
			row = append(row, string(plot))
		}
		plotGrid = append(plotGrid, row)
	}
	for row, plots := range plotGrid {
		for col, plot := range plots {
			pos := Position{row, col}
			regions, exists := allRegions[plot]
			if !exists { //No regions with this plant type exist yet
				region := make(map[Position]bool)
				region[pos] = true
				allRegions[plot] = []map[Position]bool{region}
			} else {
				contiguousRegions := make(map[int]bool)
				//Find all
				for index, region := range regions {
					directions := []Position{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
					for _, direction := range directions {
						p := Position{pos.row + direction.row, pos.col + direction.col}
						if withinBoundsPos(p, rows, cols) && region[p] {
							// A plot's neighbour is in the same region - add to set of contiguous regions
							contiguousRegions[index] = true
							region[pos] = true //Just do this for all regions - we'll merge later if we find more than one
						}
					}
				}
				if len(contiguousRegions) == 0 {
					//This plot start a new region
					region := make(map[Position]bool)
					region[pos] = true
					allRegions[plot] = append(allRegions[plot], region)
				} else if len(contiguousRegions) > 1 {
					//There are more than one disjoint regions neighbouring this plot. Merge them
					allRegions[plot] = mergeRegions(regions, contiguousRegions)
				}
			}
		}
	}

	return allRegions
}

func countEdges(rowEdgeSegments, colEdgeSegments map[int][]int) int {
	count := doCount(rowEdgeSegments, colEdgeSegments)
	count += doCount(colEdgeSegments, rowEdgeSegments)
	return count
}

// We count edges by checking for the end of an edge, which is determined by an intersection with
// perpendicular edge segments one row or column ahead, with the edge segment originating or
// terminating on the column or row.
func doCount(edgeSegments, perpendicularEdgeSegments map[int][]int) int {
	count := 0
	for dim1, dim2s := range edgeSegments {
		for _, dim2 := range dim2s {
			intersections, exists := perpendicularEdgeSegments[dim2+1]
			if exists && (slices.Contains(intersections, dim1) || (slices.Contains(intersections, dim1-1))) {
				count++
			}
		}
	}

	return count
}

func mergeRegions(regions []map[Position]bool, regionsToMerge map[int]bool) []map[Position]bool {
	var revSortRegionsToMerge []int
	for index := range maps.Keys(regionsToMerge) {
		revSortRegionsToMerge = append(revSortRegionsToMerge, index)
	}
	slices.Sort(revSortRegionsToMerge)
	slices.Reverse(revSortRegionsToMerge)

	mergedRegion := regions[revSortRegionsToMerge[len(revSortRegionsToMerge)-1]] //The lowest indexed region - merge into this
	for _, index := range revSortRegionsToMerge {
		if index != revSortRegionsToMerge[len(revSortRegionsToMerge)-1] {
			region := regions[index]
			for key := range maps.Keys(region) {
				mergedRegion[key] = true
			}
			regions = slices.Concat(regions[:index], regions[index+1:])
		}
	}
	return regions
}
