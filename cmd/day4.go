/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"
	"time"

	xmas "github.com/japorito/merry/libxmas"
	"github.com/spf13/cobra"
)

func parseAssignmentRanges(input []string) ([][]int64, error) {
	output := make([][]int64, len(input))
	var err error
	for i, assignmentPair := range input {
		output[i], err =
			xmas.ToInt64s(
				strings.FieldsFunc(assignmentPair,
					func(r rune) bool {
						return r == ',' || r == '-'
					}))
		if err != nil {
			return nil, err
		}
	}

	return output, nil
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
	Long:  `Advent of Code Day 4`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		var runall bool = Part == "*"

		input, err := xmas.ReadFileToLines(args[0])
		if err != nil {
			return err
		}
		fmt.Printf("%d assignment pairs read.\n", len(input))

		assignmentPairs, err := parseAssignmentRanges(input)
		if err != nil {
			return err
		}

		if runall || Part == "1" {
			fmt.Println("Part 1 running...")
			fmt.Printf("There are **%d** elves with cleaning assignments containing their partner's assignments.\n",
				countCompleteSectionSubsets(assignmentPairs))
		}

		if runall || Part == "2" {
			fmt.Println("Part 2 running...")
			fmt.Printf("There are **%d** elf pairs with overlapping cleaning assignments.", countOverlaps(assignmentPairs))
		}

		xmas.PrintHolidayMessage(time.Since(start))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(day4Cmd)
}
