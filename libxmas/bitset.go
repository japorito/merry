package xmas

import (
	"math/bits"
)

const bitsetBlockSize = 64

type BitSet struct {
	data     []uint64
	capacity int64
}

func (bs *BitSet) increaseCapacity(fieldsNeeded int64) *BitSet {
	bs.capacity = fieldsNeeded / bitsetBlockSize
	if fieldsNeeded%bitsetBlockSize > 0 {
		bs.capacity += 1
	}

	temp := bs.data
	bs.data = make([]uint64, bs.capacity)

	copy(bs.data, temp)

	return bs
}

func (bs *BitSet) set(idx int64, value bool) *BitSet {
	block := idx / bitsetBlockSize
	bitmask := uint64(1 << (idx % bitsetBlockSize))

	if block >= bs.capacity {
		bs.increaseCapacity(idx)
	}

	if !value {
		bitmask = ^bitmask
		bs.data[block] = bs.data[block] & bitmask
	} else {
		bs.data[block] = bs.data[block] | bitmask
	}

	return bs
}

func (bs *BitSet) Capacity() int64 {
	return bs.capacity * bitsetBlockSize
}

func (bs *BitSet) Size() int64 {
	var size int64 = 0

	for _, block := range bs.data {
		size += int64(bits.OnesCount64(block))
	}

	return size
}

func (bs *BitSet) On(idx int64) *BitSet {
	bs.set(idx, true)

	return bs
}

func (bs *BitSet) Off(idx int64) *BitSet {
	bs.set(idx, false)

	return bs
}

func (bs *BitSet) Has(idx int64) bool {
	block := idx / bitsetBlockSize
	bitmask := uint64(1 << (idx % bitsetBlockSize))

	return block < bs.capacity && (bs.data[block]&bitmask > 0)
}

func (bs *BitSet) Intersect(others ...BitSet) BitSet {
	mincapacity := bs.capacity

	for _, other := range others {
		if other.capacity < mincapacity {
			mincapacity = other.capacity
		}
	}

	intersectedSet := BitSet{
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

func (bs *BitSet) Union(others ...BitSet) BitSet {
	maxcapacity := bs.capacity

	for _, other := range others {
		if other.capacity > maxcapacity {
			maxcapacity = other.capacity
		}
	}

	unionSet := BitSet{
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

func (bs *BitSet) Subtract(other BitSet) *BitSet {
	for blockNum := range bs.data {
		if blockNum > int(other.capacity) {
			break
		}

		bs.data[blockNum] = bs.data[blockNum] & ^other.data[blockNum]
	}

	return bs
}

func (bs *BitSet) Any() bool {
	for _, block := range bs.data {
		if block > 0 {
			return true
		}
	}

	return false
}

func (bs *BitSet) Members() []int64 {
	members := make([]int64, 0, bs.Size())
	for blockNum, block := range bs.data {
		for i := int64(0); i < bitsetBlockSize; i++ {
			if (block & (1 << i)) > 0 {
				members = append(members, int64(blockNum)*bitsetBlockSize+i)
			}
		}
	}

	return members
}
