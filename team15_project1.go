package main

import (
	"bufio"
	"fmt"
	"os"
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
		//"abcdefghijklmnopqrstuvwxyz",     // invalid, triggers panic
	}

	InputFileName := "test1_bin.txt"
	//OutputFileName := "team15_out"

	inputFile, err := os.Open(InputFileName)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(inputFile)

	fmt.Println("Testing reading entire file and output conversions to console")
	readData := false
	for scanner.Scan() {
		line := scanner.Text()
		opcode := binaryConvert.BinaryStringToInt(line[:11])
		opcodeString := binaryConvert.IntToInstruction(opcode)
		insType := binaryConvert.InstructionType[opcodeString]
		fmt.Println("Instruction Type: " + insType)

		//rm := line[11:16]
		//shamt := line[16:22]
		//rn := line[22:27]
		//rd := line[27:32]
		//fmt.Printf("rm: %s\nshamt: %s\nrn: %s\nrd: %s\n", rm, shamt, rn, rd)

		if opcodeString == "BREAK" {
			readData = true
		}

		if !readData {
			fmt.Println(opcodeString)
		} else {
			fmt.Println(opcode)
		}

	}

	err = inputFile.Close()
	if err != nil {
		panic(err)
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
