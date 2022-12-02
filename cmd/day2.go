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

func makePointsSlice(length int) [][]int32 {
	return [][]int32{make([]int32, length), make([]int32, length)}
}

func getWinningMovePoints(opponentPts int32) int32 {
	winningPoints := opponentPts + 1
	if winningPoints == 4 {
		winningPoints = 1
	}

	return winningPoints
}

func roShamBo(player1, player2 int32) (int32, int32) {
	if player1 == player2 { // draw
		return 3, 3
	}

	if player1 == getWinningMovePoints(player2) {
		return 6, 0
	}

	//loss
	return 0, 6
}

func movePointCalculator(move string, oneVal rune) int32 {
	zeroVal := oneVal - 1

	return []rune(move)[0] - zeroVal
}

func cipherToPoints(cipher [][]string) [][]int32 {
	output := makePointsSlice(len(cipher))
	for idx, game := range cipher {
		output[0][idx] = movePointCalculator(game[0], 'A')
		output[1][idx] = movePointCalculator(game[1], 'X')
	}

	return output
}

func sumInt32(input []int32) int32 {
	sum := int32(0)
	for _, addend := range input {
		sum = sum + addend
	}

	return sum
}

func cipherToRoundScore(cipher [][]string) []int32 {
	const Lose, Draw, Win string = "X", "Y", "Z"

	roundScores := make([]int32, len(cipher))
	for i, game := range cipher {
		switch game[1] {
		case Draw:
			roundScores[i] = movePointCalculator(game[0], 'A') + 3
		case Win:
			roundScores[i] = getWinningMovePoints(movePointCalculator(game[0], 'A')) + 6
		case Lose:
			score := movePointCalculator(game[0], 'A') - 1
			if score == 0 {
				score = 3
			}

			roundScores[i] = score
		}
	}

	return roundScores
}

// day2Cmd represents the day2 command
var day2Cmd = &cobra.Command{
	Use:   "day2 path/to/input/file",
	Short: "Advent of Code Day 2",
	Long:  `Advent of Code Day 2: Elf Calories`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		var runall bool = Part == "*"

		input, err := xmas.ReadFileAsTokenizedStringSlice(args[0])
		if err != nil {
			return err
		}
		fmt.Printf("%d games of rock, paper, scissors in the cipher.\n", len(input))

		if runall || Part == "1" {
			fmt.Println("Part 1 running...")

			movePoints := cipherToPoints(input)

			finalPoints := makePointsSlice(len(movePoints[0]))
			for i := 0; i < len(movePoints[0]); i++ {
				finalPoints[0][i], finalPoints[1][i] = roShamBo(movePoints[0][i], movePoints[1][i])

				finalPoints[0][i] = finalPoints[0][i] + movePoints[0][i]
				finalPoints[1][i] = finalPoints[1][i] + movePoints[1][i]
			}

			myScore := sumInt32(finalPoints[1])

			fmt.Printf("My rock, paper, scissors score following the cipher would be **%d**.\n", myScore)
		}

		if runall || Part == "2" {
			fmt.Println("Part 2 running...")

			finalScore := sumInt32(cipherToRoundScore(input))
			fmt.Printf("My rock, paper, scissor score when I actually understand the cipher would be **%d**.\n", finalScore)
		}

		xmas.PrintHolidayMessage(time.Since(start))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(day2Cmd)
}
