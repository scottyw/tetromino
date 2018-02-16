package cpu

import (
	"bytes"
	"testing"

	"github.com/scottyw/goomba/mem"
)

////////////////////////////////////////////////////////////////
// Test-agnostic validation functions
////////////////////////////////////////////////////////////////

func compareCPUsAndMemory(t *testing.T, expectedCPU, actualCPU *CPU, expectedMem, actualMem mem.Memory, startAddr, length uint16) {
	compareCPUs(t, expectedCPU, actualCPU)
	actual := actualMem.ReadRegion(startAddr, length)
	expected := expectedMem.ReadRegion(startAddr, length)
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
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.addHL()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestAddSP(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.addSP()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestAnd(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.and()
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
		{CPU{pc: 0xabcd, sp: 0x1234}, CPU{pc: 0x1af2, sp: 0x1232}},
	} {
		actual := mem.NewMemory()
		test.cpu.call("", 0x1af2, actual)
		expected := mem.NewMemory()
		*expected.Read(0x1233) = 0xab
		*expected.Read(0x1234) = 0xcd
		compareCPUsAndMemory(t, &test.expectedCPU, &test.cpu, expected, actual, 0x1232, 0xf)
	}
}

func TestXccf(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ccf()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXcp(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.cp()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXcpAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.cpAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
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

func TestXdaa(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.daa()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXdec(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.dec()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXdec16(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.dec16()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXdecSP(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.decSP()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXdecAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.decAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXdi(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.di()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXei(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ei()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXhalt(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.halt()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXinc(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.inc()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXinc16(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.inc16()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXincSP(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.incSP()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXincAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.incAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXjp(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.jp()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXjr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.jr()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXld(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ld()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXld16(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ld16()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXldFromAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldFromAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXldToAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldToAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXldhFromAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldhFromAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXldhToAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldhToAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXldAFromAddrC(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldAFromAddrC()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXldAToAddrC(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldAToAddrC()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXldSP(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldSP()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXldHLToSP(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldHLToSP()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXldSPToAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldSPToAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXldSPToHL(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldSPToHL()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXlddFromAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.lddFromAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXlddToAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.lddToAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXldiFromAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldiFromAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXldiToAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldiToAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXnop(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.nop()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXor(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.or()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXorAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.orAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestPop(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{b: 0xff, c: 0x11, sp: 0x1232}, CPU{b: 0x1a, c: 0xf2, sp: 0x1234}},
	} {
		actual := mem.NewMemory()
		*actual.Read(0x1233) = 0x1a
		*actual.Read(0x1234) = 0xf2
		test.cpu.pop(test.cpu.bc(), actual)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestPush(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{b: 0x1a, c: 0xf2, sp: 0x1234}, CPU{b: 0x1a, c: 0xf2, sp: 0x1232}},
	} {
		actual := mem.NewMemory()
		test.cpu.push(test.cpu.bc(), actual)
		expected := mem.NewMemory()
		*expected.Read(0x1233) = 0x1a
		*expected.Read(0x1234) = 0xf2
		compareCPUsAndMemory(t, &test.expectedCPU, &test.cpu, expected, actual, 0x1232, 0xf)
	}
}

func TestXres(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.res()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXresAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.resAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestRet(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{pc: 0xabab, sp: 0x1232}, CPU{pc: 0x1af2, sp: 0x1234}},
	} {
		mem := mem.NewMemory()
		*mem.Read(0x1233) = 0x1a
		*mem.Read(0x1234) = 0xf2
		test.cpu.ret("", mem)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestReti(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{pc: 0xabab, sp: 0x1232}, CPU{pc: 0x1af2, sp: 0x1234, ime: true}},
	} {
		mem := mem.NewMemory()
		*mem.Read(0x1233) = 0x1a
		*mem.Read(0x1234) = 0xf2
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

func TestXrlAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.rlAddr()
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

func TestXrlcAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.rlcAddr()
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

func TestXrrAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.rrAddr()
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

func TestXrrcAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.rrcAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXrst(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.rst()
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
		test.cpu.set(test.pos, &test.cpu.l)
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXsla(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.sla()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXslaAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.slaAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXsra(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.sra()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXsraAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.sraAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXsrl(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.srl()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXsrlAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.srlAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
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
		actual := mem.NewMemory()
		*actual.Read(0x1233) = 0xba
		test.cpu.swapAddr(0x1233, actual)
		expected := mem.NewMemory()
		*expected.Read(0x1233) = 0xab
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXscf(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.scf()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXstop(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.stop()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXsbc(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.sbc()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXsbcAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.sbcAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXsub(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.sub()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXsubAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.subAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXxor(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.xor()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}

func TestXxorAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.xorAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu)
	}
}
