package bit

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	nbits := 65
	s := NewSet(WithInitialBits(nbits))
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

	//-------------------------
	// SetRange
	// ------------------------
	s.SetRange(50, nbits)
	for i := 50; i < nbits; i++ {
		trueIndexes[i] = true
	}
	checkBits(t, s, trueIndexes)

	// ------------------------
	// ClearRange
	// ------------------------
	s.ClearRange(51, nbits)
	for i := 51; i < nbits; i++ {
		delete(trueIndexes, i)
	}
	checkBits(t, s, trueIndexes)

	assert.True(t, s.Get(50))
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
