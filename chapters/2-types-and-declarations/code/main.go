package main

import (
	"fmt"
)

const (
	x int = 5
	y     = 5
)

func main() {

	var bin int = 0b100
	var hex int = 0x10
	var oct int = 0o3
	var exp float64 = 10e2 // 100
	var consty float64 = y // x would not work

	fmt.Println(bin, hex, oct, exp, consty)

	var char rune = 'a'
	var raw string = `\hi \n`

	fmt.Println(char, raw)

	var str string  // ""
	var integer int // 0
	var booly bool  // false

	fmt.Println(str, integer, booly)

	var _u8 byte
	var _i32 rune

	fmt.Println(_u8-1, _i32)
}
