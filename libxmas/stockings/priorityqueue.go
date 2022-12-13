package stockings

import (
	"fmt"
	"math"
)

type queueItem[T comparable] struct {
	item     T
	priority int
}

func NewMaxPriorityQueue[T comparable](capacity int, prioritize func(item T) int) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		data: NewBinaryHeap(capacity, func(a, b *queueItem[T]) bool {
			return a.priority > b.priority
		}),
		prioritize: prioritize,
		members:    make(map[T]*queueItem[T]),
	}
}

func NewMinPriorityQueue[T comparable](capacity int, prioritize func(item T) int) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		data: NewBinaryHeap(capacity, func(a, b *queueItem[T]) bool {
			//fmt.Println(a.item, b.item)
			return a.priority < b.priority
		}),
		prioritize: prioritize,
		members:    make(map[T]*queueItem[T]),
	}
}

type PriorityQueue[T comparable] struct {
	data       *BinaryHeap[*queueItem[T]]
	prioritize func(T) int
	members    map[T]*queueItem[T]
}

func (q *PriorityQueue[T]) GetNext() T {
	next := q.data.ExtractTop()

	delete(q.members, next.item)

	return next.item
}

func (q *PriorityQueue[T]) Has(item T) bool {
	_, ok := q.members[item]

	return ok
}

func (q *PriorityQueue[T]) IndexOf(key T) int {
	return q.data.IndexOf(q.members[key])
}

func (q *PriorityQueue[T]) Size() int {
	return q.data.size
}

func (q *PriorityQueue[T]) GetPriority(item T) int {
	if qItem, ok := q.members[item]; ok {
		return qItem.priority
	}

	return math.MaxInt
}

func (q *PriorityQueue[T]) TryIncreasePriority(item T) bool {
	if qItem, ok := q.members[item]; ok {
		newPriority := &queueItem[T]{item: item, priority: q.prioritize(item)}

		if q.data.sort(newPriority, qItem) {
			qItem.priority = newPriority.priority
			q.data.IncreaseKeyValue(qItem)

			return true
		}
	}

	panic("failed") //return false
}

func (q *PriorityQueue[T]) Print() {
	fromCurrent, printNext := 2, 0
	for i := 0; i < q.Size(); i++ {
		fmt.Print(q.data.data[i].priority, q.data.data[i].item, ",\t")
		if i == printNext {
			fmt.Println()
			printNext += fromCurrent
			fromCurrent *= 2
		}
	}
	fmt.Println()
}

func (q *PriorityQueue[T]) Add(item T) *PriorityQueue[T] {
	if _, ok := q.members[item]; !ok {
		qItem := &queueItem[T]{
			item:     item,
			priority: q.prioritize(item),
		}

		q.members[item] = qItem

		q.data.InsertKey(qItem)
		if q.data.size < 72 {
			fmt.Println("adding", *qItem)
		}
	} else {
		panic("didn't add to queue")
	}

	return q
}
