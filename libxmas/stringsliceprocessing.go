package xmas

import (
	"fmt"
	"strconv"
	"strings"
)

func Tokenize(input []string) [][]string {
	var output [][]string
	for _, inputLine := range input {
		tokens := strings.Fields(inputLine)

		output = append(output, tokens)
	}

	return output
}

func ToInt64s(input []string) ([]int64, error) {
	var output []int64
	for lineIdx, inputLine := range input {
		parsedNum, err := strconv.ParseInt(inputLine, 10, 0)
		if err != nil {
			fmt.Printf("Encountered error parsing int at line: %d\n", lineIdx)

			return nil, err
		}

		output = append(output, parsedNum)
	}

	return output, nil
}

func ToBools(input []string) ([][]bool, error) {
	return CharToBools(input, '0', '1')
}

func CharToBools(input []string, trueVal, falseVal rune) ([][]bool, error) {
	var output [][]bool
	for _, inputLine := range input {
		var outputLine []bool

		for _, in := range inputLine {
			switch in {
			case trueVal:
				outputLine = append(outputLine, true)
			case falseVal:
				outputLine = append(outputLine, false)
			default:
				return nil, fmt.Errorf("unexpected character %v", in)
			}
		}

		output = append(output, outputLine)
	}

	return output, nil
}

func BreakToBlocks(input []string) [][]string {
	var output [][]string
	var outputBlock []string
	for _, inputLine := range input {
		inputLine = strings.TrimSpace(inputLine)

		if len(inputLine) > 0 {
			outputBlock = append(outputBlock, inputLine)
		} else if len(outputBlock) > 0 { // if anything has been read to this block
			output = append(output, outputBlock)

			outputBlock = make([]string, 0)
		} // else line is empty but nothing has yet been read. Continue to skip newlines until we read something
	}

	return output
}
