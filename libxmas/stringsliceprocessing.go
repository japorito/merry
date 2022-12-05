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

func ToRunes(input []string) [][]rune {
	output := make([][]rune, len(input))
	for i, line := range input {
		output[i] = []rune(line)
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

func BinaryStringToUint64s(input []string) ([]uint64, error) {
	var output []uint64
	for lineIdx, inputLine := range input {
		parsedNum, err := strconv.ParseUint(inputLine, 2, 64)
		if err != nil {
			fmt.Printf("Encountered error parsing int at line: %d\n", lineIdx)

			return nil, err
		}

		output = append(output, parsedNum)
	}

	return output, nil
}

func BitAbstractionToUint64s(input []string, zeroVal, oneVal string) ([]uint64, error) {
	for idx := range input {
		input[idx] = strings.ReplaceAll(input[idx], zeroVal, "0")
		input[idx] = strings.ReplaceAll(input[idx], oneVal, "1")
	}

	return BinaryStringToUint64s(input)
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
		if len(strings.TrimSpace(inputLine)) > 0 {
			outputBlock = append(outputBlock, inputLine)
		} else if len(outputBlock) > 0 { // if anything has been read to this block
			output = append(output, outputBlock)

			outputBlock = make([]string, 0)
		} // else line is empty but nothing has yet been read. Continue to skip newlines until we read something
	}

	if len(outputBlock) > 0 {
		output = append(output, outputBlock)
	}

	return output
}

func BreakToInt64Blocks(input []string) ([][]int64, error) {
	stringBlocks := BreakToBlocks(input)

	var output [][]int64
	for _, block := range stringBlocks {
		iBlock, err := ToInt64s(block)
		if err != nil {
			return nil, err
		}

		output = append(output, iBlock)
	}

	return output, nil
}
