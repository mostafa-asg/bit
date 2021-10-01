package bit

const (
	// determines the minimum allocation bits on initialization or on grows
	// if you change this, you should also change the array item type which is uint64
	minBits = 64
)

type Set struct {
	arr []uint64
}

func NewSet(options ...Option) *Set {
	ops := &Options{
		nbits: 64,
	}

	for _, option := range options {
		option(ops)
	}

	return newSet(ops)
}

func newSet(opts *Options) *Set {
	return &Set{
		arr: make([]uint64, howManyUint64(opts.nbits)),
	}
}

// ValueOf returns a new bit set containing all the bits in the given long array.
func ValueOf(nums []uint64) *Set {
	set := NewSet(WithInitialBits(len(nums) * minBits))

	for i, num := range nums {
		set.arr[i] = num
	}

	return set
}

// locate find the index within the array
func locate(index int) (arrIndex int, bitIndex int) {
	arrIndex = index / minBits
	bitIndex = index - (arrIndex * minBits)
	return
}

// howManyUint64 returns how many uint64 is needed for storing N bits of data
func howManyUint64(nbits int) int {
	if nbits <= 0 {
		return 0
	}

	return (nbits-1)/minBits + 1
}

// Clear sets the bit specified by the index to false.
func (set *Set) Clear(index int) {
	arrIndex, bitIndex := locate(index)
	set.arr[arrIndex] = set.arr[arrIndex] & (^(1 << bitIndex))
}

// ClearRange sets the bits from the specified fromIndex (inclusive)
// to the specified toIndex (exclusive) to false.
func (set *Set) ClearRange(fromIndex int, toIndex int) {
	for i := fromIndex; i < toIndex; i++ {
		set.Clear(i)
	}
}

// Set sets the bit at the specified index to true.
func (set *Set) Set(index int) {
	arrIndex, bitIndex := locate(index)
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
	arrIndex, bitIndex := locate(index)
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
