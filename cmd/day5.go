/*
Copyright © 2022 Jacob Saporito
*/
package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/japorito/merry/libxmas/sleigh"
	"github.com/japorito/merry/libxmas/stockings"
	"github.com/japorito/merry/libxmas/toybag"
	"github.com/japorito/merry/libxmas/xmas"
	"github.com/spf13/cobra"
)

type instruction struct {
	command string
	arg     int
	src     string
	dst     string
}

func parseLabels(stacklisting []string) []string {
	labelidx := len(stacklisting) - 1

	return strings.Fields(stacklisting[labelidx])
}

func parseStacks(stacklisting []string) map[string]*stockings.Stack[rune] {
	crateRuneMatrix := sleigh.StandardizeDimensions(sleigh.ToRunes(stacklisting), ' ')
	crateRuneMatrix, _ = sleigh.TransposeMatrix(crateRuneMatrix)
	sleigh.ReverseHorizontal(&crateRuneMatrix)

	stacks := make(map[string]*stockings.Stack[rune])

	for _, row := range crateRuneMatrix {
		if row[0] != ' ' {
			stack := &stockings.Stack[rune]{}

			// some stacks are shorter than others, so remove ' ' trailing runes
			crates := []rune(strings.TrimSpace(string(row[1:])))
			stack.Push(crates...)
			stacks[string(row[0])] = stack
		}
	}

	return stacks
}

func parseInstructions(instructionset []string) []instruction {
	instructions := make([]instruction, len(instructionset))
	for idx, instructionline := range instructionset {
		fields := strings.Fields(instructionline)

		inst := instruction{
			command: fields[0],
			src:     fields[3],
			dst:     fields[5],
		}
		inst.arg, _ = strconv.Atoi(fields[1]) //assuming well-formed input

		instructions[idx] = inst
	}

	return instructions
}

func applyInstructions(stacks map[string]*stockings.Stack[rune], instructions []instruction) {
	for _, inst := range instructions {
		for i := 0; i < inst.arg; i++ {
			fromstack := stacks[inst.src]
			tostack := stacks[inst.dst]

			tostack.Push(fromstack.Pop())
		}
	}
}

func applyInstructionsCrateMove9001(stacks map[string]*stockings.Stack[rune], instructions []instruction) {
	for _, inst := range instructions {
		fromstack := stacks[inst.src]
		tostack := stacks[inst.dst]

		tostack.Push(fromstack.Top(inst.arg)...)
	}
}

func printTop(intro string, labels []string, stacks map[string]*stockings.Stack[rune]) {
	fmt.Print(intro, "**")
	for _, label := range labels {
		stack := stacks[label]
		fmt.Print(string(stack.Peek()))
	}

	fmt.Println("**")
}

// day5Cmd represents the day5 command
var day5Cmd = &cobra.Command{
	Use:   "day5 path/to/input/file",
	Short: "AoC Day 5",
	Long:  `Advent of Code Day 5: Supply Stacks`,
	Run: func(cmd *cobra.Command, args []string) {
		if input := toybag.ReadToStringSliceBlocks(args...); input != nil {
			maxstackheight := len(input[0]) - 1
			fmt.Printf("%d max stack height read and %d crane instructions read.\n", maxstackheight, len(input[1]))

			defer xmas.PrintHolidayMessage(time.Now())

			labels := parseLabels(input[0])
			instructions := parseInstructions(input[1])

			if Parts.Has(1) {
				fmt.Println("Part 1 running...")

				stacks := parseStacks(input[0])

				applyInstructions(stacks, instructions)

				printTop("After crane operations, the top of the stacks are: ", labels, stacks)
			}

			if Parts.Has(2) {
				fmt.Println("Part 2 running...")

				stacks := parseStacks(input[0])

				applyInstructionsCrateMove9001(stacks, instructions)

				printTop("After CrateMover 9001 operations, the top of the stacks are: ", labels, stacks)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(day5Cmd)
}
