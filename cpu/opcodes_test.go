package cpu

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDispatchOneByteInstruction(t *testing.T) {
	for _, test := range dispatchOneByteInstructionTests {
		initialCPU := test.actualCPU
		(&test.actualCPU).dispatchOneByteInstruction(test.mem, test.instruction)
		compareCPUs(t, opcodes[test.instruction], initialCPU, test.expectedCPU, test.actualCPU, test.mem)
	}
}

func TestDispatchTwoByteInstruction(t *testing.T) {
	for _, test := range dispatchTwoByteInstructionTests {
		initialCPU := test.actualCPU
		(&test.actualCPU).dispatchTwoByteInstruction(test.mem, test.instruction, test.u8)
		compareCPUs(t, opcodes[test.instruction], initialCPU, test.expectedCPU, test.actualCPU, test.mem)
	}
}

func TestDispatchThreeByteInstruction(t *testing.T) {
	for _, test := range dispatchThreeByteInstructionTests {
		initialCPU := test.actualCPU
		(&test.actualCPU).dispatchThreeByteInstruction(test.mem, test.instruction, test.u16)
		compareCPUs(t, opcodes[test.instruction], initialCPU, test.expectedCPU, test.actualCPU, test.mem)
	}
}

func TestDispatchPrefixedInstruction(t *testing.T) {
	for _, test := range dispatchPrefixedInstructionTests {
		initialCPU := test.actualCPU
		(&test.actualCPU).dispatchPrefixedInstruction(test.mem, test.instruction)
		compareCPUs(t, prefixedOpcodes[test.instruction], initialCPU, test.expectedCPU, test.actualCPU, test.mem)
	}
}

func validateFlags(t *testing.T, opcode opcodeMetadata, actualCPU CPU) {
	if len(opcode.Flags) > 0 {
		actualFlags := []bool{actualCPU.zf, actualCPU.nf, actualCPU.hf, actualCPU.cf}
		for i, flag := range opcode.Flags {
			switch flag {
			case "0":
				if actualFlags[i] {
					t.Error("Flags do not match ", opcode.Flags)
					t.Error("  Actual   : ", actualCPU)
				}
			case "1":
				if !actualFlags[i] {
					t.Error("Flags do not match ", opcode.Flags)
					t.Error("  Actual   : ", actualCPU)
				}
			}
		}
	}
}

func validateMemory(t *testing.T, opcode opcodeMetadata, mem testableMemory) {
	if mem.expected != nil && !reflect.DeepEqual(mem.actual, mem.expected) {
		t.Error("Memory does not match for: ", opcode)
		t.Error("  Expected : ", mem.expected)
		t.Error("  Actual   : ", mem.actual)
	}
}

func compareCPUs(t *testing.T, opcode opcodeMetadata, initialCPU, expectedCPU, actualCPU CPU, mem *testableMemory) {
	if actualCPU != expectedCPU {
		t.Error("CPUs do not match for: ", opcode)
		t.Error("  Initial  : ", initialCPU)
		t.Error("  Expected : ", expectedCPU)
		t.Error("  Actual   : ", actualCPU)
	}
	validateFlags(t, opcode, actualCPU)
	if mem != nil {
		validateMemory(t, opcode, *mem)
	}
}

type dispatchOneByteInstructionTestTable struct {
	instruction uint8
	actualCPU   CPU
	expectedCPU CPU
	mem         *testableMemory
}

type dispatchTwoByteInstructionTestTable struct {
	instruction uint8
	u8          uint8
	actualCPU   CPU
	expectedCPU CPU
	mem         *testableMemory
}

type dispatchThreeByteInstructionTestTable struct {
	instruction uint8
	u16         uint16
	actualCPU   CPU
	expectedCPU CPU
	mem         *testableMemory
}

type dispatchPrefixedInstructionTestTable struct {
	instruction uint8
	actualCPU   CPU
	expectedCPU CPU
	mem         *testableMemory
}

type testableMemory struct {
	actual   map[uint16]byte
	expected map[uint16]byte
}

// Read a byte from the chosen memory location
func (mem testableMemory) Read(addr uint16) byte {
	return mem.actual[addr]
}

// Write a byte to the chosen memory location
func (mem testableMemory) Write(addr uint16, b byte) {
	if mem.expected == nil {
		panic("Writes are not permitted if there is no 'expected' memory")
	}
	mem.actual[addr] = b
}

func (mem testableMemory) GenerateCrashReport() {
	fmt.Println("TestMemory crash: ", mem.actual, mem.expected)
}

var dispatchOneByteInstructionTests = []dispatchOneByteInstructionTestTable{

	// ADC A A [Z 0 H C]
	// {0x8f,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
	// {0x8f,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

	// ADC A B [Z 0 H C]
	// {0x88,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
	// {0x88,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

	// ADC A C [Z 0 H C]
	// {0x89,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
	// {0x89,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

	// ADC A D [Z 0 H C]
	// {0x8a,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
	// {0x8a,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

	// ADC A d8 [Z 0 H C]
	// {0xce,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
	// {0xce,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

	// ADC A E [Z 0 H C]
	// {0x8b,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
	// {0x8b,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

	// ADC A H [Z 0 H C]
	// {0x8c,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
	// {0x8c,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

	// ADC A (HL) [Z 0 H C]
	// {0x8e,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
	// {0x8e,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

	// ADC A L [Z 0 H C]
	// {0x8d,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
	// {0x8d,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

	// ADD A A [Z 0 H C]
	{0x87, CPU{a: 0x12}, CPU{a: 0x24}, nil},
	{0x87, CPU{a: 0xa3}, CPU{a: 0x46, cf: true}, nil},
	{0x87, CPU{a: 0x1a}, CPU{a: 0x34, hf: true}, nil},
	{0x87, CPU{a: 0x00}, CPU{a: 0x00, zf: true}, nil},

	// ADD A B [Z 0 H C]
	{0x80, CPU{a: 0x1a, b: 0x22}, CPU{a: 0x3c, b: 0x22}, nil},
	{0x80, CPU{a: 0x1a, b: 0xf2}, CPU{a: 0x0c, b: 0xf2, cf: true}, nil},
	{0x80, CPU{a: 0x1a, b: 0x2b}, CPU{a: 0x45, b: 0x2b, hf: true}, nil},
	{0x80, CPU{a: 0x00, b: 0x00}, CPU{a: 0x00, b: 0x00, zf: true}, nil}}

// ADD A C [Z 0 H C]
// {0x81,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x81,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// ADD A D [Z 0 H C]
// {0x82,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x82,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// ADD A E [Z 0 H C]
// {0x83,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x83,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// ADD A H [Z 0 H C]
// {0x84,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x84,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// ADD A L [Z 0 H C]
// {0x85,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x85,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// ADD A (HL) [Z 0 H C]
// {0x86,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x86,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// ADD A d8 [Z 0 H C]
// {0xc6,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc6,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// ADD HL BC [- 0 H C]
// {0x09,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x09,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// ADD HL DE [- 0 H C]
// {0x19,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x19,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// ADD HL HL [- 0 H C]
// {0x29,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x29,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// ADD HL SP [- 0 H C]
// {0x39,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x39,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// ADD SP r8 [0 0 H C]
// {0xe8,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe8,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// AND A  [Z 0 1 0]
// {0xa7,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa7,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// AND B  [Z 0 1 0]
// {0xa0,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa0,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// AND C  [Z 0 1 0]
// {0xa1,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa1,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// AND D  [Z 0 1 0]
// {0xa2,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa2,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// AND d8  [Z 0 1 0]
// {0xe6,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe6,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// AND E  [Z 0 1 0]
// {0xa3,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa3,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// AND H  [Z 0 1 0]
// {0xa4,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa4,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// AND (HL)  [Z 0 1 0]
// {0xa6,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa6,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// AND L  [Z 0 1 0]
// {0xa5,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa5,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// CALL a16  []
// {0xcd,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xcd,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// CALL C a16 []
// {0xdc,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xdc,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// CALL NC a16 []
// {0xd4,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd4,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// CALL NZ a16 []
// {0xc4,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc4,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// CALL Z a16 []
// {0xcc,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xcc,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// CCF   [- 0 0 C]
// {0x3f,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x3f,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// CP A  [Z 1 H C]
// {0xbf,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xbf,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// CP B  [Z 1 H C]
// {0xb8,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb8,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// CP C  [Z 1 H C]
// {0xb9,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb9,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// CP D  [Z 1 H C]
// {0xba,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xba,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// CP d8  [Z 1 H C]
// {0xfe,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xfe,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// CP E  [Z 1 H C]
// {0xbb,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xbb,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// CP H  [Z 1 H C]
// {0xbc,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xbc,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// CP (HL)  [Z 1 H C]
// {0xbe,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xbe,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// CP L  [Z 1 H C]
// {0xbd,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xbd,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// CPL   [- 1 1 -]
// {0x2f,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x2f,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// DAA   [Z - 0 C]
// {0x27,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x27,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// DEC A  [Z 1 H -]
// {0x3d,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x3d,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// DEC B  [Z 1 H -]
// {0x05,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x05,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// DEC BC  []
// {0x0b,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0b,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// DEC C  [Z 1 H -]
// {0x0d,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0d,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// DEC D  [Z 1 H -]
// {0x15,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x15,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// DEC DE  []
// {0x1b,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1b,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// DEC E  [Z 1 H -]
// {0x1d,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1d,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// DEC H  [Z 1 H -]
// {0x25,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x25,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// DEC (HL)  [Z 1 H -]
// {0x35,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x35,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// DEC HL  []
// {0x2b,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x2b,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// DEC L  [Z 1 H -]
// {0x2d,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x2d,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// DEC SP  []
// {0x3b,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x3b,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// DI   []
// {0xf3,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf3,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// EI   []
// {0xfb,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xfb,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// HALT   []
// {0x76,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x76,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// INC A  [Z 0 H -]
// {0x3c,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x3c,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// INC B  [Z 0 H -]
// {0x04,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x04,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// INC BC  []
// {0x03,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x03,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// INC C  [Z 0 H -]
// {0x0c,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0c,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// INC D  [Z 0 H -]
// {0x14,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x14,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// INC DE  []
// {0x13,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x13,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// INC E  [Z 0 H -]
// {0x1c,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1c,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// INC H  [Z 0 H -]
// {0x24,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x24,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// INC (HL)  [Z 0 H -]
// {0x34,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x34,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// INC HL  []
// {0x23,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x23,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// INC L  [Z 0 H -]
// {0x2c,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x2c,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// INC SP  []
// {0x33,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x33,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// JP a16  []
// {0xc3,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc3,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// JP C a16 []
// {0xda,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xda,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// JP NC a16 []
// {0xd2,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd2,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// JP Z a16 []
// {0xca,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xca,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// JP NZ a16 []
// {0xc2,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc2,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// JP (HL)  []
// {0xe9,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe9,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// JR C r8 []
// {0x38,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x38,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// JR NC r8 []
// {0x30,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x30,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// JR NZ r8 []
// {0x20,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x20,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// JR r8  []
// {0x18,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x18,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// JR Z r8 []
// {0x28,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x28,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD (a16) A []
// {0xea,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xea,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD (a16) SP []
// {0x08,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x08,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD (C) A []
// {0xe2,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe2,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD (DE) A []
// {0x12,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x12,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // case 0x32 :  cpu. ld(cpu.(hl-), cpu.a) // LD (HL-) A []
// // case 0x22 :  cpu. ld(cpu.(hl+), cpu.a) // LD (HL+) A []
// LD A (a16) []
// {0xfa,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xfa,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD A (C) []
// {0xf2,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf2,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD A (DE) []
// {0x1a,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1a,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // case 0x3a :  cpu. ld(cpu.a, cpu.(hl-)) // LD A (HL-) []
// // case 0x2a :  cpu. ld(cpu.a, cpu.(hl+)) // LD A (HL+) []
// LD A A []
// {0x7f,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x7f,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD A B []
// {0x78,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x78,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD A (BC) []
// {0x0a,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0a,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD A C []
// {0x79,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x79,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD A D []
// {0x7a,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x7a,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD A d8 []
// {0x3e,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x3e,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD A E []
// {0x7b,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x7b,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD A H []
// {0x7c,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x7c,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD A (HL) []
// {0x7e,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x7e,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD A L []
// {0x7d,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x7d,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD B A []
// {0x47,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x47,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD B B []
// {0x40,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x40,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD B C []
// {0x41,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x41,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD B D []
// {0x42,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x42,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD B d8 []
// {0x06,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x06,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD B E []
// {0x43,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x43,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD B H []
// {0x44,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x44,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD B (HL) []
// {0x46,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x46,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD B L []
// {0x45,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x45,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD (BC) A []
// {0x02,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x02,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD BC d16 []
// {0x01,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x01,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD C A []
// {0x4f,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x4f,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD C B []
// {0x48,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x48,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD C C []
// {0x49,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x49,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD C D []
// {0x4a,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x4a,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD C d8 []
// {0x0e,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0e,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD C E []
// {0x4b,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x4b,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD C H []
// {0x4c,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x4c,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD C (HL) []
// {0x4e,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x4e,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD C L []
// {0x4d,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x4d,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD D A []
// {0x57,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x57,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD D B []
// {0x50,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x50,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD D C []
// {0x51,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x51,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD D D []
// {0x52,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x52,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD D d8 []
// {0x16,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x16,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD D E []
// {0x53,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x53,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD D H []
// {0x54,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x54,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD D (HL) []
// {0x56,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x56,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD D L []
// {0x55,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x55,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD DE d16 []
// {0x11,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x11,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD E A []
// {0x5f,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x5f,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD E B []
// {0x58,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x58,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD E C []
// {0x59,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x59,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD E D []
// {0x5a,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x5a,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD E d8 []
// {0x1e,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1e,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD E E []
// {0x5b,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x5b,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD E H []
// {0x5c,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x5c,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD E (HL) []
// {0x5e,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x5e,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD E L []
// {0x5d,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x5d,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD H A []
// {0x67,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x67,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD H B []
// {0x60,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x60,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD H C []
// {0x61,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x61,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD H D []
// {0x62,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x62,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD H d8 []
// {0x26,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x26,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD H E []
// {0x63,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x63,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD H H []
// {0x64,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x64,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD H (HL) []
// {0x66,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x66,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD H L []
// {0x65,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x65,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD (HL) A []
// {0x77,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x77,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD (HL) B []
// {0x70,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x70,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD (HL) C []
// {0x71,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x71,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD (HL) D []
// {0x72,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x72,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD HL d16 []
// {0x21,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x21,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD (HL) d8 []
// {0x36,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x36,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD (HL) E []
// {0x73,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x73,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD (HL) H []
// {0x74,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x74,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD (HL) L []
// {0x75,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x75,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD HL SP+r8 [0 0 H C]
// {0xf8,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf8,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD L A []
// {0x6f,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x6f,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD L B []
// {0x68,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x68,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD L C []
// {0x69,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x69,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD L D []
// {0x6a,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x6a,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD L d8 []
// {0x2e,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x2e,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD L E []
// {0x6b,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x6b,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD L H []
// {0x6c,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x6c,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD L (HL) []
// {0x6e,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x6e,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD L L []
// {0x6d,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x6d,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD SP d16 []
// {0x31,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x31,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LD SP HL []
// {0xf9,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf9,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LDH (a8) A []
// {0xe0,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe0,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// LDH A (a8) []
// {0xf0,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf0,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// NOP   []
// {0x00,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x00,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// OR A  [Z 0 0 0]
// {0xb7,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb7,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// OR B  [Z 0 0 0]
// {0xb0,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb0,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// OR C  [Z 0 0 0]
// {0xb1,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb1,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// OR D  [Z 0 0 0]
// {0xb2,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb2,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// OR d8  [Z 0 0 0]
// {0xf6,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf6,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// OR E  [Z 0 0 0]
// {0xb3,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb3,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// OR H  [Z 0 0 0]
// {0xb4,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb4,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// OR (HL)  [Z 0 0 0]
// {0xb6,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb6,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// OR L  [Z 0 0 0]
// {0xb5,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb5,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// POP AF  [Z N H C]
// {0xf1,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf1,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// POP BC  []
// {0xc1,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc1,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// POP DE  []
// {0xd1,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd1,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// POP HL  []
// {0xe1,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe1,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// PREFIX CB  []
// {0xcb,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xcb,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// PUSH AF  []
// {0xf5,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf5,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// PUSH BC  []
// {0xc5,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc5,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// PUSH DE  []
// {0xd5,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd5,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// PUSH HL  []
// {0xe5,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe5,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RET   []
// {0xc9,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc9,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RET C  []
// {0xd8,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd8,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RET NC  []
// {0xd0,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd0,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RET NZ  []
// {0xc0,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc0,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RET Z  []
// {0xc8,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc8,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RETI   []
// {0xd9,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd9,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RLA   [0 0 0 C]
// {0x17,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x17,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RLCA   [0 0 0 C]
// {0x07,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x07,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RRA   [0 0 0 C]
// {0x1f,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1f,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RRCA   [0 0 0 C]
// {0x0f,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0f,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// FIXME
// // case 0xc7 :  cpu. rst(cpu.00h) // RST 00H  []
// // case 0xcf :  cpu. rst(cpu.08h) // RST 08H  []
// // case 0xd7 :  cpu. rst(cpu.10h) // RST 10H  []
// // case 0xdf :  cpu. rst(cpu.18h) // RST 18H  []
// // case 0xe7 :  cpu. rst(cpu.20h) // RST 20H  []
// // case 0xef :  cpu. rst(cpu.28h) // RST 28H  []
// // case 0xf7 :  cpu. rst(cpu.30h) // RST 30H  []
// // case 0xff :  cpu. rst(cpu.38h) // RST 38H  []

// SBC A A [Z 1 H C]
// {0x9f,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9f,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SBC A B [Z 1 H C]
// {0x98,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x98,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SBC A C [Z 1 H C]
// {0x99,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x99,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SBC A D [Z 1 H C]
// {0x9a,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9a,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SBC A d8 [Z 1 H C]
// {0xde,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xde,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SBC A E [Z 1 H C]
// {0x9b,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9b,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SBC A H [Z 1 H C]
// {0x9c,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9c,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SBC A (HL) [Z 1 H C]
// {0x9e,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9e,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SBC A L [Z 1 H C]
// {0x9d,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9d,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SCF   [- 0 0 1]
// {0x37,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x37,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// STOP 0  []
// {0x10,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x10,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SUB A  [Z 1 H C]
// {0x97,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x97,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SUB B  [Z 1 H C]
// {0x90,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x90,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SUB C  [Z 1 H C]
// {0x91,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x91,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SUB D  [Z 1 H C]
// {0x92,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x92,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SUB d8  [Z 1 H C]
// {0xd6,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd6,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SUB E  [Z 1 H C]
// {0x93,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x93,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SUB H  [Z 1 H C]
// {0x94,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x94,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SUB (HL)  [Z 1 H C]
// {0x96,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x96,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SUB L  [Z 1 H C]
// {0x95,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x95,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// XOR A  [Z 0 0 0]
// {0xaf,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xaf,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// XOR B  [Z 0 0 0]
// {0xa8,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa8,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// XOR C  [Z 0 0 0]
// {0xa9,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa9,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// XOR D  [Z 0 0 0]
// {0xaa,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xaa,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// XOR d8  [Z 0 0 0]
// {0xee,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xee,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// XOR E  [Z 0 0 0]
// {0xab,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xab,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// XOR H  [Z 0 0 0]
// {0xac,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xac,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// XOR (HL)  [Z 0 0 0]
// {0xae,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xae,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// XOR L  [Z 0 0 0]
// {0xad,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xad,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 0 A []
// {0x87,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x87,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 0 B []
// {0x80,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x80,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 0 C []
// {0x81,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x81,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 0 D []
// {0x82,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x82,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 0 E []
// {0x83,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x83,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 0 H []
// {0x84,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x84,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 0 (HL) []
// {0x86,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x86,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 0 L []
// {0x85,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x85,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 1 A []
// {0x8f,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x8f,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 1 B []
// {0x88,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x88,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 1 C []
// {0x89,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x89,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 1 D []
// {0x8a,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x8a,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 1 E []
// {0x8b,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x8b,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 1 H []
// {0x8c,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x8c,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 1 (HL) []
// {0x8e,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x8e,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 1 L []
// {0x8d,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x8d,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 2 A []
// {0x97,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x97,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 2 B []
// {0x90,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x90,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 2 C []
// {0x91,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x91,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 2 D []
// {0x92,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x92,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 2 E []
// {0x93,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x93,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 2 H []
// {0x94,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x94,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 2 (HL) []
// {0x96,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x96,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 2 L []
// {0x95,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x95,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 3 A []
// {0x9f,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9f,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 3 B []
// {0x98,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x98,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 3 C []
// {0x99,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x99,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 3 D []
// {0x9a,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9a,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 3 E []
// {0x9b,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9b,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 3 H []
// {0x9c,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9c,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 3 (HL) []
// {0x9e,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9e,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 3 L []
// {0x9d,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9d,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 4 A []
// {0xa7,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa7,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 4 B []
// {0xa0,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa0,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 4 C []
// {0xa1,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa1,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 4 D []
// {0xa2,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa2,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 4 E []
// {0xa3,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa3,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 4 H []
// {0xa4,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa4,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 4 (HL) []
// {0xa6,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa6,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 4 L []
// {0xa5,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa5,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 5 A []
// {0xaf,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xaf,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 5 B []
// {0xa8,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa8,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 5 C []
// {0xa9,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa9,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 5 D []
// {0xaa,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xaa,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 5 E []
// {0xab,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xab,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 5 H []
// {0xac,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xac,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 5 (HL) []
// {0xae,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xae,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 5 L []
// {0xad,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xad,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 6 A []
// {0xb7,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb7,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 6 B []
// {0xb0,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb0,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 6 C []
// {0xb1,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb1,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 6 D []
// {0xb2,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb2,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 6 E []
// {0xb3,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb3,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 6 H []
// {0xb4,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb4,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 6 (HL) []
// {0xb6,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb6,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 6 L []
// {0xb5,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb5,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 7 A []
// {0xbf,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xbf,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 7 B []
// {0xb8,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb8,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 7 C []
// {0xb9,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb9,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 7 D []
// {0xba,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xba,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 7 E []
// {0xbb,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xbb,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 7 H []
// {0xbc,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xbc,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 7 (HL) []
// {0xbe,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xbe,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RES 7 L []
// {0xbd,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xbd,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RL A  [Z 0 0 C]
// {0x17,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x17,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RL B  [Z 0 0 C]
// {0x10,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x10,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RL C  [Z 0 0 C]
// {0x11,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x11,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RL D  [Z 0 0 C]
// {0x12,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x12,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RL E  [Z 0 0 C]
// {0x13,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x13,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RL H  [Z 0 0 C]
// {0x14,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x14,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RL (HL)  [Z 0 0 C]
// {0x16,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x16,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RL L  [Z 0 0 C]
// {0x15,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x15,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RLC A  [Z 0 0 C]
// {0x07,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x07,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RLC B  [Z 0 0 C]
// {0x00,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x00,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RLC C  [Z 0 0 C]
// {0x01,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x01,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RLC D  [Z 0 0 C]
// {0x02,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x02,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RLC E  [Z 0 0 C]
// {0x03,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x03,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RLC H  [Z 0 0 C]
// {0x04,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x04,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RLC (HL)  [Z 0 0 C]
// {0x06,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x06,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RLC L  [Z 0 0 C]
// {0x05,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x05,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RR A  [Z 0 0 C]
// {0x1f,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1f,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RR B  [Z 0 0 C]
// {0x18,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x18,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RR C  [Z 0 0 C]
// {0x19,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x19,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RR D  [Z 0 0 C]
// {0x1a,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1a,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RR E  [Z 0 0 C]
// {0x1b,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1b,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RR H  [Z 0 0 C]
// {0x1c,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1c,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RR (HL)  [Z 0 0 C]
// {0x1e,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1e,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RR L  [Z 0 0 C]
// {0x1d,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1d,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RRC A  [Z 0 0 C]
// {0x0f,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0f,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RRC B  [Z 0 0 C]
// {0x08,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x08,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RRC C  [Z 0 0 C]
// {0x09,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x09,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RRC D  [Z 0 0 C]
// {0x0a,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0a,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RRC E  [Z 0 0 C]
// {0x0b,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0b,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RRC H  [Z 0 0 C]
// {0x0c,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0c,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RRC (HL)  [Z 0 0 C]
// {0x0e,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0e,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// RRC L  [Z 0 0 C]
// {0x0d,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0d,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 0 A []
// {0xc7,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc7,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 0 B []
// {0xc0,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc0,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 0 C []
// {0xc1,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc1,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 0 D []
// {0xc2,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc2,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 0 E []
// {0xc3,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc3,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 0 H []
// {0xc4,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc4,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 0 (HL) []
// {0xc6,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc6,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 0 L []
// {0xc5,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc5,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 1 A []
// {0xcf,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xcf,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 1 B []
// {0xc8,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc8,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 1 C []
// {0xc9,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc9,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 1 D []
// {0xca,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xca,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 1 E []
// {0xcb,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xcb,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 1 H []
// {0xcc,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xcc,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 1 (HL) []
// {0xce,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xce,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 1 L []
// {0xcd,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xcd,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 2 A []
// {0xd7,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd7,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 2 B []
// {0xd0,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd0,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 2 C []
// {0xd1,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd1,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 2 D []
// {0xd2,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd2,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 2 E []
// {0xd3,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd3,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 2 H []
// {0xd4,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd4,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 2 (HL) []
// {0xd6,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd6,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 2 L []
// {0xd5,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd5,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 3 A []
// {0xdf,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xdf,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 3 B []
// {0xd8,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd8,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 3 C []
// {0xd9,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd9,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 3 D []
// {0xda,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xda,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 3 E []
// {0xdb,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xdb,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 3 H []
// {0xdc,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xdc,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 3 (HL) []
// {0xde,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xde,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 3 L []
// {0xdd,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xdd,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 4 A []
// {0xe7,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe7,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 4 B []
// {0xe0,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe0,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 4 C []
// {0xe1,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe1,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 4 D []
// {0xe2,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe2,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 4 E []
// {0xe3,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe3,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 4 H []
// {0xe4,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe4,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 4 (HL) []
// {0xe6,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe6,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 4 L []
// {0xe5,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe5,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 5 A []
// {0xef,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xef,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 5 B []
// {0xe8,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe8,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 5 C []
// {0xe9,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe9,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 5 D []
// {0xea,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xea,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 5 E []
// {0xeb,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xeb,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 5 H []
// {0xec,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xec,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 5 (HL) []
// {0xee,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xee,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 5 L []
// {0xed,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xed,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 6 A []
// {0xf7,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf7,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 6 B []
// {0xf0,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf0,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 6 C []
// {0xf1,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf1,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 6 D []
// {0xf2,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf2,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 6 E []
// {0xf3,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf3,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 6 H []
// {0xf4,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf4,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 6 (HL) []
// {0xf6,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf6,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 6 L []
// {0xf5,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf5,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 7 A []
// {0xff,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xff,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 7 B []
// {0xf8,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf8,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 7 C []
// {0xf9,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf9,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 7 D []
// {0xfa,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xfa,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 7 E []
// {0xfb,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xfb,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 7 H []
// {0xfc,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xfc,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 7 (HL) []
// {0xfe,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xfe,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SET 7 L []
// {0xfd,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xfd,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SLA A  [Z 0 0 C]
// {0x27,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x27,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SLA B  [Z 0 0 C]
// {0x20,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x20,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SLA C  [Z 0 0 C]
// {0x21,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x21,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SLA D  [Z 0 0 C]
// {0x22,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x22,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SLA E  [Z 0 0 C]
// {0x23,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x23,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SLA H  [Z 0 0 C]
// {0x24,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x24,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SLA (HL)  [Z 0 0 C]
// {0x26,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x26,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SLA L  [Z 0 0 C]
// {0x25,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x25,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SRA A  [Z 0 0 0]
// {0x2f,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x2f,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SRA B  [Z 0 0 0]
// {0x28,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x28,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SRA C  [Z 0 0 0]
// {0x29,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x29,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SRA D  [Z 0 0 0]
// {0x2a,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x2a,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SRA E  [Z 0 0 0]
// {0x2b,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x2b,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SRA H  [Z 0 0 0]
// {0x2c,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x2c,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SRA (HL)  [Z 0 0 0]
// {0x2e,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x2e,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SRA L  [Z 0 0 0]
// {0x2d,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x2d,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SRL A  [Z 0 0 C]
// {0x3f,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x3f,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SRL B  [Z 0 0 C]
// {0x38,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x38,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SRL C  [Z 0 0 C]
// {0x39,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x39,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SRL D  [Z 0 0 C]
// {0x3a,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x3a,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SRL E  [Z 0 0 C]
// {0x3b,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x3b,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SRL H  [Z 0 0 C]
// {0x3c,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x3c,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SRL (HL)  [Z 0 0 C]
// {0x3e,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x3e,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SRL L  [Z 0 0 C]
// {0x3d,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x3d,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SWAP A  [Z 0 0 0]
// {0x37,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x37,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SWAP B  [Z 0 0 0]
// {0x30,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x30,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SWAP C  [Z 0 0 0]
// {0x31,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x31,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SWAP D  [Z 0 0 0]
// {0x32,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x32,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SWAP E  [Z 0 0 0]
// {0x33,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x33,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SWAP H  [Z 0 0 0]
// {0x34,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x34,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SWAP (HL)  [Z 0 0 0]
// {0x36,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x36,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// SWAP L  [Z 0 0 0]
// {0x35,  CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x35,  CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil}}

// {"ADD", "A", "C",  CPU{a: 0x1a, c: 0x22}, CPU{a: 0x3c, c: 0x22}, map[string]bool{"Z": false, "H": false, "C": false}, nil},
// {"ADD", "A", "C",  CPU{a: 0x1a, c: 0xf2}, CPU{a: 0x0c, c: 0xf2}, map[string]bool{"Z": false, "H": false, "C": true}, nil},
// {"ADD", "A", "C",  CPU{a: 0x1a, c: 0x2b}, CPU{a: 0x45, c: 0x2b}, map[string]bool{"Z": false, "H": true, "C": false}, nil},
// {"ADD", "A", "C",  CPU{a: 0x00, c: 0x00}, CPU{a: 0x00, c: 0x00}, map[string]bool{"Z": true, "H": false, "C": false}, nil},
// {"ADD", "A", "D",  CPU{a: 0x1a, d: 0x22}, CPU{a: 0x3c, d: 0x22}, map[string]bool{"Z": false, "H": false, "C": false}, nil},
// {"ADD", "A", "D",  CPU{a: 0x1a, d: 0xf2}, CPU{a: 0x0c, d: 0xf2}, map[string]bool{"Z": false, "H": false, "C": true}, nil},
// {"ADD", "A", "D",  CPU{a: 0x1a, d: 0x2b}, CPU{a: 0x45, d: 0x2b}, map[string]bool{"Z": false, "H": true, "C": false}, nil},
// {"ADD", "A", "D",  CPU{a: 0x00, d: 0x00}, CPU{a: 0x00, d: 0x00}, map[string]bool{"Z": true, "H": false, "C": false}, nil},
// {"ADD", "A", "E",  CPU{a: 0x1a, e: 0x22}, CPU{a: 0x3c, e: 0x22}, map[string]bool{"Z": false, "H": false, "C": false}, nil},
// {"ADD", "A", "E",  CPU{a: 0x1a, e: 0xf2}, CPU{a: 0x0c, e: 0xf2}, map[string]bool{"Z": false, "H": false, "C": true}, nil},
// {"ADD", "A", "E",  CPU{a: 0x1a, e: 0x2b}, CPU{a: 0x45, e: 0x2b}, map[string]bool{"Z": false, "H": true, "C": false}, nil},
// {"ADD", "A", "E",  CPU{a: 0x00, e: 0x00}, CPU{a: 0x00, e: 0x00}, map[string]bool{"Z": true, "H": false, "C": false}, nil},
// {"ADD", "A", "H",  CPU{a: 0x1a, h: 0x22}, CPU{a: 0x3c, h: 0x22}, map[string]bool{"Z": false, "H": false, "C": false}, nil},
// {"ADD", "A", "H",  CPU{a: 0x1a, h: 0xf2}, CPU{a: 0x0c, h: 0xf2}, map[string]bool{"Z": false, "H": false, "C": true}, nil},
// {"ADD", "A", "H",  CPU{a: 0x1a, h: 0x2b}, CPU{a: 0x45, h: 0x2b}, map[string]bool{"Z": false, "H": true, "C": false}, nil},
// {"ADD", "A", "H",  CPU{a: 0x00, h: 0x00}, CPU{a: 0x00, h: 0x00}, map[string]bool{"Z": true, "H": false, "C": false}, nil},
// {"ADD", "A", "L",  CPU{a: 0x1a, l: 0x22}, CPU{a: 0x3c, l: 0x22}, map[string]bool{"Z": false, "H": false, "C": false}, nil},
// {"ADD", "A", "L",  CPU{a: 0x1a, l: 0xf2}, CPU{a: 0x0c, l: 0xf2}, map[string]bool{"Z": false, "H": false, "C": true}, nil},
// {"ADD", "A", "L",  CPU{a: 0x1a, l: 0x2b}, CPU{a: 0x45, l: 0x2b}, map[string]bool{"Z": false, "H": true, "C": false}, nil},
// {"ADD", "A", "L",  CPU{a: 0x00, l: 0x00}, CPU{a: 0x00, l: 0x00}, map[string]bool{"Z": true, "H": false, "C": false}, nil},
// {"ADD", "A", "d8", 0x22, 0x0000, CPU{a: 0x1a}, CPU{a: 0x3c}, map[string]bool{"Z": false, "H": false, "C": false}, nil},
// {"ADD", "A", "d8", 0xf2, 0x0000, CPU{a: 0x1a}, CPU{a: 0x0c}, map[string]bool{"Z": false, "H": false, "C": true}, nil},
// {"ADD", "A", "d8", 0x2b, 0x0000, CPU{a: 0x1a}, CPU{a: 0x45}, map[string]bool{"Z": false, "H": true, "C": false}, nil},
// {"ADD", "A", "d8",  CPU{a: 0x00}, CPU{a: 0x00}, map[string]bool{"Z": true, "H": false, "C": false}, nil},

///////////////////////

// {"ADC", "A", "(HL)",  CPU{a: 0x1a, h: 0xa7, l: 0xf8, f: 0x00}, CPU{a: 0x3c, h: 0xa7, l: 0xf8, f: 0x00}, map[string]bool{"Z": false, "H": false, "C": false},  &testableMemory{actual:  map[uint16]byte{0xa7f8: 0x22}}},
// {"ADC", "A", "(HL)",  CPU{a: 0x1a, h: 0xa7, l: 0xf8, f: 0x00}, CPU{a: 0x0c, h: 0xa7, l: 0xf8, f: 0x00}, map[string]bool{"Z": false, "H": false, "C": true},  &testableMemory{actual:  map[uint16]byte{0xa7f8: 0xf2}}},
// {"ADC", "A", "(HL)",  CPU{a: 0x1a, h: 0xa7, l: 0xf8, f: 0x00}, CPU{a: 0x45, h: 0xa7, l: 0xf8, f: 0x00}, map[string]bool{"Z": false, "H": true, "C": false},  &testableMemory{actual:  map[uint16]byte{0xa7f8: 0x2b}}},
// {"ADC", "A", "(HL)", 0x00, 0x00p00, CPU{a: 0x00, h: 0xa7, l: 0xf8, f: 0x00}, CPU{a: 0x00, h: 0xa7, l: 0xf8, f: 0x00}, map[string]bool{"Z": true, "H": false, "C": false},  &testableMemory{actual:  map[uint16]byte{0xa7f8: 0x00}}},
// {"ADC", "A", "(HL)",  CPU{a: 0x1a, h: 0xa7, l: 0xf8, f: cFlag}, CPU{a: 0x3d, h: 0xa7, l: 0xf8, f: cFlag}, map[string]bool{"Z": false, "H": false, "C": false},  &testableMemory{actual:  map[uint16]byte{0xa7f8: 0x22}}},
// {"ADC", "A", "(HL)",  CPU{a: 0x1a, h: 0xa7, l: 0xf8, f: cFlag}, CPU{a: 0x0d, h: 0xa7, l: 0xf8, f: cFlag}, map[string]bool{"Z": false, "H": false, "C": true},  &testableMemory{actual:  map[uint16]byte{0xa7f8: 0xf2}}},
// {"ADC", "A", "(HL)",  CPU{a: 0x1a, h: 0xa7, l: 0xf8, f: cFlag}, CPU{a: 0x46, h: 0xa7, l: 0xf8, f: cFlag}, map[string]bool{"Z": false, "H": true, "C": false},  &testableMemory{actual:  map[uint16]byte{0xa7f8: 0x2b}}},
// {"ADC", "A", "(HL)",  CPU{a: 0x00, h: 0xa7, l: 0xf8, f: cFlag}, CPU{a: 0x00, h: 0xa7, l: 0xf8, f: cFlag}, map[string]bool{"Z": true, "H": true, "C": true},  &testableMemory{actual:  map[uint16]byte{0xa7f8: 0xff}}},

// {"ADC", "A", "A",  CPU{a: 0x12, f: 0x00}, CPU{a: 0x24, f: 0x00}, map[string]bool{"Z": false, "H": false, "C": false}, nil},
// {"ADC", "A", "A",  CPU{a: 0xa3, f: 0x00}, CPU{a: 0x46, f: 0x00}, map[string]bool{"Z": false, "H": false, "C": true}, nil},
// {"ADC", "A", "A",  CPU{a: 0x1a, f: 0x00}, CPU{a: 0x34, f: 0x00}, map[string]bool{"Z": false, "H": true, "C": false}, nil},
// {"ADC", "A", "A",  CPU{a: 0x00, f: 0x00}, CPU{a: 0x00, f: 0x00}, map[string]bool{"Z": true, "H": false, "C": false}, nil},
// {"ADC", "A", "A",  CPU{a: 0x12, f: cFlag}, CPU{a: 0x25, f: cFlag}, map[string]bool{"Z": false, "H": false, "C": false}, nil},
// {"ADC", "A", "A",  CPU{a: 0xa3, f: cFlag}, CPU{a: 0x47, f: cFlag}, map[string]bool{"Z": false, "H": false, "C": true}, nil},
// {"ADC", "A", "A",  CPU{a: 0x1a, f: cFlag}, CPU{a: 0x35, f: cFlag}, map[string]bool{"Z": false, "H": true, "C": false}, nil},

///////////////////////

///////////////////////
// GOOD

// {"ADD", "A", "(HL)",  CPU{a: 0x1a, h: 0xa7, l: 0xf8}, CPU{a: 0x3c, h: 0xa7, l: 0xf8}, map[string]bool{"Z": false, "H": false, "C": false},  &testableMemory{actual:  map[uint16]byte{0xa7f8: 0x22}}},
// {"ADD", "A", "(HL)",  CPU{a: 0x1a, h: 0xa7, l: 0xf8}, CPU{a: 0x0c, h: 0xa7, l: 0xf8}, map[string]bool{"Z": false, "H": false, "C": true},  &testableMemory{actual:  map[uint16]byte{0xa7f8: 0xf2}}},
// {"ADD", "A", "(HL)",  CPU{a: 0x1a, h: 0xa7, l: 0xf8}, CPU{a: 0x45, h: 0xa7, l: 0xf8}, map[string]bool{"Z": false, "H": true, "C": false},  &testableMemory{actual:  map[uint16]byte{0xa7f8: 0x2b}}},
// {"ADD", "A", "(HL)",  CPU{a: 0x00, h: 0xa7, l: 0xf8}, CPU{a: 0x00, h: 0xa7, l: 0xf8}, map[string]bool{"Z": true, "H": false, "C": false},  &testableMemory{actual:  map[uint16]byte{0xa7f8: 0x00}}},

///////////////////////

///////////////////////
// GOOD

// 	{"AND", "A", "",  CPU{a: 0x1a}, CPU{a: 0x1a}, map[string]bool{"Z": false}, nil},
// 	{"AND", "A", "",  CPU{a: 0x00}, CPU{a: 0x00}, map[string]bool{"Z": true}, nil},
// 	{"AND", "B", "",  CPU{a: 0x1a, b: 0x33}, CPU{a: 0x12, b: 0x33}, map[string]bool{"Z": false}, nil},
// 	{"AND", "B", "",  CPU{a: 0x1a, b: 0x21}, CPU{a: 0x00, b: 0x21}, map[string]bool{"Z": true}, nil},
// 	{"AND", "C", "",  CPU{a: 0x1a, c: 0x33}, CPU{a: 0x12, c: 0x33}, map[string]bool{"Z": false}, nil},
// 	{"AND", "C", "",  CPU{a: 0x1a, c: 0x21}, CPU{a: 0x00, c: 0x21}, map[string]bool{"Z": true}, nil},
// 	{"AND", "D", "",  CPU{a: 0x1a, d: 0x33}, CPU{a: 0x12, d: 0x33}, map[string]bool{"Z": false}, nil},
// 	{"AND", "D", "",  CPU{a: 0x1a, d: 0x21}, CPU{a: 0x00, d: 0x21}, map[string]bool{"Z": true}, nil},
// 	{"AND", "E", "",  CPU{a: 0x1a, e: 0x33}, CPU{a: 0x12, e: 0x33}, map[string]bool{"Z": false}, nil},
// 	{"AND", "E", "",  CPU{a: 0x1a, e: 0x21}, CPU{a: 0x00, e: 0x21}, map[string]bool{"Z": true}, nil},
// 	{"AND", "H", "",  CPU{a: 0x1a, h: 0x33}, CPU{a: 0x12, h: 0x33}, map[string]bool{"Z": false}, nil},
// 	{"AND", "H", "",  CPU{a: 0x1a, h: 0x21}, CPU{a: 0x00, h: 0x21}, map[string]bool{"Z": true}, nil},
// 	{"AND", "L", "",  CPU{a: 0x1a, l: 0x33}, CPU{a: 0x12, l: 0x33}, map[string]bool{"Z": false}, nil},
// 	{"AND", "L", "",  CPU{a: 0x1a, l: 0x21}, CPU{a: 0x00, l: 0x21}, map[string]bool{"Z": true}, nil},
// 	{"AND", "d8", "", 0x33, 0x0000, CPU{a: 0x1a}, CPU{a: 0x12}, map[string]bool{"Z": false}, nil},
// 	{"AND", "d8", "", 0x21, 0x0000, CPU{a: 0x1a}, CPU{a: 0x00}, map[string]bool{"Z": true}, nil}}

///

var dispatchTwoByteInstructionTests = []dispatchTwoByteInstructionTestTable{

//

}

var dispatchThreeByteInstructionTests = []dispatchThreeByteInstructionTestTable{

//

}

var dispatchPrefixedInstructionTests = []dispatchPrefixedInstructionTestTable{

	// BIT 6 (HL) [Z 0 1 -]
	{0x76, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8, zf: false, nf: false, hf: true}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x40}}},
	{0x76, CPU{h: 0xa7, l: 0xf8, zf: true, nf: true, hf: true, cf: true}, CPU{h: 0xa7, l: 0xf8, zf: true, nf: false, hf: true, cf: true}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0xbf}}},

	// BIT 1 E [Z 0 1 -]
	{0x4b, CPU{e: 0x02}, CPU{e: 0x02, zf: false, nf: false, hf: true}, nil},
	{0x4b, CPU{e: 0xfd, zf: true, nf: true, hf: true, cf: true}, CPU{e: 0xfd, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 2 D [Z 0 1 -]
	{0x52, CPU{d: 0x04}, CPU{d: 0x04, zf: false, nf: false, hf: true}, nil},
	{0x52, CPU{d: 0xfb, zf: true, nf: true, hf: true, cf: true}, CPU{d: 0xfb, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 3 (HL) [Z 0 1 -]
	{0x5e, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8, zf: false, nf: false, hf: true}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x08}}},
	{0x5e, CPU{h: 0xa7, l: 0xf8, zf: true, nf: true, hf: true, cf: true}, CPU{h: 0xa7, l: 0xf8, zf: true, nf: false, hf: true, cf: true}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0xf7}}},

	// BIT 7 D [Z 0 1 -]
	{0x7a, CPU{d: 0x80}, CPU{d: 0x80, zf: false, nf: false, hf: true}, nil},
	{0x7a, CPU{d: 0x7f, zf: true, nf: true, hf: true, cf: true}, CPU{d: 0x7f, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 6 H [Z 0 1 -]
	{0x74, CPU{h: 0x40}, CPU{h: 0x40, zf: false, nf: false, hf: true}, nil},
	{0x74, CPU{h: 0xbf, zf: true, nf: true, hf: true, cf: true}, CPU{h: 0xbf, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 3 B [Z 0 1 -]
	{0x58, CPU{b: 0x08}, CPU{b: 0x08, zf: false, nf: false, hf: true}, nil},
	{0x58, CPU{b: 0xf7, zf: true, nf: true, hf: true, cf: true}, CPU{b: 0xf7, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 5 (HL) [Z 0 1 -]
	{0x6e, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8, zf: false, nf: false, hf: true}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x20}}},
	{0x6e, CPU{h: 0xa7, l: 0xf8, zf: true, nf: true, hf: true, cf: true}, CPU{h: 0xa7, l: 0xf8, zf: true, nf: false, hf: true, cf: true}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0xdf}}},

	// BIT 4 H [Z 0 1 -]
	{0x64, CPU{h: 0x10}, CPU{h: 0x10, zf: false, nf: false, hf: true}, nil},
	{0x64, CPU{h: 0xef, zf: true, nf: true, hf: true, cf: true}, CPU{h: 0xef, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 6 A [Z 0 1 -]
	{0x77, CPU{a: 0x40}, CPU{a: 0x40, zf: false, nf: false, hf: true}, nil},
	{0x77, CPU{a: 0xbf, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0xbf, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 1 (HL) [Z 0 1 -]
	{0x4e, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8, zf: false, nf: false, hf: true}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x02}}},
	{0x4e, CPU{h: 0xa7, l: 0xf8, zf: true, nf: true, hf: true, cf: true}, CPU{h: 0xa7, l: 0xf8, zf: true, nf: false, hf: true, cf: true}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0xfd}}},

	// BIT 2 E [Z 0 1 -]
	{0x53, CPU{e: 0x04}, CPU{e: 0x04, zf: false, nf: false, hf: true}, nil},
	{0x53, CPU{e: 0xfb, zf: true, nf: true, hf: true, cf: true}, CPU{e: 0xfb, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 4 (HL) [Z 0 1 -]
	{0x66, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8, zf: false, nf: false, hf: true}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x10}}},
	{0x66, CPU{h: 0xa7, l: 0xf8, zf: true, nf: true, hf: true, cf: true}, CPU{h: 0xa7, l: 0xf8, zf: true, nf: false, hf: true, cf: true}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0xef}}},

	// BIT 0 L [Z 0 1 -]
	{0x45, CPU{l: 0x01}, CPU{l: 0x01, zf: false, nf: false, hf: true}, nil},
	{0x45, CPU{l: 0xfe, zf: true, nf: true, hf: true, cf: true}, CPU{l: 0xfe, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 2 L [Z 0 1 -]
	{0x55, CPU{l: 0x04}, CPU{l: 0x04, zf: false, nf: false, hf: true}, nil},
	{0x55, CPU{l: 0xfb, zf: true, nf: true, hf: true, cf: true}, CPU{l: 0xfb, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 3 C [Z 0 1 -]
	{0x59, CPU{c: 0x08}, CPU{c: 0x08, zf: false, nf: false, hf: true}, nil},
	{0x59, CPU{c: 0xf7, zf: true, nf: true, hf: true, cf: true}, CPU{c: 0xf7, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 5 A [Z 0 1 -]
	{0x6f, CPU{a: 0x20}, CPU{a: 0x20, zf: false, nf: false, hf: true}, nil},
	{0x6f, CPU{a: 0xdf, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0xdf, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 2 B [Z 0 1 -]
	{0x50, CPU{b: 0x04}, CPU{b: 0x04, zf: false, nf: false, hf: true}, nil},
	{0x50, CPU{b: 0xfb, zf: true, nf: true, hf: true, cf: true}, CPU{b: 0xfb, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 4 D [Z 0 1 -]
	{0x62, CPU{d: 0x10}, CPU{d: 0x10, zf: false, nf: false, hf: true}, nil},
	{0x62, CPU{d: 0xef, zf: true, nf: true, hf: true, cf: true}, CPU{d: 0xef, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 6 B [Z 0 1 -]
	{0x70, CPU{b: 0x40}, CPU{b: 0x40, zf: false, nf: false, hf: true}, nil},
	{0x70, CPU{b: 0xbf, zf: true, nf: true, hf: true, cf: true}, CPU{b: 0xbf, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 3 L [Z 0 1 -]
	{0x5d, CPU{l: 0x08}, CPU{l: 0x08, zf: false, nf: false, hf: true}, nil},
	{0x5d, CPU{l: 0xf7, zf: true, nf: true, hf: true, cf: true}, CPU{l: 0xf7, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 5 E [Z 0 1 -]
	{0x6b, CPU{e: 0x20}, CPU{e: 0x20, zf: false, nf: false, hf: true}, nil},
	{0x6b, CPU{e: 0xdf, zf: true, nf: true, hf: true, cf: true}, CPU{e: 0xdf, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 7 E [Z 0 1 -]
	{0x7b, CPU{e: 0x80}, CPU{e: 0x80, zf: false, nf: false, hf: true}, nil},
	{0x7b, CPU{e: 0x7f, zf: true, nf: true, hf: true, cf: true}, CPU{e: 0x7f, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 7 C [Z 0 1 -]
	{0x79, CPU{c: 0x80}, CPU{c: 0x80, zf: false, nf: false, hf: true}, nil},
	{0x79, CPU{c: 0x7f, zf: true, nf: true, hf: true, cf: true}, CPU{c: 0x7f, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 1 B [Z 0 1 -]
	{0x48, CPU{b: 0x02}, CPU{b: 0x02, zf: false, nf: false, hf: true}, nil},
	{0x48, CPU{b: 0xfd, zf: true, nf: true, hf: true, cf: true}, CPU{b: 0xfd, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 6 L [Z 0 1 -]
	{0x75, CPU{l: 0x40}, CPU{l: 0x40, zf: false, nf: false, hf: true}, nil},
	{0x75, CPU{l: 0xbf, zf: true, nf: true, hf: true, cf: true}, CPU{l: 0xbf, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 6 E [Z 0 1 -]
	{0x73, CPU{e: 0x40}, CPU{e: 0x40, zf: false, nf: false, hf: true}, nil},
	{0x73, CPU{e: 0xbf, zf: true, nf: true, hf: true, cf: true}, CPU{e: 0xbf, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 2 (HL) [Z 0 1 -]
	{0x56, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8, zf: false, nf: false, hf: true}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x04}}},
	{0x56, CPU{h: 0xa7, l: 0xf8, zf: true, nf: true, hf: true, cf: true}, CPU{h: 0xa7, l: 0xf8, zf: true, nf: false, hf: true, cf: true}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0xfb}}},

	// BIT 6 C [Z 0 1 -]
	{0x71, CPU{c: 0x40}, CPU{c: 0x40, zf: false, nf: false, hf: true}, nil},
	{0x71, CPU{c: 0xbf, zf: true, nf: true, hf: true, cf: true}, CPU{c: 0xbf, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 0 (HL) [Z 0 1 -]
	{0x46, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8, zf: false, nf: false, hf: true}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x01}}},
	{0x46, CPU{h: 0xa7, l: 0xf8, zf: true, nf: true, hf: true, cf: true}, CPU{h: 0xa7, l: 0xf8, zf: true, nf: false, hf: true, cf: true}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0xfe}}},

	// BIT 1 D [Z 0 1 -]
	{0x4a, CPU{d: 0x02}, CPU{d: 0x02, zf: false, nf: false, hf: true}, nil},
	{0x4a, CPU{d: 0xfd, zf: true, nf: true, hf: true, cf: true}, CPU{d: 0xfd, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 1 L [Z 0 1 -]
	{0x4d, CPU{l: 0x02}, CPU{l: 0x02, zf: false, nf: false, hf: true}, nil},
	{0x4d, CPU{l: 0xfd, zf: true, nf: true, hf: true, cf: true}, CPU{l: 0xfd, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 5 C [Z 0 1 -]
	{0x69, CPU{c: 0x20}, CPU{c: 0x20, zf: false, nf: false, hf: true}, nil},
	{0x69, CPU{c: 0xdf, zf: true, nf: true, hf: true, cf: true}, CPU{c: 0xdf, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 4 E [Z 0 1 -]
	{0x63, CPU{e: 0x10}, CPU{e: 0x10, zf: false, nf: false, hf: true}, nil},
	{0x63, CPU{e: 0xef, zf: true, nf: true, hf: true, cf: true}, CPU{e: 0xef, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 0 H [Z 0 1 -]
	{0x44, CPU{h: 0x01}, CPU{h: 0x01, zf: false, nf: false, hf: true}, nil},
	{0x44, CPU{h: 0xfe, zf: true, nf: true, hf: true, cf: true}, CPU{h: 0xfe, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 5 B [Z 0 1 -]
	{0x68, CPU{b: 0x20}, CPU{b: 0x20, zf: false, nf: false, hf: true}, nil},
	{0x68, CPU{b: 0xdf, zf: true, nf: true, hf: true, cf: true}, CPU{b: 0xdf, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 2 H [Z 0 1 -]
	{0x54, CPU{h: 0x04}, CPU{h: 0x04, zf: false, nf: false, hf: true}, nil},
	{0x54, CPU{h: 0xfb, zf: true, nf: true, hf: true, cf: true}, CPU{h: 0xfb, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 4 B [Z 0 1 -]
	{0x60, CPU{b: 0x10}, CPU{b: 0x10, zf: false, nf: false, hf: true}, nil},
	{0x60, CPU{b: 0xef, zf: true, nf: true, hf: true, cf: true}, CPU{b: 0xef, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 0 A [Z 0 1 -]
	{0x47, CPU{a: 0x01}, CPU{a: 0x01, zf: false, nf: false, hf: true}, nil},
	{0x47, CPU{a: 0xfe, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0xfe, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 1 H [Z 0 1 -]
	{0x4c, CPU{h: 0x02}, CPU{h: 0x02, zf: false, nf: false, hf: true}, nil},
	{0x4c, CPU{h: 0xfd, zf: true, nf: true, hf: true, cf: true}, CPU{h: 0xfd, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 7 H [Z 0 1 -]
	{0x7c, CPU{h: 0x80}, CPU{h: 0x80, zf: false, nf: false, hf: true}, nil},
	{0x7c, CPU{h: 0x7f, zf: true, nf: true, hf: true, cf: true}, CPU{h: 0x7f, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 0 D [Z 0 1 -]
	{0x42, CPU{d: 0x01}, CPU{d: 0x01, zf: false, nf: false, hf: true}, nil},
	{0x42, CPU{d: 0xfe, zf: true, nf: true, hf: true, cf: true}, CPU{d: 0xfe, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 1 C [Z 0 1 -]
	{0x49, CPU{c: 0x02}, CPU{c: 0x02, zf: false, nf: false, hf: true}, nil},
	{0x49, CPU{c: 0xfd, zf: true, nf: true, hf: true, cf: true}, CPU{c: 0xfd, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 7 B [Z 0 1 -]
	{0x78, CPU{b: 0x80}, CPU{b: 0x80, zf: false, nf: false, hf: true}, nil},
	{0x78, CPU{b: 0x7f, zf: true, nf: true, hf: true, cf: true}, CPU{b: 0x7f, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 0 E [Z 0 1 -]
	{0x43, CPU{e: 0x01}, CPU{e: 0x01, zf: false, nf: false, hf: true}, nil},
	{0x43, CPU{e: 0xfe, zf: true, nf: true, hf: true, cf: true}, CPU{e: 0xfe, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 3 A [Z 0 1 -]
	{0x5f, CPU{a: 0x08}, CPU{a: 0x08, zf: false, nf: false, hf: true}, nil},
	{0x5f, CPU{a: 0xf7, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0xf7, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 6 D [Z 0 1 -]
	{0x72, CPU{d: 0x40}, CPU{d: 0x40, zf: false, nf: false, hf: true}, nil},
	{0x72, CPU{d: 0xbf, zf: true, nf: true, hf: true, cf: true}, CPU{d: 0xbf, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 7 (HL) [Z 0 1 -]
	{0x7e, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8, zf: false, nf: false, hf: true}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x80}}},
	{0x7e, CPU{h: 0xa7, l: 0xf8, zf: true, nf: true, hf: true, cf: true}, CPU{h: 0xa7, l: 0xf8, zf: true, nf: false, hf: true, cf: true}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x7f}}},

	// BIT 5 D [Z 0 1 -]
	{0x6a, CPU{d: 0x20}, CPU{d: 0x20, zf: false, nf: false, hf: true}, nil},
	{0x6a, CPU{d: 0xdf, zf: true, nf: true, hf: true, cf: true}, CPU{d: 0xdf, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 2 A [Z 0 1 -]
	{0x57, CPU{a: 0x04}, CPU{a: 0x04, zf: false, nf: false, hf: true}, nil},
	{0x57, CPU{a: 0xfb, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0xfb, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 1 A [Z 0 1 -]
	{0x4f, CPU{a: 0x02}, CPU{a: 0x02, zf: false, nf: false, hf: true}, nil},
	{0x4f, CPU{a: 0xfd, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0xfd, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 3 D [Z 0 1 -]
	{0x5a, CPU{d: 0x08}, CPU{d: 0x08, zf: false, nf: false, hf: true}, nil},
	{0x5a, CPU{d: 0xf7, zf: true, nf: true, hf: true, cf: true}, CPU{d: 0xf7, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 3 H [Z 0 1 -]
	{0x5c, CPU{h: 0x08}, CPU{h: 0x08, zf: false, nf: false, hf: true}, nil},
	{0x5c, CPU{h: 0xf7, zf: true, nf: true, hf: true, cf: true}, CPU{h: 0xf7, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 0 B [Z 0 1 -]
	{0x40, CPU{b: 0x01}, CPU{b: 0x01, zf: false, nf: false, hf: true}, nil},
	{0x40, CPU{b: 0xfe, zf: true, nf: true, hf: true, cf: true}, CPU{b: 0xfe, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 7 A [Z 0 1 -]
	{0x7f, CPU{a: 0x80}, CPU{a: 0x80, zf: false, nf: false, hf: true}, nil},
	{0x7f, CPU{a: 0x7f, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x7f, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 4 C [Z 0 1 -]
	{0x61, CPU{c: 0x10}, CPU{c: 0x10, zf: false, nf: false, hf: true}, nil},
	{0x61, CPU{c: 0xef, zf: true, nf: true, hf: true, cf: true}, CPU{c: 0xef, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 5 H [Z 0 1 -]
	{0x6c, CPU{h: 0x20}, CPU{h: 0x20, zf: false, nf: false, hf: true}, nil},
	{0x6c, CPU{h: 0xdf, zf: true, nf: true, hf: true, cf: true}, CPU{h: 0xdf, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 4 L [Z 0 1 -]
	{0x65, CPU{l: 0x10}, CPU{l: 0x10, zf: false, nf: false, hf: true}, nil},
	{0x65, CPU{l: 0xef, zf: true, nf: true, hf: true, cf: true}, CPU{l: 0xef, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 0 C [Z 0 1 -]
	{0x41, CPU{c: 0x01}, CPU{c: 0x01, zf: false, nf: false, hf: true}, nil},
	{0x41, CPU{c: 0xfe, zf: true, nf: true, hf: true, cf: true}, CPU{c: 0xfe, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 7 L [Z 0 1 -]
	{0x7d, CPU{l: 0x80}, CPU{l: 0x80, zf: false, nf: false, hf: true}, nil},
	{0x7d, CPU{l: 0x7f, zf: true, nf: true, hf: true, cf: true}, CPU{l: 0x7f, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 2 C [Z 0 1 -]
	{0x51, CPU{c: 0x04}, CPU{c: 0x04, zf: false, nf: false, hf: true}, nil},
	{0x51, CPU{c: 0xfb, zf: true, nf: true, hf: true, cf: true}, CPU{c: 0xfb, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 3 E [Z 0 1 -]
	{0x5b, CPU{e: 0x08}, CPU{e: 0x08, zf: false, nf: false, hf: true}, nil},
	{0x5b, CPU{e: 0xf7, zf: true, nf: true, hf: true, cf: true}, CPU{e: 0xf7, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 4 A [Z 0 1 -]
	{0x67, CPU{a: 0x10}, CPU{a: 0x10, zf: false, nf: false, hf: true}, nil},
	{0x67, CPU{a: 0xef, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0xef, zf: true, nf: false, hf: true, cf: true}, nil},

	// BIT 5 L [Z 0 1 -]
	{0x6d, CPU{l: 0x20}, CPU{l: 0x20, zf: false, nf: false, hf: true}, nil},
	{0x6d, CPU{l: 0xdf, zf: true, nf: true, hf: true, cf: true}, CPU{l: 0xdf, zf: true, nf: false, hf: true, cf: true}, nil},

	// RES 1 A []
	{0x8f, CPU{a: 0x02}, CPU{a: 0x00}, nil},

	// RES 3 H []
	{0x9c, CPU{h: 0x08}, CPU{h: 0x00}, nil},

	// RES 3 L []
	{0x9d, CPU{l: 0x08}, CPU{l: 0x00}, nil},

	// RES 3 E []
	{0x9b, CPU{e: 0x08}, CPU{e: 0x00}, nil},

	// RES 4 C []
	{0xa1, CPU{c: 0x10}, CPU{c: 0x00}, nil},

	// RES 3 (HL) []
	{0x9e, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x08}, expected: map[uint16]byte{0xa7f8: 0x00}}},

	// RES 7 L []
	{0xbd, CPU{l: 0x80}, CPU{l: 0x00}, nil},

	// RES 5 (HL) []
	{0xae, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x20}, expected: map[uint16]byte{0xa7f8: 0x00}}},

	// RES 7 B []
	{0xb8, CPU{b: 0x80}, CPU{b: 0x00}, nil},

	// RES 1 D []
	{0x8a, CPU{d: 0x02}, CPU{d: 0x00}, nil},

	// RES 0 E []
	{0x83, CPU{e: 0x01}, CPU{e: 0x00}, nil},

	// RES 2 C []
	{0x91, CPU{c: 0x04}, CPU{c: 0x00}, nil},

	// RES 2 A []
	{0x97, CPU{a: 0x04}, CPU{a: 0x00}, nil},

	// RES 4 L []
	{0xa5, CPU{l: 0x10}, CPU{l: 0x00}, nil},

	// RES 7 E []
	{0xbb, CPU{e: 0x80}, CPU{e: 0x00}, nil},

	// RES 1 E []
	{0x8b, CPU{e: 0x02}, CPU{e: 0x00}, nil},

	// RES 5 B []
	{0xa8, CPU{b: 0x20}, CPU{b: 0x00}, nil},

	// RES 7 D []
	{0xba, CPU{d: 0x80}, CPU{d: 0x00}, nil},

	// RES 6 B []
	{0xb0, CPU{b: 0x40}, CPU{b: 0x00}, nil},

	// RES 0 L []
	{0x85, CPU{l: 0x01}, CPU{l: 0x00}, nil},

	// RES 7 A []
	{0xbf, CPU{a: 0x80}, CPU{a: 0x00}, nil},

	// RES 1 L []
	{0x8d, CPU{l: 0x02}, CPU{l: 0x00}, nil},

	// RES 1 C []
	{0x89, CPU{c: 0x02}, CPU{c: 0x00}, nil},

	// RES 2 L []
	{0x95, CPU{l: 0x04}, CPU{l: 0x00}, nil},

	// RES 1 B []
	{0x88, CPU{b: 0x02}, CPU{b: 0x00}, nil},

	// RES 6 E []
	{0xb3, CPU{e: 0x40}, CPU{e: 0x00}, nil},

	// RES 1 H []
	{0x8c, CPU{h: 0x02}, CPU{h: 0x00}, nil},

	// RES 2 D []
	{0x92, CPU{d: 0x04}, CPU{d: 0x00}, nil},

	// RES 2 H []
	{0x94, CPU{h: 0x04}, CPU{h: 0x00}, nil},

	// RES 5 L []
	{0xad, CPU{l: 0x20}, CPU{l: 0x00}, nil},

	// RES 7 (HL) []
	{0xbe, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x80}, expected: map[uint16]byte{0xa7f8: 0x00}}},

	// RES 2 B []
	{0x90, CPU{b: 0x04}, CPU{b: 0x00}, nil},

	// RES 6 H []
	{0xb4, CPU{h: 0x40}, CPU{h: 0x00}, nil},

	// RES 6 C []
	{0xb1, CPU{c: 0x40}, CPU{c: 0x00}, nil},

	// RES 7 H []
	{0xbc, CPU{h: 0x80}, CPU{h: 0x00}, nil},

	// RES 5 E []
	{0xab, CPU{e: 0x20}, CPU{e: 0x00}, nil},

	// RES 6 L []
	{0xb5, CPU{l: 0x40}, CPU{l: 0x00}, nil},

	// RES 0 D []
	{0x82, CPU{d: 0x01}, CPU{d: 0x00}, nil},

	// RES 1 (HL) []
	{0x8e, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x02}, expected: map[uint16]byte{0xa7f8: 0x00}}},

	// RES 4 A []
	{0xa7, CPU{a: 0x10}, CPU{a: 0x00}, nil},

	// RES 4 D []
	{0xa2, CPU{d: 0x10}, CPU{d: 0x00}, nil},

	// RES 3 B []
	{0x98, CPU{b: 0x08}, CPU{b: 0x00}, nil},

	// RES 6 D []
	{0xb2, CPU{d: 0x40}, CPU{d: 0x00}, nil},

	// RES 5 D []
	{0xaa, CPU{d: 0x20}, CPU{d: 0x00}, nil},

	// RES 3 A []
	{0x9f, CPU{a: 0x08}, CPU{a: 0x00}, nil},

	// RES 5 H []
	{0xac, CPU{h: 0x20}, CPU{h: 0x00}, nil},

	// RES 4 E []
	{0xa3, CPU{e: 0x10}, CPU{e: 0x00}, nil},

	// RES 6 (HL) []
	{0xb6, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x40}, expected: map[uint16]byte{0xa7f8: 0x00}}},

	// RES 3 C []
	{0x99, CPU{c: 0x08}, CPU{c: 0x00}, nil},

	// RES 7 C []
	{0xb9, CPU{c: 0x80}, CPU{c: 0x00}, nil},

	// RES 4 B []
	{0xa0, CPU{b: 0x10}, CPU{b: 0x00}, nil},

	// RES 0 A []
	{0x87, CPU{a: 0x01}, CPU{a: 0x00}, nil},

	// RES 2 (HL) []
	{0x96, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x04}, expected: map[uint16]byte{0xa7f8: 0x00}}},

	// RES 5 C []
	{0xa9, CPU{c: 0x20}, CPU{c: 0x00}, nil},

	// RES 0 B []
	{0x80, CPU{b: 0x01}, CPU{b: 0x00}, nil},

	// RES 3 D []
	{0x9a, CPU{d: 0x08}, CPU{d: 0x00}, nil},

	// RES 4 H []
	{0xa4, CPU{h: 0x10}, CPU{h: 0x00}, nil},

	// RES 0 C []
	{0x81, CPU{c: 0x01}, CPU{c: 0x00}, nil},

	// RES 0 (HL) []
	{0x86, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x01}, expected: map[uint16]byte{0xa7f8: 0x00}}},

	// RES 4 (HL) []
	{0xa6, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x10}, expected: map[uint16]byte{0xa7f8: 0x00}}},

	// RES 5 A []
	{0xaf, CPU{a: 0x20}, CPU{a: 0x00}, nil},

	// RES 6 A []
	{0xb7, CPU{a: 0x40}, CPU{a: 0x00}, nil},

	// RES 0 H []
	{0x84, CPU{h: 0x01}, CPU{h: 0x00}, nil},

	// RES 2 E []
	{0x93, CPU{e: 0x04}, CPU{e: 0x00}, nil},

	// SET 4 E []
	{0xe3, CPU{e: 0x00}, CPU{e: 0x10}, nil},

	// SET 6 H []
	{0xf4, CPU{h: 0x00}, CPU{h: 0x40}, nil},

	// SET 4 C []
	{0xe1, CPU{c: 0x00}, CPU{c: 0x10}, nil},

	// SET 5 C []
	{0xe9, CPU{c: 0x00}, CPU{c: 0x20}, nil},

	// SET 0 D []
	{0xc2, CPU{d: 0x00}, CPU{d: 0x01}, nil},

	// SET 5 E []
	{0xeb, CPU{e: 0x00}, CPU{e: 0x20}, nil},

	// SET 6 L []
	{0xf5, CPU{l: 0x00}, CPU{l: 0x40}, nil},

	// SET 1 E []
	{0xcb, CPU{e: 0x00}, CPU{e: 0x02}, nil},

	// SET 7 L []
	{0xfd, CPU{l: 0x00}, CPU{l: 0x80}, nil},

	// SET 0 H []
	{0xc4, CPU{h: 0x00}, CPU{h: 0x01}, nil},

	// SET 1 D []
	{0xca, CPU{d: 0x00}, CPU{d: 0x02}, nil},

	// SET 0 E []
	{0xc3, CPU{e: 0x00}, CPU{e: 0x01}, nil},

	// SET 3 B []
	{0xd8, CPU{b: 0x00}, CPU{b: 0x08}, nil},

	// SET 1 B []
	{0xc8, CPU{b: 0x00}, CPU{b: 0x02}, nil},

	// SET 7 C []
	{0xf9, CPU{c: 0x00}, CPU{c: 0x80}, nil},

	// SET 7 A []
	{0xff, CPU{a: 0x00}, CPU{a: 0x80}, nil},

	// SET 1 C []
	{0xc9, CPU{c: 0x00}, CPU{c: 0x02}, nil},

	// SET 4 H []
	{0xe4, CPU{h: 0x00}, CPU{h: 0x10}, nil},

	// SET 5 H []
	{0xec, CPU{h: 0x00}, CPU{h: 0x20}, nil},

	// SET 3 L []
	{0xdd, CPU{l: 0x00}, CPU{l: 0x08}, nil},

	// SET 2 B []
	{0xd0, CPU{b: 0x00}, CPU{b: 0x04}, nil},

	// SET 5 B []
	{0xe8, CPU{b: 0x00}, CPU{b: 0x20}, nil},

	// SET 5 (HL) []
	{0xee, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x00}, expected: map[uint16]byte{0xa7f8: 0x20}}},

	// SET 0 (HL) []
	{0xc6, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x00}, expected: map[uint16]byte{0xa7f8: 0x01}}},

	// SET 2 C []
	{0xd1, CPU{c: 0x00}, CPU{c: 0x04}, nil},

	// SET 2 D []
	{0xd2, CPU{d: 0x00}, CPU{d: 0x04}, nil},

	// SET 1 (HL) []
	{0xce, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x00}, expected: map[uint16]byte{0xa7f8: 0x02}}},

	// SET 7 B []
	{0xf8, CPU{b: 0x00}, CPU{b: 0x80}, nil},

	// SET 2 H []
	{0xd4, CPU{h: 0x00}, CPU{h: 0x04}, nil},

	// SET 2 L []
	{0xd5, CPU{l: 0x00}, CPU{l: 0x04}, nil},

	// SET 3 (HL) []
	{0xde, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x00}, expected: map[uint16]byte{0xa7f8: 0x08}}},

	// SET 2 (HL) []
	{0xd6, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x00}, expected: map[uint16]byte{0xa7f8: 0x04}}},

	// SET 6 (HL) []
	{0xf6, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x00}, expected: map[uint16]byte{0xa7f8: 0x40}}},

	// SET 4 A []
	{0xe7, CPU{a: 0x00}, CPU{a: 0x10}, nil},

	// SET 1 L []
	{0xcd, CPU{l: 0x00}, CPU{l: 0x02}, nil},

	// SET 2 A []
	{0xd7, CPU{a: 0x00}, CPU{a: 0x04}, nil},

	// SET 4 (HL) []
	{0xe6, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x00}, expected: map[uint16]byte{0xa7f8: 0x10}}},

	// SET 7 (HL) []
	{0xfe, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x00}, expected: map[uint16]byte{0xa7f8: 0x80}}},

	// SET 4 D []
	{0xe2, CPU{d: 0x00}, CPU{d: 0x10}, nil},

	// SET 0 C []
	{0xc1, CPU{c: 0x00}, CPU{c: 0x01}, nil},

	// SET 0 A []
	{0xc7, CPU{a: 0x00}, CPU{a: 0x01}, nil},

	// SET 7 D []
	{0xfa, CPU{d: 0x00}, CPU{d: 0x80}, nil},

	// SET 1 A []
	{0xcf, CPU{a: 0x00}, CPU{a: 0x02}, nil},

	// SET 5 D []
	{0xea, CPU{d: 0x00}, CPU{d: 0x20}, nil},

	// SET 6 B []
	{0xf0, CPU{b: 0x00}, CPU{b: 0x40}, nil},

	// SET 6 C []
	{0xf1, CPU{c: 0x00}, CPU{c: 0x40}, nil},

	// SET 6 D []
	{0xf2, CPU{d: 0x00}, CPU{d: 0x40}, nil},

	// SET 3 A []
	{0xdf, CPU{a: 0x00}, CPU{a: 0x08}, nil},

	// SET 0 B []
	{0xc0, CPU{b: 0x00}, CPU{b: 0x01}, nil},

	// SET 4 L []
	{0xe5, CPU{l: 0x00}, CPU{l: 0x10}, nil},

	// SET 4 B []
	{0xe0, CPU{b: 0x00}, CPU{b: 0x10}, nil},

	// SET 3 E []
	{0xdb, CPU{e: 0x00}, CPU{e: 0x08}, nil},

	// SET 7 H []
	{0xfc, CPU{h: 0x00}, CPU{h: 0x80}, nil},

	// SET 1 H []
	{0xcc, CPU{h: 0x00}, CPU{h: 0x02}, nil},

	// SET 3 C []
	{0xd9, CPU{c: 0x00}, CPU{c: 0x08}, nil},

	// SET 5 L []
	{0xed, CPU{l: 0x00}, CPU{l: 0x20}, nil},

	// SET 5 A []
	{0xef, CPU{a: 0x00}, CPU{a: 0x20}, nil},

	// SET 6 A []
	{0xf7, CPU{a: 0x00}, CPU{a: 0x40}, nil},

	// SET 3 D []
	{0xda, CPU{d: 0x00}, CPU{d: 0x08}, nil},

	// SET 3 H []
	{0xdc, CPU{h: 0x00}, CPU{h: 0x08}, nil},

	// SET 6 E []
	{0xf3, CPU{e: 0x00}, CPU{e: 0x40}, nil},

	// SET 7 E []
	{0xfb, CPU{e: 0x00}, CPU{e: 0x80}, nil},

	// SET 0 L []
	{0xc5, CPU{l: 0x00}, CPU{l: 0x01}, nil},

	// SET 2 E []
	{0xd3, CPU{e: 0x00}, CPU{e: 0x04}, nil}}
