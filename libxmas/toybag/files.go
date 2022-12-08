// Package toybag is how Santa brings his input data
package toybag

import (
	"bufio"
	"fmt"

	"github.com/japorito/merry/libxmas/sleigh"
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

	return sleigh.ToRunes(inputLines)
}

func ReadAsInt64Slice(args ...string) []int64 {
	inputLines := ReadToLines(args...)

	return sleigh.ToInt64s(inputLines)
}

func ReadAsBinaryUInt64Slice(args ...string) []uint64 {
	inputLines := ReadToLines(args...)

	return sleigh.BinaryStringToUint64s(inputLines)
}

func ReadAsBitAbstractionUInt64Slice(zeroVal, oneVal string, args ...string) []uint64 {
	inputLines := ReadToLines(args...)

	return sleigh.BitAbstractionToUint64s(inputLines, zeroVal, oneVal)
}

func ReadAsTokenizedStringSlice(args ...string) [][]string {
	inputLines := ReadToLines(args...)

	return sleigh.Tokenize(inputLines)
}

func ReadAsCharacterBooleanSlice(args ...string) [][]bool {
	inputLines := ReadToLines(args...)

	return sleigh.ToBools(inputLines)
}

func ReadToStringSliceBlocks(args ...string) [][]string {
	inputLines := ReadToLines(args...)

	return sleigh.BreakToBlocks(inputLines)
}

func ReadToInt64SliceBlocks(args ...string) [][]int64 {
	inputLines := ReadToLines(args...)

	return sleigh.BreakToInt64Blocks(inputLines)
}
