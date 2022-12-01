package xmas

import "sort"

func SortInt64Asc(input []int64) {
	sort.Slice(input, func(i, j int) bool { return input[i] < input[j] })
}

func SortInt64Desc(input []int64) {
	sort.Slice(input, func(i, j int) bool { return input[i] > input[j] })
}

func SumInt64(input []int64) int64 {
	sum := int64(0)
	for _, addend := range input {
		sum = sum + addend
	}

	return sum
}
