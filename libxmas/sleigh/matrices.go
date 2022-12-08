package sleigh

import "fmt"

func CreateMatrix[T any](numRows, numCols int) [][]T {
	out := make([][]T, numRows)
	for i := 0; i < numRows; i++ {
		out[i] = make([]T, numCols)
	}

	return out
}

func FlipMatrix[T any](in [][]T) ([][]T, error) {
	rowCount := len(in)
	if rowCount == 0 {
		return in, nil //got empty, return empty
	}

	rowLength := len(in[0])
	out := CreateMatrix[T](rowLength, rowCount)
	for rowNum, row := range in {
		if rowLength != len(row) {
			return nil, fmt.Errorf("dimensions are not standard for each input row. not a matrix")
		}

		for colNum, col := range row {
			out[colNum][rowNum] = col
		}
	}

	return out, nil
}
