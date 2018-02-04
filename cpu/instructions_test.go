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

func Testcall(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.call()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testccf(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ccf()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testcp(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.cp()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestcpAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.cpAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testcpl(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.cpl()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testdaa(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.daa()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testdec(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.dec()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testdec16(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.dec16()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestdecSP(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.decSP()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestdecAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.decAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testdi(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.di()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testei(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ei()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testhalt(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.halt()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testinc(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.inc()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testinc16(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.inc16()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestincSP(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.incSP()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestincAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.incAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testjp(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.jp()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testjr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.jr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testld(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ld()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testld16(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ld16()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestldFromAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldFromAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestldToAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldToAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestldhFromAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldhFromAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestldhToAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldhToAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestldAFromAddrC(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldAFromAddrC()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestldAToAddrC(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldAToAddrC()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestldSP(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldSP()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestldHLToSP(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldHLToSP()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestldSPToAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldSPToAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestldSPToHL(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldSPToHL()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestlddFromAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.lddFromAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestlddToAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.lddToAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestldiFromAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldiFromAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestldiToAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ldiToAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testnop(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.nop()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testor(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.or()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestorAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.orAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testpop(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.pop()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testpush(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.push()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testres(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.res()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestresAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.resAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testret(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.ret()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testreti(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.reti()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testrl(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.rl()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestrlAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.rlAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testrla(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.rla()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testrlc(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.rlc()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testrlca(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.rlca()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestrlcAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.rlcAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testrr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.rr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testrra(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.rra()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestrrAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.rrAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testrrc(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.rrc()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testrrca(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.rrca()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestrrcAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.rrcAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testrst(t *testing.T) {
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

func Testsla(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.sla()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestslaAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.slaAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testsra(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.sra()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestsraAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.sraAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testsrl(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.srl()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestsrlAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.srlAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testswap(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.swap()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestswapAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.swapAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testscf(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.scf()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Teststop(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.stop()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testsbc(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.sbc()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestsbcAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.sbcAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testsub(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.sub()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestsubAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.subAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func Testxor(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.xor()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}

func TestxorAddr(t *testing.T) {
	for _, test := range []struct{ cpu, expectedCPU CPU }{
		{CPU{}, CPU{}},
	} {
		// test.cpu.xorAddr()
		compareCPUs(t, &test.expectedCPU, &test.cpu, nil)
	}
}
