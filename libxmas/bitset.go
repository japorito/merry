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
	data        []uint64
	capacity    int64
	leftMostIdx T
}

func (bs *BitSet[T]) blockNumFromIndex(idx T) int64 {
	return int64(idx-bs.leftMostIdx) / bitsetBlockSize
}

func (bs *BitSet[T]) blockOffsetFromIndex(idx T) int64 {
	if bs.capacity == 0 {
		return 0
	}

	return int64(idx-bs.leftMostIdx) % bitsetBlockSize
}

func (bs *BitSet[T]) indexFromBlockAndOffset(blockNum, offset int64) int64 {
	return (int64(bs.leftMostIdx) + (bitsetBlockSize * blockNum)) + offset
}

func (bs *BitSet[T]) rightMostIdx() T {
	return T(int64(bs.leftMostIdx) + (bitsetBlockSize * int64(bs.capacity)) - 1)
}

func (bs *BitSet[T]) assureCapacity(idx T) *BitSet[T] {
	if bs.capacity == 0 {
		//initialize set
		bs.data = make([]uint64, 1)
		bs.capacity = 1
		if idx < 0 {
			bs.leftMostIdx = (((idx + 1) / bitsetBlockSize) - 1) * bitsetBlockSize
		} else {
			bs.leftMostIdx = (idx / bitsetBlockSize) * bitsetBlockSize
		}

		return bs
	}

	rightMostIdx := bs.rightMostIdx()
	initialCapacity := bs.capacity

	copyBlockOffset := int64(0)
	if idx < bs.leftMostIdx {
		// If we're -1 left of leftMostIdx, -1 / blocksize = 0...
		// but we still need to add a block, so add 1 at end of calculation
		// accounting for zero-indexing
		// If we're -64 left of leftMostIdx, -64 / blocksize = -1...
		// need to add 1 to relative index to account for zero-indexing
		copyBlockOffset = ((int64(idx-bs.leftMostIdx) + 1) / (-1 * bitsetBlockSize)) + 1
		bs.capacity = bs.capacity + copyBlockOffset
		bs.leftMostIdx = bs.leftMostIdx - T(bitsetBlockSize*copyBlockOffset)
	} else if idx > rightMostIdx {
		bs.capacity = bs.blockNumFromIndex(idx) + 1
	}

	if bs.capacity != initialCapacity {
		temp := bs.data
		bs.data = make([]uint64, bs.capacity)

		copy(bs.data[copyBlockOffset:], temp)
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
	bs.assureCapacity(idx)

	block := bs.blockNumFromIndex(idx)
	bitmask := uint64(1 << bs.blockOffsetFromIndex(idx))

	bs.data[block] = bs.data[block] | bitmask

	return bs
}

func (bs *BitSet[T]) Off(idx T) *BitSet[T] {
	block := bs.blockNumFromIndex(idx)
	bitmask := uint64(1 << bs.blockOffsetFromIndex(idx))

	if block >= bs.capacity || block < 0 {
		// outside of capacity so implicitly off
		return bs
	}

	bitmask = ^bitmask
	bs.data[block] = bs.data[block] & bitmask

	return bs
}

func (bs *BitSet[T]) Has(idx T) bool {
	block := bs.blockNumFromIndex(idx)
	bitmask := uint64(1 << bs.blockOffsetFromIndex(idx))

	return block < bs.capacity &&
		block >= 0 &&
		(bs.data[block]&bitmask > 0)
}

func (bs *BitSet[T]) Intersect(others ...BitSet[T]) BitSet[T] {
	maxLeftMostIdx := bs.leftMostIdx
	minRightMostIdx := bs.rightMostIdx()

	for _, other := range others {
		if other.leftMostIdx > maxLeftMostIdx {
			maxLeftMostIdx = other.leftMostIdx
		}

		if minRightMostIdx < other.rightMostIdx() {
			minRightMostIdx = other.rightMostIdx()
		}

		if minRightMostIdx <= maxLeftMostIdx {
			return BitSet[T]{}
		}
	}

	capacityNeeded := int64((1 + minRightMostIdx - maxLeftMostIdx) / bitsetBlockSize)
	intersectedSet := BitSet[T]{
		data:        make([]uint64, capacityNeeded),
		capacity:    capacityNeeded,
		leftMostIdx: maxLeftMostIdx,
	}

	leftBlock := bs.blockNumFromIndex(maxLeftMostIdx)
	copy(intersectedSet.data, bs.data[leftBlock:leftBlock+intersectedSet.capacity])

	for _, other := range others {
		leftBlock = other.blockNumFromIndex(maxLeftMostIdx)
		for i, block := range other.data[leftBlock : leftBlock+intersectedSet.capacity] {
			intersectedSet.data[i] = intersectedSet.data[i] & block
		}
	}

	return intersectedSet
}

func (bs *BitSet[T]) Union(others ...BitSet[T]) BitSet[T] {
	minLeftIdx, maxRightIdx := bs.leftMostIdx, bs.rightMostIdx()

	for _, other := range others {
		if minLeftIdx > other.leftMostIdx {
			minLeftIdx = other.leftMostIdx
		}

		if maxRightIdx < other.rightMostIdx() {
			maxRightIdx = other.rightMostIdx()
		}
	}

	newCapacity := int64((1 + maxRightIdx - minLeftIdx) / bitsetBlockSize)
	unionSet := BitSet[T]{
		data:        make([]uint64, newCapacity),
		capacity:    newCapacity,
		leftMostIdx: minLeftIdx,
	}

	leftBlock := unionSet.blockNumFromIndex(bs.leftMostIdx)
	copy(unionSet.data[leftBlock:], bs.data)

	for _, other := range others {
		leftBlock = unionSet.blockNumFromIndex(other.leftMostIdx)
		for i, block := range other.data {
			unionSet.data[leftBlock+int64(i)] = unionSet.data[leftBlock+int64(i)] | block
		}
	}

	return unionSet
}

func (bs *BitSet[T]) Subtract(other BitSet[T]) *BitSet[T] {
	otherOffset := int64(bs.leftMostIdx-other.leftMostIdx) / bitsetBlockSize
	for blockNum := range bs.data {
		otherBlockNum := blockNum + int(otherOffset)
		if otherBlockNum >= 0 {
			if otherBlockNum >= int(other.capacity) {
				break
			}

			bs.data[blockNum] = bs.data[blockNum] & ^other.data[otherBlockNum]
		}
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
				members = append(members, T(bs.indexFromBlockAndOffset(int64(blockNum), i)))
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
