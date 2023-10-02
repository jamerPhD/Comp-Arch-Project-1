package binaryConvert

import "strconv"

// InstructionType
/*
Maps instruction to its type
parameter: instruction string (ie AND)
returns: instruction type string (ie R)
*/
var InstructionType = map[string]string{
	"B":    "B",
	"AND":  "R",
	"ADD":  "R",
	"ADDI": "I",
	"ORR":  "R",
	"CBZ":  "CB",
	"CBNZ": "CB",
	"SUB,": "R",
	"SUBI": "I",
	"MOVZ": "IM",
	"MOVK": "IM",
	"LSR":  "R",
	"LSL":  "R",
	"STUR": "D",
	"LDUR": "D",
	"ASR":  "R",
	"EOR":  "R",
}

// IntToInstruction
/*
converts an integer value (from a binary instruction) into a string OPCODE
parameter: int32
return: string
*/
func IntToInstruction(value int32) string {
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

	return "???" // Instruction not found

}

// BinaryStringToInt
/*
converts a string of 11 characters (binary number) to a base 10 integer
parameter: string
return: int
*/
func BinaryStringToInt(binary string) int32 {
	if i64, err := strconv.ParseInt(binary, 2, 32); err != nil {
		return (BinaryStringToInt(twosComplement(binary)) + 1) * -1
	} else {
		return int32(i64)
	}
}

// twosComplement
/*
helper function for binaryStringToInt
this is called if the ParseInt on the binary string overflows in binaryStringToInt
inverts the bits then returns the inverted string
can also throw panic if the string contains non-binary characters for some reason
parameter: string
return: string with inverted bits
*/
func twosComplement(binary string) string {
	var complement = ""

	for i := 0; i < len(binary); i++ {
		if binary[i] == '0' {
			complement += "1"
		} else if binary[i] == '1' {
			complement += "0"
		} else {
			panic("Binary string contains invalid character")
		}
	}

	return complement
}
