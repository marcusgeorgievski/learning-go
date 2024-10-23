package main

import (
	"fmt"
	"slices"
)

func main() {
	// Arrays()
	// Slices1()
	// Slices2()
	// Strings1()
	// Ex2()
	Ex3()
}

func Arrays() {
	var arr [2]int
	var arr2 [2]int
	var arr3 = [3]int{1, 2: 99}
	var arr4 = [...]int{1, 2, 3, 4}

	fmt.Println(arr, arr3, arr4)
	fmt.Println(arr == arr2)
}

func Slices1() {
	var snil []int
	var s = []int{1, 2, 3}
	var ss = []int{1, 2, 3}

	fmt.Println(snil)
	fmt.Println(slices.Equal(s, ss))
	s = append(s, ss...)
	fmt.Println(s)
}
func Slices2() {
	x := make([]int, 0, 5)
	x = append(x, 1, 2, 3, 4)

	y := x[2:4:5]

	y = append(y, 5)

	fmt.Println(x, len(x), cap(x))
	fmt.Println(y, len(y), cap(y))
}

func Strings1() {
	var s string = "hello 😀"
	var s2 string = s[3:7]
	var b byte = s[6]
	var bs []byte = []byte(s[6:])

	fmt.Println(s, s2, b, bs)
}

func Ex2() {
	s := "Hi 😀 and 😁"
	rs := []rune(s)

	// print emoji
	fmt.Println(string(rs[3]))
}

type Employee struct {
	firstName string
	lastName  string
	id        int
}

func Ex3() {
	e1 := Employee{
		firstName: "M",
		lastName:  "G",
		id:        1,
	}

	e2 := Employee{
		"E",
		"M",
		2,
	}

	var e3 Employee

	e3.firstName = "A"
	e3.lastName = "B"
	e3.id = 3

	fmt.Println(e1, e2, e3)
}
