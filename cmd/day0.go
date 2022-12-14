/*
Copyright © 2022 Jacob Saporito
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/japorito/merry/libxmas/toybag"
	"github.com/japorito/merry/libxmas/xmas"
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
	Run: func(cmd *cobra.Command, args []string) {
		if input := toybag.ReadToLines(args...); input != nil {
			fmt.Printf("%d input lines read.\n", len(input))

			defer xmas.PrintHolidayMessage(time.Now())

			if Parts.Has(1) {
				fmt.Println("Part 1 running...")
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
