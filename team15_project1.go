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
	defer inputFile.Close()

	outputFile, err := os.Create(*OutputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)

	programCounter := 96
	instructionCounter := 0

	//Start reading instructions
	for scanner.Scan() {
		line := scanner.Text()
		opcode := binaryConvert.BinaryStringToInt(line[:11])
		opcodeString := binaryConvert.IntToInstruction(opcode)

		// Once we see BREAK we'll exit this loop and read the rest of the file as data only
		if opcodeString == "BREAK" {

			fmt.Fprintf(outputFile, "%s%s %s %s %s %s %s    %d    %s\n", line[:1], line[1:6], line[6:11], line[11:16], line[16:21], line[21:26], line[26:], programCounter, opcodeString)
			break
		}

		insType := binaryConvert.GetInstructionType(opcodeString)

		switch insType {
		case "R":
			//rm, rn, rd etc are labels given in the lecture 7 slides
			rm := line[11:16]
			rn := line[22:27]
			rd := line[27:32]
			fmt.Fprintf(outputFile, "%s    %d     %s    R%d,    R%d,    R%d\n", line[:11]+" "+line[11:16]+" "+line[16:21]+" "+line[21:26]+" "+line[26:31]+" "+line[31:], programCounter, opcodeString, binaryConvert.BinaryStringToInt(rd), binaryConvert.BinaryStringToInt(rn), binaryConvert.BinaryStringToInt(rm))
		case "I":
			immediate := binaryConvert.BinaryStringToInt(line[10:22])
			rn := binaryConvert.BinaryStringToInt(line[22:27])
			rd := binaryConvert.BinaryStringToInt(line[27:32])
			fmt.Fprintf(outputFile, "%s      %d     %s   R%d,    R%d,    #%d\n", line[:10]+" "+line[10:22]+" "+line[22:27]+" "+line[27:32], programCounter, opcodeString, rd, rn, immediate)
		case "IM":
			immediate := binaryConvert.BinaryStringToInt(line[10:22])
			shiftCode := binaryConvert.BinaryStringToInt(line[22:24])
			rd := binaryConvert.BinaryStringToInt(line[27:32])
			shiftType := ""

			switch shiftCode {
			case 0:
				shiftType = ""
			case 1:
				shiftType = ", LSL    "
			case 2:
				shiftType = ", LSR"
			case 3:
				shiftType = ", ASR"
			}

			fmt.Fprintf(outputFile, "%s     %d    %s    R%d,    #%d%s\n", line[:10]+" "+line[10:22]+" "+line[22:24]+" "+line[24:27]+" "+line[27:32], programCounter, opcodeString, rd, immediate, shiftType)
		case "CB":
			offset := binaryConvert.BinaryStringToInt(line[8:27])
			conditional := binaryConvert.BinaryStringToInt(line[27:32])
			fmt.Fprintf(outputFile, "%s       %d     %s    R%d,    #%d\n", line[:8]+" "+line[8:27]+" "+line[27:32], programCounter, opcodeString, offset, conditional)
		case "B":
			offset := binaryConvert.BinaryStringToInt(line[6:32])
			fmt.Fprintf(outputFile, "%s        %d     %s    #%d\n", line[:6]+" "+line[6:32], programCounter, opcodeString, offset)
		case "D":
			address := binaryConvert.BinaryStringToInt(line[11:20])
			//op2 := binaryConvert.BinaryStringToInt(line[20:22])
			rn := binaryConvert.BinaryStringToInt(line[22:27])
			rt := binaryConvert.BinaryStringToInt(line[27:32])
			fmt.Fprintf(outputFile, "%s     %d     %s    R%d,   [R%d, #%d]\n", line[:11]+" "+line[11:20]+" "+line[20:22]+" "+line[22:27]+" "+line[27:32], programCounter, opcodeString, rt, rn, address)
		case "NOP":
			fmt.Fprintf(outputFile, "%s     %d    %s\n", line, programCounter, opcodeString)
		default: // Instruction cannot be identified
			fmt.Fprintf(outputFile, "%s   Unknown Value   %d\n", line, programCounter)
		}

		programCounter += 4
		instructionCounter++
	}

	for scanner.Scan() {
		line := scanner.Text()
		data := binaryConvert.BinaryStringToInt(line)
		fmt.Fprintf(outputFile, "%s         %d    %d\n", line, programCounter, data)
		programCounter += 4
	}
}
