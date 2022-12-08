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

func findNUniqueIndex(datastream []rune, unique int64) int64 {
	letterSet := xmas.BitSet[rune]{}
	for i := unique; i < int64(len(datastream)); i++ {
		for _, letter := range datastream[i-unique : i] {
			letterSet.On(letter)
		}

		if letterSet.Size() == unique {
			return i
		}
		letterSet.Clear()
	}

	fmt.Println("Not found!")
	return -1
}

// day6Cmd represents the day6 command
var day6Cmd = &cobra.Command{
	Use:   "day6 path/to/input/file",
	Short: "AoC Day 6",
	Long:  `Advent of Code Day 6: Tuning Trouble`,
	Run: func(cmd *cobra.Command, args []string) {
		if input := sleigh.ReadToRuneSliceLines(args...); input != nil {
			fmt.Printf("%d communications datastreams read.\n", len(input))

			defer xmas.PrintHolidayMessage(time.Now())

			if Parts.Has(1) {
				fmt.Println("Part 1 running...")
				fmt.Printf("First start-of-packet marker found at position %d.\n", findNUniqueIndex(input[0], 4))
			}

			if Parts.Has(2) {
				fmt.Println("Part 2 running...")
				fmt.Printf("First start-of-message marker found at position %d.\n", findNUniqueIndex(input[0], 14))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(day6Cmd)
}
