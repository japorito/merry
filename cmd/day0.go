/*
Copyright © 2022 Jacob Saporito
*/
package cmd

import (
	"fmt"
	"time"

	xmas "github.com/japorito/merry/libxmas"
	"github.com/spf13/cobra"
)

// day1Cmd represents the day1 command
var day0Cmd = &cobra.Command{
	Use:   "day0 path/to/input/file",
	Short: "AoC Day 0",
	Long: `Advent of Code Day 0 which provides a
very simple outline that can be used for
future (real) days. Requires a filename
argument to process.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		defer xmas.PrintHolidayMessage(time.Now())

		if input := xmas.ReadFileAsInt64Slice(args[0]); input != nil {
			fmt.Printf("%d input lines read.\n", len(input))

			if Parts.Has(1) {
				fmt.Println("Part 1 running...")

				set1 := xmas.BitSet[int]{}
				set1.On(0)
				set1.On(-64)
				fmt.Println(set1.Members())

				set1.On(-1)
				set1.On(64)
				set1.Off(-64)
				fmt.Println(set1.Members())

				set2 := xmas.BitSet[int]{}
				set3 := xmas.BitSet[int]{}
				set4 := xmas.BitSet[int]{}

				set2.On(128)
				set3.On(-128)
				set4.On(64)

				set5 := set1.Union(set2, set3)
				fmt.Println(set5.Members())

				set1.Subtract(set4)
				fmt.Println(set4.Capacity(), set4.Members(), set1.Members())
			}

			if Parts.Has(2) {
				fmt.Println("Part 2 running...")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(day0Cmd)
}
