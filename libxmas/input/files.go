// Package sleigh brings Santa's input data
package sleigh

import (
	"bufio"
	"fmt"

	xmas "github.com/japorito/merry/libxmas"
)

func ReadToLines(args ...string) []string {
	input, err := InputSource.ResolveInput(args...)
	if err != nil {
		fmt.Println("Failed to read input", args)
		return nil
	}
	defer input.Close()

	inputScanner := bufio.NewScanner(input)
	inputScanner.Split(bufio.ScanLines)

	var inputLines []string
	for inputScanner.Scan() {
		inputLines = append(inputLines, inputScanner.Text())
	}

	return inputLines
}

func ReadToRuneSliceLines(args ...string) [][]rune {
	inputLines := ReadToLines(args...)

	return xmas.ToRunes(inputLines)
}

func ReadAsInt64Slice(args ...string) []int64 {
	inputLines := ReadToLines(args...)

	return xmas.ToInt64s(inputLines)
}

func ReadAsBinaryUInt64Slice(args ...string) []uint64 {
	inputLines := ReadToLines(args...)

	return xmas.BinaryStringToUint64s(inputLines)
}

func ReadAsBitAbstractionUInt64Slice(zeroVal, oneVal string, args ...string) []uint64 {
	inputLines := ReadToLines(args...)

	return xmas.BitAbstractionToUint64s(inputLines, zeroVal, oneVal)
}

func ReadAsTokenizedStringSlice(args ...string) [][]string {
	inputLines := ReadToLines(args...)

	return xmas.Tokenize(inputLines)
}

func ReadAsCharacterBooleanSlice(args ...string) [][]bool {
	inputLines := ReadToLines(args...)

	return xmas.ToBools(inputLines)
}

func ReadToStringSliceBlocks(args ...string) [][]string {
	inputLines := ReadToLines(args...)

	return xmas.BreakToBlocks(inputLines)
}

func ReadToInt64SliceBlocks(args ...string) [][]int64 {
	inputLines := ReadToLines(args...)

	return xmas.BreakToInt64Blocks(inputLines)
}
