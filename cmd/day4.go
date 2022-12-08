/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"
	"time"

	xmas "github.com/japorito/merry/libxmas"
	sleigh "github.com/japorito/merry/libxmas/input"
	"github.com/spf13/cobra"
)

func parseAssignmentRanges(input []string) [][]int64 {
	output := make([][]int64, len(input))
	for i, assignmentPair := range input {
		output[i] =
			xmas.ToInt64s(
				strings.FieldsFunc(assignmentPair,
					func(r rune) bool {
						return r == ',' || r == '-'
					}))
	}

	return output
}

func sectionRangeContains(range1, range2 []int64) bool {
	return range1[0] <= range2[0] && range1[1] >= range2[1]
}

func eitherIsSubset(range1, range2 []int64) bool {
	return sectionRangeContains(range1, range2) || sectionRangeContains(range2, range1)
}

func countCompleteSectionSubsets(assignmentPairs [][]int64) int64 {
	var subsets int64 = 0

	for _, rangePair := range assignmentPairs {
		if eitherIsSubset(rangePair[:2], rangePair[2:]) {
			subsets++
		}
	}

	return subsets
}

func sectionRangeOverlaps(range1, range2 []int64) bool {
	return (range1[0] <= range2[0] && range1[1] >= range2[0]) ||
		(range2[0] <= range1[0] && range2[1] >= range1[0])
}

func countOverlaps(assignmentPairs [][]int64) int64 {
	var overlaps int64 = 0

	for _, rangePair := range assignmentPairs {
		if sectionRangeOverlaps(rangePair[:2], rangePair[2:]) {
			overlaps++
		}
	}

	return overlaps
}

// day4Cmd represents the day4 command
var day4Cmd = &cobra.Command{
	Use:   "day4 path/to/input/file",
	Short: "AoC Day 4",
	Long:  `Advent of Code Day 4: Camp Cleanup`,
	Run: func(cmd *cobra.Command, args []string) {
		defer xmas.PrintHolidayMessage(time.Now())

		if input := sleigh.ReadToLines(args...); input != nil {
			fmt.Printf("%d assignment pairs read.\n", len(input))

			assignmentPairs := parseAssignmentRanges(input)

			if Parts.Has(1) {
				fmt.Println("Part 1 running...")
				fmt.Printf("There are **%d** elves with cleaning assignments containing their partner's assignments.\n",
					countCompleteSectionSubsets(assignmentPairs))
			}

			if Parts.Has(2) {
				fmt.Println("Part 2 running...")
				fmt.Printf("There are **%d** elf pairs with overlapping cleaning assignments.\n", countOverlaps(assignmentPairs))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(day4Cmd)
}
