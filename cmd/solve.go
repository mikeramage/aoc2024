/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/mikeramage/aoc2024/internal"
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
			internal.Day1,
			internal.Day2,
			internal.Day3,
			internal.Day4,
			internal.Day5,
			internal.Day6,
			internal.Day7,
			internal.Day8,
			internal.Day9,
			internal.Day10,
			internal.Day11,
			internal.Day12,
			internal.Day13,
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
