package cpu

import (
	"fmt"
	"strings"
)

// {"AND", "(HL)", "", 0x00, 0x0000, cpu{a: 0x12, h: 0xaa, l: 0xbb}, cpu{a: 0x10, h: 0xaa, l: 0xbb}, map[string]bool{"Z": false}, testMemory{mem: map[uint16]byte{0xaabb: 0x34}}},

var format = `{"%s", "%s", "%s", 0x%02x, 0x%04x, cpu{ %s }, cpu{%s}, map[string]bool{ %s }, %s},`

func generateCPUByOperand(operand string) string {
	switch {
	case strings.Contains(operand, "BC"):
		return "b:0xb3,c:0xc4"
	case strings.Contains(operand, "DE"):
		return "d:0xd5,e:0xe6"
	case strings.Contains(operand, "HL"):
		return "h:0xa7,l:0xf8"
	case strings.Contains(operand, "A"):
		return "a:0x1a"
	case strings.Contains(operand, "B"):
		return "b:0x2b"
	case strings.Contains(operand, "C"):
		return "c:0x3c"
	case strings.Contains(operand, "D"):
		return "d:0x4d"
	case strings.Contains(operand, "E"):
		return "e:0x5e"
	case strings.Contains(operand, "F"):
		return "f:0x6f"
	case strings.Contains(operand, "H"):
		return "h:0x71"
	case strings.Contains(operand, "L"):
		return "l:0x82"
	}
	return ""
}

func generateCPU(opcode opcodeMetadata) string {
	op1 := generateCPUByOperand(opcode.Operand1)
	op2 := generateCPUByOperand(opcode.Operand2)
	if op2 == "" {
		if op1 == "" {
			return ""
		}
		return op1
	}
	if op1 == "" {
		return op2
	}
	if op1 == op2 {
		return op1
	}
	return op1 + "," + op2
}

func generateMemory(opcode opcodeMetadata) string {
	switch {
	case strings.Contains(opcode.Operand1, "(BC") || strings.Contains(opcode.Operand2, "(BC"):
		return "testMemory{mem: map[uint16]byte{0xb3c4: 0xaa}}"
	case strings.Contains(opcode.Operand1, "(DE") || strings.Contains(opcode.Operand2, "(DE"):
		return "testMemory{mem: map[uint16]byte{0xd5e6: 0xbb}}"
	case strings.Contains(opcode.Operand1, "(HL") || strings.Contains(opcode.Operand2, "(HL"):
		return "testMemory{mem: map[uint16]byte{0xa7f8: 0xcc}}"
	case strings.Contains(opcode.Operand1, "(A") || strings.Contains(opcode.Operand2, "(A"):
		return "testMemory{mem: map[uint16]byte{0xff1a: 0x22}}"
	case strings.Contains(opcode.Operand1, "(B") || strings.Contains(opcode.Operand2, "(B"):
		return "testMemory{mem: map[uint16]byte{0xff2b: 0x33}}"
	case strings.Contains(opcode.Operand1, "(C") || strings.Contains(opcode.Operand2, "(C"):
		return "testMemory{mem: map[uint16]byte{0xff3c: 0x44}}"
	case strings.Contains(opcode.Operand1, "(D") || strings.Contains(opcode.Operand2, "(D"):
		return "testMemory{mem: map[uint16]byte{0xff4d: 0x55}}"
	case strings.Contains(opcode.Operand1, "(E") || strings.Contains(opcode.Operand2, "(E"):
		return "testMemory{mem: map[uint16]byte{0xff5e: 0x66}}"
	case strings.Contains(opcode.Operand1, "(F") || strings.Contains(opcode.Operand2, "(F"):
		return "testMemory{mem: map[uint16]byte{0xff6f: 0x77}}"
	case strings.Contains(opcode.Operand1, "(H") || strings.Contains(opcode.Operand2, "(H"):
		return "testMemory{mem: map[uint16]byte{0xff71: 0x88}}"
	case strings.Contains(opcode.Operand1, "(L") || strings.Contains(opcode.Operand2, "(L"):
		return "testMemory{mem: map[uint16]byte{0xff82: 0x99}}"
	}
	return "nil"
}

func generateFlags(opcode opcodeMetadata) []string {
	flags := []string{}
	if opcode.Flags == nil {
		flags = append(flags, "")
		return flags
	}
	z := opcode.Flags[0] == "Z"
	n := opcode.Flags[1] == "N"
	h := opcode.Flags[2] == "H"
	c := opcode.Flags[3] == "C"
	switch {
	case z && n && h && c:
		flags = append(flags, `"Z":false,"N":false,"H":false,"C":false`)
		flags = append(flags, `"Z":true,"N":false,"H":false,"C":false`)
		flags = append(flags, `"Z":false,"N":true,"H":false,"C":false`)
		flags = append(flags, `"Z":false,"N":false,"H":true,"C":false`)
		flags = append(flags, `"Z":false,"N":false,"H":false,"C":true`)
	case z && n && h:
		flags = append(flags, `"Z":false,"N":false,"H":false`)
		flags = append(flags, `"Z":true,"N":false,"H":false`)
		flags = append(flags, `"Z":false,"N":true,"H":false`)
		flags = append(flags, `"Z":false,"N":false,"H":true`)
	case z && n && c:
		flags = append(flags, `"Z":false,"N":false,"C":false`)
		flags = append(flags, `"Z":true,"N":false,"C":false`)
		flags = append(flags, `"Z":false,"N":true,"C":false`)
		flags = append(flags, `"Z":false,"N":false,"C":true`)
	case z && h && c:
		flags = append(flags, `"Z":false,"H":false,"C":false`)
		flags = append(flags, `"Z":true,"H":false,"C":false`)
		flags = append(flags, `"Z":false,"H":true,"C":false`)
		flags = append(flags, `"Z":false,"H":false,"C":true`)
	case n && h && c:
		flags = append(flags, `"N":false,"H":false,"C":false`)
		flags = append(flags, `"N":true,"H":false,"C":false`)
		flags = append(flags, `"N":false,"H":true,"C":false`)
		flags = append(flags, `"N":false,"H":false,"C":true`)
	case z && n:
		flags = append(flags, `"Z":false,"N":false`)
		flags = append(flags, `"Z":true,"N":false`)
		flags = append(flags, `"Z":false,"N":true`)
	case z && h:
		flags = append(flags, `"Z":false,"H":false`)
		flags = append(flags, `"Z":true,"H":false`)
		flags = append(flags, `"Z":false,"H":true`)
	case z && c:
		flags = append(flags, `"Z":false,"C":false`)
		flags = append(flags, `"Z":true,"C":false`)
		flags = append(flags, `"Z":false,"C":true`)
	case n && h:
		flags = append(flags, `"N":false,"H":false`)
		flags = append(flags, `"N":true,"H":false`)
		flags = append(flags, `"N":false,"H":true`)
	case n && c:
		flags = append(flags, `"N":false,"C":false`)
		flags = append(flags, `"N":true,"C":false`)
		flags = append(flags, `"N":false,"C":true`)
	case h && c:
		flags = append(flags, `"H":false,"C":false`)
		flags = append(flags, `"H":true,"C":false`)
		flags = append(flags, `"H":false,"C":true`)
	case z:
		flags = append(flags, `"Z":false`)
		flags = append(flags, `"Z":true`)
	case n:
		flags = append(flags, `"N":false`)
		flags = append(flags, `"N":true`)
	case h:
		flags = append(flags, `"H":false`)
		flags = append(flags, `"H":true`)
	case c:
		flags = append(flags, `"C":false`)
		flags = append(flags, `"C":true`)
	}
	if len(flags) == 0 {
		flags = append(flags, "")
	}
	return flags
}

func generateOpcode(opcode opcodeMetadata) {
	if opcode.Mnemonic == "" {
		// No opcode for this array index
		return
	}
	cpu := generateCPU(opcode)
	memory := generateMemory(opcode)
	flags := generateFlags(opcode)
	if len(flags) == 0 {
		panic("Need flags")
	}
	for _, f := range flags {
		fmt.Printf(format+"\n", opcode.Mnemonic, opcode.Operand1, opcode.Operand2, 0, 0, cpu, cpu, f, memory)
	}
}

// GenerateTestData generates Go source for use in opcode unit tests
func GenerateTestData() {
	for _, opcode := range opcodes {
		generateOpcode(opcode)
	}
	for _, opcode := range prefixedOpcodes {
		generateOpcode(opcode)
	}
}
