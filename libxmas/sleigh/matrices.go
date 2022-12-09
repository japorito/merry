package sleigh

import "fmt"

func Reverse[T any](in *[]T) {
	for i, j := 0, len(*in)-1; i < j; i, j = i+1, j-1 {
		(*in)[i], (*in)[j] = (*in)[j], (*in)[i]
	}
}

func ReverseVertical[T any](in *[][]T) {
	Reverse(in)
}

func ReverseHorizontal[T any](in *[][]T) {
	for i := range *in {
		Reverse(&(*in)[i])
	}
}

func CreateMatrix[T any](numRows, numCols int) [][]T {
	out := make([][]T, numRows)
	for i := 0; i < numRows; i++ {
		out[i] = make([]T, numCols)
	}

	return out
}

func StandardizeDimensions[T any](in [][]T, fill T) [][]T {
	if len(in) == 0 {
		return in // got empty, return empty
	}

	maxWidth := 0
	for _, row := range in {
		rowWidth := len(row)
		if rowWidth > maxWidth {
			maxWidth = rowWidth
		}
	}

	out := make([][]T, len(in))
	for rowNum, row := range in {
		rowLength := len(row)
		if rowLength < maxWidth {
			out[rowNum] = make([]T, maxWidth)
			copy(out[rowNum], row)

			for i := rowLength; i < maxWidth; i++ {
				out[rowNum][i] = fill
			}
		} else {
			out[rowNum] = row
		}
	}

	return out
}

func TransposeMatrix[T any](in [][]T) ([][]T, error) {
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
