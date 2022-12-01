/*
Copyright © 2022 Jacob Saporito
*/
package cmd

import (
	"fmt"

	xmas "github.com/japorito/merry/libxmas"
	"github.com/spf13/cobra"
)

// day1Cmd represents the day1 command
var day1Cmd = &cobra.Command{
	Use:   "day1",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var runall bool = Part == "*"

		input, err := xmas.ReadFileToInt64SliceBlocks(args[0])
		if err != nil {
			return err
		}
		fmt.Printf("%d input blocks read.\n", len(input))

		if runall || Part == "1" {
			fmt.Println("Part 1 running...")

			var maxSum int64 = 0
			for _, elf := range input {
				var sum int64 = 0
				for _, snack := range elf {
					sum = sum + snack
				}

				if sum > maxSum {
					maxSum = sum
				}
			}

			fmt.Printf("Answer 1: %d\n", maxSum)
		}

		if runall || Part == "2" {
			fmt.Println("Part 2 running...")
		}

		xmas.PrintHolidayMessage()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(day1Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// day1Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// day1Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
