package day17

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/mikeramage/aoc2024/utils"
)

func Day17() (int, int) {
	var part1, part2, A, B, C int
	var program []int

	lines := utils.Lines("./input/day17.txt")
	for _, line := range lines {
		if strings.HasPrefix(line, "Register A") {
			A = parseRegister(line)
		} else if strings.HasPrefix(line, "Register B") {
			B = parseRegister(line)
		} else if strings.HasPrefix(line, "Register C") {
			C = parseRegister(line)
		} else if strings.HasPrefix(line, "Program") {
			program = parseProgram(line)
		}
	}

	registers := []*int{&A, &B, &C}
	output := run(program, registers)
	part1Str := strconv.Itoa(output[0])
	for i := 1; i < len(output); i++ {
		part1Str = strings.Join([]string{part1Str, strconv.Itoa(output[i])}, ",")
	}

	fmt.Println(part1Str)

	//And now onto part 2. Sigh. See the README if you want an inadequate explanation
	programIndex := len(program) - 1
	shiftFactor := 3 * programIndex
	candidates := []int{1 << shiftFactor} //Candidate solutions

	for programIndex >= 0 {
		shiftFactor = 3 * programIndex //*3 to convert from binary to octal
		var newCandidates []int
		for _, candidate := range candidates {
			for attempt := 0; attempt < 8; attempt++ {
				A = candidate
				B, C = 0, 0
				output = run(program, registers)

				if output[programIndex] == program[programIndex] {
					//Match
					newCandidates = append(newCandidates, candidate)
				}

				if attempt == 7 {
					//Finished looking for matches - don't want to increment candidate in this case
					break
				}

				candidate += (1 << shiftFactor)
			}
		}
		candidates = newCandidates
		programIndex--
	}

	part2 = slices.Min(candidates)

	return part1, part2
}

func parseRegister(line string) int {
	ss := strings.Split(line, ":")
	s := strings.TrimSpace(ss[1])
	register, _ := strconv.Atoi(s)
	return register
}

func parseProgram(line string) []int {
	var program []int
	ss := strings.Split(line, ":")
	s := strings.TrimSpace(ss[1])
	dss := strings.Split(s, ",")
	for _, ds := range dss {
		d, _ := strconv.Atoi(ds)
		program = append(program, d)
	}
	return program
}

func run(program []int, registers []*int) []int {
	var output []int
	var iPtr int
	for iPtr < len(program)-1 {
		oldIptr := iPtr
		instruction := program[iPtr]
		operand := program[iPtr+1]
		// for _, r := range registers {
		// 	fmt.Printf("%v, ", *r)
		// }
		// fmt.Println(instruction, operand, len(program), iPtr)

		switch instruction {
		case 0:
			adv(registers, operand)
		case 1:
			bxl(registers, operand)
		case 2:
			bst(registers, operand)
		case 3:
			jnz(registers, operand, &iPtr)
		case 4:
			bxc(registers, operand)
		case 5:
			output = out(registers, operand, output)
		case 6:
			bdv(registers, operand)
		case 7:
			cdv(registers, operand)
		default:
			panic(fmt.Sprintf("Unrecognised instruction: %v", instruction))
		}

		if iPtr == oldIptr { //Probably don't need to distinguish between a jump of 0 and not jumping
			iPtr += 2
		}
	}

	// for _, r := range registers {
	// 	fmt.Printf("%v, ", *r)
	// }
	// fmt.Println(len(program), iPtr)

	return output
}

func toComboOperand(registers []*int, operand int) int {
	switch operand {
	case 0, 1, 2, 3:
		return operand
	case 4, 5, 6:
		return *registers[operand-4]
	default:
		panic(fmt.Sprintf("Unrecognised operand: %v", operand))
	}
}

func adv(registers []*int, operand int) {
	dv(registers, 0, operand)
}

func bxl(registers []*int, operand int) {
	*registers[1] ^= operand
}

func bst(registers []*int, operand int) {
	*registers[1] = toComboOperand(registers, operand) % 8
}

func jnz(registers []*int, operand int, iPtr *int) {
	if *registers[0] == 0 {
		return
	}
	*iPtr = operand
}

func bxc(registers []*int, _ int) {
	*registers[1] ^= *registers[2]
}

func out(registers []*int, operand int, output []int) []int {
	output = append(output, toComboOperand(registers, operand)%8)
	return output
}

func bdv(registers []*int, operand int) {
	dv(registers, 1, operand)
}

func cdv(registers []*int, operand int) {
	dv(registers, 2, operand)
}

func dv(registers []*int, rIx int, operand int) {
	num := *registers[0]
	denom := 1 << toComboOperand(registers, operand)
	*registers[rIx] = num / denom
}
