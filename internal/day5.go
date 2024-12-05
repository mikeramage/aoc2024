package internal

import (
	"slices"
	"strconv"
	"strings"
)

func Day5() (int, int) {

	lines := Lines("./input/day5.txt")
	precedenceMap := make(map[int]map[int]bool)
	var updates [][]int
	var part1, part2 int
	parsingRules := true
	for _, line := range lines {
		if parsingRules {
			if strings.TrimSpace(line) == "" {
				//Blank line - switch to parsing updates
				parsingRules = false
			} else {
				rule := strings.Split(line, "|")
				precedent, _ := strconv.Atoi(rule[0]) //ignore Atoi errors - input guarantees success
				subsequent, _ := strconv.Atoi(rule[1])
				if precedenceMap[subsequent] == nil {
					precedenceMap[subsequent] = make(map[int]bool)
				}
				precedenceMap[subsequent][precedent] = true
			}
		} else {
			//Parsing updates
			var update []int
			for _, pageTxt := range strings.Split(line, ",") {
				page, _ := strconv.Atoi(pageTxt) //ignore errors
				update = append(update, page)
			}
			updates = append(updates, update)
		}
	}

	for _, update := range updates {
		if isValidUpdate(update, precedenceMap) {
			part1 += getMiddleElement(update)
		} else {
			orderedUpdate := orderUpdate(update, precedenceMap)
			part2 += getMiddleElement(orderedUpdate)
		}
	}

	return part1, part2
}

func isValidUpdate(update []int, precedenceMap map[int]map[int]bool) bool {
	for i, page := range update {
		precedents, exists := precedenceMap[page]
		if exists {
			for _, otherPage := range update[i+1:] {
				if precedents[otherPage] {
					//otherPage is a prerequisite of this one - not valid
					return false
				}
			}
		}
	}
	return true
}

func getMiddleElement(intSlice []int) int {
	//Assumes slice is non-zero, odd length - which is guaranteed by this problem
	//If slice was even length it would sort of work, returning the earlier of the two middle elements
	return intSlice[len(intSlice)/2]
}

func orderUpdate(update []int, precedenceMap map[int]map[int]bool) []int {
	var orderedUpdate []int
	remaining := slices.Clone(update) //Good practice not to modify input.
	for len(remaining) > 0 {          // keep going until no elements left in update
		found := false
		for i, page := range remaining {
			if canPlacePage(page, remaining, precedenceMap) {
				orderedUpdate = append(orderedUpdate, page)
				remaining = slices.Delete(remaining, i, i+1)
				found = true
				break
			}
		}

		if !found {
			panic("No way to reorder updates")
		}
	}
	return orderedUpdate
}

func canPlacePage(page int, remaining []int, precedenceMap map[int]map[int]bool) bool {
	precedents, exists := precedenceMap[page]
	if !exists {
		//No requirements on this page - can place anywhere!
		return true
	}
	for _, otherPage := range remaining {
		if precedents[otherPage] {
			//Page can't be placed yet. A subsequent page is a prerequisite
			return false
		}
	}
	return true
}
