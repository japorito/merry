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

func getWinningMovePoints(opponentPts int32) int32 {
	return (opponentPts % 3) + 1 // 1 "Rock" beats 3 "Scissors"
}

func roShamBo(player1, player2 int32) (int32, int32) {
	const DrawPoints, WinPoints, LosePoints int32 = 3, 6, 0

	switch player1 {
	case player2:
		return DrawPoints + player1, DrawPoints + player2
	case getWinningMovePoints(player2):
		return WinPoints + player1, LosePoints + player2
	default:
		return LosePoints + player1, WinPoints + player2
	}
}

func movePointCalculator(move string, oneVal rune) int32 {
	zeroVal := oneVal - 1

	return []rune(move)[0] - zeroVal
}

func misunderstoodCipherToPoints(cipher [][]string) int32 {
	pts := makePointsSlice(len(cipher))
	for idx, game := range cipher {
		pts[0][idx], pts[1][idx] = roShamBo(movePointCalculator(game[0], 'A'), movePointCalculator(game[1], 'X'))
	}

	return sumInt32(pts[1])
}

func cipherToPoints(cipher [][]string) int32 {
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
				score = 3 // 1 "Rock" beats 3 "Scissors"
			}

			roundScores[i] = score
		}
	}

	return sumInt32(roundScores)
}

func sumInt32(input []int32) int32 {
	sum := int32(0)
	for _, addend := range input {
		sum += addend
	}

	return sum
}

func makePointsSlice(length int) [][]int32 {
	return [][]int32{make([]int32, length), make([]int32, length)}
}

// day2Cmd represents the day2 command
var day2Cmd = &cobra.Command{
	Use:   "day2 path/to/input/file",
	Short: "Advent of Code Day 2",
	Long:  `Advent of Code Day 2: Rock Paper Scissors`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		state := NewMerryState(args)
		Parts := state.Parts
		defer xmas.PrintHolidayMessage(time.Now())

		if input := xmas.ReadFileAsTokenizedStringSlice(args[0]); input != nil {
			fmt.Printf("%d games of rock, paper, scissors in the cipher.\n", len(input))

			if Parts.Has(1) {
				fmt.Println("Part 1 running...")
				fmt.Printf("My rock, paper, scissors score following the cipher would be **%d**.\n", misunderstoodCipherToPoints(input))
			}

			if Parts.Has(2) {
				fmt.Println("Part 2 running...")
				fmt.Printf("My rock, paper, scissor score when I actually understand the cipher would be **%d**.\n", cipherToPoints(input))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(day2Cmd)
}
