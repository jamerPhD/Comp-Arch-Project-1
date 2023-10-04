package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"project1/binaryConvert"
)

func main() {
	InputFileName := flag.String("i", "", "Gets the input file name")
	//OutputFileName := flag.String("o", "", "Gets the output file name")

	flag.Parse()

	if flag.NArg() != 0 {
		os.Exit(200)
	}

	fmt.Println(InputFileName)
	fmt.Println(*InputFileName)
	//InputFileName := "test1_bin.txt"
	//OutputFileName := "team15_out"

	inputFile, err := os.Open(*InputFileName)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(inputFile)

	programCounter := 96
	//Start reading instructions
	for scanner.Scan() {
		line := scanner.Text()
		opcode := binaryConvert.BinaryStringToInt(line[:11])
		opcodeString := binaryConvert.IntToInstruction(opcode)
		insType := binaryConvert.GetInstructionType(opcodeString)

		//TODO
		//Put file write statements into this switch
		//James H: I added some print statements just to demo how to use my package
		switch insType {
		case "R":
			//rm, rn, rd etc are labels given in the lecture 7 slides
			rm := line[11:16]
			//shamt := line[16:22]
			rn := line[22:27]
			rd := line[27:32]
			fmt.Printf("%s  R%d, R%d, R%d\n", opcodeString, binaryConvert.BinaryStringToInt(rd), binaryConvert.BinaryStringToInt(rn), binaryConvert.BinaryStringToInt(rm))
		case "I":
			//immediate := line[10:22]
			//rn := line[22:27]
			//rd := line[27:32]
		case "IM":
		case "CB":
		case "B":
			offset := binaryConvert.BinaryStringToInt(line[6:32])
			fmt.Printf("B #%d\n", offset)
		}

		programCounter += 4

		// Once we see BREAK we'll exit this loop and read the rest of the file as data only
		if opcodeString == "BREAK" {
			break
		}

		fmt.Println("Instruction Type: " + insType)
		fmt.Println("OPCODE: " + opcodeString)
	}

	// Read data
	//TODO
	//Write data to file
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
