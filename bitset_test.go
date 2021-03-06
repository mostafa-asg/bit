package bit

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumberOfUint64Needed(t *testing.T) {
	testCases := []struct {
		nbits    int
		expected int
	}{
		{
			nbits:    1,
			expected: 1,
		},
		{
			nbits:    64,
			expected: 1,
		},
		{
			nbits:    65,
			expected: 2,
		},
		{
			nbits:    0,
			expected: 0,
		},
		{
			nbits:    -5,
			expected: 0,
		},
	}

	for _, test := range testCases {
		result := howManyUint64(test.nbits)
		assert.Equal(t, test.expected, result)
	}
}

func TestSet(t *testing.T) {
	nbits := 65
	s, err := NewSet(WithInitialBits(nbits))
	if err != nil {
		t.FailNow()
	}

	trueIndexes := make(map[int]bool)

	// all bits are initially false
	checkBits(t, s, trueIndexes)

	// ------------------------
	// Set
	// ------------------------
	s.Set(12)
	trueIndexes[12] = true
	checkBits(t, s, trueIndexes)

	s.Set(29)
	trueIndexes[29] = true
	checkBits(t, s, trueIndexes)

	s.Set(60)
	trueIndexes[60] = true
	checkBits(t, s, trueIndexes)

	s.Set(64)
	trueIndexes[64] = true
	checkBits(t, s, trueIndexes)

	// expansion test
	s.Set(3 * minBits)
	trueIndexes[3*minBits] = true
	checkBits(t, s, trueIndexes)

	// should have no side effect
	s.Clear(-1)
	s.Set(-1)
	checkBits(t, s, trueIndexes)
	// ------------------------
	// Clear
	// ------------------------
	s.Clear(12)
	delete(trueIndexes, 12)
	checkBits(t, s, trueIndexes)

	s.Clear(29)
	delete(trueIndexes, 29)
	checkBits(t, s, trueIndexes)

	s.Clear(60)
	delete(trueIndexes, 60)
	checkBits(t, s, trueIndexes)

	s.Clear(64)
	delete(trueIndexes, 64)
	checkBits(t, s, trueIndexes)

	s.Clear(3 * minBits)
	delete(trueIndexes, 3*minBits)
	checkBits(t, s, trueIndexes)
	//-------------------------
	// SetRange
	// ------------------------
	s.SetRange(50, 4*minBits)
	for i := 50; i < 4*minBits; i++ {
		trueIndexes[i] = true
	}
	checkBits(t, s, trueIndexes)

	// ------------------------
	// ClearRange
	// ------------------------
	s.ClearRange(51, 4*minBits)
	for i := 51; i < 4*minBits; i++ {
		delete(trueIndexes, i)
	}
	checkBits(t, s, trueIndexes)

	assert.True(t, s.Get(50))

	//-------------------------
	// SetRangeValue(true)
	// ------------------------
	s.SetRangeValue(0, 4, true)
	for i := 0; i < 4; i++ {
		trueIndexes[i] = true
	}
	checkBits(t, s, trueIndexes)

	// SetRangeValue(false)
	s.SetRangeValue(0, 4, false)
	for i := 0; i < 4; i++ {
		delete(trueIndexes, i)
	}
	checkBits(t, s, trueIndexes)
}

func TestGet(t *testing.T) {
	set, _ := NewSet()
	set.Set(0)
	set.Set(1)
	assert.True(t, set.Get(0))
	assert.True(t, set.Get(1))

	// outside boundary
	assert.False(t, set.Get(-1))
	assert.False(t, set.Get(set.Size()))
}

func TestCardinality(t *testing.T) {
	testCases := []struct {
		set         *Set
		cardinality int
	}{
		{
			set:         ValueOf([]uint64{15}),
			cardinality: 4,
		},
		{
			set:         ValueOf([]uint64{1, 1}),
			cardinality: 2,
		},
		{
			set:         ValueOf([]uint64{^uint64(0), ^uint64(0)}),
			cardinality: 128,
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.cardinality, test.set.Cardinality())
	}
}

func TestAnd(t *testing.T) {
	testCases := []struct {
		set1     *Set
		set2     *Set
		expected *Set
	}{
		{
			set1:     ValueOf([]uint64{15}),
			set2:     ValueOf([]uint64{10}),
			expected: ValueOf([]uint64{10}),
		},
		{
			set1:     ValueOf([]uint64{15, 32}),
			set2:     ValueOf([]uint64{10}),
			expected: ValueOf([]uint64{10, 32}),
		},
		{
			set1:     ValueOf([]uint64{15}),
			set2:     ValueOf([]uint64{10, 32}),
			expected: ValueOf([]uint64{10}),
		},
	}

	for _, test := range testCases {
		result := test.set1.And(test.set2)
		assert.True(t, test.expected.Equal(result))
	}
}

func TestXor(t *testing.T) {
	testCases := []struct {
		set1     *Set
		set2     *Set
		expected *Set
	}{
		{
			set1:     ValueOf([]uint64{15}),
			set2:     ValueOf([]uint64{10}),
			expected: ValueOf([]uint64{5}),
		},
		{
			set1:     ValueOf([]uint64{15, 32}),
			set2:     ValueOf([]uint64{10}),
			expected: ValueOf([]uint64{5, 32}),
		},
		{
			set1:     ValueOf([]uint64{15}),
			set2:     ValueOf([]uint64{10, 32}),
			expected: ValueOf([]uint64{5}),
		},
	}

	for _, test := range testCases {
		result := test.set1.Xor(test.set2)
		assert.True(t, test.expected.Equal(result))
	}
}

func TestAndNot(t *testing.T) {
	testCases := []struct {
		set1     *Set
		set2     *Set
		expected *Set
	}{
		{
			set1:     ValueOf([]uint64{7}),
			set2:     ValueOf([]uint64{13}),
			expected: ValueOf([]uint64{2}),
		},
		{
			set1:     ValueOf([]uint64{0}),
			set2:     ValueOf([]uint64{15}),
			expected: ValueOf([]uint64{0}),
		},
		{
			set1:     ValueOf([]uint64{15}),
			set2:     ValueOf([]uint64{0}),
			expected: ValueOf([]uint64{15}),
		},
		{
			set1:     ValueOf([]uint64{15}),
			set2:     ValueOf([]uint64{15}),
			expected: ValueOf([]uint64{0}),
		},
	}

	for _, test := range testCases {
		result := test.set1.AndNot(test.set2)
		assert.True(t, test.expected.Equal(result))
	}
}
func TestEqual(t *testing.T) {
	testCases := []struct {
		set1     *Set
		set2     *Set
		expected bool
	}{
		{
			set1:     ValueOf([]uint64{83, 12}),
			set2:     ValueOf([]uint64{83, 12}),
			expected: true,
		},
		{
			set1:     ValueOf([]uint64{83, 12}),
			set2:     ValueOf([]uint64{83, 11}),
			expected: false,
		},
		{
			set1:     ValueOf([]uint64{11, 20, 0, 0}),
			set2:     ValueOf([]uint64{11, 20}),
			expected: true,
		},
		{
			set1:     ValueOf([]uint64{11, 20, 0, 0, 1}),
			set2:     ValueOf([]uint64{11, 20}),
			expected: false,
		},
		{
			set1:     ValueOf([]uint64{66, 90}),
			set2:     ValueOf([]uint64{66, 90, 0, 0}),
			expected: true,
		},
		{
			set1:     ValueOf([]uint64{66, 90}),
			set2:     ValueOf([]uint64{66, 90, 0, 0, 2}),
			expected: false,
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.expected, test.set1.Equal(test.set2))
		assert.Equal(t, test.expected, test.set2.Equal(test.set1))
	}
}

func checkBits(t *testing.T, s *Set, trueIndexes map[int]bool) {
	for i := 0; i < s.Size(); i++ {
		_, mustTrue := trueIndexes[i]
		if mustTrue {
			assert.True(t, s.Get(i), fmt.Sprintf("index %d must be true", i))
		} else {
			assert.False(t, s.Get(i))
		}
	}
}

func TestClearAll(t *testing.T) {
	testCases := []struct {
		set      *Set
		expected *Set
	}{
		{
			set:      ValueOf([]uint64{6361}),
			expected: ValueOf([]uint64{0}),
		},
		{
			set:      ValueOf([]uint64{83, 12}),
			expected: ValueOf([]uint64{0, 0}),
		},
	}

	for _, test := range testCases {
		assert.True(t, test.expected.Equal(test.set.ClearAll()))
	}
}

func TestFlip(t *testing.T) {
	s := ValueOf([]uint64{14})

	// no changes will apply
	s.Flip(-4)
	assert.True(t, ValueOf([]uint64{14}).Equal(s))

	s.Flip(0)
	assert.True(t, ValueOf([]uint64{15}).Equal(s))

	s.Flip(0)
	assert.True(t, ValueOf([]uint64{14}).Equal(s))

	s.Flip(3)
	assert.True(t, ValueOf([]uint64{6}).Equal(s))

	s.Flip(3)
	assert.True(t, ValueOf([]uint64{14}).Equal(s))

	// expansion test
	s.Flip(72)
	assert.True(t, ValueOf([]uint64{14, 256}).Equal(s))
}

func TestFlipRange(t *testing.T) {
	s, err := NewSet()
	if err != nil {
		t.Fail()
	}
	// also expansion test
	s.FlipRange(63, 65)
	expected := ValueOf([]uint64{1 << 63, 1})
	assert.True(t, expected.Equal(s))

	s.FlipRange(63, 65)
	expected = ValueOf([]uint64{0, 0})
	assert.True(t, expected.Equal(s))
}

func TestArrayExpansion(t *testing.T) {
	s, err := NewSet()
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 1, len(s.arr))

	// no need to expansion
	s.expandIfNeeded(0)
	assert.Equal(t, 1, len(s.arr))

	// request access to index 3
	s.expandIfNeeded(3)
	assert.Equal(t, 4, len(s.arr))

	// no need to expansion
	s.expandIfNeeded(2)
	assert.Equal(t, 4, len(s.arr))
}

func TestIntersects(t *testing.T) {
	testCases := []struct {
		set1     *Set
		set2     *Set
		expected bool
	}{
		{
			set1:     ValueOf([]uint64{8}),
			set2:     ValueOf([]uint64{2}),
			expected: false,
		},
		{
			set1:     ValueOf([]uint64{8, 1, 4}),
			set2:     ValueOf([]uint64{2, 0, 12}),
			expected: true,
		},
		{
			set1:     ValueOf([]uint64{0, 0, 10}),
			set2:     ValueOf([]uint64{8}),
			expected: true,
		},
		{
			set1:     ValueOf([]uint64{10}),
			set2:     ValueOf([]uint64{0, 0, 0, 0, 8}),
			expected: true,
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.expected, test.set1.Intersects(test.set2))
	}
}

func TestIsEmpty(t *testing.T) {
	testCases := []struct {
		set1     *Set
		expected bool
	}{
		{
			set1:     ValueOf([]uint64{0}),
			expected: true,
		},
		{
			set1:     ValueOf([]uint64{0, 0, 0, 0}),
			expected: true,
		},
		{
			set1:     ValueOf([]uint64{0, 0, 10, 0, 0}),
			expected: false,
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.expected, test.set1.IsEmpty())
	}

	// tests creation via NewSet
	s, err := NewSet(WithInitialBits(3 * 64))
	if err != nil {
		t.Fail()
	}
	assert.True(t, s.IsEmpty())
	s.Set(4 * 64)
	assert.False(t, s.IsEmpty())
	s.Clear(4 * 64)
	assert.True(t, s.IsEmpty())
}

func TestNextClearBit(t *testing.T) {
	testCases := []struct {
		set1        *Set
		fromIndex   int
		expected    int
		expectError bool
	}{
		{
			set1:     ValueOf([]uint64{8}),
			expected: 0,
		},
		{
			set1:     ValueOf([]uint64{1}),
			expected: 1,
		},
		{
			set1:     ValueOf([]uint64{15}),
			expected: 4,
		},
		{
			set1:     ValueOf([]uint64{^uint64(0), ^uint64(0)}),
			expected: 128,
		},
		{
			set1:      ValueOf([]uint64{10}),
			fromIndex: 1,
			expected:  2,
		},
		{
			// outside boundary check
			set1:        ValueOf([]uint64{10}),
			fromIndex:   -1,
			expectError: true,
		},
		{
			// outside boundary check
			set1:      ValueOf([]uint64{15}),
			fromIndex: 70,
			expected:  70,
		},
	}

	for _, test := range testCases {
		index, err := test.set1.NextClearBit(test.fromIndex)
		if err != nil {
			if test.expectError {
				// everything is ok, continue
				continue
			}

			t.FailNow()
		}

		assert.Equal(t, test.expected, index)
	}
}

func TestNextSetBit(t *testing.T) {
	testCases := []struct {
		set1        *Set
		fromIndex   int
		expected    int
		expectError bool
	}{
		{
			set1:     ValueOf([]uint64{8}),
			expected: 3,
		},
		{
			set1:     ValueOf([]uint64{1}),
			expected: 0,
		},
		{
			set1:     ValueOf([]uint64{0, 1}),
			expected: 64,
		},
		{
			set1:      ValueOf([]uint64{11}),
			fromIndex: 2,
			expected:  3,
		},
		{
			set1:      ValueOf([]uint64{0}),
			fromIndex: 0,
			expected:  -1,
		},
		{
			// outside boundary check
			set1:        ValueOf([]uint64{10}),
			fromIndex:   -1,
			expectError: true,
		},
		{
			// outside boundary check
			set1:      ValueOf([]uint64{15}),
			fromIndex: 70,
			expected:  -1,
		},
	}

	for _, test := range testCases {
		index, err := test.set1.NextSetBit(test.fromIndex)
		if err != nil {
			if test.expectError {
				// everything is ok, continue
				continue
			}

			t.FailNow()
		}

		assert.Equal(t, test.expected, index)
	}
}

func TestPreviousClearBit(t *testing.T) {
	testCases := []struct {
		set1        *Set
		fromIndex   int
		expected    int
		expectError bool
	}{
		{
			set1:      ValueOf([]uint64{8}),
			fromIndex: 3,
			expected:  2,
		},
		{
			set1:     ValueOf([]uint64{1}),
			expected: -1,
		},
		{
			set1:      ValueOf([]uint64{0, 1}),
			fromIndex: 64,
			expected:  63,
		},
		{
			set1:      ValueOf([]uint64{^uint64(0)}),
			fromIndex: 63,
			expected:  -1,
		},
		{
			// outside boundary check
			set1:        ValueOf([]uint64{10}),
			fromIndex:   -1,
			expected:    -1,
			expectError: false,
		},
		{
			// outside boundary check
			set1:        ValueOf([]uint64{15}),
			fromIndex:   -2,
			expectError: true,
		},
		{
			// outside boundary check
			set1:      ValueOf([]uint64{15}),
			fromIndex: 70,
			expected:  70,
		},
	}

	for _, test := range testCases {
		index, err := test.set1.PreviousClearBit(test.fromIndex)
		if err != nil {
			if test.expectError {
				// everything is ok, continue
				continue
			}

			t.FailNow()
		}

		assert.Equal(t, test.expected, index)
	}
}

func TestPreviousSetBit(t *testing.T) {
	testCases := []struct {
		set1        *Set
		fromIndex   int
		expected    int
		expectError bool
	}{
		{
			set1:      ValueOf([]uint64{0}),
			fromIndex: 63,
			expected:  -1,
		},
		{
			set1:      ValueOf([]uint64{8}),
			fromIndex: 3,
			expected:  3,
		},
		{
			set1:      ValueOf([]uint64{9}),
			fromIndex: 2,
			expected:  0,
		},
		{
			// outside boundary check
			set1:      ValueOf([]uint64{0, 1, 0}),
			fromIndex: 500,
			expected:  64,
		},
		{
			// outside boundary check
			set1:        ValueOf([]uint64{10}),
			fromIndex:   -1,
			expected:    -1,
			expectError: false,
		},
		{
			// outside boundary check
			set1:        ValueOf([]uint64{15}),
			fromIndex:   -2,
			expectError: true,
		},
	}

	for _, test := range testCases {
		index, err := test.set1.PreviousSetBit(test.fromIndex)
		if err != nil {
			if test.expectError {
				// everything is ok, continue
				continue
			}

			t.FailNow()
		}

		assert.Equal(t, test.expected, index)
	}
}

func TestSize(t *testing.T) {
	testCases := []struct {
		set      *Set
		expected int
	}{
		{
			set:      ValueOf([]uint64{8}),
			expected: minBits,
		},
		{
			set:      ValueOf([]uint64{8, 0, 1}),
			expected: 3 * minBits,
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.expected, test.set.Size())
	}

	set, err := NewSet()
	if err != nil {
		t.FailNow()
	}
	assert.Equal(t, minBits, set.Size())

	set, err = NewSet(WithInitialBits(80))
	if err != nil {
		t.FailNow()
	}
	assert.Equal(t, 2*minBits, set.Size())
}

func TestToArray(t *testing.T) {
	testCases := []struct {
		set      *Set
		expected []uint64
	}{
		{
			set:      ValueOf([]uint64{8, 7, 12}),
			expected: []uint64{8, 7, 12},
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.expected, test.set.ToArray())
	}

	set, err := NewSet()
	if err != nil {
		t.FailNow()
	}
	set.Set(4)
	assert.Equal(t, []uint64{16}, set.ToArray())
}

func TestBytes(t *testing.T) {
	testCases := []struct {
		set      *Set
		expected []byte
	}{
		{
			set:      ValueOf([]uint64{0}),
			expected: []byte{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			set:      ValueOf([]uint64{8, 10}),
			expected: []byte{8, 0, 0, 0, 0, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0},
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.expected, test.set.Bytes())
	}
}

func TestFromByteArray(t *testing.T) {
	testCases := []struct {
		arr      []byte
		expected *Set
	}{
		{
			arr:      []byte{0},
			expected: ValueOf([]uint64{0}),
		},
		{
			arr:      []byte{8, 0, 0, 0},
			expected: ValueOf([]uint64{8}),
		},
		{
			arr:      []byte{255, 0, 0, 0, 0, 0, 0, 0},
			expected: ValueOf([]uint64{255}),
		},
		{
			arr:      []byte{255, 0, 0, 0, 0, 0, 0, 0, 7},
			expected: ValueOf([]uint64{255, 7}),
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.expected, FromByteArray(test.arr))
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		set      *Set
		expected string
	}{
		{
			set:      ValueOf([]uint64{0}),
			expected: "{}",
		},
		{
			set:      ValueOf([]uint64{8, 7}),
			expected: "{3, 64, 65, 66}",
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.expected, test.set.String())
	}
}

func TestLength(t *testing.T) {
	testCases := []struct {
		set      *Set
		expected int
	}{
		{
			set:      ValueOf([]uint64{0}),
			expected: 0,
		},
		{
			set:      ValueOf([]uint64{9}),
			expected: 3 + 1,
		},
		{
			set:      ValueOf([]uint64{9, 2}),
			expected: 65 + 1,
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.expected, test.set.Length())
	}
}

func TestOr(t *testing.T) {
	testCases := []struct {
		set1     *Set
		set2     *Set
		expected *Set
	}{
		{
			set1:     ValueOf([]uint64{15}),
			set2:     ValueOf([]uint64{10}),
			expected: ValueOf([]uint64{15}),
		},
		{
			set1:     ValueOf([]uint64{15, 32}),
			set2:     ValueOf([]uint64{10}),
			expected: ValueOf([]uint64{15, 32}),
		},
		{
			set1:     ValueOf([]uint64{15}),
			set2:     ValueOf([]uint64{10, 32}),
			expected: ValueOf([]uint64{15}),
		},
	}

	for _, test := range testCases {
		result := test.set1.Or(test.set2)
		assert.True(t, test.expected.Equal(result))
	}
}

func TestGetRange(t *testing.T) {
	testCases := []struct {
		set       *Set
		fromIndex int
		toIndex   int
		expected  *Set
	}{
		{
			set:       ValueOf([]uint64{10}),
			fromIndex: 1,
			toIndex:   4,
			expected:  ValueOf([]uint64{5}),
		},
		{
			set:       ValueOf([]uint64{2, 1}),
			fromIndex: 1,
			toIndex:   65,
			expected:  ValueOf([]uint64{(1 << 63) + 1}),
		},
		{
			set:       ValueOf([]uint64{15}),
			fromIndex: 4,
			toIndex:   10,
			expected:  ValueOf([]uint64{0}),
		},
	}

	for _, test := range testCases {
		result := test.set.GetRange(test.fromIndex, test.toIndex)
		assert.True(t, test.expected.Equal(result))
	}
}
