package stockings

import (
	"container/heap"
)

func NewBinaryHeap[T comparable](capacity int, sort func(a, b T) bool) *BinaryHeap[T] {
	h := BinaryHeap[T]{
		data:      make([]T, 0, capacity),
		memberMap: make(map[T]bool),
		sort:      sort,
	}

	heap.Init(&h)

	return &h
}

type BinaryHeap[T comparable] struct {
	data      []T
	memberMap map[T]bool
	sort      func(a, b T) bool
}

func (heap *BinaryHeap[T]) Len() int {
	return len(heap.data)
}

func (heap *BinaryHeap[T]) Less(i, j int) bool {
	return heap.sort(heap.data[i], heap.data[j])
}

func (h *BinaryHeap[T]) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	h.memberMap[x.(T)] = true
	h.data = append(h.data, x.(T))

	heap.Fix(h, len(h.data)-1)
}

func (heap *BinaryHeap[T]) Pop() any {
	old := heap.data
	n := len(old)
	x := old[n-1]
	heap.data = old[0 : n-1]
	delete(heap.memberMap, x)
	return x
}

func (heap *BinaryHeap[T]) Swap(i, j int) {
	heap.data[i], heap.data[j] = heap.data[j], heap.data[i]
}

func (heap *BinaryHeap[T]) IndexOf(key T) int {
	if heap.memberMap[key] {
		for i := range heap.data {
			if heap.data[i] == key {
				return i
			}
		}
	}

	return -1
}

func (heap *BinaryHeap[T]) Has(key T) bool {
	return heap.memberMap[key]
}

func (h *BinaryHeap[T]) IncreaseKeyValue(key T) *BinaryHeap[T] {
	if i := h.IndexOf(key); i >= 0 {
		heap.Fix(h, i)
	}

	return h
}
