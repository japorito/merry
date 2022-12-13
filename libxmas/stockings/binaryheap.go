package stockings

import "fmt"

func NewBinaryHeap[T comparable](capacity int, sort func(a, b T) bool) *BinaryHeap[T] {
	return &BinaryHeap[T]{
		data:      make([]T, capacity),
		memberMap: make(map[T]bool),
		sort:      sort,
	}
}

type BinaryHeap[T comparable] struct {
	data      []T
	memberMap map[T]bool
	sort      func(a, b T) bool
	size      int
}

func (heap *BinaryHeap[T]) parent(i int) int {
	return (i - 1) / 2
}

func (heap *BinaryHeap[T]) left(i int) int {
	return i*2 + 1
}

func (heap *BinaryHeap[T]) right(i int) int {
	return i*2 + 2
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

func (heap *BinaryHeap[T]) reroot(root int) *BinaryHeap[T] {
	rootLeft, rootRight := heap.left(root), heap.right(root)
	newRoot := root

	if rootLeft < heap.size && heap.sort(heap.data[rootLeft], heap.data[newRoot]) {
		if root == 0 && heap.size < 72 {
			fmt.Println("here left")
		}
		newRoot = rootLeft
	}

	if rootRight < heap.size && heap.sort(heap.data[rootRight], heap.data[newRoot]) {
		if root == 0 && heap.size < 72 {
			fmt.Println("here right")
		}
		newRoot = rootLeft
	}

	if newRoot != root {
		if root == 0 && heap.size < 72 {
			fmt.Println(heap.data[root], heap.data[newRoot])
		}
		temp := heap.data[root]
		heap.data[root] = heap.data[newRoot]
		heap.data[newRoot] = temp
		heap.reroot(newRoot)
	}

	return heap
}

func (heap *BinaryHeap[T]) Has(key T) bool {
	return heap.memberMap[key]
}

func (heap *BinaryHeap[T]) InsertKey(key T) *BinaryHeap[T] {
	if heap.Has(key) {
		// No need to add key that already exists
		panic("didn't insert heap")
		//return heap
	}

	// ensure sufficient capacity
	heap.size++
	if heap.size == len(heap.data) {
		newData := make([]T, len(heap.data)*2)
		copy(newData, heap.data)
		heap.data = newData
	}

	// Put new element at bottom of heap
	i := heap.size - 1
	heap.data[i] = key
	heap.memberMap[key] = true

	var p int
	for p = heap.parent(i); i >= 0 && heap.sort(key, heap.data[p]); p = heap.parent(i) {
		//swap
		//fmt.Println(p, i)
		heap.data[i], heap.data[p] = heap.data[p], heap.data[i]
		i = p
	}
	//fmt.Println("after", p, i, heap.data[i], heap.data[p])

	return heap
}

func (heap *BinaryHeap[T]) ExtractTop() T {
	if heap.size <= 0 {
		panic("tried to take top of empty heap")
	} else if heap.size == 1 {
		d := heap.data[0]

		delete(heap.memberMap, d)
		heap.size = 0

		return d
	}

	d := heap.data[0]

	delete(heap.memberMap, d)
	heap.data[0] = heap.data[heap.size-1]
	heap.size--

	if heap.size < 72 {
		fmt.Println("extracting", d) //, heap.data[0:heap.size])
	}

	heap.reroot(0)

	if heap.size < 72 {
		fmt.Println("new root", heap.data[0]) //, heap.data[0:heap.size])
	}

	return d
}

func (heap *BinaryHeap[T]) IncreaseKeyValue(key T) *BinaryHeap[T] {
	if i := heap.IndexOf(key); i >= 0 {
		for p := heap.parent(i); i > 0 && heap.sort(key, heap.data[p]); p = heap.parent(i) {
			//swap
			heap.data[i], heap.data[p] = heap.data[p], heap.data[i]

			i = p
		}
	} else {
		panic("what happened?")
	}

	return heap
}
