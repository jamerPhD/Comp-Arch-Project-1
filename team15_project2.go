package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

// InstructionType
// Maps instruction to its type
// parameter: instruction string (ie AND)
// returns: instruction type string (ie R)
var instructionType = map[string]string{
	"B":    "B",
	"AND":  "R",
	"ADD":  "R",
	"ADDI": "I",
	"ORR":  "R",
	"CBZ":  "CB",
	"CBNZ": "CB",
	"SUB":  "R",
	"SUBI": "I",
	"MOVZ": "IM",
	"MOVK": "IM",
	"LSR":  "RL",
	"LSL":  "RL",
	"STUR": "D",
	"LDUR": "D",
	"ASR":  "RL",
	"EOR":  "R",
	"NOP":  "NOP",
}

// GetInstructionType
// Accessor for instructionType map
// Returns a string of the instruction type
func GetInstructionType(key string) string {
	if key == "Unknown Value" {
		return key
	} else {
		return instructionType[key]
	}

}

// IntToInstruction
// converts an integer value (from a binary instruction) into a string OPCODE
// parameter: int32
// return: string
func IntToInstruction(value int32) string {
	if value == 0 {
		return "NOP"
	}
	if value >= 160 && value <= 191 {
		return "B"
	}
	if value == 1104 {
		return "AND"
	}
	if value == 1112 {
		return "ADD"
	}
	if value == 1160 || value == 1161 {
		return "ADDI"
	}
	if value == 1360 {
		return "ORR"
	}
	if value >= 1440 && value <= 1447 {
		return "CBZ"
	}
	if value >= 1448 && value <= 1455 {
		return "CBNZ"
	}
	if value == 1624 {
		return "SUB"
	}
	if value == 1672 || value == 1673 {
		return "SUBI"
	}
	if value >= 1684 && value <= 1687 {
		return "MOVZ"
	}
	if value == 1690 {
		return "LSR"
	}
	if value == 1691 {
		return "LSL"
	}
	if value == 1692 {
		return "ASR"
	}
	if value == 1872 {
		return "EOR"
	}
	if value >= 1940 && value <= 1943 {
		return "MOVK"
	}
	if value == 1984 {
		return "STUR"
	}
	if value == 1986 {
		return "LDUR"
	}
	if value == 2038 {
		return "BREAK"
	}

	return "Unknown Instruction" // Instruction not found

}

// BinaryStringToInt
// converts a string of 11 characters (binary number) to a base 10 integer
// parameter: string
// return: int
func BinaryStringToInt(binary string) int32 {
	bitSize := 26
	if len(binary) > 26 {
		bitSize = 32
	}
	if i64, err := strconv.ParseInt(binary, 2, bitSize); err != nil {
		return (BinaryStringToInt(twosComplement(binary)) + 1) * -1
	} else {
		return int32(i64)
	}
}

// twosComplement
// helper function for binaryStringToInt
// this is called if the ParseInt on the binary string overflows in binaryStringToInt
// inverts the bits then returns the inverted string
// can also throw panic if the string contains non-binary characters for some reason
// parameter: string
// return: string with inverted bits
func twosComplement(binary string) string {
	var complement = ""

	for i := 0; i < len(binary); i++ {
		if binary[i] == '0' {
			complement += "1"
		} else if binary[i] == '1' {
			complement += "0"
		} else {
			panic("Binary string contains invalid character: " + binary)
		}
	}

	return complement
}

type Instruction struct {
	instructionName string // Name of instruction, ie ADD, SUB
	instructionType string // Type of instruction, ie R, IM
	instructionInfo string // Formatted instruction info, ie ADD R1, R2, R3
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

	var registers [32]int64
	var memory []int32
	var startOfData int

	InputFileName := flag.String("i", "", "Gets the input file name")
	OutputFileName := flag.String("o", "", "Gets the output file name")
	flag.Parse()

	if flag.NArg() != 0 {
		os.Exit(200)
	}

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

	//Start reading instructions
	for scanner.Scan() {
		line := scanner.Text()
		opcode := BinaryStringToInt(line[:11])
		opcodeString := IntToInstruction(opcode)

		// Once we see BREAK we'll exit this loop and read the rest of the file as data only
		if opcodeString == "BREAK" {
			var BreakInstruction Instruction
			BreakInstruction.instructionName = "BREAK"
			BreakInstruction.instructionType = "BREAK"
			BreakInstruction.instructionInfo = "BREAK"
			instructionQueue = append(instructionQueue, BreakInstruction)
			fmt.Fprintf(outputFile, "%s%s %s %s %s %s %s\t%d\t%s\n", line[:1], line[1:6], line[6:11], line[11:16], line[16:21], line[21:26], line[26:], programCounter, opcodeString)
			break
		}

		insType := GetInstructionType(opcodeString)

		switch insType {
		case "R":
			//rm, rn, rd etc are labels given in the lecture 7 slides
			rm := line[11:16]
			rn := line[22:27]
			rd := line[27:32]

			var registerInstruction Instruction
			registerInstruction.instructionName = opcodeString
			registerInstruction.instructionType = insType
			registerInstruction.rd = BinaryStringToInt(rd)
			registerInstruction.rn = BinaryStringToInt(rn)
			registerInstruction.rm = BinaryStringToInt(rm)
			registerInstruction.instructionInfo = fmt.Sprintf("%s R%d, R%d, R%d", opcodeString, registerInstruction.rd, registerInstruction.rn, registerInstruction.rm)
			instructionQueue = append(instructionQueue, registerInstruction)

			fmt.Fprintf(outputFile, "%s\t%d\t%s R%d, R%d, R%d\n", line[:11]+" "+line[11:16]+" "+line[16:22]+" "+line[22:27]+" "+line[27:32], programCounter, opcodeString, BinaryStringToInt(rd), BinaryStringToInt(rn), BinaryStringToInt(rm))
		case "RL":
			//rm, rn, rd etc are labels given in the lecture 7 slides
			immediate := BinaryStringToInt(line[16:22])
			rn := line[22:27]
			rd := line[27:32]

			var shiftInstruction Instruction
			shiftInstruction.instructionName = opcodeString
			shiftInstruction.instructionType = insType
			shiftInstruction.immediate = immediate
			shiftInstruction.rn = BinaryStringToInt(rn)
			shiftInstruction.rd = BinaryStringToInt(rd)
			shiftInstruction.instructionInfo = fmt.Sprintf("%s R%d, R%d, #%d", opcodeString, shiftInstruction.rd, shiftInstruction.rn, immediate)
			instructionQueue = append(instructionQueue, shiftInstruction)

			fmt.Fprintf(outputFile, "%s\t%d\t%s R%d, R%d, #%d\n", line[:11]+" "+line[11:16]+" "+line[16:22]+" "+line[22:27]+" "+line[27:32], programCounter, opcodeString, BinaryStringToInt(rd), BinaryStringToInt(rn), immediate)
		case "I":
			immediate := BinaryStringToInt(line[10:22])
			rn := BinaryStringToInt(line[22:27])
			rd := BinaryStringToInt(line[27:32])

			var immediateInstruction Instruction
			immediateInstruction.instructionName = opcodeString
			immediateInstruction.instructionType = insType
			immediateInstruction.rd = rd
			immediateInstruction.rn = rn
			immediateInstruction.immediate = immediate
			immediateInstruction.instructionInfo = fmt.Sprintf("%s R%d, R%d, #%d", opcodeString, rd, rn, immediate)
			instructionQueue = append(instructionQueue, immediateInstruction)

			fmt.Fprintf(outputFile, "%s %s %s %s\t%d\t%s R%d, R%d, #%d\n", line[:10], line[10:22], line[22:27], line[27:32], programCounter, opcodeString, rd, rn, immediate)
		case "IM":
			immediate := BinaryStringToInt(line[11:27])
			shiftCode := BinaryStringToInt(line[9:11])
			rd := BinaryStringToInt(line[27:32])
			shiftType := shiftCode * 16

			var IMInstruction Instruction
			IMInstruction.instructionName = opcodeString
			IMInstruction.instructionType = insType
			IMInstruction.immediate = immediate
			IMInstruction.shiftCode = shiftCode
			IMInstruction.rd = rd
			IMInstruction.shiftType = shiftType
			IMInstruction.instructionInfo = fmt.Sprintf("%s R%d, %d LSL %d", opcodeString, rd, immediate, shiftType)
			instructionQueue = append(instructionQueue, IMInstruction)

			fmt.Fprintf(outputFile, "%s %s %s %s\t%d\t%s R%d, %d LSL %d\n", line[:10], line[10:12], line[12:28], line[28:], programCounter, opcodeString, rd, immediate, shiftType)
		case "CB":
			conditional := BinaryStringToInt(line[27:32])
			offset := BinaryStringToInt(line[8:27])

			var CBInstruction Instruction
			CBInstruction.instructionName = opcodeString
			CBInstruction.instructionType = insType
			CBInstruction.conditional = conditional
			CBInstruction.offset = offset
			CBInstruction.instructionInfo = fmt.Sprintf("%s R%d, #%d", opcodeString, conditional, offset)
			instructionQueue = append(instructionQueue, CBInstruction)

			fmt.Fprintf(outputFile, "%s\t%d\t%s R%d, #%d\n", line[:8]+" "+line[8:27]+" "+line[27:32], programCounter, opcodeString, conditional, offset)
		case "B":
			offset := BinaryStringToInt(line[6:32])

			var BInstruction Instruction
			BInstruction.instructionName = opcodeString
			BInstruction.instructionType = insType
			BInstruction.offset = offset
			BInstruction.instructionInfo = fmt.Sprintf("B #%d", offset)
			instructionQueue = append(instructionQueue, BInstruction)

			fmt.Fprintf(outputFile, "%s %s\t%d\t%s #%d\n", line[:6], line[6:32], programCounter, opcodeString, offset)
		case "D":
			address := BinaryStringToInt(line[11:20])
			//op2 := BinaryStringToInt(line[20:22])
			rn := BinaryStringToInt(line[22:27])
			rt := BinaryStringToInt(line[27:32])
			rd := BinaryStringToInt(line[27:32])

			var DInstruction Instruction
			DInstruction.instructionName = opcodeString
			DInstruction.instructionType = insType
			DInstruction.rn = rn
			DInstruction.rd = rd
			DInstruction.address = address
			DInstruction.instructionInfo = fmt.Sprintf("%s R%d, [R%d, #%d]", opcodeString, rd, rn, address)
			instructionQueue = append(instructionQueue, DInstruction)

			fmt.Fprintf(outputFile, "%s %s %s %s %s\t%d\t%s R%d, [R%d, #%d]\n", line[:11], line[11:20], line[20:22], line[22:27], line[27:32], programCounter, opcodeString, rt, rn, address)
		case "NOP":
			fmt.Fprintf(outputFile, "%s\t%d\t%s\n", line, programCounter, opcodeString)
		default: // Instruction cannot be identified
			fmt.Fprintf(outputFile, "%s\tUnknown Value\t%d\n", line, programCounter)
		}

		programCounter += 4
	}

	// Read file until the end as data, write data to file
	startOfData = programCounter + 4
	for scanner.Scan() {
		programCounter += 4
		line := scanner.Text()
		data := BinaryStringToInt(line)
		memory = append(memory, data)
		fmt.Fprintf(outputFile, "%s\t%d\t%d\n", line, programCounter, data)
	}

	// Loop to read through instructionQueue, execute instructions, and write to file
	cycleCounter := 1
	programCounter = 96

	for i := 0; i < len(instructionQueue); i++ {
		fmt.Fprintln(outputFile2, "=====================")
		fmt.Fprintf(outputFile2, "cycle:%d\t%d\t%s\n\n", cycleCounter, programCounter, instructionQueue[i].instructionInfo)

		switch instructionQueue[i].instructionType {
		case "R":
			// rd = rn + rm
			rd := int(instructionQueue[i].rd)
			rn := int(instructionQueue[i].rn)
			rm := int(instructionQueue[i].rm)

			switch instructionQueue[i].instructionName {
			case "EOR":
				registers[rd] = registers[rn] ^ registers[rm]
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
				uRn := uint32(registers[rn]) >> immediate
				finalRn := int64(uRn)
				registers[rd] = finalRn
				//registers[rd] = registers[rn] >> immediate
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
				registers[rd] = registers[rn] + int64(immediate)
			case "SUBI":
				registers[rd] = registers[rn] - int64(immediate)
			}
		case "IM":
			rd := int(instructionQueue[i].rd)
			immediate := instructionQueue[i].immediate
			shiftType := instructionQueue[i].shiftType

			switch instructionQueue[i].instructionName {
			case "MOVK":
				registers[rd] = registers[rd] + int64(immediate<<shiftType)
			case "MOVZ":
				registers[rd] = int64(immediate << shiftType)
			}
		case "CB":
			conditional := registers[instructionQueue[i].conditional]
			switch instructionQueue[i].instructionName {
			case "CBZ":
				if conditional == 0 {
					programCounter += int(instructionQueue[i].offset*4) - 4
					i += int(instructionQueue[i].offset) - 1
				}
			case "CBNZ":
				if conditional != 0 {
					programCounter += int(instructionQueue[i].offset*4) - 4
					i += int(instructionQueue[i].offset) - 1
				}
			}
		case "B":
			programCounter += int(instructionQueue[i].offset*4) - 4
			i += int(instructionQueue[i].offset) - 1
		case "D":
			switch instructionQueue[i].instructionName {
			case "STUR":
				rd := int(instructionQueue[i].rd)
				rn := int(registers[instructionQueue[i].rn])
				offset := int(instructionQueue[i].address * 4)
				memLoc := (rn + offset - startOfData) / 4
				for memLoc > len(memory)-1 {
					//append 0 to memory
					memory = append(memory, 0)
				}
				memory[memLoc] = int32(registers[rd])
			case "LDUR":
				rd := instructionQueue[i].rd
				rn := registers[instructionQueue[i].rn]
				offset := instructionQueue[i].address * 4
				//rd = rn + offset
				memLoc := (int(rn+int64(offset)) - startOfData) / 4
				registers[rd] = int64(memory[memLoc])
			}
		case "NOP":
			fmt.Println("NOP")
		case "BREAK":
			break
		default: // Instruction cannot be identified
			fmt.Println("Instruction not identified")
		}

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
		memAddress := startOfData
		if len(memory) != 0 {
			for i := range memory {
				if i%8 == 0 {
					fmt.Fprintf(outputFile2, "%d:", memAddress)
				}
				fmt.Fprintf(outputFile2, "%d\t", memory[i])
				if i%8 == 7 {
					fmt.Fprintln(outputFile2)
					memAddress += 32
				}
			}
			// If len(memory) isn't divisible by 8, print the remaining 0s
			for i := len(memory) % 8; i < 8; i++ {
				fmt.Fprintf(outputFile2, "0\t")
			}
		}

		fmt.Fprintln(outputFile2)

		cycleCounter += 1
		programCounter += 4
	}
}
