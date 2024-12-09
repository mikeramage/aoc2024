package internal

import (
	"slices"
	"strconv"
)

const space = -1

func Day9() (int, int) {
	diskLayout, blockMapKeys, spaceMapKeys, blockMap, spaceMap := createDataStructures()

	diskLayoutCopy := make([]int, len(diskLayout))
	copy(diskLayoutCopy, diskLayout)

	part1 := calculatePart1(diskLayout)
	part2 := calculatePart2(diskLayoutCopy, blockMapKeys, spaceMapKeys, blockMap, spaceMap)

	//Different approach for part 2

	return part1, part2
}

func calculatePart1(diskLayout []int) int {
	// Part1: Two pointers solution, one at beginning, one at end. Stop when pointers cross.
	part1 := 0
	spaceIndex := 0
	blockIndex := len(diskLayout) - 1
	for {
		for diskLayout[spaceIndex] != space {
			spaceIndex++
		}
		for diskLayout[blockIndex] == space {
			blockIndex--
		}

		if spaceIndex > blockIndex {
			break
		}

		diskLayout[spaceIndex] = diskLayout[blockIndex]
		diskLayout[blockIndex] = space
	}

	for i, blockId := range diskLayout {
		if blockId == space {
			break
		}
		part1 += i * blockId
	}

	return part1
}

func calculatePart2(diskLayout, blockMapKeys, spaceMapKeys []int, blockMap, spaceMap map[int]int) int {
	for _, blockMapKey := range slices.Backward(blockMapKeys) {
		blockSize := blockMap[blockMapKey]
		//Is there room?
		for i, spaceMapKey := range spaceMapKeys { //The keyset won't be kept up to date but that's OK as the underlying map will
			spaceSize := spaceMap[spaceMapKey]
			if blockSize <= spaceSize && blockMapKey > spaceMapKey {
				//Fits
				blockId := diskLayout[blockMapKey]
				for j := 0; j < blockSize; j++ {
					diskLayout[spaceMapKey+j] = blockId
					diskLayout[blockMapKey+j] = space
				}

				spaceRemaining := spaceSize - blockSize
				delete(spaceMap, spaceMapKey)
				if spaceRemaining == 0 {
					spaceMapKeys = slices.Delete(spaceMapKeys, i, i+1)
				} else {
					spaceMap[spaceMapKey+blockSize] = spaceRemaining
					spaceMapKeys[i] = spaceMapKey + blockSize
				}
				break
			}
		}
	}

	part2 := 0
	for i, blockId := range diskLayout {
		if blockId == space {
			continue
		}
		part2 += i * blockId
	}

	return part2
}

func createDataStructures() ([]int, []int, []int, map[int]int, map[int]int) {
	input := Lines("./input/day9.txt")
	var diskMap []int
	for _, char := range input[0] {
		digit, _ := strconv.Atoi(string(char))
		diskMap = append(diskMap, digit)
	}

	blockMap := make(map[int]int) //Maps index in diskLayout to the size of the contiguous block starting at that index
	spaceMap := make(map[int]int) //Similar but map to contiguous sections of space.
	var spaceMapKeys []int        //Want to iterate keys over insertion order, so maintain separate arrays tracking that
	var blockMapKeys []int
	// Obtain disk layout from the mapping
	var diskLayout []int
	for i, size := range diskMap {
		if i%2 == 0 {
			//File
			if size > 0 {
				blockMap[len(diskLayout)] = size
				blockMapKeys = append(blockMapKeys, len(diskLayout))
			}
			for j := 0; j < size; j++ {
				diskLayout = append(diskLayout, i/2)
			}

		} else {
			//Space
			if size > 0 {
				spaceMap[len(diskLayout)] = size
				spaceMapKeys = append(spaceMapKeys, len(diskLayout))
			}
			for j := 0; j < size; j++ {
				diskLayout = append(diskLayout, space)
			}

		}
	}
	return diskLayout, blockMapKeys, spaceMapKeys, blockMap, spaceMap
}
