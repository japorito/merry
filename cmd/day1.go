/*
Copyright © 2022 Jacob Saporito
*/
package cmd

import (
	"fmt"
	"time"

	xmas "github.com/japorito/merry/libxmas"
	sleigh "github.com/japorito/merry/libxmas/input"
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
	Long:  `Advent of Code Day 1: Calorie Counting`,
	Run: func(cmd *cobra.Command, args []string) {
		defer xmas.PrintHolidayMessage(time.Now())

		if input := sleigh.ReadToInt64SliceBlocks(args...); input != nil {
			fmt.Printf("%d elves are carrying snacks.\n", len(input))

			calories := calculateElfCalories(input)
			xmas.SortInt64Desc(calories)

			if Parts.Has(1) {
				fmt.Println("Part 1 running...")
				fmt.Printf("Answer 1: The top snack stash has **%d** calories.\n", calories[0])
			}

			if Parts.Has(2) {
				fmt.Println("Part 2 running...")
				fmt.Printf("Answer 2: The top 3 snack stashes have **%d** total calories.\n", xmas.SumInt64(calories[:3]))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(day1Cmd)
}
