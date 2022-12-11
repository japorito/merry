/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/japorito/merry/libxmas/sleigh"
	"github.com/japorito/merry/libxmas/toybag"
	"github.com/japorito/merry/libxmas/xmas"
	"github.com/spf13/cobra"
)

type Monkey struct {
	name            string
	items           []*int
	inspect         func(int) int
	test            func(int) bool
	trueTarget      int
	falseTarget     int
	inspectionCount *int
}

func product(a, b int) int {
	return a * b
}

func sum(a, b int) int {
	return a + b
}

func createInspect(fDef []string) func(int) int {
	arg1, op, arg2 := fDef[3], fDef[4], fDef[5]

	var mathFunc func(int, int) int
	switch op {
	case "*":
		mathFunc = product
	case "+":
		mathFunc = sum
	}

	arg1Old := (arg1 == "old")
	arg2Old := (arg2 == "old")
	if arg1Old && arg2Old {
		return func(old int) int {
			return mathFunc(old, old)
		}
	} else if arg1Old {
		b, _ := strconv.Atoi(arg2)
		return func(old int) int {
			return mathFunc(old, b)
		}
	} else if arg2Old {
		a, _ := strconv.Atoi(arg1)
		return func(old int) int {
			return mathFunc(a, old)
		}
	} else {
		a, _ := strconv.Atoi(arg1)
		b, _ := strconv.Atoi(arg2)
		return func(int) int {
			return mathFunc(a, b)
		}
	}
}

func createTest(testDef []string, denominatorProduct *int) func(int) bool {
	arg, _ := strconv.Atoi(testDef[3])
	*denominatorProduct = *denominatorProduct * arg

	return func(tested int) bool {
		return tested%arg == 0
	}
}

func getTarget(targetDef []string) int {
	target, _ := strconv.Atoi(targetDef[5])
	return target
}

func parseMonkeys(input [][]string) ([]*Monkey, int) {
	monkeys := make([]*Monkey, len(input))

	denominatorProduct := new(int)
	*denominatorProduct = 1
	for i, block := range input {
		monkeyDef := sleigh.Tokenize(block)

		monkey := &Monkey{}
		monkey.name = strings.Join(monkeyDef[0], " ")
		for _, item := range monkeyDef[1][2:] {
			itemNo, _ := strconv.Atoi(strings.TrimRight(item, ","))
			monkey.items = append(monkey.items, &itemNo)
		}
		monkey.inspect = createInspect(monkeyDef[2])
		monkey.inspectionCount = new(int)
		monkey.test = createTest(monkeyDef[3], denominatorProduct)
		monkey.trueTarget = getTarget(monkeyDef[4])
		monkey.falseTarget = getTarget(monkeyDef[5])

		monkeys[i] = monkey
	}

	return monkeys, *denominatorProduct
}

func sighOfRelief(item int) int {
	return item / 3
}

func simianSimulation(monkeys []*Monkey, rounds int, worryManagement func(int) int) {
	for i := 0; i < rounds; i++ {
		for _, monkey := range monkeys {
			*monkey.inspectionCount += len(monkey.items)
			for _, item := range monkey.items {
				*item = monkey.inspect(*item)
				*item = worryManagement(*item)

				if monkey.test(*item) {
					monkeys[monkey.trueTarget].items = append(monkeys[monkey.trueTarget].items, item)
				} else {
					monkeys[monkey.falseTarget].items = append(monkeys[monkey.falseTarget].items, item)
				}
			}

			monkey.items = monkey.items[:0]
		}
	}
}

func calculateMonkeyBusiness(monkeys []*Monkey) int {
	monkeysCopy := make([]*Monkey, len(monkeys))
	copy(monkeysCopy, monkeys)
	sort.SliceStable(monkeysCopy, func(i, j int) bool {
		return *monkeysCopy[j].inspectionCount < *monkeysCopy[i].inspectionCount
	})

	return *monkeysCopy[0].inspectionCount * *monkeysCopy[1].inspectionCount
}

// day11Cmd represents the day11 command
var day11Cmd = &cobra.Command{
	Use:   "day11 path/to/input/file",
	Short: "AoC Day 11",
	Long:  `Advent of Code Day 11: Monkey in the Middle`,
	Run: func(cmd *cobra.Command, args []string) {
		if input := toybag.ReadToStringSliceBlocks(args...); input != nil {
			fmt.Printf("%d input lines read.\n", len(input))

			defer xmas.PrintHolidayMessage(time.Now())

			if Parts.Has(1) {
				fmt.Println("Part 1 running...")

				monkeys, _ := parseMonkeys(input)
				simianSimulation(monkeys, 20, sighOfRelief)

				fmt.Printf("The monkey business at round 20 is **%d**.\n", calculateMonkeyBusiness(monkeys))
			}

			if Parts.Has(2) {
				fmt.Println("Part 2 running...")

				monkeys, denominatorProduct := parseMonkeys(input)
				simianSimulation(monkeys, 10000, func(worry int) int {
					return worry % denominatorProduct
				})

				fmt.Printf("The monkey business at round 10000 is **%d**.\n", calculateMonkeyBusiness(monkeys))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(day11Cmd)
}
