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
		days := []func() (int, int){
			internal.Day1,
			internal.Day2,
			internal.Day3,
			internal.Day4,
			internal.Day5,
			internal.Day6,
			internal.Day7,
			internal.Day8,
		}
		totalStartTime := time.Now()
		for i, day := range days {
			startTime := time.Now()
			part1, part2 := day()
			endTime := time.Now()
			elapsed := endTime.Sub(startTime)
			fmt.Printf("Day%v took %v\n", i+1, elapsed)
			fmt.Println("  Part 1:", part1)
			fmt.Println("  Part 2:", part2)
			fmt.Println()
		}
		totalEndTime := time.Now()
		totalElapsed := totalEndTime.Sub(totalStartTime)
		fmt.Println("All days took", totalElapsed)
	},
}

func init() {
	rootCmd.AddCommand(solveCmd)
	// TODO add option to only run a specified set of days (single, comma-separated, range, etc.)
}
