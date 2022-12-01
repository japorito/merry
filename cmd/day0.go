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
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		var runall bool = Part == "*"

		input, err := xmas.ReadFileAsInt64Slice(args[0])
		if err != nil {
			return err
		}
		fmt.Printf("%d input lines read.\n", len(input))

		if runall || Part == "1" {
			fmt.Println("Part 1 running...")
		}

		if runall || Part == "2" {
			fmt.Println("Part 2 running...")
		}

		xmas.PrintHolidayMessage(time.Since(start))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(day0Cmd)
}
