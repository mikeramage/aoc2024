package day7

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mikeramage/aoc2024/utils"
)

type Operator string

const (
	plus   = "+"
	mul    = "*"
	concat = "||"
)

func Day7() (int, int) {

	equations := utils.Lines("./input/day7.txt")

	ops := []Operator{plus, mul, concat}
	lhss, rhss := parseEquations(equations)
	part1 := countValidEquations(lhss, rhss, ops[:2])
	part2 := countValidEquations(lhss, rhss, ops)

	return part1, part2
}

func parseEquations(equations []string) ([]int, [][]int) {
	var lhss []int
	var rhss [][]int
	for _, equation := range equations {
		parts := strings.Split(equation, ":")
		lhs, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		lhss = append(lhss, lhs)
		var operands []int
		for _, operandAsString := range strings.Split(strings.TrimSpace(parts[1]), " ") {
			operand, _ := strconv.Atoi(operandAsString)
			operands = append(operands, operand)
		}
		rhss = append(rhss, operands)
	}
	return lhss, rhss
}

func countValidEquations(lhss []int, rhss [][]int, ops []Operator) int {
	var count int
	for i, rhs := range rhss {
		lhs := lhss[i]
		if len(rhs) == 0 {
			continue //Weird but OK
		}

		if rhs[0] == lhs {
			count += lhs
			continue
		}

		currentResults := []int{pop(&rhs)}
		var nextResults []int
	Outer:
		for len(rhs) > 0 {
			secondOperand := pop(&rhs)
			for _, firstOperand := range currentResults {
				for _, operator := range ops {
					result := doOp(firstOperand, secondOperand, operator)
					if result == lhs {
						count += lhs
						break Outer
					}
					nextResults = append(nextResults, result)
				}
			}
			currentResults = nextResults
			nextResults = make([]int, 0)
		}
	}

	return count
}

func doOp(first, second int, op Operator) int {
	switch op {
	case "+":
		return first + second
	case "*":
		return first * second
	case "||":
		operands := []string{strconv.Itoa(first), strconv.Itoa(second)}
		asText := strings.Join(operands, "")
		answer, _ := strconv.Atoi(asText)
		return answer
	default:
		panic(fmt.Sprintln("Unknown operator!", op))
	}

}

func pop(list *[]int) int {
	if len(*list) == 0 {
		panic("Can't pop from an empty list")
	}
	val := (*list)[0]
	*list = (*list)[1:]
	return val
}
