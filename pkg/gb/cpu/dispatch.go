package cpu

import (
	"github.com/scottyw/tetromino/pkg/gb/mem"
)

// Dispatch determines how CPU instructions are dispatched
type Dispatch struct {
	cpu               *CPU
	mem               *mem.Memory
	hwr               *mem.HardwareRegisters
	normal            [256][]func()
	prefix            [256][]func()
	steps             *[]func()
	stepIndex         int
	altStepIndex      int
	handlingInterrupt bool
	Mooneye           bool
	u8                uint8
	u16               uint16
	mem8              uint8
	mem16             uint16
}

// NewDispatch returns a Dispatch instance bringing the CPU and memory together
func NewDispatch(cpu *CPU, mem *mem.Memory, hwr *mem.HardwareRegisters) *Dispatch {
	initialSteps := []func(){}
	dispatch := &Dispatch{
		cpu:   cpu,
		mem:   mem,
		hwr:   hwr,
		steps: &initialSteps,
	}
	dispatch.initialize(cpu, mem)
	return dispatch
}

// TestA returns the value of register a for test purposes
func (d *Dispatch) TestA() uint8 {
	return d.cpu.a
}

// Start the CPU again on button press
func (d *Dispatch) Start() {
	d.cpu.stopped = false
}

func nop() {
	// Do nothing
}

func readFirstParam() {
	// FIXME this function only exists to heep track of what steps comprise each instruction
}

func readSecondParam() {
	// FIXME this function only exists to heep track of what steps comprise each instruction
}

func readMemory() {
	// FIXME this function only exists to heep track of what steps comprise each instruction
}

func writeMemory() {
	// FIXME this function only exists to heep track of what steps comprise each instruction
}

func (d *Dispatch) initialize(cpu *CPU, mem *mem.Memory) {
	var normalExec = [256]func(){}
	var prefixExec = [256]func(){}

	d.normal[0x00] = []func(){nop} // NOP          1 [4]

	normalExec[0x01] = func() { cpu.ld16(&cpu.b, &cpu.c, d.u16) }                // LD BC d16 []
	d.normal[0x01] = []func(){readFirstParam, readSecondParam, normalExec[0x01]} // LD BC d16    3 [12]

	normalExec[0x02] = func() { cpu.ldA16U8(cpu.bc(), cpu.a, mem) } // LD (BC) A []
	d.normal[0x02] = []func(){normalExec[0x02], writeMemory}        // LD (BC) A    1 [8]

	normalExec[0x03] = func() { cpu.inc16(&cpu.b, &cpu.c) } // INC BC  []
	d.normal[0x03] = []func(){nop, normalExec[0x03]}        // INC BC       1 [8]

	normalExec[0x04] = func() { cpu.inc(&cpu.b) } // INC B  [Z 0 H -]
	d.normal[0x04] = []func(){normalExec[0x04]}   // INC B        1 [4]

	normalExec[0x05] = func() { cpu.dec(&cpu.b) } // DEC B  [Z 1 H -]
	d.normal[0x05] = []func(){normalExec[0x05]}   // DEC B        1 [4]

	normalExec[0x06] = func() { cpu.ld(&cpu.b, d.u8) }          // LD B d8 []
	d.normal[0x06] = []func(){readFirstParam, normalExec[0x06]} // LD B d8      2 [8]

	normalExec[0x07] = func() { cpu.rlca() }    // RLCA   [0 0 0 C]
	d.normal[0x07] = []func(){normalExec[0x07]} // RLCA         1 [4]

	normalExec[0x08] = func() { cpu.ldA16SP(d.u16, mem) }                                          // LD (a16) SP []
	d.normal[0x08] = []func(){readFirstParam, readSecondParam, nop, normalExec[0x08], writeMemory} // LD (a16) SP  3 [20]

	normalExec[0x09] = func() { cpu.addHL(cpu.bc()) } // ADD HL BC [- 0 H C]
	d.normal[0x09] = []func(){nop, normalExec[0x09]}  // ADD HL BC    1 [8]

	normalExec[0x0a] = func() { cpu.ldR8A16(&cpu.a, cpu.bc(), mem) } // LD A (BC) []
	d.normal[0x0a] = []func(){readMemory, normalExec[0x0a]}          // LD A (BC)    1 [8]

	normalExec[0x0b] = func() { cpu.dec16(&cpu.b, &cpu.c) } // DEC BC  []
	d.normal[0x0b] = []func(){nop, normalExec[0x0b]}        // DEC BC       1 [8]

	normalExec[0x0c] = func() { cpu.inc(&cpu.c) } // INC C  [Z 0 H -]
	d.normal[0x0c] = []func(){normalExec[0x0c]}   // INC C        1 [4]

	normalExec[0x0d] = func() { cpu.dec(&cpu.c) } // DEC C  [Z 1 H -]
	d.normal[0x0d] = []func(){normalExec[0x0d]}   // DEC C        1 [4]

	normalExec[0x0e] = func() { cpu.ld(&cpu.c, d.u8) }          // LD C d8 []
	d.normal[0x0e] = []func(){readFirstParam, normalExec[0x0e]} // LD C d8      2 [8]

	normalExec[0x0f] = func() { cpu.rrca() }    // RRCA   [0 0 0 C]
	d.normal[0x0f] = []func(){normalExec[0x0f]} // RRCA         1 [4]

	normalExec[0x10] = func() { cpu.stop() }    // STOP 0  []
	d.normal[0x10] = []func(){normalExec[0x10]} // STOP 0       1 [4]

	normalExec[0x11] = func() { cpu.ld16(&cpu.d, &cpu.e, d.u16) }                // LD DE d16 []
	d.normal[0x11] = []func(){readFirstParam, readSecondParam, normalExec[0x11]} // LD DE d16    3 [12]

	normalExec[0x12] = func() { cpu.ldA16U8(cpu.de(), cpu.a, mem) } // LD (DE) A []
	d.normal[0x12] = []func(){normalExec[0x12], writeMemory}        // LD (DE) A    1 [8]

	normalExec[0x13] = func() { cpu.inc16(&cpu.d, &cpu.e) } // INC DE  []
	d.normal[0x13] = []func(){nop, normalExec[0x13]}        // INC DE       1 [8]

	normalExec[0x14] = func() { cpu.inc(&cpu.d) } // INC D  [Z 0 H -]
	d.normal[0x14] = []func(){normalExec[0x14]}   // INC D        1 [4]

	normalExec[0x15] = func() { cpu.dec(&cpu.d) } // DEC D  [Z 1 H -]
	d.normal[0x15] = []func(){normalExec[0x15]}   // DEC D        1 [4]

	normalExec[0x16] = func() { cpu.ld(&cpu.d, d.u8) }          // LD D d8 []
	d.normal[0x16] = []func(){readFirstParam, normalExec[0x16]} // LD D d8      2 [8]

	normalExec[0x17] = func() { cpu.rla() }     // RLA   [0 0 0 C]
	d.normal[0x17] = []func(){normalExec[0x17]} // RLA          1 [4]

	normalExec[0x18] = func() { cpu.jr("", int8(d.u8)) }             // JR r8  []
	d.normal[0x18] = []func(){readFirstParam, nop, normalExec[0x18]} // JR r8        2 [12]

	normalExec[0x19] = func() { cpu.addHL(cpu.de()) } // ADD HL DE [- 0 H C]
	d.normal[0x19] = []func(){nop, normalExec[0x19]}  // ADD HL DE    1 [8]

	normalExec[0x1a] = func() { cpu.ldR8A16(&cpu.a, cpu.de(), mem) } // LD A (DE) []
	d.normal[0x1a] = []func(){readMemory, normalExec[0x1a]}          // LD A (DE)    1 [8]

	normalExec[0x1b] = func() { cpu.dec16(&cpu.d, &cpu.e) } // DEC DE  []
	d.normal[0x1b] = []func(){nop, normalExec[0x1b]}        // DEC DE       1 [8]

	normalExec[0x1c] = func() { cpu.inc(&cpu.e) } // INC E  [Z 0 H -]
	d.normal[0x1c] = []func(){normalExec[0x1c]}   // INC E        1 [4]

	normalExec[0x1d] = func() { cpu.dec(&cpu.e) } // DEC E  [Z 1 H -]
	d.normal[0x1d] = []func(){normalExec[0x1d]}   // DEC E        1 [4]

	normalExec[0x1e] = func() { cpu.ld(&cpu.e, d.u8) }          // LD E d8 []
	d.normal[0x1e] = []func(){readFirstParam, normalExec[0x1e]} // LD E d8      2 [8]

	normalExec[0x1f] = func() { cpu.rra() }     // RRA   [0 0 0 C]
	d.normal[0x1f] = []func(){normalExec[0x1f]} // RRA          1 [4]

	normalExec[0x20] = func() { cpu.jr("NZ", int8(d.u8)) }           // JR NZ r8 []
	d.normal[0x20] = []func(){readFirstParam, normalExec[0x20], nop} // JR NZ r8     2 [12 8]

	normalExec[0x21] = func() { cpu.ld16(&cpu.h, &cpu.l, d.u16) }                // LD HL d16 []
	d.normal[0x21] = []func(){readFirstParam, readSecondParam, normalExec[0x21]} // LD HL d16    3 [12]

	normalExec[0x22] = func() { cpu.ldiA16A(mem) }           // LD (HL+) A []
	d.normal[0x22] = []func(){normalExec[0x22], writeMemory} // LD (HL+) A   1 [8]

	normalExec[0x23] = func() { cpu.inc16(&cpu.h, &cpu.l) } // INC HL  []
	d.normal[0x23] = []func(){nop, normalExec[0x23]}        // INC HL       1 [8]

	normalExec[0x24] = func() { cpu.inc(&cpu.h) } // INC H  [Z 0 H -]
	d.normal[0x24] = []func(){normalExec[0x24]}   // INC H        1 [4]

	normalExec[0x25] = func() { cpu.dec(&cpu.h) } // DEC H  [Z 1 H -]
	d.normal[0x25] = []func(){normalExec[0x25]}   // DEC H        1 [4]

	normalExec[0x26] = func() { cpu.ld(&cpu.h, d.u8) }          // LD H d8 []
	d.normal[0x26] = []func(){readFirstParam, normalExec[0x26]} // LD H d8      2 [8]

	normalExec[0x27] = func() { cpu.daa() }     // DAA   [Z - 0 C]
	d.normal[0x27] = []func(){normalExec[0x27]} // DAA          1 [4]

	normalExec[0x28] = func() { cpu.jr("Z", int8(d.u8)) }            // JR Z r8 []
	d.normal[0x28] = []func(){readFirstParam, normalExec[0x28], nop} // JR Z r8      2 [12 8]

	normalExec[0x29] = func() { cpu.addHL(cpu.hl()) } // ADD HL HL [- 0 H C]
	d.normal[0x29] = []func(){nop, normalExec[0x29]}  // ADD HL HL    1 [8]

	normalExec[0x2a] = func() { cpu.ldiAA16(mem) }          // LD A (HL+) []
	d.normal[0x2a] = []func(){readMemory, normalExec[0x2a]} // LD A (HL+)   1 [8]

	normalExec[0x2b] = func() { cpu.dec16(&cpu.h, &cpu.l) } // DEC HL  []
	d.normal[0x2b] = []func(){nop, normalExec[0x2b]}        // DEC HL       1 [8]

	normalExec[0x2c] = func() { cpu.inc(&cpu.l) } // INC L  [Z 0 H -]
	d.normal[0x2c] = []func(){normalExec[0x2c]}   // INC L        1 [4]

	normalExec[0x2d] = func() { cpu.dec(&cpu.l) } // DEC L  [Z 1 H -]
	d.normal[0x2d] = []func(){normalExec[0x2d]}   // DEC L        1 [4]

	normalExec[0x2e] = func() { cpu.ld(&cpu.l, d.u8) }          // LD L d8 []
	d.normal[0x2e] = []func(){readFirstParam, normalExec[0x2e]} // LD L d8      2 [8]

	normalExec[0x2f] = func() { cpu.cpl() }     // CPL   [- 1 1 -]
	d.normal[0x2f] = []func(){normalExec[0x2f]} // CPL          1 [4]

	normalExec[0x30] = func() { cpu.jr("NC", int8(d.u8)) }           // JR NC r8 []
	d.normal[0x30] = []func(){readFirstParam, normalExec[0x30], nop} // JR NC r8     2 [12 8]

	normalExec[0x31] = func() { cpu.ldSP(d.u16) }                                // LD SP d16 []
	d.normal[0x31] = []func(){readFirstParam, readSecondParam, normalExec[0x31]} // LD SP d16    3 [12]

	normalExec[0x32] = func() { cpu.lddA16A(mem) }           // LD (HL-) A []
	d.normal[0x32] = []func(){normalExec[0x32], writeMemory} // LD (HL-) A   1 [8]

	normalExec[0x33] = func() { cpu.incSP() }        // INC SP  []
	d.normal[0x33] = []func(){nop, normalExec[0x33]} // INC SP       1 [8]

	normalExec[0x34] = func() { cpu.incAddr(cpu.hl(), mem) }             // INC (HL)  [Z 0 H -]
	d.normal[0x34] = []func(){readMemory, normalExec[0x34], writeMemory} // INC (HL)     1 [12]

	normalExec[0x35] = func() { cpu.decAddr(cpu.hl(), mem) }             // DEC (HL)  [Z 1 H -]
	d.normal[0x35] = []func(){readMemory, normalExec[0x35], writeMemory} // DEC (HL)     1 [12]

	normalExec[0x36] = func() { cpu.ldA16U8(cpu.hl(), d.u8, mem) }           // LD (HL) d8 []
	d.normal[0x36] = []func(){readFirstParam, normalExec[0x36], writeMemory} // LD (HL) d8   2 [12]

	normalExec[0x37] = func() { cpu.scf() }     // SCF   [- 0 0 1]
	d.normal[0x37] = []func(){normalExec[0x37]} // SCF          1 [4]

	normalExec[0x38] = func() { cpu.jr("C", int8(d.u8)) }            // JR C r8 []
	d.normal[0x38] = []func(){readFirstParam, normalExec[0x38], nop} // JR C r8      2 [12 8]

	normalExec[0x39] = func() { cpu.addHL(cpu.sp) }  // ADD HL SP [- 0 H C]
	d.normal[0x39] = []func(){nop, normalExec[0x39]} // ADD HL SP    1 [8]

	normalExec[0x3a] = func() { cpu.lddAA16(mem) }          // LD A (HL-) []
	d.normal[0x3a] = []func(){readMemory, normalExec[0x3a]} // LD A (HL-)   1 [8]

	normalExec[0x3b] = func() { cpu.decSP() }        // DEC SP  []
	d.normal[0x3b] = []func(){nop, normalExec[0x3b]} // DEC SP       1 [8]

	normalExec[0x3c] = func() { cpu.inc(&cpu.a) } // INC A  [Z 0 H -]
	d.normal[0x3c] = []func(){normalExec[0x3c]}   // INC A        1 [4]

	normalExec[0x3d] = func() { cpu.dec(&cpu.a) } // DEC A  [Z 1 H -]
	d.normal[0x3d] = []func(){normalExec[0x3d]}   // DEC A        1 [4]

	normalExec[0x3e] = func() { cpu.ld(&cpu.a, d.u8) }          // LD A d8 []
	d.normal[0x3e] = []func(){readFirstParam, normalExec[0x3e]} // LD A d8      2 [8]

	normalExec[0x3f] = func() { cpu.ccf() }     // CCF   [- 0 0 C]
	d.normal[0x3f] = []func(){normalExec[0x3f]} // CCF          1 [4]

	normalExec[0x40] = func() { cpu.ld(&cpu.b, cpu.b) } // LD B B []
	d.normal[0x40] = []func(){normalExec[0x40]}         // LD B B       1 [4]

	normalExec[0x41] = func() { cpu.ld(&cpu.b, cpu.c) } // LD B C []
	d.normal[0x41] = []func(){normalExec[0x41]}         // LD B C       1 [4]

	normalExec[0x42] = func() { cpu.ld(&cpu.b, cpu.d) } // LD B D []
	d.normal[0x42] = []func(){normalExec[0x42]}         // LD B D       1 [4]

	normalExec[0x43] = func() { cpu.ld(&cpu.b, cpu.e) } // LD B E []
	d.normal[0x43] = []func(){normalExec[0x43]}         // LD B E       1 [4]

	normalExec[0x44] = func() { cpu.ld(&cpu.b, cpu.h) } // LD B H []
	d.normal[0x44] = []func(){normalExec[0x44]}         // LD B H       1 [4]

	normalExec[0x45] = func() { cpu.ld(&cpu.b, cpu.l) } // LD B L []
	d.normal[0x45] = []func(){normalExec[0x45]}         // LD B L       1 [4]

	normalExec[0x46] = func() { cpu.ldR8A16(&cpu.b, cpu.hl(), mem) } // LD B (HL) []
	d.normal[0x46] = []func(){readMemory, normalExec[0x46]}          // LD B (HL)    1 [8]

	normalExec[0x47] = func() { cpu.ld(&cpu.b, cpu.a) } // LD B A []
	d.normal[0x47] = []func(){normalExec[0x47]}         // LD B A       1 [4]

	normalExec[0x48] = func() { cpu.ld(&cpu.c, cpu.b) } // LD C B []
	d.normal[0x48] = []func(){normalExec[0x48]}         // LD C B       1 [4]

	normalExec[0x49] = func() { cpu.ld(&cpu.c, cpu.c) } // LD C C []
	d.normal[0x49] = []func(){normalExec[0x49]}         // LD C C       1 [4]

	normalExec[0x4a] = func() { cpu.ld(&cpu.c, cpu.d) } // LD C D []
	d.normal[0x4a] = []func(){normalExec[0x4a]}         // LD C D       1 [4]

	normalExec[0x4b] = func() { cpu.ld(&cpu.c, cpu.e) } // LD C E []
	d.normal[0x4b] = []func(){normalExec[0x4b]}         // LD C E       1 [4]

	normalExec[0x4c] = func() { cpu.ld(&cpu.c, cpu.h) } // LD C H []
	d.normal[0x4c] = []func(){normalExec[0x4c]}         // LD C H       1 [4]

	normalExec[0x4d] = func() { cpu.ld(&cpu.c, cpu.l) } // LD C L []
	d.normal[0x4d] = []func(){normalExec[0x4d]}         // LD C L       1 [4]

	normalExec[0x4e] = func() { cpu.ldR8A16(&cpu.c, cpu.hl(), mem) } // LD C (HL) []
	d.normal[0x4e] = []func(){readMemory, normalExec[0x4e]}          // LD C (HL)    1 [8]

	normalExec[0x4f] = func() { cpu.ld(&cpu.c, cpu.a) } // LD C A []
	d.normal[0x4f] = []func(){normalExec[0x4f]}         // LD C A       1 [4]

	normalExec[0x50] = func() { cpu.ld(&cpu.d, cpu.b) } // LD D B []
	d.normal[0x50] = []func(){normalExec[0x50]}         // LD D B       1 [4]

	normalExec[0x51] = func() { cpu.ld(&cpu.d, cpu.c) } // LD D C []
	d.normal[0x51] = []func(){normalExec[0x51]}         // LD D C       1 [4]

	normalExec[0x52] = func() { cpu.ld(&cpu.d, cpu.d) } // LD D D []
	d.normal[0x52] = []func(){normalExec[0x52]}         // LD D D       1 [4]

	normalExec[0x53] = func() { cpu.ld(&cpu.d, cpu.e) } // LD D E []
	d.normal[0x53] = []func(){normalExec[0x53]}         // LD D E       1 [4]

	normalExec[0x54] = func() { cpu.ld(&cpu.d, cpu.h) } // LD D H []
	d.normal[0x54] = []func(){normalExec[0x54]}         // LD D H       1 [4]

	normalExec[0x55] = func() { cpu.ld(&cpu.d, cpu.l) } // LD D L []
	d.normal[0x55] = []func(){normalExec[0x55]}         // LD D L       1 [4]

	normalExec[0x56] = func() { cpu.ldR8A16(&cpu.d, cpu.hl(), mem) } // LD D (HL) []
	d.normal[0x56] = []func(){readMemory, normalExec[0x56]}          // LD D (HL)    1 [8]

	normalExec[0x57] = func() { cpu.ld(&cpu.d, cpu.a) } // LD D A []
	d.normal[0x57] = []func(){normalExec[0x57]}         // LD D A       1 [4]

	normalExec[0x58] = func() { cpu.ld(&cpu.e, cpu.b) } // LD E B []
	d.normal[0x58] = []func(){normalExec[0x58]}         // LD E B       1 [4]

	normalExec[0x59] = func() { cpu.ld(&cpu.e, cpu.c) } // LD E C []
	d.normal[0x59] = []func(){normalExec[0x59]}         // LD E C       1 [4]

	normalExec[0x5a] = func() { cpu.ld(&cpu.e, cpu.d) } // LD E D []
	d.normal[0x5a] = []func(){normalExec[0x5a]}         // LD E D       1 [4]

	normalExec[0x5b] = func() { cpu.ld(&cpu.e, cpu.e) } // LD E E []
	d.normal[0x5b] = []func(){normalExec[0x5b]}         // LD E E       1 [4]

	normalExec[0x5c] = func() { cpu.ld(&cpu.e, cpu.h) } // LD E H []
	d.normal[0x5c] = []func(){normalExec[0x5c]}         // LD E H       1 [4]

	normalExec[0x5d] = func() { cpu.ld(&cpu.e, cpu.l) } // LD E L []
	d.normal[0x5d] = []func(){normalExec[0x5d]}         // LD E L       1 [4]

	normalExec[0x5e] = func() { cpu.ldR8A16(&cpu.e, cpu.hl(), mem) } // LD E (HL) []
	d.normal[0x5e] = []func(){readMemory, normalExec[0x5e]}          // LD E (HL)    1 [8]

	normalExec[0x5f] = func() { cpu.ld(&cpu.e, cpu.a) } // LD E A []
	d.normal[0x5f] = []func(){normalExec[0x5f]}         // LD E A       1 [4]

	normalExec[0x60] = func() { cpu.ld(&cpu.h, cpu.b) } // LD H B []
	d.normal[0x60] = []func(){normalExec[0x60]}         // LD H B       1 [4]

	normalExec[0x61] = func() { cpu.ld(&cpu.h, cpu.c) } // LD H C []
	d.normal[0x61] = []func(){normalExec[0x61]}         // LD H C       1 [4]

	normalExec[0x62] = func() { cpu.ld(&cpu.h, cpu.d) } // LD H D []
	d.normal[0x62] = []func(){normalExec[0x62]}         // LD H D       1 [4]

	normalExec[0x63] = func() { cpu.ld(&cpu.h, cpu.e) } // LD H E []
	d.normal[0x63] = []func(){normalExec[0x63]}         // LD H E       1 [4]

	normalExec[0x64] = func() { cpu.ld(&cpu.h, cpu.h) } // LD H H []
	d.normal[0x64] = []func(){normalExec[0x64]}         // LD H H       1 [4]

	normalExec[0x65] = func() { cpu.ld(&cpu.h, cpu.l) } // LD H L []
	d.normal[0x65] = []func(){normalExec[0x65]}         // LD H L       1 [4]

	normalExec[0x66] = func() { cpu.ldR8A16(&cpu.h, cpu.hl(), mem) } // LD H (HL) []
	d.normal[0x66] = []func(){readMemory, normalExec[0x66]}          // LD H (HL)    1 [8]

	normalExec[0x67] = func() { cpu.ld(&cpu.h, cpu.a) } // LD H A []
	d.normal[0x67] = []func(){normalExec[0x67]}         // LD H A       1 [4]

	normalExec[0x68] = func() { cpu.ld(&cpu.l, cpu.b) } // LD L B []
	d.normal[0x68] = []func(){normalExec[0x68]}         // LD L B       1 [4]

	normalExec[0x69] = func() { cpu.ld(&cpu.l, cpu.c) } // LD L C []
	d.normal[0x69] = []func(){normalExec[0x69]}         // LD L C       1 [4]

	normalExec[0x6a] = func() { cpu.ld(&cpu.l, cpu.d) } // LD L D []
	d.normal[0x6a] = []func(){normalExec[0x6a]}         // LD L D       1 [4]

	normalExec[0x6b] = func() { cpu.ld(&cpu.l, cpu.e) } // LD L E []
	d.normal[0x6b] = []func(){normalExec[0x6b]}         // LD L E       1 [4]

	normalExec[0x6c] = func() { cpu.ld(&cpu.l, cpu.h) } // LD L H []
	d.normal[0x6c] = []func(){normalExec[0x6c]}         // LD L H       1 [4]

	normalExec[0x6d] = func() { cpu.ld(&cpu.l, cpu.l) } // LD L L []
	d.normal[0x6d] = []func(){normalExec[0x6d]}         // LD L L       1 [4]

	normalExec[0x6e] = func() { cpu.ldR8A16(&cpu.l, cpu.hl(), mem) } // LD L (HL) []
	d.normal[0x6e] = []func(){readMemory, normalExec[0x6e]}          // LD L (HL)    1 [8]

	normalExec[0x6f] = func() { cpu.ld(&cpu.l, cpu.a) } // LD L A []
	d.normal[0x6f] = []func(){normalExec[0x6f]}         // LD L A       1 [4]

	normalExec[0x70] = func() { cpu.ldA16U8(cpu.hl(), cpu.b, mem) } // LD (HL) B []
	d.normal[0x70] = []func(){normalExec[0x70], writeMemory}        // LD (HL) B    1 [8]

	normalExec[0x71] = func() { cpu.ldA16U8(cpu.hl(), cpu.c, mem) } // LD (HL) C []
	d.normal[0x71] = []func(){normalExec[0x71], writeMemory}        // LD (HL) C    1 [8]

	normalExec[0x72] = func() { cpu.ldA16U8(cpu.hl(), cpu.d, mem) } // LD (HL) D []
	d.normal[0x72] = []func(){normalExec[0x72], writeMemory}        // LD (HL) D    1 [8]

	normalExec[0x73] = func() { cpu.ldA16U8(cpu.hl(), cpu.e, mem) } // LD (HL) E []
	d.normal[0x73] = []func(){normalExec[0x73], writeMemory}        // LD (HL) E    1 [8]

	normalExec[0x74] = func() { cpu.ldA16U8(cpu.hl(), cpu.h, mem) } // LD (HL) H []
	d.normal[0x74] = []func(){normalExec[0x74], writeMemory}        // LD (HL) H    1 [8]

	normalExec[0x75] = func() { cpu.ldA16U8(cpu.hl(), cpu.l, mem) } // LD (HL) L []
	d.normal[0x75] = []func(){normalExec[0x75], writeMemory}        // LD (HL) L    1 [8]

	normalExec[0x76] = func() { cpu.halt() }    // HALT   []
	d.normal[0x76] = []func(){normalExec[0x76]} // HALT         1 [4]

	normalExec[0x77] = func() { cpu.ldA16U8(cpu.hl(), cpu.a, mem) } // LD (HL) A []
	d.normal[0x77] = []func(){normalExec[0x77], writeMemory}        // LD (HL) A    1 [8]

	normalExec[0x78] = func() { cpu.ld(&cpu.a, cpu.b) } // LD A B []
	d.normal[0x78] = []func(){normalExec[0x78]}         // LD A B       1 [4]

	normalExec[0x79] = func() { cpu.ld(&cpu.a, cpu.c) } // LD A C []
	d.normal[0x79] = []func(){normalExec[0x79]}         // LD A C       1 [4]

	normalExec[0x7a] = func() { cpu.ld(&cpu.a, cpu.d) } // LD A D []
	d.normal[0x7a] = []func(){normalExec[0x7a]}         // LD A D       1 [4]

	normalExec[0x7b] = func() { cpu.ld(&cpu.a, cpu.e) } // LD A E []
	d.normal[0x7b] = []func(){normalExec[0x7b]}         // LD A E       1 [4]

	normalExec[0x7c] = func() { cpu.ld(&cpu.a, cpu.h) } // LD A H []
	d.normal[0x7c] = []func(){normalExec[0x7c]}         // LD A H       1 [4]

	normalExec[0x7d] = func() { cpu.ld(&cpu.a, cpu.l) } // LD A L []
	d.normal[0x7d] = []func(){normalExec[0x7d]}         // LD A L       1 [4]

	normalExec[0x7e] = func() { cpu.ldR8A16(&cpu.a, cpu.hl(), mem) } // LD A (HL) []
	d.normal[0x7e] = []func(){readMemory, normalExec[0x7e]}          // LD A (HL)    1 [8]

	normalExec[0x7f] = func() { cpu.ld(&cpu.a, cpu.a) } // LD A A []
	d.normal[0x7f] = []func(){normalExec[0x7f]}         // LD A A       1 [4]

	normalExec[0x80] = func() { cpu.add(cpu.b) } // ADD A B [Z 0 H C]
	d.normal[0x80] = []func(){normalExec[0x80]}  // ADD A B      1 [4]

	normalExec[0x81] = func() { cpu.add(cpu.c) } // ADD A C [Z 0 H C]
	d.normal[0x81] = []func(){normalExec[0x81]}  // ADD A C      1 [4]

	normalExec[0x82] = func() { cpu.add(cpu.d) } // ADD A D [Z 0 H C]
	d.normal[0x82] = []func(){normalExec[0x82]}  // ADD A D      1 [4]

	normalExec[0x83] = func() { cpu.add(cpu.e) } // ADD A E [Z 0 H C]
	d.normal[0x83] = []func(){normalExec[0x83]}  // ADD A E      1 [4]

	normalExec[0x84] = func() { cpu.add(cpu.h) } // ADD A H [Z 0 H C]
	d.normal[0x84] = []func(){normalExec[0x84]}  // ADD A H      1 [4]

	normalExec[0x85] = func() { cpu.add(cpu.l) } // ADD A L [Z 0 H C]
	d.normal[0x85] = []func(){normalExec[0x85]}  // ADD A L      1 [4]

	normalExec[0x86] = func() { cpu.addAddr(cpu.hl(), mem) } // ADD A (HL) [Z 0 H C]
	d.normal[0x86] = []func(){readMemory, normalExec[0x86]}  // ADD A (HL)   1 [8]

	normalExec[0x87] = func() { cpu.add(cpu.a) } // ADD A A [Z 0 H C]
	d.normal[0x87] = []func(){normalExec[0x87]}  // ADD A A      1 [4]

	normalExec[0x88] = func() { cpu.adc(cpu.b) } // ADC A B [Z 0 H C]
	d.normal[0x88] = []func(){normalExec[0x88]}  // ADC A B      1 [4]

	normalExec[0x89] = func() { cpu.adc(cpu.c) } // ADC A C [Z 0 H C]
	d.normal[0x89] = []func(){normalExec[0x89]}  // ADC A C      1 [4]

	normalExec[0x8a] = func() { cpu.adc(cpu.d) } // ADC A D [Z 0 H C]
	d.normal[0x8a] = []func(){normalExec[0x8a]}  // ADC A D      1 [4]

	normalExec[0x8b] = func() { cpu.adc(cpu.e) } // ADC A E [Z 0 H C]
	d.normal[0x8b] = []func(){normalExec[0x8b]}  // ADC A E      1 [4]

	normalExec[0x8c] = func() { cpu.adc(cpu.h) } // ADC A H [Z 0 H C]
	d.normal[0x8c] = []func(){normalExec[0x8c]}  // ADC A H      1 [4]

	normalExec[0x8d] = func() { cpu.adc(cpu.l) } // ADC A L [Z 0 H C]
	d.normal[0x8d] = []func(){normalExec[0x8d]}  // ADC A L      1 [4]

	normalExec[0x8e] = func() { cpu.adcAddr(cpu.hl(), mem) } // ADC A (HL) [Z 0 H C]
	d.normal[0x8e] = []func(){readMemory, normalExec[0x8e]}  // ADC A (HL)   1 [8]

	normalExec[0x8f] = func() { cpu.adc(cpu.a) } // ADC A A [Z 0 H C]
	d.normal[0x8f] = []func(){normalExec[0x8f]}  // ADC A A      1 [4]

	normalExec[0x90] = func() { cpu.sub(cpu.b) } // SUB B  [Z 1 H C]
	d.normal[0x90] = []func(){normalExec[0x90]}  // SUB B        1 [4]

	normalExec[0x91] = func() { cpu.sub(cpu.c) } // SUB C  [Z 1 H C]
	d.normal[0x91] = []func(){normalExec[0x91]}  // SUB C        1 [4]

	normalExec[0x92] = func() { cpu.sub(cpu.d) } // SUB D  [Z 1 H C]
	d.normal[0x92] = []func(){normalExec[0x92]}  // SUB D        1 [4]

	normalExec[0x93] = func() { cpu.sub(cpu.e) } // SUB E  [Z 1 H C]
	d.normal[0x93] = []func(){normalExec[0x93]}  // SUB E        1 [4]

	normalExec[0x94] = func() { cpu.sub(cpu.h) } // SUB H  [Z 1 H C]
	d.normal[0x94] = []func(){normalExec[0x94]}  // SUB H        1 [4]

	normalExec[0x95] = func() { cpu.sub(cpu.l) } // SUB L  [Z 1 H C]
	d.normal[0x95] = []func(){normalExec[0x95]}  // SUB L        1 [4]

	normalExec[0x96] = func() { cpu.subAddr(cpu.hl(), mem) } // SUB (HL)  [Z 1 H C]
	d.normal[0x96] = []func(){readMemory, normalExec[0x96]}  // SUB (HL)     1 [8]

	normalExec[0x97] = func() { cpu.sub(cpu.a) } // SUB A  [Z 1 H C]
	d.normal[0x97] = []func(){normalExec[0x97]}  // SUB A        1 [4]

	normalExec[0x98] = func() { cpu.sbc(cpu.b) } // SBC A B [Z 1 H C]
	d.normal[0x98] = []func(){normalExec[0x98]}  // SBC A B      1 [4]

	normalExec[0x99] = func() { cpu.sbc(cpu.c) } // SBC A C [Z 1 H C]
	d.normal[0x99] = []func(){normalExec[0x99]}  // SBC A C      1 [4]

	normalExec[0x9a] = func() { cpu.sbc(cpu.d) } // SBC A D [Z 1 H C]
	d.normal[0x9a] = []func(){normalExec[0x9a]}  // SBC A D      1 [4]

	normalExec[0x9b] = func() { cpu.sbc(cpu.e) } // SBC A E [Z 1 H C]
	d.normal[0x9b] = []func(){normalExec[0x9b]}  // SBC A E      1 [4]

	normalExec[0x9c] = func() { cpu.sbc(cpu.h) } // SBC A H [Z 1 H C]
	d.normal[0x9c] = []func(){normalExec[0x9c]}  // SBC A H      1 [4]

	normalExec[0x9d] = func() { cpu.sbc(cpu.l) } // SBC A L [Z 1 H C]
	d.normal[0x9d] = []func(){normalExec[0x9d]}  // SBC A L      1 [4]

	normalExec[0x9e] = func() { cpu.sbcAddr(cpu.hl(), mem) } // SBC A (HL) [Z 1 H C]
	d.normal[0x9e] = []func(){readMemory, normalExec[0x9e]}  // SBC A (HL)   1 [8]

	normalExec[0x9f] = func() { cpu.sbc(cpu.a) } // SBC A A [Z 1 H C]
	d.normal[0x9f] = []func(){normalExec[0x9f]}  // SBC A A      1 [4]

	normalExec[0xa0] = func() { cpu.and(cpu.b) } // AND B  [Z 0 1 0]
	d.normal[0xa0] = []func(){normalExec[0xa0]}  // AND B        1 [4]

	normalExec[0xa1] = func() { cpu.and(cpu.c) } // AND C  [Z 0 1 0]
	d.normal[0xa1] = []func(){normalExec[0xa1]}  // AND C        1 [4]

	normalExec[0xa2] = func() { cpu.and(cpu.d) } // AND D  [Z 0 1 0]
	d.normal[0xa2] = []func(){normalExec[0xa2]}  // AND D        1 [4]

	normalExec[0xa3] = func() { cpu.and(cpu.e) } // AND E  [Z 0 1 0]
	d.normal[0xa3] = []func(){normalExec[0xa3]}  // AND E        1 [4]

	normalExec[0xa4] = func() { cpu.and(cpu.h) } // AND H  [Z 0 1 0]
	d.normal[0xa4] = []func(){normalExec[0xa4]}  // AND H        1 [4]

	normalExec[0xa5] = func() { cpu.and(cpu.l) } // AND L  [Z 0 1 0]
	d.normal[0xa5] = []func(){normalExec[0xa5]}  // AND L        1 [4]

	normalExec[0xa6] = func() { cpu.andAddr(cpu.hl(), mem) } // AND (HL)  [Z 0 1 0]
	d.normal[0xa6] = []func(){readMemory, normalExec[0xa6]}  // AND (HL)     1 [8]

	normalExec[0xa7] = func() { cpu.and(cpu.a) } // AND A  [Z 0 1 0]
	d.normal[0xa7] = []func(){normalExec[0xa7]}  // AND A        1 [4]

	normalExec[0xa8] = func() { cpu.xor(cpu.b) } // XOR B  [Z 0 0 0]
	d.normal[0xa8] = []func(){normalExec[0xa8]}  // XOR B        1 [4]

	normalExec[0xa9] = func() { cpu.xor(cpu.c) } // XOR C  [Z 0 0 0]
	d.normal[0xa9] = []func(){normalExec[0xa9]}  // XOR C        1 [4]

	normalExec[0xaa] = func() { cpu.xor(cpu.d) } // XOR D  [Z 0 0 0]
	d.normal[0xaa] = []func(){normalExec[0xaa]}  // XOR D        1 [4]

	normalExec[0xab] = func() { cpu.xor(cpu.e) } // XOR E  [Z 0 0 0]
	d.normal[0xab] = []func(){normalExec[0xab]}  // XOR E        1 [4]

	normalExec[0xac] = func() { cpu.xor(cpu.h) } // XOR H  [Z 0 0 0]
	d.normal[0xac] = []func(){normalExec[0xac]}  // XOR H        1 [4]

	normalExec[0xad] = func() { cpu.xor(cpu.l) } // XOR L  [Z 0 0 0]
	d.normal[0xad] = []func(){normalExec[0xad]}  // XOR L        1 [4]

	normalExec[0xae] = func() { cpu.xorAddr(cpu.hl(), mem) } // XOR (HL)  [Z 0 0 0]
	d.normal[0xae] = []func(){readMemory, normalExec[0xae]}  // XOR (HL)     1 [8]

	normalExec[0xaf] = func() { cpu.xor(cpu.a) } // XOR A  [Z 0 0 0]
	d.normal[0xaf] = []func(){normalExec[0xaf]}  // XOR A        1 [4]

	normalExec[0xb0] = func() { cpu.or(cpu.b) } // OR B  [Z 0 0 0]
	d.normal[0xb0] = []func(){normalExec[0xb0]} // OR B         1 [4]

	normalExec[0xb1] = func() { cpu.or(cpu.c) } // OR C  [Z 0 0 0]
	d.normal[0xb1] = []func(){normalExec[0xb1]} // OR C         1 [4]

	normalExec[0xb2] = func() { cpu.or(cpu.d) } // OR D  [Z 0 0 0]
	d.normal[0xb2] = []func(){normalExec[0xb2]} // OR D         1 [4]

	normalExec[0xb3] = func() { cpu.or(cpu.e) } // OR E  [Z 0 0 0]
	d.normal[0xb3] = []func(){normalExec[0xb3]} // OR E         1 [4]

	normalExec[0xb4] = func() { cpu.or(cpu.h) } // OR H  [Z 0 0 0]
	d.normal[0xb4] = []func(){normalExec[0xb4]} // OR H         1 [4]

	normalExec[0xb5] = func() { cpu.or(cpu.l) } // OR L  [Z 0 0 0]
	d.normal[0xb5] = []func(){normalExec[0xb5]} // OR L         1 [4]

	normalExec[0xb6] = func() { cpu.orAddr(cpu.hl(), mem) } // OR (HL)  [Z 0 0 0]
	d.normal[0xb6] = []func(){readMemory, normalExec[0xb6]} // OR (HL)      1 [8]

	normalExec[0xb7] = func() { cpu.or(cpu.a) } // OR A  [Z 0 0 0]
	d.normal[0xb7] = []func(){normalExec[0xb7]} // OR A         1 [4]

	normalExec[0xb8] = func() { cpu.cp(cpu.b) } // CP B  [Z 1 H C]
	d.normal[0xb8] = []func(){normalExec[0xb8]} // CP B         1 [4]

	normalExec[0xb9] = func() { cpu.cp(cpu.c) } // CP C  [Z 1 H C]
	d.normal[0xb9] = []func(){normalExec[0xb9]} // CP C         1 [4]

	normalExec[0xba] = func() { cpu.cp(cpu.d) } // CP D  [Z 1 H C]
	d.normal[0xba] = []func(){normalExec[0xba]} // CP D         1 [4]

	normalExec[0xbb] = func() { cpu.cp(cpu.e) } // CP E  [Z 1 H C]
	d.normal[0xbb] = []func(){normalExec[0xbb]} // CP E         1 [4]

	normalExec[0xbc] = func() { cpu.cp(cpu.h) } // CP H  [Z 1 H C]
	d.normal[0xbc] = []func(){normalExec[0xbc]} // CP H         1 [4]

	normalExec[0xbd] = func() { cpu.cp(cpu.l) } // CP L  [Z 1 H C]
	d.normal[0xbd] = []func(){normalExec[0xbd]} // CP L         1 [4]

	normalExec[0xbe] = func() { cpu.cpAddr(cpu.hl(), mem) } // CP (HL)  [Z 1 H C]
	d.normal[0xbe] = []func(){readMemory, normalExec[0xbe]} // CP (HL)      1 [8]

	normalExec[0xbf] = func() { cpu.cp(cpu.a) } // CP A  [Z 1 H C]
	d.normal[0xbf] = []func(){normalExec[0xbf]} // CP A         1 [4]

	normalExec[0xc0] = func() { cpu.ret("NZ", mem) }                // RET NZ  []
	d.normal[0xc0] = []func(){nop, normalExec[0xc0], nop, nop, nop} // RET NZ       1 [20 8]

	normalExec[0xc1] = func() { cpu.pop(&cpu.b, &cpu.c, mem) } // POP BC  []
	d.normal[0xc1] = []func(){nop, nop, normalExec[0xc1]}      // POP BC       1 [12]

	normalExec[0xc2] = func() { cpu.jp("NZ", d.u16) }                                 // JP NZ a16 []
	d.normal[0xc2] = []func(){readFirstParam, readSecondParam, normalExec[0xc2], nop} // JP NZ a16    3 [16 12]

	normalExec[0xc3] = func() { cpu.jp("", d.u16) }                                   // JP a16  []
	d.normal[0xc3] = []func(){readFirstParam, readSecondParam, nop, normalExec[0xc3]} // JP a16       3 [16]

	normalExec[0xc4] = func() { cpu.call("NZ", d.u16, mem) }                                    // CALL NZ a16 []
	d.normal[0xc4] = []func(){readFirstParam, readSecondParam, normalExec[0xc4], nop, nop, nop} // CALL NZ a16  3 [24 12]

	// PUSH BC      1 [16]
	d.normal[0xc5] = []func(){
		nop,
		nop,
		func() {
			cpu.sp--
			mem.Write(cpu.sp, cpu.b)
		},
		func() {
			cpu.sp--
			mem.Write(cpu.sp, cpu.c)
		},
	}

	normalExec[0xc6] = func() { cpu.add(d.u8) }                 // ADD A d8 [Z 0 H C]
	d.normal[0xc6] = []func(){readFirstParam, normalExec[0xc6]} // ADD A d8     2 [8]

	normalExec[0xc7] = func() { cpu.rst(0x0000, mem) }         // RST 00H  []
	d.normal[0xc7] = []func(){nop, nop, nop, normalExec[0xc7]} // RST 00H      1 [16]

	normalExec[0xc8] = func() { cpu.ret("Z", mem) }                 // RET Z  []
	d.normal[0xc8] = []func(){nop, normalExec[0xc8], nop, nop, nop} // RET Z        1 [20 8]

	normalExec[0xc9] = func() { cpu.ret("", mem) }             // RET   []
	d.normal[0xc9] = []func(){nop, nop, nop, normalExec[0xc9]} // RET          1 [16]

	normalExec[0xca] = func() { cpu.jp("Z", d.u16) }                                  // JP Z a16 []
	d.normal[0xca] = []func(){readFirstParam, readSecondParam, normalExec[0xca], nop} // JP Z a16     3 [16 12]

	normalExec[0xcc] = func() { cpu.call("Z", d.u16, mem) }                                     // CALL Z a16 []
	d.normal[0xcc] = []func(){readFirstParam, readSecondParam, normalExec[0xcc], nop, nop, nop} // CALL Z a16   3 [24 12]

	normalExec[0xcd] = func() { cpu.call("", d.u16, mem) }                                      // CALL a16  []
	d.normal[0xcd] = []func(){readFirstParam, readSecondParam, nop, nop, nop, normalExec[0xcd]} // CALL a16     3 [24]

	normalExec[0xce] = func() { cpu.adc(d.u8) }                 // ADC A d8 [Z 0 H C]
	d.normal[0xce] = []func(){readFirstParam, normalExec[0xce]} // ADC A d8     2 [8]

	normalExec[0xcf] = func() { cpu.rst(0x0008, mem) }         // RST 08H  []
	d.normal[0xcf] = []func(){nop, nop, nop, normalExec[0xcf]} // RST 08H      1 [16]

	normalExec[0xd0] = func() { cpu.ret("NC", mem) }                // RET NC  []
	d.normal[0xd0] = []func(){nop, normalExec[0xd0], nop, nop, nop} // RET NC       1 [20 8]

	normalExec[0xd1] = func() { cpu.pop(&cpu.d, &cpu.e, mem) } // POP DE  []
	d.normal[0xd1] = []func(){nop, nop, normalExec[0xd1]}      // POP DE       1 [12]

	normalExec[0xd2] = func() { cpu.jp("NC", d.u16) }                                 // JP NC a16 []
	d.normal[0xd2] = []func(){readFirstParam, readSecondParam, normalExec[0xd2], nop} // JP NC a16    3 [16 12]

	normalExec[0xd4] = func() { cpu.call("NC", d.u16, mem) }                                    // CALL NC a16 []
	d.normal[0xd4] = []func(){readFirstParam, readSecondParam, normalExec[0xd4], nop, nop, nop} // CALL NC a16  3 [24 12]

	// PUSH DE      1 [16]
	d.normal[0xd5] = []func(){
		nop,
		nop,
		func() {
			cpu.sp--
			mem.Write(cpu.sp, cpu.d)
		},
		func() {
			cpu.sp--
			mem.Write(cpu.sp, cpu.e)
		},
	}

	normalExec[0xd6] = func() { cpu.sub(d.u8) }                 // SUB d8  [Z 1 H C]
	d.normal[0xd6] = []func(){readFirstParam, normalExec[0xd6]} // SUB d8       2 [8]

	normalExec[0xd7] = func() { cpu.rst(0x0010, mem) }         // RST 10H  []
	d.normal[0xd7] = []func(){nop, nop, nop, normalExec[0xd7]} // RST 10H      1 [16]

	normalExec[0xd8] = func() { cpu.ret("C", mem) }                 // RET C  []
	d.normal[0xd8] = []func(){nop, normalExec[0xd8], nop, nop, nop} // RET C        1 [20 8]

	normalExec[0xd9] = func() { cpu.reti(mem) }                // RETI   []
	d.normal[0xd9] = []func(){nop, nop, nop, normalExec[0xd9]} // RETI         1 [16]

	normalExec[0xda] = func() { cpu.jp("C", d.u16) }                                  // JP C a16 []
	d.normal[0xda] = []func(){readFirstParam, readSecondParam, normalExec[0xda], nop} // JP C a16     3 [16 12]

	normalExec[0xdc] = func() { cpu.call("C", d.u16, mem) }                                     // CALL C a16 []
	d.normal[0xdc] = []func(){readFirstParam, readSecondParam, normalExec[0xdc], nop, nop, nop} // CALL C a16   3 [24 12]

	normalExec[0xde] = func() { cpu.sbc(d.u8) }                 // SBC A d8 [Z 1 H C]
	d.normal[0xde] = []func(){readFirstParam, normalExec[0xde]} // SBC A d8     2 [8]

	normalExec[0xdf] = func() { cpu.rst(0x0018, mem) }         // RST 18H  []
	d.normal[0xdf] = []func(){nop, nop, nop, normalExec[0xdf]} // RST 18H      1 [16]

	// LDH (a8) A   2 [12]
	d.normal[0xe0] = []func(){
		readFirstParam,
		func() { cpu.ld(&d.mem8, cpu.a) },
		func() { mem.Write(uint16(0xff00+uint16(d.u8)), d.mem8) },
	}

	normalExec[0xe1] = func() { cpu.pop(&cpu.h, &cpu.l, mem) } // POP HL  []
	d.normal[0xe1] = []func(){nop, nop, normalExec[0xe1]}      // POP HL       1 [12]

	// LD (C) A     1 [8]
	d.normal[0xe2] = []func(){
		func() { cpu.ld(&d.mem8, cpu.a) },
		func() { mem.Write(uint16(0xff00+uint16(cpu.c)), d.mem8) },
	}

	// PUSH HL      1 [16]
	d.normal[0xe5] = []func(){
		nop,
		nop,
		func() {
			cpu.sp--
			mem.Write(cpu.sp, cpu.h)
		},
		func() {
			cpu.sp--
			mem.Write(cpu.sp, cpu.l)
		},
	}

	normalExec[0xe6] = func() { cpu.and(d.u8) }                 // AND d8  [Z 0 1 0]
	d.normal[0xe6] = []func(){readFirstParam, normalExec[0xe6]} // AND d8       2 [8]

	normalExec[0xe7] = func() { cpu.rst(0x0020, mem) }         // RST 20H  []
	d.normal[0xe7] = []func(){nop, nop, nop, normalExec[0xe7]} // RST 20H      1 [16]

	normalExec[0xe8] = func() { cpu.addSP(int8(d.u8)) }                   // ADD SP r8 [0 0 H C]
	d.normal[0xe8] = []func(){readFirstParam, nop, nop, normalExec[0xe8]} // ADD SP r8    2 [16]

	normalExec[0xe9] = func() { cpu.jp("", cpu.hl()) } // JP (HL)  []
	d.normal[0xe9] = []func(){normalExec[0xe9]}        // JP (HL)      1 [4]

	normalExec[0xea] = func() { cpu.ldA16U8(d.u16, cpu.a, mem) }                              // LD (a16) A []
	d.normal[0xea] = []func(){readFirstParam, readSecondParam, normalExec[0xea], writeMemory} // LD (a16) A   3 [16]

	normalExec[0xee] = func() { cpu.xor(d.u8) }                 // XOR d8  [Z 0 0 0]
	d.normal[0xee] = []func(){readFirstParam, normalExec[0xee]} // XOR d8       2 [8]

	normalExec[0xef] = func() { cpu.rst(0x0028, mem) }         // RST 28H  []
	d.normal[0xef] = []func(){nop, nop, nop, normalExec[0xef]} // RST 28H      1 [16]

	// LDH A (a8)   2 [12]
	d.normal[0xf0] = []func(){
		readFirstParam,
		nop,
		func() {
			addr := uint16(0xff00 + uint16(d.u8))
			cpu.ld(&cpu.a, mem.Read(addr))
		},
	}

	normalExec[0xf1] = func() { cpu.popAF(mem) }          // POP AF  [Z N H C]
	d.normal[0xf1] = []func(){nop, nop, normalExec[0xf1]} // POP AF       1 [12]

	// LD A (C)     1 [8]
	d.normal[0xf2] = []func(){
		nop,
		func() {
			addr := uint16(0xff00 + uint16(cpu.c))
			cpu.ld(&cpu.a, mem.Read(addr))
		},
	}

	normalExec[0xf3] = func() { cpu.di() }      // DI   []
	d.normal[0xf3] = []func(){normalExec[0xf3]} // DI           1 [4]

	// PUSH AF      1 [16]
	d.normal[0xf5] = []func(){
		nop,
		nop,
		func() {
			cpu.sp--
			mem.Write(cpu.sp, cpu.a)
		},
		func() {
			cpu.sp--
			mem.Write(cpu.sp, cpu.f)
		},
	}

	normalExec[0xf6] = func() { cpu.or(d.u8) }                  // OR d8  [Z 0 0 0]
	d.normal[0xf6] = []func(){readFirstParam, normalExec[0xf6]} // OR d8        2 [8]

	normalExec[0xf7] = func() { cpu.rst(0x0030, mem) }         // RST 30H  []
	d.normal[0xf7] = []func(){nop, nop, nop, normalExec[0xf7]} // RST 30H      1 [16]

	normalExec[0xf8] = func() { cpu.ldHLSP(int8(d.u8)) }             // LD HL SP+r8 [0 0 H C]
	d.normal[0xf8] = []func(){readFirstParam, nop, normalExec[0xf8]} // LD HL SP+r8  2 [12]

	normalExec[0xf9] = func() { cpu.ldSPHL() }       // LD SP HL []
	d.normal[0xf9] = []func(){nop, normalExec[0xf9]} // LD SP HL     1 [8]

	normalExec[0xfa] = func() { cpu.ldR8A16(&cpu.a, d.u16, mem) }                            // LD A (a16) []
	d.normal[0xfa] = []func(){readFirstParam, readSecondParam, readMemory, normalExec[0xfa]} // LD A (a16)   3 [16]

	normalExec[0xfb] = func() { cpu.ei() }      // EI   []
	d.normal[0xfb] = []func(){normalExec[0xfb]} // EI           1 [4]

	normalExec[0xfe] = func() { cpu.cp(d.u8) }                  // CP d8  [Z 1 H C]
	d.normal[0xfe] = []func(){readFirstParam, normalExec[0xfe]} // CP d8        2 [8]

	normalExec[0xff] = func() { cpu.rst(0x0038, mem) }         // RST 38H  []
	d.normal[0xff] = []func(){nop, nop, nop, normalExec[0xff]} // RST 38H      1 [16]

	prefixExec[0x00] = func() { cpu.rlc(&cpu.b) }    // RLC B  [Z 0 0 C]
	d.prefix[0x00] = []func(){nop, prefixExec[0x00]} // RLC B        2 [8]

	prefixExec[0x01] = func() { cpu.rlc(&cpu.c) }    // RLC C  [Z 0 0 C]
	d.prefix[0x01] = []func(){nop, prefixExec[0x01]} // RLC C        2 [8]

	prefixExec[0x02] = func() { cpu.rlc(&cpu.d) }    // RLC D  [Z 0 0 C]
	d.prefix[0x02] = []func(){nop, prefixExec[0x02]} // RLC D        2 [8]

	prefixExec[0x03] = func() { cpu.rlc(&cpu.e) }    // RLC E  [Z 0 0 C]
	d.prefix[0x03] = []func(){nop, prefixExec[0x03]} // RLC E        2 [8]

	prefixExec[0x04] = func() { cpu.rlc(&cpu.h) }    // RLC H  [Z 0 0 C]
	d.prefix[0x04] = []func(){nop, prefixExec[0x04]} // RLC H        2 [8]

	prefixExec[0x05] = func() { cpu.rlc(&cpu.l) }    // RLC L  [Z 0 0 C]
	d.prefix[0x05] = []func(){nop, prefixExec[0x05]} // RLC L        2 [8]

	prefixExec[0x06] = func() { cpu.rlcAddr(cpu.hl(), mem) }                  // RLC (HL)  [Z 0 0 C]
	d.prefix[0x06] = []func(){nop, readMemory, prefixExec[0x06], writeMemory} // RLC (HL)     2 [16]

	prefixExec[0x07] = func() { cpu.rlc(&cpu.a) }    // RLC A  [Z 0 0 C]
	d.prefix[0x07] = []func(){nop, prefixExec[0x07]} // RLC A        2 [8]

	prefixExec[0x08] = func() { cpu.rrc(&cpu.b) }    // RRC B  [Z 0 0 C]
	d.prefix[0x08] = []func(){nop, prefixExec[0x08]} // RRC B        2 [8]

	prefixExec[0x09] = func() { cpu.rrc(&cpu.c) }    // RRC C  [Z 0 0 C]
	d.prefix[0x09] = []func(){nop, prefixExec[0x09]} // RRC C        2 [8]

	prefixExec[0x0a] = func() { cpu.rrc(&cpu.d) }    // RRC D  [Z 0 0 C]
	d.prefix[0x0a] = []func(){nop, prefixExec[0x0a]} // RRC D        2 [8]

	prefixExec[0x0b] = func() { cpu.rrc(&cpu.e) }    // RRC E  [Z 0 0 C]
	d.prefix[0x0b] = []func(){nop, prefixExec[0x0b]} // RRC E        2 [8]

	prefixExec[0x0c] = func() { cpu.rrc(&cpu.h) }    // RRC H  [Z 0 0 C]
	d.prefix[0x0c] = []func(){nop, prefixExec[0x0c]} // RRC H        2 [8]

	prefixExec[0x0d] = func() { cpu.rrc(&cpu.l) }    // RRC L  [Z 0 0 C]
	d.prefix[0x0d] = []func(){nop, prefixExec[0x0d]} // RRC L        2 [8]

	prefixExec[0x0e] = func() { cpu.rrcAddr(cpu.hl(), mem) }                  // RRC (HL)  [Z 0 0 C]
	d.prefix[0x0e] = []func(){nop, readMemory, prefixExec[0x0e], writeMemory} // RRC (HL)     2 [16]

	prefixExec[0x0f] = func() { cpu.rrc(&cpu.a) }    // RRC A  [Z 0 0 C]
	d.prefix[0x0f] = []func(){nop, prefixExec[0x0f]} // RRC A        2 [8]

	prefixExec[0x10] = func() { cpu.rl(&cpu.b) }     // RL B  [Z 0 0 C]
	d.prefix[0x10] = []func(){nop, prefixExec[0x10]} // RL B         2 [8]

	prefixExec[0x11] = func() { cpu.rl(&cpu.c) }     // RL C  [Z 0 0 C]
	d.prefix[0x11] = []func(){nop, prefixExec[0x11]} // RL C         2 [8]

	prefixExec[0x12] = func() { cpu.rl(&cpu.d) }     // RL D  [Z 0 0 C]
	d.prefix[0x12] = []func(){nop, prefixExec[0x12]} // RL D         2 [8]

	prefixExec[0x13] = func() { cpu.rl(&cpu.e) }     // RL E  [Z 0 0 C]
	d.prefix[0x13] = []func(){nop, prefixExec[0x13]} // RL E         2 [8]

	prefixExec[0x14] = func() { cpu.rl(&cpu.h) }     // RL H  [Z 0 0 C]
	d.prefix[0x14] = []func(){nop, prefixExec[0x14]} // RL H         2 [8]

	prefixExec[0x15] = func() { cpu.rl(&cpu.l) }     // RL L  [Z 0 0 C]
	d.prefix[0x15] = []func(){nop, prefixExec[0x15]} // RL L         2 [8]

	prefixExec[0x16] = func() { cpu.rlAddr(cpu.hl(), mem) }                   // RL (HL)  [Z 0 0 C]
	d.prefix[0x16] = []func(){nop, readMemory, prefixExec[0x16], writeMemory} // RL (HL)      2 [16]

	prefixExec[0x17] = func() { cpu.rl(&cpu.a) }     // RL A  [Z 0 0 C]
	d.prefix[0x17] = []func(){nop, prefixExec[0x17]} // RL A         2 [8]

	prefixExec[0x18] = func() { cpu.rr(&cpu.b) }     // RR B  [Z 0 0 C]
	d.prefix[0x18] = []func(){nop, prefixExec[0x18]} // RR B         2 [8]

	prefixExec[0x19] = func() { cpu.rr(&cpu.c) }     // RR C  [Z 0 0 C]
	d.prefix[0x19] = []func(){nop, prefixExec[0x19]} // RR C         2 [8]

	prefixExec[0x1a] = func() { cpu.rr(&cpu.d) }     // RR D  [Z 0 0 C]
	d.prefix[0x1a] = []func(){nop, prefixExec[0x1a]} // RR D         2 [8]

	prefixExec[0x1b] = func() { cpu.rr(&cpu.e) }     // RR E  [Z 0 0 C]
	d.prefix[0x1b] = []func(){nop, prefixExec[0x1b]} // RR E         2 [8]

	prefixExec[0x1c] = func() { cpu.rr(&cpu.h) }     // RR H  [Z 0 0 C]
	d.prefix[0x1c] = []func(){nop, prefixExec[0x1c]} // RR H         2 [8]

	prefixExec[0x1d] = func() { cpu.rr(&cpu.l) }     // RR L  [Z 0 0 C]
	d.prefix[0x1d] = []func(){nop, prefixExec[0x1d]} // RR L         2 [8]

	prefixExec[0x1e] = func() { cpu.rrAddr(cpu.hl(), mem) }                   // RR (HL)  [Z 0 0 C]
	d.prefix[0x1e] = []func(){nop, readMemory, prefixExec[0x1e], writeMemory} // RR (HL)      2 [16]

	prefixExec[0x1f] = func() { cpu.rr(&cpu.a) }     // RR A  [Z 0 0 C]
	d.prefix[0x1f] = []func(){nop, prefixExec[0x1f]} // RR A         2 [8]

	prefixExec[0x20] = func() { cpu.sla(&cpu.b) }    // SLA B  [Z 0 0 C]
	d.prefix[0x20] = []func(){nop, prefixExec[0x20]} // SLA B        2 [8]

	prefixExec[0x21] = func() { cpu.sla(&cpu.c) }    // SLA C  [Z 0 0 C]
	d.prefix[0x21] = []func(){nop, prefixExec[0x21]} // SLA C        2 [8]

	prefixExec[0x22] = func() { cpu.sla(&cpu.d) }    // SLA D  [Z 0 0 C]
	d.prefix[0x22] = []func(){nop, prefixExec[0x22]} // SLA D        2 [8]

	prefixExec[0x23] = func() { cpu.sla(&cpu.e) }    // SLA E  [Z 0 0 C]
	d.prefix[0x23] = []func(){nop, prefixExec[0x23]} // SLA E        2 [8]

	prefixExec[0x24] = func() { cpu.sla(&cpu.h) }    // SLA H  [Z 0 0 C]
	d.prefix[0x24] = []func(){nop, prefixExec[0x24]} // SLA H        2 [8]

	prefixExec[0x25] = func() { cpu.sla(&cpu.l) }    // SLA L  [Z 0 0 C]
	d.prefix[0x25] = []func(){nop, prefixExec[0x25]} // SLA L        2 [8]

	prefixExec[0x26] = func() { cpu.slaAddr(cpu.hl(), mem) }                  // SLA (HL)  [Z 0 0 C]
	d.prefix[0x26] = []func(){nop, readMemory, prefixExec[0x26], writeMemory} // SLA (HL)     2 [16]

	prefixExec[0x27] = func() { cpu.sla(&cpu.a) }    // SLA A  [Z 0 0 C]
	d.prefix[0x27] = []func(){nop, prefixExec[0x27]} // SLA A        2 [8]

	prefixExec[0x28] = func() { cpu.sra(&cpu.b) }    // SRA B  [Z 0 0 C]
	d.prefix[0x28] = []func(){nop, prefixExec[0x28]} // SRA B        2 [8]

	prefixExec[0x29] = func() { cpu.sra(&cpu.c) }    // SRA C  [Z 0 0 C]
	d.prefix[0x29] = []func(){nop, prefixExec[0x29]} // SRA C        2 [8]

	prefixExec[0x2a] = func() { cpu.sra(&cpu.d) }    // SRA D  [Z 0 0 C]
	d.prefix[0x2a] = []func(){nop, prefixExec[0x2a]} // SRA D        2 [8]

	prefixExec[0x2b] = func() { cpu.sra(&cpu.e) }    // SRA E  [Z 0 0 C]
	d.prefix[0x2b] = []func(){nop, prefixExec[0x2b]} // SRA E        2 [8]

	prefixExec[0x2c] = func() { cpu.sra(&cpu.h) }    // SRA H  [Z 0 0 C]
	d.prefix[0x2c] = []func(){nop, prefixExec[0x2c]} // SRA H        2 [8]

	prefixExec[0x2d] = func() { cpu.sra(&cpu.l) }    // SRA L  [Z 0 0 C]
	d.prefix[0x2d] = []func(){nop, prefixExec[0x2d]} // SRA L        2 [8]

	prefixExec[0x2e] = func() { cpu.sraAddr(cpu.hl(), mem) }                  // SRA (HL)  [Z 0 0 C]
	d.prefix[0x2e] = []func(){nop, readMemory, prefixExec[0x2e], writeMemory} // SRA (HL)     2 [16]

	prefixExec[0x2f] = func() { cpu.sra(&cpu.a) }    // SRA A  [Z 0 0 C]
	d.prefix[0x2f] = []func(){nop, prefixExec[0x2f]} // SRA A        2 [8]

	prefixExec[0x30] = func() { cpu.swap(&cpu.b) }   // SWAP B  [Z 0 0 0]
	d.prefix[0x30] = []func(){nop, prefixExec[0x30]} // SWAP B       2 [8]

	prefixExec[0x31] = func() { cpu.swap(&cpu.c) }   // SWAP C  [Z 0 0 0]
	d.prefix[0x31] = []func(){nop, prefixExec[0x31]} // SWAP C       2 [8]

	prefixExec[0x32] = func() { cpu.swap(&cpu.d) }   // SWAP D  [Z 0 0 0]
	d.prefix[0x32] = []func(){nop, prefixExec[0x32]} // SWAP D       2 [8]

	prefixExec[0x33] = func() { cpu.swap(&cpu.e) }   // SWAP E  [Z 0 0 0]
	d.prefix[0x33] = []func(){nop, prefixExec[0x33]} // SWAP E       2 [8]

	prefixExec[0x34] = func() { cpu.swap(&cpu.h) }   // SWAP H  [Z 0 0 0]
	d.prefix[0x34] = []func(){nop, prefixExec[0x34]} // SWAP H       2 [8]

	prefixExec[0x35] = func() { cpu.swap(&cpu.l) }   // SWAP L  [Z 0 0 0]
	d.prefix[0x35] = []func(){nop, prefixExec[0x35]} // SWAP L       2 [8]

	prefixExec[0x36] = func() { cpu.swapAddr(cpu.hl(), mem) }                 // SWAP (HL)  [Z 0 0 0]
	d.prefix[0x36] = []func(){nop, readMemory, prefixExec[0x36], writeMemory} // SWAP (HL)    2 [16]

	prefixExec[0x37] = func() { cpu.swap(&cpu.a) }   // SWAP A  [Z 0 0 0]
	d.prefix[0x37] = []func(){nop, prefixExec[0x37]} // SWAP A       2 [8]

	prefixExec[0x38] = func() { cpu.srl(&cpu.b) }    // SRL B  [Z 0 0 C]
	d.prefix[0x38] = []func(){nop, prefixExec[0x38]} // SRL B        2 [8]

	prefixExec[0x39] = func() { cpu.srl(&cpu.c) }    // SRL C  [Z 0 0 C]
	d.prefix[0x39] = []func(){nop, prefixExec[0x39]} // SRL C        2 [8]

	prefixExec[0x3a] = func() { cpu.srl(&cpu.d) }    // SRL D  [Z 0 0 C]
	d.prefix[0x3a] = []func(){nop, prefixExec[0x3a]} // SRL D        2 [8]

	prefixExec[0x3b] = func() { cpu.srl(&cpu.e) }    // SRL E  [Z 0 0 C]
	d.prefix[0x3b] = []func(){nop, prefixExec[0x3b]} // SRL E        2 [8]

	prefixExec[0x3c] = func() { cpu.srl(&cpu.h) }    // SRL H  [Z 0 0 C]
	d.prefix[0x3c] = []func(){nop, prefixExec[0x3c]} // SRL H        2 [8]

	prefixExec[0x3d] = func() { cpu.srl(&cpu.l) }    // SRL L  [Z 0 0 C]
	d.prefix[0x3d] = []func(){nop, prefixExec[0x3d]} // SRL L        2 [8]

	prefixExec[0x3e] = func() { cpu.srlAddr(cpu.hl(), mem) }                  // SRL (HL)  [Z 0 0 C]
	d.prefix[0x3e] = []func(){nop, readMemory, prefixExec[0x3e], writeMemory} // SRL (HL)     2 [16]

	prefixExec[0x3f] = func() { cpu.srl(&cpu.a) }    // SRL A  [Z 0 0 C]
	d.prefix[0x3f] = []func(){nop, prefixExec[0x3f]} // SRL A        2 [8]

	prefixExec[0x40] = func() { cpu.bit(0, cpu.b) }  // BIT 0 B [Z 0 1 -]
	d.prefix[0x40] = []func(){nop, prefixExec[0x40]} // BIT 0 B      2 [8]

	prefixExec[0x41] = func() { cpu.bit(0, cpu.c) }  // BIT 0 C [Z 0 1 -]
	d.prefix[0x41] = []func(){nop, prefixExec[0x41]} // BIT 0 C      2 [8]

	prefixExec[0x42] = func() { cpu.bit(0, cpu.d) }  // BIT 0 D [Z 0 1 -]
	d.prefix[0x42] = []func(){nop, prefixExec[0x42]} // BIT 0 D      2 [8]

	prefixExec[0x43] = func() { cpu.bit(0, cpu.e) }  // BIT 0 E [Z 0 1 -]
	d.prefix[0x43] = []func(){nop, prefixExec[0x43]} // BIT 0 E      2 [8]

	prefixExec[0x44] = func() { cpu.bit(0, cpu.h) }  // BIT 0 H [Z 0 1 -]
	d.prefix[0x44] = []func(){nop, prefixExec[0x44]} // BIT 0 H      2 [8]

	prefixExec[0x45] = func() { cpu.bit(0, cpu.l) }  // BIT 0 L [Z 0 1 -]
	d.prefix[0x45] = []func(){nop, prefixExec[0x45]} // BIT 0 L      2 [8]

	prefixExec[0x46] = func() { cpu.bitAddr(0, cpu.hl(), mem) }  // BIT 0 (HL) [Z 0 1 -]
	d.prefix[0x46] = []func(){nop, readMemory, prefixExec[0x46]} // BIT 0 (HL)   2 [12]

	prefixExec[0x47] = func() { cpu.bit(0, cpu.a) }  // BIT 0 A [Z 0 1 -]
	d.prefix[0x47] = []func(){nop, prefixExec[0x47]} // BIT 0 A      2 [8]

	prefixExec[0x48] = func() { cpu.bit(1, cpu.b) }  // BIT 1 B [Z 0 1 -]
	d.prefix[0x48] = []func(){nop, prefixExec[0x48]} // BIT 1 B      2 [8]

	prefixExec[0x49] = func() { cpu.bit(1, cpu.c) }  // BIT 1 C [Z 0 1 -]
	d.prefix[0x49] = []func(){nop, prefixExec[0x49]} // BIT 1 C      2 [8]

	prefixExec[0x4a] = func() { cpu.bit(1, cpu.d) }  // BIT 1 D [Z 0 1 -]
	d.prefix[0x4a] = []func(){nop, prefixExec[0x4a]} // BIT 1 D      2 [8]

	prefixExec[0x4b] = func() { cpu.bit(1, cpu.e) }  // BIT 1 E [Z 0 1 -]
	d.prefix[0x4b] = []func(){nop, prefixExec[0x4b]} // BIT 1 E      2 [8]

	prefixExec[0x4c] = func() { cpu.bit(1, cpu.h) }  // BIT 1 H [Z 0 1 -]
	d.prefix[0x4c] = []func(){nop, prefixExec[0x4c]} // BIT 1 H      2 [8]

	prefixExec[0x4d] = func() { cpu.bit(1, cpu.l) }  // BIT 1 L [Z 0 1 -]
	d.prefix[0x4d] = []func(){nop, prefixExec[0x4d]} // BIT 1 L      2 [8]

	prefixExec[0x4e] = func() { cpu.bitAddr(1, cpu.hl(), mem) }  // BIT 1 (HL) [Z 0 1 -]
	d.prefix[0x4e] = []func(){nop, readMemory, prefixExec[0x4e]} // BIT 1 (HL)   2 [12]

	prefixExec[0x4f] = func() { cpu.bit(1, cpu.a) }  // BIT 1 A [Z 0 1 -]
	d.prefix[0x4f] = []func(){nop, prefixExec[0x4f]} // BIT 1 A      2 [8]

	prefixExec[0x50] = func() { cpu.bit(2, cpu.b) }  // BIT 2 B [Z 0 1 -]
	d.prefix[0x50] = []func(){nop, prefixExec[0x50]} // BIT 2 B      2 [8]

	prefixExec[0x51] = func() { cpu.bit(2, cpu.c) }  // BIT 2 C [Z 0 1 -]
	d.prefix[0x51] = []func(){nop, prefixExec[0x51]} // BIT 2 C      2 [8]

	prefixExec[0x52] = func() { cpu.bit(2, cpu.d) }  // BIT 2 D [Z 0 1 -]
	d.prefix[0x52] = []func(){nop, prefixExec[0x52]} // BIT 2 D      2 [8]

	prefixExec[0x53] = func() { cpu.bit(2, cpu.e) }  // BIT 2 E [Z 0 1 -]
	d.prefix[0x53] = []func(){nop, prefixExec[0x53]} // BIT 2 E      2 [8]

	prefixExec[0x54] = func() { cpu.bit(2, cpu.h) }  // BIT 2 H [Z 0 1 -]
	d.prefix[0x54] = []func(){nop, prefixExec[0x54]} // BIT 2 H      2 [8]

	prefixExec[0x55] = func() { cpu.bit(2, cpu.l) }  // BIT 2 L [Z 0 1 -]
	d.prefix[0x55] = []func(){nop, prefixExec[0x55]} // BIT 2 L      2 [8]

	prefixExec[0x56] = func() { cpu.bitAddr(2, cpu.hl(), mem) }  // BIT 2 (HL) [Z 0 1 -]
	d.prefix[0x56] = []func(){nop, readMemory, prefixExec[0x56]} // BIT 2 (HL)   2 [12]

	prefixExec[0x57] = func() { cpu.bit(2, cpu.a) }  // BIT 2 A [Z 0 1 -]
	d.prefix[0x57] = []func(){nop, prefixExec[0x57]} // BIT 2 A      2 [8]

	prefixExec[0x58] = func() { cpu.bit(3, cpu.b) }  // BIT 3 B [Z 0 1 -]
	d.prefix[0x58] = []func(){nop, prefixExec[0x58]} // BIT 3 B      2 [8]

	prefixExec[0x59] = func() { cpu.bit(3, cpu.c) }  // BIT 3 C [Z 0 1 -]
	d.prefix[0x59] = []func(){nop, prefixExec[0x59]} // BIT 3 C      2 [8]

	prefixExec[0x5a] = func() { cpu.bit(3, cpu.d) }  // BIT 3 D [Z 0 1 -]
	d.prefix[0x5a] = []func(){nop, prefixExec[0x5a]} // BIT 3 D      2 [8]

	prefixExec[0x5b] = func() { cpu.bit(3, cpu.e) }  // BIT 3 E [Z 0 1 -]
	d.prefix[0x5b] = []func(){nop, prefixExec[0x5b]} // BIT 3 E      2 [8]

	prefixExec[0x5c] = func() { cpu.bit(3, cpu.h) }  // BIT 3 H [Z 0 1 -]
	d.prefix[0x5c] = []func(){nop, prefixExec[0x5c]} // BIT 3 H      2 [8]

	prefixExec[0x5d] = func() { cpu.bit(3, cpu.l) }  // BIT 3 L [Z 0 1 -]
	d.prefix[0x5d] = []func(){nop, prefixExec[0x5d]} // BIT 3 L      2 [8]

	prefixExec[0x5e] = func() { cpu.bitAddr(3, cpu.hl(), mem) }  // BIT 3 (HL) [Z 0 1 -]
	d.prefix[0x5e] = []func(){nop, readMemory, prefixExec[0x5e]} // BIT 3 (HL)   2 [12]

	prefixExec[0x5f] = func() { cpu.bit(3, cpu.a) }  // BIT 3 A [Z 0 1 -]
	d.prefix[0x5f] = []func(){nop, prefixExec[0x5f]} // BIT 3 A      2 [8]

	prefixExec[0x60] = func() { cpu.bit(4, cpu.b) }  // BIT 4 B [Z 0 1 -]
	d.prefix[0x60] = []func(){nop, prefixExec[0x60]} // BIT 4 B      2 [8]

	prefixExec[0x61] = func() { cpu.bit(4, cpu.c) }  // BIT 4 C [Z 0 1 -]
	d.prefix[0x61] = []func(){nop, prefixExec[0x61]} // BIT 4 C      2 [8]

	prefixExec[0x62] = func() { cpu.bit(4, cpu.d) }  // BIT 4 D [Z 0 1 -]
	d.prefix[0x62] = []func(){nop, prefixExec[0x62]} // BIT 4 D      2 [8]

	prefixExec[0x63] = func() { cpu.bit(4, cpu.e) }  // BIT 4 E [Z 0 1 -]
	d.prefix[0x63] = []func(){nop, prefixExec[0x63]} // BIT 4 E      2 [8]

	prefixExec[0x64] = func() { cpu.bit(4, cpu.h) }  // BIT 4 H [Z 0 1 -]
	d.prefix[0x64] = []func(){nop, prefixExec[0x64]} // BIT 4 H      2 [8]

	prefixExec[0x65] = func() { cpu.bit(4, cpu.l) }  // BIT 4 L [Z 0 1 -]
	d.prefix[0x65] = []func(){nop, prefixExec[0x65]} // BIT 4 L      2 [8]

	prefixExec[0x66] = func() { cpu.bitAddr(4, cpu.hl(), mem) }  // BIT 4 (HL) [Z 0 1 -]
	d.prefix[0x66] = []func(){nop, readMemory, prefixExec[0x66]} // BIT 4 (HL)   2 [12]

	prefixExec[0x67] = func() { cpu.bit(4, cpu.a) }  // BIT 4 A [Z 0 1 -]
	d.prefix[0x67] = []func(){nop, prefixExec[0x67]} // BIT 4 A      2 [8]

	prefixExec[0x68] = func() { cpu.bit(5, cpu.b) }  // BIT 5 B [Z 0 1 -]
	d.prefix[0x68] = []func(){nop, prefixExec[0x68]} // BIT 5 B      2 [8]

	prefixExec[0x69] = func() { cpu.bit(5, cpu.c) }  // BIT 5 C [Z 0 1 -]
	d.prefix[0x69] = []func(){nop, prefixExec[0x69]} // BIT 5 C      2 [8]

	prefixExec[0x6a] = func() { cpu.bit(5, cpu.d) }  // BIT 5 D [Z 0 1 -]
	d.prefix[0x6a] = []func(){nop, prefixExec[0x6a]} // BIT 5 D      2 [8]

	prefixExec[0x6b] = func() { cpu.bit(5, cpu.e) }  // BIT 5 E [Z 0 1 -]
	d.prefix[0x6b] = []func(){nop, prefixExec[0x6b]} // BIT 5 E      2 [8]

	prefixExec[0x6c] = func() { cpu.bit(5, cpu.h) }  // BIT 5 H [Z 0 1 -]
	d.prefix[0x6c] = []func(){nop, prefixExec[0x6c]} // BIT 5 H      2 [8]

	prefixExec[0x6d] = func() { cpu.bit(5, cpu.l) }  // BIT 5 L [Z 0 1 -]
	d.prefix[0x6d] = []func(){nop, prefixExec[0x6d]} // BIT 5 L      2 [8]

	prefixExec[0x6e] = func() { cpu.bitAddr(5, cpu.hl(), mem) }  // BIT 5 (HL) [Z 0 1 -]
	d.prefix[0x6e] = []func(){nop, readMemory, prefixExec[0x6e]} // BIT 5 (HL)   2 [12]

	prefixExec[0x6f] = func() { cpu.bit(5, cpu.a) }  // BIT 5 A [Z 0 1 -]
	d.prefix[0x6f] = []func(){nop, prefixExec[0x6f]} // BIT 5 A      2 [8]

	prefixExec[0x70] = func() { cpu.bit(6, cpu.b) }  // BIT 6 B [Z 0 1 -]
	d.prefix[0x70] = []func(){nop, prefixExec[0x70]} // BIT 6 B      2 [8]

	prefixExec[0x71] = func() { cpu.bit(6, cpu.c) }  // BIT 6 C [Z 0 1 -]
	d.prefix[0x71] = []func(){nop, prefixExec[0x71]} // BIT 6 C      2 [8]

	prefixExec[0x72] = func() { cpu.bit(6, cpu.d) }  // BIT 6 D [Z 0 1 -]
	d.prefix[0x72] = []func(){nop, prefixExec[0x72]} // BIT 6 D      2 [8]

	prefixExec[0x73] = func() { cpu.bit(6, cpu.e) }  // BIT 6 E [Z 0 1 -]
	d.prefix[0x73] = []func(){nop, prefixExec[0x73]} // BIT 6 E      2 [8]

	prefixExec[0x74] = func() { cpu.bit(6, cpu.h) }  // BIT 6 H [Z 0 1 -]
	d.prefix[0x74] = []func(){nop, prefixExec[0x74]} // BIT 6 H      2 [8]

	prefixExec[0x75] = func() { cpu.bit(6, cpu.l) }  // BIT 6 L [Z 0 1 -]
	d.prefix[0x75] = []func(){nop, prefixExec[0x75]} // BIT 6 L      2 [8]

	prefixExec[0x76] = func() { cpu.bitAddr(6, cpu.hl(), mem) }  // BIT 6 (HL) [Z 0 1 -]
	d.prefix[0x76] = []func(){nop, readMemory, prefixExec[0x76]} // BIT 6 (HL)   2 [12]

	prefixExec[0x77] = func() { cpu.bit(6, cpu.a) }  // BIT 6 A [Z 0 1 -]
	d.prefix[0x77] = []func(){nop, prefixExec[0x77]} // BIT 6 A      2 [8]

	prefixExec[0x78] = func() { cpu.bit(7, cpu.b) }  // BIT 7 B [Z 0 1 -]
	d.prefix[0x78] = []func(){nop, prefixExec[0x78]} // BIT 7 B      2 [8]

	prefixExec[0x79] = func() { cpu.bit(7, cpu.c) }  // BIT 7 C [Z 0 1 -]
	d.prefix[0x79] = []func(){nop, prefixExec[0x79]} // BIT 7 C      2 [8]

	prefixExec[0x7a] = func() { cpu.bit(7, cpu.d) }  // BIT 7 D [Z 0 1 -]
	d.prefix[0x7a] = []func(){nop, prefixExec[0x7a]} // BIT 7 D      2 [8]

	prefixExec[0x7b] = func() { cpu.bit(7, cpu.e) }  // BIT 7 E [Z 0 1 -]
	d.prefix[0x7b] = []func(){nop, prefixExec[0x7b]} // BIT 7 E      2 [8]

	prefixExec[0x7c] = func() { cpu.bit(7, cpu.h) }  // BIT 7 H [Z 0 1 -]
	d.prefix[0x7c] = []func(){nop, prefixExec[0x7c]} // BIT 7 H      2 [8]

	prefixExec[0x7d] = func() { cpu.bit(7, cpu.l) }  // BIT 7 L [Z 0 1 -]
	d.prefix[0x7d] = []func(){nop, prefixExec[0x7d]} // BIT 7 L      2 [8]

	prefixExec[0x7e] = func() { cpu.bitAddr(7, cpu.hl(), mem) }  // BIT 7 (HL) [Z 0 1 -]
	d.prefix[0x7e] = []func(){nop, readMemory, prefixExec[0x7e]} // BIT 7 (HL)   2 [12]

	prefixExec[0x7f] = func() { cpu.bit(7, cpu.a) }  // BIT 7 A [Z 0 1 -]
	d.prefix[0x7f] = []func(){nop, prefixExec[0x7f]} // BIT 7 A      2 [8]

	prefixExec[0x80] = func() { cpu.res(0, &cpu.b) } // RES 0 B []
	d.prefix[0x80] = []func(){nop, prefixExec[0x80]} // RES 0 B      2 [8]

	prefixExec[0x81] = func() { cpu.res(0, &cpu.c) } // RES 0 C []
	d.prefix[0x81] = []func(){nop, prefixExec[0x81]} // RES 0 C      2 [8]

	prefixExec[0x82] = func() { cpu.res(0, &cpu.d) } // RES 0 D []
	d.prefix[0x82] = []func(){nop, prefixExec[0x82]} // RES 0 D      2 [8]

	prefixExec[0x83] = func() { cpu.res(0, &cpu.e) } // RES 0 E []
	d.prefix[0x83] = []func(){nop, prefixExec[0x83]} // RES 0 E      2 [8]

	prefixExec[0x84] = func() { cpu.res(0, &cpu.h) } // RES 0 H []
	d.prefix[0x84] = []func(){nop, prefixExec[0x84]} // RES 0 H      2 [8]

	prefixExec[0x85] = func() { cpu.res(0, &cpu.l) } // RES 0 L []
	d.prefix[0x85] = []func(){nop, prefixExec[0x85]} // RES 0 L      2 [8]

	prefixExec[0x86] = func() { cpu.resAddr(0, cpu.hl(), mem) }        // RES 0 (HL) []
	d.prefix[0x86] = []func(){nop, nop, prefixExec[0x86], writeMemory} // RES 0 (HL)   2 [16]

	prefixExec[0x87] = func() { cpu.res(0, &cpu.a) } // RES 0 A []
	d.prefix[0x87] = []func(){nop, prefixExec[0x87]} // RES 0 A      2 [8]

	prefixExec[0x88] = func() { cpu.res(1, &cpu.b) } // RES 1 B []
	d.prefix[0x88] = []func(){nop, prefixExec[0x88]} // RES 1 B      2 [8]

	prefixExec[0x89] = func() { cpu.res(1, &cpu.c) } // RES 1 C []
	d.prefix[0x89] = []func(){nop, prefixExec[0x89]} // RES 1 C      2 [8]

	prefixExec[0x8a] = func() { cpu.res(1, &cpu.d) } // RES 1 D []
	d.prefix[0x8a] = []func(){nop, prefixExec[0x8a]} // RES 1 D      2 [8]

	prefixExec[0x8b] = func() { cpu.res(1, &cpu.e) } // RES 1 E []
	d.prefix[0x8b] = []func(){nop, prefixExec[0x8b]} // RES 1 E      2 [8]

	prefixExec[0x8c] = func() { cpu.res(1, &cpu.h) } // RES 1 H []
	d.prefix[0x8c] = []func(){nop, prefixExec[0x8c]} // RES 1 H      2 [8]

	prefixExec[0x8d] = func() { cpu.res(1, &cpu.l) } // RES 1 L []
	d.prefix[0x8d] = []func(){nop, prefixExec[0x8d]} // RES 1 L      2 [8]

	prefixExec[0x8e] = func() { cpu.resAddr(1, cpu.hl(), mem) }        // RES 1 (HL) []
	d.prefix[0x8e] = []func(){nop, nop, prefixExec[0x8e], writeMemory} // RES 1 (HL)   2 [16]

	prefixExec[0x8f] = func() { cpu.res(1, &cpu.a) } // RES 1 A []
	d.prefix[0x8f] = []func(){nop, prefixExec[0x8f]} // RES 1 A      2 [8]

	prefixExec[0x90] = func() { cpu.res(2, &cpu.b) } // RES 2 B []
	d.prefix[0x90] = []func(){nop, prefixExec[0x90]} // RES 2 B      2 [8]

	prefixExec[0x91] = func() { cpu.res(2, &cpu.c) } // RES 2 C []
	d.prefix[0x91] = []func(){nop, prefixExec[0x91]} // RES 2 C      2 [8]

	prefixExec[0x92] = func() { cpu.res(2, &cpu.d) } // RES 2 D []
	d.prefix[0x92] = []func(){nop, prefixExec[0x92]} // RES 2 D      2 [8]

	prefixExec[0x93] = func() { cpu.res(2, &cpu.e) } // RES 2 E []
	d.prefix[0x93] = []func(){nop, prefixExec[0x93]} // RES 2 E      2 [8]

	prefixExec[0x94] = func() { cpu.res(2, &cpu.h) } // RES 2 H []
	d.prefix[0x94] = []func(){nop, prefixExec[0x94]} // RES 2 H      2 [8]

	prefixExec[0x95] = func() { cpu.res(2, &cpu.l) } // RES 2 L []
	d.prefix[0x95] = []func(){nop, prefixExec[0x95]} // RES 2 L      2 [8]

	prefixExec[0x96] = func() { cpu.resAddr(2, cpu.hl(), mem) }        // RES 2 (HL) []
	d.prefix[0x96] = []func(){nop, nop, prefixExec[0x96], writeMemory} // RES 2 (HL)   2 [16]

	prefixExec[0x97] = func() { cpu.res(2, &cpu.a) } // RES 2 A []
	d.prefix[0x97] = []func(){nop, prefixExec[0x97]} // RES 2 A      2 [8]

	prefixExec[0x98] = func() { cpu.res(3, &cpu.b) } // RES 3 B []
	d.prefix[0x98] = []func(){nop, prefixExec[0x98]} // RES 3 B      2 [8]

	prefixExec[0x99] = func() { cpu.res(3, &cpu.c) } // RES 3 C []
	d.prefix[0x99] = []func(){nop, prefixExec[0x99]} // RES 3 C      2 [8]

	prefixExec[0x9a] = func() { cpu.res(3, &cpu.d) } // RES 3 D []
	d.prefix[0x9a] = []func(){nop, prefixExec[0x9a]} // RES 3 D      2 [8]

	prefixExec[0x9b] = func() { cpu.res(3, &cpu.e) } // RES 3 E []
	d.prefix[0x9b] = []func(){nop, prefixExec[0x9b]} // RES 3 E      2 [8]

	prefixExec[0x9c] = func() { cpu.res(3, &cpu.h) } // RES 3 H []
	d.prefix[0x9c] = []func(){nop, prefixExec[0x9c]} // RES 3 H      2 [8]

	prefixExec[0x9d] = func() { cpu.res(3, &cpu.l) } // RES 3 L []
	d.prefix[0x9d] = []func(){nop, prefixExec[0x9d]} // RES 3 L      2 [8]

	prefixExec[0x9e] = func() { cpu.resAddr(3, cpu.hl(), mem) }        // RES 3 (HL) []
	d.prefix[0x9e] = []func(){nop, nop, prefixExec[0x9e], writeMemory} // RES 3 (HL)   2 [16]

	prefixExec[0x9f] = func() { cpu.res(3, &cpu.a) } // RES 3 A []
	d.prefix[0x9f] = []func(){nop, prefixExec[0x9f]} // RES 3 A      2 [8]

	prefixExec[0xa0] = func() { cpu.res(4, &cpu.b) } // RES 4 B []
	d.prefix[0xa0] = []func(){nop, prefixExec[0xa0]} // RES 4 B      2 [8]

	prefixExec[0xa1] = func() { cpu.res(4, &cpu.c) } // RES 4 C []
	d.prefix[0xa1] = []func(){nop, prefixExec[0xa1]} // RES 4 C      2 [8]

	prefixExec[0xa2] = func() { cpu.res(4, &cpu.d) } // RES 4 D []
	d.prefix[0xa2] = []func(){nop, prefixExec[0xa2]} // RES 4 D      2 [8]

	prefixExec[0xa3] = func() { cpu.res(4, &cpu.e) } // RES 4 E []
	d.prefix[0xa3] = []func(){nop, prefixExec[0xa3]} // RES 4 E      2 [8]

	prefixExec[0xa4] = func() { cpu.res(4, &cpu.h) } // RES 4 H []
	d.prefix[0xa4] = []func(){nop, prefixExec[0xa4]} // RES 4 H      2 [8]

	prefixExec[0xa5] = func() { cpu.res(4, &cpu.l) } // RES 4 L []
	d.prefix[0xa5] = []func(){nop, prefixExec[0xa5]} // RES 4 L      2 [8]

	prefixExec[0xa6] = func() { cpu.resAddr(4, cpu.hl(), mem) }        // RES 4 (HL) []
	d.prefix[0xa6] = []func(){nop, nop, prefixExec[0xa6], writeMemory} // RES 4 (HL)   2 [16]

	prefixExec[0xa7] = func() { cpu.res(4, &cpu.a) } // RES 4 A []
	d.prefix[0xa7] = []func(){nop, prefixExec[0xa7]} // RES 4 A      2 [8]

	prefixExec[0xa8] = func() { cpu.res(5, &cpu.b) } // RES 5 B []
	d.prefix[0xa8] = []func(){nop, prefixExec[0xa8]} // RES 5 B      2 [8]

	prefixExec[0xa9] = func() { cpu.res(5, &cpu.c) } // RES 5 C []
	d.prefix[0xa9] = []func(){nop, prefixExec[0xa9]} // RES 5 C      2 [8]

	prefixExec[0xaa] = func() { cpu.res(5, &cpu.d) } // RES 5 D []
	d.prefix[0xaa] = []func(){nop, prefixExec[0xaa]} // RES 5 D      2 [8]

	prefixExec[0xab] = func() { cpu.res(5, &cpu.e) } // RES 5 E []
	d.prefix[0xab] = []func(){nop, prefixExec[0xab]} // RES 5 E      2 [8]

	prefixExec[0xac] = func() { cpu.res(5, &cpu.h) } // RES 5 H []
	d.prefix[0xac] = []func(){nop, prefixExec[0xac]} // RES 5 H      2 [8]

	prefixExec[0xad] = func() { cpu.res(5, &cpu.l) } // RES 5 L []
	d.prefix[0xad] = []func(){nop, prefixExec[0xad]} // RES 5 L      2 [8]

	prefixExec[0xae] = func() { cpu.resAddr(5, cpu.hl(), mem) }        // RES 5 (HL) []
	d.prefix[0xae] = []func(){nop, nop, prefixExec[0xae], writeMemory} // RES 5 (HL)   2 [16]

	prefixExec[0xaf] = func() { cpu.res(5, &cpu.a) } // RES 5 A []
	d.prefix[0xaf] = []func(){nop, prefixExec[0xaf]} // RES 5 A      2 [8]

	prefixExec[0xb0] = func() { cpu.res(6, &cpu.b) } // RES 6 B []
	d.prefix[0xb0] = []func(){nop, prefixExec[0xb0]} // RES 6 B      2 [8]

	prefixExec[0xb1] = func() { cpu.res(6, &cpu.c) } // RES 6 C []
	d.prefix[0xb1] = []func(){nop, prefixExec[0xb1]} // RES 6 C      2 [8]

	prefixExec[0xb2] = func() { cpu.res(6, &cpu.d) } // RES 6 D []
	d.prefix[0xb2] = []func(){nop, prefixExec[0xb2]} // RES 6 D      2 [8]

	prefixExec[0xb3] = func() { cpu.res(6, &cpu.e) } // RES 6 E []
	d.prefix[0xb3] = []func(){nop, prefixExec[0xb3]} // RES 6 E      2 [8]

	prefixExec[0xb4] = func() { cpu.res(6, &cpu.h) } // RES 6 H []
	d.prefix[0xb4] = []func(){nop, prefixExec[0xb4]} // RES 6 H      2 [8]

	prefixExec[0xb5] = func() { cpu.res(6, &cpu.l) } // RES 6 L []
	d.prefix[0xb5] = []func(){nop, prefixExec[0xb5]} // RES 6 L      2 [8]

	prefixExec[0xb6] = func() { cpu.resAddr(6, cpu.hl(), mem) }        // RES 6 (HL) []
	d.prefix[0xb6] = []func(){nop, nop, prefixExec[0xb6], writeMemory} // RES 6 (HL)   2 [16]

	prefixExec[0xb7] = func() { cpu.res(6, &cpu.a) } // RES 6 A []
	d.prefix[0xb7] = []func(){nop, prefixExec[0xb7]} // RES 6 A      2 [8]

	prefixExec[0xb8] = func() { cpu.res(7, &cpu.b) } // RES 7 B []
	d.prefix[0xb8] = []func(){nop, prefixExec[0xb8]} // RES 7 B      2 [8]

	prefixExec[0xb9] = func() { cpu.res(7, &cpu.c) } // RES 7 C []
	d.prefix[0xb9] = []func(){nop, prefixExec[0xb9]} // RES 7 C      2 [8]

	prefixExec[0xba] = func() { cpu.res(7, &cpu.d) } // RES 7 D []
	d.prefix[0xba] = []func(){nop, prefixExec[0xba]} // RES 7 D      2 [8]

	prefixExec[0xbb] = func() { cpu.res(7, &cpu.e) } // RES 7 E []
	d.prefix[0xbb] = []func(){nop, prefixExec[0xbb]} // RES 7 E      2 [8]

	prefixExec[0xbc] = func() { cpu.res(7, &cpu.h) } // RES 7 H []
	d.prefix[0xbc] = []func(){nop, prefixExec[0xbc]} // RES 7 H      2 [8]

	prefixExec[0xbd] = func() { cpu.res(7, &cpu.l) } // RES 7 L []
	d.prefix[0xbd] = []func(){nop, prefixExec[0xbd]} // RES 7 L      2 [8]

	prefixExec[0xbe] = func() { cpu.resAddr(7, cpu.hl(), mem) }        // RES 7 (HL) []
	d.prefix[0xbe] = []func(){nop, nop, prefixExec[0xbe], writeMemory} // RES 7 (HL)   2 [16]

	prefixExec[0xbf] = func() { cpu.res(7, &cpu.a) } // RES 7 A []
	d.prefix[0xbf] = []func(){nop, prefixExec[0xbf]} // RES 7 A      2 [8]

	prefixExec[0xc0] = func() { cpu.set(0, &cpu.b) } // SET 0 B []
	d.prefix[0xc0] = []func(){nop, prefixExec[0xc0]} // SET 0 B      2 [8]

	prefixExec[0xc1] = func() { cpu.set(0, &cpu.c) } // SET 0 C []
	d.prefix[0xc1] = []func(){nop, prefixExec[0xc1]} // SET 0 C      2 [8]

	prefixExec[0xc2] = func() { cpu.set(0, &cpu.d) } // SET 0 D []
	d.prefix[0xc2] = []func(){nop, prefixExec[0xc2]} // SET 0 D      2 [8]

	prefixExec[0xc3] = func() { cpu.set(0, &cpu.e) } // SET 0 E []
	d.prefix[0xc3] = []func(){nop, prefixExec[0xc3]} // SET 0 E      2 [8]

	prefixExec[0xc4] = func() { cpu.set(0, &cpu.h) } // SET 0 H []
	d.prefix[0xc4] = []func(){nop, prefixExec[0xc4]} // SET 0 H      2 [8]

	prefixExec[0xc5] = func() { cpu.set(0, &cpu.l) } // SET 0 L []
	d.prefix[0xc5] = []func(){nop, prefixExec[0xc5]} // SET 0 L      2 [8]

	prefixExec[0xc6] = func() { cpu.setAddr(0, cpu.hl(), mem) }        // SET 0 (HL) []
	d.prefix[0xc6] = []func(){nop, nop, prefixExec[0xc6], writeMemory} // SET 0 (HL)   2 [16]

	prefixExec[0xc7] = func() { cpu.set(0, &cpu.a) } // SET 0 A []
	d.prefix[0xc7] = []func(){nop, prefixExec[0xc7]} // SET 0 A      2 [8]

	prefixExec[0xc8] = func() { cpu.set(1, &cpu.b) } // SET 1 B []
	d.prefix[0xc8] = []func(){nop, prefixExec[0xc8]} // SET 1 B      2 [8]

	prefixExec[0xc9] = func() { cpu.set(1, &cpu.c) } // SET 1 C []
	d.prefix[0xc9] = []func(){nop, prefixExec[0xc9]} // SET 1 C      2 [8]

	prefixExec[0xca] = func() { cpu.set(1, &cpu.d) } // SET 1 D []
	d.prefix[0xca] = []func(){nop, prefixExec[0xca]} // SET 1 D      2 [8]

	prefixExec[0xcb] = func() { cpu.set(1, &cpu.e) } // SET 1 E []
	d.prefix[0xcb] = []func(){nop, prefixExec[0xcb]} // SET 1 E      2 [8]

	prefixExec[0xcc] = func() { cpu.set(1, &cpu.h) } // SET 1 H []
	d.prefix[0xcc] = []func(){nop, prefixExec[0xcc]} // SET 1 H      2 [8]

	prefixExec[0xcd] = func() { cpu.set(1, &cpu.l) } // SET 1 L []
	d.prefix[0xcd] = []func(){nop, prefixExec[0xcd]} // SET 1 L      2 [8]

	prefixExec[0xce] = func() { cpu.setAddr(1, cpu.hl(), mem) }        // SET 1 (HL) []
	d.prefix[0xce] = []func(){nop, nop, prefixExec[0xce], writeMemory} // SET 1 (HL)   2 [16]

	prefixExec[0xcf] = func() { cpu.set(1, &cpu.a) } // SET 1 A []
	d.prefix[0xcf] = []func(){nop, prefixExec[0xcf]} // SET 1 A      2 [8]

	prefixExec[0xd0] = func() { cpu.set(2, &cpu.b) } // SET 2 B []
	d.prefix[0xd0] = []func(){nop, prefixExec[0xd0]} // SET 2 B      2 [8]

	prefixExec[0xd1] = func() { cpu.set(2, &cpu.c) } // SET 2 C []
	d.prefix[0xd1] = []func(){nop, prefixExec[0xd1]} // SET 2 C      2 [8]

	prefixExec[0xd2] = func() { cpu.set(2, &cpu.d) } // SET 2 D []
	d.prefix[0xd2] = []func(){nop, prefixExec[0xd2]} // SET 2 D      2 [8]

	prefixExec[0xd3] = func() { cpu.set(2, &cpu.e) } // SET 2 E []
	d.prefix[0xd3] = []func(){nop, prefixExec[0xd3]} // SET 2 E      2 [8]

	prefixExec[0xd4] = func() { cpu.set(2, &cpu.h) } // SET 2 H []
	d.prefix[0xd4] = []func(){nop, prefixExec[0xd4]} // SET 2 H      2 [8]

	prefixExec[0xd5] = func() { cpu.set(2, &cpu.l) } // SET 2 L []
	d.prefix[0xd5] = []func(){nop, prefixExec[0xd5]} // SET 2 L      2 [8]

	prefixExec[0xd6] = func() { cpu.setAddr(2, cpu.hl(), mem) }        // SET 2 (HL) []
	d.prefix[0xd6] = []func(){nop, nop, prefixExec[0xd6], writeMemory} // SET 2 (HL)   2 [16]

	prefixExec[0xd7] = func() { cpu.set(2, &cpu.a) } // SET 2 A []
	d.prefix[0xd7] = []func(){nop, prefixExec[0xd7]} // SET 2 A      2 [8]

	prefixExec[0xd8] = func() { cpu.set(3, &cpu.b) } // SET 3 B []
	d.prefix[0xd8] = []func(){nop, prefixExec[0xd8]} // SET 3 B      2 [8]

	prefixExec[0xd9] = func() { cpu.set(3, &cpu.c) } // SET 3 C []
	d.prefix[0xd9] = []func(){nop, prefixExec[0xd9]} // SET 3 C      2 [8]

	prefixExec[0xda] = func() { cpu.set(3, &cpu.d) } // SET 3 D []
	d.prefix[0xda] = []func(){nop, prefixExec[0xda]} // SET 3 D      2 [8]

	prefixExec[0xdb] = func() { cpu.set(3, &cpu.e) } // SET 3 E []
	d.prefix[0xdb] = []func(){nop, prefixExec[0xdb]} // SET 3 E      2 [8]

	prefixExec[0xdc] = func() { cpu.set(3, &cpu.h) } // SET 3 H []
	d.prefix[0xdc] = []func(){nop, prefixExec[0xdc]} // SET 3 H      2 [8]

	prefixExec[0xdd] = func() { cpu.set(3, &cpu.l) } // SET 3 L []
	d.prefix[0xdd] = []func(){nop, prefixExec[0xdd]} // SET 3 L      2 [8]

	prefixExec[0xde] = func() { cpu.setAddr(3, cpu.hl(), mem) }        // SET 3 (HL) []
	d.prefix[0xde] = []func(){nop, nop, prefixExec[0xde], writeMemory} // SET 3 (HL)   2 [16]

	prefixExec[0xdf] = func() { cpu.set(3, &cpu.a) } // SET 3 A []
	d.prefix[0xdf] = []func(){nop, prefixExec[0xdf]} // SET 3 A      2 [8]

	prefixExec[0xe0] = func() { cpu.set(4, &cpu.b) } // SET 4 B []
	d.prefix[0xe0] = []func(){nop, prefixExec[0xe0]} // SET 4 B      2 [8]

	prefixExec[0xe1] = func() { cpu.set(4, &cpu.c) } // SET 4 C []
	d.prefix[0xe1] = []func(){nop, prefixExec[0xe1]} // SET 4 C      2 [8]

	prefixExec[0xe2] = func() { cpu.set(4, &cpu.d) } // SET 4 D []
	d.prefix[0xe2] = []func(){nop, prefixExec[0xe2]} // SET 4 D      2 [8]

	prefixExec[0xe3] = func() { cpu.set(4, &cpu.e) } // SET 4 E []
	d.prefix[0xe3] = []func(){nop, prefixExec[0xe3]} // SET 4 E      2 [8]

	prefixExec[0xe4] = func() { cpu.set(4, &cpu.h) } // SET 4 H []
	d.prefix[0xe4] = []func(){nop, prefixExec[0xe4]} // SET 4 H      2 [8]

	prefixExec[0xe5] = func() { cpu.set(4, &cpu.l) } // SET 4 L []
	d.prefix[0xe5] = []func(){nop, prefixExec[0xe5]} // SET 4 L      2 [8]

	prefixExec[0xe6] = func() { cpu.setAddr(4, cpu.hl(), mem) }        // SET 4 (HL) []
	d.prefix[0xe6] = []func(){nop, nop, prefixExec[0xe6], writeMemory} // SET 4 (HL)   2 [16]

	prefixExec[0xe7] = func() { cpu.set(4, &cpu.a) } // SET 4 A []
	d.prefix[0xe7] = []func(){nop, prefixExec[0xe7]} // SET 4 A      2 [8]

	prefixExec[0xe8] = func() { cpu.set(5, &cpu.b) } // SET 5 B []
	d.prefix[0xe8] = []func(){nop, prefixExec[0xe8]} // SET 5 B      2 [8]

	prefixExec[0xe9] = func() { cpu.set(5, &cpu.c) } // SET 5 C []
	d.prefix[0xe9] = []func(){nop, prefixExec[0xe9]} // SET 5 C      2 [8]

	prefixExec[0xea] = func() { cpu.set(5, &cpu.d) } // SET 5 D []
	d.prefix[0xea] = []func(){nop, prefixExec[0xea]} // SET 5 D      2 [8]

	prefixExec[0xeb] = func() { cpu.set(5, &cpu.e) } // SET 5 E []
	d.prefix[0xeb] = []func(){nop, prefixExec[0xeb]} // SET 5 E      2 [8]

	prefixExec[0xec] = func() { cpu.set(5, &cpu.h) } // SET 5 H []
	d.prefix[0xec] = []func(){nop, prefixExec[0xec]} // SET 5 H      2 [8]

	prefixExec[0xed] = func() { cpu.set(5, &cpu.l) } // SET 5 L []
	d.prefix[0xed] = []func(){nop, prefixExec[0xed]} // SET 5 L      2 [8]

	prefixExec[0xee] = func() { cpu.setAddr(5, cpu.hl(), mem) }        // SET 5 (HL) []
	d.prefix[0xee] = []func(){nop, nop, prefixExec[0xee], writeMemory} // SET 5 (HL)   2 [16]

	prefixExec[0xef] = func() { cpu.set(5, &cpu.a) } // SET 5 A []
	d.prefix[0xef] = []func(){nop, prefixExec[0xef]} // SET 5 A      2 [8]

	prefixExec[0xf0] = func() { cpu.set(6, &cpu.b) } // SET 6 B []
	d.prefix[0xf0] = []func(){nop, prefixExec[0xf0]} // SET 6 B      2 [8]

	prefixExec[0xf1] = func() { cpu.set(6, &cpu.c) } // SET 6 C []
	d.prefix[0xf1] = []func(){nop, prefixExec[0xf1]} // SET 6 C      2 [8]

	prefixExec[0xf2] = func() { cpu.set(6, &cpu.d) } // SET 6 D []
	d.prefix[0xf2] = []func(){nop, prefixExec[0xf2]} // SET 6 D      2 [8]

	prefixExec[0xf3] = func() { cpu.set(6, &cpu.e) } // SET 6 E []
	d.prefix[0xf3] = []func(){nop, prefixExec[0xf3]} // SET 6 E      2 [8]

	prefixExec[0xf4] = func() { cpu.set(6, &cpu.h) } // SET 6 H []
	d.prefix[0xf4] = []func(){nop, prefixExec[0xf4]} // SET 6 H      2 [8]

	prefixExec[0xf5] = func() { cpu.set(6, &cpu.l) } // SET 6 L []
	d.prefix[0xf5] = []func(){nop, prefixExec[0xf5]} // SET 6 L      2 [8]

	prefixExec[0xf6] = func() { cpu.setAddr(6, cpu.hl(), mem) }        // SET 6 (HL) []
	d.prefix[0xf6] = []func(){nop, nop, prefixExec[0xf6], writeMemory} // SET 6 (HL)   2 [16]

	prefixExec[0xf7] = func() { cpu.set(6, &cpu.a) } // SET 6 A []
	d.prefix[0xf7] = []func(){nop, prefixExec[0xf7]} // SET 6 A      2 [8]

	prefixExec[0xf8] = func() { cpu.set(7, &cpu.b) } // SET 7 B []
	d.prefix[0xf8] = []func(){nop, prefixExec[0xf8]} // SET 7 B      2 [8]

	prefixExec[0xf9] = func() { cpu.set(7, &cpu.c) } // SET 7 C []
	d.prefix[0xf9] = []func(){nop, prefixExec[0xf9]} // SET 7 C      2 [8]

	prefixExec[0xfa] = func() { cpu.set(7, &cpu.d) } // SET 7 D []
	d.prefix[0xfa] = []func(){nop, prefixExec[0xfa]} // SET 7 D      2 [8]

	prefixExec[0xfb] = func() { cpu.set(7, &cpu.e) } // SET 7 E []
	d.prefix[0xfb] = []func(){nop, prefixExec[0xfb]} // SET 7 E      2 [8]

	prefixExec[0xfc] = func() { cpu.set(7, &cpu.h) } // SET 7 H []
	d.prefix[0xfc] = []func(){nop, prefixExec[0xfc]} // SET 7 H      2 [8]

	prefixExec[0xfd] = func() { cpu.set(7, &cpu.l) } // SET 7 L []
	d.prefix[0xfd] = []func(){nop, prefixExec[0xfd]} // SET 7 L      2 [8]

	prefixExec[0xfe] = func() { cpu.setAddr(7, cpu.hl(), mem) }        // SET 7 (HL) []
	d.prefix[0xfe] = []func(){nop, nop, prefixExec[0xfe], writeMemory} // SET 7 (HL)   2 [16]

	prefixExec[0xff] = func() { cpu.set(7, &cpu.a) } // SET 7 A []
	d.prefix[0xff] = []func(){nop, prefixExec[0xff]} // SET 7 A      2 [8]
}
