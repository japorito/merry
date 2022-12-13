/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"math"
	"time"

	"github.com/japorito/merry/libxmas/stockings"
	"github.com/japorito/merry/libxmas/toybag"
	"github.com/japorito/merry/libxmas/xmas"
	"github.com/spf13/cobra"
)

type aStarNode struct {
	value          rune
	generatedScore int
	heuristicScore int
	predecessor    Coordinate
	self           Coordinate
}

func manhattanDistance(a, b Coordinate) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func letterDistance(a rune, b rune) int {
	return int(b - a)
}

func heuristic(src, dst Coordinate, m [][]rune) int {
	dManhattan := manhattanDistance(src, dst)
	dLetter := letterDistance(m[src.y][src.x], m[dst.y][dst.x])

	if dLetter > dManhattan {
		return dLetter
	}

	return dManhattan
}

func findAndReplaceRune(r, repl rune, m [][]rune) Coordinate {
	for i := range m {
		for j := range m[i] {
			if m[i][j] == r {
				m[i][j] = repl

				return Coordinate{
					x: j,
					y: i,
				}
			}
		}
	}

	return Coordinate{-1, -1}
}

func createAStarMap(end Coordinate, m [][]rune) [][]*aStarNode {
	heightMap := make([][]*aStarNode, len(m))
	for row := range m {
		heightMap[row] = make([]*aStarNode, len(m[row]))
		for col := range m[row] {
			selfCoords := Coordinate{x: col, y: row}
			heightMap[row][col] = &aStarNode{
				value:          m[row][col],
				generatedScore: math.MaxInt,
				heuristicScore: heuristic(selfCoords, end, m),
				self:           selfCoords,
			}
		}
	}

	return heightMap
}

func resetAStar(heightMap [][]*aStarNode) {
	for _, row := range heightMap {
		for _, node := range row {
			node.generatedScore = math.MaxInt
		}
	}
}

func heightClimbable(src, dst rune) bool {
	return (dst - src) <= 1
}

func getNeighbors(self Coordinate, heightMap [][]*aStarNode) []*aStarNode {
	var neighbors []*aStarNode
	selfNode := heightMap[self.y][self.x]

	if row := self.y - 1; row >= 0 &&
		heightClimbable(selfNode.value, heightMap[row][self.x].value) {
		neighbors = append(neighbors, heightMap[row][self.x])
	}
	if row := self.y + 1; row < len(heightMap) &&
		heightClimbable(selfNode.value, heightMap[row][self.x].value) {
		neighbors = append(neighbors, heightMap[row][self.x])
	}
	if col := self.x - 1; col >= 0 &&
		heightClimbable(selfNode.value, heightMap[self.y][col].value) {
		neighbors = append(neighbors, heightMap[self.y][col])
	}
	if col := self.x + 1; col < len(heightMap[self.y]) &&
		heightClimbable(selfNode.value, heightMap[self.y][col].value) {
		neighbors = append(neighbors, heightMap[self.y][col])
	}

	return neighbors
}

func runAStar(start, end Coordinate, heightMap [][]*aStarNode) {
	startNode := heightMap[start.y][start.x]
	startNode.generatedScore = 0

	endNode := heightMap[end.y][end.x]

	processQueue := stockings.NewMinPriorityQueue(32, func(item *aStarNode) int {
		return item.generatedScore + item.heuristicScore
	})
	processQueue.Add(startNode)

	var current *aStarNode
	for current = processQueue.GetNext(); current != endNode; current = processQueue.GetNext() {
		nextScore := current.generatedScore + 1

		for _, neighbor := range getNeighbors(current.self, heightMap) {
			if neighbor.generatedScore > nextScore {
				neighbor.generatedScore = nextScore
				neighbor.predecessor = current.self

				if processQueue.Has(neighbor) {
					processQueue.TryIncreasePriority(neighbor)
				} else {
					processQueue.Add(neighbor)
				}
			}
		}

		if processQueue.Size() == 0 {
			break
		}
	}
}

// day12Cmd represents the day12 command
var day12Cmd = &cobra.Command{
	Use:   "day12 path/to/input/file",
	Short: "AoC Day 12",
	Long:  `Advent of Code Day 12: `,
	Run: func(cmd *cobra.Command, args []string) {
		if input := toybag.ReadToRuneSliceLines(args...); input != nil {
			fmt.Printf("%d input lines read.\n", len(input))

			defer xmas.PrintHolidayMessage(time.Now())

			start := findAndReplaceRune('S', 'a', input)
			end := findAndReplaceRune('E', 'z', input)

			heightMap := createAStarMap(end, input)

			if Parts.Has(1) {
				fmt.Println("Part 1 running...")

				runAStar(start, end, heightMap)

				fmt.Printf("The fastest route to the end point will get there in **%d** steps.\n",
					heightMap[end.y][end.x].generatedScore)
			}

			if Parts.Has(2) {
				fmt.Println("Part 2 running...")

				min := math.MaxInt
				for row := range input {
					for col, val := range input[row] {
						if val == 'a' {
							resetAStar(heightMap)
							runAStar(Coordinate{x: col, y: row}, end, heightMap)

							if heightMap[end.y][end.x].generatedScore < min {
								min = heightMap[end.y][end.x].generatedScore
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
