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
	OutputFileName := flag.String("o", "", "Gets the output file name")

	flag.Parse()

	if flag.NArg() != 0 {
		os.Exit(200)
	}
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

		// Once we see BREAK we'll exit this loop and read the rest of the file as data only
		if opcodeString == "BREAK" {
			break
		}

		insType := binaryConvert.GetInstructionType(opcodeString)

		switch insType {
		case "R":
			//rm, rn, rd etc are labels given in the lecture 7 slides
			rm := line[11:16]
			rn := line[22:27]
			rd := line[27:32]
			fmt.Printf("%s  R%d, R%d, R%d\n", opcodeString, binaryConvert.BinaryStringToInt(rd), binaryConvert.BinaryStringToInt(rn), binaryConvert.BinaryStringToInt(rm))
		case "I":
			immediate := binaryConvert.BinaryStringToInt(line[10:22])
			rn := binaryConvert.BinaryStringToInt(line[22:27])
			rd := binaryConvert.BinaryStringToInt(line[27:32])
			fmt.Printf("%s  R%d, R%d, R%d\n", opcodeString, rd, rn, immediate)
		case "IM":
			immediate := binaryConvert.BinaryStringToInt(line[10:22])
			shiftCode := binaryConvert.BinaryStringToInt(line[22:24])
			rd := binaryConvert.BinaryStringToInt(line[27:32])
			fmt.Printf("%s R%d, #%d, LSL #%d\n", opcodeString, rd, immediate, shiftCode)
		case "CB":
			offset := binaryConvert.BinaryStringToInt(line[8:27])
			conditional := binaryConvert.BinaryStringToInt(line[27:32])
			fmt.Printf("%s R%d, #%d\n", opcodeString, offset, conditional)
		case "B":
			offset := binaryConvert.BinaryStringToInt(line[6:32])
			fmt.Printf("B #%d\n", offset)
		case "D":
			address := binaryConvert.BinaryStringToInt(line[11:20])
			//op2 := binaryConvert.BinaryStringToInt(line[20:22])
			rn := binaryConvert.BinaryStringToInt(line[22:27])
			rt := binaryConvert.BinaryStringToInt(line[27:32])
			fmt.Printf("%s R%d, [R%d, #%d]\n", opcodeString, rt, rn, address)
		}

		programCounter += 4
	}

	OutputFile, err := os.Create(*OutputFileName)
	if err != nil {
		panic(err)
	}
	err = OutputFile.Close()

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
