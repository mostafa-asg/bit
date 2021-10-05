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
