package cpu

var normal [256][]func()
var prefix [256][]func()
var veryShortInterrupt []func()
var shortInterrupt []func()
var longInterrupt []func()
var earlyCheck [256]func(int) bool

func nop() {
	// Do nothing
}

func isFinishedSpecial(check func() bool, early, last int) func(int) bool {
	return func(currentCycle int) bool {
		// fmt.Println("isFinishedSpecial", check(), early, last, currentCycle)
		return (currentCycle == early && check()) || currentCycle == last
	}
}

func isFinished(last int) func(int) bool {
	return func(currentCycle int) bool {
		// fmt.Println("isFinished", last)
		return currentCycle == last
	}
}

func (cpu *CPU) Initialize() {

	cpu.next()

	mapper := cpu.mapper

	veryShortInterrupt = []func(){cpu.handleInterrupt}
	shortInterrupt = []func(){nop, nop, nop, nop, cpu.handleInterrupt}
	longInterrupt = []func(){nop, nop, nop, nop, nop, cpu.handleInterrupt}

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

	earlyCheck[0x20] = isFinishedSpecial(cpu.zf, 1, 2)
	earlyCheck[0x28] = isFinishedSpecial(cpu.nzf, 1, 2)
	earlyCheck[0x30] = isFinishedSpecial(cpu.cf, 1, 2)
	earlyCheck[0x38] = isFinishedSpecial(cpu.ncf, 1, 2)
	earlyCheck[0xc0] = isFinishedSpecial(cpu.zf, 1, 4)
	earlyCheck[0xc2] = isFinishedSpecial(cpu.zf, 2, 3)
	earlyCheck[0xc4] = isFinishedSpecial(cpu.zf, 2, 5)
	earlyCheck[0xc8] = isFinishedSpecial(cpu.nzf, 1, 4)
	earlyCheck[0xca] = isFinishedSpecial(cpu.nzf, 2, 3)
	earlyCheck[0xcc] = isFinishedSpecial(cpu.nzf, 2, 5)
	earlyCheck[0xd0] = isFinishedSpecial(cpu.cf, 1, 4)
	earlyCheck[0xd2] = isFinishedSpecial(cpu.cf, 2, 3)
	earlyCheck[0xd4] = isFinishedSpecial(cpu.cf, 2, 5)
	earlyCheck[0xd8] = isFinishedSpecial(cpu.ncf, 1, 4)
	earlyCheck[0xda] = isFinishedSpecial(cpu.ncf, 2, 3)
	earlyCheck[0xdc] = isFinishedSpecial(cpu.ncf, 2, 5)

	// NOP          1 [4]
	normal[0x00] = []func(){nop}

	// LD BC d16 [] 3 [12]
	normal[0x01] = []func(){cpu.readParamA, cpu.readParamB, cpu.ldBCU16}

	// LD (BC) A [] 1 [8]
	normal[0x02] = []func(){nop, cpu.ldBCA(mapper)}

	// INC BC  [] 1 [8]
	normal[0x03] = []func(){nop, cpu.incBC}

	// INC B  [Z 0 H -] 1 [4]
	normal[0x04] = []func(){cpu.incB}

	// DEC B  [Z 1 H -] 1 [4]
	normal[0x05] = []func(){cpu.decB}

	// LD B d8 [] 2 [8]
	normal[0x06] = []func(){cpu.readParamA, cpu.ldBU}

	// RLCA   [0 0 0 C] 1 [4]
	normal[0x07] = []func(){cpu.rlca}

	// LD (a16) SP [] 3 [20]
	normal[0x08] = []func(){cpu.readParamA, cpu.readParamB, nop, cpu.writeLowSP(mapper), cpu.writeHighSP(mapper)}

	// ADD HL BC [- 0 H C] 1 [8]
	normal[0x09] = []func(){nop, cpu.addHLBC}

	// LD A (BC) [] 1 [8]
	normal[0x0a] = []func(){nop, cpu.ldABC(mapper)}

	// DEC BC  [] 1 [8]
	normal[0x0b] = []func(){nop, cpu.decBC}

	// INC C  [Z 0 H -] 1 [4]
	normal[0x0c] = []func(){cpu.incC}

	// DEC C  [Z 1 H -] 1 [4]
	normal[0x0d] = []func(){cpu.decC}

	// LD C d8 [] 2 [8]
	normal[0x0e] = []func(){cpu.readParamA, cpu.ldCU}

	// RRCA   [0 0 0 C] 1 [4]
	normal[0x0f] = []func(){cpu.rrca}

	// STOP 0  [] 1 [4]
	normal[0x10] = []func(){cpu.stop}

	// LD DE d16 [] 3 [12]
	normal[0x11] = []func(){cpu.readParamA, cpu.readParamB, cpu.ldDEU16}

	// LD (DE) A [] 1 [8]
	normal[0x12] = []func(){nop, cpu.ldDEA(mapper)}

	// INC DE  [] 1 [8]
	normal[0x13] = []func(){nop, cpu.incDE}

	// INC D  [Z 0 H -] 1 [4]
	normal[0x14] = []func(){cpu.incD}

	// DEC D  [Z 1 H -] 1 [4]
	normal[0x15] = []func(){cpu.decD}

	// LD D d8 [] 2 [8]
	normal[0x16] = []func(){cpu.readParamA, cpu.ldDU}

	// RLA   [0 0 0 C] 1 [4]
	normal[0x17] = []func(){cpu.rla}

	// JR r8  [] 2 [12]
	normal[0x18] = []func(){nop, cpu.readParamA, cpu.jr}

	// ADD HL DE [- 0 H C] 1 [8]
	normal[0x19] = []func(){nop, cpu.addHLDE}

	// LD A (DE) [] 1 [8]
	normal[0x1a] = []func(){nop, cpu.ldADE(mapper)}

	// DEC DE  [] 1 [8]
	normal[0x1b] = []func(){nop, cpu.decDE}

	// INC E  [Z 0 H -] 1 [4]
	normal[0x1c] = []func(){cpu.incE}

	// DEC E  [Z 1 H -] 1 [4]
	normal[0x1d] = []func(){cpu.decE}

	// LD E d8 [] 2 [8]
	normal[0x1e] = []func(){cpu.readParamA, cpu.ldEU}

	// RRA   [0 0 0 C] 1 [4]
	normal[0x1f] = []func(){cpu.rra}

	// JR NZ r8 [] 2 [12 8]
	normal[0x20] = []func(){nop, cpu.readParamA, cpu.jr}

	// LD HL d16 [] 3 [12]
	normal[0x21] = []func(){cpu.readParamA, cpu.readParamB, cpu.ldHLU16}

	// LD (HL+) A [] 1 [8]
	normal[0x22] = []func(){nop, cpu.ldHLIA(mapper)}

	// INC HL  [] 1 [8]
	normal[0x23] = []func(){nop, cpu.incHL}

	// INC H  [Z 0 H -] 1 [4]
	normal[0x24] = []func(){cpu.incH}

	// DEC H  [Z 1 H -] 1 [4]
	normal[0x25] = []func(){cpu.decH}

	// LD H d8 [] 2 [8]
	normal[0x26] = []func(){cpu.readParamA, cpu.ldHU}

	// DAA   [Z - 0 C] 1 [4]
	normal[0x27] = []func(){cpu.daa}

	// JR Z r8 [] 2 [12 8]
	normal[0x28] = []func(){nop, cpu.readParamA, cpu.jr}

	// ADD HL HL [- 0 H C] 1 [8]
	normal[0x29] = []func(){nop, cpu.addHLHL}

	// LD A (HL+) [] 1 [8]
	normal[0x2a] = []func(){nop, cpu.ldAHLI(mapper)}

	// DEC HL  [] 1 [8]
	normal[0x2b] = []func(){nop, cpu.decHL}

	// INC L  [Z 0 H -] 1 [4]
	normal[0x2c] = []func(){cpu.incL}

	// DEC L  [Z 1 H -] 1 [4]
	normal[0x2d] = []func(){cpu.decL}

	// LD L d8 [] 2 [8]
	normal[0x2e] = []func(){cpu.readParamA, cpu.ldLU}

	// CPL   [- 1 1 -] 1 [4]
	normal[0x2f] = []func(){cpu.cpl}

	// JR NC r8 [] 2 [12 8]
	normal[0x30] = []func(){nop, cpu.readParamA, cpu.jr}

	// LD SP d16 [] 3 [12]
	normal[0x31] = []func(){cpu.readParamA, cpu.readParamB, cpu.ldSPU16}

	// LD (HL-) A [] 1 [8]
	normal[0x32] = []func(){nop, cpu.ldHLDA(mapper)}

	// INC SP  [] 1 [8]
	normal[0x33] = []func(){nop, cpu.incSP}

	// INC (HL)  [Z 0 H -] 1 [12]
	normal[0x34] = []func(){nop, cpu.ldMHL(mapper), cpu.incM(mapper)}

	// DEC (HL)  [Z 1 H -] 1 [12]
	normal[0x35] = []func(){nop, cpu.ldMHL(mapper), cpu.decM(mapper)}

	// LD (HL) d8 [] 2 [12]
	normal[0x36] = []func(){cpu.readParamA, nop, cpu.ldHLU8(mapper)}

	// SCF   [- 0 0 1] 1 [4]
	normal[0x37] = []func(){cpu.scf}

	// JR C r8 [] 2 [12 8]
	normal[0x38] = []func(){nop, cpu.readParamA, cpu.jr}

	// ADD HL SP [- 0 H C] 1 [8]
	normal[0x39] = []func(){nop, cpu.addHLSP}

	// LD A (HL-) [] 1 [8]
	normal[0x3a] = []func(){nop, cpu.ldAHLD(mapper)}

	// DEC SP  [] 1 [8]
	normal[0x3b] = []func(){nop, cpu.decSP}

	// INC A  [Z 0 H -] 1 [4]
	normal[0x3c] = []func(){cpu.incA}

	// DEC A  [Z 1 H -] 1 [4]
	normal[0x3d] = []func(){cpu.decA}

	// LD A d8 [] 2 [8]
	normal[0x3e] = []func(){cpu.readParamA, cpu.ldAU}

	// CCF   [- 0 0 C] 1 [4]
	normal[0x3f] = []func(){cpu.ccf}

	// LD B B [] 1 [4]
	normal[0x40] = []func(){cpu.mooneye}

	// LD B C [] 1 [4]
	normal[0x41] = []func(){cpu.ldBC}

	// LD B D [] 1 [4]
	normal[0x42] = []func(){cpu.ldBD}

	// LD B E [] 1 [4]
	normal[0x43] = []func(){cpu.ldBE}

	// LD B H [] 1 [4]
	normal[0x44] = []func(){cpu.ldBH}

	// LD B L [] 1 [4]
	normal[0x45] = []func(){cpu.ldBL}

	// LD B (HL) [] 1 [8]
	normal[0x46] = []func(){nop, cpu.ldBHL(mapper)}

	// LD B A [] 1 [4]
	normal[0x47] = []func(){cpu.ldBA}

	// LD C B [] 1 [4]
	normal[0x48] = []func(){cpu.ldCB}

	// LD C C [] 1 [4]
	normal[0x49] = []func(){nop}

	// LD C D [] 1 [4]
	normal[0x4a] = []func(){cpu.ldCD}

	// LD C E [] 1 [4]
	normal[0x4b] = []func(){cpu.ldCE}

	// LD C H [] 1 [4]
	normal[0x4c] = []func(){cpu.ldCH}

	// LD C L [] 1 [4]
	normal[0x4d] = []func(){cpu.ldCL}

	// LD C (HL) [] 1 [8]
	normal[0x4e] = []func(){nop, cpu.ldCHL(mapper)}

	// LD C A [] 1 [4]
	normal[0x4f] = []func(){cpu.ldCA}

	// LD D B [] 1 [4]
	normal[0x50] = []func(){cpu.ldDB}

	// LD D C [] 1 [4]
	normal[0x51] = []func(){cpu.ldDC}

	// LD D D [] 1 [4]
	normal[0x52] = []func(){nop}

	// LD D E [] 1 [4]
	normal[0x53] = []func(){cpu.ldDE}

	// LD D H [] 1 [4]
	normal[0x54] = []func(){cpu.ldDH}

	// LD D L [] 1 [4]
	normal[0x55] = []func(){cpu.ldDL}

	// LD D (HL) [] 1 [8]
	normal[0x56] = []func(){nop, cpu.ldDHL(mapper)}

	// LD D A [] 1 [4]
	normal[0x57] = []func(){cpu.ldDA}

	// LD E B [] 1 [4]
	normal[0x58] = []func(){cpu.ldEB}

	// LD E C [] 1 [4]
	normal[0x59] = []func(){cpu.ldEC}

	// LD E D [] 1 [4]
	normal[0x5a] = []func(){cpu.ldED}

	// LD E E [] 1 [4]
	normal[0x5b] = []func(){nop}

	// LD E H [] 1 [4]
	normal[0x5c] = []func(){cpu.ldEH}

	// LD E L [] 1 [4]
	normal[0x5d] = []func(){cpu.ldEL}

	// LD E (HL) [] 1 [8]
	normal[0x5e] = []func(){nop, cpu.ldEHL(mapper)}

	// LD E A [] 1 [4]
	normal[0x5f] = []func(){cpu.ldEA}

	// LD H B [] 1 [4]
	normal[0x60] = []func(){cpu.ldHB}

	// LD H C [] 1 [4]
	normal[0x61] = []func(){cpu.ldHC}

	// LD H D [] 1 [4]
	normal[0x62] = []func(){cpu.ldHD}

	// LD H E [] 1 [4]
	normal[0x63] = []func(){cpu.ldHE}

	// LD H H [] 1 [4]
	normal[0x64] = []func(){nop}

	// LD H L [] 1 [4]
	normal[0x65] = []func(){cpu.ldHL}

	// LD H (HL) [] 1 [8]
	normal[0x66] = []func(){nop, cpu.ldHHL(mapper)}

	// LD H A [] 1 [4]
	normal[0x67] = []func(){cpu.ldHA}

	// LD L B [] 1 [4]
	normal[0x68] = []func(){cpu.ldLB}

	// LD L C [] 1 [4]
	normal[0x69] = []func(){cpu.ldLC}

	// LD L D [] 1 [4]
	normal[0x6a] = []func(){cpu.ldLD}

	// LD L E [] 1 [4]
	normal[0x6b] = []func(){cpu.ldLE}

	// LD L H [] 1 [4]
	normal[0x6c] = []func(){cpu.ldLH}

	// LD L L [] 1 [4]
	normal[0x6d] = []func(){nop}

	// LD L (HL) [] 1 [8]
	normal[0x6e] = []func(){nop, cpu.ldLHL(mapper)}

	// LD L A [] 1 [4]
	normal[0x6f] = []func(){cpu.ldLA}

	// LD (HL) B [] 1 [8]
	normal[0x70] = []func(){nop, cpu.ldHLB(mapper)}

	// LD (HL) C [] 1 [8]
	normal[0x71] = []func(){nop, cpu.ldHLC(mapper)}

	// LD (HL) D [] 1 [8]
	normal[0x72] = []func(){nop, cpu.ldHLD(mapper)}

	// LD (HL) E [] 1 [8]
	normal[0x73] = []func(){nop, cpu.ldHLE(mapper)}

	// LD (HL) H [] 1 [8]
	normal[0x74] = []func(){nop, cpu.ldHLH(mapper)}

	// LD (HL) L [] 1 [8]
	normal[0x75] = []func(){nop, cpu.ldHLL(mapper)}

	// HALT   [] 1 [4]
	normal[0x76] = []func(){cpu.halt()}

	// LD (HL) A [] 1 [8]
	normal[0x77] = []func(){nop, cpu.ldHLA(mapper)}

	// LD A B [] 1 [4]
	normal[0x78] = []func(){cpu.ldAB}

	// LD A C [] 1 [4]
	normal[0x79] = []func(){cpu.ldAC}

	// LD A D [] 1 [4]
	normal[0x7a] = []func(){cpu.ldAD}

	// LD A E [] 1 [4]
	normal[0x7b] = []func(){cpu.ldAE}

	// LD A H [] 1 [4]
	normal[0x7c] = []func(){cpu.ldAH}

	// LD A L [] 1 [4]
	normal[0x7d] = []func(){cpu.ldAL}

	// LD A (HL) [] 1 [8]
	normal[0x7e] = []func(){nop, cpu.ldAHL(mapper)}

	// LD A A [] 1 [4]
	normal[0x7f] = []func(){nop}

	// ADD A B [Z 0 H C] 1 [4]
	normal[0x80] = []func(){cpu.addB}

	// ADD A C [Z 0 H C] 1 [4]
	normal[0x81] = []func(){cpu.addC}

	// ADD A D [Z 0 H C] 1 [4]
	normal[0x82] = []func(){cpu.addD}

	// ADD A E [Z 0 H C] 1 [4]
	normal[0x83] = []func(){cpu.addE}

	// ADD A H [Z 0 H C] 1 [4]
	normal[0x84] = []func(){cpu.addH}

	// ADD A L [Z 0 H C] 1 [4]
	normal[0x85] = []func(){cpu.addL}

	// ADD A (HL) [Z 0 H C] 1 [8]
	normal[0x86] = []func(){nop, cpu.addM(mapper)}

	// ADD A A [Z 0 H C] 1 [4]
	normal[0x87] = []func(){cpu.addA}

	// ADC A B [Z 0 H C] 1 [4]
	normal[0x88] = []func(){cpu.adcB}

	// ADC A C [Z 0 H C] 1 [4]
	normal[0x89] = []func(){cpu.adcC}

	// ADC A D [Z 0 H C] 1 [4]
	normal[0x8a] = []func(){cpu.adcD}

	// ADC A E [Z 0 H C] 1 [4]
	normal[0x8b] = []func(){cpu.adcE}

	// ADC A H [Z 0 H C] 1 [4]
	normal[0x8c] = []func(){cpu.adcH}

	// ADC A L [Z 0 H C] 1 [4]
	normal[0x8d] = []func(){cpu.adcL}

	// ADC A (HL) [Z 0 H C] 1 [8]
	normal[0x8e] = []func(){nop, cpu.adcM(mapper)}

	// ADC A A [Z 0 H C] 1 [4]
	normal[0x8f] = []func(){cpu.adcA}

	// SUB B  [Z 1 H C] 1 [4]
	normal[0x90] = []func(){cpu.subB}

	// SUB C  [Z 1 H C] 1 [4]
	normal[0x91] = []func(){cpu.subC}

	// SUB D  [Z 1 H C] 1 [4]
	normal[0x92] = []func(){cpu.subD}

	// SUB E  [Z 1 H C] 1 [4]
	normal[0x93] = []func(){cpu.subE}

	// SUB H  [Z 1 H C] 1 [4]
	normal[0x94] = []func(){cpu.subH}

	// SUB L  [Z 1 H C] 1 [4]
	normal[0x95] = []func(){cpu.subL}

	// SUB (HL)  [Z 1 H C] 1 [8]
	normal[0x96] = []func(){nop, cpu.subM(mapper)}

	// SUB A  [Z 1 H C] 1 [4]
	normal[0x97] = []func(){cpu.subA}

	// SBC A B [Z 1 H C] 1 [4]
	normal[0x98] = []func(){cpu.sbcB}

	// SBC A C [Z 1 H C] 1 [4]
	normal[0x99] = []func(){cpu.sbcC}

	// SBC A D [Z 1 H C] 1 [4]
	normal[0x9a] = []func(){cpu.sbcD}

	// SBC A E [Z 1 H C] 1 [4]
	normal[0x9b] = []func(){cpu.sbcE}

	// SBC A H [Z 1 H C] 1 [4]
	normal[0x9c] = []func(){cpu.sbcH}

	// SBC A L [Z 1 H C] 1 [4]
	normal[0x9d] = []func(){cpu.sbcL}

	// SBC A (HL) [Z 1 H C] 1 [8]
	normal[0x9e] = []func(){nop, cpu.sbcM(mapper)}

	// SBC A A [Z 1 H C] 1 [4]
	normal[0x9f] = []func(){cpu.sbcA}

	// AND B  [Z 0 1 0] 1 [4]
	normal[0xa0] = []func(){cpu.andB}

	// AND C  [Z 0 1 0] 1 [4]
	normal[0xa1] = []func(){cpu.andC}

	// AND D  [Z 0 1 0] 1 [4]
	normal[0xa2] = []func(){cpu.andD}

	// AND E  [Z 0 1 0] 1 [4]
	normal[0xa3] = []func(){cpu.andE}

	// AND H  [Z 0 1 0] 1 [4]
	normal[0xa4] = []func(){cpu.andH}

	// AND L  [Z 0 1 0] 1 [4]
	normal[0xa5] = []func(){cpu.andL}

	// AND (HL)  [Z 0 1 0] 1 [8]
	normal[0xa6] = []func(){nop, cpu.andM(mapper)}

	// AND A  [Z 0 1 0] 1 [4]
	normal[0xa7] = []func(){cpu.andA}

	// XOR B  [Z 0 0 0] 1 [4]
	normal[0xa8] = []func(){cpu.xorB}

	// XOR C  [Z 0 0 0] 1 [4]
	normal[0xa9] = []func(){cpu.xorC}

	// XOR D  [Z 0 0 0] 1 [4]
	normal[0xaa] = []func(){cpu.xorD}

	// XOR E  [Z 0 0 0] 1 [4]
	normal[0xab] = []func(){cpu.xorE}

	// XOR H  [Z 0 0 0] 1 [4]
	normal[0xac] = []func(){cpu.xorH}

	// XOR L  [Z 0 0 0] 1 [4]
	normal[0xad] = []func(){cpu.xorL}

	// XOR (HL)  [Z 0 0 0] 1 [8]
	normal[0xae] = []func(){nop, cpu.xorM(mapper)}

	// XOR A  [Z 0 0 0] 1 [4]
	normal[0xaf] = []func(){cpu.xorA}

	// OR B  [Z 0 0 0] 1 [4]
	normal[0xb0] = []func(){cpu.orB}

	// OR C  [Z 0 0 0] 1 [4]
	normal[0xb1] = []func(){cpu.orC}

	// OR D  [Z 0 0 0] 1 [4]
	normal[0xb2] = []func(){cpu.orD}

	// OR E  [Z 0 0 0] 1 [4]
	normal[0xb3] = []func(){cpu.orE}

	// OR H  [Z 0 0 0] 1 [4]
	normal[0xb4] = []func(){cpu.orH}

	// OR L  [Z 0 0 0] 1 [4]
	normal[0xb5] = []func(){cpu.orL}

	// OR (HL)  [Z 0 0 0] 1 [8]
	normal[0xb6] = []func(){nop, cpu.orM(mapper)}

	// OR A  [Z 0 0 0] 1 [4]
	normal[0xb7] = []func(){cpu.orA}

	// CP B  [Z 1 H C] 1 [4]
	normal[0xb8] = []func(){cpu.cpB}

	// CP C  [Z 1 H C] 1 [4]
	normal[0xb9] = []func(){cpu.cpC}

	// CP D  [Z 1 H C] 1 [4]
	normal[0xba] = []func(){cpu.cpD}

	// CP E  [Z 1 H C] 1 [4]
	normal[0xbb] = []func(){cpu.cpE}

	// CP H  [Z 1 H C] 1 [4]
	normal[0xbc] = []func(){cpu.cpH}

	// CP L  [Z 1 H C] 1 [4]
	normal[0xbd] = []func(){cpu.cpL}

	// CP (HL)  [Z 1 H C] 1 [8]
	normal[0xbe] = []func(){nop, cpu.cpM(mapper)}

	// CP A  [Z 1 H C] 1 [4]
	normal[0xbf] = []func(){cpu.cpA}

	// RET NZ  [] 1 [20 8]
	normal[0xc0] = []func(){nop, nop, cpu.pop(mapper, &cpu.m8a), cpu.pop(mapper, &cpu.m8b), cpu.ret}

	// POP BC  [] 1 [12]
	normal[0xc1] = []func(){nop, cpu.pop(mapper, &cpu.c), cpu.pop(mapper, &cpu.b)}

	// JP NZ a16 [] 3 [16 12]
	normal[0xc2] = []func(){nop, cpu.readParamA, cpu.readParamB, cpu.jp}

	// JP a16  [] 3 [16]
	normal[0xc3] = []func(){nop, cpu.readParamA, cpu.readParamB, cpu.jp}

	// CALL NZ a16 [] 3 [24 12]
	normal[0xc4] = []func(){nop, cpu.readParamA, cpu.readParamB, cpu.call, cpu.push(mapper, &cpu.m8b), cpu.push(mapper, &cpu.m8a)}

	// PUSH BC      1 [16]
	normal[0xc5] = []func(){nop, nop, cpu.push(mapper, &cpu.b), cpu.push(mapper, &cpu.c)}

	// ADD A d8 [Z 0 H C] 2 [8]
	normal[0xc6] = []func(){cpu.readParamA, cpu.addU}

	// RST 00H  [] 1 [16]
	normal[0xc7] = []func(){nop, cpu.rst(0x0000), cpu.push(mapper, &cpu.m8b), cpu.push(mapper, &cpu.m8a)}

	// RET Z  [] 1 [20 8]
	normal[0xc8] = []func(){nop, nop, cpu.pop(mapper, &cpu.m8a), cpu.pop(mapper, &cpu.m8b), cpu.ret}

	// RET   [] 1 [16]
	normal[0xc9] = []func(){nop, cpu.pop(mapper, &cpu.m8a), cpu.pop(mapper, &cpu.m8b), cpu.ret}

	// JP Z a16 [] 3 [16 12]
	normal[0xca] = []func(){nop, cpu.readParamA, cpu.readParamB, cpu.jp}

	// CALL Z a16 [] 3 [24 12]
	normal[0xcc] = []func(){nop, cpu.readParamA, cpu.readParamB, cpu.call, cpu.push(mapper, &cpu.m8b), cpu.push(mapper, &cpu.m8a)}

	// CALL a16  [] 3 [24]
	normal[0xcd] = []func(){nop, cpu.readParamA, cpu.readParamB, cpu.call, cpu.push(mapper, &cpu.m8b), cpu.push(mapper, &cpu.m8a)}

	// ADC A d8 [Z 0 H C] 2 [8]
	normal[0xce] = []func(){cpu.readParamA, cpu.adcU}

	// RST 08H  [] 1 [16]
	normal[0xcf] = []func(){nop, cpu.rst(0x0008), cpu.push(mapper, &cpu.m8b), cpu.push(mapper, &cpu.m8a)}

	// RET NC  [] 1 [20 8]
	normal[0xd0] = []func(){nop, nop, cpu.pop(mapper, &cpu.m8a), cpu.pop(mapper, &cpu.m8b), cpu.ret}

	// POP DE  [] 1 [12]
	normal[0xd1] = []func(){nop, cpu.pop(mapper, &cpu.e), cpu.pop(mapper, &cpu.d)}

	// JP NC a16 [] 3 [16 12]
	normal[0xd2] = []func(){nop, cpu.readParamA, cpu.readParamB, cpu.jp}

	// CALL NC a16 [] 3 [24 12]
	normal[0xd4] = []func(){nop, cpu.readParamA, cpu.readParamB, cpu.call, cpu.push(mapper, &cpu.m8b), cpu.push(mapper, &cpu.m8a)}

	// PUSH DE      1 [16]
	normal[0xd5] = []func(){nop, nop, cpu.push(mapper, &cpu.d), cpu.push(mapper, &cpu.e)}

	// SUB d8  [Z 1 H C] 2 [8]
	normal[0xd6] = []func(){cpu.readParamA, cpu.subU}

	// RST 10H  [] 1 [16]
	normal[0xd7] = []func(){nop, cpu.rst(0x0010), cpu.push(mapper, &cpu.m8b), cpu.push(mapper, &cpu.m8a)}

	// RET C  [] 1 [20 8]
	normal[0xd8] = []func(){nop, nop, cpu.pop(mapper, &cpu.m8a), cpu.pop(mapper, &cpu.m8b), cpu.ret}

	// RETI   [] 1 [16]
	normal[0xd9] = []func(){nop, cpu.pop(mapper, &cpu.m8a), cpu.pop(mapper, &cpu.m8b), cpu.reti}

	// JP C a16 [] 3 [16 12]
	normal[0xda] = []func(){nop, cpu.readParamA, cpu.readParamB, cpu.jp}

	// CALL C a16 [] 3 [24 12]
	normal[0xdc] = []func(){nop, cpu.readParamA, cpu.readParamB, cpu.call, cpu.push(mapper, &cpu.m8b), cpu.push(mapper, &cpu.m8a)}

	// SBC A d8 [Z 1 H C] 2 [8]
	normal[0xde] = []func(){cpu.readParamA, cpu.sbcU}

	// RST 18H  [] 1 [16]
	normal[0xdf] = []func(){nop, cpu.rst(0x0018), cpu.push(mapper, &cpu.m8b), cpu.push(mapper, &cpu.m8a)}

	// LDH (a8) A   2 [12]
	normal[0xe0] = []func(){cpu.readParamA, nop, cpu.ldUXA(mapper)}

	// POP HL  [] 1 [12]
	normal[0xe1] = []func(){nop, cpu.pop(mapper, &cpu.l), cpu.pop(mapper, &cpu.h)}

	// LD (C) A     1 [8]
	normal[0xe2] = []func(){nop, cpu.ldCXA(mapper)}

	// PUSH HL      1 [16]
	normal[0xe5] = []func(){nop, nop, cpu.push(mapper, &cpu.h), cpu.push(mapper, &cpu.l)}

	// AND d8  [Z 0 1 0] 2 [8]
	normal[0xe6] = []func(){cpu.readParamA, cpu.andU}

	// RST 20H  [] 1 [16]
	normal[0xe7] = []func(){nop, cpu.rst(0x0020), cpu.push(mapper, &cpu.m8b), cpu.push(mapper, &cpu.m8a)}

	// ADD SP r8 [0 0 H C] 2 [16]
	normal[0xe8] = []func(){nop, cpu.readParamA, cpu.addSP, nop}

	// JP (HL)  [] 1 [4]
	normal[0xe9] = []func(){cpu.jpHL}

	// LD (a16) A [] 3 [16]
	normal[0xea] = []func(){cpu.readParamA, cpu.readParamB, nop, cpu.ldUX16A(mapper)}

	// XOR d8  [Z 0 0 0] 2 [8]
	normal[0xee] = []func(){cpu.readParamA, cpu.xorU}

	// RST 28H  [] 1 [16]
	normal[0xef] = []func(){nop, cpu.rst(0x0028), cpu.push(mapper, &cpu.m8b), cpu.push(mapper, &cpu.m8a)}

	// LDH A (a8)   2 [12]
	normal[0xf0] = []func(){cpu.readParamA, nop, cpu.ldAUX(mapper)}

	// POP AF  [Z N H C] 1 [12]
	normal[0xf1] = []func(){nop, cpu.popF(mapper), cpu.pop(mapper, &cpu.a)}

	// LD A (C)     1 [8]
	normal[0xf2] = []func(){nop, cpu.ldACX(mapper)}

	// DI   [] 1 [4]
	normal[0xf3] = []func(){cpu.di}

	// PUSH AF      1 [16]
	normal[0xf5] = []func(){nop, nop, cpu.push(mapper, &cpu.a), cpu.push(mapper, &cpu.f)}

	// OR d8  [Z 0 0 0] 2 [8]
	normal[0xf6] = []func(){cpu.readParamA, cpu.orU}

	// RST 30H  [] 1 [16]
	normal[0xf7] = []func(){nop, cpu.rst(0x0030), cpu.push(mapper, &cpu.m8b), cpu.push(mapper, &cpu.m8a)}

	// LD HL SP+r8 [0 0 H C] 2 [12]
	normal[0xf8] = []func(){cpu.readParamA, nop, cpu.ldHLSP}

	// LD SP HL [] 1 [8]
	normal[0xf9] = []func(){nop, cpu.ldSPHL}

	// LD A (a16) [] 3 [16]
	normal[0xfa] = []func(){cpu.readParamA, cpu.readParamB, nop, cpu.ldAUX16(mapper)}

	// EI   [] 1 [4]
	normal[0xfb] = []func(){cpu.ei}

	// CP d8  [Z 1 H C] 2 [8]
	normal[0xfe] = []func(){cpu.readParamA, cpu.cpU}

	// RST 38H  [] 1 [16]
	normal[0xff] = []func(){nop, cpu.rst(0x0038), cpu.push(mapper, &cpu.m8b), cpu.push(mapper, &cpu.m8a)}

	// RLC B  [Z 0 0 C] 2 [8]
	prefix[0x00] = []func(){nop, cpu.rlcB}

	// RLC C  [Z 0 0 C] 2 [8]
	prefix[0x01] = []func(){nop, cpu.rlcC}

	// RLC D  [Z 0 0 C] 2 [8]
	prefix[0x02] = []func(){nop, cpu.rlcD}

	// RLC E  [Z 0 0 C] 2 [8]
	prefix[0x03] = []func(){nop, cpu.rlcE}

	// RLC H  [Z 0 0 C] 2 [8]
	prefix[0x04] = []func(){nop, cpu.rlcH}

	// RLC L  [Z 0 0 C] 2 [8]
	prefix[0x05] = []func(){nop, cpu.rlcL}

	// RLC (HL)  [Z 0 0 C] 2 [16]
	prefix[0x06] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.rlcM(mapper)}

	// RLC A  [Z 0 0 C] 2 [8]
	prefix[0x07] = []func(){nop, cpu.rlcA}

	// RRC B  [Z 0 0 C] 2 [8]
	prefix[0x08] = []func(){nop, cpu.rrcB}

	// RRC C  [Z 0 0 C] 2 [8]
	prefix[0x09] = []func(){nop, cpu.rrcC}

	// RRC D  [Z 0 0 C] 2 [8]
	prefix[0x0a] = []func(){nop, cpu.rrcD}

	// RRC E  [Z 0 0 C] 2 [8]
	prefix[0x0b] = []func(){nop, cpu.rrcE}

	// RRC H  [Z 0 0 C] 2 [8]
	prefix[0x0c] = []func(){nop, cpu.rrcH}

	// RRC L  [Z 0 0 C] 2 [8]
	prefix[0x0d] = []func(){nop, cpu.rrcL}

	// RRC (HL)  [Z 0 0 C] 2 [16]
	prefix[0x0e] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.rrcM(mapper)}

	// RRC A  [Z 0 0 C] 2 [8]
	prefix[0x0f] = []func(){nop, cpu.rrcA}

	// RL B  [Z 0 0 C] 2 [8]
	prefix[0x10] = []func(){nop, cpu.rlB}

	// RL C  [Z 0 0 C] 2 [8]
	prefix[0x11] = []func(){nop, cpu.rlC}

	// RL D  [Z 0 0 C] 2 [8]
	prefix[0x12] = []func(){nop, cpu.rlD}

	// RL E  [Z 0 0 C] 2 [8]
	prefix[0x13] = []func(){nop, cpu.rlE}

	// RL H  [Z 0 0 C] 2 [8]
	prefix[0x14] = []func(){nop, cpu.rlH}

	// RL L  [Z 0 0 C] 2 [8]
	prefix[0x15] = []func(){nop, cpu.rlL}

	// RL (HL)  [Z 0 0 C] 2 [16]
	prefix[0x16] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.rlM(mapper)}

	// RL A  [Z 0 0 C] 2 [8]
	prefix[0x17] = []func(){nop, cpu.rlA}

	// RR B  [Z 0 0 C] 2 [8]
	prefix[0x18] = []func(){nop, cpu.rrB}

	// RR C  [Z 0 0 C] 2 [8]
	prefix[0x19] = []func(){nop, cpu.rrC}

	// RR D  [Z 0 0 C] 2 [8]
	prefix[0x1a] = []func(){nop, cpu.rrD}

	// RR E  [Z 0 0 C] 2 [8]
	prefix[0x1b] = []func(){nop, cpu.rrE}

	// RR H  [Z 0 0 C] 2 [8]
	prefix[0x1c] = []func(){nop, cpu.rrH}

	// RR L  [Z 0 0 C] 2 [8]
	prefix[0x1d] = []func(){nop, cpu.rrL}

	// RR (HL)  [Z 0 0 C] 2 [16]
	prefix[0x1e] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.rrM(mapper)}

	// RR A  [Z 0 0 C] 2 [8]
	prefix[0x1f] = []func(){nop, cpu.rrA}

	// SLA B  [Z 0 0 C] 2 [8]
	prefix[0x20] = []func(){nop, cpu.slaB}

	// SLA C  [Z 0 0 C] 2 [8]
	prefix[0x21] = []func(){nop, cpu.slaC}

	// SLA D  [Z 0 0 C] 2 [8]
	prefix[0x22] = []func(){nop, cpu.slaD}

	// SLA E  [Z 0 0 C] 2 [8]
	prefix[0x23] = []func(){nop, cpu.slaE}

	// SLA H  [Z 0 0 C] 2 [8]
	prefix[0x24] = []func(){nop, cpu.slaH}

	// SLA L  [Z 0 0 C] 2 [8]
	prefix[0x25] = []func(){nop, cpu.slaL}

	// SLA (HL)  [Z 0 0 C] 2 [16]
	prefix[0x26] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.slaM(mapper)}

	// SLA A  [Z 0 0 C] 2 [8]
	prefix[0x27] = []func(){nop, cpu.slaA}

	// SRA B  [Z 0 0 C] 2 [8]
	prefix[0x28] = []func(){nop, cpu.sraB}

	// SRA C  [Z 0 0 C] 2 [8]
	prefix[0x29] = []func(){nop, cpu.sraC}

	// SRA D  [Z 0 0 C] 2 [8]
	prefix[0x2a] = []func(){nop, cpu.sraD}

	// SRA E  [Z 0 0 C] 2 [8]
	prefix[0x2b] = []func(){nop, cpu.sraE}

	// SRA H  [Z 0 0 C] 2 [8]
	prefix[0x2c] = []func(){nop, cpu.sraH}

	// SRA L  [Z 0 0 C] 2 [8]
	prefix[0x2d] = []func(){nop, cpu.sraL}

	// SRA (HL)  [Z 0 0 C] 2 [16]
	prefix[0x2e] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.sraM(mapper)}

	// SRA A  [Z 0 0 C] 2 [8]
	prefix[0x2f] = []func(){nop, cpu.sraA}

	// SWAP B  [Z 0 0 0] 2 [8]
	prefix[0x30] = []func(){nop, cpu.swapB}

	// SWAP C  [Z 0 0 0] 2 [8]
	prefix[0x31] = []func(){nop, cpu.swapC}

	// SWAP D  [Z 0 0 0] 2 [8]
	prefix[0x32] = []func(){nop, cpu.swapD}

	// SWAP E  [Z 0 0 0] 2 [8]
	prefix[0x33] = []func(){nop, cpu.swapE}

	// SWAP H  [Z 0 0 0] 2 [8]
	prefix[0x34] = []func(){nop, cpu.swapH}

	// SWAP L  [Z 0 0 0] 2 [8]
	prefix[0x35] = []func(){nop, cpu.swapL}

	// SWAP (HL)  [Z 0 0 0] 2 [16]
	prefix[0x36] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.swapM(mapper)}

	// SWAP A  [Z 0 0 0] 2 [8]
	prefix[0x37] = []func(){nop, cpu.swapA}

	// SRL B  [Z 0 0 C] 2 [8]
	prefix[0x38] = []func(){nop, cpu.srlB}

	// SRL C  [Z 0 0 C] 2 [8]
	prefix[0x39] = []func(){nop, cpu.srlC}

	// SRL D  [Z 0 0 C] 2 [8]
	prefix[0x3a] = []func(){nop, cpu.srlD}

	// SRL E  [Z 0 0 C] 2 [8]
	prefix[0x3b] = []func(){nop, cpu.srlE}

	// SRL H  [Z 0 0 C] 2 [8]
	prefix[0x3c] = []func(){nop, cpu.srlH}

	// SRL L  [Z 0 0 C] 2 [8]
	prefix[0x3d] = []func(){nop, cpu.srlL}

	// SRL (HL)  [Z 0 0 C] 2 [16]
	prefix[0x3e] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.srlM(mapper)}

	// SRL A  [Z 0 0 C] 2 [8]
	prefix[0x3f] = []func(){nop, cpu.srlA}

	// BIT 0 B [Z 0 1 -] 2 [8]
	prefix[0x40] = []func(){nop, cpu.bit(0, &cpu.b)}

	// BIT 0 C [Z 0 1 -] 2 [8]
	prefix[0x41] = []func(){nop, cpu.bit(0, &cpu.c)}

	// BIT 0 D [Z 0 1 -] 2 [8]
	prefix[0x42] = []func(){nop, cpu.bit(0, &cpu.d)}

	// BIT 0 E [Z 0 1 -] 2 [8]
	prefix[0x43] = []func(){nop, cpu.bit(0, &cpu.e)}

	// BIT 0 H [Z 0 1 -] 2 [8]
	prefix[0x44] = []func(){nop, cpu.bit(0, &cpu.h)}

	// BIT 0 L [Z 0 1 -] 2 [8]
	prefix[0x45] = []func(){nop, cpu.bit(0, &cpu.l)}

	// BIT 0 (HL) [Z 0 1 -] 2 [12]
	prefix[0x46] = []func(){nop, nop, cpu.bitM(mapper, 0)}

	// BIT 0 A [Z 0 1 -] 2 [8]
	prefix[0x47] = []func(){nop, cpu.bit(0, &cpu.a)}

	// BIT 1 B [Z 0 1 -] 2 [8]
	prefix[0x48] = []func(){nop, cpu.bit(1, &cpu.b)}

	// BIT 1 C [Z 0 1 -] 2 [8]
	prefix[0x49] = []func(){nop, cpu.bit(1, &cpu.c)}

	// BIT 1 D [Z 0 1 -] 2 [8]
	prefix[0x4a] = []func(){nop, cpu.bit(1, &cpu.d)}

	// BIT 1 E [Z 0 1 -] 2 [8]
	prefix[0x4b] = []func(){nop, cpu.bit(1, &cpu.e)}

	// BIT 1 H [Z 0 1 -] 2 [8]
	prefix[0x4c] = []func(){nop, cpu.bit(1, &cpu.h)}

	// BIT 1 L [Z 0 1 -] 2 [8]
	prefix[0x4d] = []func(){nop, cpu.bit(1, &cpu.l)}

	// BIT 1 (HL) [Z 0 1 -] 2 [12]
	prefix[0x4e] = []func(){nop, nop, cpu.bitM(mapper, 1)}

	// BIT 1 A [Z 0 1 -] 2 [8]
	prefix[0x4f] = []func(){nop, cpu.bit(1, &cpu.a)}

	// BIT 2 B [Z 0 1 -] 2 [8]
	prefix[0x50] = []func(){nop, cpu.bit(2, &cpu.b)}

	// BIT 2 C [Z 0 1 -] 2 [8]
	prefix[0x51] = []func(){nop, cpu.bit(2, &cpu.c)}

	// BIT 2 D [Z 0 1 -] 2 [8]
	prefix[0x52] = []func(){nop, cpu.bit(2, &cpu.d)}

	// BIT 2 E [Z 0 1 -] 2 [8]
	prefix[0x53] = []func(){nop, cpu.bit(2, &cpu.e)}

	// BIT 2 H [Z 0 1 -] 2 [8]
	prefix[0x54] = []func(){nop, cpu.bit(2, &cpu.h)}

	// BIT 2 L [Z 0 1 -] 2 [8]
	prefix[0x55] = []func(){nop, cpu.bit(2, &cpu.l)}

	// BIT 2 (HL) [Z 0 1 -] 2 [12]
	prefix[0x56] = []func(){nop, nop, cpu.bitM(mapper, 2)}

	// BIT 2 A [Z 0 1 -] 2 [8]
	prefix[0x57] = []func(){nop, cpu.bit(2, &cpu.a)}

	// BIT 3 B [Z 0 1 -] 2 [8]
	prefix[0x58] = []func(){nop, cpu.bit(3, &cpu.b)}

	// BIT 3 C [Z 0 1 -] 2 [8]
	prefix[0x59] = []func(){nop, cpu.bit(3, &cpu.c)}

	// BIT 3 D [Z 0 1 -] 2 [8]
	prefix[0x5a] = []func(){nop, cpu.bit(3, &cpu.d)}

	// BIT 3 E [Z 0 1 -] 2 [8]
	prefix[0x5b] = []func(){nop, cpu.bit(3, &cpu.e)}

	// BIT 3 H [Z 0 1 -] 2 [8]
	prefix[0x5c] = []func(){nop, cpu.bit(3, &cpu.h)}

	// BIT 3 L [Z 0 1 -] 2 [8]
	prefix[0x5d] = []func(){nop, cpu.bit(3, &cpu.l)}

	// BIT 3 (HL) [Z 0 1 -] 2 [12]
	prefix[0x5e] = []func(){nop, nop, cpu.bitM(mapper, 3)}

	// BIT 3 A [Z 0 1 -] 2 [8]
	prefix[0x5f] = []func(){nop, cpu.bit(3, &cpu.a)}

	// BIT 4 B [Z 0 1 -] 2 [8]
	prefix[0x60] = []func(){nop, cpu.bit(4, &cpu.b)}

	// BIT 4 C [Z 0 1 -] 2 [8]
	prefix[0x61] = []func(){nop, cpu.bit(4, &cpu.c)}

	// BIT 4 D [Z 0 1 -] 2 [8]
	prefix[0x62] = []func(){nop, cpu.bit(4, &cpu.d)}

	// BIT 4 E [Z 0 1 -] 2 [8]
	prefix[0x63] = []func(){nop, cpu.bit(4, &cpu.e)}

	// BIT 4 H [Z 0 1 -] 2 [8]
	prefix[0x64] = []func(){nop, cpu.bit(4, &cpu.h)}

	// BIT 4 L [Z 0 1 -] 2 [8]
	prefix[0x65] = []func(){nop, cpu.bit(4, &cpu.l)}

	// BIT 4 (HL) [Z 0 1 -] 2 [12]
	prefix[0x66] = []func(){nop, nop, cpu.bitM(mapper, 4)}

	// BIT 4 A [Z 0 1 -] 2 [8]
	prefix[0x67] = []func(){nop, cpu.bit(4, &cpu.a)}

	// BIT 5 B [Z 0 1 -] 2 [8]
	prefix[0x68] = []func(){nop, cpu.bit(5, &cpu.b)}

	// BIT 5 C [Z 0 1 -] 2 [8]
	prefix[0x69] = []func(){nop, cpu.bit(5, &cpu.c)}

	// BIT 5 D [Z 0 1 -] 2 [8]
	prefix[0x6a] = []func(){nop, cpu.bit(5, &cpu.d)}

	// BIT 5 E [Z 0 1 -] 2 [8]
	prefix[0x6b] = []func(){nop, cpu.bit(5, &cpu.e)}

	// BIT 5 H [Z 0 1 -] 2 [8]
	prefix[0x6c] = []func(){nop, cpu.bit(5, &cpu.h)}

	// BIT 5 L [Z 0 1 -] 2 [8]
	prefix[0x6d] = []func(){nop, cpu.bit(5, &cpu.l)}

	// BIT 5 (HL) [Z 0 1 -] 2 [12]
	prefix[0x6e] = []func(){nop, nop, cpu.bitM(mapper, 5)}

	// BIT 5 A [Z 0 1 -] 2 [8]
	prefix[0x6f] = []func(){nop, cpu.bit(5, &cpu.a)}

	// BIT 6 B [Z 0 1 -] 2 [8]
	prefix[0x70] = []func(){nop, cpu.bit(6, &cpu.b)}

	// BIT 6 C [Z 0 1 -] 2 [8]
	prefix[0x71] = []func(){nop, cpu.bit(6, &cpu.c)}

	// BIT 6 D [Z 0 1 -] 2 [8]
	prefix[0x72] = []func(){nop, cpu.bit(6, &cpu.d)}

	// BIT 6 E [Z 0 1 -] 2 [8]
	prefix[0x73] = []func(){nop, cpu.bit(6, &cpu.e)}

	// BIT 6 H [Z 0 1 -] 2 [8]
	prefix[0x74] = []func(){nop, cpu.bit(6, &cpu.h)}

	// BIT 6 L [Z 0 1 -] 2 [8]
	prefix[0x75] = []func(){nop, cpu.bit(6, &cpu.l)}

	// BIT 6 (HL) [Z 0 1 -] 2 [12]
	prefix[0x76] = []func(){nop, nop, cpu.bitM(mapper, 6)}

	// BIT 6 A [Z 0 1 -] 2 [8]
	prefix[0x77] = []func(){nop, cpu.bit(6, &cpu.a)}

	// BIT 7 B [Z 0 1 -] 2 [8]
	prefix[0x78] = []func(){nop, cpu.bit(7, &cpu.b)}

	// BIT 7 C [Z 0 1 -] 2 [8]
	prefix[0x79] = []func(){nop, cpu.bit(7, &cpu.c)}

	// BIT 7 D [Z 0 1 -] 2 [8]
	prefix[0x7a] = []func(){nop, cpu.bit(7, &cpu.d)}

	// BIT 7 E [Z 0 1 -] 2 [8]
	prefix[0x7b] = []func(){nop, cpu.bit(7, &cpu.e)}

	// BIT 7 H [Z 0 1 -] 2 [8]
	prefix[0x7c] = []func(){nop, cpu.bit(7, &cpu.h)}

	// BIT 7 L [Z 0 1 -] 2 [8]
	prefix[0x7d] = []func(){nop, cpu.bit(7, &cpu.l)}

	// BIT 7 (HL) [Z 0 1 -] 2 [12]
	prefix[0x7e] = []func(){nop, nop, cpu.bitM(mapper, 7)}

	// BIT 7 A [Z 0 1 -] 2 [8]
	prefix[0x7f] = []func(){nop, cpu.bit(7, &cpu.a)}

	// RES 0 B [] 2 [8]
	prefix[0x80] = []func(){nop, cpu.res(0, &cpu.b)}

	// RES 0 C [] 2 [8]
	prefix[0x81] = []func(){nop, cpu.res(0, &cpu.c)}

	// RES 0 D [] 2 [8]
	prefix[0x82] = []func(){nop, cpu.res(0, &cpu.d)}

	// RES 0 E [] 2 [8]
	prefix[0x83] = []func(){nop, cpu.res(0, &cpu.e)}

	// RES 0 H [] 2 [8]
	prefix[0x84] = []func(){nop, cpu.res(0, &cpu.h)}

	// RES 0 L [] 2 [8]
	prefix[0x85] = []func(){nop, cpu.res(0, &cpu.l)}

	// RES 0 (HL) [] 2 [16]
	prefix[0x86] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.resM(mapper, 0)}

	// RES 0 A [] 2 [8]
	prefix[0x87] = []func(){nop, cpu.res(0, &cpu.a)}

	// RES 1 B [] 2 [8]
	prefix[0x88] = []func(){nop, cpu.res(1, &cpu.b)}

	// RES 1 C [] 2 [8]
	prefix[0x89] = []func(){nop, cpu.res(1, &cpu.c)}

	// RES 1 D [] 2 [8]
	prefix[0x8a] = []func(){nop, cpu.res(1, &cpu.d)}

	// RES 1 E [] 2 [8]
	prefix[0x8b] = []func(){nop, cpu.res(1, &cpu.e)}

	// RES 1 H [] 2 [8]
	prefix[0x8c] = []func(){nop, cpu.res(1, &cpu.h)}

	// RES 1 L [] 2 [8]
	prefix[0x8d] = []func(){nop, cpu.res(1, &cpu.l)}

	// RES 1 (HL) [] 2 [16]
	prefix[0x8e] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.resM(mapper, 1)}

	// RES 1 A [] 2 [8]
	prefix[0x8f] = []func(){nop, cpu.res(1, &cpu.a)}

	// RES 2 B [] 2 [8]
	prefix[0x90] = []func(){nop, cpu.res(2, &cpu.b)}

	// RES 2 C [] 2 [8]
	prefix[0x91] = []func(){nop, cpu.res(2, &cpu.c)}

	// RES 2 D [] 2 [8]
	prefix[0x92] = []func(){nop, cpu.res(2, &cpu.d)}

	// RES 2 E [] 2 [8]
	prefix[0x93] = []func(){nop, cpu.res(2, &cpu.e)}

	// RES 2 H [] 2 [8]
	prefix[0x94] = []func(){nop, cpu.res(2, &cpu.h)}

	// RES 2 L [] 2 [8]
	prefix[0x95] = []func(){nop, cpu.res(2, &cpu.l)}

	// RES 2 (HL) [] 2 [16]
	prefix[0x96] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.resM(mapper, 2)}

	// RES 2 A [] 2 [8]
	prefix[0x97] = []func(){nop, cpu.res(2, &cpu.a)}

	// RES 3 B [] 2 [8]
	prefix[0x98] = []func(){nop, cpu.res(3, &cpu.b)}

	// RES 3 C [] 2 [8]
	prefix[0x99] = []func(){nop, cpu.res(3, &cpu.c)}

	// RES 3 D [] 2 [8]
	prefix[0x9a] = []func(){nop, cpu.res(3, &cpu.d)}

	// RES 3 E [] 2 [8]
	prefix[0x9b] = []func(){nop, cpu.res(3, &cpu.e)}

	// RES 3 H [] 2 [8]
	prefix[0x9c] = []func(){nop, cpu.res(3, &cpu.h)}

	// RES 3 L [] 2 [8]
	prefix[0x9d] = []func(){nop, cpu.res(3, &cpu.l)}

	// RES 3 (HL) [] 2 [16]
	prefix[0x9e] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.resM(mapper, 3)}

	// RES 3 A [] 2 [8]
	prefix[0x9f] = []func(){nop, cpu.res(3, &cpu.a)}

	// RES 4 B [] 2 [8]
	prefix[0xa0] = []func(){nop, cpu.res(4, &cpu.b)}

	// RES 4 C [] 2 [8]
	prefix[0xa1] = []func(){nop, cpu.res(4, &cpu.c)}

	// RES 4 D [] 2 [8]
	prefix[0xa2] = []func(){nop, cpu.res(4, &cpu.d)}

	// RES 4 E [] 2 [8]
	prefix[0xa3] = []func(){nop, cpu.res(4, &cpu.e)}

	// RES 4 H [] 2 [8]
	prefix[0xa4] = []func(){nop, cpu.res(4, &cpu.h)}

	// RES 4 L [] 2 [8]
	prefix[0xa5] = []func(){nop, cpu.res(4, &cpu.l)}

	// RES 4 (HL) [] 2 [16]
	prefix[0xa6] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.resM(mapper, 4)}

	// RES 4 A [] 2 [8]
	prefix[0xa7] = []func(){nop, cpu.res(4, &cpu.a)}

	// RES 5 B [] 2 [8]
	prefix[0xa8] = []func(){nop, cpu.res(5, &cpu.b)}

	// RES 5 C [] 2 [8]
	prefix[0xa9] = []func(){nop, cpu.res(5, &cpu.c)}

	// RES 5 D [] 2 [8]
	prefix[0xaa] = []func(){nop, cpu.res(5, &cpu.d)}

	// RES 5 E [] 2 [8]
	prefix[0xab] = []func(){nop, cpu.res(5, &cpu.e)}

	// RES 5 H [] 2 [8]
	prefix[0xac] = []func(){nop, cpu.res(5, &cpu.h)}

	// RES 5 L [] 2 [8]
	prefix[0xad] = []func(){nop, cpu.res(5, &cpu.l)}

	// RES 5 (HL) [] 2 [16]
	prefix[0xae] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.resM(mapper, 5)}

	// RES 5 A [] 2 [8]
	prefix[0xaf] = []func(){nop, cpu.res(5, &cpu.a)}

	// RES 6 B [] 2 [8]
	prefix[0xb0] = []func(){nop, cpu.res(6, &cpu.b)}

	// RES 6 C [] 2 [8]
	prefix[0xb1] = []func(){nop, cpu.res(6, &cpu.c)}

	// RES 6 D [] 2 [8]
	prefix[0xb2] = []func(){nop, cpu.res(6, &cpu.d)}

	// RES 6 E [] 2 [8]
	prefix[0xb3] = []func(){nop, cpu.res(6, &cpu.e)}

	// RES 6 H [] 2 [8]
	prefix[0xb4] = []func(){nop, cpu.res(6, &cpu.h)}

	// RES 6 L [] 2 [8]
	prefix[0xb5] = []func(){nop, cpu.res(6, &cpu.l)}

	// RES 6 (HL) [] 2 [16]
	prefix[0xb6] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.resM(mapper, 6)}

	// RES 6 A [] 2 [8]
	prefix[0xb7] = []func(){nop, cpu.res(6, &cpu.a)}

	// RES 7 B [] 2 [8]
	prefix[0xb8] = []func(){nop, cpu.res(7, &cpu.b)}

	// RES 7 C [] 2 [8]
	prefix[0xb9] = []func(){nop, cpu.res(7, &cpu.c)}

	// RES 7 D [] 2 [8]
	prefix[0xba] = []func(){nop, cpu.res(7, &cpu.d)}

	// RES 7 E [] 2 [8]
	prefix[0xbb] = []func(){nop, cpu.res(7, &cpu.e)}

	// RES 7 H [] 2 [8]
	prefix[0xbc] = []func(){nop, cpu.res(7, &cpu.h)}

	// RES 7 L [] 2 [8]
	prefix[0xbd] = []func(){nop, cpu.res(7, &cpu.l)}

	// RES 7 (HL) [] 2 [16]
	prefix[0xbe] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.resM(mapper, 7)}

	// RES 7 A [] 2 [8]
	prefix[0xbf] = []func(){nop, cpu.res(7, &cpu.a)}

	// SET 0 B [] 2 [8]
	prefix[0xc0] = []func(){nop, cpu.set(0, &cpu.b)}

	// SET 0 C [] 2 [8]
	prefix[0xc1] = []func(){nop, cpu.set(0, &cpu.c)}

	// SET 0 D [] 2 [8]
	prefix[0xc2] = []func(){nop, cpu.set(0, &cpu.d)}

	// SET 0 E [] 2 [8]
	prefix[0xc3] = []func(){nop, cpu.set(0, &cpu.e)}

	// SET 0 H [] 2 [8]
	prefix[0xc4] = []func(){nop, cpu.set(0, &cpu.h)}

	// SET 0 L [] 2 [8]
	prefix[0xc5] = []func(){nop, cpu.set(0, &cpu.l)}

	// SET 0 (HL) [] 2 [16]
	prefix[0xc6] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.setM(mapper, 0)}

	// SET 0 A [] 2 [8]
	prefix[0xc7] = []func(){nop, cpu.set(0, &cpu.a)}

	// SET 1 B [] 2 [8]
	prefix[0xc8] = []func(){nop, cpu.set(1, &cpu.b)}

	// SET 1 C [] 2 [8]
	prefix[0xc9] = []func(){nop, cpu.set(1, &cpu.c)}

	// SET 1 D [] 2 [8]
	prefix[0xca] = []func(){nop, cpu.set(1, &cpu.d)}

	// SET 1 E [] 2 [8]
	prefix[0xcb] = []func(){nop, cpu.set(1, &cpu.e)}

	// SET 1 H [] 2 [8]
	prefix[0xcc] = []func(){nop, cpu.set(1, &cpu.h)}

	// SET 1 L [] 2 [8]
	prefix[0xcd] = []func(){nop, cpu.set(1, &cpu.l)}

	// SET 1 (HL) [] 2 [16]
	prefix[0xce] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.setM(mapper, 1)}

	// SET 1 A [] 2 [8]
	prefix[0xcf] = []func(){nop, cpu.set(1, &cpu.a)}

	// SET 2 B [] 2 [8]
	prefix[0xd0] = []func(){nop, cpu.set(2, &cpu.b)}

	// SET 2 C [] 2 [8]
	prefix[0xd1] = []func(){nop, cpu.set(2, &cpu.c)}

	// SET 2 D [] 2 [8]
	prefix[0xd2] = []func(){nop, cpu.set(2, &cpu.d)}

	// SET 2 E [] 2 [8]
	prefix[0xd3] = []func(){nop, cpu.set(2, &cpu.e)}

	// SET 2 H [] 2 [8]
	prefix[0xd4] = []func(){nop, cpu.set(2, &cpu.h)}

	// SET 2 L [] 2 [8]
	prefix[0xd5] = []func(){nop, cpu.set(2, &cpu.l)}

	// SET 2 (HL) [] 2 [16]
	prefix[0xd6] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.setM(mapper, 2)}

	// SET 2 A [] 2 [8]
	prefix[0xd7] = []func(){nop, cpu.set(2, &cpu.a)}

	// SET 3 B [] 2 [8]
	prefix[0xd8] = []func(){nop, cpu.set(3, &cpu.b)}

	// SET 3 C [] 2 [8]
	prefix[0xd9] = []func(){nop, cpu.set(3, &cpu.c)}

	// SET 3 D [] 2 [8]
	prefix[0xda] = []func(){nop, cpu.set(3, &cpu.d)}

	// SET 3 E [] 2 [8]
	prefix[0xdb] = []func(){nop, cpu.set(3, &cpu.e)}

	// SET 3 H [] 2 [8]
	prefix[0xdc] = []func(){nop, cpu.set(3, &cpu.h)}

	// SET 3 L [] 2 [8]
	prefix[0xdd] = []func(){nop, cpu.set(3, &cpu.l)}

	// SET 3 (HL) [] 2 [16]
	prefix[0xde] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.setM(mapper, 3)}

	// SET 3 A [] 2 [8]
	prefix[0xdf] = []func(){nop, cpu.set(3, &cpu.a)}

	// SET 4 B [] 2 [8]
	prefix[0xe0] = []func(){nop, cpu.set(4, &cpu.b)}

	// SET 4 C [] 2 [8]
	prefix[0xe1] = []func(){nop, cpu.set(4, &cpu.c)}

	// SET 4 D [] 2 [8]
	prefix[0xe2] = []func(){nop, cpu.set(4, &cpu.d)}

	// SET 4 E [] 2 [8]
	prefix[0xe3] = []func(){nop, cpu.set(4, &cpu.e)}

	// SET 4 H [] 2 [8]
	prefix[0xe4] = []func(){nop, cpu.set(4, &cpu.h)}

	// SET 4 L [] 2 [8]
	prefix[0xe5] = []func(){nop, cpu.set(4, &cpu.l)}

	// SET 4 (HL) [] 2 [16]
	prefix[0xe6] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.setM(mapper, 4)}

	// SET 4 A [] 2 [8]
	prefix[0xe7] = []func(){nop, cpu.set(4, &cpu.a)}

	// SET 5 B [] 2 [8]
	prefix[0xe8] = []func(){nop, cpu.set(5, &cpu.b)}

	// SET 5 C [] 2 [8]
	prefix[0xe9] = []func(){nop, cpu.set(5, &cpu.c)}

	// SET 5 D [] 2 [8]
	prefix[0xea] = []func(){nop, cpu.set(5, &cpu.d)}

	// SET 5 E [] 2 [8]
	prefix[0xeb] = []func(){nop, cpu.set(5, &cpu.e)}

	// SET 5 H [] 2 [8]
	prefix[0xec] = []func(){nop, cpu.set(5, &cpu.h)}

	// SET 5 L [] 2 [8]
	prefix[0xed] = []func(){nop, cpu.set(5, &cpu.l)}

	// SET 5 (HL) [] 2 [16]
	prefix[0xee] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.setM(mapper, 5)}

	// SET 5 A [] 2 [8]
	prefix[0xef] = []func(){nop, cpu.set(5, &cpu.a)}

	// SET 6 B [] 2 [8]
	prefix[0xf0] = []func(){nop, cpu.set(6, &cpu.b)}

	// SET 6 C [] 2 [8]
	prefix[0xf1] = []func(){nop, cpu.set(6, &cpu.c)}

	// SET 6 D [] 2 [8]
	prefix[0xf2] = []func(){nop, cpu.set(6, &cpu.d)}

	// SET 6 E [] 2 [8]
	prefix[0xf3] = []func(){nop, cpu.set(6, &cpu.e)}

	// SET 6 H [] 2 [8]
	prefix[0xf4] = []func(){nop, cpu.set(6, &cpu.h)}

	// SET 6 L [] 2 [8]
	prefix[0xf5] = []func(){nop, cpu.set(6, &cpu.l)}

	// SET 6 (HL) [] 2 [16]
	prefix[0xf6] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.setM(mapper, 6)}

	// SET 6 A [] 2 [8]
	prefix[0xf7] = []func(){nop, cpu.set(6, &cpu.a)}

	// SET 7 B [] 2 [8]
	prefix[0xf8] = []func(){nop, cpu.set(7, &cpu.b)}

	// SET 7 C [] 2 [8]
	prefix[0xf9] = []func(){nop, cpu.set(7, &cpu.c)}

	// SET 7 D [] 2 [8]
	prefix[0xfa] = []func(){nop, cpu.set(7, &cpu.d)}

	// SET 7 E [] 2 [8]
	prefix[0xfb] = []func(){nop, cpu.set(7, &cpu.e)}

	// SET 7 H [] 2 [8]
	prefix[0xfc] = []func(){nop, cpu.set(7, &cpu.h)}

	// SET 7 L [] 2 [8]
	prefix[0xfd] = []func(){nop, cpu.set(7, &cpu.l)}

	// SET 7 (HL) [] 2 [16]
	prefix[0xfe] = []func(){nop, nop, cpu.ldMHL(mapper), cpu.setM(mapper, 7)}

	// SET 7 A [] 2 [8]
	prefix[0xff] = []func(){nop, cpu.set(7, &cpu.a)}
}
