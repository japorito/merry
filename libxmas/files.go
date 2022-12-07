package xmas

import (
	"bufio"
	"fmt"
	"os"
)

func ReadFileToLines(filePath string) []string {
	inputFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Failed te read file", filePath)
		return nil
	}
	defer inputFile.Close()

	inputScanner := bufio.NewScanner(inputFile)
	inputScanner.Split(bufio.ScanLines)

	var inputLines []string
	for inputScanner.Scan() {
		inputLines = append(inputLines, inputScanner.Text())
	}

	return inputLines
}

func ReadFileToRuneSliceLines(filePath string) [][]rune {
	inputLines := ReadFileToLines(filePath)

	return ToRunes(inputLines)
}

func ReadFileAsInt64Slice(filePath string) []int64 {
	inputLines := ReadFileToLines(filePath)

	return ToInt64s(inputLines)
}

func ReadFileAsBinaryUInt64Slice(filePath string) []uint64 {
	inputLines := ReadFileToLines(filePath)

	return BinaryStringToUint64s(inputLines)
}

func ReadFileAsBitAbstractionUInt64Slice(filePath, zeroVal, oneVal string) []uint64 {
	inputLines := ReadFileToLines(filePath)

	return BitAbstractionToUint64s(inputLines, zeroVal, oneVal)
}

func ReadFileAsTokenizedStringSlice(filePath string) [][]string {
	inputLines := ReadFileToLines(filePath)

	return Tokenize(inputLines)
}

func ReadFileAsCharacterBooleanSlice(filePath string) [][]bool {
	inputLines := ReadFileToLines(filePath)

	return ToBools(inputLines)
}

func ReadFileToStringSliceBlocks(filePath string) [][]string {
	inputLines := ReadFileToLines(filePath)

	return BreakToBlocks(inputLines)
}

func ReadFileToInt64SliceBlocks(filePath string) [][]int64 {
	inputLines := ReadFileToLines(filePath)

	return BreakToInt64Blocks(inputLines)
}
