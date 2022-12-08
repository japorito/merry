/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	xmas "github.com/japorito/merry/libxmas"
	sleigh "github.com/japorito/merry/libxmas/input"
	"github.com/spf13/cobra"
)

type tree struct {
	height      rune
	visible     bool
	scenicScore int
}

func createForest(input [][]rune) [][]*tree {
	rowCount := len(input)
	rowLength := len(input[0])
	trees := make([][]*tree, rowCount)

	for rowNum, row := range input {
		trees[rowNum] = make([]*tree, rowLength)
		for colNum, height := range row {
			trees[rowNum][colNum] = &tree{
				height: height,
				// outside edge always visible
				visible:     rowNum == 0 || colNum == 0 || colNum == (rowLength-1) || rowNum == (rowCount-1),
				scenicScore: 1,
			}
		}
	}

	return trees
}

func determineVisibility(forest [][]*tree) [][]*tree {
	rowCount := len(forest)
	rowLength := len(forest[0])

	// right-left
	for _, row := range forest {
		maxLeft, maxRight := '0', '0'
		for colNum, tree := range row {
			if tree.height > maxLeft {
				maxLeft = tree.height
				tree.visible = true
			}

			reverseTree := row[(rowLength-1)-colNum]
			if reverseTree.height > maxRight {
				maxRight = reverseTree.height
				reverseTree.visible = true
			}
		}
	}

	// top-bottom
	for i := 0; i < rowLength; i++ {
		maxTop, maxBottom := '0', '0'
		for rowNum := range forest {
			tree := forest[rowNum][i]
			reverseTree := forest[-1+rowCount-rowNum][i]

			if tree.height > maxTop {
				maxTop = tree.height
				tree.visible = true
			}

			if reverseTree.height > maxBottom {
				maxBottom = reverseTree.height
				reverseTree.visible = true
			}
		}
	}

	return forest
}

func checkHeightBlocksView(forest [][]*tree, col, row int, height rune) bool {
	treeHeight := forest[row][col].height
	return treeHeight >= height
}

func calculateScenicScore(forest [][]*tree, col, row, cols, rows int, out chan int) {
	if col == 0 || row == 0 || row == (rows-1) || col == (cols-1) {
		out <- 0
		return
	}

	height := forest[row][col].height
	upCount, downCount, leftCount, rightCount := 0, 0, 0, 0
	i := 0
	for i = row - 1; i > 0; i-- {
		if checkHeightBlocksView(forest, col, i, height) {
			break
		}
	}
	upCount = row - i
	for i = row + 1; i < (rows - 1); i++ {
		if checkHeightBlocksView(forest, col, i, height) {
			break
		}
	}
	downCount = i - row
	for i = col - 1; i > 0; i-- {
		if checkHeightBlocksView(forest, i, row, height) {
			break
		}
	}
	leftCount = col - i
	for i = col + 1; i < (cols - 1); i++ {
		if checkHeightBlocksView(forest, i, row, height) {
			break
		}
	}
	rightCount = i - col

	out <- upCount * downCount * leftCount * rightCount
}

func findHighestScenicScore(forest [][]*tree) int {
	scores := make(chan int, 15)
	rows, cols := len(forest), len(forest[0])

	for row := range forest {
		for col := range forest[row] {
			go calculateScenicScore(forest, col, row, cols, rows, scores)
		}
	}

	maxScore := -1
	for i := 0; i < (len(forest) * len(forest[0])); i++ {
		score := <-scores
		if score > maxScore {
			maxScore = score
		}
	}

	return maxScore
}

func countVisible(forest [][]*tree) int {
	sum := 0
	for _, row := range forest {
		for _, tree := range row {
			if tree.visible {
				sum++
			}
		}
	}

	return sum
}

// day8Cmd represents the day8 command
var day8Cmd = &cobra.Command{
	Use:   "day8 path/to/input/file",
	Short: "AoC Day 8",
	Long:  `Advent of Code Day 8: Treetop Tree House`,
	Run: func(cmd *cobra.Command, args []string) {
		if input := sleigh.ReadToRuneSliceLines(args...); input != nil {
			fmt.Printf("%d rows of trees read.\n", len(input))

			defer xmas.PrintHolidayMessage(time.Now())

			forest := createForest(input)
			determineVisibility(forest)

			if Parts.Has(1) {
				fmt.Println("Part 1 running...")
				fmt.Printf("Count of outwardly-visible trees in forest is **%d**.\n", countVisible(forest))
			}

			if Parts.Has(2) {
				fmt.Println("Part 2 running...")
				fmt.Printf("Top tree has Scenic Score: **%d**\n", findHighestScenicScore(forest))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(day8Cmd)
}
