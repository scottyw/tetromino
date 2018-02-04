package cpu

import (
	"fmt"
	"reflect"
	"testing"
)

////////////////////////////////////////////////////////////////
// Testable memory implementation
////////////////////////////////////////////////////////////////

type testableMemory struct {
	actual   map[uint16]*byte
	expected map[uint16]*byte
}

// Read a byte from the chosen memory location
func (mem testableMemory) Read(addr uint16) *byte {
	result, present := mem.actual[addr]
	if !present {
		b := byte(0)
		mem.actual[addr] = &b
		return &b
	}
	return result
}

// GenerateCrashReport to show memory state
func (mem testableMemory) GenerateCrashReport() {
	fmt.Println("TestMemory crash: ", mem.actual, mem.expected)
}

func bytePtr(b byte) *byte {
	return &b
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
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
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

func TestAddHL(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.addHL()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestAddSP(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.addSP()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestAnd(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.and()
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
		test.cpu.bit(2, test.cpu.c)
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXcall(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.call()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXccf(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ccf()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXcp(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.cp()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXcpAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.cpAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXcpl(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.cpl()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXdaa(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.daa()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXdec(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.dec()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXdec16(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.dec16()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXdecSP(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.decSP()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXdecAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.decAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXdi(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.di()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXei(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ei()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXhalt(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.halt()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXinc(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.inc()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXinc16(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.inc16()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXincSP(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.incSP()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXincAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.incAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXjp(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.jp()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXjr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.jr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXld(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ld()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXld16(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ld16()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXldFromAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldFromAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXldToAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldToAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXldhFromAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldhFromAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXldhToAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldhToAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXldAFromAddrC(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldAFromAddrC()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXldAToAddrC(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldAToAddrC()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXldSP(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldSP()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXldHLToSP(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldHLToSP()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXldSPToAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldSPToAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXldSPToHL(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldSPToHL()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXlddFromAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.lddFromAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXlddToAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.lddToAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXldiFromAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldiFromAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXldiToAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldiToAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXnop(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.nop()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXor(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.or()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXorAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.orAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXpop(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.pop()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXpush(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.push()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXres(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.res()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXresAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.resAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXret(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ret()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXreti(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.reti()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
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
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXrlAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.rlAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
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
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
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
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
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
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXrlcAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.rlcAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
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
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
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
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXrrAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.rrAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
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
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
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
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXrrcAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.rrcAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXrst(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.rst()
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

func TestXsla(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.sla()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXslaAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.slaAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXsra(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.sra()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXsraAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.sraAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXsrl(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.srl()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXsrlAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.srlAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXswap(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.swap()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXswapAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.swapAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXscf(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.scf()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXstop(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.stop()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXsbc(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.sbc()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXsbcAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.sbcAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXsub(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.sub()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXsubAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.subAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXxor(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.xor()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestXxorAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.xorAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}
