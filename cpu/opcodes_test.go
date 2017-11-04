package cpu

import (
	"fmt"
	"testing"

	"github.com/scottyw/goomba/mem"
)

type testMemory struct {
	mem map[uint16]byte
}

// MemoryDump the contents of the whole address space to file
func (mem testMemory) MemoryDump() {
	fmt.Println(mem.mem)
}

// Read a byte from the chosen memory location
func (mem testMemory) Read(addr uint16) byte {
	val, present := mem.mem[addr]
	if !present {
		panic(fmt.Sprintf("Read memory failure: addr=0x%04x mem=%v", addr, mem))
	}
	return val
}

// Write a byte to the chosen memory location
func (mem testMemory) Write(addr uint16, b byte) {
	expected, present := mem.mem[addr]
	if !present || expected != b {
		panic(fmt.Sprintf("Write memory failure: addr=0x%04x mem=%v", addr, mem))
	}
}

type testTable struct {
	instruction uint8
	u8          uint8
	u16         uint16
	actualCPU   CPU
	expectedCPU CPU
	mem         mem.Memory
}

func (cpu *CPU) checkFlags(op opcodeMetadata, flagsToUpdate map[string]bool) {
	if len(op.Flags) > 0 {
		//for i, flag := range op.Flags {
		// 			switch flag {
		// 			case "0":
		// 			case "1":
		// default:
		// 	 	}
		// 	 	}
	}
}

func compareCPUs(t *testing.T, i int, test testTable) {
	if test.actualCPU != test.expectedCPU {
		t.Error("( TEST", i, ") CPUs do not match for: ", test)
		t.Error("  Expected : ", test.expectedCPU)
		t.Error("  Actual   : ", test.actualCPU)
	}
	// checkFlags() //FIXME
}

func TestDispatch(t *testing.T) {
	for i, test := range dispatchTests {
		(&test.actualCPU).dispatch(test.instruction)
		compareCPUs(t, i, test)
	}
}

func TestDispatch8(t *testing.T) {
	for i, test := range dispatch8Tests {
		(&test.actualCPU).dispatch8(test.instruction, test.u8)
		compareCPUs(t, i, test)
	}
}

var dispatch8Tests = []testTable{

//

}

var dispatchTests = []testTable{

	// // ADC A A [Z 0 H C]
	// {0x8f, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
	// {0x8f, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

	// // ADC A B [Z 0 H C]
	// {0x88, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
	// {0x88, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

	// // ADC A C [Z 0 H C]
	// {0x89, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
	// {0x89, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

	// // ADC A D [Z 0 H C]
	// {0x8a, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
	// {0x8a, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

	// // ADC A d8 [Z 0 H C]
	// {0xce, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
	// {0xce, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

	// // ADC A E [Z 0 H C]
	// {0x8b, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
	// {0x8b, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

	// // ADC A H [Z 0 H C]
	// {0x8c, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
	// {0x8c, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

	// // ADC A (HL) [Z 0 H C]
	// {0x8e, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
	// {0x8e, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

	// // ADC A L [Z 0 H C]
	// {0x8d, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
	// {0x8d, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

	// ADD A A [Z 0 H C]
	{0x87, 0x00, 0x0000, CPU{a: 0x12}, CPU{a: 0x24}, nil},
	{0x87, 0x00, 0x0000, CPU{a: 0xa3}, CPU{a: 0x46, cf: true}, nil},
	{0x87, 0x00, 0x0000, CPU{a: 0x1a}, CPU{a: 0x34, hf: true}, nil},
	{0x87, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00, zf: true}, nil},

	// ADD A B [Z 0 H C]
	{0x80, 0x00, 0x0000, CPU{a: 0x1a, b: 0x22}, CPU{a: 0x3c, b: 0x22}, nil},
	{0x80, 0x00, 0x0000, CPU{a: 0x1a, b: 0xf2}, CPU{a: 0x0c, b: 0xf2, cf: true}, nil},
	{0x80, 0x00, 0x0000, CPU{a: 0x1a, b: 0x2b}, CPU{a: 0x45, b: 0x2b, hf: true}, nil},
	{0x80, 0x00, 0x0000, CPU{a: 0x00, b: 0x00}, CPU{a: 0x00, b: 0x00, zf: true}, nil}}

// // ADD A C [Z 0 H C]
// {0x81, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x81, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // ADD A D [Z 0 H C]
// {0x82, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x82, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // ADD A E [Z 0 H C]
// {0x83, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x83, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // ADD A H [Z 0 H C]
// {0x84, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x84, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // ADD A L [Z 0 H C]
// {0x85, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x85, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // ADD A (HL) [Z 0 H C]
// {0x86, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x86, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // ADD A d8 [Z 0 H C]
// {0xc6, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc6, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // ADD HL BC [- 0 H C]
// {0x09, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x09, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // ADD HL DE [- 0 H C]
// {0x19, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x19, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // ADD HL HL [- 0 H C]
// {0x29, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x29, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // ADD HL SP [- 0 H C]
// {0x39, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x39, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // ADD SP r8 [0 0 H C]
// {0xe8, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe8, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // AND A  [Z 0 1 0]
// {0xa7, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa7, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // AND B  [Z 0 1 0]
// {0xa0, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa0, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // AND C  [Z 0 1 0]
// {0xa1, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa1, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // AND D  [Z 0 1 0]
// {0xa2, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa2, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // AND d8  [Z 0 1 0]
// {0xe6, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe6, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // AND E  [Z 0 1 0]
// {0xa3, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa3, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // AND H  [Z 0 1 0]
// {0xa4, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa4, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // AND (HL)  [Z 0 1 0]
// {0xa6, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa6, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // AND L  [Z 0 1 0]
// {0xa5, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa5, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // CALL a16  []
// {0xcd, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xcd, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // CALL C a16 []
// {0xdc, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xdc, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // CALL NC a16 []
// {0xd4, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd4, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // CALL NZ a16 []
// {0xc4, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc4, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // CALL Z a16 []
// {0xcc, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xcc, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // CCF   [- 0 0 C]
// {0x3f, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x3f, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // CP A  [Z 1 H C]
// {0xbf, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xbf, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // CP B  [Z 1 H C]
// {0xb8, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb8, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // CP C  [Z 1 H C]
// {0xb9, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb9, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // CP D  [Z 1 H C]
// {0xba, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xba, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // CP d8  [Z 1 H C]
// {0xfe, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xfe, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // CP E  [Z 1 H C]
// {0xbb, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xbb, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // CP H  [Z 1 H C]
// {0xbc, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xbc, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // CP (HL)  [Z 1 H C]
// {0xbe, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xbe, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // CP L  [Z 1 H C]
// {0xbd, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xbd, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // CPL   [- 1 1 -]
// {0x2f, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x2f, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // DAA   [Z - 0 C]
// {0x27, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x27, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // DEC A  [Z 1 H -]
// {0x3d, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x3d, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // DEC B  [Z 1 H -]
// {0x05, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x05, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // DEC BC  []
// {0x0b, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0b, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // DEC C  [Z 1 H -]
// {0x0d, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0d, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // DEC D  [Z 1 H -]
// {0x15, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x15, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // DEC DE  []
// {0x1b, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1b, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // DEC E  [Z 1 H -]
// {0x1d, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1d, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // DEC H  [Z 1 H -]
// {0x25, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x25, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // DEC (HL)  [Z 1 H -]
// {0x35, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x35, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // DEC HL  []
// {0x2b, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x2b, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // DEC L  [Z 1 H -]
// {0x2d, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x2d, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // DEC SP  []
// {0x3b, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x3b, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // DI   []
// {0xf3, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf3, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // EI   []
// {0xfb, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xfb, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // HALT   []
// {0x76, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x76, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // INC A  [Z 0 H -]
// {0x3c, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x3c, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // INC B  [Z 0 H -]
// {0x04, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x04, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // INC BC  []
// {0x03, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x03, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // INC C  [Z 0 H -]
// {0x0c, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0c, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // INC D  [Z 0 H -]
// {0x14, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x14, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // INC DE  []
// {0x13, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x13, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // INC E  [Z 0 H -]
// {0x1c, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1c, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // INC H  [Z 0 H -]
// {0x24, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x24, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // INC (HL)  [Z 0 H -]
// {0x34, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x34, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // INC HL  []
// {0x23, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x23, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // INC L  [Z 0 H -]
// {0x2c, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x2c, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // INC SP  []
// {0x33, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x33, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // JP a16  []
// {0xc3, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc3, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // JP C a16 []
// {0xda, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xda, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // JP NC a16 []
// {0xd2, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd2, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // JP Z a16 []
// {0xca, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xca, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // JP NZ a16 []
// {0xc2, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc2, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // JP (HL)  []
// {0xe9, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe9, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // JR C r8 []
// {0x38, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x38, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // JR NC r8 []
// {0x30, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x30, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // JR NZ r8 []
// {0x20, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x20, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // JR r8  []
// {0x18, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x18, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // JR Z r8 []
// {0x28, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x28, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD (a16) A []
// {0xea, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xea, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD (a16) SP []
// {0x08, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x08, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD (C) A []
// {0xe2, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe2, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD (DE) A []
// {0x12, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x12, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // // case 0x32 :  cpu. ld(cpu.(hl-), cpu.a) // LD (HL-) A []
// // // case 0x22 :  cpu. ld(cpu.(hl+), cpu.a) // LD (HL+) A []
// // LD A (a16) []
// {0xfa, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xfa, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD A (C) []
// {0xf2, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf2, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD A (DE) []
// {0x1a, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1a, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // // case 0x3a :  cpu. ld(cpu.a, cpu.(hl-)) // LD A (HL-) []
// // // case 0x2a :  cpu. ld(cpu.a, cpu.(hl+)) // LD A (HL+) []
// // LD A A []
// {0x7f, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x7f, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD A B []
// {0x78, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x78, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD A (BC) []
// {0x0a, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0a, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD A C []
// {0x79, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x79, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD A D []
// {0x7a, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x7a, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD A d8 []
// {0x3e, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x3e, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD A E []
// {0x7b, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x7b, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD A H []
// {0x7c, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x7c, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD A (HL) []
// {0x7e, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x7e, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD A L []
// {0x7d, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x7d, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD B A []
// {0x47, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x47, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD B B []
// {0x40, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x40, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD B C []
// {0x41, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x41, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD B D []
// {0x42, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x42, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD B d8 []
// {0x06, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x06, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD B E []
// {0x43, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x43, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD B H []
// {0x44, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x44, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD B (HL) []
// {0x46, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x46, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD B L []
// {0x45, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x45, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD (BC) A []
// {0x02, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x02, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD BC d16 []
// {0x01, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x01, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD C A []
// {0x4f, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x4f, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD C B []
// {0x48, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x48, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD C C []
// {0x49, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x49, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD C D []
// {0x4a, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x4a, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD C d8 []
// {0x0e, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0e, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD C E []
// {0x4b, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x4b, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD C H []
// {0x4c, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x4c, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD C (HL) []
// {0x4e, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x4e, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD C L []
// {0x4d, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x4d, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD D A []
// {0x57, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x57, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD D B []
// {0x50, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x50, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD D C []
// {0x51, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x51, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD D D []
// {0x52, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x52, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD D d8 []
// {0x16, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x16, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD D E []
// {0x53, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x53, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD D H []
// {0x54, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x54, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD D (HL) []
// {0x56, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x56, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD D L []
// {0x55, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x55, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD DE d16 []
// {0x11, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x11, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD E A []
// {0x5f, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x5f, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD E B []
// {0x58, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x58, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD E C []
// {0x59, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x59, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD E D []
// {0x5a, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x5a, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD E d8 []
// {0x1e, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1e, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD E E []
// {0x5b, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x5b, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD E H []
// {0x5c, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x5c, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD E (HL) []
// {0x5e, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x5e, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD E L []
// {0x5d, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x5d, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD H A []
// {0x67, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x67, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD H B []
// {0x60, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x60, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD H C []
// {0x61, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x61, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD H D []
// {0x62, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x62, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD H d8 []
// {0x26, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x26, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD H E []
// {0x63, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x63, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD H H []
// {0x64, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x64, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD H (HL) []
// {0x66, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x66, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD H L []
// {0x65, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x65, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD (HL) A []
// {0x77, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x77, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD (HL) B []
// {0x70, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x70, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD (HL) C []
// {0x71, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x71, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD (HL) D []
// {0x72, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x72, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD HL d16 []
// {0x21, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x21, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD (HL) d8 []
// {0x36, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x36, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD (HL) E []
// {0x73, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x73, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD (HL) H []
// {0x74, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x74, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD (HL) L []
// {0x75, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x75, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD HL SP+r8 [0 0 H C]
// {0xf8, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf8, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD L A []
// {0x6f, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x6f, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD L B []
// {0x68, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x68, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD L C []
// {0x69, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x69, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD L D []
// {0x6a, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x6a, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD L d8 []
// {0x2e, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x2e, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD L E []
// {0x6b, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x6b, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD L H []
// {0x6c, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x6c, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD L (HL) []
// {0x6e, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x6e, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD L L []
// {0x6d, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x6d, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD SP d16 []
// {0x31, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x31, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LD SP HL []
// {0xf9, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf9, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LDH (a8) A []
// {0xe0, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe0, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // LDH A (a8) []
// {0xf0, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf0, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // NOP   []
// {0x00, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x00, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // OR A  [Z 0 0 0]
// {0xb7, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb7, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // OR B  [Z 0 0 0]
// {0xb0, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb0, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // OR C  [Z 0 0 0]
// {0xb1, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb1, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // OR D  [Z 0 0 0]
// {0xb2, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb2, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // OR d8  [Z 0 0 0]
// {0xf6, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf6, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // OR E  [Z 0 0 0]
// {0xb3, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb3, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // OR H  [Z 0 0 0]
// {0xb4, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb4, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // OR (HL)  [Z 0 0 0]
// {0xb6, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb6, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // OR L  [Z 0 0 0]
// {0xb5, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb5, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // POP AF  [Z N H C]
// {0xf1, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf1, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // POP BC  []
// {0xc1, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc1, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // POP DE  []
// {0xd1, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd1, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // POP HL  []
// {0xe1, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe1, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // PREFIX CB  []
// {0xcb, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xcb, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // PUSH AF  []
// {0xf5, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf5, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // PUSH BC  []
// {0xc5, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc5, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // PUSH DE  []
// {0xd5, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd5, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // PUSH HL  []
// {0xe5, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe5, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RET   []
// {0xc9, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc9, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RET C  []
// {0xd8, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd8, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RET NC  []
// {0xd0, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd0, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RET NZ  []
// {0xc0, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc0, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RET Z  []
// {0xc8, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc8, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RETI   []
// {0xd9, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd9, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RLA   [0 0 0 C]
// {0x17, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x17, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RLCA   [0 0 0 C]
// {0x07, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x07, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RRA   [0 0 0 C]
// {0x1f, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1f, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RRCA   [0 0 0 C]
// {0x0f, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0f, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// FIXME
// // // case 0xc7 :  cpu. rst(cpu.00h) // RST 00H  []
// // // case 0xcf :  cpu. rst(cpu.08h) // RST 08H  []
// // // case 0xd7 :  cpu. rst(cpu.10h) // RST 10H  []
// // // case 0xdf :  cpu. rst(cpu.18h) // RST 18H  []
// // // case 0xe7 :  cpu. rst(cpu.20h) // RST 20H  []
// // // case 0xef :  cpu. rst(cpu.28h) // RST 28H  []
// // // case 0xf7 :  cpu. rst(cpu.30h) // RST 30H  []
// // // case 0xff :  cpu. rst(cpu.38h) // RST 38H  []

// // SBC A A [Z 1 H C]
// {0x9f, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9f, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SBC A B [Z 1 H C]
// {0x98, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x98, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SBC A C [Z 1 H C]
// {0x99, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x99, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SBC A D [Z 1 H C]
// {0x9a, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9a, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SBC A d8 [Z 1 H C]
// {0xde, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xde, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SBC A E [Z 1 H C]
// {0x9b, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9b, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SBC A H [Z 1 H C]
// {0x9c, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9c, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SBC A (HL) [Z 1 H C]
// {0x9e, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9e, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SBC A L [Z 1 H C]
// {0x9d, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9d, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SCF   [- 0 0 1]
// {0x37, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x37, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // STOP 0  []
// {0x10, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x10, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SUB A  [Z 1 H C]
// {0x97, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x97, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SUB B  [Z 1 H C]
// {0x90, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x90, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SUB C  [Z 1 H C]
// {0x91, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x91, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SUB D  [Z 1 H C]
// {0x92, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x92, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SUB d8  [Z 1 H C]
// {0xd6, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd6, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SUB E  [Z 1 H C]
// {0x93, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x93, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SUB H  [Z 1 H C]
// {0x94, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x94, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SUB (HL)  [Z 1 H C]
// {0x96, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x96, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SUB L  [Z 1 H C]
// {0x95, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x95, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // XOR A  [Z 0 0 0]
// {0xaf, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xaf, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // XOR B  [Z 0 0 0]
// {0xa8, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa8, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // XOR C  [Z 0 0 0]
// {0xa9, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa9, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // XOR D  [Z 0 0 0]
// {0xaa, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xaa, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // XOR d8  [Z 0 0 0]
// {0xee, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xee, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // XOR E  [Z 0 0 0]
// {0xab, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xab, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // XOR H  [Z 0 0 0]
// {0xac, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xac, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // XOR (HL)  [Z 0 0 0]
// {0xae, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xae, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // XOR L  [Z 0 0 0]
// {0xad, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xad, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 0 A [Z 0 1 -]
// {0x47, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x47, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 0 B [Z 0 1 -]
// {0x40, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x40, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 0 C [Z 0 1 -]
// {0x41, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x41, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 0 D [Z 0 1 -]
// {0x42, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x42, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 0 E [Z 0 1 -]
// {0x43, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x43, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 0 H [Z 0 1 -]
// {0x44, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x44, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 0 (HL) [Z 0 1 -]
// {0x46, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x46, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 0 L [Z 0 1 -]
// {0x45, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x45, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 1 A [Z 0 1 -]
// {0x4f, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x4f, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 1 B [Z 0 1 -]
// {0x48, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x48, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 1 C [Z 0 1 -]
// {0x49, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x49, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 1 D [Z 0 1 -]
// {0x4a, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x4a, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 1 E [Z 0 1 -]
// {0x4b, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x4b, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 1 H [Z 0 1 -]
// {0x4c, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x4c, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 1 (HL) [Z 0 1 -]
// {0x4e, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x4e, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 1 L [Z 0 1 -]
// {0x4d, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x4d, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 2 A [Z 0 1 -]
// {0x57, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x57, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 2 B [Z 0 1 -]
// {0x50, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x50, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 2 C [Z 0 1 -]
// {0x51, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x51, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 2 D [Z 0 1 -]
// {0x52, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x52, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 2 E [Z 0 1 -]
// {0x53, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x53, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 2 H [Z 0 1 -]
// {0x54, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x54, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 2 (HL) [Z 0 1 -]
// {0x56, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x56, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 2 L [Z 0 1 -]
// {0x55, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x55, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 3 A [Z 0 1 -]
// {0x5f, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x5f, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 3 B [Z 0 1 -]
// {0x58, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x58, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 3 C [Z 0 1 -]
// {0x59, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x59, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 3 D [Z 0 1 -]
// {0x5a, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x5a, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 3 E [Z 0 1 -]
// {0x5b, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x5b, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 3 H [Z 0 1 -]
// {0x5c, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x5c, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 3 (HL) [Z 0 1 -]
// {0x5e, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x5e, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 3 L [Z 0 1 -]
// {0x5d, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x5d, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 4 A [Z 0 1 -]
// {0x67, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x67, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 4 B [Z 0 1 -]
// {0x60, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x60, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 4 C [Z 0 1 -]
// {0x61, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x61, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 4 D [Z 0 1 -]
// {0x62, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x62, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 4 E [Z 0 1 -]
// {0x63, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x63, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 4 H [Z 0 1 -]
// {0x64, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x64, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 4 (HL) [Z 0 1 -]
// {0x66, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x66, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 4 L [Z 0 1 -]
// {0x65, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x65, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 5 A [Z 0 1 -]
// {0x6f, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x6f, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 5 B [Z 0 1 -]
// {0x68, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x68, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 5 C [Z 0 1 -]
// {0x69, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x69, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 5 D [Z 0 1 -]
// {0x6a, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x6a, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 5 E [Z 0 1 -]
// {0x6b, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x6b, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 5 H [Z 0 1 -]
// {0x6c, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x6c, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 5 (HL) [Z 0 1 -]
// {0x6e, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x6e, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 5 L [Z 0 1 -]
// {0x6d, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x6d, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 6 A [Z 0 1 -]
// {0x77, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x77, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 6 B [Z 0 1 -]
// {0x70, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x70, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 6 C [Z 0 1 -]
// {0x71, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x71, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 6 D [Z 0 1 -]
// {0x72, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x72, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 6 E [Z 0 1 -]
// {0x73, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x73, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 6 H [Z 0 1 -]
// {0x74, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x74, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 6 (HL) [Z 0 1 -]
// {0x76, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x76, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 6 L [Z 0 1 -]
// {0x75, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x75, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 7 A [Z 0 1 -]
// {0x7f, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x7f, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 7 B [Z 0 1 -]
// {0x78, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x78, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 7 C [Z 0 1 -]
// {0x79, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x79, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 7 D [Z 0 1 -]
// {0x7a, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x7a, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 7 E [Z 0 1 -]
// {0x7b, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x7b, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 7 H [Z 0 1 -]
// {0x7c, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x7c, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 7 (HL) [Z 0 1 -]
// {0x7e, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x7e, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // BIT 7 L [Z 0 1 -]
// {0x7d, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x7d, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 0 A []
// {0x87, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x87, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 0 B []
// {0x80, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x80, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 0 C []
// {0x81, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x81, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 0 D []
// {0x82, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x82, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 0 E []
// {0x83, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x83, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 0 H []
// {0x84, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x84, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 0 (HL) []
// {0x86, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x86, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 0 L []
// {0x85, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x85, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 1 A []
// {0x8f, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x8f, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 1 B []
// {0x88, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x88, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 1 C []
// {0x89, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x89, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 1 D []
// {0x8a, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x8a, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 1 E []
// {0x8b, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x8b, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 1 H []
// {0x8c, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x8c, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 1 (HL) []
// {0x8e, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x8e, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 1 L []
// {0x8d, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x8d, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 2 A []
// {0x97, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x97, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 2 B []
// {0x90, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x90, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 2 C []
// {0x91, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x91, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 2 D []
// {0x92, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x92, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 2 E []
// {0x93, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x93, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 2 H []
// {0x94, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x94, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 2 (HL) []
// {0x96, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x96, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 2 L []
// {0x95, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x95, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 3 A []
// {0x9f, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9f, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 3 B []
// {0x98, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x98, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 3 C []
// {0x99, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x99, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 3 D []
// {0x9a, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9a, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 3 E []
// {0x9b, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9b, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 3 H []
// {0x9c, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9c, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 3 (HL) []
// {0x9e, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9e, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 3 L []
// {0x9d, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x9d, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 4 A []
// {0xa7, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa7, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 4 B []
// {0xa0, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa0, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 4 C []
// {0xa1, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa1, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 4 D []
// {0xa2, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa2, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 4 E []
// {0xa3, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa3, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 4 H []
// {0xa4, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa4, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 4 (HL) []
// {0xa6, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa6, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 4 L []
// {0xa5, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa5, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 5 A []
// {0xaf, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xaf, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 5 B []
// {0xa8, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa8, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 5 C []
// {0xa9, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xa9, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 5 D []
// {0xaa, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xaa, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 5 E []
// {0xab, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xab, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 5 H []
// {0xac, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xac, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 5 (HL) []
// {0xae, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xae, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 5 L []
// {0xad, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xad, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 6 A []
// {0xb7, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb7, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 6 B []
// {0xb0, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb0, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 6 C []
// {0xb1, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb1, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 6 D []
// {0xb2, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb2, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 6 E []
// {0xb3, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb3, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 6 H []
// {0xb4, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb4, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 6 (HL) []
// {0xb6, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb6, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 6 L []
// {0xb5, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb5, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 7 A []
// {0xbf, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xbf, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 7 B []
// {0xb8, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb8, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 7 C []
// {0xb9, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xb9, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 7 D []
// {0xba, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xba, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 7 E []
// {0xbb, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xbb, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 7 H []
// {0xbc, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xbc, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 7 (HL) []
// {0xbe, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xbe, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RES 7 L []
// {0xbd, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xbd, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RL A  [Z 0 0 C]
// {0x17, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x17, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RL B  [Z 0 0 C]
// {0x10, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x10, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RL C  [Z 0 0 C]
// {0x11, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x11, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RL D  [Z 0 0 C]
// {0x12, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x12, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RL E  [Z 0 0 C]
// {0x13, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x13, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RL H  [Z 0 0 C]
// {0x14, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x14, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RL (HL)  [Z 0 0 C]
// {0x16, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x16, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RL L  [Z 0 0 C]
// {0x15, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x15, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RLC A  [Z 0 0 C]
// {0x07, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x07, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RLC B  [Z 0 0 C]
// {0x00, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x00, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RLC C  [Z 0 0 C]
// {0x01, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x01, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RLC D  [Z 0 0 C]
// {0x02, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x02, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RLC E  [Z 0 0 C]
// {0x03, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x03, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RLC H  [Z 0 0 C]
// {0x04, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x04, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RLC (HL)  [Z 0 0 C]
// {0x06, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x06, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RLC L  [Z 0 0 C]
// {0x05, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x05, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RR A  [Z 0 0 C]
// {0x1f, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1f, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RR B  [Z 0 0 C]
// {0x18, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x18, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RR C  [Z 0 0 C]
// {0x19, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x19, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RR D  [Z 0 0 C]
// {0x1a, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1a, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RR E  [Z 0 0 C]
// {0x1b, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1b, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RR H  [Z 0 0 C]
// {0x1c, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1c, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RR (HL)  [Z 0 0 C]
// {0x1e, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1e, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RR L  [Z 0 0 C]
// {0x1d, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x1d, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RRC A  [Z 0 0 C]
// {0x0f, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0f, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RRC B  [Z 0 0 C]
// {0x08, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x08, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RRC C  [Z 0 0 C]
// {0x09, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x09, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RRC D  [Z 0 0 C]
// {0x0a, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0a, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RRC E  [Z 0 0 C]
// {0x0b, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0b, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RRC H  [Z 0 0 C]
// {0x0c, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0c, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RRC (HL)  [Z 0 0 C]
// {0x0e, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0e, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // RRC L  [Z 0 0 C]
// {0x0d, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x0d, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 0 A []
// {0xc7, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc7, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 0 B []
// {0xc0, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc0, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 0 C []
// {0xc1, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc1, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 0 D []
// {0xc2, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc2, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 0 E []
// {0xc3, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc3, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 0 H []
// {0xc4, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc4, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 0 (HL) []
// {0xc6, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc6, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 0 L []
// {0xc5, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc5, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 1 A []
// {0xcf, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xcf, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 1 B []
// {0xc8, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc8, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 1 C []
// {0xc9, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xc9, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 1 D []
// {0xca, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xca, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 1 E []
// {0xcb, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xcb, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 1 H []
// {0xcc, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xcc, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 1 (HL) []
// {0xce, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xce, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 1 L []
// {0xcd, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xcd, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 2 A []
// {0xd7, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd7, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 2 B []
// {0xd0, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd0, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 2 C []
// {0xd1, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd1, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 2 D []
// {0xd2, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd2, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 2 E []
// {0xd3, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd3, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 2 H []
// {0xd4, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd4, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 2 (HL) []
// {0xd6, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd6, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 2 L []
// {0xd5, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd5, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 3 A []
// {0xdf, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xdf, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 3 B []
// {0xd8, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd8, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 3 C []
// {0xd9, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xd9, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 3 D []
// {0xda, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xda, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 3 E []
// {0xdb, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xdb, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 3 H []
// {0xdc, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xdc, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 3 (HL) []
// {0xde, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xde, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 3 L []
// {0xdd, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xdd, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 4 A []
// {0xe7, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe7, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 4 B []
// {0xe0, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe0, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 4 C []
// {0xe1, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe1, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 4 D []
// {0xe2, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe2, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 4 E []
// {0xe3, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe3, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 4 H []
// {0xe4, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe4, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 4 (HL) []
// {0xe6, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe6, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 4 L []
// {0xe5, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe5, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 5 A []
// {0xef, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xef, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 5 B []
// {0xe8, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe8, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 5 C []
// {0xe9, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xe9, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 5 D []
// {0xea, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xea, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 5 E []
// {0xeb, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xeb, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 5 H []
// {0xec, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xec, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 5 (HL) []
// {0xee, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xee, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 5 L []
// {0xed, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xed, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 6 A []
// {0xf7, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf7, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 6 B []
// {0xf0, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf0, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 6 C []
// {0xf1, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf1, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 6 D []
// {0xf2, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf2, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 6 E []
// {0xf3, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf3, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 6 H []
// {0xf4, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf4, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 6 (HL) []
// {0xf6, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf6, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 6 L []
// {0xf5, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf5, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 7 A []
// {0xff, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xff, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 7 B []
// {0xf8, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf8, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 7 C []
// {0xf9, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xf9, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 7 D []
// {0xfa, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xfa, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 7 E []
// {0xfb, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xfb, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 7 H []
// {0xfc, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xfc, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 7 (HL) []
// {0xfe, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xfe, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SET 7 L []
// {0xfd, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0xfd, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SLA A  [Z 0 0 C]
// {0x27, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x27, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SLA B  [Z 0 0 C]
// {0x20, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x20, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SLA C  [Z 0 0 C]
// {0x21, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x21, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SLA D  [Z 0 0 C]
// {0x22, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x22, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SLA E  [Z 0 0 C]
// {0x23, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x23, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SLA H  [Z 0 0 C]
// {0x24, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x24, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SLA (HL)  [Z 0 0 C]
// {0x26, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x26, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SLA L  [Z 0 0 C]
// {0x25, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x25, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SRA A  [Z 0 0 0]
// {0x2f, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x2f, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SRA B  [Z 0 0 0]
// {0x28, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x28, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SRA C  [Z 0 0 0]
// {0x29, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x29, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SRA D  [Z 0 0 0]
// {0x2a, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x2a, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SRA E  [Z 0 0 0]
// {0x2b, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x2b, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SRA H  [Z 0 0 0]
// {0x2c, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x2c, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SRA (HL)  [Z 0 0 0]
// {0x2e, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x2e, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SRA L  [Z 0 0 0]
// {0x2d, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x2d, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SRL A  [Z 0 0 C]
// {0x3f, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x3f, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SRL B  [Z 0 0 C]
// {0x38, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x38, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SRL C  [Z 0 0 C]
// {0x39, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x39, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SRL D  [Z 0 0 C]
// {0x3a, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x3a, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SRL E  [Z 0 0 C]
// {0x3b, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x3b, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SRL H  [Z 0 0 C]
// {0x3c, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x3c, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SRL (HL)  [Z 0 0 C]
// {0x3e, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x3e, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SRL L  [Z 0 0 C]
// {0x3d, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x3d, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SWAP A  [Z 0 0 0]
// {0x37, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x37, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SWAP B  [Z 0 0 0]
// {0x30, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x30, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SWAP C  [Z 0 0 0]
// {0x31, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x31, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SWAP D  [Z 0 0 0]
// {0x32, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x32, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SWAP E  [Z 0 0 0]
// {0x33, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x33, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SWAP H  [Z 0 0 0]
// {0x34, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x34, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SWAP (HL)  [Z 0 0 0]
// {0x36, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x36, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil},

// // SWAP L  [Z 0 0 0]
// {0x35, 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, nil},
// {0x35, 0x00, 0x0000, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, CPU{a: 0x00, zf: true, nf: true, hf: true, cf: true}, nil}}

// {"ADD", "A", "C", 0x00, 0x0000, CPU{a: 0x1a, c: 0x22}, CPU{a: 0x3c, c: 0x22}, map[string]bool{"Z": false, "H": false, "C": false}, nil},
// {"ADD", "A", "C", 0x00, 0x0000, CPU{a: 0x1a, c: 0xf2}, CPU{a: 0x0c, c: 0xf2}, map[string]bool{"Z": false, "H": false, "C": true}, nil},
// {"ADD", "A", "C", 0x00, 0x0000, CPU{a: 0x1a, c: 0x2b}, CPU{a: 0x45, c: 0x2b}, map[string]bool{"Z": false, "H": true, "C": false}, nil},
// {"ADD", "A", "C", 0x00, 0x0000, CPU{a: 0x00, c: 0x00}, CPU{a: 0x00, c: 0x00}, map[string]bool{"Z": true, "H": false, "C": false}, nil},
// {"ADD", "A", "D", 0x00, 0x0000, CPU{a: 0x1a, d: 0x22}, CPU{a: 0x3c, d: 0x22}, map[string]bool{"Z": false, "H": false, "C": false}, nil},
// {"ADD", "A", "D", 0x00, 0x0000, CPU{a: 0x1a, d: 0xf2}, CPU{a: 0x0c, d: 0xf2}, map[string]bool{"Z": false, "H": false, "C": true}, nil},
// {"ADD", "A", "D", 0x00, 0x0000, CPU{a: 0x1a, d: 0x2b}, CPU{a: 0x45, d: 0x2b}, map[string]bool{"Z": false, "H": true, "C": false}, nil},
// {"ADD", "A", "D", 0x00, 0x0000, CPU{a: 0x00, d: 0x00}, CPU{a: 0x00, d: 0x00}, map[string]bool{"Z": true, "H": false, "C": false}, nil},
// {"ADD", "A", "E", 0x00, 0x0000, CPU{a: 0x1a, e: 0x22}, CPU{a: 0x3c, e: 0x22}, map[string]bool{"Z": false, "H": false, "C": false}, nil},
// {"ADD", "A", "E", 0x00, 0x0000, CPU{a: 0x1a, e: 0xf2}, CPU{a: 0x0c, e: 0xf2}, map[string]bool{"Z": false, "H": false, "C": true}, nil},
// {"ADD", "A", "E", 0x00, 0x0000, CPU{a: 0x1a, e: 0x2b}, CPU{a: 0x45, e: 0x2b}, map[string]bool{"Z": false, "H": true, "C": false}, nil},
// {"ADD", "A", "E", 0x00, 0x0000, CPU{a: 0x00, e: 0x00}, CPU{a: 0x00, e: 0x00}, map[string]bool{"Z": true, "H": false, "C": false}, nil},
// {"ADD", "A", "H", 0x00, 0x0000, CPU{a: 0x1a, h: 0x22}, CPU{a: 0x3c, h: 0x22}, map[string]bool{"Z": false, "H": false, "C": false}, nil},
// {"ADD", "A", "H", 0x00, 0x0000, CPU{a: 0x1a, h: 0xf2}, CPU{a: 0x0c, h: 0xf2}, map[string]bool{"Z": false, "H": false, "C": true}, nil},
// {"ADD", "A", "H", 0x00, 0x0000, CPU{a: 0x1a, h: 0x2b}, CPU{a: 0x45, h: 0x2b}, map[string]bool{"Z": false, "H": true, "C": false}, nil},
// {"ADD", "A", "H", 0x00, 0x0000, CPU{a: 0x00, h: 0x00}, CPU{a: 0x00, h: 0x00}, map[string]bool{"Z": true, "H": false, "C": false}, nil},
// {"ADD", "A", "L", 0x00, 0x0000, CPU{a: 0x1a, l: 0x22}, CPU{a: 0x3c, l: 0x22}, map[string]bool{"Z": false, "H": false, "C": false}, nil},
// {"ADD", "A", "L", 0x00, 0x0000, CPU{a: 0x1a, l: 0xf2}, CPU{a: 0x0c, l: 0xf2}, map[string]bool{"Z": false, "H": false, "C": true}, nil},
// {"ADD", "A", "L", 0x00, 0x0000, CPU{a: 0x1a, l: 0x2b}, CPU{a: 0x45, l: 0x2b}, map[string]bool{"Z": false, "H": true, "C": false}, nil},
// {"ADD", "A", "L", 0x00, 0x0000, CPU{a: 0x00, l: 0x00}, CPU{a: 0x00, l: 0x00}, map[string]bool{"Z": true, "H": false, "C": false}, nil},
// {"ADD", "A", "d8", 0x22, 0x0000, CPU{a: 0x1a}, CPU{a: 0x3c}, map[string]bool{"Z": false, "H": false, "C": false}, nil},
// {"ADD", "A", "d8", 0xf2, 0x0000, CPU{a: 0x1a}, CPU{a: 0x0c}, map[string]bool{"Z": false, "H": false, "C": true}, nil},
// {"ADD", "A", "d8", 0x2b, 0x0000, CPU{a: 0x1a}, CPU{a: 0x45}, map[string]bool{"Z": false, "H": true, "C": false}, nil},
// {"ADD", "A", "d8", 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, map[string]bool{"Z": true, "H": false, "C": false}, nil},

///////////////////////

// {"ADC", "A", "(HL)", 0x00, 0x0000, CPU{a: 0x1a, h: 0xa7, l: 0xf8, f: 0x00}, CPU{a: 0x3c, h: 0xa7, l: 0xf8, f: 0x00}, map[string]bool{"Z": false, "H": false, "C": false}, testMemory{mem: map[uint16]byte{0xa7f8: 0x22}}},
// {"ADC", "A", "(HL)", 0x00, 0x0000, CPU{a: 0x1a, h: 0xa7, l: 0xf8, f: 0x00}, CPU{a: 0x0c, h: 0xa7, l: 0xf8, f: 0x00}, map[string]bool{"Z": false, "H": false, "C": true}, testMemory{mem: map[uint16]byte{0xa7f8: 0xf2}}},
// {"ADC", "A", "(HL)", 0x00, 0x0000, CPU{a: 0x1a, h: 0xa7, l: 0xf8, f: 0x00}, CPU{a: 0x45, h: 0xa7, l: 0xf8, f: 0x00}, map[string]bool{"Z": false, "H": true, "C": false}, testMemory{mem: map[uint16]byte{0xa7f8: 0x2b}}},
// {"ADC", "A", "(HL)", 0x00, 0x00p00, CPU{a: 0x00, h: 0xa7, l: 0xf8, f: 0x00}, CPU{a: 0x00, h: 0xa7, l: 0xf8, f: 0x00}, map[string]bool{"Z": true, "H": false, "C": false}, testMemory{mem: map[uint16]byte{0xa7f8: 0x00}}},
// {"ADC", "A", "(HL)", 0x00, 0x0000, CPU{a: 0x1a, h: 0xa7, l: 0xf8, f: cFlag}, CPU{a: 0x3d, h: 0xa7, l: 0xf8, f: cFlag}, map[string]bool{"Z": false, "H": false, "C": false}, testMemory{mem: map[uint16]byte{0xa7f8: 0x22}}},
// {"ADC", "A", "(HL)", 0x00, 0x0000, CPU{a: 0x1a, h: 0xa7, l: 0xf8, f: cFlag}, CPU{a: 0x0d, h: 0xa7, l: 0xf8, f: cFlag}, map[string]bool{"Z": false, "H": false, "C": true}, testMemory{mem: map[uint16]byte{0xa7f8: 0xf2}}},
// {"ADC", "A", "(HL)", 0x00, 0x0000, CPU{a: 0x1a, h: 0xa7, l: 0xf8, f: cFlag}, CPU{a: 0x46, h: 0xa7, l: 0xf8, f: cFlag}, map[string]bool{"Z": false, "H": true, "C": false}, testMemory{mem: map[uint16]byte{0xa7f8: 0x2b}}},
// {"ADC", "A", "(HL)", 0x00, 0x0000, CPU{a: 0x00, h: 0xa7, l: 0xf8, f: cFlag}, CPU{a: 0x00, h: 0xa7, l: 0xf8, f: cFlag}, map[string]bool{"Z": true, "H": true, "C": true}, testMemory{mem: map[uint16]byte{0xa7f8: 0xff}}},

// {"ADC", "A", "A", 0x00, 0x0000, CPU{a: 0x12, f: 0x00}, CPU{a: 0x24, f: 0x00}, map[string]bool{"Z": false, "H": false, "C": false}, nil},
// {"ADC", "A", "A", 0x00, 0x0000, CPU{a: 0xa3, f: 0x00}, CPU{a: 0x46, f: 0x00}, map[string]bool{"Z": false, "H": false, "C": true}, nil},
// {"ADC", "A", "A", 0x00, 0x0000, CPU{a: 0x1a, f: 0x00}, CPU{a: 0x34, f: 0x00}, map[string]bool{"Z": false, "H": true, "C": false}, nil},
// {"ADC", "A", "A", 0x00, 0x0000, CPU{a: 0x00, f: 0x00}, CPU{a: 0x00, f: 0x00}, map[string]bool{"Z": true, "H": false, "C": false}, nil},
// {"ADC", "A", "A", 0x00, 0x0000, CPU{a: 0x12, f: cFlag}, CPU{a: 0x25, f: cFlag}, map[string]bool{"Z": false, "H": false, "C": false}, nil},
// {"ADC", "A", "A", 0x00, 0x0000, CPU{a: 0xa3, f: cFlag}, CPU{a: 0x47, f: cFlag}, map[string]bool{"Z": false, "H": false, "C": true}, nil},
// {"ADC", "A", "A", 0x00, 0x0000, CPU{a: 0x1a, f: cFlag}, CPU{a: 0x35, f: cFlag}, map[string]bool{"Z": false, "H": true, "C": false}, nil},

///////////////////////

///////////////////////
// GOOD

// {"ADD", "A", "(HL)", 0x00, 0x0000, CPU{a: 0x1a, h: 0xa7, l: 0xf8}, CPU{a: 0x3c, h: 0xa7, l: 0xf8}, map[string]bool{"Z": false, "H": false, "C": false}, testMemory{mem: map[uint16]byte{0xa7f8: 0x22}}},
// {"ADD", "A", "(HL)", 0x00, 0x0000, CPU{a: 0x1a, h: 0xa7, l: 0xf8}, CPU{a: 0x0c, h: 0xa7, l: 0xf8}, map[string]bool{"Z": false, "H": false, "C": true}, testMemory{mem: map[uint16]byte{0xa7f8: 0xf2}}},
// {"ADD", "A", "(HL)", 0x00, 0x0000, CPU{a: 0x1a, h: 0xa7, l: 0xf8}, CPU{a: 0x45, h: 0xa7, l: 0xf8}, map[string]bool{"Z": false, "H": true, "C": false}, testMemory{mem: map[uint16]byte{0xa7f8: 0x2b}}},
// {"ADD", "A", "(HL)", 0x00, 0x0000, CPU{a: 0x00, h: 0xa7, l: 0xf8}, CPU{a: 0x00, h: 0xa7, l: 0xf8}, map[string]bool{"Z": true, "H": false, "C": false}, testMemory{mem: map[uint16]byte{0xa7f8: 0x00}}},

///////////////////////

///////////////////////
// GOOD

// 	{"AND", "A", "", 0x00, 0x0000, CPU{a: 0x1a}, CPU{a: 0x1a}, map[string]bool{"Z": false}, nil},
// 	{"AND", "A", "", 0x00, 0x0000, CPU{a: 0x00}, CPU{a: 0x00}, map[string]bool{"Z": true}, nil},
// 	{"AND", "B", "", 0x00, 0x0000, CPU{a: 0x1a, b: 0x33}, CPU{a: 0x12, b: 0x33}, map[string]bool{"Z": false}, nil},
// 	{"AND", "B", "", 0x00, 0x0000, CPU{a: 0x1a, b: 0x21}, CPU{a: 0x00, b: 0x21}, map[string]bool{"Z": true}, nil},
// 	{"AND", "C", "", 0x00, 0x0000, CPU{a: 0x1a, c: 0x33}, CPU{a: 0x12, c: 0x33}, map[string]bool{"Z": false}, nil},
// 	{"AND", "C", "", 0x00, 0x0000, CPU{a: 0x1a, c: 0x21}, CPU{a: 0x00, c: 0x21}, map[string]bool{"Z": true}, nil},
// 	{"AND", "D", "", 0x00, 0x0000, CPU{a: 0x1a, d: 0x33}, CPU{a: 0x12, d: 0x33}, map[string]bool{"Z": false}, nil},
// 	{"AND", "D", "", 0x00, 0x0000, CPU{a: 0x1a, d: 0x21}, CPU{a: 0x00, d: 0x21}, map[string]bool{"Z": true}, nil},
// 	{"AND", "E", "", 0x00, 0x0000, CPU{a: 0x1a, e: 0x33}, CPU{a: 0x12, e: 0x33}, map[string]bool{"Z": false}, nil},
// 	{"AND", "E", "", 0x00, 0x0000, CPU{a: 0x1a, e: 0x21}, CPU{a: 0x00, e: 0x21}, map[string]bool{"Z": true}, nil},
// 	{"AND", "H", "", 0x00, 0x0000, CPU{a: 0x1a, h: 0x33}, CPU{a: 0x12, h: 0x33}, map[string]bool{"Z": false}, nil},
// 	{"AND", "H", "", 0x00, 0x0000, CPU{a: 0x1a, h: 0x21}, CPU{a: 0x00, h: 0x21}, map[string]bool{"Z": true}, nil},
// 	{"AND", "L", "", 0x00, 0x0000, CPU{a: 0x1a, l: 0x33}, CPU{a: 0x12, l: 0x33}, map[string]bool{"Z": false}, nil},
// 	{"AND", "L", "", 0x00, 0x0000, CPU{a: 0x1a, l: 0x21}, CPU{a: 0x00, l: 0x21}, map[string]bool{"Z": true}, nil},
// 	{"AND", "d8", "", 0x33, 0x0000, CPU{a: 0x1a}, CPU{a: 0x12}, map[string]bool{"Z": false}, nil},
// 	{"AND", "d8", "", 0x21, 0x0000, CPU{a: 0x1a}, CPU{a: 0x00}, map[string]bool{"Z": true}, nil}}

///////////////////////
