package bit_test

import (
	"fmt"

	"github.com/mostafa-asg/bit"
)

// ExampleNewSet shows how to create a bit set and do some basic operation on it
func ExampleNewSet() {
	set1, _ := bit.NewSet()

	// how many bits allocated in the memory?
	fmt.Println(set1.Size())

	// set the index number 4
	set1.Set(4)

	// get the value of index 3
	fmt.Println(set1.Get(3))

	// get the value of index 4
	fmt.Println(set1.Get(4))

	// clear the value at index 4
	set1.Clear(4)
	fmt.Println(set1.Get(4))

	// Output:
	// 64
	// false
	// true
	// false
}

// ExampleValueOf shows how to create a set based on bytes of numbers
// For instance the binary representation of 9 is 1001
// So bits index 0 and 3 are set to true
func ExampleValueOf() {
	set1 := bit.ValueOf([]uint64{9})

	// how many bits allocated in the memory?
	fmt.Println(set1.Size())

	fmt.Println(set1.Get(0))
	fmt.Println(set1.Get(3))

	set1.Set(1)
	set1.Set(2)

	// print all the set bits
	fmt.Println(set1)

	// Output:
	// 64
	// true
	// true
	// {0, 1, 2, 3}
}
