/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/mikeramage/aoc2024/day1"
	"github.com/mikeramage/aoc2024/day10"
	"github.com/mikeramage/aoc2024/day11"
	"github.com/mikeramage/aoc2024/day12"
	"github.com/mikeramage/aoc2024/day13"
	"github.com/mikeramage/aoc2024/day14"
	"github.com/mikeramage/aoc2024/day15"
	"github.com/mikeramage/aoc2024/day16"
	"github.com/mikeramage/aoc2024/day17"
	"github.com/mikeramage/aoc2024/day18"
	"github.com/mikeramage/aoc2024/day19"
	"github.com/mikeramage/aoc2024/day2"
	"github.com/mikeramage/aoc2024/day20"
	"github.com/mikeramage/aoc2024/day21"
	"github.com/mikeramage/aoc2024/day22"
	"github.com/mikeramage/aoc2024/day23"
	"github.com/mikeramage/aoc2024/day24"
	"github.com/mikeramage/aoc2024/day25"
	"github.com/mikeramage/aoc2024/day3"
	"github.com/mikeramage/aoc2024/day4"
	"github.com/mikeramage/aoc2024/day5"
	"github.com/mikeramage/aoc2024/day6"
	"github.com/mikeramage/aoc2024/day7"
	"github.com/mikeramage/aoc2024/day8"
	"github.com/mikeramage/aoc2024/day9"
	"github.com/spf13/cobra"
)

// solveCmd represents the solve command
var solveCmd = &cobra.Command{
	Use:   "solve",
	Short: "Solve Advent of Code",
	Long: `Solve Advent of Code for the specified day(s) or range of days, outputs 
the associated solutions and visualizations`,
	Run: func(cmd *cobra.Command, args []string) {
		solutions := []func() (int, int){
			day1.Day1,
			day2.Day2,
			day3.Day3,
			day4.Day4,
			day5.Day5,
			day6.Day6,
			day7.Day7,
			day8.Day8,
			day9.Day9,
			day10.Day10,
			day11.Day11,
			day12.Day12,
			day13.Day13,
			day14.Day14,
			day15.Day15,
			day16.Day16,
			day17.Day17,
			day18.Day18,
			day19.Day19,
			day20.Day20,
			day21.Day21,
			day22.Day22,
			day23.Day23,
			day24.Day24,
			day25.Day25,
		}

		if day == -1 {
			totalStartTime := time.Now()
			for i, solution := range solutions {
				doDay(i+1, solution)
			}
			totalEndTime := time.Now()
			totalElapsed := totalEndTime.Sub(totalStartTime)
			fmt.Println("All days took", totalElapsed)
		} else {
			doDay(day, solutions[day-1])
		}
	},
}

func doDay(day int, solution func() (int, int)) {
	startTime := time.Now()
	part1, part2 := solution()
	endTime := time.Now()
	elapsed := endTime.Sub(startTime)
	fmt.Printf("Day%v took %v\n", day, elapsed)
	fmt.Println("  Part 1:", part1)
	fmt.Println("  Part 2:", part2)
	fmt.Println()
}

func init() {
	rootCmd.AddCommand(solveCmd)
	solveCmd.Flags().IntVarP(&day, "day", "d", -1, "Solve the specified day only")
}
