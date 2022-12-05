/*
Copyright © 2022 Jacob Saporito
*/
package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	xmas "github.com/japorito/merry/libxmas"
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

	labelline := stacklisting[labelidx]

	return strings.Fields(labelline)
}

func parseStacks(labels, stacklisting []string) map[string]*xmas.Stack[rune] {
	maxstackheight := len(stacklisting) - 1

	labelline := stacklisting[maxstackheight]
	stackdescriptions := stacklisting[:maxstackheight]

	stackcol := make(map[string]int)
	for _, label := range labels {
		stackcol[label] = strings.Index(labelline, label)
	}

	stacks := make(map[string]*xmas.Stack[rune], len(labels))
	for stackname, stackidx := range stackcol {
		for i := len(stackdescriptions) - 1; i >= 0; i-- {
			line := stackdescriptions[i]
			layerdescription := []rune(line)

			if stackidx < len(layerdescription) && layerdescription[stackidx] != ' ' {
				stack, ok := stacks[stackname]
				if !ok {
					stack = &xmas.Stack[rune]{}
					stacks[stackname] = stack
				}
				stack.Push(layerdescription[stackidx])
			}
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

func applyInstructions(stacks map[string]*xmas.Stack[rune], instructions []instruction) {
	for _, inst := range instructions {
		for i := 0; i < inst.arg; i++ {
			fromstack := stacks[inst.src]
			tostack := stacks[inst.dst]

			tostack.Push(fromstack.Pop())
		}
	}
}

func applyInstructionsCrateMove9001(stacks map[string]*xmas.Stack[rune], instructions []instruction) {
	for _, inst := range instructions {
		fromstack := stacks[inst.src]
		tostack := stacks[inst.dst]

		tostack.Push(fromstack.Top(inst.arg)...)
	}
}

func printTop(intro string, labels []string, stacks map[string]*xmas.Stack[rune]) {
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
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		var runall bool = Part == "*"

		input, err := xmas.ReadFileToStringSliceBlocks(args[0])
		if err != nil {
			return err
		}
		maxstackheight := len(input[0]) - 1
		fmt.Printf("%d max stack height read and %d crane instructions read.\n", maxstackheight, len(input[1]))

		labels := parseLabels(input[0])
		instructions := parseInstructions(input[1])

		if runall || Part == "1" {
			fmt.Println("Part 1 running...")

			stacks := parseStacks(labels, input[0])

			applyInstructions(stacks, instructions)

			printTop("After crane operations, the top of the stacks are: ", labels, stacks)
		}

		if runall || Part == "2" {
			fmt.Println("Part 2 running...")

			stacks := parseStacks(labels, input[0])

			applyInstructionsCrateMove9001(stacks, instructions)

			printTop("After CrateMover 9001 operations, the top of the stacks are: ", labels, stacks)
		}

		xmas.PrintHolidayMessage(time.Since(start))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(day5Cmd)
}
