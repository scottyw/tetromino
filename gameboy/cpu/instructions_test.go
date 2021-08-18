package cpu

import (
	"testing"
)

////////////////////////////////////////////////////////////////
// Test-agnostic validation functions
////////////////////////////////////////////////////////////////

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
		{CPU{a: 0x1a, c: 0xf2}, CPU{a: 0x0c, c: 0xf2, f: 0x10}},
		{CPU{a: 0x1a, c: 0x2b}, CPU{a: 0x45, c: 0x2b, f: 0x20}},
		{CPU{a: 0x00, c: 0x00}, CPU{a: 0x00, c: 0x00, f: 0x80}},
		{CPU{a: 0x1a, c: 0x22, f: 0x10}, CPU{a: 0x3d, c: 0x22}},
		{CPU{a: 0x1a, c: 0xf2, f: 0x10}, CPU{a: 0x0d, c: 0xf2, f: 0x10}},
		{CPU{a: 0x1a, c: 0x2b, f: 0x10}, CPU{a: 0x46, c: 0x2b, f: 0x20}},
		{CPU{a: 0xff, c: 0x00, f: 0x10}, CPU{a: 0x00, c: 0x00, f: 0xb0}},
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
		{CPU{a: 0xa3}, CPU{a: 0x46, f: 0x10}},
		{CPU{a: 0x1a}, CPU{a: 0x34, f: 0x20}},
		{CPU{a: 0x00}, CPU{a: 0x00, f: 0x80}},
	} {
		test.cpu.add(test.cpu.a)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
	for _, test := range []struct {
		cpu, expectedCPU CPU
	}{
		{CPU{a: 0x1a, c: 0x22}, CPU{a: 0x3c, c: 0x22}},
		{CPU{a: 0x1a, c: 0xf2}, CPU{a: 0x0c, c: 0xf2, f: 0x10}},
		{CPU{a: 0x1a, c: 0x2b}, CPU{a: 0x45, c: 0x2b, f: 0x20}},
		{CPU{a: 0x00, c: 0x00}, CPU{a: 0x00, c: 0x00, f: 0x80}},
	} {
		test.cpu.add(test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestAddHL(t *testing.T) {
	for _, test := range []struct {
		cpu, expectedCPU CPU
	}{
		{CPU{h: 0x1f, l: 0xb2}, CPU{h: 0x51, l: 0xc4, f: 0x20}},
		{CPU{h: 0xd1, l: 0xb2}, CPU{h: 0x03, l: 0xc4, f: 0x10}},
		{CPU{h: 0xcd, l: 0xee}, CPU{h: 0x00, l: 0x00, f: 0x30}},
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
		{CPU{c: 0x04, f: 0x80}, CPU{c: 0x04, f: 0x20}},
		{CPU{c: 0xf0, f: 0x50}, CPU{c: 0xf0, f: 0xb0}},
	} {
		test.cpu.bit(2, &test.cpu.c)()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestCpl(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{a: 0xb1}, CPU{a: 0x4e, f: 0x60}},
	} {
		test.cpu.cpl()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestRes(t *testing.T) {
	// Flags: [- - - -]
	for _, test := range []struct {
		cpu, expectedCPU CPU
	}{
		{CPU{c: 0x0a}, CPU{c: 0x0a}},
		{CPU{c: 0x0e, f: 0xf0}, CPU{c: 0x0a, f: 0xf0}},
	} {
		test.cpu.res(2, &test.cpu.c)()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestRl(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{c: 0xa8}, CPU{c: 0x50, f: 0x10}},
		{CPU{c: 0xa8, f: 0x10}, CPU{c: 0x51, f: 0x10}},
		{CPU{c: 0x15}, CPU{c: 0x2a}},
		{CPU{c: 0x15, f: 0x10}, CPU{c: 0x2b}},
		{CPU{c: 0x00}, CPU{c: 0x00, f: 0x80}},
	} {
		test.cpu.rl(&test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestRla(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{a: 0xa8}, CPU{a: 0x50, f: 0x10}},
		{CPU{a: 0xa8, f: 0x10}, CPU{a: 0x51, f: 0x10}},
		{CPU{a: 0x15}, CPU{a: 0x2a}},
		{CPU{a: 0x15, f: 0x10}, CPU{a: 0x2b}},
		{CPU{a: 0x00}, CPU{a: 0x00}},
	} {
		test.cpu.rla()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestRlc(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{c: 0xa8}, CPU{c: 0x51, f: 0x10}},
		{CPU{c: 0xa8, f: 0x10}, CPU{c: 0x51, f: 0x10}},
		{CPU{c: 0x15}, CPU{c: 0x2a}},
		{CPU{c: 0x15, f: 0x10}, CPU{c: 0x2a}},
		{CPU{c: 0x00}, CPU{c: 0x00, f: 0x80}},
	} {
		test.cpu.rlc(&test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestRlca(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{a: 0xa8}, CPU{a: 0x51, f: 0x10}},
		{CPU{a: 0xa8, f: 0x10}, CPU{a: 0x51, f: 0x10}},
		{CPU{a: 0x15}, CPU{a: 0x2a}},
		{CPU{a: 0x15, f: 0x10}, CPU{a: 0x2a}},
		{CPU{a: 0x00}, CPU{a: 0x00}},
	} {
		test.cpu.rlca()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestRr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{c: 0x15}, CPU{c: 0x0a, f: 0x10}},
		{CPU{c: 0x15, f: 0x10}, CPU{c: 0x8a, f: 0x10}},
		{CPU{c: 0xa8}, CPU{c: 0x54}},
		{CPU{c: 0xa8, f: 0x10}, CPU{c: 0xd4}},
		{CPU{c: 0x00}, CPU{c: 0x00, f: 0x80}},
	} {
		test.cpu.rr(&test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestRra(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{a: 0x15}, CPU{a: 0x0a, f: 0x10}},
		{CPU{a: 0x15, f: 0x10}, CPU{a: 0x8a, f: 0x10}},
		{CPU{a: 0xa8}, CPU{a: 0x54}},
		{CPU{a: 0xa8, f: 0x10}, CPU{a: 0xd4}},
		{CPU{a: 0x00}, CPU{a: 0x00, f: 0x00}},
	} {
		test.cpu.rra()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestRrc(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{c: 0x15}, CPU{c: 0x8a, f: 0x10}},
		{CPU{c: 0x15, f: 0x10}, CPU{c: 0x8a, f: 0x10}},
		{CPU{c: 0xa8}, CPU{c: 0x54}},
		{CPU{c: 0xa8, f: 0x10}, CPU{c: 0x54}},
		{CPU{c: 0x00}, CPU{c: 0x00, f: 0x80}},
	} {
		test.cpu.rrc(&test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestRrca(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{a: 0x15}, CPU{a: 0x8a, f: 0x10}},
		{CPU{a: 0x15, f: 0x10}, CPU{a: 0x8a, f: 0x10}},
		{CPU{a: 0xa8}, CPU{a: 0x54}},
		{CPU{a: 0xa8, f: 0x10}, CPU{a: 0x54}},
		{CPU{a: 0x00}, CPU{a: 0x00, f: 0x00}},
	} {
		test.cpu.rrca()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
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
		test.cpu.set(test.pos, &test.cpu.l)()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestSla(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{c: 0xa9}, CPU{c: 0x52, f: 0x10}},
		{CPU{c: 0x15, f: 0x10}, CPU{c: 0x2a}},
		{CPU{c: 0x00}, CPU{c: 0x00, f: 0x80}},
	} {
		test.cpu.sla(&test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestSra(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{c: 0x55}, CPU{c: 0x2a, f: 0x10}},
		{CPU{c: 0xa8, f: 0x10}, CPU{c: 0xd4}},
		{CPU{c: 0x00}, CPU{c: 0x00, f: 0x80}},
	} {
		test.cpu.sra(&test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestSrl(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{c: 0x55}, CPU{c: 0x2a, f: 0x10}},
		{CPU{c: 0xa8, f: 0x10}, CPU{c: 0x54}},
		{CPU{c: 0x00}, CPU{c: 0x00, f: 0x80}},
	} {
		test.cpu.srl(&test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestSwap(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{c: 0x15, f: 0x10}, CPU{c: 0x51}},
		{CPU{c: 0x00}, CPU{c: 0x00, f: 0x80}},
	} {
		test.cpu.swap(&test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}
