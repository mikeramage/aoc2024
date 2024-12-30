package day23

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/mikeramage/aoc2024/utils"
)

func Day23() (int, int) {
	var part1, part2 int
	lines := utils.Lines("./input/day23.txt")
	links := make(map[string][]string)
	for _, line := range lines {
		computers := strings.Split(line, "-")
		links[computers[0]] = append(links[computers[0]], computers[1])
		links[computers[1]] = append(links[computers[1]], computers[0])
	}

	threesomes := make(map[string]bool)          //Unique threesome names, 6-byte string of concatenated, alphabetically ordered computer names
	threesomesIAmIn := make(map[string][]string) //Keyed by computer, list of the threesomes they are in
	seenAndDone := make(map[string]bool)         // When building threesomes, don't consider computers for which all possible threesomes have been constructed
	largestGroups := make(map[string][]string)   // Largest group containing the computer mapped to the members in the group (including this one)
	for first, firstLinks := range links {
		for i := 0; i < len(firstLinks); i++ {
			second := firstLinks[i]
			group := []string{first, second}
			for j := i + 1; j < len(firstLinks); j++ {
				third := firstLinks[j]
				if !seenAndDone[second] && !seenAndDone[third] {
					if slices.Contains(links[second], third) {
						makeThreesome(first, second, third, threesomes, threesomesIAmIn)
					}
					fullyConnected := true
					for _, member := range group {
						if !slices.Contains(links[member], third) {
							fullyConnected = false
						}
					}
					if fullyConnected {
						group = append(group, third)
					}
				}
			}

			if len(group) > len(largestGroups[first]) {
				largestGroups[first] = group
			}
		}
		seenAndDone[first] = true
	}

	for threesome := range threesomes {
		if threesome[0] == 't' || threesome[2] == 't' || threesome[4] == 't' {
			part1++
		}
	}

	largestGroup := slices.MaxFunc(slices.Collect(maps.Values(largestGroups)), func(a, b []string) int {
		return cmp.Compare(len(a), len(b))
	})

	fmt.Println(groupPassword(largestGroup))

	return part1, part2
}

func makeThreesome(first, second, third string, threesomes map[string]bool, threesomesIAmIn map[string][]string) {
	threesome := threesomeName(first, second, third)
	threesomes[threesome] = true
	threesomesIAmIn[first] = append(threesomesIAmIn[first], threesome)
	threesomesIAmIn[second] = append(threesomesIAmIn[second], threesome)
	threesomesIAmIn[second] = append(threesomesIAmIn[second], threesome)
}

func threesomeName(first, second, third string) string {
	names := []string{first, second, third}
	slices.Sort(names)
	return strings.Join(names, "")
}

func groupPassword(group []string) string {
	slices.Sort(group)
	return strings.Join(group, ",")
}
