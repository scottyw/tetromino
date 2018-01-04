package cpu

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/scottyw/goomba/mem"
)

////////////////////////////////////////////////////////////////
// Testable memory implementation
////////////////////////////////////////////////////////////////

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

// GenerateCrashReport to show memory state
func (mem testableMemory) GenerateCrashReport() {
	fmt.Println("TestMemory crash: ", mem.actual, mem.expected)
}

////////////////////////////////////////////////////////////////
// Test-agnostic validation functions
////////////////////////////////////////////////////////////////

func validateMemory(t *testing.T, mem testableMemory) {
	if mem.expected != nil && !reflect.DeepEqual(mem.actual, mem.expected) {
		t.Error("Memory does not match")
		t.Error("  Expected : ", mem.expected)
		t.Error("  Actual   : ", mem.actual)
	}
}

func compareCPUs(t *testing.T, expectedCPU, actualCPU *CPU, mem *testableMemory) {
	if *actualCPU != *expectedCPU {
		t.Error("CPUs do not match")
		t.Error("  Expected : ", expectedCPU)
		t.Error("  Actual   : ", actualCPU)
	}
	if mem != nil {
		validateMemory(t, *mem)
	}
}

////////////////////////////////////////////////////////////////
// Opcode unit tests
////////////////////////////////////////////////////////////////

func TestAdd(t *testing.T) {
	// Flags: [Z 0 H C]
	for _, test := range []struct {
		cpu, expectedCPU CPU
	}{
		{CPU{a: 0x12}, CPU{a: 0x24}},
		{CPU{a: 0xa3}, CPU{a: 0x46, cf: true}},
		{CPU{a: 0x1a}, CPU{a: 0x34, hf: true}},
		{CPU{a: 0x00}, CPU{a: 0x00, zf: true}},
	} {
		test.cpu.add(test.cpu.a)
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
	for _, test := range []struct {
		cpu, expectedCPU CPU
	}{
		{CPU{a: 0x1a, c: 0x22}, CPU{a: 0x3c, c: 0x22}},
		{CPU{a: 0x1a, c: 0xf2}, CPU{a: 0x0c, c: 0xf2, cf: true}},
		{CPU{a: 0x1a, c: 0x2b}, CPU{a: 0x45, c: 0x2b, hf: true}},
		{CPU{a: 0x00, c: 0x00}, CPU{a: 0x00, c: 0x00, zf: true}},
	} {
		test.cpu.add(test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestBit(t *testing.T) {
	// Flags: [Z 0 1 -]
	for _, test := range []struct {
		cpu, expectedCPU CPU
	}{
		{CPU{c: 0x04}, CPU{c: 0x04, zf: false, nf: false, hf: true}},
		{CPU{c: 0xfb, zf: true, nf: true, hf: true, cf: true}, CPU{c: 0xfb, zf: true, nf: false, hf: true, cf: true}},
	} {
		test.cpu.bit(2, &test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestBitHL(t *testing.T) {
	// Flags: [Z 0 1 -]
	for _, test := range []struct {
		cpu, expectedCPU CPU
		mem              mem.Memory
	}{
		{CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8, hf: true}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x04}}},
		{CPU{h: 0xa7, l: 0xf8, zf: true, nf: true, hf: true, cf: true}, CPU{h: 0xa7, l: 0xf8, zf: true, nf: false, hf: true, cf: true}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0xfb}}},
	} {
		test.cpu.bitHL(test.mem, 2)
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestSet(t *testing.T) {
	// Flags: []
	for _, test := range []struct {
		pos              uint8
		cpu, expectedCPU CPU
	}{
		{4, CPU{l: 0x00}, CPU{l: 0x10}},
		{5, CPU{l: 0x20}, CPU{l: 0x20}},
		{2, CPU{l: 0x11}, CPU{l: 0x15}},
	} {
		test.cpu.set(test.pos, &test.cpu.l)
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

////////
////////
////////
////////

// // RES 1 A []
// {0x8f, CPU{a: 0x02}, CPU{a: 0x00}, nil},

// // RES 3 H []
// {0x9c, CPU{h: 0x08}, CPU{h: 0x00}, nil},

// // RES 3 L []
// {0x9d, CPU{l: 0x08}, CPU{l: 0x00}, nil},

// // RES 3 E []
// {0x9b, CPU{e: 0x08}, CPU{e: 0x00}, nil},

// // RES 4 C []
// {0xa1, CPU{c: 0x10}, CPU{c: 0x00}, nil},

// // RES 3 (HL) []
// {0x9e, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x08}, expected: map[uint16]byte{0xa7f8: 0x00}}},

// // RES 7 L []
// {0xbd, CPU{l: 0x80}, CPU{l: 0x00}, nil},

// // RES 5 (HL) []
// {0xae, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x20}, expected: map[uint16]byte{0xa7f8: 0x00}}},

// // RES 7 B []
// {0xb8, CPU{b: 0x80}, CPU{b: 0x00}, nil},

// // RES 1 D []
// {0x8a, CPU{d: 0x02}, CPU{d: 0x00}, nil},

// // RES 0 E []
// {0x83, CPU{e: 0x01}, CPU{e: 0x00}, nil},

// // RES 2 C []
// {0x91, CPU{c: 0x04}, CPU{c: 0x00}, nil},

// // RES 2 A []
// {0x97, CPU{a: 0x04}, CPU{a: 0x00}, nil},

// // RES 4 L []
// {0xa5, CPU{l: 0x10}, CPU{l: 0x00}, nil},

// // RES 7 E []
// {0xbb, CPU{e: 0x80}, CPU{e: 0x00}, nil},

// // RES 1 E []
// {0x8b, CPU{e: 0x02}, CPU{e: 0x00}, nil},

// // RES 5 B []
// {0xa8, CPU{b: 0x20}, CPU{b: 0x00}, nil},

// // RES 7 D []
// {0xba, CPU{d: 0x80}, CPU{d: 0x00}, nil},

// // RES 6 B []
// {0xb0, CPU{b: 0x40}, CPU{b: 0x00}, nil},

// // RES 0 L []
// {0x85, CPU{l: 0x01}, CPU{l: 0x00}, nil},

// // RES 7 A []
// {0xbf, CPU{a: 0x80}, CPU{a: 0x00}, nil},

// // RES 1 L []
// {0x8d, CPU{l: 0x02}, CPU{l: 0x00}, nil},

// // RES 1 C []
// {0x89, CPU{c: 0x02}, CPU{c: 0x00}, nil},

// // RES 2 L []
// {0x95, CPU{l: 0x04}, CPU{l: 0x00}, nil},

// // RES 1 B []
// {0x88, CPU{b: 0x02}, CPU{b: 0x00}, nil},

// // RES 6 E []
// {0xb3, CPU{e: 0x40}, CPU{e: 0x00}, nil},

// // RES 1 H []
// {0x8c, CPU{h: 0x02}, CPU{h: 0x00}, nil},

// // RES 2 D []
// {0x92, CPU{d: 0x04}, CPU{d: 0x00}, nil},

// // RES 2 H []
// {0x94, CPU{h: 0x04}, CPU{h: 0x00}, nil},

// // RES 5 L []
// {0xad, CPU{l: 0x20}, CPU{l: 0x00}, nil},

// // RES 7 (HL) []
// {0xbe, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x80}, expected: map[uint16]byte{0xa7f8: 0x00}}},

// // RES 2 B []
// {0x90, CPU{b: 0x04}, CPU{b: 0x00}, nil},

// // RES 6 H []
// {0xb4, CPU{h: 0x40}, CPU{h: 0x00}, nil},

// // RES 6 C []
// {0xb1, CPU{c: 0x40}, CPU{c: 0x00}, nil},

// // RES 7 H []
// {0xbc, CPU{h: 0x80}, CPU{h: 0x00}, nil},

// // RES 5 E []
// {0xab, CPU{e: 0x20}, CPU{e: 0x00}, nil},

// // RES 6 L []
// {0xb5, CPU{l: 0x40}, CPU{l: 0x00}, nil},

// // RES 0 D []
// {0x82, CPU{d: 0x01}, CPU{d: 0x00}, nil},

// // RES 1 (HL) []
// {0x8e, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x02}, expected: map[uint16]byte{0xa7f8: 0x00}}},

// // RES 4 A []
// {0xa7, CPU{a: 0x10}, CPU{a: 0x00}, nil},

// // RES 4 D []
// {0xa2, CPU{d: 0x10}, CPU{d: 0x00}, nil},

// // RES 3 B []
// {0x98, CPU{b: 0x08}, CPU{b: 0x00}, nil},

// // RES 6 D []
// {0xb2, CPU{d: 0x40}, CPU{d: 0x00}, nil},

// // RES 5 D []
// {0xaa, CPU{d: 0x20}, CPU{d: 0x00}, nil},

// // RES 3 A []
// {0x9f, CPU{a: 0x08}, CPU{a: 0x00}, nil},

// // RES 5 H []
// {0xac, CPU{h: 0x20}, CPU{h: 0x00}, nil},

// // RES 4 E []
// {0xa3, CPU{e: 0x10}, CPU{e: 0x00}, nil},

// // RES 6 (HL) []
// {0xb6, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x40}, expected: map[uint16]byte{0xa7f8: 0x00}}},

// // RES 3 C []
// {0x99, CPU{c: 0x08}, CPU{c: 0x00}, nil},

// // RES 7 C []
// {0xb9, CPU{c: 0x80}, CPU{c: 0x00}, nil},

// // RES 4 B []
// {0xa0, CPU{b: 0x10}, CPU{b: 0x00}, nil},

// // RES 0 A []
// {0x87, CPU{a: 0x01}, CPU{a: 0x00}, nil},

// // RES 2 (HL) []
// {0x96, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x04}, expected: map[uint16]byte{0xa7f8: 0x00}}},

// // RES 5 C []
// {0xa9, CPU{c: 0x20}, CPU{c: 0x00}, nil},

// // RES 0 B []
// {0x80, CPU{b: 0x01}, CPU{b: 0x00}, nil},

// // RES 3 D []
// {0x9a, CPU{d: 0x08}, CPU{d: 0x00}, nil},

// // RES 4 H []
// {0xa4, CPU{h: 0x10}, CPU{h: 0x00}, nil},

// // RES 0 C []
// {0x81, CPU{c: 0x01}, CPU{c: 0x00}, nil},

// // RES 0 (HL) []
// {0x86, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x01}, expected: map[uint16]byte{0xa7f8: 0x00}}},

// // RES 4 (HL) []
// {0xa6, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x10}, expected: map[uint16]byte{0xa7f8: 0x00}}},

// // RES 5 A []
// {0xaf, CPU{a: 0x20}, CPU{a: 0x00}, nil},

// // RES 6 A []
// {0xb7, CPU{a: 0x40}, CPU{a: 0x00}, nil},

// // RES 0 H []
// {0x84, CPU{h: 0x01}, CPU{h: 0x00}, nil},

// // RES 2 E []
// {0x93, CPU{e: 0x04}, CPU{e: 0x00}, nil},

// // SET 4 E []
// {0xe3, CPU{e: 0x00}, CPU{e: 0x10}, nil},

// // SET 6 H []
// {0xf4, CPU{h: 0x00}, CPU{h: 0x40}, nil},

// // SET 4 C []
// {0xe1, CPU{c: 0x00}, CPU{c: 0x10}, nil},

// // SET 5 C []
// {0xe9, CPU{c: 0x00}, CPU{c: 0x20}, nil},

// // SET 0 D []
// {0xc2, CPU{d: 0x00}, CPU{d: 0x01}, nil},

// // SET 5 E []
// {0xeb, CPU{e: 0x00}, CPU{e: 0x20}, nil},

// // SET 6 L []
// {0xf5, CPU{l: 0x00}, CPU{l: 0x40}, nil},

// // SET 1 E []
// {0xcb, CPU{e: 0x00}, CPU{e: 0x02}, nil},

// // SET 7 L []
// {0xfd, CPU{l: 0x00}, CPU{l: 0x80}, nil},

// // SET 0 H []
// {0xc4, CPU{h: 0x00}, CPU{h: 0x01}, nil},

// // SET 1 D []
// {0xca, CPU{d: 0x00}, CPU{d: 0x02}, nil},

// // SET 0 E []
// {0xc3, CPU{e: 0x00}, CPU{e: 0x01}, nil},

// // SET 3 B []
// {0xd8, CPU{b: 0x00}, CPU{b: 0x08}, nil},

// // SET 1 B []
// {0xc8, CPU{b: 0x00}, CPU{b: 0x02}, nil},

// // SET 7 C []
// {0xf9, CPU{c: 0x00}, CPU{c: 0x80}, nil},

// // SET 7 A []
// {0xff, CPU{a: 0x00}, CPU{a: 0x80}, nil},

// // SET 1 C []
// {0xc9, CPU{c: 0x00}, CPU{c: 0x02}, nil},

// // SET 4 H []
// {0xe4, CPU{h: 0x00}, CPU{h: 0x10}, nil},

// // SET 5 H []
// {0xec, CPU{h: 0x00}, CPU{h: 0x20}, nil},

// // SET 3 L []
// {0xdd, CPU{l: 0x00}, CPU{l: 0x08}, nil},

// // SET 2 B []
// {0xd0, CPU{b: 0x00}, CPU{b: 0x04}, nil},

// // SET 5 B []
// {0xe8, CPU{b: 0x00}, CPU{b: 0x20}, nil},

// // SET 5 (HL) []
// {0xee, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x00}, expected: map[uint16]byte{0xa7f8: 0x20}}},

// // SET 0 (HL) []
// {0xc6, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x00}, expected: map[uint16]byte{0xa7f8: 0x01}}},

// // SET 2 C []
// {0xd1, CPU{c: 0x00}, CPU{c: 0x04}, nil},

// // SET 2 D []
// {0xd2, CPU{d: 0x00}, CPU{d: 0x04}, nil},

// // SET 1 (HL) []
// {0xce, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x00}, expected: map[uint16]byte{0xa7f8: 0x02}}},

// // SET 7 B []
// {0xf8, CPU{b: 0x00}, CPU{b: 0x80}, nil},

// // SET 2 H []
// {0xd4, CPU{h: 0x00}, CPU{h: 0x04}, nil},

// // SET 2 L []
// {0xd5, CPU{l: 0x00}, CPU{l: 0x04}, nil},

// // SET 3 (HL) []
// {0xde, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x00}, expected: map[uint16]byte{0xa7f8: 0x08}}},

// // SET 2 (HL) []
// {0xd6, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x00}, expected: map[uint16]byte{0xa7f8: 0x04}}},

// // SET 6 (HL) []
// {0xf6, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x00}, expected: map[uint16]byte{0xa7f8: 0x40}}},

// // SET 4 A []
// {0xe7, CPU{a: 0x00}, CPU{a: 0x10}, nil},

// // SET 1 L []
// {0xcd, CPU{l: 0x00}, CPU{l: 0x02}, nil},

// // SET 2 A []
// {0xd7, CPU{a: 0x00}, CPU{a: 0x04}, nil},

// // SET 4 (HL) []
// {0xe6, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x00}, expected: map[uint16]byte{0xa7f8: 0x10}}},

// // SET 7 (HL) []
// {0xfe, CPU{h: 0xa7, l: 0xf8}, CPU{h: 0xa7, l: 0xf8}, &testableMemory{actual: map[uint16]byte{0xa7f8: 0x00}, expected: map[uint16]byte{0xa7f8: 0x80}}},

// // SET 4 D []
// {0xe2, CPU{d: 0x00}, CPU{d: 0x10}, nil},

// // SET 0 C []
// {0xc1, CPU{c: 0x00}, CPU{c: 0x01}, nil},

// // SET 0 A []
// {0xc7, CPU{a: 0x00}, CPU{a: 0x01}, nil},

// // SET 7 D []
// {0xfa, CPU{d: 0x00}, CPU{d: 0x80}, nil},

// // SET 1 A []
// {0xcf, CPU{a: 0x00}, CPU{a: 0x02}, nil},

// // SET 5 D []
// {0xea, CPU{d: 0x00}, CPU{d: 0x20}, nil},

// // SET 6 B []
// {0xf0, CPU{b: 0x00}, CPU{b: 0x40}, nil},

// // SET 6 C []
// {0xf1, CPU{c: 0x00}, CPU{c: 0x40}, nil},

// // SET 6 D []
// {0xf2, CPU{d: 0x00}, CPU{d: 0x40}, nil},

// // SET 3 A []
// {0xdf, CPU{a: 0x00}, CPU{a: 0x08}, nil},

// // SET 0 B []
// {0xc0, CPU{b: 0x00}, CPU{b: 0x01}, nil},

// // SET 4 L []
// {0xe5, CPU{l: 0x00}, CPU{l: 0x10}, nil},

// // SET 4 B []
// {0xe0, CPU{b: 0x00}, CPU{b: 0x10}, nil},

// // SET 3 E []
// {0xdb, CPU{e: 0x00}, CPU{e: 0x08}, nil},

// // SET 7 H []
// {0xfc, CPU{h: 0x00}, CPU{h: 0x80}, nil},

// // SET 1 H []
// {0xcc, CPU{h: 0x00}, CPU{h: 0x02}, nil},

// // SET 3 C []
// {0xd9, CPU{c: 0x00}, CPU{c: 0x08}, nil},

// // SET 5 L []
// {0xed, CPU{l: 0x00}, CPU{l: 0x20}, nil},

// // SET 5 A []
// {0xef, CPU{a: 0x00}, CPU{a: 0x20}, nil},

// // SET 6 A []
// {0xf7, CPU{a: 0x00}, CPU{a: 0x40}, nil},

// // SET 3 D []
// {0xda, CPU{d: 0x00}, CPU{d: 0x08}, nil},

// // SET 3 H []
// {0xdc, CPU{h: 0x00}, CPU{h: 0x08}, nil},

// // SET 6 E []
// {0xf3, CPU{e: 0x00}, CPU{e: 0x40}, nil},

// // SET 7 E []
// {0xfb, CPU{e: 0x00}, CPU{e: 0x80}, nil},

// // SET 0 L []
// {0xc5, CPU{l: 0x00}, CPU{l: 0x01}, nil},

// // SET 2 E []
// {0xd3, CPU{e: 0x00}, CPU{e: 0x04}, nil},
