package day24

import (
	"cmp"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/mikeramage/aoc2024/utils"
)

type Wire struct {
	name                string
	initialValue, value WireValue // default to -1 to say input not arrived yet; store off initial value for easy reset
}

type Gate struct {
	input1, input2, output *Wire
	category               Category
	state                  GateState
}

func (g *Gate) do() {
	if g.state != awaitingInput {
		panic("Gate already done!")
	}
	switch g.category {
	case AND:
		g.output.value = g.input1.value & g.input2.value
	case OR:
		g.output.value = g.input1.value | g.input2.value
	case XOR:
		g.output.value = g.input1.value ^ g.input2.value
	default:
		panic(fmt.Sprintln("Gate with invalid category:", g.category))
	}
	g.state = done
}

func (g *Gate) reset() {
	g.state = awaitingInput
	g.output.value = g.output.initialValue
	g.input1.value = g.input1.initialValue
	g.input2.value = g.input2.initialValue
}

func (g *Gate) String() string {
	return fmt.Sprintf("s: %v, c: %v, i1: %v, %v, i2: %v, %v, o: %v, %v\n", g.state, g.category, g.input1.name, g.input1.value, g.input2.name, g.input2.value, g.output.name, g.output.value)
}

type GateState int

const (
	awaitingInput GateState = 0
	done          GateState = 1
)

type Category int

const (
	AND Category = 0
	OR  Category = 1
	XOR Category = 2
)

type WireValue int

const (
	off       WireValue = 0
	on        WireValue = 1
	undefined WireValue = -1
)

var reGateSpec = regexp.MustCompile(`(\w+) (\w+) (\w+) -> (\w+)`)

func Day24() (int, int) {
	var part1, part2 int

	lines := utils.Lines("./input/day24.txt")
	wires := make(map[string]*Wire)
	var zWires []*Wire
	var xWires []*Wire
	var yWires []*Wire
	var gates []*Gate
	outputToGate := make(map[string]*Gate)
	parsingWires := true
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			parsingWires = false
			continue
		}

		if parsingWires {
			wire := parseWire(line)
			wires[wire.name] = &wire
			continue
		}

		//Parsing gates
		gate := parseGate(line, wires)
		gates = append(gates, &gate)
		outputToGate[gate.output.name] = &gate
	}

	for name, wire := range wires {
		if strings.HasPrefix(name, "z") {
			zWires = append(zWires, wire)
			continue
		}
		if strings.HasPrefix(name, "x") {
			xWires = append(xWires, wire)
			continue
		}
		if strings.HasPrefix(name, "y") {
			yWires = append(yWires, wire)
			continue
		}
	}

	slices.SortFunc(zWires, func(a, b *Wire) int { return cmp.Compare(a.name, b.name) })
	slices.SortFunc(xWires, func(a, b *Wire) int { return cmp.Compare(a.name, b.name) })
	slices.SortFunc(yWires, func(a, b *Wire) int { return cmp.Compare(a.name, b.name) })
	slices.SortFunc(gates, func(a, b *Gate) int { return cmp.Compare(a.output.name, b.output.name) })

	changesOccurred := true
	for changesOccurred { //keep going until no further changes occur (could optimize by just doing till all z wires have defined output, but this should be sufficient)
		changesOccurred = false
		for _, gate := range gates {
			if gate.state == awaitingInput && gate.input1.value != undefined && gate.input2.value != undefined {
				gate.do()
				changesOccurred = true
			}
		}
	}

	part1 = toDecimal(zWires)

	// 222 choose 8 groups of gates ~1e14, i.e. not brute-forceable. Strategy is to see which bits are wrong

	//Figure out which bits in the answer are wrong and work out what gates contribute to those.
	actual := toWireValues(part1)
	xDec := toDecimal(xWires)
	yDec := toDecimal(yWires)
	expected := toWireValues(xDec + yDec)
	var wrongBits []int
	// wrongBits := toWireValues(part1 ^ (xDec + yDec))
	// N.B. expected and actual could have different lengths. If that's the case there's a mismatch in
	// the most significant digit.
	for i, bit := range actual {
		if i == len(expected) && bit == 1 {
			wrongBits = append(wrongBits, i)
		} else if bit != expected[i] {
			wrongBits = append(wrongBits, i)
		}
	}

	//Observations:
	// - All Z wires must be outputs of a XOR operation because that's how gate addition works:
	// you XOR the corresponding X and Y and then XOR the result of that with the output of some gate
	// that calulates the "carry" due to operations up to the current result
	// - There are 4 bad Z wires because they're not XORs - 09, 21, 39 and 45. They must each be swapped with
	// the output in a XOR gate, moreover a XOR which doesn't include any X or Y as the inputs (because the final XOR does
	// not include the XOR of X and Y directly). We can ignore 45 since it's the last. We're now only looking for 3 gates.
	// - The first bad Z wire appears at 09. Therefore we can ignore any gates that contribute to Zs below that because if
	// we didn't Z09 would feed into Z0X, X<9 and that's causally wrong.
	//
	// This reduces the possible space of swaps

	seen := make(map[string]bool) //Track which gates we've already considered - they're the ones we can ignore

	// for _, index := range wrongBits {
	for index := 0; index < 9; index++ {
		if seen[zWires[index].name] {
			continue
		}
		earlyGate := outputToGate[zWires[index].name] //Must always exist - z wires are always outputs to some gate
		toExplore := []*Gate{earlyGate}
		for len(toExplore) > 0 {
			earlyGate = toExplore[len(toExplore)-1]  //...
			toExplore = toExplore[:len(toExplore)-1] //... Pop
			seen[earlyGate.output.name] = true
			// Find the gates for which inputs to this gate are outputs (if they exist) and add them to the toExplore stack if we haven't seen them
			cg1, exists := outputToGate[earlyGate.input1.name]
			if exists && !seen[cg1.output.name] && !slices.Contains(toExplore, cg1) {
				toExplore = append(toExplore, cg1)
			}
			cg2, exists := outputToGate[earlyGate.input2.name]
			if exists && !seen[cg2.output.name] && !slices.Contains(toExplore, cg2) {
				toExplore = append(toExplore, cg2)
			}
		}
	}

	reZ := regexp.MustCompile(`z\d+`)
	reX := regexp.MustCompile(`x\d+`)
	reY := regexp.MustCompile(`y\d+`)
	var candidateGates []*Gate
	for _, gate := range gates {
		if gate.category == XOR && !reZ.MatchString(gate.output.name) && !reX.MatchString(gate.input1.name) && !reY.MatchString(gate.input1.name) && !reX.MatchString(gate.input2.name) && !reY.MatchString(gate.input2.name) {
			candidateGates = append(candidateGates, gate)
		}
	}
	//Exactly 3 candidate gates so try all 3 combinations.

	solutionNames := []string{"z09", "z21", "z39", candidateGates[0].output.name, candidateGates[1].output.name, candidateGates[2].output.name}
	// combos := [][]int{{1, 2, 0}, {2, 0, 1}, {0, 1, 2}, {2, 1, 0}, {1, 0, 2}, {0, 2, 1}}
	// for _, combo := range combos {
	resetAll(gates) // resets all gates and wires to their initial parsed state
	swap(outputToGate["z09"], outputToGate["gwh"])
	swap(outputToGate["z21"], outputToGate["rcb"])
	swap(outputToGate["z39"], outputToGate["jct"])
	changesOccurred = true
	for changesOccurred { //keep going until no further changes occur (could optimize by just doing till all z wires have defined output, but this should be sufficient)
		changesOccurred = false
		for _, gate := range gates {
			if gate.state == awaitingInput && gate.input1.value != undefined && gate.input2.value != undefined {
				gate.do()
				changesOccurred = true
			}
		}
	}

	xDec = toDecimal(xWires)
	yDec = toDecimal(yWires)

	//Find the other one - just try everything.
	for i := 0; i < len(gates); i++ {
		for j := i; j < len(gates); j++ {

			testCases := [][]int{{xDec, yDec}, {345256476810, 1148596753211},
				{2548201458879, 3336500147449}}
			allPassed := true
			swap(gates[i], gates[j])

			for _, test := range testCases {
				resetAll(gates)
				xBits := toWireValues(test[0])
				for i := 0; i < len(xWires); i++ {
					if i < len(xBits) {
						xWires[i].value = WireValue(xBits[i])
					} else {
						xWires[i].value = off
					}
				}

				yBits := toWireValues(test[1])
				for i := 0; i < len(yWires); i++ {
					if i < len(yBits) {
						yWires[i].value = WireValue(yBits[i])
					} else {
						yWires[i].value = off
					}
				}
				changesOccurred = true
				for changesOccurred { //keep going until no further changes occur (could optimize by just doing till all z wires have defined output, but this should be sufficient)
					changesOccurred = false
					for _, gate := range gates {
						if gate.state == awaitingInput && gate.input1.value != undefined && gate.input2.value != undefined {
							gate.do()
							changesOccurred = true
						}
					}
				}
				x := test[0]
				y := test[1]
				z := x + y
				zDec := toDecimal(zWires)
				if zDec != z {
					allPassed = false
				}
			}
			if allPassed {
				solutionNames = append(solutionNames, gates[i].output.name)
				solutionNames = append(solutionNames, gates[j].output.name)
			}

			swap(gates[j], gates[i])

		}
	}

	slices.Sort(solutionNames)
	fmt.Println(solutionNames)

	return part1, part2
}

func swap(gate1, gate2 *Gate) {
	tmp := gate1.output
	gate1.output = gate2.output
	gate2.output = tmp
}

func resetAll(gates []*Gate) {
	for _, gate := range gates {
		gate.reset()
	}
}

func toDecimal(zWires []*Wire) int {
	//assumes wires are in sorted order with no gaps z00, z01, z02 ... z45
	decimal := 0
	for i, wire := range zWires {
		decimal += int(wire.value) << i
	}
	return decimal
}

// Essentially reverses toDecimal. Note that the first element of the array is the least
// significant bit so you read left to right.
func toWireValues(decimal int) []byte {
	var wireValues []byte
	for ; decimal > 0; decimal = decimal >> 1 {
		wireValues = append(wireValues, byte(decimal&1))
	}
	return wireValues
}

func parseWire(wireSpec string) Wire {
	components := strings.Split(wireSpec, ": ")
	name := components[0]
	valueInt, _ := strconv.Atoi(components[1])
	value := parseWireValue(valueInt)
	return Wire{name: name, value: value, initialValue: value}
}

func parseGate(gateSpec string, wires map[string]*Wire) Gate {
	components := reGateSpec.FindStringSubmatch(gateSpec)
	input1Name := components[1]
	categoryName := components[2]
	input2Name := components[3]
	outputName := components[4]
	input1, exists := wires[input1Name]
	if !exists {
		input1 = &Wire{name: input1Name, value: undefined, initialValue: undefined}
		wires[input1Name] = input1
	}
	input2, exists := wires[input2Name]
	if !exists {
		input2 = &Wire{name: input2Name, value: undefined, initialValue: undefined}
		wires[input2Name] = input2
	}
	output, exists := wires[outputName]
	if !exists {
		output = &Wire{name: outputName, value: undefined, initialValue: undefined}
		wires[outputName] = output
	}
	category := parseCategory(categoryName)
	return Gate{input1: input1, input2: input2, output: output, category: category, state: awaitingInput}
}

func parseCategory(categoryName string) Category {
	switch categoryName {
	case "AND":
		return AND
	case "OR":
		return OR
	case "XOR":
		return XOR
	default:
		panic(fmt.Sprintln("Unrecognized category:", categoryName))
	}
}

func parseWireValue(state int) WireValue {
	switch state {
	case 0:
		return off
	case 1:
		return on
	default:
		return undefined
	}
}
