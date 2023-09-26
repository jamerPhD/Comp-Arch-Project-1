package main

import (
	"fmt"
	"project1/binaryConvert"
)

func main() {

	var testStrings = []string{
		"10001010000",                      // "AND" == 1104
		"10101010000",                      // "ORR" == 1360
		"11010100101",                      // Just a random number 1701
		"110101",                           // Random smaller number (53)
		"11111111111111111111111111111111", // -1 w/ 2s complement
		"11111111111111111111111111111110", // -2 w/ 2s complement
		"11111111111111111111111111111101", // -3 w/ 2s complement
		//"abcdefghijklmnopqrstuvwxyz",     // invalid, throws panic
	}

	fmt.Println("Testing binaryStringToInt")
	for i := 0; i < len(testStrings); i++ {
		fmt.Println(binaryConvert.BinaryStringToInt(testStrings[i]))
	}

	fmt.Println("Testing binary string to instruction conversion")
	for i := 0; i < len(testStrings); i++ {
		x := binaryConvert.BinaryStringToInt(testStrings[i])

		fmt.Println(binaryConvert.IntToInstruction(x))
	}

}
