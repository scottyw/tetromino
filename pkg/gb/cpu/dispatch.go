package cpu

import "fmt"
import "github.com/scottyw/tetromino/pkg/gb/mem"

func (cpu *CPU) dispatchOneByteInstruction(mem *mem.Memory, instruction uint8) {
	switch instruction {
	case 0x8f:
		cpu.adc(cpu.a) // ADC A A [Z 0 H C]
	case 0x88:
		cpu.adc(cpu.b) // ADC A B [Z 0 H C]
	case 0x89:
		cpu.adc(cpu.c) // ADC A C [Z 0 H C]
	case 0x8a:
		cpu.adc(cpu.d) // ADC A D [Z 0 H C]
	case 0x8b:
		cpu.adc(cpu.e) // ADC A E [Z 0 H C]
	case 0x8c:
		cpu.adc(cpu.h) // ADC A H [Z 0 H C]
	case 0x8e:
		cpu.adcAddr(cpu.hl(), mem) // ADC A (HL) [Z 0 H C]
	case 0x8d:
		cpu.adc(cpu.l) // ADC A L [Z 0 H C]
	case 0x87:
		cpu.add(cpu.a) // ADD A A [Z 0 H C]
	case 0x80:
		cpu.add(cpu.b) // ADD A B [Z 0 H C]
	case 0x81:
		cpu.add(cpu.c) // ADD A C [Z 0 H C]
	case 0x82:
		cpu.add(cpu.d) // ADD A D [Z 0 H C]
	case 0x83:
		cpu.add(cpu.e) // ADD A E [Z 0 H C]
	case 0x84:
		cpu.add(cpu.h) // ADD A H [Z 0 H C]
	case 0x85:
		cpu.add(cpu.l) // ADD A L [Z 0 H C]
	case 0x86:
		cpu.addAddr(cpu.hl(), mem) // ADD A (HL) [Z 0 H C]
	case 0x09:
		cpu.addHL(cpu.bc()) // ADD HL BC [- 0 H C]
	case 0x19:
		cpu.addHL(cpu.de()) // ADD HL DE [- 0 H C]
	case 0x29:
		cpu.addHL(cpu.hl()) // ADD HL HL [- 0 H C]
	case 0x39:
		cpu.addHL(cpu.sp) // ADD HL SP [- 0 H C]
	case 0xa7:
		cpu.and(cpu.a) // AND A  [Z 0 1 0]
	case 0xa0:
		cpu.and(cpu.b) // AND B  [Z 0 1 0]
	case 0xa1:
		cpu.and(cpu.c) // AND C  [Z 0 1 0]
	case 0xa2:
		cpu.and(cpu.d) // AND D  [Z 0 1 0]
	case 0xa3:
		cpu.and(cpu.e) // AND E  [Z 0 1 0]
	case 0xa4:
		cpu.and(cpu.h) // AND H  [Z 0 1 0]
	case 0xa6:
		cpu.andAddr(cpu.hl(), mem) // AND (HL)  [Z 0 1 0]
	case 0xa5:
		cpu.and(cpu.l) // AND L  [Z 0 1 0]
	case 0x3f:
		cpu.ccf() // CCF   [- 0 0 C]
	case 0xbf:
		cpu.cp(cpu.a) // CP A  [Z 1 H C]
	case 0xb8:
		cpu.cp(cpu.b) // CP B  [Z 1 H C]
	case 0xb9:
		cpu.cp(cpu.c) // CP C  [Z 1 H C]
	case 0xba:
		cpu.cp(cpu.d) // CP D  [Z 1 H C]
	case 0xbb:
		cpu.cp(cpu.e) // CP E  [Z 1 H C]
	case 0xbc:
		cpu.cp(cpu.h) // CP H  [Z 1 H C]
	case 0xbe:
		cpu.cpAddr(cpu.hl(), mem) // CP (HL)  [Z 1 H C]
	case 0xbd:
		cpu.cp(cpu.l) // CP L  [Z 1 H C]
	case 0x2f:
		cpu.cpl() // CPL   [- 1 1 -]
	case 0x27:
		cpu.daa() // DAA   [Z - 0 C]
	case 0x3d:
		cpu.dec(&cpu.a) // DEC A  [Z 1 H -]
	case 0x05:
		cpu.dec(&cpu.b) // DEC B  [Z 1 H -]
	case 0x0b:
		cpu.dec16(&cpu.b, &cpu.c) // DEC BC  []
	case 0x0d:
		cpu.dec(&cpu.c) // DEC C  [Z 1 H -]
	case 0x15:
		cpu.dec(&cpu.d) // DEC D  [Z 1 H -]
	case 0x1b:
		cpu.dec16(&cpu.d, &cpu.e) // DEC DE  []
	case 0x1d:
		cpu.dec(&cpu.e) // DEC E  [Z 1 H -]
	case 0x25:
		cpu.dec(&cpu.h) // DEC H  [Z 1 H -]
	case 0x35:
		cpu.decAddr(cpu.hl(), mem) // DEC (HL)  [Z 1 H -]
	case 0x2b:
		cpu.dec16(&cpu.h, &cpu.l) // DEC HL  []
	case 0x2d:
		cpu.dec(&cpu.l) // DEC L  [Z 1 H -]
	case 0x3b:
		cpu.decSP() // DEC SP  []
	case 0xf3:
		cpu.di() // DI   []
	case 0xfb:
		cpu.ei() // EI   []
	case 0x76:
		cpu.halt() // HALT   []
	case 0x3c:
		cpu.inc(&cpu.a) // INC A  [Z 0 H -]
	case 0x04:
		cpu.inc(&cpu.b) // INC B  [Z 0 H -]
	case 0x03:
		cpu.inc16(&cpu.b, &cpu.c) // INC BC  []
	case 0x0c:
		cpu.inc(&cpu.c) // INC C  [Z 0 H -]
	case 0x14:
		cpu.inc(&cpu.d) // INC D  [Z 0 H -]
	case 0x13:
		cpu.inc16(&cpu.d, &cpu.e) // INC DE  []
	case 0x1c:
		cpu.inc(&cpu.e) // INC E  [Z 0 H -]
	case 0x24:
		cpu.inc(&cpu.h) // INC H  [Z 0 H -]
	case 0x34:
		cpu.incAddr(cpu.hl(), mem) // INC (HL)  [Z 0 H -]
	case 0x23:
		cpu.inc16(&cpu.h, &cpu.l) // INC HL  []
	case 0x2c:
		cpu.inc(&cpu.l) // INC L  [Z 0 H -]
	case 0x33:
		cpu.incSP() // INC SP  []
	case 0xe9:
		cpu.jp("", cpu.hl()) // JP (HL)  []
	case 0xe2:
		cpu.ldhA8A(cpu.c, mem) // LD (C) A []
	case 0x12:
		cpu.ldA16U8(cpu.de(), cpu.a, mem) // LD (DE) A []
	case 0xf2:
		cpu.ldhAA8(cpu.c, mem) // LD A (C) []
	case 0x1a:
		cpu.ldR8A16(&cpu.a, cpu.de(), mem) // LD A (DE) []
	case 0x3a:
		cpu.lddAA16(mem) // LD A (HL-) []
	case 0x32:
		cpu.lddA16A(mem) // LD (HL-) A []
	case 0x2a:
		cpu.ldiAA16(mem) // LD A (HL+) []
	case 0x22:
		cpu.ldiA16A(mem) // LD (HL+) A []
	case 0x7f:
		cpu.ld(&cpu.a, cpu.a) // LD A A []
	case 0x78:
		cpu.ld(&cpu.a, cpu.b) // LD A B []
	case 0x0a:
		cpu.ldR8A16(&cpu.a, cpu.bc(), mem) // LD A (BC) []
	case 0x79:
		cpu.ld(&cpu.a, cpu.c) // LD A C []
	case 0x7a:
		cpu.ld(&cpu.a, cpu.d) // LD A D []
	case 0x7b:
		cpu.ld(&cpu.a, cpu.e) // LD A E []
	case 0x7c:
		cpu.ld(&cpu.a, cpu.h) // LD A H []
	case 0x7e:
		cpu.ldR8A16(&cpu.a, cpu.hl(), mem) // LD A (HL) []
	case 0x7d:
		cpu.ld(&cpu.a, cpu.l) // LD A L []
	case 0x47:
		cpu.ld(&cpu.b, cpu.a) // LD B A []
	case 0x40:
		cpu.ld(&cpu.b, cpu.b) // LD B B []
		// Mooneye uses this instruction as a magic breakpoint
		// to indicate that a test rom has completeed
		cpu.Mooneye = true
	case 0x41:
		cpu.ld(&cpu.b, cpu.c) // LD B C []
	case 0x42:
		cpu.ld(&cpu.b, cpu.d) // LD B D []
	case 0x43:
		cpu.ld(&cpu.b, cpu.e) // LD B E []
	case 0x44:
		cpu.ld(&cpu.b, cpu.h) // LD B H []
	case 0x46:
		cpu.ldR8A16(&cpu.b, cpu.hl(), mem) // LD B (HL) []
	case 0x45:
		cpu.ld(&cpu.b, cpu.l) // LD B L []
	case 0x02:
		cpu.ldA16U8(cpu.bc(), cpu.a, mem) // LD (BC) A []
	case 0x4f:
		cpu.ld(&cpu.c, cpu.a) // LD C A []
	case 0x48:
		cpu.ld(&cpu.c, cpu.b) // LD C B []
	case 0x49:
		cpu.ld(&cpu.c, cpu.c) // LD C C []
	case 0x4a:
		cpu.ld(&cpu.c, cpu.d) // LD C D []
	case 0x4b:
		cpu.ld(&cpu.c, cpu.e) // LD C E []
	case 0x4c:
		cpu.ld(&cpu.c, cpu.h) // LD C H []
	case 0x4e:
		cpu.ldR8A16(&cpu.c, cpu.hl(), mem) // LD C (HL) []
	case 0x4d:
		cpu.ld(&cpu.c, cpu.l) // LD C L []
	case 0x57:
		cpu.ld(&cpu.d, cpu.a) // LD D A []
	case 0x50:
		cpu.ld(&cpu.d, cpu.b) // LD D B []
	case 0x51:
		cpu.ld(&cpu.d, cpu.c) // LD D C []
	case 0x52:
		cpu.ld(&cpu.d, cpu.d) // LD D D []
	case 0x53:
		cpu.ld(&cpu.d, cpu.e) // LD D E []
	case 0x54:
		cpu.ld(&cpu.d, cpu.h) // LD D H []
	case 0x56:
		cpu.ldR8A16(&cpu.d, cpu.hl(), mem) // LD D (HL) []
	case 0x55:
		cpu.ld(&cpu.d, cpu.l) // LD D L []
	case 0x5f:
		cpu.ld(&cpu.e, cpu.a) // LD E A []
	case 0x58:
		cpu.ld(&cpu.e, cpu.b) // LD E B []
	case 0x59:
		cpu.ld(&cpu.e, cpu.c) // LD E C []
	case 0x5a:
		cpu.ld(&cpu.e, cpu.d) // LD E D []
	case 0x5b:
		cpu.ld(&cpu.e, cpu.e) // LD E E []
	case 0x5c:
		cpu.ld(&cpu.e, cpu.h) // LD E H []
	case 0x5e:
		cpu.ldR8A16(&cpu.e, cpu.hl(), mem) // LD E (HL) []
	case 0x5d:
		cpu.ld(&cpu.e, cpu.l) // LD E L []
	case 0x67:
		cpu.ld(&cpu.h, cpu.a) // LD H A []
	case 0x60:
		cpu.ld(&cpu.h, cpu.b) // LD H B []
	case 0x61:
		cpu.ld(&cpu.h, cpu.c) // LD H C []
	case 0x62:
		cpu.ld(&cpu.h, cpu.d) // LD H D []
	case 0x63:
		cpu.ld(&cpu.h, cpu.e) // LD H E []
	case 0x64:
		cpu.ld(&cpu.h, cpu.h) // LD H H []
	case 0x66:
		cpu.ldR8A16(&cpu.h, cpu.hl(), mem) // LD H (HL) []
	case 0x65:
		cpu.ld(&cpu.h, cpu.l) // LD H L []
	case 0x77:
		cpu.ldA16U8(cpu.hl(), cpu.a, mem) // LD (HL) A []
	case 0x70:
		cpu.ldA16U8(cpu.hl(), cpu.b, mem) // LD (HL) B []
	case 0x71:
		cpu.ldA16U8(cpu.hl(), cpu.c, mem) // LD (HL) C []
	case 0x72:
		cpu.ldA16U8(cpu.hl(), cpu.d, mem) // LD (HL) D []
	case 0x73:
		cpu.ldA16U8(cpu.hl(), cpu.e, mem) // LD (HL) E []
	case 0x74:
		cpu.ldA16U8(cpu.hl(), cpu.h, mem) // LD (HL) H []
	case 0x75:
		cpu.ldA16U8(cpu.hl(), cpu.l, mem) // LD (HL) L []
	case 0x6f:
		cpu.ld(&cpu.l, cpu.a) // LD L A []
	case 0x68:
		cpu.ld(&cpu.l, cpu.b) // LD L B []
	case 0x69:
		cpu.ld(&cpu.l, cpu.c) // LD L C []
	case 0x6a:
		cpu.ld(&cpu.l, cpu.d) // LD L D []
	case 0x6b:
		cpu.ld(&cpu.l, cpu.e) // LD L E []
	case 0x6c:
		cpu.ld(&cpu.l, cpu.h) // LD L H []
	case 0x6e:
		cpu.ldR8A16(&cpu.l, cpu.hl(), mem) // LD L (HL) []
	case 0x6d:
		cpu.ld(&cpu.l, cpu.l) // LD L L []
	case 0xf9:
		cpu.ldSPHL() // LD SP HL []
	case 0x00:
		cpu.nop() // NOP   []
	case 0xb7:
		cpu.or(cpu.a) // OR A  [Z 0 0 0]
	case 0xb0:
		cpu.or(cpu.b) // OR B  [Z 0 0 0]
	case 0xb1:
		cpu.or(cpu.c) // OR C  [Z 0 0 0]
	case 0xb2:
		cpu.or(cpu.d) // OR D  [Z 0 0 0]
	case 0xb3:
		cpu.or(cpu.e) // OR E  [Z 0 0 0]
	case 0xb4:
		cpu.or(cpu.h) // OR H  [Z 0 0 0]
	case 0xb6:
		cpu.orAddr(cpu.hl(), mem) // OR (HL)  [Z 0 0 0]
	case 0xb5:
		cpu.or(cpu.l) // OR L  [Z 0 0 0]
	case 0xf1:
		cpu.popAF(mem) // POP AF  [Z N H C]
	case 0xc1:
		cpu.pop(&cpu.b, &cpu.c, mem) // POP BC  []
	case 0xd1:
		cpu.pop(&cpu.d, &cpu.e, mem) // POP DE  []
	case 0xe1:
		cpu.pop(&cpu.h, &cpu.l, mem) // POP HL  []
	case 0xf5:
		cpu.push(cpu.a, cpu.f, mem) // PUSH AF  []
	case 0xc5:
		cpu.push(cpu.b, cpu.c, mem) // PUSH BC  []
	case 0xd5:
		cpu.push(cpu.d, cpu.e, mem) // PUSH DE  []
	case 0xe5:
		cpu.push(cpu.h, cpu.l, mem) // PUSH HL  []
	case 0xc9:
		cpu.ret("", mem) // RET   []
	case 0xd8:
		cpu.ret("C", mem) // RET C  []
	case 0xd0:
		cpu.ret("NC", mem) // RET NC  []
	case 0xc0:
		cpu.ret("NZ", mem) // RET NZ  []
	case 0xc8:
		cpu.ret("Z", mem) // RET Z  []
	case 0xd9:
		cpu.reti(mem) // RETI   []
	case 0x17:
		cpu.rla() // RLA   [0 0 0 C]
	case 0x07:
		cpu.rlca() // RLCA   [0 0 0 C]
	case 0x1f:
		cpu.rra() // RRA   [0 0 0 C]
	case 0x0f:
		cpu.rrca() // RRCA   [0 0 0 C]
	case 0xc7:
		cpu.rst(0x0000, mem) // RST 00H  []
	case 0xcf:
		cpu.rst(0x0008, mem) // RST 08H  []
	case 0xd7:
		cpu.rst(0x0010, mem) // RST 10H  []
	case 0xdf:
		cpu.rst(0x0018, mem) // RST 18H  []
	case 0xe7:
		cpu.rst(0x0020, mem) // RST 20H  []
	case 0xef:
		cpu.rst(0x0028, mem) // RST 28H  []
	case 0xf7:
		cpu.rst(0x0030, mem) // RST 30H  []
	case 0xff:
		cpu.rst(0x0038, mem) // RST 38H  []
	case 0x9f:
		cpu.sbc(cpu.a) // SBC A A [Z 1 H C]
	case 0x98:
		cpu.sbc(cpu.b) // SBC A B [Z 1 H C]
	case 0x99:
		cpu.sbc(cpu.c) // SBC A C [Z 1 H C]
	case 0x9a:
		cpu.sbc(cpu.d) // SBC A D [Z 1 H C]
	case 0x9b:
		cpu.sbc(cpu.e) // SBC A E [Z 1 H C]
	case 0x9c:
		cpu.sbc(cpu.h) // SBC A H [Z 1 H C]
	case 0x9e:
		cpu.sbcAddr(cpu.hl(), mem) // SBC A (HL) [Z 1 H C]
	case 0x9d:
		cpu.sbc(cpu.l) // SBC A L [Z 1 H C]
	case 0x37:
		cpu.scf() // SCF   [- 0 0 1]
	case 0x10:
		cpu.stop() // STOP 0  []
	case 0x97:
		cpu.sub(cpu.a) // SUB A  [Z 1 H C]
	case 0x90:
		cpu.sub(cpu.b) // SUB B  [Z 1 H C]
	case 0x91:
		cpu.sub(cpu.c) // SUB C  [Z 1 H C]
	case 0x92:
		cpu.sub(cpu.d) // SUB D  [Z 1 H C]
	case 0x93:
		cpu.sub(cpu.e) // SUB E  [Z 1 H C]
	case 0x94:
		cpu.sub(cpu.h) // SUB H  [Z 1 H C]
	case 0x96:
		cpu.subAddr(cpu.hl(), mem) // SUB (HL)  [Z 1 H C]
	case 0x95:
		cpu.sub(cpu.l) // SUB L  [Z 1 H C]
	case 0xaf:
		cpu.xor(cpu.a) // XOR A  [Z 0 0 0]
	case 0xa8:
		cpu.xor(cpu.b) // XOR B  [Z 0 0 0]
	case 0xa9:
		cpu.xor(cpu.c) // XOR C  [Z 0 0 0]
	case 0xaa:
		cpu.xor(cpu.d) // XOR D  [Z 0 0 0]
	case 0xab:
		cpu.xor(cpu.e) // XOR E  [Z 0 0 0]
	case 0xac:
		cpu.xor(cpu.h) // XOR H  [Z 0 0 0]
	case 0xae:
		cpu.xorAddr(cpu.hl(), mem) // XOR (HL)  [Z 0 0 0]
	case 0xad:
		cpu.xor(cpu.l) // XOR L  [Z 0 0 0]
	default:
		panic(fmt.Sprintf("Missing dispatchOneByteInstruction: %s %02x", instructionMetadata[instruction].Mnemonic, instruction))
	}
}

func (cpu *CPU) dispatchTwoByteInstruction(mem *mem.Memory, instruction, u8 uint8) {
	switch instruction {
	case 0xce:
		cpu.adc(u8) // ADC A d8 [Z 0 H C]
	case 0xc6:
		cpu.add(u8) // ADD A d8 [Z 0 H C]
	case 0xe6:
		cpu.and(u8) // AND d8  [Z 0 1 0]
	case 0xe8:
		cpu.addSP(int8(u8)) // ADD SP r8 [0 0 H C]
	case 0xfe:
		cpu.cp(u8) // CP d8  [Z 1 H C]
	case 0x38:
		cpu.jr("C", int8(u8)) // JR C r8 []
	case 0x30:
		cpu.jr("NC", int8(u8)) // JR NC r8 []
	case 0x20:
		cpu.jr("NZ", int8(u8)) // JR NZ r8 []
	case 0x18:
		cpu.jr("", int8(u8)) // JR r8  []
	case 0x28:
		cpu.jr("Z", int8(u8)) // JR Z r8 []
	case 0x3e:
		cpu.ld(&cpu.a, u8) // LD A d8 []
	case 0x06:
		cpu.ld(&cpu.b, u8) // LD B d8 []
	case 0x0e:
		cpu.ld(&cpu.c, u8) // LD C d8 []
	case 0x16:
		cpu.ld(&cpu.d, u8) // LD D d8 []
	case 0x1e:
		cpu.ld(&cpu.e, u8) // LD E d8 []
	case 0x26:
		cpu.ld(&cpu.h, u8) // LD H d8 []
	case 0x2e:
		cpu.ld(&cpu.l, u8) // LD L d8 []
	case 0x36:
		cpu.ldA16U8(cpu.hl(), u8, mem) // LD (HL) d8 []
	case 0xf8:
		cpu.ldHLSP(int8(u8)) // LD HL SP+r8 [0 0 H C]
	case 0xe0:
		cpu.ldhA8A(u8, mem) // LDH (a8) A []
	case 0xf0:
		cpu.ldhAA8(u8, mem) // LDH A (a8) []
	case 0xf6:
		cpu.or(u8) // OR d8  [Z 0 0 0]
	case 0xde:
		cpu.sbc(u8) // SBC A d8 [Z 1 H C]
	case 0xd6:
		cpu.sub(u8) // SUB d8  [Z 1 H C]
	case 0xee:
		cpu.xor(u8) // XOR d8  [Z 0 0 0]
	default:
		panic(fmt.Sprintf("Missing dispatchTwoByteInstruction: %s %02x %02x", instructionMetadata[instruction].Mnemonic, instruction, u8))
	}
}

func (cpu *CPU) dispatchThreeByteInstruction(mem *mem.Memory, instruction uint8, u16 uint16) {
	switch instruction {
	case 0xcd:
		cpu.call("", u16, mem) // CALL a16  []
	case 0xdc:
		cpu.call("C", u16, mem) // CALL C a16 []
	case 0xd4:
		cpu.call("NC", u16, mem) // CALL NC a16 []
	case 0xc4:
		cpu.call("NZ", u16, mem) // CALL NZ a16 []
	case 0xcc:
		cpu.call("Z", u16, mem) // CALL Z a16 []
	case 0xc3:
		cpu.jp("", u16) // JP a16  []
	case 0xda:
		cpu.jp("C", u16) // JP C a16 []
	case 0xd2:
		cpu.jp("NC", u16) // JP NC a16 []
	case 0xca:
		cpu.jp("Z", u16) // JP Z a16 []
	case 0xc2:
		cpu.jp("NZ", u16) // JP NZ a16 []
	case 0xea:
		cpu.ldA16U8(u16, cpu.a, mem) // LD (a16) A []
	case 0x08:
		cpu.ldA16SP(u16, mem) // LD (a16) SP []
	case 0xfa:
		cpu.ldR8A16(&cpu.a, u16, mem) // LD A (a16) []
	case 0x01:
		cpu.ld16(&cpu.b, &cpu.c, u16) // LD BC d16 []
	case 0x11:
		cpu.ld16(&cpu.d, &cpu.e, u16) // LD DE d16 []
	case 0x21:
		cpu.ld16(&cpu.h, &cpu.l, u16) // LD HL d16 []
	case 0x31:
		cpu.ldSP(u16) // LD SP d16 []
	default:
		panic(fmt.Sprintf("Missing dispatchThreeByteInstruction: %s %02x %04x", instructionMetadata[instruction].Mnemonic, instruction, u16))
	}
}

func (cpu *CPU) dispatchPrefixedInstruction(mem *mem.Memory, instruction uint8) {
	switch instruction {
	case 0x47:
		cpu.bit(0, cpu.a) // BIT 0 A [Z 0 1 -]
	case 0x40:
		cpu.bit(0, cpu.b) // BIT 0 B [Z 0 1 -]
	case 0x41:
		cpu.bit(0, cpu.c) // BIT 0 C [Z 0 1 -]
	case 0x42:
		cpu.bit(0, cpu.d) // BIT 0 D [Z 0 1 -]
	case 0x43:
		cpu.bit(0, cpu.e) // BIT 0 E [Z 0 1 -]
	case 0x44:
		cpu.bit(0, cpu.h) // BIT 0 H [Z 0 1 -]
	case 0x46:
		cpu.bitAddr(0, cpu.hl(), mem) // BIT 0 (HL) [Z 0 1 -]
	case 0x45:
		cpu.bit(0, cpu.l) // BIT 0 L [Z 0 1 -]
	case 0x4f:
		cpu.bit(1, cpu.a) // BIT 1 A [Z 0 1 -]
	case 0x48:
		cpu.bit(1, cpu.b) // BIT 1 B [Z 0 1 -]
	case 0x49:
		cpu.bit(1, cpu.c) // BIT 1 C [Z 0 1 -]
	case 0x4a:
		cpu.bit(1, cpu.d) // BIT 1 D [Z 0 1 -]
	case 0x4b:
		cpu.bit(1, cpu.e) // BIT 1 E [Z 0 1 -]
	case 0x4c:
		cpu.bit(1, cpu.h) // BIT 1 H [Z 0 1 -]
	case 0x4e:
		cpu.bitAddr(1, cpu.hl(), mem) // BIT 1 (HL) [Z 0 1 -]
	case 0x4d:
		cpu.bit(1, cpu.l) // BIT 1 L [Z 0 1 -]
	case 0x57:
		cpu.bit(2, cpu.a) // BIT 2 A [Z 0 1 -]
	case 0x50:
		cpu.bit(2, cpu.b) // BIT 2 B [Z 0 1 -]
	case 0x51:
		cpu.bit(2, cpu.c) // BIT 2 C [Z 0 1 -]
	case 0x52:
		cpu.bit(2, cpu.d) // BIT 2 D [Z 0 1 -]
	case 0x53:
		cpu.bit(2, cpu.e) // BIT 2 E [Z 0 1 -]
	case 0x54:
		cpu.bit(2, cpu.h) // BIT 2 H [Z 0 1 -]
	case 0x56:
		cpu.bitAddr(2, cpu.hl(), mem) // BIT 2 (HL) [Z 0 1 -]
	case 0x55:
		cpu.bit(2, cpu.l) // BIT 2 L [Z 0 1 -]
	case 0x5f:
		cpu.bit(3, cpu.a) // BIT 3 A [Z 0 1 -]
	case 0x58:
		cpu.bit(3, cpu.b) // BIT 3 B [Z 0 1 -]
	case 0x59:
		cpu.bit(3, cpu.c) // BIT 3 C [Z 0 1 -]
	case 0x5a:
		cpu.bit(3, cpu.d) // BIT 3 D [Z 0 1 -]
	case 0x5b:
		cpu.bit(3, cpu.e) // BIT 3 E [Z 0 1 -]
	case 0x5c:
		cpu.bit(3, cpu.h) // BIT 3 H [Z 0 1 -]
	case 0x5e:
		cpu.bitAddr(3, cpu.hl(), mem) // BIT 3 (HL) [Z 0 1 -]
	case 0x5d:
		cpu.bit(3, cpu.l) // BIT 3 L [Z 0 1 -]
	case 0x67:
		cpu.bit(4, cpu.a) // BIT 4 A [Z 0 1 -]
	case 0x60:
		cpu.bit(4, cpu.b) // BIT 4 B [Z 0 1 -]
	case 0x61:
		cpu.bit(4, cpu.c) // BIT 4 C [Z 0 1 -]
	case 0x62:
		cpu.bit(4, cpu.d) // BIT 4 D [Z 0 1 -]
	case 0x63:
		cpu.bit(4, cpu.e) // BIT 4 E [Z 0 1 -]
	case 0x64:
		cpu.bit(4, cpu.h) // BIT 4 H [Z 0 1 -]
	case 0x66:
		cpu.bitAddr(4, cpu.hl(), mem) // BIT 4 (HL) [Z 0 1 -]
	case 0x65:
		cpu.bit(4, cpu.l) // BIT 4 L [Z 0 1 -]
	case 0x6f:
		cpu.bit(5, cpu.a) // BIT 5 A [Z 0 1 -]
	case 0x68:
		cpu.bit(5, cpu.b) // BIT 5 B [Z 0 1 -]
	case 0x69:
		cpu.bit(5, cpu.c) // BIT 5 C [Z 0 1 -]
	case 0x6a:
		cpu.bit(5, cpu.d) // BIT 5 D [Z 0 1 -]
	case 0x6b:
		cpu.bit(5, cpu.e) // BIT 5 E [Z 0 1 -]
	case 0x6c:
		cpu.bit(5, cpu.h) // BIT 5 H [Z 0 1 -]
	case 0x6e:
		cpu.bitAddr(5, cpu.hl(), mem) // BIT 5 (HL) [Z 0 1 -]
	case 0x6d:
		cpu.bit(5, cpu.l) // BIT 5 L [Z 0 1 -]
	case 0x77:
		cpu.bit(6, cpu.a) // BIT 6 A [Z 0 1 -]
	case 0x70:
		cpu.bit(6, cpu.b) // BIT 6 B [Z 0 1 -]
	case 0x71:
		cpu.bit(6, cpu.c) // BIT 6 C [Z 0 1 -]
	case 0x72:
		cpu.bit(6, cpu.d) // BIT 6 D [Z 0 1 -]
	case 0x73:
		cpu.bit(6, cpu.e) // BIT 6 E [Z 0 1 -]
	case 0x74:
		cpu.bit(6, cpu.h) // BIT 6 H [Z 0 1 -]
	case 0x76:
		cpu.bitAddr(6, cpu.hl(), mem) // BIT 6 (HL) [Z 0 1 -]
	case 0x75:
		cpu.bit(6, cpu.l) // BIT 6 L [Z 0 1 -]
	case 0x7f:
		cpu.bit(7, cpu.a) // BIT 7 A [Z 0 1 -]
	case 0x78:
		cpu.bit(7, cpu.b) // BIT 7 B [Z 0 1 -]
	case 0x79:
		cpu.bit(7, cpu.c) // BIT 7 C [Z 0 1 -]
	case 0x7a:
		cpu.bit(7, cpu.d) // BIT 7 D [Z 0 1 -]
	case 0x7b:
		cpu.bit(7, cpu.e) // BIT 7 E [Z 0 1 -]
	case 0x7c:
		cpu.bit(7, cpu.h) // BIT 7 H [Z 0 1 -]
	case 0x7e:
		cpu.bitAddr(7, cpu.hl(), mem) // BIT 7 (HL) [Z 0 1 -]
	case 0x7d:
		cpu.bit(7, cpu.l) // BIT 7 L [Z 0 1 -]
	case 0x87:
		cpu.res(0, &cpu.a) // RES 0 A []
	case 0x80:
		cpu.res(0, &cpu.b) // RES 0 B []
	case 0x81:
		cpu.res(0, &cpu.c) // RES 0 C []
	case 0x82:
		cpu.res(0, &cpu.d) // RES 0 D []
	case 0x83:
		cpu.res(0, &cpu.e) // RES 0 E []
	case 0x84:
		cpu.res(0, &cpu.h) // RES 0 H []
	case 0x86:
		cpu.resAddr(0, cpu.hl(), mem) // RES 0 (HL) []
	case 0x85:
		cpu.res(0, &cpu.l) // RES 0 L []
	case 0x8f:
		cpu.res(1, &cpu.a) // RES 1 A []
	case 0x88:
		cpu.res(1, &cpu.b) // RES 1 B []
	case 0x89:
		cpu.res(1, &cpu.c) // RES 1 C []
	case 0x8a:
		cpu.res(1, &cpu.d) // RES 1 D []
	case 0x8b:
		cpu.res(1, &cpu.e) // RES 1 E []
	case 0x8c:
		cpu.res(1, &cpu.h) // RES 1 H []
	case 0x8e:
		cpu.resAddr(1, cpu.hl(), mem) // RES 1 (HL) []
	case 0x8d:
		cpu.res(1, &cpu.l) // RES 1 L []
	case 0x97:
		cpu.res(2, &cpu.a) // RES 2 A []
	case 0x90:
		cpu.res(2, &cpu.b) // RES 2 B []
	case 0x91:
		cpu.res(2, &cpu.c) // RES 2 C []
	case 0x92:
		cpu.res(2, &cpu.d) // RES 2 D []
	case 0x93:
		cpu.res(2, &cpu.e) // RES 2 E []
	case 0x94:
		cpu.res(2, &cpu.h) // RES 2 H []
	case 0x96:
		cpu.resAddr(2, cpu.hl(), mem) // RES 2 (HL) []
	case 0x95:
		cpu.res(2, &cpu.l) // RES 2 L []
	case 0x9f:
		cpu.res(3, &cpu.a) // RES 3 A []
	case 0x98:
		cpu.res(3, &cpu.b) // RES 3 B []
	case 0x99:
		cpu.res(3, &cpu.c) // RES 3 C []
	case 0x9a:
		cpu.res(3, &cpu.d) // RES 3 D []
	case 0x9b:
		cpu.res(3, &cpu.e) // RES 3 E []
	case 0x9c:
		cpu.res(3, &cpu.h) // RES 3 H []
	case 0x9e:
		cpu.resAddr(3, cpu.hl(), mem) // RES 3 (HL) []
	case 0x9d:
		cpu.res(3, &cpu.l) // RES 3 L []
	case 0xa7:
		cpu.res(4, &cpu.a) // RES 4 A []
	case 0xa0:
		cpu.res(4, &cpu.b) // RES 4 B []
	case 0xa1:
		cpu.res(4, &cpu.c) // RES 4 C []
	case 0xa2:
		cpu.res(4, &cpu.d) // RES 4 D []
	case 0xa3:
		cpu.res(4, &cpu.e) // RES 4 E []
	case 0xa4:
		cpu.res(4, &cpu.h) // RES 4 H []
	case 0xa6:
		cpu.resAddr(4, cpu.hl(), mem) // RES 4 (HL) []
	case 0xa5:
		cpu.res(4, &cpu.l) // RES 4 L []
	case 0xaf:
		cpu.res(5, &cpu.a) // RES 5 A []
	case 0xa8:
		cpu.res(5, &cpu.b) // RES 5 B []
	case 0xa9:
		cpu.res(5, &cpu.c) // RES 5 C []
	case 0xaa:
		cpu.res(5, &cpu.d) // RES 5 D []
	case 0xab:
		cpu.res(5, &cpu.e) // RES 5 E []
	case 0xac:
		cpu.res(5, &cpu.h) // RES 5 H []
	case 0xae:
		cpu.resAddr(5, cpu.hl(), mem) // RES 5 (HL) []
	case 0xad:
		cpu.res(5, &cpu.l) // RES 5 L []
	case 0xb7:
		cpu.res(6, &cpu.a) // RES 6 A []
	case 0xb0:
		cpu.res(6, &cpu.b) // RES 6 B []
	case 0xb1:
		cpu.res(6, &cpu.c) // RES 6 C []
	case 0xb2:
		cpu.res(6, &cpu.d) // RES 6 D []
	case 0xb3:
		cpu.res(6, &cpu.e) // RES 6 E []
	case 0xb4:
		cpu.res(6, &cpu.h) // RES 6 H []
	case 0xb6:
		cpu.resAddr(6, cpu.hl(), mem) // RES 6 (HL) []
	case 0xb5:
		cpu.res(6, &cpu.l) // RES 6 L []
	case 0xbf:
		cpu.res(7, &cpu.a) // RES 7 A []
	case 0xb8:
		cpu.res(7, &cpu.b) // RES 7 B []
	case 0xb9:
		cpu.res(7, &cpu.c) // RES 7 C []
	case 0xba:
		cpu.res(7, &cpu.d) // RES 7 D []
	case 0xbb:
		cpu.res(7, &cpu.e) // RES 7 E []
	case 0xbc:
		cpu.res(7, &cpu.h) // RES 7 H []
	case 0xbe:
		cpu.resAddr(7, cpu.hl(), mem) // RES 7 (HL) []
	case 0xbd:
		cpu.res(7, &cpu.l) // RES 7 L []
	case 0x17:
		cpu.rl(&cpu.a) // RL A  [Z 0 0 C]
	case 0x10:
		cpu.rl(&cpu.b) // RL B  [Z 0 0 C]
	case 0x11:
		cpu.rl(&cpu.c) // RL C  [Z 0 0 C]
	case 0x12:
		cpu.rl(&cpu.d) // RL D  [Z 0 0 C]
	case 0x13:
		cpu.rl(&cpu.e) // RL E  [Z 0 0 C]
	case 0x14:
		cpu.rl(&cpu.h) // RL H  [Z 0 0 C]
	case 0x16:
		cpu.rlAddr(cpu.hl(), mem) // RL (HL)  [Z 0 0 C]
	case 0x15:
		cpu.rl(&cpu.l) // RL L  [Z 0 0 C]
	case 0x07:
		cpu.rlc(&cpu.a) // RLC A  [Z 0 0 C]
	case 0x00:
		cpu.rlc(&cpu.b) // RLC B  [Z 0 0 C]
	case 0x01:
		cpu.rlc(&cpu.c) // RLC C  [Z 0 0 C]
	case 0x02:
		cpu.rlc(&cpu.d) // RLC D  [Z 0 0 C]
	case 0x03:
		cpu.rlc(&cpu.e) // RLC E  [Z 0 0 C]
	case 0x04:
		cpu.rlc(&cpu.h) // RLC H  [Z 0 0 C]
	case 0x06:
		cpu.rlcAddr(cpu.hl(), mem) // RLC (HL)  [Z 0 0 C]
	case 0x05:
		cpu.rlc(&cpu.l) // RLC L  [Z 0 0 C]
	case 0x1f:
		cpu.rr(&cpu.a) // RR A  [Z 0 0 C]
	case 0x18:
		cpu.rr(&cpu.b) // RR B  [Z 0 0 C]
	case 0x19:
		cpu.rr(&cpu.c) // RR C  [Z 0 0 C]
	case 0x1a:
		cpu.rr(&cpu.d) // RR D  [Z 0 0 C]
	case 0x1b:
		cpu.rr(&cpu.e) // RR E  [Z 0 0 C]
	case 0x1c:
		cpu.rr(&cpu.h) // RR H  [Z 0 0 C]
	case 0x1e:
		cpu.rrAddr(cpu.hl(), mem) // RR (HL)  [Z 0 0 C]
	case 0x1d:
		cpu.rr(&cpu.l) // RR L  [Z 0 0 C]
	case 0x0f:
		cpu.rrc(&cpu.a) // RRC A  [Z 0 0 C]
	case 0x08:
		cpu.rrc(&cpu.b) // RRC B  [Z 0 0 C]
	case 0x09:
		cpu.rrc(&cpu.c) // RRC C  [Z 0 0 C]
	case 0x0a:
		cpu.rrc(&cpu.d) // RRC D  [Z 0 0 C]
	case 0x0b:
		cpu.rrc(&cpu.e) // RRC E  [Z 0 0 C]
	case 0x0c:
		cpu.rrc(&cpu.h) // RRC H  [Z 0 0 C]
	case 0x0e:
		cpu.rrcAddr(cpu.hl(), mem) // RRC (HL)  [Z 0 0 C]
	case 0x0d:
		cpu.rrc(&cpu.l) // RRC L  [Z 0 0 C]
	case 0xc7:
		cpu.set(0, &cpu.a) // SET 0 A []
	case 0xc0:
		cpu.set(0, &cpu.b) // SET 0 B []
	case 0xc1:
		cpu.set(0, &cpu.c) // SET 0 C []
	case 0xc2:
		cpu.set(0, &cpu.d) // SET 0 D []
	case 0xc3:
		cpu.set(0, &cpu.e) // SET 0 E []
	case 0xc4:
		cpu.set(0, &cpu.h) // SET 0 H []
	case 0xc6:
		cpu.setAddr(0, cpu.hl(), mem) // SET 0 (HL) []
	case 0xc5:
		cpu.set(0, &cpu.l) // SET 0 L []
	case 0xcf:
		cpu.set(1, &cpu.a) // SET 1 A []
	case 0xc8:
		cpu.set(1, &cpu.b) // SET 1 B []
	case 0xc9:
		cpu.set(1, &cpu.c) // SET 1 C []
	case 0xca:
		cpu.set(1, &cpu.d) // SET 1 D []
	case 0xcb:
		cpu.set(1, &cpu.e) // SET 1 E []
	case 0xcc:
		cpu.set(1, &cpu.h) // SET 1 H []
	case 0xce:
		cpu.setAddr(1, cpu.hl(), mem) // SET 1 (HL) []
	case 0xcd:
		cpu.set(1, &cpu.l) // SET 1 L []
	case 0xd7:
		cpu.set(2, &cpu.a) // SET 2 A []
	case 0xd0:
		cpu.set(2, &cpu.b) // SET 2 B []
	case 0xd1:
		cpu.set(2, &cpu.c) // SET 2 C []
	case 0xd2:
		cpu.set(2, &cpu.d) // SET 2 D []
	case 0xd3:
		cpu.set(2, &cpu.e) // SET 2 E []
	case 0xd4:
		cpu.set(2, &cpu.h) // SET 2 H []
	case 0xd6:
		cpu.setAddr(2, cpu.hl(), mem) // SET 2 (HL) []
	case 0xd5:
		cpu.set(2, &cpu.l) // SET 2 L []
	case 0xdf:
		cpu.set(3, &cpu.a) // SET 3 A []
	case 0xd8:
		cpu.set(3, &cpu.b) // SET 3 B []
	case 0xd9:
		cpu.set(3, &cpu.c) // SET 3 C []
	case 0xda:
		cpu.set(3, &cpu.d) // SET 3 D []
	case 0xdb:
		cpu.set(3, &cpu.e) // SET 3 E []
	case 0xdc:
		cpu.set(3, &cpu.h) // SET 3 H []
	case 0xde:
		cpu.setAddr(3, cpu.hl(), mem) // SET 3 (HL) []
	case 0xdd:
		cpu.set(3, &cpu.l) // SET 3 L []
	case 0xe7:
		cpu.set(4, &cpu.a) // SET 4 A []
	case 0xe0:
		cpu.set(4, &cpu.b) // SET 4 B []
	case 0xe1:
		cpu.set(4, &cpu.c) // SET 4 C []
	case 0xe2:
		cpu.set(4, &cpu.d) // SET 4 D []
	case 0xe3:
		cpu.set(4, &cpu.e) // SET 4 E []
	case 0xe4:
		cpu.set(4, &cpu.h) // SET 4 H []
	case 0xe6:
		cpu.setAddr(4, cpu.hl(), mem) // SET 4 (HL) []
	case 0xe5:
		cpu.set(4, &cpu.l) // SET 4 L []
	case 0xef:
		cpu.set(5, &cpu.a) // SET 5 A []
	case 0xe8:
		cpu.set(5, &cpu.b) // SET 5 B []
	case 0xe9:
		cpu.set(5, &cpu.c) // SET 5 C []
	case 0xea:
		cpu.set(5, &cpu.d) // SET 5 D []
	case 0xeb:
		cpu.set(5, &cpu.e) // SET 5 E []
	case 0xec:
		cpu.set(5, &cpu.h) // SET 5 H []
	case 0xee:
		cpu.setAddr(5, cpu.hl(), mem) // SET 5 (HL) []
	case 0xed:
		cpu.set(5, &cpu.l) // SET 5 L []
	case 0xf7:
		cpu.set(6, &cpu.a) // SET 6 A []
	case 0xf0:
		cpu.set(6, &cpu.b) // SET 6 B []
	case 0xf1:
		cpu.set(6, &cpu.c) // SET 6 C []
	case 0xf2:
		cpu.set(6, &cpu.d) // SET 6 D []
	case 0xf3:
		cpu.set(6, &cpu.e) // SET 6 E []
	case 0xf4:
		cpu.set(6, &cpu.h) // SET 6 H []
	case 0xf6:
		cpu.setAddr(6, cpu.hl(), mem) // SET 6 (HL) []
	case 0xf5:
		cpu.set(6, &cpu.l) // SET 6 L []
	case 0xff:
		cpu.set(7, &cpu.a) // SET 7 A []
	case 0xf8:
		cpu.set(7, &cpu.b) // SET 7 B []
	case 0xf9:
		cpu.set(7, &cpu.c) // SET 7 C []
	case 0xfa:
		cpu.set(7, &cpu.d) // SET 7 D []
	case 0xfb:
		cpu.set(7, &cpu.e) // SET 7 E []
	case 0xfc:
		cpu.set(7, &cpu.h) // SET 7 H []
	case 0xfe:
		cpu.setAddr(7, cpu.hl(), mem) // SET 7 (HL) []
	case 0xfd:
		cpu.set(7, &cpu.l) // SET 7 L []
	case 0x27:
		cpu.sla(&cpu.a) // SLA A  [Z 0 0 C]
	case 0x20:
		cpu.sla(&cpu.b) // SLA B  [Z 0 0 C]
	case 0x21:
		cpu.sla(&cpu.c) // SLA C  [Z 0 0 C]
	case 0x22:
		cpu.sla(&cpu.d) // SLA D  [Z 0 0 C]
	case 0x23:
		cpu.sla(&cpu.e) // SLA E  [Z 0 0 C]
	case 0x24:
		cpu.sla(&cpu.h) // SLA H  [Z 0 0 C]
	case 0x26:
		cpu.slaAddr(cpu.hl(), mem) // SLA (HL)  [Z 0 0 C]
	case 0x25:
		cpu.sla(&cpu.l) // SLA L  [Z 0 0 C]
	case 0x2f:
		cpu.sra(&cpu.a) // SRA A  [Z 0 0 C]
	case 0x28:
		cpu.sra(&cpu.b) // SRA B  [Z 0 0 C]
	case 0x29:
		cpu.sra(&cpu.c) // SRA C  [Z 0 0 C]
	case 0x2a:
		cpu.sra(&cpu.d) // SRA D  [Z 0 0 C]
	case 0x2b:
		cpu.sra(&cpu.e) // SRA E  [Z 0 0 C]
	case 0x2c:
		cpu.sra(&cpu.h) // SRA H  [Z 0 0 C]
	case 0x2e:
		cpu.sraAddr(cpu.hl(), mem) // SRA (HL)  [Z 0 0 C]
	case 0x2d:
		cpu.sra(&cpu.l) // SRA L  [Z 0 0 C]
	case 0x3f:
		cpu.srl(&cpu.a) // SRL A  [Z 0 0 C]
	case 0x38:
		cpu.srl(&cpu.b) // SRL B  [Z 0 0 C]
	case 0x39:
		cpu.srl(&cpu.c) // SRL C  [Z 0 0 C]
	case 0x3a:
		cpu.srl(&cpu.d) // SRL D  [Z 0 0 C]
	case 0x3b:
		cpu.srl(&cpu.e) // SRL E  [Z 0 0 C]
	case 0x3c:
		cpu.srl(&cpu.h) // SRL H  [Z 0 0 C]
	case 0x3e:
		cpu.srlAddr(cpu.hl(), mem) // SRL (HL)  [Z 0 0 C]
	case 0x3d:
		cpu.srl(&cpu.l) // SRL L  [Z 0 0 C]
	case 0x37:
		cpu.swap(&cpu.a) // SWAP A  [Z 0 0 0]
	case 0x30:
		cpu.swap(&cpu.b) // SWAP B  [Z 0 0 0]
	case 0x31:
		cpu.swap(&cpu.c) // SWAP C  [Z 0 0 0]
	case 0x32:
		cpu.swap(&cpu.d) // SWAP D  [Z 0 0 0]
	case 0x33:
		cpu.swap(&cpu.e) // SWAP E  [Z 0 0 0]
	case 0x34:
		cpu.swap(&cpu.h) // SWAP H  [Z 0 0 0]
	case 0x36:
		cpu.swapAddr(cpu.hl(), mem) // SWAP (HL)  [Z 0 0 0]
	case 0x35:
		cpu.swap(&cpu.l) // SWAP L  [Z 0 0 0]
	default:
		panic(fmt.Sprintf("Missing dispatchPrefixedInstruction: %cb02x", instruction))
	}
}
