package day13

import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"

	"github.com/mikeramage/aoc2024/utils"
)

type Button struct {
	X, Y, tokens int
}

type Prize struct {
	X, Y int
}

type ClawMachine struct {
	A, B  Button
	prize Prize
}

var reButton = regexp.MustCompile(`Button (A|B): X\+(\d+), Y\+(\d+)`)
var rePrize = regexp.MustCompile(`Prize: X=(\d+), Y=(\d+)`)

func Day13() (int, int) {
	var part1, part2 int
	var buttonA, buttonB Button
	var prize Prize
	var clawMachines []ClawMachine

	lines := utils.Lines("./input/day13.txt")
	for _, line := range lines {
		if match := reButton.FindStringSubmatch(line); match != nil {
			X, _ := strconv.Atoi(match[2])
			Y, _ := strconv.Atoi(match[3])
			if match[1] == "A" {
				buttonA = Button{X: X, Y: Y, tokens: 3}
			} else if match[1] == "B" {
				buttonB = Button{X: X, Y: Y, tokens: 1}
			} else {
				panic("Button not A or B!")
			}

		} else if match := rePrize.FindStringSubmatch(line); match != nil {
			X, _ := strconv.Atoi(match[1])
			Y, _ := strconv.Atoi(match[2])
			prize = Prize{X: X, Y: Y}
		} else {
			//Assume blank line, but don't bother to check
			clawMachines = append(clawMachines, ClawMachine{buttonA, buttonB, prize})
		}
	}

	//We'll have missed the last one - add it now
	clawMachines = append(clawMachines, ClawMachine{buttonA, buttonB, prize})

	for _, clawMachine := range clawMachines {
		minTokensRequired := 0
		for aPresses := 1; aPresses <= 100; aPresses++ {
			for bPresses := 1; bPresses <= 100; bPresses++ {
				if aPresses*clawMachine.A.X+bPresses*clawMachine.B.X == clawMachine.prize.X &&
					aPresses*clawMachine.A.Y+bPresses*clawMachine.B.Y == clawMachine.prize.Y {
					tokensRequired := aPresses*clawMachine.A.tokens + bPresses*clawMachine.B.tokens
					if minTokensRequired == 0 {
						minTokensRequired = tokensRequired
					} else {
						minTokensRequired = min(minTokensRequired, tokensRequired)
					}
				}
			}
		}
		part1 += minTokensRequired
	}

	c := 10000000000000
	part2Big := big.NewInt(0)
	for _, clawMachine := range clawMachines {
		// There's a unique solution 2 equations, 2 unknowns and I've checked they're not linearly dependent (i.e. denominator Bx*Ay - Ax*By != 0) in all cases - so just use maths.
		// Gotta be cheating, right? Right???
		aPresses := big.NewInt(-1)
		bPresses := big.NewInt(-1)

		bNum := big.NewInt(int64((clawMachine.prize.X+c)*clawMachine.A.Y - (clawMachine.prize.Y+c)*clawMachine.A.X))
		bDenom := big.NewInt(int64(clawMachine.B.X*clawMachine.A.Y - clawMachine.A.X*clawMachine.B.Y))
		bMod := big.NewInt(-1)
		zero := big.NewInt(0)
		if bDenom.Cmp(zero) != 0 {
			bPresses.DivMod(bNum, bDenom, bMod)
		}

		aNum := big.NewInt(int64(clawMachine.prize.X + c))
		aNumB := big.NewInt(int64(clawMachine.B.X))
		aNumB.Mul(aNumB, bPresses)
		aNum.Sub(aNum, aNumB)
		aDenom := big.NewInt(int64(clawMachine.A.X))
		aMod := big.NewInt(-1)
		aPresses.DivMod(aNum, aDenom, aMod)

		three := big.NewInt(3)
		if bMod.Cmp(zero) == 0 && aMod.Cmp(zero) == 0 && bPresses.Cmp(zero) > 0 && aPresses.Cmp(zero) > 0 {
			//Integer solutions found
			part2Big.Add(part2Big, aPresses.Add(aPresses.Mul(aPresses, three), bPresses))
		}
	}

	fmt.Println(part2Big)
	part2 = int(part2Big.Int64())

	return part1, part2
}
