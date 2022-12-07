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
	data      []uint64
	capacity  int64
	zeroBlock int64
}

func (bs *BitSet[T]) blockNumFromIndex(idx T) int64 {

	if idx < 0 {
		// -1 through -64 should return -1, shift +1 so floor division will work
		return bs.zeroBlock + (int64(idx+1) / bitsetBlockSize) - 1
	}

	// 0 through 63 should return 0
	return bs.zeroBlock + (int64(idx) / bitsetBlockSize)
}

func (bs *BitSet[T]) blockOffsetFromIndex(idx T) int64 {
	offset := int64(idx) % bitsetBlockSize
	if idx < 0 {
		return (bitsetBlockSize + offset) % bitsetBlockSize
	}

	return offset
}

func (bs *BitSet[T]) indexFromBlockAndOffset(blockNum, offset int64) int64 {
	effectiveBlock := blockNum - bs.zeroBlock
	minBlockVal := effectiveBlock * bitsetBlockSize

	return minBlockVal + offset
}

func (bs *BitSet[T]) assureCapacity(idx T) *BitSet[T] {
	blockNum := bs.blockNumFromIndex(idx)
	initialCapacity := bs.capacity

	copyBlockOffset := int64(0)
	if blockNum < 0 {
		copyBlockOffset = -1 * blockNum
		bs.capacity = bs.capacity + copyBlockOffset
		bs.zeroBlock = bs.zeroBlock + copyBlockOffset
	} else if blockNum >= bs.capacity {
		bs.capacity = blockNum + 1 // account for 0-indexing
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
	minNegativeBlocks := bs.zeroBlock
	minPositiveBlocks := bs.capacity - bs.zeroBlock

	for _, other := range others {
		otherPositiveBlocks := other.capacity - other.zeroBlock
		if otherPositiveBlocks < minPositiveBlocks {
			minPositiveBlocks = otherPositiveBlocks
		}

		if other.zeroBlock < minNegativeBlocks {
			minNegativeBlocks = other.zeroBlock
		}
	}

	newCapacity := minNegativeBlocks + minPositiveBlocks
	intersectedSet := BitSet[T]{
		data:      make([]uint64, newCapacity),
		capacity:  newCapacity,
		zeroBlock: minNegativeBlocks,
	}

	copy(intersectedSet.data, bs.data[bs.zeroBlock-minNegativeBlocks:bs.zeroBlock+minPositiveBlocks])

	for _, other := range others {
		for i, block := range other.data[other.zeroBlock-minNegativeBlocks : other.zeroBlock+minPositiveBlocks] {
			intersectedSet.data[i] = intersectedSet.data[i] & block
		}
	}

	return intersectedSet
}

func (bs *BitSet[T]) Union(others ...BitSet[T]) BitSet[T] {
	maxNegativeBlocks := bs.zeroBlock
	maxPositiveBlocks := bs.capacity - bs.zeroBlock

	for _, other := range others {
		otherPositiveBlocks := other.capacity - other.zeroBlock
		if otherPositiveBlocks > maxPositiveBlocks {
			maxPositiveBlocks = otherPositiveBlocks
		}

		if other.zeroBlock > maxNegativeBlocks {
			maxNegativeBlocks = other.zeroBlock
		}
	}

	newCapacity := maxNegativeBlocks + maxPositiveBlocks
	unionSet := BitSet[T]{
		data:      make([]uint64, newCapacity),
		capacity:  newCapacity,
		zeroBlock: maxNegativeBlocks,
	}

	copy(unionSet.data[unionSet.zeroBlock-bs.zeroBlock:], bs.data)

	for _, other := range others {
		otherZeroOffset := unionSet.zeroBlock - other.zeroBlock
		for i, block := range other.data {
			unionSet.data[otherZeroOffset+int64(i)] = unionSet.data[otherZeroOffset+int64(i)] | block
		}
	}

	return unionSet
}

func (bs *BitSet[T]) Subtract(other BitSet[T]) *BitSet[T] {
	otherOffset := other.zeroBlock - bs.zeroBlock
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
