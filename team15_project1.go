package main

import (
	"bufio"
	"fmt"
	"os"
	"project1/binaryConvert"
)

func main() {
	InputFileName := "test1_bin.txt"
	//OutputFileName := "team15_out"

	inputFile, err := os.Open(InputFileName)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(inputFile)

	fmt.Println("Testing reading entire file and output conversions to console")

	//// When readData == true the loop will start reading data rather than instructions
	//readData := false
	for scanner.Scan() {
		line := scanner.Text()
		opcode := binaryConvert.BinaryStringToInt(line[:11])
		opcodeString := binaryConvert.IntToInstruction(opcode)
		insType := binaryConvert.GetInstructionType(opcodeString)
		//rm := line[11:16]
		//shamt := line[16:22]
		//rn := line[22:27]
		//rd := line[27:32]
		//fmt.Printf("rm: %s\nshamt: %s\nrn: %s\nrd: %s\n", rm, shamt, rn, rd)

		// Once we see BREAK we'll exit this loop and read the rest of the file as data only
		if opcodeString == "BREAK" {
			break
		}
		fmt.Println("Instruction Type: " + insType)
		fmt.Println("OPCODE: " + opcodeString)
	}

	// Read data
	for scanner.Scan() {
		line := scanner.Text()
		data := binaryConvert.BinaryStringToInt(line)
		fmt.Println(data)
	}

	err = inputFile.Close()
	if err != nil {
		panic(err)
	}

}
