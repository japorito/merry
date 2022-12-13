/*
Copyright © 2022 Jacob Saporito
*/
package cmd

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/japorito/merry/libxmas/toybag"
	"github.com/japorito/merry/libxmas/xmas"
	"github.com/spf13/cobra"
)

type distressCall []string

func (call *distressCall) Len() int {
	return len(*call)
}

func (call *distressCall) Less(i, j int) bool {
	return compare((*call)[i], (*call)[j]) > 0
}

func (call *distressCall) Swap(i, j int) {
	(*call)[i], (*call)[j] = (*call)[j], (*call)[i]
}

func isList(data string) bool {
	d := []rune(data)

	return len(d) > 0 && d[0] == '['
}

func parsePacketList(packetData string) []string {
	if !isList(packetData) {
		return []string{packetData}
	}

	packet := []rune(packetData)
	output := make([]string, 0, 4)

	packet = packet[1 : len(packet)-1] // strip []

	startCurrent := 0
	for braces, i := 0, 0; i < len(packet); i++ {
		switch packet[i] {
		case '[':
			braces++
		case ']':
			braces--
		case ',':
			if braces == 0 {
				output = append(output, string(packet[startCurrent:i]))
				startCurrent = i + 1
			}
		}
	}
	output = append(output, string(packet[startCurrent:]))

	return output
}

func compare(left, right string) int {
	if leftEmpty, rightEmpty := left == "", right == ""; leftEmpty && rightEmpty {
		return 0
	} else if leftEmpty {
		return 1 // left ran out first
	} else if rightEmpty {
		return -1 // right ran out first
	}

	if !isList(left) && !isList(right) {
		lNum, _ := strconv.Atoi(string(left))
		rNum, _ := strconv.Atoi(string(right))

		if lNum < rNum {
			return 1
		} else if lNum == rNum {
			// undecided
			return 0
		} else {
			return -1
		}
	}

	parsedL, parsedR := parsePacketList(left), parsePacketList(right)

	for i, dataL := range parsedL {
		if i >= len(parsedR) {
			return -1
		}

		entryCompare := compare(dataL, parsedR[i])
		if entryCompare != 0 {
			return entryCompare
		}
	}

	if len(parsedL) < len(parsedR) {
		return 1
	}

	return 0
}

func findCorrectIndices(input [][]string) int {
	sum := 0

	for i, pair := range input {
		pairCompare := compare(pair[0], pair[1])
		if pairCompare > 0 {
			sum += (i + 1)
		}
	}

	return sum
}

func recompilePairs(pairs [][]string, additional ...string) distressCall {
	output := make([]string, 0, len(pairs)*2+len(additional))

	for _, pair := range pairs {
		output = append(output, pair...)
	}

	output = append(output, additional...)

	return output
}

func findDecoderKey(input [][]string, dividers ...string) int {
	packetList := recompilePairs(input, dividers...)

	sort.Sort(&packetList)

	decoder := 1
	for _, divider := range dividers {
		dividerIdx := sort.Search(len(packetList), func(i int) bool {
			return compare(packetList[i], divider) <= 0
		})
		decoder *= (dividerIdx + 1)
	}

	return decoder
}

// day13Cmd represents the day13 command
var day13Cmd = &cobra.Command{
	Use:   "day13 path/to/input/file",
	Short: "AoC Day 13",
	Long:  `Advent of Code Day 13: `,
	Run: func(cmd *cobra.Command, args []string) {
		if input := toybag.ReadToStringSliceBlocks(args...); input != nil {
			fmt.Printf("%d packet pairs read read.\n", len(input))

			defer xmas.PrintHolidayMessage(time.Now())

			if Parts.Has(1) {
				fmt.Println("Part 1 running...")

				fmt.Printf("The sum of the indices of the correct pairs are **%d**.\n", findCorrectIndices(input))
			}

			if Parts.Has(2) {
				fmt.Println("Part 2 running...")

				decoder := findDecoderKey(input, "[[2]]", "[[6]]")

				fmt.Printf("The decoder key for the unscrambled input is **%d**.\n", decoder)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(day13Cmd)
}
