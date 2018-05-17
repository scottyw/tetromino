package cpu

import (
	"bytes"
	"testing"

	"github.com/scottyw/tetromino/pkg/gb/mem"
)

////////////////////////////////////////////////////////////////
// Test-agnostic validation functions
////////////////////////////////////////////////////////////////

func ReadRegion(memory *mem.Memory, startAddr, length uint16) []byte {
	region := make([]byte, length)
	for i := uint16(0); i < length; i++ {
		region[i] = memory.Read(startAddr + i)
	}
	return region
}

func compareCPUsAndMemory(t *testing.T, expectedCPU, actualCPU *CPU, expectedMem, actualMem *mem.Memory, startAddr, length uint16) {
	compareCPUs(t, expectedCPU, actualCPU)
	actual := ReadRegion(actualMem, startAddr, length)
	expected := ReadRegion(expectedMem, startAddr, length)
	if bytes.Compare(actual, expected) != 0 {
		t.Error("Memory does not match")
		t.Error("  Expected : ", expected)
		t.Error("  Actual   : ", actual)
	}
}

func compareCPUs(t *testing.T, expectedCPU, actualCPU *CPU) {
	if *actualCPU != *expectedCPU {
		t.Error("CPUs do not match")
		t.Error("  Expected : ", expectedCPU)
		t.Error("  Actual   : ", actualCPU)
	}
}

////////////////////////////////////////////////////////////////
// Opcode unit tests
////////////////////////////////////////////////////////////////

func TestAdc(t *testing.T) {
	// Flags: [Z 0 H C]
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{a: 0x1a, c: 0x22}, CPU{a: 0x3c, c: 0x22}},
		{CPU{a: 0x1a, c: 0xf2}, CPU{a: 0x0c, c: 0xf2, cf: true}},
		{CPU{a: 0x1a, c: 0x2b}, CPU{a: 0x45, c: 0x2b, hf: true}},
		{CPU{a: 0x00, c: 0x00}, CPU{a: 0x00, c: 0x00, zf: true}},
		{CPU{a: 0x1a, c: 0x22, cf: true}, CPU{a: 0x3d, c: 0x22}},
		{CPU{a: 0x1a, c: 0xf2, cf: true}, CPU{a: 0x0d, c: 0xf2, cf: true}},
		{CPU{a: 0x1a, c: 0x2b, cf: true}, CPU{a: 0x46, c: 0x2b, hf: true}},
		{CPU{a: 0xff, c: 0x00, cf: true}, CPU{a: 0x00, c: 0x00, zf: true, hf: true, cf: true}},
	} {
		test.cpu.adc(test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

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
		compareCPUs(t, &test.expectedCPU, &test.cpu)
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
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestAddHL(t *testing.T) {
	for _, test := range []struct {
		cpu, expectedCPU CPU
	}{
		{CPU{h: 0x1f, l: 0xb2}, CPU{h: 0x51, l: 0xc4, hf: true}},
		{CPU{h: 0xd1, l: 0xb2}, CPU{h: 0x03, l: 0xc4, cf: true}},
		{CPU{h: 0xcd, l: 0xee}, CPU{h: 0x00, l: 0x00, hf: true, cf: true}},
	} {
		test.cpu.addHL(0x3212)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
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
		test.cpu.bit(2, test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestCall(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{pc: 0xabcd, sp: 0x8642}, CPU{pc: 0x1af2, sp: 0x8640}},
	} {
		actual := mem.NewMemory(mem.NewHardwareRegisters(nil, nil), "")
		test.cpu.call("", 0x1af2, actual)
		expected := mem.NewMemory(mem.NewHardwareRegisters(nil, nil), "")
		expected.Write(0x8641, 0xab)
		expected.Write(0x8642, 0xcd)
		compareCPUsAndMemory(t, &test.expectedCPU, &test.cpu, expected, actual, 0x8640, 0xf)
	}
}

func TestCpl(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{a: 0xb1}, CPU{a: 0x4e, nf: true, hf: true}},
	} {
		test.cpu.cpl()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestPop(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{b: 0xff, c: 0x11, sp: 0x8640}, CPU{b: 0x1a, c: 0xf2, sp: 0x8642}},
	} {
		actual := mem.NewMemory(mem.NewHardwareRegisters(nil, nil), "")
		actual.Write(0x8641, 0x1a)
		actual.Write(0x8642, 0xf2)
		test.cpu.pop(&test.cpu.b, &test.cpu.c, actual)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestPush(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{b: 0x1a, c: 0xf2, sp: 0x8642}, CPU{b: 0x1a, c: 0xf2, sp: 0x8640}},
	} {
		actual := mem.NewMemory(mem.NewHardwareRegisters(nil, nil), "")
		test.cpu.push(test.cpu.b, test.cpu.c, actual)
		expected := mem.NewMemory(mem.NewHardwareRegisters(nil, nil), "")
		expected.Write(0x8641, 0x1a)
		expected.Write(0x8642, 0xf2)
		compareCPUsAndMemory(t, &test.expectedCPU, &test.cpu, expected, actual, 0x8640, 0xf)
	}
}

func TestRes(t *testing.T) {
	// Flags: [- - - -]
	for _, test := range []struct {
		cpu, expectedCPU CPU
	}{
		{CPU{c: 0x0a}, CPU{c: 0x0a}},
		{CPU{c: 0x0e, zf: true, nf: true, hf: true, cf: true}, CPU{c: 0x0a, zf: true, nf: true, hf: true, cf: true}},
	} {
		test.cpu.res(2, &test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestRet(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{pc: 0xabab, sp: 0x8640}, CPU{pc: 0x1af2, sp: 0x8642}},
	} {
		mem := mem.NewMemory(mem.NewHardwareRegisters(nil, nil), "")
		mem.Write(0x8641, 0x1a)
		mem.Write(0x8642, 0xf2)
		test.cpu.ret("", mem)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestReti(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{pc: 0xabab, sp: 0x8640}, CPU{pc: 0x1af2, sp: 0x8642, ime: true}},
	} {
		mem := mem.NewMemory(mem.NewHardwareRegisters(nil, nil), "")
		mem.Write(0x8641, 0x1a)
		mem.Write(0x8642, 0xf2)
		test.cpu.reti(mem)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestRl(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{c: 0xa8, cf: false}, CPU{c: 0x50, cf: true}},
		{CPU{c: 0xa8, cf: true}, CPU{c: 0x51, cf: true}},
		{CPU{c: 0x15, cf: false}, CPU{c: 0x2a, cf: false}},
		{CPU{c: 0x15, cf: true}, CPU{c: 0x2b, cf: false}},
		{CPU{c: 0x00}, CPU{c: 0x00, zf: true}},
	} {
		test.cpu.rl(&test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestRla(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{a: 0xa8, cf: false}, CPU{a: 0x50, cf: true}},
		{CPU{a: 0xa8, cf: true}, CPU{a: 0x51, cf: true}},
		{CPU{a: 0x15, cf: false}, CPU{a: 0x2a, cf: false}},
		{CPU{a: 0x15, cf: true}, CPU{a: 0x2b, cf: false}},
		{CPU{a: 0x00}, CPU{a: 0x00}},
	} {
		test.cpu.rla()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestRlc(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{c: 0xa8, cf: false}, CPU{c: 0x51, cf: true}},
		{CPU{c: 0xa8, cf: true}, CPU{c: 0x51, cf: true}},
		{CPU{c: 0x15, cf: false}, CPU{c: 0x2a, cf: false}},
		{CPU{c: 0x15, cf: true}, CPU{c: 0x2a, cf: false}},
		{CPU{c: 0x00}, CPU{c: 0x00, zf: true}},
	} {
		test.cpu.rlc(&test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestRlca(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{a: 0xa8, cf: false}, CPU{a: 0x51, cf: true}},
		{CPU{a: 0xa8, cf: true}, CPU{a: 0x51, cf: true}},
		{CPU{a: 0x15, cf: false}, CPU{a: 0x2a, cf: false}},
		{CPU{a: 0x15, cf: true}, CPU{a: 0x2a, cf: false}},
		{CPU{a: 0x00}, CPU{a: 0x00}},
	} {
		test.cpu.rlca()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestRr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{c: 0x15, cf: false}, CPU{c: 0x0a, cf: true}},
		{CPU{c: 0x15, cf: true}, CPU{c: 0x8a, cf: true}},
		{CPU{c: 0xa8, cf: false}, CPU{c: 0x54, cf: false}},
		{CPU{c: 0xa8, cf: true}, CPU{c: 0xd4, cf: false}},
		{CPU{c: 0x00}, CPU{c: 0x00, zf: true}},
	} {
		test.cpu.rr(&test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestRra(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{a: 0x15, cf: false}, CPU{a: 0x0a, cf: true}},
		{CPU{a: 0x15, cf: true}, CPU{a: 0x8a, cf: true}},
		{CPU{a: 0xa8, cf: false}, CPU{a: 0x54, cf: false}},
		{CPU{a: 0xa8, cf: true}, CPU{a: 0xd4, cf: false}},
		{CPU{a: 0x00}, CPU{a: 0x00, zf: true}},
	} {
		test.cpu.rra()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestRrc(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{c: 0x15, cf: false}, CPU{c: 0x8a, cf: true}},
		{CPU{c: 0x15, cf: true}, CPU{c: 0x8a, cf: true}},
		{CPU{c: 0xa8, cf: false}, CPU{c: 0x54, cf: false}},
		{CPU{c: 0xa8, cf: true}, CPU{c: 0x54, cf: false}},
		{CPU{c: 0x00}, CPU{c: 0x00, zf: true}},
	} {
		test.cpu.rrc(&test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestRrca(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{a: 0x15, cf: false}, CPU{a: 0x8a, cf: true}},
		{CPU{a: 0x15, cf: true}, CPU{a: 0x8a, cf: true}},
		{CPU{a: 0xa8, cf: false}, CPU{a: 0x54, cf: false}},
		{CPU{a: 0xa8, cf: true}, CPU{a: 0x54, cf: false}},
		{CPU{a: 0x00}, CPU{a: 0x00, zf: true}},
	} {
		test.cpu.rrca()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestRst(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{pc: 0xabcd, sp: 0x8642}, CPU{pc: 0x0008, sp: 0x8640}},
	} {
		actual := mem.NewMemory(mem.NewHardwareRegisters(nil, nil), "")
		test.cpu.rst(0x0008, actual)
		expected := mem.NewMemory(mem.NewHardwareRegisters(nil, nil), "")
		expected.Write(0x8641, 0xab)
		expected.Write(0x8642, 0xcd)
		compareCPUsAndMemory(t, &test.expectedCPU, &test.cpu, expected, actual, 0x8640, 0xf)
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
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestSla(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{c: 0xa9, cf: false}, CPU{c: 0x52, cf: true}},
		{CPU{c: 0x15, cf: true}, CPU{c: 0x2a, cf: false}},
		{CPU{c: 0x00}, CPU{c: 0x00, zf: true}},
	} {
		test.cpu.sla(&test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestSlaAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{cf: false}, CPU{cf: true}},
	} {
		actual := mem.NewMemory(mem.NewHardwareRegisters(nil, nil), "")
		actual.Write(0x8642, 0xa9)
		test.cpu.slaAddr(0x8642, actual)
		expected := mem.NewMemory(mem.NewHardwareRegisters(nil, nil), "")
		expected.Write(0x8642, 0x52)
		compareCPUsAndMemory(t, &test.expectedCPU, &test.cpu, expected, actual, 0x8642, 0x1)
	}
}

func TestSra(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{c: 0x55, cf: false}, CPU{c: 0x2a, cf: true}},
		{CPU{c: 0xa8, cf: true}, CPU{c: 0xd4, cf: false}},
		{CPU{c: 0x00}, CPU{c: 0x00, zf: true}},
	} {
		test.cpu.sra(&test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestSraAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{cf: false}, CPU{cf: true}},
	} {
		actual := mem.NewMemory(mem.NewHardwareRegisters(nil, nil), "")
		actual.Write(0x8642, 0x55)
		test.cpu.sraAddr(0x8642, actual)
		expected := mem.NewMemory(mem.NewHardwareRegisters(nil, nil), "")
		expected.Write(0x8642, 0x2a)
		compareCPUsAndMemory(t, &test.expectedCPU, &test.cpu, expected, actual, 0x8642, 0x1)
	}
}

func TestSrl(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{c: 0x55, cf: false}, CPU{c: 0x2a, cf: true}},
		{CPU{c: 0xa8, cf: true}, CPU{c: 0x54, cf: false}},
		{CPU{c: 0x00}, CPU{c: 0x00, zf: true}},
	} {
		test.cpu.srl(&test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestSrlAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{cf: true}, CPU{}},
	} {
		actual := mem.NewMemory(mem.NewHardwareRegisters(nil, nil), "")
		actual.Write(0x8642, 0xa8)
		test.cpu.srlAddr(0x8642, actual)
		expected := mem.NewMemory(mem.NewHardwareRegisters(nil, nil), "")
		expected.Write(0x8642, 0x54)
		compareCPUsAndMemory(t, &test.expectedCPU, &test.cpu, expected, actual, 0x8642, 0x1)
	}
}

func TestSwap(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{c: 0x15, cf: true}, CPU{c: 0x51}},
		{CPU{c: 0x00}, CPU{c: 0x00, zf: true}},
	} {
		test.cpu.swap(&test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestSwapAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{cf: true}, CPU{}},
	} {
		actual := mem.NewMemory(mem.NewHardwareRegisters(nil, nil), "")
		actual.Write(0x8641, 0xba)
		test.cpu.swapAddr(0x8641, actual)
		expected := mem.NewMemory(mem.NewHardwareRegisters(nil, nil), "")
		expected.Write(0x8641, 0xab)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}
