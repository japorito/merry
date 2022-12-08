/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	xmas "github.com/japorito/merry/libxmas"
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

func calculateScenicScore(forest [][]*tree, x, y int, out chan int) {
	rowCount, rowLength := len(forest), len(forest[0])
	if x == 0 || y == 0 || y == (rowCount-1) || x == (rowLength-1) {
		out <- 0
		return
	}

	height := forest[y][x].height

	upCount, downCount, leftCount, rightCount := 0, 0, 0, 0
	lookUp, lookDown, lookLeft, lookRight := true, true, true, true
	for i := 1; i < rowCount; i++ {
		yUp, yDown := y-i, y+i
		if yUp < 0 && yDown > rowCount {
			// hit vertical boundaries
			break
		}

		if !lookUp && !lookDown {
			break
		}

		if lookUp && forest[yUp][x].height <= height {
			upCount = i
		}

		if yUp <= 0 || forest[yUp][x].height >= height {
			lookUp = false
		}

		if lookDown && forest[yDown][x].height <= height {
			downCount = i
		}

		if yDown >= (rowCount-1) || forest[yDown][x].height >= height {
			lookDown = false
		}
	}

	for i := 1; i < rowCount; i++ {
		xLeft, xRight := x-i, x+i
		if xLeft < 0 && xRight > rowCount {
			// hit vertical boundaries
			break
		}

		if !lookLeft && !lookRight {
			break
		}

		if lookLeft && forest[y][xLeft].height <= height {
			leftCount = i
		}

		if xLeft <= 0 || forest[y][xLeft].height >= height {
			lookLeft = false
		}

		if lookRight && forest[y][xRight].height <= height {
			rightCount = i
		}

		if xRight >= (rowLength-1) || forest[y][xRight].height >= height {
			lookRight = false
		}
	}

	out <- upCount * downCount * leftCount * rightCount
}

func findHighestScenicScore(forest [][]*tree) int {
	scores := make(chan int, 20)

	for y := range forest {
		for x := range forest[y] {
			go calculateScenicScore(forest, x, y, scores)
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
	Long:  `Advent of Code Day 8: `,
	//Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		state := NewMerryState(args)
		if input := xmas.ReadFileToRuneSliceLines(args[0]); input != nil {
			fmt.Printf("%d input lines read.\n", len(input))

			defer xmas.PrintHolidayMessage(time.Now())

			forest := createForest(input)
			determineVisibility(forest)

			if state.Parts.Has(1) {
				fmt.Println("Part 1 running...")
				fmt.Printf("Count of outwardly-visible trees in forest is **%d**.\n", countVisible(forest))
			}

			if state.Parts.Has(2) {
				fmt.Println("Part 2 running...")
				fmt.Println(findHighestScenicScore(forest))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(day8Cmd)
}
