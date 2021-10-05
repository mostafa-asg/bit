package bit

import "errors"

const (
	// determines the minimum allocation bits on initialization or on grows
	// if you change this, you should also change the array item type which is uint64
	minBits = 64
)

type Set struct {
	arr []uint64
}

func NewSet(options ...Option) (*Set, error) {
	ops := &Options{
		nbits: 64,
	}

	for _, option := range options {
		option(ops)
	}

	return newSet(ops)
}

func newSet(opts *Options) (*Set, error) {
	if opts.nbits < 0 {
		return nil, errors.New("Number of bits is negative")
	}

	if opts.nbits == 0 {
		opts.nbits = minBits
	}

	return &Set{
		arr: make([]uint64, howManyUint64(opts.nbits)),
	}, nil
}

// ValueOf returns a new bit set containing all the bits in the given long array.
func ValueOf(nums []uint64) *Set {
	set, _ := NewSet(WithInitialBits(len(nums) * minBits))

	for i, num := range nums {
		set.arr[i] = num
	}

	return set
}

// Flip sets the bit at the specified index to the complement of its current value.
func (set *Set) Flip(index int) *Set {
	arrIndex, bitIndex := set.locate(index)
	set.arr[arrIndex] = set.arr[arrIndex] ^ (1 << bitIndex)
	return set
}

// FlipRange sets each bit from the specified fromIndex (inclusive)
// to the specified toIndex (exclusive) to the complement of its current value.
func (set *Set) FlipRange(fromIndex int, toIndex int) *Set {
	for i := fromIndex; i < toIndex; i++ {
		set.Flip(i)
	}
	return set
}

// Clear sets the bit specified by the index to false.
func (set *Set) Clear(index int) {
	arrIndex, bitIndex := set.locate(index)
	set.arr[arrIndex] = set.arr[arrIndex] & (^(1 << bitIndex))
}

// ClearRange sets the bits from the specified fromIndex (inclusive)
// to the specified toIndex (exclusive) to false.
func (set *Set) ClearRange(fromIndex int, toIndex int) {
	for i := fromIndex; i < toIndex; i++ {
		set.Clear(i)
	}
}

// ClearAll sets all of the bits in this BitSet to false.
func (set *Set) ClearAll() *Set {
	for i := range set.arr {
		set.arr[i] = 0
	}
	return set
}

// Set sets the bit at the specified index to true.
func (set *Set) Set(index int) {
	arrIndex, bitIndex := set.locate(index)
	set.arr[arrIndex] = set.arr[arrIndex] | (1 << bitIndex)
}

// SetRange sets the bits from the specified fromIndex (inclusive)
// to the specified toIndex (exclusive) to true.
func (set *Set) SetRange(fromIndex int, toIndex int) {
	for i := fromIndex; i < toIndex; i++ {
		set.Set(i)
	}
}

// SetValue sets the bit at the specified index to the specified value.
func (set *Set) SetValue(index int, value bool) {
	switch value {
	case true:
		set.Set(index)
	case false:
		set.Clear(index)
	}
}

// Get returns the value of the bit with the specified index.
func (set *Set) Get(index int) bool {
	arrIndex, bitIndex := set.locate(index)
	value := set.arr[arrIndex] & (1 << bitIndex)
	return value > 0
}

// Size returns the number of bits of space actually in use by this BitSet
// to represent bit values.
func (set *Set) Size() int {
	return len(set.arr) * minBits
}

// Cardinality returns the number of bits set to true in this BitSet.
func (set *Set) Cardinality() int {
	count := 0

	for _, item := range set.arr {
		n := item
		for n > 0 {
			n = n & (n - 1)
			count++
		}
	}

	return count
}

// And performs a logical AND of this target bit set with the argument bit set.
func (set *Set) And(otherSet *Set) *Set {
	length := min(len(set.arr), len(otherSet.arr))

	for i := 0; i < length; i++ {
		set.arr[i] &= otherSet.arr[i]
	}

	return set
}

// AndNot clears all of the bits in this BitSet whose corresponding bit is
// set in the specified BitSet.
func (set *Set) AndNot(otherSet *Set) *Set {
	tmp := set.Clone()
	tmp.And(otherSet)
	return set.Xor(tmp)
}

// Xor performs a logical XOR of this bit set with the bit set argument.
func (set *Set) Xor(otherSet *Set) *Set {
	length := min(len(set.arr), len(otherSet.arr))

	for i := 0; i < length; i++ {
		set.arr[i] ^= otherSet.arr[i]
	}

	return set
}

// Equal checks equality between this set and the other set passed in the argument.
func (set *Set) Equal(otherSet *Set) bool {
	length := min(len(set.arr), len(otherSet.arr))

	for i := 0; i < length; i++ {
		if set.arr[i] != otherSet.arr[i] {
			return false
		}
	}

	// if array's length is different, only they are equal if array's items equal to zero
	var arr []uint64
	if len(set.arr) > len(otherSet.arr) {
		arr = set.arr[length:]
	} else if len(otherSet.arr) > len(set.arr) {
		arr = otherSet.arr[length:]
	}

	if len(arr) > 0 {
		for i := 0; i < len(arr); i++ {
			if arr[i] != 0 {
				return false
			}
		}
	}

	return true
}

// Clone creates a new copy of the current set
func (set *Set) Clone() *Set {
	copySet, _ := NewSet(WithInitialBits(len(set.arr) * minBits))

	for i, item := range set.arr {
		copySet.arr[i] = item
	}

	return copySet
}

// Intersects returns true if the specified BitSet has any bits set to true that
// are also set to true in this BitSet.
func (set *Set) Intersects(otherSet *Set) bool {

	for _, item := range set.arr {
		for _, otherItem := range otherSet.arr {
			if item&otherItem > 0 {
				return true
			}
		}
	}

	return false
}

// IsEmpty returns true if this BitSet contains no bits that are set to true.
func (set *Set) IsEmpty() bool {
	for _, item := range set.arr {
		if item > 0 {
			return false
		}
	}

	return true
}

func (set *Set) expandIfNeeded(arrIndex int) {
	lastIndexNum := len(set.arr) - 1
	if arrIndex > lastIndexNum {
		itemsNeeded := arrIndex - lastIndexNum
		for i := 1; i <= itemsNeeded; i++ {
			set.arr = append(set.arr, uint64(0))
		}
	}
}

// locate find the index within the array, it also expands the array
// if index is out of range
func (set *Set) locate(index int) (arrIndex int, bitIndex int) {
	arrIndex = index / minBits
	bitIndex = index - (arrIndex * minBits)

	set.expandIfNeeded(arrIndex)

	return
}

// howManyUint64 returns how many uint64 is needed for storing N bits of data
func howManyUint64(nbits int) int {
	if nbits <= 0 {
		return 0
	}

	return (nbits-1)/minBits + 1
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
