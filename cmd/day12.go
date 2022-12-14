/*
Copyright © 2022 Jacob Saporito
*/
package cmd

import (
	"fmt"
	"math"
	"time"

	"github.com/japorito/merry/libxmas/rudolph"
	"github.com/japorito/merry/libxmas/toybag"
	"github.com/japorito/merry/libxmas/xmas"
	"github.com/spf13/cobra"
)

type HeightMap rudolph.AStarMap[rune]

func (hMap *HeightMap) Get(coords rudolph.Coordinate) *rudolph.MapNode[rune] {
	return hMap.MapNodes[coords.Row][coords.Col]
}

func (hMap *HeightMap) GetMap() [][]*rudolph.MapNode[rune] {
	return hMap.MapNodes
}

func (hMap *HeightMap) SetMap(m [][]*rudolph.MapNode[rune]) {
	hMap.MapNodes = m
}

func (hMap *HeightMap) Heuristic(src, dst rudolph.Coordinate) int {
	dManhattan := rudolph.ManhattanDistance(src, dst)
	dLetter := rudolph.DifferenceDistance(hMap.MapNodes[src.Row][src.Col].Value,
		hMap.MapNodes[dst.Row][dst.Col].Value)

	if dLetter > dManhattan {
		return dLetter
	}

	return dManhattan
}

func (hMap *HeightMap) TravelCost(src, dst rudolph.Coordinate) int {
	return hMap.MapNodes[src.Row][src.Col].TravelScore + 1
}

func (hMap *HeightMap) GetConnectedNodes(src rudolph.Coordinate) []rudolph.Coordinate {
	var neighbors []rudolph.Coordinate
	selfNode := hMap.MapNodes[src.Row][src.Col]

	if row := src.Row - 1; row >= 0 &&
		heightClimbable(selfNode.Value, hMap.MapNodes[row][src.Col].Value) {
		neighbors = append(neighbors, hMap.MapNodes[row][src.Col].Self)
	}
	if row := src.Row + 1; row < len(hMap.MapNodes) &&
		heightClimbable(selfNode.Value, hMap.MapNodes[row][src.Col].Value) {
		neighbors = append(neighbors, hMap.MapNodes[row][src.Col].Self)
	}
	if col := src.Col - 1; col >= 0 &&
		heightClimbable(selfNode.Value, hMap.MapNodes[src.Row][col].Value) {
		neighbors = append(neighbors, hMap.MapNodes[src.Row][col].Self)
	}
	if col := src.Col + 1; col < len(hMap.MapNodes[src.Row]) &&
		heightClimbable(selfNode.Value, hMap.MapNodes[src.Row][col].Value) {
		neighbors = append(neighbors, hMap.MapNodes[src.Row][col].Self)
	}

	return neighbors

}

func findAndReplaceRune(r, repl rune, m [][]rune) rudolph.Coordinate {
	for i := range m {
		for j := range m[i] {
			if m[i][j] == r {
				m[i][j] = repl

				return rudolph.Coordinate{
					Col: j,
					Row: i,
				}
			}
		}
	}

	return rudolph.ErrorCoordinate
}

func heightClimbable(src, dst rune) bool {
	return (dst - src) <= 1
}

// day12Cmd represents the day12 command
var day12Cmd = &cobra.Command{
	Use:   "day12 path/to/input/file",
	Short: "AoC Day 12",
	Long:  `Advent of Code Day 12: Hill Climbing Algorithm`,
	Run: func(cmd *cobra.Command, args []string) {
		if input := toybag.ReadToRuneSliceLines(args...); input != nil {
			fmt.Printf("%dX%d height map read.\n", len(input), len(input[0]))

			defer xmas.PrintHolidayMessage(time.Now())

			start := findAndReplaceRune('S', 'a', input)
			end := findAndReplaceRune('E', 'z', input)

			var heightMap rudolph.IAStarMap[rune] = &HeightMap{}
			rudolph.AStarInit(end, input, heightMap)

			if Parts.Has(1) {
				fmt.Println("Part 1 running...")

				rudolph.AStarRun(start, end, heightMap)

				fmt.Printf("The fastest route to the end point will get there in **%d** steps.\n",
					heightMap.Get(end).TravelScore)
			}

			if Parts.Has(2) {
				fmt.Println("Part 2 running...")

				min := math.MaxInt
				for row := range input {
					for col, val := range input[row] {
						if val == 'a' {
							rudolph.AStarReset(heightMap)
							rudolph.AStarRun(rudolph.Coordinate{Col: col, Row: row}, end, heightMap)

							if endNode := heightMap.Get(end); endNode.TravelScore < min {
								min = endNode.TravelScore
							}
						}
					}
				}

				fmt.Printf("The shortest hiking route is **%d** steps.\n", min)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(day12Cmd)
}
