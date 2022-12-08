/*
Copyright © 2022 Jacob Saporito

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"os"

	"github.com/japorito/merry/libxmas/stockings"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "merry",
	Short: "Advent of Code application for 2022",
	Long: `Advent of Code application for 2022
	
Planned use will look something like:
merry day1 [--part=1] /path/to/input`,
}

var Parts stockings.BitSet[int]

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	var part string

	rootCmd.PersistentFlags().StringVarP(&part, "Part", "p", "*", "Which puzzle part to run.")

	allParts := (part == "*")
	if allParts || part == "1" {
		Parts.On(1)
	}

	if allParts || part == "2" {
		Parts.On(2)
	}
}
