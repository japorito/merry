/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"math"
	"strconv"
	"time"

	xmas "github.com/japorito/merry/libxmas"
	"github.com/spf13/cobra"
)

const hardDriveSize int = 70000000
const updateSize int = 30000000

type WeightedNode interface {
	Name() string
	Weight() int
	Children() map[string]*WeightedNode
}

type File struct {
	name string
	size int
}

func (f File) Name() string {
	return f.name
}

func (f File) Weight() int {
	return f.size
}

func (f File) Children() map[string]*WeightedNode {
	return nil
}

type Directory struct {
	name     string
	children map[string]*WeightedNode
}

func (d Directory) Name() string {
	return d.name
}

func (d Directory) Weight() int {
	weight := 0

	for _, child := range d.children {
		weight += (*child).Weight()
	}

	return weight
}

func (d Directory) Children() map[string]*WeightedNode {
	return d.children
}

func changeDirectories(target string, pathToRoot *xmas.Stack[*WeightedNode]) {
	current := pathToRoot.Peek()
	if target == ".." {
		pathToRoot.Pop()
	} else {
		pathToRoot.Push((*current).Children()[target])
	}
}

func buildFileTree(commandLog [][]string) WeightedNode {
	pathToRoot := xmas.Stack[*WeightedNode]{}

	var root WeightedNode = Directory{
		name:     "/",
		children: make(map[string]*WeightedNode),
	}
	pathToRoot.Push(&root)

	// first cd to root initialized above
	for _, logLine := range commandLog[1:] {
		switch logLine[1] {
		case "ls":
			// only ls in currentDir... safely skip this line
		case "cd":
			changeDirectories(logLine[2], &pathToRoot)
		default:
			filename := logLine[1]

			if _, ok := (*pathToRoot.Peek()).Children()[filename]; !ok {
				var newNode WeightedNode

				switch logLine[0] {
				case "dir":
					newNode = Directory{
						name:     filename,
						children: make(map[string]*WeightedNode),
					}
				default:
					filesize, _ := strconv.Atoi(logLine[0])
					newNode = File{
						name: filename,
						size: filesize,
					}
				}

				(*pathToRoot.Peek()).Children()[filename] = &newNode
			}
		}
	}

	return root
}

func PrintTree(root WeightedNode, indent int) {
	for i := 0; i < indent; i++ {
		fmt.Print(" ")
	}

	fmt.Print(root.Name())
	fmt.Println()

	for _, child := range root.Children() {
		PrintTree(*child, indent+1)
	}
}

func combinedSize(root WeightedNode, sum *int) {
	currentSize := root.Weight()
	if currentSize < 100000 {
		*sum += currentSize
	}

	for _, child := range root.Children() {
		switch (*child).(type) {
		case Directory:
			combinedSize(*child, sum)
		}
	}
}
func smallestDeleteableDirSize(root WeightedNode, minSatisfyingSize int) int {
	size := root.Weight()
	minSize := math.MaxInt
	if size > minSatisfyingSize {
		minSize = size
	}

	for _, child := range root.Children() {
		switch (*child).(type) {
		case Directory:
			size = smallestDeleteableDirSize(*child, minSatisfyingSize)

			if size > minSatisfyingSize && size < minSize {
				minSize = size
			}
		}
	}

	return minSize
}

func spaceToFree(root WeightedNode) int {
	usedSpace := root.Weight()
	freeSpace := hardDriveSize - usedSpace

	spaceNeeded := updateSize - freeSpace
	return spaceNeeded
}

// day7Cmd represents the day7 command
var day7Cmd = &cobra.Command{
	Use:   "day7 path/to/input/file",
	Short: "AoC Day 7",
	Long:  `Advent of Code Day 7`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		var runall bool = Part == "*"

		input, err := xmas.ReadFileAsTokenizedStringSlice(args[0])
		if err != nil {
			return err
		}
		fmt.Printf("%d terminal output lines read.\n", len(input))

		root := buildFileTree(input)

		if runall || Part == "1" {
			fmt.Println("Part 1 running...")

			sum := 0
			combinedSize(root, &sum)

			fmt.Printf("Combined size of directories smaller than 100,000 is **%d**.\n", sum)
		}

		if runall || Part == "2" {
			fmt.Println("Part 2 running...")

			fmt.Printf("The size of the smallest directory we could delete to install the update is **%d**.\n",
				smallestDeleteableDirSize(root, spaceToFree(root)))
		}

		xmas.PrintHolidayMessage(time.Since(start))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(day7Cmd)
}
