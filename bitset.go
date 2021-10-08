package bit

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
)

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

// FromByteArray returns a new bit set containing all the bits in the given byte array.
func FromByteArray(nums []byte) *Set {
	arr := make([]uint64, 0)
	size := 8

	if rem := len(nums) % size; rem != 0 {
		total := size - rem
		for i := 1; i <= total; i++ {
			nums = append(nums, 0) // LittleEndian.Uint64 expected 8 bytes
		}
	}

	for len(nums) > 0 {
		arr = append(arr, binary.LittleEndian.Uint64(nums[0:size]))
		nums = nums[size:]
	}

	return ValueOf(arr)
}

// Flip sets the bit at the specified index to the complement of its current value.
// If index is negative no change will happen.
func (set *Set) Flip(index int) *Set {
	if index < 0 {
		// do nothing
		return set
	}

	arrIndex, bitIndex := set.locate(index)
	set.arr[arrIndex] = set.arr[arrIndex] ^ (1 << bitIndex)
	return set
}

// FlipRange sets each bit from the specified fromIndex (inclusive)
// to the specified toIndex (exclusive) to the complement of its current value.
func (set *Set) FlipRange(fromIndex int, toIndex int) *Set {
	if fromIndex < 0 {
		fromIndex = 0
	}

	for i := fromIndex; i < toIndex; i++ {
		set.Flip(i)
	}
	return set
}

// Clear sets the bit specified by the index to false.
// If index is negative no change will happen.
func (set *Set) Clear(index int) {
	if index < 0 {
		// do nothing
		return
	}

	arrIndex, bitIndex := set.locate(index)
	set.arr[arrIndex] = set.arr[arrIndex] & (^(1 << bitIndex))
}

// ClearRange sets the bits from the specified fromIndex (inclusive)
// to the specified toIndex (exclusive) to false.
func (set *Set) ClearRange(fromIndex int, toIndex int) {
	if fromIndex < 0 {
		fromIndex = 0
	}

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
// If index is negative no change will happen.
func (set *Set) Set(index int) {
	if index < 0 {
		// do nothing
		return
	}

	arrIndex, bitIndex := set.locate(index)
	set.arr[arrIndex] = set.arr[arrIndex] | (1 << bitIndex)
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

// SetRange sets the bits from the specified fromIndex (inclusive)
// to the specified toIndex (exclusive) to true.
func (set *Set) SetRange(fromIndex int, toIndex int) {
	if fromIndex < 0 {
		// do nothing
		return
	}

	for i := fromIndex; i < toIndex; i++ {
		set.Set(i)
	}
}

// SetRangeValue sets the bits from the specified fromIndex (inclusive)
// to the specified toIndex (exclusive) to the specified value.
func (set *Set) SetRangeValue(fromIndex int, toIndex int, value bool) {
	if fromIndex < 0 {
		fromIndex = 0
	}

	for i := fromIndex; i < toIndex; i++ {
		set.SetValue(i, value)
	}
}

// Get returns the value of the bit with the specified index.
// If index is negative, always false will be returned
func (set *Set) Get(index int) bool {
	if index < 0 || index >= set.Size() {
		// outside boundary
		return false
	}

	arrIndex, bitIndex := set.locate(index)
	value := set.arr[arrIndex] & (1 << bitIndex)
	return value > 0
}

// GetRange returns a new BitSet composed of bits from this BitSet from fromIndex (inclusive)
// to toIndex (exclusive).
func (set *Set) GetRange(fromIndex int, toIndex int) *Set {
	result, _ := NewSet()
	resultIndex := 0

	if fromIndex < 0 {
		fromIndex = 0
	}

	for i := fromIndex; i < toIndex; i++ {
		result.SetValue(resultIndex, set.Get(i))
		resultIndex++
	}

	return result
}

// Size returns the number of bits of space actually in use by this BitSet
// to represent bit values.
func (set *Set) Size() int {
	return len(set.arr) * minBits
}

// Length returns the "logical size" of this BitSet: the index of the
// highest set bit in the BitSet plus one.
// Returns zero if the BitSet contains no set bits.
func (set *Set) Length() int {
	index, _ := set.PreviousSetBit(set.Size() - 1)
	if index == -1 {
		return 0
	}
	return index + 1
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
// This bit set is modified so that each bit in it has the value true if and only if
// it both initially had the value true and the corresponding bit in the bit set
// argument also had the value true.
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

// Or performs a logical OR of this bit set with the bit set argument.
// This bit set is modified so that a bit in it has the value true if and only if
// it either already had the value true or the corresponding bit in the
// bit set argument has the value true.
func (set *Set) Or(otherSet *Set) *Set {
	length := min(len(set.arr), len(otherSet.arr))

	for i := 0; i < length; i++ {
		set.arr[i] |= otherSet.arr[i]
	}

	return set
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

// ToArray returns a new array containing all the bits in this bit set.
func (set *Set) ToArray() []uint64 {
	result := make([]uint64, len(set.arr))

	copy(result, set.arr)

	return result
}

// Bytes returns a new byte array containing all the bits in this bit set.
func (set *Set) Bytes() []byte {
	buf := new(bytes.Buffer)

	for _, item := range set.arr {
		binary.Write(buf, binary.LittleEndian, item)
	}

	return buf.Bytes()
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

// NextClearBit returns the index of the first bit that is set to false that
// occurs on or after the specified starting index.
func (set *Set) NextClearBit(fromIndex int) (int, error) {
	return set.nextBitIndex(fromIndex, false)
}

// NextSetBit returns the index of the first bit that is set to true that
// occurs on or after the specified starting index.
func (set *Set) NextSetBit(fromIndex int) (int, error) {
	return set.nextBitIndex(fromIndex, true)
}

// PreviousClearBit returns the index of the nearest bit that is set to false that
// occurs on or before the specified starting index.
// If no such bit exists, or if -1 is given as the starting index, then -1 is returned.
func (set *Set) PreviousClearBit(fromIndex int) (int, error) {
	return set.previousBitIndex(fromIndex, false)
}

// PreviousSetBit returns the index of the nearest bit that is set to true that
// occurs on or before the specified starting index.
// If no such bit exists, or if -1 is given as the starting index, then -1 is returned.
func (set *Set) PreviousSetBit(fromIndex int) (int, error) {
	return set.previousBitIndex(fromIndex, true)
}

// String returns a string representation of this bit set. For every index
// for which this BitSet contains a bit in the set state, the decimal
// representation of that index is included in the result. Such indices are
// listed in order from lowest to highest, separated by ", " (a comma and a space) and
// surrounded by braces, resulting in the usual mathematical notation for a
// set of integers.
func (set *Set) String() string {
	b := bytes.Buffer{}

	index := 0
	for {
		i, _ := set.NextSetBit(index)
		if i == -1 {
			break
		}
		b.WriteString(", ")
		b.WriteString(strconv.Itoa(i))

		index = i + 1
	}

	content := b.String()
	if len(content) > 0 {
		content = content[2:] // remove first ", "
	}

	return "{" + content + "}"
}

func (set *Set) previousBitIndex(fromIndex int, value bool) (int, error) {
	if fromIndex < -1 {
		return -1, fmt.Errorf("Index is negative: %d", fromIndex)
	}

	if fromIndex == -1 {
		return -1, nil
	}

	lastIndex := len(set.arr)*minBits - 1

	// outside boundery check
	if fromIndex > lastIndex {
		if value {
			// we know all bits are clear, no need to search
			fromIndex = lastIndex
		} else {
			// all is clear outside boundary
			return fromIndex, nil
		}
	}

	for i := fromIndex; i >= 0; i-- {
		if set.Get(i) == value {
			return i, nil
		}
	}

	return -1, nil
}

func (set *Set) nextBitIndex(fromIndex int, value bool) (int, error) {
	if fromIndex < 0 {
		return -1, fmt.Errorf("Index should be positive: %d", fromIndex)
	}

	lastIndex := len(set.arr)*minBits - 1

	// outside boundery check
	if fromIndex > lastIndex {
		if value {
			return -1, nil // there is no set bit outside boundary
		}

		// all is clear outside boundary
		return fromIndex, nil
	}

	for i := fromIndex; i <= lastIndex; i++ {
		if set.Get(i) == value {
			return i, nil
		}
	}

	if value {
		// no set bit found
		return -1, nil
	}

	// value equals to false
	return lastIndex + 1, nil
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
