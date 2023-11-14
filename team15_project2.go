package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"team15_project1/binaryConvert"
)

type Instruction struct {
	instructionName string
	instructionType string
	rm              int32
	rn              int32
	rd              int32
	immediate       int32
	shiftCode       int32
	shiftType       int32
	conditional     int32
	offset          int32
	address         int32
	op2             int32
}

func main() {
	var instructionQueue []Instruction

	var registers [32]int32
	var memory [1024]int32

	// Pre-filling registers with arbitrary data for testing purposes
	for i := range registers {
		registers[i] = int32(i)
	}
	for i := range memory {
		memory[i] = 0
	}

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

	outputFile, err := os.Create(*OutputFileName + "_dis.txt")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	outputFile2, err := os.Create(*OutputFileName + "_sim.txt")
	if err != nil {
		panic(err)
	}
	defer outputFile2.Close()

	scanner := bufio.NewScanner(inputFile)

	programCounter := 96
	instructionCounter := 0
	cycleCounter := 1

	//Start reading instructions
	for scanner.Scan() {
		line := scanner.Text()
		opcode := binaryConvert.BinaryStringToInt(line[:11])
		opcodeString := binaryConvert.IntToInstruction(opcode)

		// Once we see BREAK we'll exit this loop and read the rest of the file as data only
		if opcodeString == "BREAK" {
			var BreakInstruction Instruction
			BreakInstruction.instructionName = "BREAK"
			BreakInstruction.instructionType = "BREAK"
			fmt.Fprintf(outputFile, "%s%s %s %s %s %s %s\t%d\t%s\n", line[:1], line[1:6], line[6:11], line[11:16], line[16:21], line[21:26], line[26:], programCounter, opcodeString)
			break
		}

		insType := binaryConvert.GetInstructionType(opcodeString)

		switch insType {
		case "R":
			//rm, rn, rd etc are labels given in the lecture 7 slides
			rm := line[11:16]
			rn := line[22:27]
			rd := line[27:32]
			fmt.Fprintf(outputFile, "%s\t%d\t%s R%d, R%d, R%d\n", line[:11]+" "+line[11:16]+" "+line[16:22]+" "+line[22:27]+" "+line[27:32], programCounter, opcodeString, binaryConvert.BinaryStringToInt(rd), binaryConvert.BinaryStringToInt(rn), binaryConvert.BinaryStringToInt(rm))

			var registerInstruction Instruction
			registerInstruction.instructionName = opcodeString
			registerInstruction.instructionType = insType
			registerInstruction.rd = binaryConvert.BinaryStringToInt(rd)
			registerInstruction.rn = binaryConvert.BinaryStringToInt(rn)
			registerInstruction.rm = binaryConvert.BinaryStringToInt(rm)
			instructionQueue = append(instructionQueue, registerInstruction)

		case "RL":
			//rm, rn, rd etc are labels given in the lecture 7 slides
			immediate := binaryConvert.BinaryStringToInt(line[16:22])
			rn := line[22:27]
			rd := line[27:32]

			fmt.Fprintf(outputFile, "%s\t%d\t%s R%d, R%d, #%d\n", line[:11]+" "+line[11:16]+" "+line[16:22]+" "+line[22:27]+" "+line[27:32], programCounter, opcodeString, binaryConvert.BinaryStringToInt(rn), binaryConvert.BinaryStringToInt(rd), immediate)

			var shiftInstruction Instruction
			shiftInstruction.instructionName = opcodeString
			shiftInstruction.instructionType = insType
			shiftInstruction.immediate = immediate
			shiftInstruction.rn = binaryConvert.BinaryStringToInt(rn)
			shiftInstruction.rd = binaryConvert.BinaryStringToInt(rd)
			instructionQueue = append(instructionQueue, shiftInstruction)

		case "I":
			immediate := binaryConvert.BinaryStringToInt(line[10:22])
			rn := binaryConvert.BinaryStringToInt(line[22:27])
			rd := binaryConvert.BinaryStringToInt(line[27:32])
			fmt.Fprintf(outputFile, "%s %s %s %s\t%d\t%s R%d, R%d, #%d\n", line[:10], line[10:22], line[22:27], line[27:32], programCounter, opcodeString, rd, rn, immediate)

			var immediateInstruction Instruction
			immediateInstruction.instructionName = opcodeString
			immediateInstruction.instructionType = insType
			immediateInstruction.rd = rd
			immediateInstruction.rn = rn
			immediateInstruction.immediate = immediate
			instructionQueue = append(instructionQueue, immediateInstruction)

		case "IM":
			immediate := binaryConvert.BinaryStringToInt(line[11:27])
			shiftCode := binaryConvert.BinaryStringToInt(line[9:11])
			rd := binaryConvert.BinaryStringToInt(line[27:32])
			shiftType := shiftCode * 16

			var IMInstruction Instruction
			IMInstruction.instructionName = opcodeString
			IMInstruction.instructionType = insType
			IMInstruction.immediate = immediate
			IMInstruction.shiftCode = shiftCode
			IMInstruction.rd = rd
			IMInstruction.shiftType = shiftType
			instructionQueue = append(instructionQueue, IMInstruction)

			fmt.Fprintf(outputFile, "%s %s %s %s\t%d\t%s R%d, %d LSL %d\n", line[:10], line[10:12], line[12:28], line[28:], programCounter, opcodeString, rd, immediate, shiftType)
		case "CB":
			conditional := binaryConvert.BinaryStringToInt(line[27:32])
			offset := binaryConvert.BinaryStringToInt(line[8:27])
			fmt.Fprintf(outputFile, "%s\t%d\t%s R%d, #%d\n", line[:8]+" "+line[8:27]+" "+line[27:32], programCounter, opcodeString, conditional, offset)
		case "B":
			offset := binaryConvert.BinaryStringToInt(line[6:32])
			fmt.Fprintf(outputFile, "%s %s\t%d\t%s #%d\n", line[:6], line[6:32], programCounter, opcodeString, offset)
		case "D":
			address := binaryConvert.BinaryStringToInt(line[11:20])
			//op2 := binaryConvert.BinaryStringToInt(line[20:22])
			rn := binaryConvert.BinaryStringToInt(line[22:27])
			rt := binaryConvert.BinaryStringToInt(line[27:32])
			fmt.Fprintf(outputFile, "%s %s %s %s %s\t%d\t%s R%d, [R%d, #%d]\n", line[:11], line[11:20], line[20:22], line[22:27], line[27:32], programCounter, opcodeString, rt, rn, address)
		case "NOP":
			fmt.Fprintf(outputFile, "%s\t%d\t%s\n", line, programCounter, opcodeString)
		default: // Instruction cannot be identified
			fmt.Fprintf(outputFile, "%s\tUnknown Value\t%d\n", line, programCounter)
		}

		programCounter += 4
		instructionCounter++
	}

	// Read file until the end as data, write data to file
	for scanner.Scan() {
		programCounter += 4
		line := scanner.Text()
		data := binaryConvert.BinaryStringToInt(line)
		fmt.Fprintf(outputFile, "%s\t%d\t%d\n", line, programCounter, data)
	}

	// Loop to read through instructionQueue, execute instructions, and write to file
	for i := range instructionQueue {
		switch instructionQueue[i].instructionType {
		case "R":
			// rd = rn + rm
			rd := int(instructionQueue[i].rd)
			rn := int(instructionQueue[i].rn)
			rm := int(instructionQueue[i].rm)

			switch instructionQueue[i].instructionName {
			case "ADD":
				registers[rd] = registers[rn] + registers[rm]
			case "SUB":
				registers[rd] = registers[rn] - registers[rm]
			case "AND":
				registers[rd] = registers[rn] & registers[rm]
			case "ORR":
				registers[rd] = registers[rn] | registers[rm]
			}

		case "RL":
			//rd = rn shift by immediate
			rd := int(instructionQueue[i].rd)
			rn := int(instructionQueue[i].rn)
			immediate := instructionQueue[i].immediate

			switch instructionQueue[i].instructionName {
			case "LSL":
				registers[rd] = registers[rn] << immediate
			case "LSR":
				registers[rd] = registers[rn] >> immediate
			case "ASR":
				registers[rd] = registers[rn] >> immediate
			}
		case "I":
			//rd = rn + immediate
			rd := int(instructionQueue[i].rd)
			rn := int(instructionQueue[i].rn)
			immediate := instructionQueue[i].immediate

			switch instructionQueue[i].instructionName {
			case "ADDI":
				registers[rd] = registers[rn] + immediate
			case "SUBI":
				registers[rd] = registers[rn] - immediate
			}
		case "IM":

		case "CB":

		case "B":

		case "D":

		case "NOP":
			fmt.Println("NOP")
		case "BREAK":
			break
		default: // Instruction cannot be identified

		}
	}

	for i := range instructionQueue {
		fmt.Fprintln(outputFile2, "====================")
		fmt.Fprintln(outputFile2, "Cycle:", cycleCounter, "\t", instructionQueue[i].instructionName)
		fmt.Fprintln(outputFile2, "registers:")
		for i := 0; i < 32; i += 8 {
			fmt.Fprintf(outputFile2, "r%02d:\t", i)
			for j := i; j < i+8; j++ {
				fmt.Fprintf(outputFile2, "%d\t", registers[j])
			}
			fmt.Fprintln(outputFile2)

		}
		fmt.Fprintln(outputFile2)
		fmt.Fprintln(outputFile2, "data:")
		for i := 0; i < 8; i += 1 {
			//fmt.Fprint(outputFile2, programCounter, ":", 0)
			for j := i; j < i+1; j++ {
				fmt.Fprintf(outputFile2, "%d\t", memory[j])
			}
			fmt.Fprint(outputFile2)
		}
		fmt.Fprintln(outputFile2)

		cycleCounter += 1

		if instructionQueue[i].instructionName == "BREAK" {
			break
		}
	}
}
