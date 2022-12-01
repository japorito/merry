package xmas

import (
	"bufio"
	"os"
)

func ReadFileToLines(filePath string) ([]string, error) {
	inputFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	inputScanner := bufio.NewScanner(inputFile)
	inputScanner.Split(bufio.ScanLines)

	var inputLines []string
	for inputScanner.Scan() {
		inputLines = append(inputLines, inputScanner.Text())
	}

	return inputLines, nil
}

func ReadFileAsInt64Slice(filePath string) ([]int64, error) {
	inputLines, err := ReadFileToLines(filePath)
	if err != nil {
		return nil, err
	}

	return ToInt64s(inputLines)
}

func ReadFileAsTokenizedStringSlice(filePath string) ([][]string, error) {
	inputLines, err := ReadFileToLines(filePath)
	if err != nil {
		return nil, err
	}

	return Tokenize(inputLines), nil
}

func ReadFileAsCharacterBooleanSlice(filePath string) ([][]bool, error) {
	inputLines, err := ReadFileToLines(filePath)
	if err != nil {
		return nil, err
	}

	return ToBools(inputLines)
}

func ReadFileToStringSliceBlocks(filePath string) ([][]string, error) {
	inputLines, err := ReadFileToLines(filePath)
	if err != nil {
		return nil, err
	}

	return BreakToBlocks(inputLines), nil
}

func ReadFileToInt64SliceBlocks(filePath string) ([][]int64, error) {
	inputLines, err := ReadFileToLines(filePath)
	if err != nil {
		return nil, err
	}

	return BreakToInt64Blocks(inputLines)
}
