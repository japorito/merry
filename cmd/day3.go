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

func runeToPriority(input [][]rune) [][]int64 {
	const lowercaseZero int32 = 'a' - 1
	const uppercaseZero int32 = 'A' - 1
	const alphabetLength = 26

	priorities := make([][]int64, len(input))
	for i, runes := range input {
		priorities[i] = make([]int64, len(runes))
		for j, letter := range runes {
			if letter <= 'Z' {
				priorities[i][j] = int64(alphabetLength + letter - uppercaseZero)
			} else {
				priorities[i][j] = int64(letter - lowercaseZero)
			}
		}
	}

	return priorities
}

// find items repeated in all slices
func repeatedPriorities(itemlists ...xmas.BitSet) []int64 {
	repeatedItems := itemlists[0]

	repeatedItems = repeatedItems.Intersect(itemlists[1:]...)

	return repeatedItems.Members()
}

func createSet(items []int64) xmas.BitSet {
	set := xmas.BitSet{}
	for _, item := range items {
		set.On(item)
	}

	return set
}

func createSets(itemgroups ...[]int64) []xmas.BitSet {
	sets := make([]xmas.BitSet, 0, len(itemgroups))

	for _, itemgroup := range itemgroups {
		sets = append(sets, createSet(itemgroup))
	}

	return sets
}

func compartmentRepeatedPriorities(rucksack []int64) []int64 {
	compartmentCapacity := len(rucksack) / 2
	return repeatedPriorities(createSets(rucksack[:compartmentCapacity], rucksack[compartmentCapacity:])...)
}

func rucksackErrors(allRucksacks [][]int64) []int64 {
	var errors []int64
	for _, rucksack := range allRucksacks {
		errors = append(errors, compartmentRepeatedPriorities(rucksack)...)
	}

	return errors
}

func findBadges(allRucksacks [][]int64) []int64 {
	var badges []int64

	for i := 3; i <= len(allRucksacks); i += 3 {
		badges = append(badges, repeatedPriorities(createSets(allRucksacks[(i-3):i]...)...)...)
	}

	return badges
}

// day3Cmd represents the day3 command
var day3Cmd = &cobra.Command{
	Use:   "day3 path/to/input/file",
	Short: "AoC Day 3",
	Long:  `Advent of Code Day 3: Needs Title`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		var runall bool = Part == "*"

		input, err := xmas.ReadFileToRuneSliceLines(args[0])
		if err != nil {
			return err
		}
		fmt.Printf("%d rucksacks packed.\n", len(input))

		priorities := runeToPriority(input)

		if runall || Part == "1" {
			fmt.Println("Part 1 running...")

			errors := rucksackErrors(priorities)

			fmt.Printf("Combined priorities of rucksacking-packing errors is **%d**\n", xmas.SumInt64(errors))
		}

		if runall || Part == "2" {
			fmt.Println("Part 2 running...")

			badges := findBadges(priorities)

			fmt.Printf("The sum of all %d badge type priorities is **%d**\n", len(badges), xmas.SumInt64(badges))
		}

		xmas.PrintHolidayMessage(time.Since(start))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(day3Cmd)
}
