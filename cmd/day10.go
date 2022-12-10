/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/japorito/merry/libxmas/toybag"
	"github.com/japorito/merry/libxmas/xmas"
	"github.com/spf13/cobra"
)

type screenState struct {
	clock, register, signal *int
}

func newState() screenState {
	var clock, register, signal *int = new(int), new(int), new(int)
	*register = 1 // initial value
	*clock = 1

	return screenState{clock: clock, register: register, signal: signal}
}

func lightPixel(state screenState) {
	scanPos := (*state.clock - 1) % 40

	if scanPos == 0 && *state.clock > 1 {
		fmt.Println()
	}

	if abs(*state.register-(scanPos)) <= 1 {
		fmt.Printf("#")
	} else {
		fmt.Printf(" ")
	}
}

func sumSignal(state screenState) {
	if *state.clock <= 220 && (*state.clock-20)%40 == 0 {
		*state.signal += (*state.register * *state.clock)
	}
}

func tick(state screenState, f func(screenState)) {
	f(state)

	*state.clock++
}

func noop(state screenState, f func(screenState)) {
	tick(state, f)
}

func addx(state screenState, amount int, f func(screenState)) {
	tick(state, f)
	tick(state, f)

	*state.register += amount
}

func run(code [][]string, f func(screenState)) screenState {
	state := newState()

	for _, command := range code {
		switch command[0] {
		case "noop":
			noop(state, f)
		case "addx":
			arg, _ := strconv.Atoi(command[1])
			addx(state, arg, f)
		}
	}

	return state
}

// day10Cmd represents the day10 command
var day10Cmd = &cobra.Command{
	Use:   "day10 path/to/input/file",
	Short: "AoC Day 10",
	Long:  `Advent of Code Day 10: `,
	Run: func(cmd *cobra.Command, args []string) {
		if input := toybag.ReadAsTokenizedStringSlice(args...); input != nil {
			fmt.Printf("%d input lines read.\n", len(input))

			defer xmas.PrintHolidayMessage(time.Now())

			if Parts.Has(1) {
				fmt.Println("Part 1 running...")
				endstate := run(input, sumSignal)

				fmt.Printf("The sum of the signal strengths are **%d**.\n", *endstate.signal)
			}

			if Parts.Has(2) {
				fmt.Println("Part 2 running...")

				run(input, lightPixel)
				fmt.Println()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(day10Cmd)
}
