/*
Copyright © 2022 Jacob Saporito
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

var directions = map[string]*Coordinate{
	"R": {x: 1, y: 0},
	"L": {x: -1, y: 0},
	"U": {x: 0, y: 1},
	"D": {x: 0, y: -1},
}

type Coordinate struct {
	x, y int
}

func moveHead(head, direction *Coordinate) {
	head.x = head.x + direction.x
	head.y = head.y + direction.y
}

func abs(num int) int {
	if num < 0 {
		num = -num
	}
	return num
}

func moveTail(head, tail *Coordinate) {
	dx, dy := (head.x - tail.x), (head.y - tail.y)

	// only move if not already close
	if abs(dx) > 1 || abs(dy) > 1 {
		if dx != 0 {
			dx = dx / abs(dx)
		}
		if dy != 0 {
			dy = dy / abs(dy)
		}

		tail.x = tail.x + dx
		tail.y = tail.y + dy
	}
}

func runMoves(input [][]string, knots int) map[string]bool {
	rope := make([]Coordinate, knots)
	head := &rope[0]
	tail := &rope[knots-1]
	tailVisits := make(map[string]bool)

	for _, move := range input {
		direction := directions[move[0]]
		count, _ := strconv.Atoi(move[1])

		for i := 0; i < count; i++ {
			moveHead(head, direction)

			for j := 1; j < knots; j++ {
				moveTail(&rope[j-1], &rope[j])
			}

			key := fmt.Sprintf("%d:%d", tail.x, tail.y)
			tailVisits[key] = true
		}
	}

	return tailVisits
}

// day9Cmd represents the day9 command
var day9Cmd = &cobra.Command{
	Use:   "day9 path/to/input/file",
	Short: "AoC Day 9",
	Long:  `Advent of Code Day 9: Rope Bridge`,
	Run: func(cmd *cobra.Command, args []string) {
		if input := toybag.ReadAsTokenizedStringSlice(args...); input != nil {
			fmt.Printf("%d input lines read.\n", len(input))

			defer xmas.PrintHolidayMessage(time.Now())

			if Parts.Has(1) {
				fmt.Println("Part 1 running...")
				visits := runMoves(input, 2)

				fmt.Printf("Rope tail visits **%d** coordinates.\n", len(visits))
			}

			if Parts.Has(2) {
				fmt.Println("Part 2 running...")
				visits := runMoves(input, 10)

				fmt.Printf("Longer rope tail visits **%d** coordinates.\n", len(visits))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(day9Cmd)
}
