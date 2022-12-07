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

func createFileListingNode(fileListing []string) *WeightedNode {
	var newNode WeightedNode
	filename := fileListing[1]
	switch fileListing[0] {
	case "dir":
		newNode = Directory{
			name:     filename,
			children: make(map[string]*WeightedNode),
		}
	default:
		filesize, _ := strconv.Atoi(fileListing[0])
		newNode = File{
			name: filename,
			size: filesize,
		}
	}

	return &newNode
}

func buildFileTree(commandLog [][]string) WeightedNode {
	pathToRoot := xmas.Stack[*WeightedNode]{}

	// initialize root directory
	var root WeightedNode = Directory{
		name:     "/",
		children: make(map[string]*WeightedNode),
	}
	pathToRoot.Push(&root)

	// first cd to root initialized above
	currentDir := pathToRoot.Peek()
	for _, logLine := range commandLog[1:] {
		switch logLine[1] {
		case "ls":
			// only ls in currentDir... safely skip this line
		case "cd":
			changeDirectories(logLine[2], &pathToRoot)
			currentDir = pathToRoot.Peek()
		default: // file or directory listing
			filename := logLine[1]
			if _, ok := (*currentDir).Children()[filename]; !ok {
				(*currentDir).Children()[filename] = createFileListingNode(logLine)
			}
		}
	}

	return root
}

// func printTree(node WeightedNode, indent int) {
// 	indentFormat := fmt.Sprintf("%%%ds", indent)
// 	indentString := fmt.Sprintf(indentFormat, "")

// 	fmt.Printf("%s - %s (%T %d)\n", indentString, node.Name(), node, node.Weight())

// 	for _, child := range node.Children() {
// 		printTree(*child, indent+1)
// 	}
// }

func walkTree(node WeightedNode, currentVal int, f func(WeightedNode, int) int) int {
	currentVal = f(node, currentVal)

	for _, child := range node.Children() {
		currentVal = walkTree(*child, currentVal, f)
	}

	return currentVal
}

func combinedSize(root WeightedNode, maxToInclude int) int {
	return walkTree(root, 0, func(node WeightedNode, sum int) int {
		switch node.(type) {
		case Directory:
			currentSize := node.Weight()
			if currentSize <= maxToInclude {
				sum += currentSize
			}
		}

		return sum
	})
}

func smallestDeleteableDirSize(root WeightedNode, minSatisfyingSize int) int {
	return walkTree(root, math.MaxInt, func(node WeightedNode, minSize int) int {
		switch node.(type) {
		case Directory:
			size := node.Weight()
			if size > minSatisfyingSize && size < minSize {
				minSize = size
			}
		}

		return minSize
	})
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
	Long:  `Advent of Code Day 7: No Space Left On Device`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		defer xmas.PrintHolidayMessage(time.Now())

		if input := xmas.ReadFileAsTokenizedStringSlice(args[0]); input != nil {
			fmt.Printf("%d terminal output lines read.\n", len(input))

			root := buildFileTree(input)

			if Parts.Has(1) {
				fmt.Println("Part 1 running...")

				fmt.Printf("Combined size of directories smaller than 100,000 is **%d**.\n", combinedSize(root, 100000))
			}

			if Parts.Has(2) {
				fmt.Println("Part 2 running...")

				fmt.Printf("The size of the smallest directory we could delete to install the update is **%d**.\n",
					smallestDeleteableDirSize(root, spaceToFree(root)))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(day7Cmd)
}
