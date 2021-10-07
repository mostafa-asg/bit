[![Go Report Card](https://goreportcard.com/badge/github.com/mostafa-asg/bit)](https://goreportcard.com/report/github.com/mostafa-asg/bit)
# BitSet
Package bit implements vector of bits that grows as needed. Basically it is a memory-effiecent version of array of boolean: `[]bool`. All the methods implemented in this package inspired by Java's [BitSet](https://docs.oracle.com/javase/8/docs/api/java/util/BitSet.html)

## Method Summary
 - [And(otherSet *Set)](#)
	 - Performs a logical **AND** of this target bit set with the argument bit set.
 - [AndNot(otherSet *Set)](#)
	 - Clears all of the bits in this BitSet whose corresponding bit is set in the specified BitSet.
 - [Cardinality()](#)
	 - Returns the number of bits set to true in this BitSet.
 - [ClearAll()](#)
	 - Sets all of the bits in this BitSet to false.
 - [Clear(index int)](#)
	 - Sets the bit specified by the index to false.
 - [ClearRange(fromIndex int, toIndex int)](#)
	 - Sets the bits from the specified fromIndex (inclusive) to the specified toIndex (exclusive) to false.
 - [Clone()](#)
	 - Cloning this BitSet produces a new BitSet that is equal to it.
 - [Equal(otherSet *Set)](#)
	 - checks equality between this set and the other set passed in the argument.
 - [Flip(index int)](#)
	 - Sets the bit at the specified index to the complement of its current value.
 - [FlipRange(fromIndex int, toIndex int)](#)
	 - Sets each bit from the specified fromIndex (inclusive) to the specified toIndex (exclusive) to the complement of its current value.
 - [Get(index int)](#)
	 - Returns the value of the bit with the specified index.
 - [Intersects(otherSet *Set)](#)
	 - Returns true if the specified BitSet has any bits set to true that are also set to true in this BitSet.
 - [IsEmpty()](#)
   - Returns true if this BitSet contains no bits that are set to true.
 - [Length()](#)
   - Returns the "logical size" of this BitSet: the index of the highest set bit in the BitSet plus one.
 - [NextClearBit(fromIndex int)](#)
   - Returns the index of the first bit that is set to false that occurs on or after the specified starting index.
 - [NextSetBit(fromIndex int)](#)
	 - Returns the index of the first bit that is set to true that occurs on or after the specified starting index.
 - [PreviousClearBit(fromIndex int)](#)
	 - Returns the index of the nearest bit that is set to false that occurs on or before the specified starting index.
 - [PreviousSetBit(fromIndex int)](#)
	 - Returns the index of the nearest bit that is set to true that occurs on or before the specified starting index. 
 - [Set(index int)](#)
	 - Sets the bit at the specified index to true.
 - [SetValue(index int, value bool)](#)
	 - Sets the bit at the specified index to the specified value.
 - [SetRange(fromIndex int, toIndex int)](#)
	 - Sets the bits from the specified fromIndex (inclusive) to the specified toIndex (exclusive) to true.
 - [SetRangeValue(fromIndex int, toIndex int, value bool)](#)
	 - Sets the bits from the specified fromIndex (inclusive) to the specified toIndex (exclusive) to the specified value.
 - [Size()](#)
	 - Returns the number of bits of space actually in use by this BitSet to represent bit values.
 - [Bytes() []byte](#)
	 - Returns a new byte array containing all the bits in this bit set. 
 - [ToArray() []uint64](#)
	 - Returns a new array containing all the bits in this bit set.
 - [String()](#)
	 - Returns a string representation of this bit set.
 - [ValueOf(nums []uint64)](#)
	 - Returns a new bit set containing all the bits in the given long array.
 - [FromByteArray(nums []byte)](#)
   - Returns a new bit set containing all the bits in the given byte array.
 - [Xor(otherSet *Set)](#)
   - Performs a logical XOR of this bit set with the bit set argument.

  

