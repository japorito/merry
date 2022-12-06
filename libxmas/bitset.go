/*
This is a bad implementation that you should not use, but it is sufficient for my purposes.
Should work fine for small, positive ints.
*/

package xmas

import (
	"math/bits"
)

const bitsetBlockSize = 64

type Integral interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint |
		~int8 | ~int16 | ~int32 | ~int64 | ~int
}

type BitSet[T Integral] struct {
	data     []uint64
	capacity int64
}

func (bs *BitSet[T]) increaseCapacity(fieldsNeeded T) *BitSet[T] {
	bs.capacity = int64(fieldsNeeded / bitsetBlockSize)
	if fieldsNeeded%bitsetBlockSize > 0 {
		bs.capacity += 1
	}

	temp := bs.data
	bs.data = make([]uint64, bs.capacity)

	copy(bs.data, temp)

	return bs
}

func (bs *BitSet[T]) set(idx T, on bool) *BitSet[T] {
	block := int64(idx / bitsetBlockSize)
	bitmask := uint64(1 << (idx % bitsetBlockSize))

	if block >= bs.capacity {
		if !on {
			// No changes necessary. It's already "off"
			return bs
		}

		bs.increaseCapacity(idx)
	}

	if !on {
		bitmask = ^bitmask
		bs.data[block] = bs.data[block] & bitmask
	} else {
		bs.data[block] = bs.data[block] | bitmask
	}

	return bs
}

func (bs *BitSet[T]) Capacity() int64 {
	return bs.capacity * bitsetBlockSize
}

func (bs *BitSet[T]) Size() int64 {
	var size int64 = 0

	for _, block := range bs.data {
		size += int64(bits.OnesCount64(block))
	}

	return size
}

func (bs *BitSet[T]) On(idx T) *BitSet[T] {
	bs.set(idx, true)

	return bs
}

func (bs *BitSet[T]) Off(idx T) *BitSet[T] {
	bs.set(idx, false)

	return bs
}

func (bs *BitSet[T]) Has(idx T) bool {
	block := int64(idx / bitsetBlockSize)
	bitmask := uint64(1 << (idx % bitsetBlockSize))

	return block < bs.capacity && (bs.data[block]&bitmask > 0)
}

func (bs *BitSet[T]) Intersect(others ...BitSet[T]) BitSet[T] {
	mincapacity := bs.capacity

	for _, other := range others {
		if other.capacity < mincapacity {
			mincapacity = other.capacity
		}
	}

	intersectedSet := BitSet[T]{
		data:     make([]uint64, mincapacity),
		capacity: mincapacity,
	}

	copy(intersectedSet.data, bs.data[:mincapacity])

	for _, other := range others {
		for i, block := range other.data {
			if i >= int(mincapacity) {
				break
			}

			intersectedSet.data[i] = intersectedSet.data[i] & block
		}
	}

	return intersectedSet
}

func (bs *BitSet[T]) Union(others ...BitSet[T]) BitSet[T] {
	maxcapacity := bs.capacity

	for _, other := range others {
		if other.capacity > maxcapacity {
			maxcapacity = other.capacity
		}
	}

	unionSet := BitSet[T]{
		data:     make([]uint64, maxcapacity),
		capacity: maxcapacity,
	}

	copy(unionSet.data, bs.data[:maxcapacity])

	for _, other := range others {
		for i, block := range other.data {
			if i >= int(maxcapacity) {
				break
			}

			unionSet.data[i] = unionSet.data[i] | block
		}
	}

	return unionSet
}

func (bs *BitSet[T]) Subtract(other BitSet[T]) *BitSet[T] {
	for blockNum := range bs.data {
		if blockNum > int(other.capacity) {
			break
		}

		bs.data[blockNum] = bs.data[blockNum] & ^other.data[blockNum]
	}

	return bs
}

func (bs *BitSet[T]) Any() bool {
	for _, block := range bs.data {
		if block > 0 {
			return true
		}
	}

	return false
}

func (bs *BitSet[T]) Members() []T {
	members := make([]T, 0, bs.Size())
	for blockNum, block := range bs.data {
		for i := int64(0); i < bitsetBlockSize; i++ {
			if (block & (1 << i)) > 0 {
				members = append(members, T(int64(blockNum)*bitsetBlockSize+i))
			}
		}
	}

	return members
}

func (bs *BitSet[T]) Clear() *BitSet[T] {
	for i := range bs.data {
		bs.data[i] = uint64(0)
	}

	return bs
}
