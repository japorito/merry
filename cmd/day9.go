/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/japorito/merry/libxmas/toybag"
	"github.com/japorito/merry/libxmas/xmas"
	"github.com/spf13/cobra"
)

// day9Cmd represents the day9 command
var day9Cmd = &cobra.Command{
	Use:   "day9 path/to/input/file",
	Short: "AoC Day 9",
	Long:  `Advent of Code Day 9: `,
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
	rootCmd.AddCommand(day9Cmd)
}
