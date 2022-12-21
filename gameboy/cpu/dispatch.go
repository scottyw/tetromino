package cpu

var normal [256][]func(cpu *CPU)
var prefix [256][]func(cpu *CPU)
var veryShortInterrupt []func(cpu *CPU)
var shortInterrupt []func(cpu *CPU)
var longInterrupt []func(cpu *CPU)
var earlyCheck [256]func(int) bool

func nop(cpu *CPU) {
	// Do nothing
}

// Read an 8-bit instruction argument
func readParamA(cpu *CPU) {
	cpu.u8a = cpu.mapper.Read(cpu.pc)
	cpu.pc++
}

// Read an additonal 8-bit instruction argument
func readParamB(cpu *CPU) {
	cpu.u8b = cpu.mapper.Read(cpu.pc)
	cpu.pc++
}

func mooneye(cpu *CPU) {
	// Mooneye uses this instruction (0x40) as a magic breakpoint
	// to indicate that a test rom has completed
	cpu.mooneyeDebugBreakpoint = true
}

func isFinishedSpecial(check func() bool, early, last int) func(int) bool {
	return func(currentCycle int) bool {
		return (currentCycle == early && check()) || currentCycle == last
	}
}

func isFinished(last int) func(int) bool {
	return func(currentCycle int) bool {
		return currentCycle == last
	}
}

//
// FIXME Cleanup leftovers from the core execution refactor
//
// This initialize() method on CPU is no longer necessary and is a leftover from the major refactor of core execution that removed dispatch as a separate emulator component
//
// Replace this method with a standard init() function that doesn't depend on a reference to the CPU
//

func (cpu *CPU) Initialize() {

	next()

	veryShortInterrupt = []func(*CPU){handleInterrupt}
	shortInterrupt = []func(*CPU){nop, nop, nop, nop, handleInterrupt}
	longInterrupt = []func(*CPU){nop, nop, nop, nop, nop, handleInterrupt}

	// earlyCycle[0x20] = 2
	// earlyCycle[0x28] = 2
	// earlyCycle[0x30] = 2
	// earlyCycle[0x38] = 2
	// earlyCycle[0xc0] = 2
	// earlyCycle[0xc2] = 3
	// earlyCycle[0xc4] = 3
	// earlyCycle[0xc8] = 2
	// earlyCycle[0xca] = 3
	// earlyCycle[0xcc] = 3
	// earlyCycle[0xd0] = 2
	// earlyCycle[0xd2] = 3
	// earlyCycle[0xd4] = 3
	// earlyCycle[0xd8] = 2
	// earlyCycle[0xda] = 3
	// earlyCycle[0xdc] = 3

	earlyCheck[0x20] = isFinishedSpecial(zf, 2, 3)
	earlyCheck[0x28] = isFinishedSpecial(nzf, 2, 3)
	earlyCheck[0x30] = isFinishedSpecial(cf, 2, 3)
	earlyCheck[0x38] = isFinishedSpecial(ncf, 2, 3)
	earlyCheck[0xc0] = isFinishedSpecial(zf, 2, 5)
	earlyCheck[0xc2] = isFinishedSpecial(zf, 3, 4)
	earlyCheck[0xc4] = isFinishedSpecial(zf, 3, 6)
	earlyCheck[0xc8] = isFinishedSpecial(nzf, 2, 5)
	earlyCheck[0xca] = isFinishedSpecial(nzf, 3, 4)
	earlyCheck[0xcc] = isFinishedSpecial(nzf, 3, 6)
	earlyCheck[0xd0] = isFinishedSpecial(cf, 2, 5)
	earlyCheck[0xd2] = isFinishedSpecial(cf, 3, 4)
	earlyCheck[0xd4] = isFinishedSpecial(cf, 3, 6)
	earlyCheck[0xd8] = isFinishedSpecial(ncf, 2, 5)
	earlyCheck[0xda] = isFinishedSpecial(ncf, 3, 4)
	earlyCheck[0xdc] = isFinishedSpecial(ncf, 3, 6)

	// NOP          1 [4]
	normal[0x00] = []func(*CPU){nop}

	// LD BC d16 [] 3 [12]
	normal[0x01] = []func(*CPU){readParamA, readParamB, ldBCU16}

	// LD (BC) A [] 1 [8]
	normal[0x02] = []func(*CPU){nop, ldBCA}

	// INC BC  [] 1 [8]
	normal[0x03] = []func(*CPU){nop, incBC}

	// INC B  [Z 0 H -] 1 [4]
	normal[0x04] = []func(*CPU){incB}

	// DEC B  [Z 1 H -] 1 [4]
	normal[0x05] = []func(*CPU){decB}

	// LD B d8 [] 2 [8]
	normal[0x06] = []func(*CPU){readParamA, ldBU}

	// RLCA   [0 0 0 C] 1 [4]
	normal[0x07] = []func(*CPU){rlca}

	// LD (a16) SP [] 3 [20]
	normal[0x08] = []func(*CPU){readParamA, readParamB, nop, writeLowSP, writeHighSP}

	// ADD HL BC [- 0 H C] 1 [8]
	normal[0x09] = []func(*CPU){nop, addHLBC}

	// LD A (BC) [] 1 [8]
	normal[0x0a] = []func(*CPU){nop, ldABC}

	// DEC BC  [] 1 [8]
	normal[0x0b] = []func(*CPU){nop, decBC}

	// INC C  [Z 0 H -] 1 [4]
	normal[0x0c] = []func(*CPU){incC}

	// DEC C  [Z 1 H -] 1 [4]
	normal[0x0d] = []func(*CPU){decC}

	// LD C d8 [] 2 [8]
	normal[0x0e] = []func(*CPU){readParamA, ldCU}

	// RRCA   [0 0 0 C] 1 [4]
	normal[0x0f] = []func(*CPU){rrca}

	// STOP 0  [] 1 [4]
	normal[0x10] = []func(*CPU){stop}

	// LD DE d16 [] 3 [12]
	normal[0x11] = []func(*CPU){readParamA, readParamB, ldDEU16}

	// LD (DE) A [] 1 [8]
	normal[0x12] = []func(*CPU){nop, ldDEA}

	// INC DE  [] 1 [8]
	normal[0x13] = []func(*CPU){nop, incDE}

	// INC D  [Z 0 H -] 1 [4]
	normal[0x14] = []func(*CPU){incD}

	// DEC D  [Z 1 H -] 1 [4]
	normal[0x15] = []func(*CPU){decD}

	// LD D d8 [] 2 [8]
	normal[0x16] = []func(*CPU){readParamA, ldDU}

	// RLA   [0 0 0 C] 1 [4]
	normal[0x17] = []func(*CPU){rla}

	// JR r8  [] 2 [12]
	normal[0x18] = []func(*CPU){nop, readParamA, jr}

	// ADD HL DE [- 0 H C] 1 [8]
	normal[0x19] = []func(*CPU){nop, addHLDE}

	// LD A (DE) [] 1 [8]
	normal[0x1a] = []func(*CPU){nop, ldADE}

	// DEC DE  [] 1 [8]
	normal[0x1b] = []func(*CPU){nop, decDE}

	// INC E  [Z 0 H -] 1 [4]
	normal[0x1c] = []func(*CPU){incE}

	// DEC E  [Z 1 H -] 1 [4]
	normal[0x1d] = []func(*CPU){decE}

	// LD E d8 [] 2 [8]
	normal[0x1e] = []func(*CPU){readParamA, ldEU}

	// RRA   [0 0 0 C] 1 [4]
	normal[0x1f] = []func(*CPU){rra}

	// JR NZ r8 [] 2 [12 8]
	normal[0x20] = []func(*CPU){nop, readParamA, jr}

	// LD HL d16 [] 3 [12]
	normal[0x21] = []func(*CPU){readParamA, readParamB, ldHLU16}

	// LD (HL+) A [] 1 [8]
	normal[0x22] = []func(*CPU){nop, ldHLIA}

	// INC HL  [] 1 [8]
	normal[0x23] = []func(*CPU){nop, incHL}

	// INC H  [Z 0 H -] 1 [4]
	normal[0x24] = []func(*CPU){incH}

	// DEC H  [Z 1 H -] 1 [4]
	normal[0x25] = []func(*CPU){decH}

	// LD H d8 [] 2 [8]
	normal[0x26] = []func(*CPU){readParamA, ldHU}

	// DAA   [Z - 0 C] 1 [4]
	normal[0x27] = []func(*CPU){daa}

	// JR Z r8 [] 2 [12 8]
	normal[0x28] = []func(*CPU){nop, readParamA, jr}

	// ADD HL HL [- 0 H C] 1 [8]
	normal[0x29] = []func(*CPU){nop, addHLHL}

	// LD A (HL+) [] 1 [8]
	normal[0x2a] = []func(*CPU){nop, ldAHLI}

	// DEC HL  [] 1 [8]
	normal[0x2b] = []func(*CPU){nop, decHL}

	// INC L  [Z 0 H -] 1 [4]
	normal[0x2c] = []func(*CPU){incL}

	// DEC L  [Z 1 H -] 1 [4]
	normal[0x2d] = []func(*CPU){decL}

	// LD L d8 [] 2 [8]
	normal[0x2e] = []func(*CPU){readParamA, ldLU}

	// CPL   [- 1 1 -] 1 [4]
	normal[0x2f] = []func(*CPU){cpl}

	// JR NC r8 [] 2 [12 8]
	normal[0x30] = []func(*CPU){nop, readParamA, jr}

	// LD SP d16 [] 3 [12]
	normal[0x31] = []func(*CPU){readParamA, readParamB, ldSPU16}

	// LD (HL-) A [] 1 [8]
	normal[0x32] = []func(*CPU){nop, ldHLDA}

	// INC SP  [] 1 [8]
	normal[0x33] = []func(*CPU){nop, incSP}

	// INC (HL)  [Z 0 H -] 1 [12]
	normal[0x34] = []func(*CPU){nop, ldMHL, incM}

	// DEC (HL)  [Z 1 H -] 1 [12]
	normal[0x35] = []func(*CPU){nop, ldMHL, decM}

	// LD (HL) d8 [] 2 [12]
	normal[0x36] = []func(*CPU){readParamA, nop, ldHLU8}

	// SCF   [- 0 0 1] 1 [4]
	normal[0x37] = []func(*CPU){scf}

	// JR C r8 [] 2 [12 8]
	normal[0x38] = []func(*CPU){nop, readParamA, jr}

	// ADD HL SP [- 0 H C] 1 [8]
	normal[0x39] = []func(*CPU){nop, addHLSP}

	// LD A (HL-) [] 1 [8]
	normal[0x3a] = []func(*CPU){nop, ldAHLD}

	// DEC SP  [] 1 [8]
	normal[0x3b] = []func(*CPU){nop, decSP}

	// INC A  [Z 0 H -] 1 [4]
	normal[0x3c] = []func(*CPU){incA}

	// DEC A  [Z 1 H -] 1 [4]
	normal[0x3d] = []func(*CPU){decA}

	// LD A d8 [] 2 [8]
	normal[0x3e] = []func(*CPU){readParamA, ldAU}

	// CCF   [- 0 0 C] 1 [4]
	normal[0x3f] = []func(*CPU){ccf}

	// LD B B [] 1 [4]
	normal[0x40] = []func(*CPU){mooneye}

	// LD B C [] 1 [4]
	normal[0x41] = []func(*CPU){ldBC}

	// LD B D [] 1 [4]
	normal[0x42] = []func(*CPU){ldBD}

	// LD B E [] 1 [4]
	normal[0x43] = []func(*CPU){ldBE}

	// LD B H [] 1 [4]
	normal[0x44] = []func(*CPU){ldBH}

	// LD B L [] 1 [4]
	normal[0x45] = []func(*CPU){ldBL}

	// LD B (HL) [] 1 [8]
	normal[0x46] = []func(*CPU){nop, ldBHL}

	// LD B A [] 1 [4]
	normal[0x47] = []func(*CPU){ldBA}

	// LD C B [] 1 [4]
	normal[0x48] = []func(*CPU){ldCB}

	// LD C C [] 1 [4]
	normal[0x49] = []func(*CPU){nop}

	// LD C D [] 1 [4]
	normal[0x4a] = []func(*CPU){ldCD}

	// LD C E [] 1 [4]
	normal[0x4b] = []func(*CPU){ldCE}

	// LD C H [] 1 [4]
	normal[0x4c] = []func(*CPU){ldCH}

	// LD C L [] 1 [4]
	normal[0x4d] = []func(*CPU){ldCL}

	// LD C (HL) [] 1 [8]
	normal[0x4e] = []func(*CPU){nop, ldCHL}

	// LD C A [] 1 [4]
	normal[0x4f] = []func(*CPU){ldCA}

	// LD D B [] 1 [4]
	normal[0x50] = []func(*CPU){ldDB}

	// LD D C [] 1 [4]
	normal[0x51] = []func(*CPU){ldDC}

	// LD D D [] 1 [4]
	normal[0x52] = []func(*CPU){nop}

	// LD D E [] 1 [4]
	normal[0x53] = []func(*CPU){ldDE}

	// LD D H [] 1 [4]
	normal[0x54] = []func(*CPU){ldDH}

	// LD D L [] 1 [4]
	normal[0x55] = []func(*CPU){ldDL}

	// LD D (HL) [] 1 [8]
	normal[0x56] = []func(*CPU){nop, ldDHL}

	// LD D A [] 1 [4]
	normal[0x57] = []func(*CPU){ldDA}

	// LD E B [] 1 [4]
	normal[0x58] = []func(*CPU){ldEB}

	// LD E C [] 1 [4]
	normal[0x59] = []func(*CPU){ldEC}

	// LD E D [] 1 [4]
	normal[0x5a] = []func(*CPU){ldED}

	// LD E E [] 1 [4]
	normal[0x5b] = []func(*CPU){nop}

	// LD E H [] 1 [4]
	normal[0x5c] = []func(*CPU){ldEH}

	// LD E L [] 1 [4]
	normal[0x5d] = []func(*CPU){ldEL}

	// LD E (HL) [] 1 [8]
	normal[0x5e] = []func(*CPU){nop, ldEHL}

	// LD E A [] 1 [4]
	normal[0x5f] = []func(*CPU){ldEA}

	// LD H B [] 1 [4]
	normal[0x60] = []func(*CPU){ldHB}

	// LD H C [] 1 [4]
	normal[0x61] = []func(*CPU){ldHC}

	// LD H D [] 1 [4]
	normal[0x62] = []func(*CPU){ldHD}

	// LD H E [] 1 [4]
	normal[0x63] = []func(*CPU){ldHE}

	// LD H H [] 1 [4]
	normal[0x64] = []func(*CPU){nop}

	// LD H L [] 1 [4]
	normal[0x65] = []func(*CPU){ldHL}

	// LD H (HL) [] 1 [8]
	normal[0x66] = []func(*CPU){nop, ldHHL}

	// LD H A [] 1 [4]
	normal[0x67] = []func(*CPU){ldHA}

	// LD L B [] 1 [4]
	normal[0x68] = []func(*CPU){ldLB}

	// LD L C [] 1 [4]
	normal[0x69] = []func(*CPU){ldLC}

	// LD L D [] 1 [4]
	normal[0x6a] = []func(*CPU){ldLD}

	// LD L E [] 1 [4]
	normal[0x6b] = []func(*CPU){ldLE}

	// LD L H [] 1 [4]
	normal[0x6c] = []func(*CPU){ldLH}

	// LD L L [] 1 [4]
	normal[0x6d] = []func(*CPU){nop}

	// LD L (HL) [] 1 [8]
	normal[0x6e] = []func(*CPU){nop, ldLHL}

	// LD L A [] 1 [4]
	normal[0x6f] = []func(*CPU){ldLA}

	// LD (HL) B [] 1 [8]
	normal[0x70] = []func(*CPU){nop, ldHLB}

	// LD (HL) C [] 1 [8]
	normal[0x71] = []func(*CPU){nop, ldHLC}

	// LD (HL) D [] 1 [8]
	normal[0x72] = []func(*CPU){nop, ldHLD}

	// LD (HL) E [] 1 [8]
	normal[0x73] = []func(*CPU){nop, ldHLE}

	// LD (HL) H [] 1 [8]
	normal[0x74] = []func(*CPU){nop, ldHLH}

	// LD (HL) L [] 1 [8]
	normal[0x75] = []func(*CPU){nop, ldHLL}

	// HALT   [] 1 [4]
	normal[0x76] = []func(*CPU){halt}

	// LD (HL) A [] 1 [8]
	normal[0x77] = []func(*CPU){nop, ldHLA}

	// LD A B [] 1 [4]
	normal[0x78] = []func(*CPU){ldAB}

	// LD A C [] 1 [4]
	normal[0x79] = []func(*CPU){ldAC}

	// LD A D [] 1 [4]
	normal[0x7a] = []func(*CPU){ldAD}

	// LD A E [] 1 [4]
	normal[0x7b] = []func(*CPU){ldAE}

	// LD A H [] 1 [4]
	normal[0x7c] = []func(*CPU){ldAH}

	// LD A L [] 1 [4]
	normal[0x7d] = []func(*CPU){ldAL}

	// LD A (HL) [] 1 [8]
	normal[0x7e] = []func(*CPU){nop, ldAHL}

	// LD A A [] 1 [4]
	normal[0x7f] = []func(*CPU){nop}

	// ADD A B [Z 0 H C] 1 [4]
	normal[0x80] = []func(*CPU){addB}

	// ADD A C [Z 0 H C] 1 [4]
	normal[0x81] = []func(*CPU){addC}

	// ADD A D [Z 0 H C] 1 [4]
	normal[0x82] = []func(*CPU){addD}

	// ADD A E [Z 0 H C] 1 [4]
	normal[0x83] = []func(*CPU){addE}

	// ADD A H [Z 0 H C] 1 [4]
	normal[0x84] = []func(*CPU){addH}

	// ADD A L [Z 0 H C] 1 [4]
	normal[0x85] = []func(*CPU){addL}

	// ADD A (HL) [Z 0 H C] 1 [8]
	normal[0x86] = []func(*CPU){nop, addM}

	// ADD A A [Z 0 H C] 1 [4]
	normal[0x87] = []func(*CPU){addA}

	// ADC A B [Z 0 H C] 1 [4]
	normal[0x88] = []func(*CPU){adcB}

	// ADC A C [Z 0 H C] 1 [4]
	normal[0x89] = []func(*CPU){adcC}

	// ADC A D [Z 0 H C] 1 [4]
	normal[0x8a] = []func(*CPU){adcD}

	// ADC A E [Z 0 H C] 1 [4]
	normal[0x8b] = []func(*CPU){adcE}

	// ADC A H [Z 0 H C] 1 [4]
	normal[0x8c] = []func(*CPU){adcH}

	// ADC A L [Z 0 H C] 1 [4]
	normal[0x8d] = []func(*CPU){adcL}

	// ADC A (HL) [Z 0 H C] 1 [8]
	normal[0x8e] = []func(*CPU){nop, adcM}

	// ADC A A [Z 0 H C] 1 [4]
	normal[0x8f] = []func(*CPU){adcA}

	// SUB B  [Z 1 H C] 1 [4]
	normal[0x90] = []func(*CPU){subB}

	// SUB C  [Z 1 H C] 1 [4]
	normal[0x91] = []func(*CPU){subC}

	// SUB D  [Z 1 H C] 1 [4]
	normal[0x92] = []func(*CPU){subD}

	// SUB E  [Z 1 H C] 1 [4]
	normal[0x93] = []func(*CPU){subE}

	// SUB H  [Z 1 H C] 1 [4]
	normal[0x94] = []func(*CPU){subH}

	// SUB L  [Z 1 H C] 1 [4]
	normal[0x95] = []func(*CPU){subL}

	// SUB (HL)  [Z 1 H C] 1 [8]
	normal[0x96] = []func(*CPU){nop, subM}

	// SUB A  [Z 1 H C] 1 [4]
	normal[0x97] = []func(*CPU){subA}

	// SBC A B [Z 1 H C] 1 [4]
	normal[0x98] = []func(*CPU){sbcB}

	// SBC A C [Z 1 H C] 1 [4]
	normal[0x99] = []func(*CPU){sbcC}

	// SBC A D [Z 1 H C] 1 [4]
	normal[0x9a] = []func(*CPU){sbcD}

	// SBC A E [Z 1 H C] 1 [4]
	normal[0x9b] = []func(*CPU){sbcE}

	// SBC A H [Z 1 H C] 1 [4]
	normal[0x9c] = []func(*CPU){sbcH}

	// SBC A L [Z 1 H C] 1 [4]
	normal[0x9d] = []func(*CPU){sbcL}

	// SBC A (HL) [Z 1 H C] 1 [8]
	normal[0x9e] = []func(*CPU){nop, sbcM}

	// SBC A A [Z 1 H C] 1 [4]
	normal[0x9f] = []func(*CPU){sbcA}

	// AND B  [Z 0 1 0] 1 [4]
	normal[0xa0] = []func(*CPU){andB}

	// AND C  [Z 0 1 0] 1 [4]
	normal[0xa1] = []func(*CPU){andC}

	// AND D  [Z 0 1 0] 1 [4]
	normal[0xa2] = []func(*CPU){andD}

	// AND E  [Z 0 1 0] 1 [4]
	normal[0xa3] = []func(*CPU){andE}

	// AND H  [Z 0 1 0] 1 [4]
	normal[0xa4] = []func(*CPU){andH}

	// AND L  [Z 0 1 0] 1 [4]
	normal[0xa5] = []func(*CPU){andL}

	// AND (HL)  [Z 0 1 0] 1 [8]
	normal[0xa6] = []func(*CPU){nop, andM}

	// AND A  [Z 0 1 0] 1 [4]
	normal[0xa7] = []func(*CPU){andA}

	// XOR B  [Z 0 0 0] 1 [4]
	normal[0xa8] = []func(*CPU){xorB}

	// XOR C  [Z 0 0 0] 1 [4]
	normal[0xa9] = []func(*CPU){xorC}

	// XOR D  [Z 0 0 0] 1 [4]
	normal[0xaa] = []func(*CPU){xorD}

	// XOR E  [Z 0 0 0] 1 [4]
	normal[0xab] = []func(*CPU){xorE}

	// XOR H  [Z 0 0 0] 1 [4]
	normal[0xac] = []func(*CPU){xorH}

	// XOR L  [Z 0 0 0] 1 [4]
	normal[0xad] = []func(*CPU){xorL}

	// XOR (HL)  [Z 0 0 0] 1 [8]
	normal[0xae] = []func(*CPU){nop, xorM}

	// XOR A  [Z 0 0 0] 1 [4]
	normal[0xaf] = []func(*CPU){xorA}

	// OR B  [Z 0 0 0] 1 [4]
	normal[0xb0] = []func(*CPU){orB}

	// OR C  [Z 0 0 0] 1 [4]
	normal[0xb1] = []func(*CPU){orC}

	// OR D  [Z 0 0 0] 1 [4]
	normal[0xb2] = []func(*CPU){orD}

	// OR E  [Z 0 0 0] 1 [4]
	normal[0xb3] = []func(*CPU){orE}

	// OR H  [Z 0 0 0] 1 [4]
	normal[0xb4] = []func(*CPU){orH}

	// OR L  [Z 0 0 0] 1 [4]
	normal[0xb5] = []func(*CPU){orL}

	// OR (HL)  [Z 0 0 0] 1 [8]
	normal[0xb6] = []func(*CPU){nop, orM}

	// OR A  [Z 0 0 0] 1 [4]
	normal[0xb7] = []func(*CPU){orA}

	// CP B  [Z 1 H C] 1 [4]
	normal[0xb8] = []func(*CPU){cpB}

	// CP C  [Z 1 H C] 1 [4]
	normal[0xb9] = []func(*CPU){cpC}

	// CP D  [Z 1 H C] 1 [4]
	normal[0xba] = []func(*CPU){cpD}

	// CP E  [Z 1 H C] 1 [4]
	normal[0xbb] = []func(*CPU){cpE}

	// CP H  [Z 1 H C] 1 [4]
	normal[0xbc] = []func(*CPU){cpH}

	// CP L  [Z 1 H C] 1 [4]
	normal[0xbd] = []func(*CPU){cpL}

	// CP (HL)  [Z 1 H C] 1 [8]
	normal[0xbe] = []func(*CPU){nop, cpM}

	// CP A  [Z 1 H C] 1 [4]
	normal[0xbf] = []func(*CPU){cpA}

	// RET NZ  [] 1 [20 8]
	normal[0xc0] = []func(*CPU){nop, nop, pop(&m8a), pop(&m8b), ret}

	// POP BC  [] 1 [12]
	normal[0xc1] = []func(*CPU){nop, pop(&c), pop(&b)}

	// JP NZ a16 [] 3 [16 12]
	normal[0xc2] = []func(*CPU){nop, readParamA, readParamB, jp}

	// JP a16  [] 3 [16]
	normal[0xc3] = []func(*CPU){nop, readParamA, readParamB, jp}

	// CALL NZ a16 [] 3 [24 12]
	normal[0xc4] = []func(*CPU){nop, readParamA, readParamB, call, push(&m8b), push(&m8a)}

	// PUSH BC      1 [16]
	normal[0xc5] = []func(*CPU){nop, nop, push(&b), push(&c)}

	// ADD A d8 [Z 0 H C] 2 [8]
	normal[0xc6] = []func(*CPU){readParamA, addU}

	// RST 00H  [] 1 [16]
	normal[0xc7] = []func(*CPU){nop, rst(0x0000), push(&m8b), push(&m8a)}

	// RET Z  [] 1 [20 8]
	normal[0xc8] = []func(*CPU){nop, nop, pop(&m8a), pop(&m8b), ret}

	// RET   [] 1 [16]
	normal[0xc9] = []func(*CPU){nop, pop(&m8a), pop(&m8b), ret}

	// JP Z a16 [] 3 [16 12]
	normal[0xca] = []func(*CPU){nop, readParamA, readParamB, jp}

	// CALL Z a16 [] 3 [24 12]
	normal[0xcc] = []func(*CPU){nop, readParamA, readParamB, call, push(&m8b), push(&m8a)}

	// CALL a16  [] 3 [24]
	normal[0xcd] = []func(*CPU){nop, readParamA, readParamB, call, push(&m8b), push(&m8a)}

	// ADC A d8 [Z 0 H C] 2 [8]
	normal[0xce] = []func(*CPU){readParamA, adcU}

	// RST 08H  [] 1 [16]
	normal[0xcf] = []func(*CPU){nop, rst(0x0008), push(&m8b), push(&m8a)}

	// RET NC  [] 1 [20 8]
	normal[0xd0] = []func(*CPU){nop, nop, pop(&m8a), pop(&m8b), ret}

	// POP DE  [] 1 [12]
	normal[0xd1] = []func(*CPU){nop, pop(&e), pop(&d)}

	// JP NC a16 [] 3 [16 12]
	normal[0xd2] = []func(*CPU){nop, readParamA, readParamB, jp}

	// CALL NC a16 [] 3 [24 12]
	normal[0xd4] = []func(*CPU){nop, readParamA, readParamB, call, push(&m8b), push(&m8a)}

	// PUSH DE      1 [16]
	normal[0xd5] = []func(*CPU){nop, nop, push(&d), push(&e)}

	// SUB d8  [Z 1 H C] 2 [8]
	normal[0xd6] = []func(*CPU){readParamA, subU}

	// RST 10H  [] 1 [16]
	normal[0xd7] = []func(*CPU){nop, rst(0x0010), push(&m8b), push(&m8a)}

	// RET C  [] 1 [20 8]
	normal[0xd8] = []func(*CPU){nop, nop, pop(&m8a), pop(&m8b), ret}

	// RETI   [] 1 [16]
	normal[0xd9] = []func(*CPU){nop, pop(&m8a), pop(&m8b), reti}

	// JP C a16 [] 3 [16 12]
	normal[0xda] = []func(*CPU){nop, readParamA, readParamB, jp}

	// CALL C a16 [] 3 [24 12]
	normal[0xdc] = []func(*CPU){nop, readParamA, readParamB, call, push(&m8b), push(&m8a)}

	// SBC A d8 [Z 1 H C] 2 [8]
	normal[0xde] = []func(*CPU){readParamA, sbcU}

	// RST 18H  [] 1 [16]
	normal[0xdf] = []func(*CPU){nop, rst(0x0018), push(&m8b), push(&m8a)}

	// LDH (a8) A   2 [12]
	normal[0xe0] = []func(*CPU){readParamA, nop, ldUXA}

	// POP HL  [] 1 [12]
	normal[0xe1] = []func(*CPU){nop, pop(&l), pop(&h)}

	// LD (C) A     1 [8]
	normal[0xe2] = []func(*CPU){nop, ldCXA}

	// PUSH HL      1 [16]
	normal[0xe5] = []func(*CPU){nop, nop, push(&h), push(&l)}

	// AND d8  [Z 0 1 0] 2 [8]
	normal[0xe6] = []func(*CPU){readParamA, andU}

	// RST 20H  [] 1 [16]
	normal[0xe7] = []func(*CPU){nop, rst(0x0020), push(&m8b), push(&m8a)}

	// ADD SP r8 [0 0 H C] 2 [16]
	normal[0xe8] = []func(*CPU){nop, readParamA, addSP, nop}

	// JP (HL)  [] 1 [4]
	normal[0xe9] = []func(*CPU){jpHL}

	// LD (a16) A [] 3 [16]
	normal[0xea] = []func(*CPU){readParamA, readParamB, nop, ldUX16A}

	// XOR d8  [Z 0 0 0] 2 [8]
	normal[0xee] = []func(*CPU){readParamA, xorU}

	// RST 28H  [] 1 [16]
	normal[0xef] = []func(*CPU){nop, rst(0x0028), push(&m8b), push(&m8a)}

	// LDH A (a8)   2 [12]
	normal[0xf0] = []func(*CPU){readParamA, nop, ldAUX}

	// POP AF  [Z N H C] 1 [12]
	normal[0xf1] = []func(*CPU){nop, popF, pop(&a)}

	// LD A (C)     1 [8]
	normal[0xf2] = []func(*CPU){nop, ldACX}

	// DI   [] 1 [4]
	normal[0xf3] = []func(*CPU){di}

	// PUSH AF      1 [16]
	normal[0xf5] = []func(*CPU){nop, nop, push(&a), push(&f)}

	// OR d8  [Z 0 0 0] 2 [8]
	normal[0xf6] = []func(*CPU){readParamA, orU}

	// RST 30H  [] 1 [16]
	normal[0xf7] = []func(*CPU){nop, rst(0x0030), push(&m8b), push(&m8a)}

	// LD HL SP+r8 [0 0 H C] 2 [12]
	normal[0xf8] = []func(*CPU){readParamA, nop, ldHLSP}

	// LD SP HL [] 1 [8]
	normal[0xf9] = []func(*CPU){nop, ldSPHL}

	// LD A (a16) [] 3 [16]
	normal[0xfa] = []func(*CPU){readParamA, readParamB, nop, ldAUX16}

	// EI   [] 1 [4]
	normal[0xfb] = []func(*CPU){ei}

	// CP d8  [Z 1 H C] 2 [8]
	normal[0xfe] = []func(*CPU){readParamA, cpU}

	// RST 38H  [] 1 [16]
	normal[0xff] = []func(*CPU){nop, rst(0x0038), push(&m8b), push(&m8a)}

	// RLC B  [Z 0 0 C] 2 [8]
	prefix[0x00] = []func(*CPU){nop, rlcB}

	// RLC C  [Z 0 0 C] 2 [8]
	prefix[0x01] = []func(*CPU){nop, rlcC}

	// RLC D  [Z 0 0 C] 2 [8]
	prefix[0x02] = []func(*CPU){nop, rlcD}

	// RLC E  [Z 0 0 C] 2 [8]
	prefix[0x03] = []func(*CPU){nop, rlcE}

	// RLC H  [Z 0 0 C] 2 [8]
	prefix[0x04] = []func(*CPU){nop, rlcH}

	// RLC L  [Z 0 0 C] 2 [8]
	prefix[0x05] = []func(*CPU){nop, rlcL}

	// RLC (HL)  [Z 0 0 C] 2 [16]
	prefix[0x06] = []func(*CPU){nop, nop, ldMHL, rlcM}

	// RLC A  [Z 0 0 C] 2 [8]
	prefix[0x07] = []func(*CPU){nop, rlcA}

	// RRC B  [Z 0 0 C] 2 [8]
	prefix[0x08] = []func(*CPU){nop, rrcB}

	// RRC C  [Z 0 0 C] 2 [8]
	prefix[0x09] = []func(*CPU){nop, rrcC}

	// RRC D  [Z 0 0 C] 2 [8]
	prefix[0x0a] = []func(*CPU){nop, rrcD}

	// RRC E  [Z 0 0 C] 2 [8]
	prefix[0x0b] = []func(*CPU){nop, rrcE}

	// RRC H  [Z 0 0 C] 2 [8]
	prefix[0x0c] = []func(*CPU){nop, rrcH}

	// RRC L  [Z 0 0 C] 2 [8]
	prefix[0x0d] = []func(*CPU){nop, rrcL}

	// RRC (HL)  [Z 0 0 C] 2 [16]
	prefix[0x0e] = []func(*CPU){nop, nop, ldMHL, rrcM}

	// RRC A  [Z 0 0 C] 2 [8]
	prefix[0x0f] = []func(*CPU){nop, rrcA}

	// RL B  [Z 0 0 C] 2 [8]
	prefix[0x10] = []func(*CPU){nop, rlB}

	// RL C  [Z 0 0 C] 2 [8]
	prefix[0x11] = []func(*CPU){nop, rlC}

	// RL D  [Z 0 0 C] 2 [8]
	prefix[0x12] = []func(*CPU){nop, rlD}

	// RL E  [Z 0 0 C] 2 [8]
	prefix[0x13] = []func(*CPU){nop, rlE}

	// RL H  [Z 0 0 C] 2 [8]
	prefix[0x14] = []func(*CPU){nop, rlH}

	// RL L  [Z 0 0 C] 2 [8]
	prefix[0x15] = []func(*CPU){nop, rlL}

	// RL (HL)  [Z 0 0 C] 2 [16]
	prefix[0x16] = []func(*CPU){nop, nop, ldMHL, rlM}

	// RL A  [Z 0 0 C] 2 [8]
	prefix[0x17] = []func(*CPU){nop, rlA}

	// RR B  [Z 0 0 C] 2 [8]
	prefix[0x18] = []func(*CPU){nop, rrB}

	// RR C  [Z 0 0 C] 2 [8]
	prefix[0x19] = []func(*CPU){nop, rrC}

	// RR D  [Z 0 0 C] 2 [8]
	prefix[0x1a] = []func(*CPU){nop, rrD}

	// RR E  [Z 0 0 C] 2 [8]
	prefix[0x1b] = []func(*CPU){nop, rrE}

	// RR H  [Z 0 0 C] 2 [8]
	prefix[0x1c] = []func(*CPU){nop, rrH}

	// RR L  [Z 0 0 C] 2 [8]
	prefix[0x1d] = []func(*CPU){nop, rrL}

	// RR (HL)  [Z 0 0 C] 2 [16]
	prefix[0x1e] = []func(*CPU){nop, nop, ldMHL, rrM}

	// RR A  [Z 0 0 C] 2 [8]
	prefix[0x1f] = []func(*CPU){nop, rrA}

	// SLA B  [Z 0 0 C] 2 [8]
	prefix[0x20] = []func(*CPU){nop, slaB}

	// SLA C  [Z 0 0 C] 2 [8]
	prefix[0x21] = []func(*CPU){nop, slaC}

	// SLA D  [Z 0 0 C] 2 [8]
	prefix[0x22] = []func(*CPU){nop, slaD}

	// SLA E  [Z 0 0 C] 2 [8]
	prefix[0x23] = []func(*CPU){nop, slaE}

	// SLA H  [Z 0 0 C] 2 [8]
	prefix[0x24] = []func(*CPU){nop, slaH}

	// SLA L  [Z 0 0 C] 2 [8]
	prefix[0x25] = []func(*CPU){nop, slaL}

	// SLA (HL)  [Z 0 0 C] 2 [16]
	prefix[0x26] = []func(*CPU){nop, nop, ldMHL, slaM}

	// SLA A  [Z 0 0 C] 2 [8]
	prefix[0x27] = []func(*CPU){nop, slaA}

	// SRA B  [Z 0 0 C] 2 [8]
	prefix[0x28] = []func(*CPU){nop, sraB}

	// SRA C  [Z 0 0 C] 2 [8]
	prefix[0x29] = []func(*CPU){nop, sraC}

	// SRA D  [Z 0 0 C] 2 [8]
	prefix[0x2a] = []func(*CPU){nop, sraD}

	// SRA E  [Z 0 0 C] 2 [8]
	prefix[0x2b] = []func(*CPU){nop, sraE}

	// SRA H  [Z 0 0 C] 2 [8]
	prefix[0x2c] = []func(*CPU){nop, sraH}

	// SRA L  [Z 0 0 C] 2 [8]
	prefix[0x2d] = []func(*CPU){nop, sraL}

	// SRA (HL)  [Z 0 0 C] 2 [16]
	prefix[0x2e] = []func(*CPU){nop, nop, ldMHL, sraM}

	// SRA A  [Z 0 0 C] 2 [8]
	prefix[0x2f] = []func(*CPU){nop, sraA}

	// SWAP B  [Z 0 0 0] 2 [8]
	prefix[0x30] = []func(*CPU){nop, swapB}

	// SWAP C  [Z 0 0 0] 2 [8]
	prefix[0x31] = []func(*CPU){nop, swapC}

	// SWAP D  [Z 0 0 0] 2 [8]
	prefix[0x32] = []func(*CPU){nop, swapD}

	// SWAP E  [Z 0 0 0] 2 [8]
	prefix[0x33] = []func(*CPU){nop, swapE}

	// SWAP H  [Z 0 0 0] 2 [8]
	prefix[0x34] = []func(*CPU){nop, swapH}

	// SWAP L  [Z 0 0 0] 2 [8]
	prefix[0x35] = []func(*CPU){nop, swapL}

	// SWAP (HL)  [Z 0 0 0] 2 [16]
	prefix[0x36] = []func(*CPU){nop, nop, ldMHL, swapM}

	// SWAP A  [Z 0 0 0] 2 [8]
	prefix[0x37] = []func(*CPU){nop, swapA}

	// SRL B  [Z 0 0 C] 2 [8]
	prefix[0x38] = []func(*CPU){nop, srlB}

	// SRL C  [Z 0 0 C] 2 [8]
	prefix[0x39] = []func(*CPU){nop, srlC}

	// SRL D  [Z 0 0 C] 2 [8]
	prefix[0x3a] = []func(*CPU){nop, srlD}

	// SRL E  [Z 0 0 C] 2 [8]
	prefix[0x3b] = []func(*CPU){nop, srlE}

	// SRL H  [Z 0 0 C] 2 [8]
	prefix[0x3c] = []func(*CPU){nop, srlH}

	// SRL L  [Z 0 0 C] 2 [8]
	prefix[0x3d] = []func(*CPU){nop, srlL}

	// SRL (HL)  [Z 0 0 C] 2 [16]
	prefix[0x3e] = []func(*CPU){nop, nop, ldMHL, srlM}

	// SRL A  [Z 0 0 C] 2 [8]
	prefix[0x3f] = []func(*CPU){nop, srlA}

	// BIT 0 B [Z 0 1 -] 2 [8]
	prefix[0x40] = []func(*CPU){nop, bit(0, &b)}

	// BIT 0 C [Z 0 1 -] 2 [8]
	prefix[0x41] = []func(*CPU){nop, bit(0, &c)}

	// BIT 0 D [Z 0 1 -] 2 [8]
	prefix[0x42] = []func(*CPU){nop, bit(0, &d)}

	// BIT 0 E [Z 0 1 -] 2 [8]
	prefix[0x43] = []func(*CPU){nop, bit(0, &e)}

	// BIT 0 H [Z 0 1 -] 2 [8]
	prefix[0x44] = []func(*CPU){nop, bit(0, &h)}

	// BIT 0 L [Z 0 1 -] 2 [8]
	prefix[0x45] = []func(*CPU){nop, bit(0, &l)}

	// BIT 0 (HL) [Z 0 1 -] 2 [12]
	prefix[0x46] = []func(*CPU){nop, nop, bitM(0)}

	// BIT 0 A [Z 0 1 -] 2 [8]
	prefix[0x47] = []func(*CPU){nop, bit(0, &a)}

	// BIT 1 B [Z 0 1 -] 2 [8]
	prefix[0x48] = []func(*CPU){nop, bit(1, &b)}

	// BIT 1 C [Z 0 1 -] 2 [8]
	prefix[0x49] = []func(*CPU){nop, bit(1, &c)}

	// BIT 1 D [Z 0 1 -] 2 [8]
	prefix[0x4a] = []func(*CPU){nop, bit(1, &d)}

	// BIT 1 E [Z 0 1 -] 2 [8]
	prefix[0x4b] = []func(*CPU){nop, bit(1, &e)}

	// BIT 1 H [Z 0 1 -] 2 [8]
	prefix[0x4c] = []func(*CPU){nop, bit(1, &h)}

	// BIT 1 L [Z 0 1 -] 2 [8]
	prefix[0x4d] = []func(*CPU){nop, bit(1, &l)}

	// BIT 1 (HL) [Z 0 1 -] 2 [12]
	prefix[0x4e] = []func(*CPU){nop, nop, bitM(1)}

	// BIT 1 A [Z 0 1 -] 2 [8]
	prefix[0x4f] = []func(*CPU){nop, bit(1, &a)}

	// BIT 2 B [Z 0 1 -] 2 [8]
	prefix[0x50] = []func(*CPU){nop, bit(2, &b)}

	// BIT 2 C [Z 0 1 -] 2 [8]
	prefix[0x51] = []func(*CPU){nop, bit(2, &c)}

	// BIT 2 D [Z 0 1 -] 2 [8]
	prefix[0x52] = []func(*CPU){nop, bit(2, &d)}

	// BIT 2 E [Z 0 1 -] 2 [8]
	prefix[0x53] = []func(*CPU){nop, bit(2, &e)}

	// BIT 2 H [Z 0 1 -] 2 [8]
	prefix[0x54] = []func(*CPU){nop, bit(2, &h)}

	// BIT 2 L [Z 0 1 -] 2 [8]
	prefix[0x55] = []func(*CPU){nop, bit(2, &l)}

	// BIT 2 (HL) [Z 0 1 -] 2 [12]
	prefix[0x56] = []func(*CPU){nop, nop, bitM(2)}

	// BIT 2 A [Z 0 1 -] 2 [8]
	prefix[0x57] = []func(*CPU){nop, bit(2, &a)}

	// BIT 3 B [Z 0 1 -] 2 [8]
	prefix[0x58] = []func(*CPU){nop, bit(3, &b)}

	// BIT 3 C [Z 0 1 -] 2 [8]
	prefix[0x59] = []func(*CPU){nop, bit(3, &c)}

	// BIT 3 D [Z 0 1 -] 2 [8]
	prefix[0x5a] = []func(*CPU){nop, bit(3, &d)}

	// BIT 3 E [Z 0 1 -] 2 [8]
	prefix[0x5b] = []func(*CPU){nop, bit(3, &e)}

	// BIT 3 H [Z 0 1 -] 2 [8]
	prefix[0x5c] = []func(*CPU){nop, bit(3, &h)}

	// BIT 3 L [Z 0 1 -] 2 [8]
	prefix[0x5d] = []func(*CPU){nop, bit(3, &l)}

	// BIT 3 (HL) [Z 0 1 -] 2 [12]
	prefix[0x5e] = []func(*CPU){nop, nop, bitM(3)}

	// BIT 3 A [Z 0 1 -] 2 [8]
	prefix[0x5f] = []func(*CPU){nop, bit(3, &a)}

	// BIT 4 B [Z 0 1 -] 2 [8]
	prefix[0x60] = []func(*CPU){nop, bit(4, &b)}

	// BIT 4 C [Z 0 1 -] 2 [8]
	prefix[0x61] = []func(*CPU){nop, bit(4, &c)}

	// BIT 4 D [Z 0 1 -] 2 [8]
	prefix[0x62] = []func(*CPU){nop, bit(4, &d)}

	// BIT 4 E [Z 0 1 -] 2 [8]
	prefix[0x63] = []func(*CPU){nop, bit(4, &e)}

	// BIT 4 H [Z 0 1 -] 2 [8]
	prefix[0x64] = []func(*CPU){nop, bit(4, &h)}

	// BIT 4 L [Z 0 1 -] 2 [8]
	prefix[0x65] = []func(*CPU){nop, bit(4, &l)}

	// BIT 4 (HL) [Z 0 1 -] 2 [12]
	prefix[0x66] = []func(*CPU){nop, nop, bitM(4)}

	// BIT 4 A [Z 0 1 -] 2 [8]
	prefix[0x67] = []func(*CPU){nop, bit(4, &a)}

	// BIT 5 B [Z 0 1 -] 2 [8]
	prefix[0x68] = []func(*CPU){nop, bit(5, &b)}

	// BIT 5 C [Z 0 1 -] 2 [8]
	prefix[0x69] = []func(*CPU){nop, bit(5, &c)}

	// BIT 5 D [Z 0 1 -] 2 [8]
	prefix[0x6a] = []func(*CPU){nop, bit(5, &d)}

	// BIT 5 E [Z 0 1 -] 2 [8]
	prefix[0x6b] = []func(*CPU){nop, bit(5, &e)}

	// BIT 5 H [Z 0 1 -] 2 [8]
	prefix[0x6c] = []func(*CPU){nop, bit(5, &h)}

	// BIT 5 L [Z 0 1 -] 2 [8]
	prefix[0x6d] = []func(*CPU){nop, bit(5, &l)}

	// BIT 5 (HL) [Z 0 1 -] 2 [12]
	prefix[0x6e] = []func(*CPU){nop, nop, bitM(5)}

	// BIT 5 A [Z 0 1 -] 2 [8]
	prefix[0x6f] = []func(*CPU){nop, bit(5, &a)}

	// BIT 6 B [Z 0 1 -] 2 [8]
	prefix[0x70] = []func(*CPU){nop, bit(6, &b)}

	// BIT 6 C [Z 0 1 -] 2 [8]
	prefix[0x71] = []func(*CPU){nop, bit(6, &c)}

	// BIT 6 D [Z 0 1 -] 2 [8]
	prefix[0x72] = []func(*CPU){nop, bit(6, &d)}

	// BIT 6 E [Z 0 1 -] 2 [8]
	prefix[0x73] = []func(*CPU){nop, bit(6, &e)}

	// BIT 6 H [Z 0 1 -] 2 [8]
	prefix[0x74] = []func(*CPU){nop, bit(6, &h)}

	// BIT 6 L [Z 0 1 -] 2 [8]
	prefix[0x75] = []func(*CPU){nop, bit(6, &l)}

	// BIT 6 (HL) [Z 0 1 -] 2 [12]
	prefix[0x76] = []func(*CPU){nop, nop, bitM(6)}

	// BIT 6 A [Z 0 1 -] 2 [8]
	prefix[0x77] = []func(*CPU){nop, bit(6, &a)}

	// BIT 7 B [Z 0 1 -] 2 [8]
	prefix[0x78] = []func(*CPU){nop, bit(7, &b)}

	// BIT 7 C [Z 0 1 -] 2 [8]
	prefix[0x79] = []func(*CPU){nop, bit(7, &c)}

	// BIT 7 D [Z 0 1 -] 2 [8]
	prefix[0x7a] = []func(*CPU){nop, bit(7, &d)}

	// BIT 7 E [Z 0 1 -] 2 [8]
	prefix[0x7b] = []func(*CPU){nop, bit(7, &e)}

	// BIT 7 H [Z 0 1 -] 2 [8]
	prefix[0x7c] = []func(*CPU){nop, bit(7, &h)}

	// BIT 7 L [Z 0 1 -] 2 [8]
	prefix[0x7d] = []func(*CPU){nop, bit(7, &l)}

	// BIT 7 (HL) [Z 0 1 -] 2 [12]
	prefix[0x7e] = []func(*CPU){nop, nop, bitM(7)}

	// BIT 7 A [Z 0 1 -] 2 [8]
	prefix[0x7f] = []func(*CPU){nop, bit(7, &a)}

	// RES 0 B [] 2 [8]
	prefix[0x80] = []func(*CPU){nop, res(0, &b)}

	// RES 0 C [] 2 [8]
	prefix[0x81] = []func(*CPU){nop, res(0, &c)}

	// RES 0 D [] 2 [8]
	prefix[0x82] = []func(*CPU){nop, res(0, &d)}

	// RES 0 E [] 2 [8]
	prefix[0x83] = []func(*CPU){nop, res(0, &e)}

	// RES 0 H [] 2 [8]
	prefix[0x84] = []func(*CPU){nop, res(0, &h)}

	// RES 0 L [] 2 [8]
	prefix[0x85] = []func(*CPU){nop, res(0, &l)}

	// RES 0 (HL) [] 2 [16]
	prefix[0x86] = []func(*CPU){nop, nop, ldMHL, resM(0)}

	// RES 0 A [] 2 [8]
	prefix[0x87] = []func(*CPU){nop, res(0, &a)}

	// RES 1 B [] 2 [8]
	prefix[0x88] = []func(*CPU){nop, res(1, &b)}

	// RES 1 C [] 2 [8]
	prefix[0x89] = []func(*CPU){nop, res(1, &c)}

	// RES 1 D [] 2 [8]
	prefix[0x8a] = []func(*CPU){nop, res(1, &d)}

	// RES 1 E [] 2 [8]
	prefix[0x8b] = []func(*CPU){nop, res(1, &e)}

	// RES 1 H [] 2 [8]
	prefix[0x8c] = []func(*CPU){nop, res(1, &h)}

	// RES 1 L [] 2 [8]
	prefix[0x8d] = []func(*CPU){nop, res(1, &l)}

	// RES 1 (HL) [] 2 [16]
	prefix[0x8e] = []func(*CPU){nop, nop, ldMHL, resM(1)}

	// RES 1 A [] 2 [8]
	prefix[0x8f] = []func(*CPU){nop, res(1, &a)}

	// RES 2 B [] 2 [8]
	prefix[0x90] = []func(*CPU){nop, res(2, &b)}

	// RES 2 C [] 2 [8]
	prefix[0x91] = []func(*CPU){nop, res(2, &c)}

	// RES 2 D [] 2 [8]
	prefix[0x92] = []func(*CPU){nop, res(2, &d)}

	// RES 2 E [] 2 [8]
	prefix[0x93] = []func(*CPU){nop, res(2, &e)}

	// RES 2 H [] 2 [8]
	prefix[0x94] = []func(*CPU){nop, res(2, &h)}

	// RES 2 L [] 2 [8]
	prefix[0x95] = []func(*CPU){nop, res(2, &l)}

	// RES 2 (HL) [] 2 [16]
	prefix[0x96] = []func(*CPU){nop, nop, ldMHL, resM(2)}

	// RES 2 A [] 2 [8]
	prefix[0x97] = []func(*CPU){nop, res(2, &a)}

	// RES 3 B [] 2 [8]
	prefix[0x98] = []func(*CPU){nop, res(3, &b)}

	// RES 3 C [] 2 [8]
	prefix[0x99] = []func(*CPU){nop, res(3, &c)}

	// RES 3 D [] 2 [8]
	prefix[0x9a] = []func(*CPU){nop, res(3, &d)}

	// RES 3 E [] 2 [8]
	prefix[0x9b] = []func(*CPU){nop, res(3, &e)}

	// RES 3 H [] 2 [8]
	prefix[0x9c] = []func(*CPU){nop, res(3, &h)}

	// RES 3 L [] 2 [8]
	prefix[0x9d] = []func(*CPU){nop, res(3, &l)}

	// RES 3 (HL) [] 2 [16]
	prefix[0x9e] = []func(*CPU){nop, nop, ldMHL, resM(3)}

	// RES 3 A [] 2 [8]
	prefix[0x9f] = []func(*CPU){nop, res(3, &a)}

	// RES 4 B [] 2 [8]
	prefix[0xa0] = []func(*CPU){nop, res(4, &b)}

	// RES 4 C [] 2 [8]
	prefix[0xa1] = []func(*CPU){nop, res(4, &c)}

	// RES 4 D [] 2 [8]
	prefix[0xa2] = []func(*CPU){nop, res(4, &d)}

	// RES 4 E [] 2 [8]
	prefix[0xa3] = []func(*CPU){nop, res(4, &e)}

	// RES 4 H [] 2 [8]
	prefix[0xa4] = []func(*CPU){nop, res(4, &h)}

	// RES 4 L [] 2 [8]
	prefix[0xa5] = []func(*CPU){nop, res(4, &l)}

	// RES 4 (HL) [] 2 [16]
	prefix[0xa6] = []func(*CPU){nop, nop, ldMHL, resM(4)}

	// RES 4 A [] 2 [8]
	prefix[0xa7] = []func(*CPU){nop, res(4, &a)}

	// RES 5 B [] 2 [8]
	prefix[0xa8] = []func(*CPU){nop, res(5, &b)}

	// RES 5 C [] 2 [8]
	prefix[0xa9] = []func(*CPU){nop, res(5, &c)}

	// RES 5 D [] 2 [8]
	prefix[0xaa] = []func(*CPU){nop, res(5, &d)}

	// RES 5 E [] 2 [8]
	prefix[0xab] = []func(*CPU){nop, res(5, &e)}

	// RES 5 H [] 2 [8]
	prefix[0xac] = []func(*CPU){nop, res(5, &h)}

	// RES 5 L [] 2 [8]
	prefix[0xad] = []func(*CPU){nop, res(5, &l)}

	// RES 5 (HL) [] 2 [16]
	prefix[0xae] = []func(*CPU){nop, nop, ldMHL, resM(5)}

	// RES 5 A [] 2 [8]
	prefix[0xaf] = []func(*CPU){nop, res(5, &a)}

	// RES 6 B [] 2 [8]
	prefix[0xb0] = []func(*CPU){nop, res(6, &b)}

	// RES 6 C [] 2 [8]
	prefix[0xb1] = []func(*CPU){nop, res(6, &c)}

	// RES 6 D [] 2 [8]
	prefix[0xb2] = []func(*CPU){nop, res(6, &d)}

	// RES 6 E [] 2 [8]
	prefix[0xb3] = []func(*CPU){nop, res(6, &e)}

	// RES 6 H [] 2 [8]
	prefix[0xb4] = []func(*CPU){nop, res(6, &h)}

	// RES 6 L [] 2 [8]
	prefix[0xb5] = []func(*CPU){nop, res(6, &l)}

	// RES 6 (HL) [] 2 [16]
	prefix[0xb6] = []func(*CPU){nop, nop, ldMHL, resM(6)}

	// RES 6 A [] 2 [8]
	prefix[0xb7] = []func(*CPU){nop, res(6, &a)}

	// RES 7 B [] 2 [8]
	prefix[0xb8] = []func(*CPU){nop, res(7, &b)}

	// RES 7 C [] 2 [8]
	prefix[0xb9] = []func(*CPU){nop, res(7, &c)}

	// RES 7 D [] 2 [8]
	prefix[0xba] = []func(*CPU){nop, res(7, &d)}

	// RES 7 E [] 2 [8]
	prefix[0xbb] = []func(*CPU){nop, res(7, &e)}

	// RES 7 H [] 2 [8]
	prefix[0xbc] = []func(*CPU){nop, res(7, &h)}

	// RES 7 L [] 2 [8]
	prefix[0xbd] = []func(*CPU){nop, res(7, &l)}

	// RES 7 (HL) [] 2 [16]
	prefix[0xbe] = []func(*CPU){nop, nop, ldMHL, resM(7)}

	// RES 7 A [] 2 [8]
	prefix[0xbf] = []func(*CPU){nop, res(7, &a)}

	// SET 0 B [] 2 [8]
	prefix[0xc0] = []func(*CPU){nop, set(0, &b)}

	// SET 0 C [] 2 [8]
	prefix[0xc1] = []func(*CPU){nop, set(0, &c)}

	// SET 0 D [] 2 [8]
	prefix[0xc2] = []func(*CPU){nop, set(0, &d)}

	// SET 0 E [] 2 [8]
	prefix[0xc3] = []func(*CPU){nop, set(0, &e)}

	// SET 0 H [] 2 [8]
	prefix[0xc4] = []func(*CPU){nop, set(0, &h)}

	// SET 0 L [] 2 [8]
	prefix[0xc5] = []func(*CPU){nop, set(0, &l)}

	// SET 0 (HL) [] 2 [16]
	prefix[0xc6] = []func(*CPU){nop, nop, ldMHL, setM(0)}

	// SET 0 A [] 2 [8]
	prefix[0xc7] = []func(*CPU){nop, set(0, &a)}

	// SET 1 B [] 2 [8]
	prefix[0xc8] = []func(*CPU){nop, set(1, &b)}

	// SET 1 C [] 2 [8]
	prefix[0xc9] = []func(*CPU){nop, set(1, &c)}

	// SET 1 D [] 2 [8]
	prefix[0xca] = []func(*CPU){nop, set(1, &d)}

	// SET 1 E [] 2 [8]
	prefix[0xcb] = []func(*CPU){nop, set(1, &e)}

	// SET 1 H [] 2 [8]
	prefix[0xcc] = []func(*CPU){nop, set(1, &h)}

	// SET 1 L [] 2 [8]
	prefix[0xcd] = []func(*CPU){nop, set(1, &l)}

	// SET 1 (HL) [] 2 [16]
	prefix[0xce] = []func(*CPU){nop, nop, ldMHL, setM(1)}

	// SET 1 A [] 2 [8]
	prefix[0xcf] = []func(*CPU){nop, set(1, &a)}

	// SET 2 B [] 2 [8]
	prefix[0xd0] = []func(*CPU){nop, set(2, &b)}

	// SET 2 C [] 2 [8]
	prefix[0xd1] = []func(*CPU){nop, set(2, &c)}

	// SET 2 D [] 2 [8]
	prefix[0xd2] = []func(*CPU){nop, set(2, &d)}

	// SET 2 E [] 2 [8]
	prefix[0xd3] = []func(*CPU){nop, set(2, &e)}

	// SET 2 H [] 2 [8]
	prefix[0xd4] = []func(*CPU){nop, set(2, &h)}

	// SET 2 L [] 2 [8]
	prefix[0xd5] = []func(*CPU){nop, set(2, &l)}

	// SET 2 (HL) [] 2 [16]
	prefix[0xd6] = []func(*CPU){nop, nop, ldMHL, setM(2)}

	// SET 2 A [] 2 [8]
	prefix[0xd7] = []func(*CPU){nop, set(2, &a)}

	// SET 3 B [] 2 [8]
	prefix[0xd8] = []func(*CPU){nop, set(3, &b)}

	// SET 3 C [] 2 [8]
	prefix[0xd9] = []func(*CPU){nop, set(3, &c)}

	// SET 3 D [] 2 [8]
	prefix[0xda] = []func(*CPU){nop, set(3, &d)}

	// SET 3 E [] 2 [8]
	prefix[0xdb] = []func(*CPU){nop, set(3, &e)}

	// SET 3 H [] 2 [8]
	prefix[0xdc] = []func(*CPU){nop, set(3, &h)}

	// SET 3 L [] 2 [8]
	prefix[0xdd] = []func(*CPU){nop, set(3, &l)}

	// SET 3 (HL) [] 2 [16]
	prefix[0xde] = []func(*CPU){nop, nop, ldMHL, setM(3)}

	// SET 3 A [] 2 [8]
	prefix[0xdf] = []func(*CPU){nop, set(3, &a)}

	// SET 4 B [] 2 [8]
	prefix[0xe0] = []func(*CPU){nop, set(4, &b)}

	// SET 4 C [] 2 [8]
	prefix[0xe1] = []func(*CPU){nop, set(4, &c)}

	// SET 4 D [] 2 [8]
	prefix[0xe2] = []func(*CPU){nop, set(4, &d)}

	// SET 4 E [] 2 [8]
	prefix[0xe3] = []func(*CPU){nop, set(4, &e)}

	// SET 4 H [] 2 [8]
	prefix[0xe4] = []func(*CPU){nop, set(4, &h)}

	// SET 4 L [] 2 [8]
	prefix[0xe5] = []func(*CPU){nop, set(4, &l)}

	// SET 4 (HL) [] 2 [16]
	prefix[0xe6] = []func(*CPU){nop, nop, ldMHL, setM(4)}

	// SET 4 A [] 2 [8]
	prefix[0xe7] = []func(*CPU){nop, set(4, &a)}

	// SET 5 B [] 2 [8]
	prefix[0xe8] = []func(*CPU){nop, set(5, &b)}

	// SET 5 C [] 2 [8]
	prefix[0xe9] = []func(*CPU){nop, set(5, &c)}

	// SET 5 D [] 2 [8]
	prefix[0xea] = []func(*CPU){nop, set(5, &d)}

	// SET 5 E [] 2 [8]
	prefix[0xeb] = []func(*CPU){nop, set(5, &e)}

	// SET 5 H [] 2 [8]
	prefix[0xec] = []func(*CPU){nop, set(5, &h)}

	// SET 5 L [] 2 [8]
	prefix[0xed] = []func(*CPU){nop, set(5, &l)}

	// SET 5 (HL) [] 2 [16]
	prefix[0xee] = []func(*CPU){nop, nop, ldMHL, setM(5)}

	// SET 5 A [] 2 [8]
	prefix[0xef] = []func(*CPU){nop, set(5, &a)}

	// SET 6 B [] 2 [8]
	prefix[0xf0] = []func(*CPU){nop, set(6, &b)}

	// SET 6 C [] 2 [8]
	prefix[0xf1] = []func(*CPU){nop, set(6, &c)}

	// SET 6 D [] 2 [8]
	prefix[0xf2] = []func(*CPU){nop, set(6, &d)}

	// SET 6 E [] 2 [8]
	prefix[0xf3] = []func(*CPU){nop, set(6, &e)}

	// SET 6 H [] 2 [8]
	prefix[0xf4] = []func(*CPU){nop, set(6, &h)}

	// SET 6 L [] 2 [8]
	prefix[0xf5] = []func(*CPU){nop, set(6, &l)}

	// SET 6 (HL) [] 2 [16]
	prefix[0xf6] = []func(*CPU){nop, nop, ldMHL, setM(6)}

	// SET 6 A [] 2 [8]
	prefix[0xf7] = []func(*CPU){nop, set(6, &a)}

	// SET 7 B [] 2 [8]
	prefix[0xf8] = []func(*CPU){nop, set(7, &b)}

	// SET 7 C [] 2 [8]
	prefix[0xf9] = []func(*CPU){nop, set(7, &c)}

	// SET 7 D [] 2 [8]
	prefix[0xfa] = []func(*CPU){nop, set(7, &d)}

	// SET 7 E [] 2 [8]
	prefix[0xfb] = []func(*CPU){nop, set(7, &e)}

	// SET 7 H [] 2 [8]
	prefix[0xfc] = []func(*CPU){nop, set(7, &h)}

	// SET 7 L [] 2 [8]
	prefix[0xfd] = []func(*CPU){nop, set(7, &l)}

	// SET 7 (HL) [] 2 [16]
	prefix[0xfe] = []func(*CPU){nop, nop, ldMHL, setM(7)}

	// SET 7 A [] 2 [8]
	prefix[0xff] = []func(*CPU){nop, set(7, &a)}
}
