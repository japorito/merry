/*
Copyright © 2022 Jacob Saporito
*/
package cmd

import (
	"fmt"

	xmas "github.com/japorito/merry/libxmas"
	"github.com/spf13/cobra"
)

func calculateElfCalories(input [][]int64) []int64 {
	var calories []int64

	for _, elf := range input {
		var sum int64 = 0
		for _, snack := range elf {
			sum = sum + snack
		}

		calories = append(calories, sum)
	}

	return calories
}

// day1Cmd represents the day1 command
var day1Cmd = &cobra.Command{
	Use:   "day1 path/to/input/file",
	Short: "Advent of Code Day 1",
	Long:  `Advent of Code Day 1: Elf Calories`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var runall bool = Part == "*"

		input, err := xmas.ReadFileToInt64SliceBlocks(args[0])
		if err != nil {
			return err
		}
		fmt.Printf("%d input blocks read.\n", len(input))

		calories := calculateElfCalories(input)
		xmas.SortInt64Desc(calories)

		if runall || Part == "1" {
			fmt.Println("Part 1 running...")
			fmt.Printf("Answer 1: %d\n", calories[0])
		}

		if runall || Part == "2" {
			fmt.Println("Part 2 running...")
			fmt.Printf("Answer 2: %d\n", xmas.SumInt64(calories[:3]))
		}

		xmas.PrintHolidayMessage()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(day1Cmd)
}
