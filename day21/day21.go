package day21

import (
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/mikeramage/aoc2024/position"
	"github.com/mikeramage/aoc2024/queue"
	"github.com/mikeramage/aoc2024/utils"
)

var (
	up    byte = '^'
	down  byte = 'v'
	left  byte = '<'
	right byte = '>'
	A     byte = 'A'
	bad   byte = 'X'
)

type BytePair struct {
	one, two byte
}

func Day21() (int, int) {
	var part1, part2 int

	lines := utils.Lines("./input/day21.txt")
	var codes []string
	var numbers []int
	for _, line := range lines {
		codes = append(codes, line)
		n, _ := strconv.Atoi(line[:3])
		numbers = append(numbers, n)
	}

	//Strategy:
	// - Start with numeric pad and generate moves required to enter code. Each move minimizes manhattan distance to next digit and presses
	//   the next digit on arrival. This will generate the series of moves for the next robot in the chain. There will be more than one
	//   solution generated as there may be 2 equivalent paths to the goal.
	// - Do the same for the next robot in the chain for each moveset generated in the previous stage, minimising the manhattan distance to
	//   the appropriate direction key. Again, preserve all solutions that have the same path cost (i.e. number of characters in the solution)
	// - And repeat.
	numberPad, numberPadMap := createNumberPad()
	directionPad, directionPadMap := createDirectionPad()

	//Precompute the possible paths on the direction pad for each combination of origin button A, ^, v, >, <
	//and destination button A, ^, v, >, <. Each set of paths is a map[BytePair]string where the first byte
	//is the origin, the second is the destination, and the string is the sequence of moves. There are multiple
	//maps because, for instance, A -> < can be given by <v< or by v<<. A-priori we don't know which combination is
	//best but we assume that down the line one will prove superior.
	directionPaths := calculateDirectionPaths(directionPad, directionPadMap)

	for c, code := range codes {
		candidateSequences := findPathsForButtonSequence(numberPad, numberPadMap, []byte(code), nil)
		var p1Candidates, p2Candidates []int
		for _, paths := range directionPaths {

			expansionMap := make(map[string]string)

			for _, seq := range candidateSequences {
				countMap := gather(seq)
				for i := 0; i < 25; i++ {
					countMap = expandAndGather(countMap, expansionMap, paths)
					if i == 1 {
						p1Candidates = append(p1Candidates, complexity(numbers[c], countMap))
					}
				}
				p2Candidates = append(p2Candidates, complexity(numbers[c], countMap))
			}
		}
		part1 += slices.Min(p1Candidates)
		part2 += slices.Min(p2Candidates)
	}

	return part1, part2
}

func calculateDirectionPaths(directionPad [][]byte, directionPadMap map[byte]position.Position) []map[BytePair]string {
	buttons := []byte{A, up, down, left, right}

	allPaths := make(map[BytePair][]string)
	indexMap := make(map[BytePair]int)
	for _, orig := range buttons {
		for _, dest := range buttons {
			key := BytePair{orig, dest}
			if orig == dest {
				allPaths[key] = []string{"A"}
				continue
			}
			origNode := &Node{pos: directionPadMap[orig], val: orig, parent: nil}
			allPaths[key] = findPathsForButtonSequence(directionPad, directionPadMap, []byte{dest}, origNode)
			if len(allPaths[key]) > 1 {
				if len(allPaths[key]) != 2 {
					panic("Unexpected number of possibilities!")
				}
				indexMap[key] = 0
			}
		}
	}

	indexKeys := slices.Collect(maps.Keys(indexMap))
	keyPermutations := 1 << len(indexKeys)

	var directionPaths []map[BytePair]string
	for i := 0; i < keyPermutations; i++ {
		assignIndexValues(indexMap, indexKeys, i)
		directionPath := make(map[BytePair]string)
		for key, paths := range allPaths {
			if len(paths) > 1 {
				directionPath[key] = paths[indexMap[key]]
				continue
			}
			directionPath[key] = paths[0]
		}
		directionPaths = append(directionPaths, directionPath)
	}

	return directionPaths
}

func assignIndexValues(indexMap map[BytePair]int, indexKeys []BytePair, value int) {
	//Convert value into an array of integers that are 0 or 1
	indices := make([]int, len(indexKeys))
	for i := 0; value > 0; i++ {
		indices[i] = value & 1
		value >>= 1
	}

	for i, indexKey := range indexKeys {
		indexMap[indexKey] = indices[i]
	}
}

func gather(sequence string) map[string]int {
	countMap := make(map[string]int)
	subsequences := strings.Split(sequence, "A")
	// By construction the last character of every sequence is always an A, and calling split using A as separator will always
	// include an instance of "" as the last element (i.e. the empty string after the last A). So ignore the
	// last character
	for _, subsequence := range subsequences[:len(subsequences)-1] {
		countMap[subsequence]++
	}
	return countMap
}

func expandAndGather(countMap map[string]int, expansionMap map[string]string, directionPaths map[BytePair]string) map[string]int {
	newCountMap := make(map[string]int)
	for subSequence, count := range countMap {
		expansions := expand(subSequence, expansionMap, directionPaths)
		subCounts := gather(expansions)
		for s, c := range subCounts {
			newCountMap[s] += count * c
		}
	}
	return newCountMap
}

func expand(subSequence string, expansionMap map[string]string, directionPaths map[BytePair]string) string {
	//Takes a subsequence like "v<<" converts to "v<<A" and expands it to e.g. "<vA<AA>>^A"

	//Try the cache first
	expansion, exists := expansionMap[subSequence]
	if exists {
		return expansion
	}

	var result string

	subSequence += string(A)

	orig := A
	for _, direction := range subSequence {
		key := BytePair{orig, byte(direction)}
		result += directionPaths[key]
		orig = byte(direction)
	}

	expansion = result
	expansionMap[subSequence] = expansion
	return expansion
}

func complexity(code int, countMap map[string]int) int {
	return code * seqLen(countMap)
}

func seqLen(countMap map[string]int) int {
	totalLen := 0
	for s, c := range countMap {
		totalLen += c * (len(s) + 1)
	}
	return totalLen
}

type Node struct {
	pos     position.Position
	val     byte
	presses int
	parent  *Node
}

func findPathsForButtonSequence(buttonPad [][]byte, buttonPadMap map[byte]position.Position, buttonSequence []byte, startingNode *Node) []string {
	var sequences []string
	if startingNode == nil {
		startingNode = &Node{pos: buttonPadMap[A], val: A, parent: nil}
	}
	var buttonSequenceNodes []*Node
	buttonSequenceNodes = append(buttonSequenceNodes, startingNode)

	for _, button := range buttonSequence {
		var newButtonSequenceNodes []*Node
		for _, node := range buttonSequenceNodes {
			newButtonSequenceNodes = append(newButtonSequenceNodes, findPathsToTargetButton(button, node, buttonPad, buttonPadMap)...)
		}
		buttonSequenceNodes = newButtonSequenceNodes
	}

	for _, node := range buttonSequenceNodes {
		var sequence []byte
		for node.parent != nil {
			for i := 0; i < node.presses; i++ {
				sequence = append(sequence, A)
			}
			dr := node.pos.Row - node.parent.pos.Row
			dc := node.pos.Col - node.parent.pos.Col
			var direction byte
			if dr == 1 && dc == 0 {
				direction = down
			} else if dr == -1 && dc == 0 {
				direction = up
			} else if dr == 0 && dc == 1 {
				direction = right
			} else if dr == 0 && dc == -1 {
				direction = left
			}
			sequence = append(sequence, direction)
			node = node.parent
		}
		slices.Reverse(sequence)
		sequences = append(sequences, string(sequence))
	}

	return sequences
}

func findPathsToTargetButton(targetButton byte, startingNode *Node, buttonPad [][]byte, buttonMap map[byte]position.Position) []*Node {
	targetPosition := buttonMap[targetButton]
	badPosition := buttonMap[bad]
	var subsequenceNodes []*Node

	frontier := queue.NewQueue[*Node]()
	frontier.Append(startingNode)

	for frontier.Len() > 0 {
		node := frontier.PopFront()

		if node.val == targetButton {
			node.presses++
			subsequenceNodes = append(subsequenceNodes, node)
			continue
		}

		for _, dPos := range position.DirectionsPos {
			nextPosition := position.Add(node.pos, dPos)
			if position.WithinBoundsPos(nextPosition, len(buttonPad), len(buttonPad[0])) && nextPosition != badPosition && mD(nextPosition, targetPosition) < mD(node.pos, targetPosition) {
				frontier.Append(&Node{pos: nextPosition, val: buttonPad[nextPosition.Row][nextPosition.Col], parent: node})
			}
		}
	}

	return subsequenceNodes
}

func createNumberPad() ([][]byte, map[byte]position.Position) {
	var numberPad [][]byte
	numberPadMap := make(map[byte]position.Position)
	numberPad = append(numberPad, []byte{'7', '8', '9'})
	numberPad = append(numberPad, []byte{'4', '5', '6'})
	numberPad = append(numberPad, []byte{'1', '2', '3'})
	numberPad = append(numberPad, []byte{bad, '0', 'A'})
	for r, row := range numberPad {
		for c, button := range row {
			numberPadMap[button] = position.NewPosition(r, c)
		}
	}

	return numberPad, numberPadMap
}

func createDirectionPad() ([][]byte, map[byte]position.Position) {
	var directionPad [][]byte
	directionPadMap := make(map[byte]position.Position)

	directionPad = append(directionPad, []byte{bad, up, A})
	directionPad = append(directionPad, []byte{left, down, right})
	for r, row := range directionPad {
		for c, b := range row {
			directionPadMap[b] = position.NewPosition(r, c)
		}
	}

	return directionPad, directionPadMap
}

func mD(pos1, pos2 position.Position) int {
	distance := utils.Abs(pos1.Row-pos2.Row) + utils.Abs(pos1.Col-pos2.Col)
	return distance
}
