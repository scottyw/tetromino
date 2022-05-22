package cpu

import (
	"github.com/scottyw/tetromino/gameboy/memory"
)

//
// Refactor FIXME
//
// It is not necessary to pass mapper as an argument to many of these methods any more because the cpu holds a direct reference to the mapper itself
//
// This is a vestigial tail from the major refactor of core execution
//

func (cpu *CPU) adcM(mapper *memory.Mapper) func() {
	return func() {
		cpu.adc(mapper.Read(cpu.hl()))
	}
}

func (cpu *CPU) adcA() { cpu.adc(cpu.a) }

func (cpu *CPU) adcB() { cpu.adc(cpu.b) }

func (cpu *CPU) adcC() { cpu.adc(cpu.c) }

func (cpu *CPU) adcD() { cpu.adc(cpu.d) }

func (cpu *CPU) adcE() { cpu.adc(cpu.e) }

func (cpu *CPU) adcH() { cpu.adc(cpu.h) }

func (cpu *CPU) adcL() { cpu.adc(cpu.l) }

func (cpu *CPU) adcU() { cpu.adc(cpu.u8a) }

func (cpu *CPU) adc(u8 uint8) {
	a := cpu.a
	cpu.a += u8
	hf := hc8(a, u8)
	cf := c8(a, u8)
	if cpu.cf() {
		cpu.a++
		hf = hf || cpu.a&0x0f == 0
		cf = cf || cpu.a == 0
	}
	// [Z 0 H C]
	cpu.setZf(cpu.a == 0)
	cpu.setNf(false)
	cpu.setHf(hf)
	cpu.setCf(cf)
}

func (cpu *CPU) addM(mapper *memory.Mapper) func() {
	return func() {
		cpu.add(mapper.Read(cpu.hl()))
	}
}

func (cpu *CPU) addA() { cpu.add(cpu.a) }

func (cpu *CPU) addB() { cpu.add(cpu.b) }

func (cpu *CPU) addC() { cpu.add(cpu.c) }

func (cpu *CPU) addD() { cpu.add(cpu.d) }

func (cpu *CPU) addE() { cpu.add(cpu.e) }

func (cpu *CPU) addH() { cpu.add(cpu.h) }

func (cpu *CPU) addL() { cpu.add(cpu.l) }

func (cpu *CPU) addU() { cpu.add(cpu.u8a) }

func (cpu *CPU) add(u8 uint8) {
	a := cpu.a
	cpu.a += u8
	// [Z 0 H C]
	cpu.setZf(cpu.a == 0)
	cpu.setNf(false)
	cpu.setHf(hc8(a, u8))
	cpu.setCf(c8(a, u8))
}

func (cpu *CPU) addHLBC() { cpu.addHL(cpu.bc()) }

func (cpu *CPU) addHLDE() { cpu.addHL(cpu.de()) }

func (cpu *CPU) addHLHL() { cpu.addHL(cpu.hl()) }

func (cpu *CPU) addHLSP() { cpu.addHL(cpu.sp) }

func (cpu *CPU) addHL(u16 uint16) {
	hl := cpu.hl()
	new := hl + u16
	cpu.h = uint8(new >> 8)
	cpu.l = uint8(new)
	// [- 0 H C]
	cpu.setNf(false)
	cpu.setHf(hc16(hl, u16))
	cpu.setCf(c16(hl, u16))
}

func (cpu *CPU) addSP() {
	i8 := int8(cpu.u8a)
	sp := cpu.sp
	cpu.sp = uint16(int(cpu.sp) + int(i8))
	// [0 0 H C]
	cpu.setZf(false)
	cpu.setNf(false)
	if i8 >= 0 {
		cpu.setHf(int(sp)&0x0f+int(i8)&0x0f > 0x0f)
		cpu.setCf(int(sp)&0xff+int(i8) > 0xff)
	} else {
		cpu.setHf(int(sp)&0x0f >= int(cpu.sp)&0x0f)
		cpu.setCf(int(sp)&0xff >= int(cpu.sp)&0xff)
	}
}

func (cpu *CPU) andM(mapper *memory.Mapper) func() {
	return func() {
		cpu.and(mapper.Read(cpu.hl()))
	}
}

func (cpu *CPU) andA() { cpu.and(cpu.a) }

func (cpu *CPU) andB() { cpu.and(cpu.b) }

func (cpu *CPU) andC() { cpu.and(cpu.c) }

func (cpu *CPU) andD() { cpu.and(cpu.d) }

func (cpu *CPU) andE() { cpu.and(cpu.e) }

func (cpu *CPU) andH() { cpu.and(cpu.h) }

func (cpu *CPU) andL() { cpu.and(cpu.l) }

func (cpu *CPU) andU() { cpu.and(cpu.u8a) }

func (cpu *CPU) and(u8 uint8) {
	cpu.a &= u8
	// [Z 0 1 0]
	cpu.setZf(cpu.a == 0)
	cpu.setNf(false)
	cpu.setHf(true)
	cpu.setCf(false)
}

func (cpu *CPU) bitM(mapper *memory.Mapper, pos uint8) func() {
	return func() {
		u8 := mapper.Read(cpu.hl())
		cpu.bit(pos, &u8)()
	}
}

func (cpu *CPU) bit(pos uint8, r8 *uint8) func() {
	return func() {
		zero := *r8&bits[pos] == 0
		// [Z 0 1 -]
		cpu.setZf(zero)
		cpu.setNf(false)
		cpu.setHf(true)
	}
}

func (cpu *CPU) call() {
	// Store the old PC to write to memory in the next steps
	cpu.m8a = uint8(cpu.pc & 0xff)
	cpu.m8b = uint8(cpu.pc >> 8)
	// Update the PC
	cpu.pc = cpu.u16()
}

func (cpu *CPU) ccf() {
	// [- 0 0 C]
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(!cpu.cf())
}

func (cpu *CPU) cpM(mapper *memory.Mapper) func() {
	return func() {
		cpu.cp(mapper.Read(cpu.hl()))
	}
}

func (cpu *CPU) cpA() { cpu.cp(cpu.a) }

func (cpu *CPU) cpB() { cpu.cp(cpu.b) }

func (cpu *CPU) cpC() { cpu.cp(cpu.c) }

func (cpu *CPU) cpD() { cpu.cp(cpu.d) }

func (cpu *CPU) cpE() { cpu.cp(cpu.e) }

func (cpu *CPU) cpH() { cpu.cp(cpu.h) }

func (cpu *CPU) cpL() { cpu.cp(cpu.l) }

func (cpu *CPU) cpU() { cpu.cp(cpu.u8a) }

func (cpu *CPU) cp(u8 uint8) {
	// [Z 1 H C]
	cpu.setZf(cpu.a == u8)
	cpu.setNf(true)
	cpu.setHf(hc8Sub(cpu.a, u8))
	cpu.setCf(c8Sub(cpu.a, u8))
}

func (cpu *CPU) cpl() {
	cpu.a = ^cpu.a
	// [- 1 1 -]
	cpu.setNf(true)
	cpu.setHf(true)
}

func (cpu *CPU) daa() {
	if cpu.nf() {
		if cpu.hf() {
			cpu.a -= 0x06
		}
		if cpu.cf() {
			cpu.a -= 0x60
		}
	} else {
		a := cpu.a
		if cpu.hf() || cpu.a&0x0f > 0x09 {
			cpu.a += 0x06
		}
		if cpu.cf() || a&0xf0 > 0x90 || cpu.a&0xf0 > 0x90 {
			cpu.a += 0x60
			cpu.setCf(true)
		}
	}
	// [Z - 0 C]
	cpu.setZf(cpu.a == 0)
	cpu.setHf(false)
}

func (cpu *CPU) decM(mapper *memory.Mapper) func() {
	return func() {
		cpu.dec(&cpu.m8a)
		mapper.Write(cpu.hl(), cpu.m8a)
	}
}

func (cpu *CPU) decA() { cpu.dec(&cpu.a) }

func (cpu *CPU) decB() { cpu.dec(&cpu.b) }

func (cpu *CPU) decC() { cpu.dec(&cpu.c) }

func (cpu *CPU) decD() { cpu.dec(&cpu.d) }

func (cpu *CPU) decE() { cpu.dec(&cpu.e) }

func (cpu *CPU) decH() { cpu.dec(&cpu.h) }

func (cpu *CPU) decL() { cpu.dec(&cpu.l) }

func (cpu *CPU) dec(r8 *uint8) {
	old := *r8
	*r8--
	// [Z 1 H -]
	cpu.setZf(*r8 == 0)
	cpu.setNf(true)
	cpu.setHf(hc8Sub(old, 1))
}

func (cpu *CPU) decBC() {
	cpu.oam.TriggerWriteCorruption(cpu.bc())
	cpu.dec16(&cpu.b, &cpu.c)
}

func (cpu *CPU) decDE() {
	cpu.oam.TriggerWriteCorruption(cpu.de())
	cpu.dec16(&cpu.d, &cpu.e)
}

func (cpu *CPU) decHL() {
	cpu.oam.TriggerWriteCorruption(cpu.hl())
	cpu.dec16(&cpu.h, &cpu.l)
}

func (cpu *CPU) dec16(msb, lsb *uint8) {
	new := uint16(*msb)<<8 + uint16(*lsb) - 1
	*msb = uint8(new >> 8)
	*lsb = uint8(new)
}

func (cpu *CPU) decSP() {
	cpu.oam.TriggerWriteCorruption(cpu.sp)
	cpu.sp--
}

func (cpu *CPU) di() {
	cpu.interrupts.Disable()
}

func (cpu *CPU) ei() {
	cpu.interrupts.Enable()
}

func (cpu *CPU) halt() func() {
	return func() {
		if cpu.interrupts.Enabled() {
			cpu.halted = true
		} else {
			if !cpu.interrupts.Pending() {
				cpu.halted = true
			} else {
				cpu.haltbug = true
			}
		}
	}
}

func (cpu *CPU) incM(mapper *memory.Mapper) func() {
	return func() {
		cpu.inc(&cpu.m8a)
		mapper.Write(cpu.hl(), cpu.m8a)
	}
}

func (cpu *CPU) incA() { cpu.inc(&cpu.a) }

func (cpu *CPU) incB() { cpu.inc(&cpu.b) }

func (cpu *CPU) incC() { cpu.inc(&cpu.c) }

func (cpu *CPU) incD() { cpu.inc(&cpu.d) }

func (cpu *CPU) incE() { cpu.inc(&cpu.e) }

func (cpu *CPU) incH() { cpu.inc(&cpu.h) }

func (cpu *CPU) incL() { cpu.inc(&cpu.l) }

func (cpu *CPU) inc(r8 *uint8) {
	old := *r8
	*r8++
	// [Z 0 H -]
	cpu.setZf(*r8 == 0)
	cpu.setNf(false)
	cpu.setHf(hc8(old, 1))
}

func (cpu *CPU) incBC() {
	cpu.oam.TriggerWriteCorruption(cpu.bc())
	cpu.inc16(&cpu.b, &cpu.c)
}

func (cpu *CPU) incDE() {
	cpu.oam.TriggerWriteCorruption(cpu.de())
	cpu.inc16(&cpu.d, &cpu.e)
}

func (cpu *CPU) incHL() {
	cpu.oam.TriggerWriteCorruption(cpu.hl())
	cpu.inc16(&cpu.h, &cpu.l)
}

func (cpu *CPU) inc16(msb, lsb *uint8) {
	new := uint16(*msb)<<8 + uint16(*lsb) + 1
	*msb = uint8(new >> 8)
	*lsb = uint8(new)
}

func (cpu *CPU) incSP() {
	cpu.oam.TriggerWriteCorruption(cpu.sp)
	cpu.sp++
}

func (cpu *CPU) jp() {
	cpu.pc = cpu.u16()
}

func (cpu *CPU) jpHL() {
	cpu.pc = cpu.hl()
}

func (cpu *CPU) jr() {
	i8 := int8(cpu.u8a)
	cpu.pc = uint16(int16(cpu.pc) + int16(i8))
}

func (cpu *CPU) ldAB() { cpu.a = cpu.b }

func (cpu *CPU) ldAC() { cpu.a = cpu.c }

func (cpu *CPU) ldAD() { cpu.a = cpu.d }

func (cpu *CPU) ldAE() { cpu.a = cpu.e }

func (cpu *CPU) ldAH() { cpu.a = cpu.h }

func (cpu *CPU) ldAL() { cpu.a = cpu.l }

func (cpu *CPU) ldAU() { cpu.a = cpu.u8a }

func (cpu *CPU) ldBA() { cpu.b = cpu.a }

func (cpu *CPU) ldBC() { cpu.b = cpu.c }

func (cpu *CPU) ldBD() { cpu.b = cpu.d }

func (cpu *CPU) ldBE() { cpu.b = cpu.e }

func (cpu *CPU) ldBH() { cpu.b = cpu.h }

func (cpu *CPU) ldBL() { cpu.b = cpu.l }

func (cpu *CPU) ldBU() { cpu.b = cpu.u8a }

func (cpu *CPU) ldCA() { cpu.c = cpu.a }

func (cpu *CPU) ldCB() { cpu.c = cpu.b }

func (cpu *CPU) ldCD() { cpu.c = cpu.d }

func (cpu *CPU) ldCE() { cpu.c = cpu.e }

func (cpu *CPU) ldCH() { cpu.c = cpu.h }

func (cpu *CPU) ldCL() { cpu.c = cpu.l }

func (cpu *CPU) ldCU() { cpu.c = cpu.u8a }

func (cpu *CPU) ldDA() { cpu.d = cpu.a }

func (cpu *CPU) ldDB() { cpu.d = cpu.b }

func (cpu *CPU) ldDC() { cpu.d = cpu.c }

func (cpu *CPU) ldDE() { cpu.d = cpu.e }

func (cpu *CPU) ldDH() { cpu.d = cpu.h }

func (cpu *CPU) ldDL() { cpu.d = cpu.l }

func (cpu *CPU) ldDU() { cpu.d = cpu.u8a }

func (cpu *CPU) ldEA() { cpu.e = cpu.a }

func (cpu *CPU) ldEB() { cpu.e = cpu.b }

func (cpu *CPU) ldEC() { cpu.e = cpu.c }

func (cpu *CPU) ldED() { cpu.e = cpu.d }

func (cpu *CPU) ldEH() { cpu.e = cpu.h }

func (cpu *CPU) ldEL() { cpu.e = cpu.l }

func (cpu *CPU) ldEU() { cpu.e = cpu.u8a }

func (cpu *CPU) ldHA() { cpu.h = cpu.a }

func (cpu *CPU) ldHB() { cpu.h = cpu.b }

func (cpu *CPU) ldHC() { cpu.h = cpu.c }

func (cpu *CPU) ldHD() { cpu.h = cpu.d }

func (cpu *CPU) ldHE() { cpu.h = cpu.e }

func (cpu *CPU) ldHL() { cpu.h = cpu.l }

func (cpu *CPU) ldHU() { cpu.h = cpu.u8a }

func (cpu *CPU) ldLA() { cpu.l = cpu.a }

func (cpu *CPU) ldLB() { cpu.l = cpu.b }

func (cpu *CPU) ldLC() { cpu.l = cpu.c }

func (cpu *CPU) ldLD() { cpu.l = cpu.d }

func (cpu *CPU) ldLE() { cpu.l = cpu.e }

func (cpu *CPU) ldLH() { cpu.l = cpu.h }

func (cpu *CPU) ldLU() { cpu.l = cpu.u8a }

func (cpu *CPU) ldHLU8(mapper *memory.Mapper) func() {
	return func() {
		mapper.Write(cpu.hl(), cpu.u8a)
	}
}

func (cpu *CPU) ldBCU16() {
	u16 := cpu.u16()
	cpu.b = uint8(u16 >> 8)
	cpu.c = uint8(u16)
}

func (cpu *CPU) ldDEU16() {
	u16 := cpu.u16()
	cpu.d = uint8(u16 >> 8)
	cpu.e = uint8(u16)
}

func (cpu *CPU) ldHLU16() {
	u16 := cpu.u16()
	cpu.h = uint8(u16 >> 8)
	cpu.l = uint8(u16)
}

func (cpu *CPU) ldSPU16() {
	cpu.sp = cpu.u16()
}

func (cpu *CPU) ldSPHL() {
	cpu.sp = cpu.hl()
}

func (cpu *CPU) ldHLSP() {
	i8 := int8(cpu.u8a)
	new := int(int(cpu.sp) + int(i8))
	cpu.h = uint8(new >> 8)
	cpu.l = uint8(new)
	// [0 0 H C]
	cpu.setZf(false)
	cpu.setNf(false)
	if i8 >= 0 {
		cpu.setHf(int(cpu.sp)&0x0f+int(i8)&0x0f > 0x0f)
		cpu.setCf(int(cpu.sp)&0xff+int(i8) > 0xff)
	} else {
		cpu.setHf(int(cpu.sp)&0x0f >= new&0x0f)
		cpu.setCf(int(cpu.sp)&0xff >= new&0xff)
	}
}

func (cpu *CPU) ldBCA(mapper *memory.Mapper) func() {
	return func() {
		mapper.Write(cpu.bc(), cpu.a)
	}
}

func (cpu *CPU) ldABC(mapper *memory.Mapper) func() {
	return func() {
		cpu.a = mapper.Read(cpu.bc())
	}
}

func (cpu *CPU) ldDEA(mapper *memory.Mapper) func() {
	return func() {
		mapper.Write(cpu.de(), cpu.a)
	}
}

func (cpu *CPU) ldADE(mapper *memory.Mapper) func() {
	return func() {
		cpu.a = mapper.Read(cpu.de())
	}
}

func (cpu *CPU) ldHLDA(mapper *memory.Mapper) func() {
	return func() {
		mapper.Write(cpu.hl(), cpu.a)
		cpu.decHL()
	}
}

func (cpu *CPU) ldAHLD(mapper *memory.Mapper) func() {
	return func() {
		cpu.a = mapper.Read(cpu.hl())
		cpu.decHL()
	}
}

func (cpu *CPU) ldHLIA(mapper *memory.Mapper) func() {
	return func() {
		mapper.Write(cpu.hl(), cpu.a)
		cpu.incHL()
	}
}

func (cpu *CPU) ldAHLI(mapper *memory.Mapper) func() {
	return func() {
		cpu.a = mapper.Read(cpu.hl())
		cpu.incHL()
	}
}

func (cpu *CPU) ldHLA(mapper *memory.Mapper) func() { return func() { mapper.Write(cpu.hl(), cpu.a) } }

func (cpu *CPU) ldHLB(mapper *memory.Mapper) func() { return func() { mapper.Write(cpu.hl(), cpu.b) } }

func (cpu *CPU) ldHLC(mapper *memory.Mapper) func() { return func() { mapper.Write(cpu.hl(), cpu.c) } }

func (cpu *CPU) ldHLD(mapper *memory.Mapper) func() { return func() { mapper.Write(cpu.hl(), cpu.d) } }

func (cpu *CPU) ldHLE(mapper *memory.Mapper) func() { return func() { mapper.Write(cpu.hl(), cpu.e) } }

func (cpu *CPU) ldHLH(mapper *memory.Mapper) func() { return func() { mapper.Write(cpu.hl(), cpu.h) } }

func (cpu *CPU) ldHLL(mapper *memory.Mapper) func() { return func() { mapper.Write(cpu.hl(), cpu.l) } }

func (cpu *CPU) ldAHL(mapper *memory.Mapper) func() { return func() { cpu.a = mapper.Read(cpu.hl()) } }

func (cpu *CPU) ldBHL(mapper *memory.Mapper) func() { return func() { cpu.b = mapper.Read(cpu.hl()) } }

func (cpu *CPU) ldCHL(mapper *memory.Mapper) func() { return func() { cpu.c = mapper.Read(cpu.hl()) } }

func (cpu *CPU) ldDHL(mapper *memory.Mapper) func() { return func() { cpu.d = mapper.Read(cpu.hl()) } }

func (cpu *CPU) ldEHL(mapper *memory.Mapper) func() { return func() { cpu.e = mapper.Read(cpu.hl()) } }

func (cpu *CPU) ldHHL(mapper *memory.Mapper) func() { return func() { cpu.h = mapper.Read(cpu.hl()) } }

func (cpu *CPU) ldLHL(mapper *memory.Mapper) func() { return func() { cpu.l = mapper.Read(cpu.hl()) } }

func (cpu *CPU) ldMHL(mapper *memory.Mapper) func() {
	return func() { cpu.m8a = mapper.Read(cpu.hl()) }
}

func (cpu *CPU) ldACX(mapper *memory.Mapper) func() {
	return func() {
		a16 := uint16(0xff00 + uint16(cpu.c))
		cpu.a = mapper.Read(a16)
	}
}

func (cpu *CPU) ldCXA(mapper *memory.Mapper) func() {
	return func() {
		a16 := uint16(0xff00 + uint16(cpu.c))
		mapper.Write(a16, cpu.a)
	}
}

func (cpu *CPU) ldAUX(mapper *memory.Mapper) func() {
	return func() {
		a16 := uint16(0xff00 + uint16(cpu.u8a))
		cpu.a = mapper.Read(a16)
	}
}

func (cpu *CPU) ldUXA(mapper *memory.Mapper) func() {
	return func() {
		a16 := uint16(0xff00 + uint16(cpu.u8a))
		mapper.Write(a16, cpu.a)
	}
}

func (cpu *CPU) ldAUX16(mapper *memory.Mapper) func() {
	return func() {
		cpu.a = mapper.Read(cpu.u16())
	}
}

func (cpu *CPU) ldUX16A(mapper *memory.Mapper) func() {
	return func() {
		mapper.Write(cpu.u16(), cpu.a)
	}
}

func (cpu *CPU) writeLowSP(mapper *memory.Mapper) func() {
	return func() {
		mapper.Write(cpu.u16(), uint8(cpu.sp))
	}
}

func (cpu *CPU) writeHighSP(mapper *memory.Mapper) func() {
	return func() {
		mapper.Write(cpu.u16()+1, uint8(cpu.sp>>8))
	}
}

func (cpu *CPU) orM(mapper *memory.Mapper) func() {
	return func() {
		cpu.or(mapper.Read(cpu.hl()))
	}
}

func (cpu *CPU) orA() { cpu.or(cpu.a) }

func (cpu *CPU) orB() { cpu.or(cpu.b) }

func (cpu *CPU) orC() { cpu.or(cpu.c) }

func (cpu *CPU) orD() { cpu.or(cpu.d) }

func (cpu *CPU) orE() { cpu.or(cpu.e) }

func (cpu *CPU) orH() { cpu.or(cpu.h) }

func (cpu *CPU) orL() { cpu.or(cpu.l) }

func (cpu *CPU) orU() { cpu.or(cpu.u8a) }

func (cpu *CPU) or(u8 uint8) {
	cpu.a |= u8
	// [Z 0 0 0]
	cpu.setZf(cpu.a == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(false)
}

func (cpu *CPU) pop(mapper *memory.Mapper, r8 *uint8) func() {
	return func() {
		*r8 = mapper.Read(cpu.sp)
		cpu.incSP()
	}
}

func (cpu *CPU) popF(mapper *memory.Mapper) func() {
	return func() {
		// Lower nibble is always zero no matter what data was written
		cpu.f = mapper.Read(cpu.sp) & 0xf0
		cpu.oam.TriggerWriteCorruption(cpu.sp)
		cpu.sp++
	}
}

func (cpu *CPU) push(mapper *memory.Mapper, r8 *uint8) func() {
	return func() {
		cpu.decSP()
		mapper.Write(cpu.sp, *r8)
	}
}

func (cpu *CPU) resM(mapper *memory.Mapper, pos uint8) func() {
	return func() {
		cpu.res(pos, &cpu.m8a)()
		mapper.Write(cpu.hl(), cpu.m8a)
	}
}

func (cpu *CPU) res(pos uint8, r8 *uint8) func() {
	return func() {
		*r8 &^= bits[pos]
	}
}

func (cpu *CPU) ret() {
	cpu.pc = uint16(cpu.m8b)<<8 | uint16(cpu.m8a)
}

func (cpu *CPU) reti() {
	cpu.ret()
	cpu.ei()
}

func (cpu *CPU) rlM(mapper *memory.Mapper) func() {
	return func() {
		cpu.rl(&cpu.m8a)
		mapper.Write(cpu.hl(), cpu.m8a)
	}
}

func (cpu *CPU) rlA() { cpu.rl(&cpu.a) }

func (cpu *CPU) rlB() { cpu.rl(&cpu.b) }

func (cpu *CPU) rlC() { cpu.rl(&cpu.c) }

func (cpu *CPU) rlD() { cpu.rl(&cpu.d) }

func (cpu *CPU) rlE() { cpu.rl(&cpu.e) }

func (cpu *CPU) rlH() { cpu.rl(&cpu.h) }

func (cpu *CPU) rlL() { cpu.rl(&cpu.l) }

func (cpu *CPU) rl(r8 *uint8) {
	cf := (*r8 & 0x80) > 0
	*r8 <<= 1
	if cpu.cf() {
		*r8 |= 0x01
	}
	// [Z 0 0 C]
	cpu.setZf(*r8 == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(cf)
}

func (cpu *CPU) rla() {
	cpu.rl(&cpu.a)
	cpu.f &^= zFlag
	// [0 0 0 C]
	cpu.setZf(false)
}

func (cpu *CPU) rlcM(mapper *memory.Mapper) func() {
	return func() {
		cpu.rlc(&cpu.m8a)
		mapper.Write(cpu.hl(), cpu.m8a)
	}
}

func (cpu *CPU) rlcA() { cpu.rlc(&cpu.a) }

func (cpu *CPU) rlcB() { cpu.rlc(&cpu.b) }

func (cpu *CPU) rlcC() { cpu.rlc(&cpu.c) }

func (cpu *CPU) rlcD() { cpu.rlc(&cpu.d) }

func (cpu *CPU) rlcE() { cpu.rlc(&cpu.e) }

func (cpu *CPU) rlcH() { cpu.rlc(&cpu.h) }

func (cpu *CPU) rlcL() { cpu.rlc(&cpu.l) }

func (cpu *CPU) rlc(r8 *uint8) {
	cf := (*r8 & 0x80) > 0
	*r8 <<= 1
	if cf {
		*r8 |= 0x01
	}
	// [Z 0 0 C]
	cpu.setZf(*r8 == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(cf)
}

func (cpu *CPU) rlca() {
	cpu.rlc(&cpu.a)
	cpu.f &^= zFlag
	// [0 0 0 C]
	cpu.setZf(false)
}

func (cpu *CPU) rrM(mapper *memory.Mapper) func() {
	return func() {
		cpu.rr(&cpu.m8a)
		mapper.Write(cpu.hl(), cpu.m8a)
	}
}

func (cpu *CPU) rrA() { cpu.rr(&cpu.a) }

func (cpu *CPU) rrB() { cpu.rr(&cpu.b) }

func (cpu *CPU) rrC() { cpu.rr(&cpu.c) }

func (cpu *CPU) rrD() { cpu.rr(&cpu.d) }

func (cpu *CPU) rrE() { cpu.rr(&cpu.e) }

func (cpu *CPU) rrH() { cpu.rr(&cpu.h) }

func (cpu *CPU) rrL() { cpu.rr(&cpu.l) }

func (cpu *CPU) rr(r8 *uint8) {
	cf := (*r8 & 0x01) > 0
	*r8 >>= 1
	if cpu.cf() {
		*r8 |= 0x80
	}
	// [Z 0 0 C]
	cpu.setZf(*r8 == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(cf)
}

func (cpu *CPU) rra() {
	cpu.rr(&cpu.a)
	// [0 0 0 C]
	cpu.setZf(false)
}

func (cpu *CPU) rrcM(mapper *memory.Mapper) func() {
	return func() {
		cpu.rrc(&cpu.m8a)
		mapper.Write(cpu.hl(), cpu.m8a)
	}
}

func (cpu *CPU) rrcA() { cpu.rrc(&cpu.a) }

func (cpu *CPU) rrcB() { cpu.rrc(&cpu.b) }

func (cpu *CPU) rrcC() { cpu.rrc(&cpu.c) }

func (cpu *CPU) rrcD() { cpu.rrc(&cpu.d) }

func (cpu *CPU) rrcE() { cpu.rrc(&cpu.e) }

func (cpu *CPU) rrcH() { cpu.rrc(&cpu.h) }

func (cpu *CPU) rrcL() { cpu.rrc(&cpu.l) }

func (cpu *CPU) rrc(r8 *uint8) {
	cf := (*r8 & 0x01) > 0
	*r8 >>= 1
	if cf {
		*r8 |= 0x80
	}
	// [Z 0 0 C]
	cpu.setZf(*r8 == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(cf)
}

func (cpu *CPU) rrca() {
	cpu.rrc(&cpu.a)
	// [0 0 0 C]
	cpu.setZf(false)
}

func (cpu *CPU) rst(a16 uint16) func() {
	return func() {
		// Store the old PC to write to memory in the next steps
		cpu.m8a = uint8(cpu.pc & 0xff)
		cpu.m8b = uint8(cpu.pc >> 8)
		// Update the PC
		cpu.pc = a16
	}
}

func (cpu *CPU) setM(mapper *memory.Mapper, pos uint8) func() {
	return func() {
		cpu.set(pos, &cpu.m8a)()
		mapper.Write(cpu.hl(), cpu.m8a)
	}
}

func (cpu *CPU) set(pos uint8, r8 *uint8) func() {
	return func() {
		*r8 |= bits[pos]
	}
}

func (cpu *CPU) slaM(mapper *memory.Mapper) func() {
	return func() {
		cpu.sla(&cpu.m8a)
		mapper.Write(cpu.hl(), cpu.m8a)
	}
}

func (cpu *CPU) slaA() { cpu.sla(&cpu.a) }

func (cpu *CPU) slaB() { cpu.sla(&cpu.b) }

func (cpu *CPU) slaC() { cpu.sla(&cpu.c) }

func (cpu *CPU) slaD() { cpu.sla(&cpu.d) }

func (cpu *CPU) slaE() { cpu.sla(&cpu.e) }

func (cpu *CPU) slaH() { cpu.sla(&cpu.h) }

func (cpu *CPU) slaL() { cpu.sla(&cpu.l) }

func (cpu *CPU) sla(r8 *uint8) {
	cf := (*r8 & 0x80) > 0
	*r8 <<= 1
	// [Z 0 0 C]
	cpu.setZf(*r8 == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(cf)
}

func (cpu *CPU) sraM(mapper *memory.Mapper) func() {
	return func() {
		cpu.sra(&cpu.m8a)
		mapper.Write(cpu.hl(), cpu.m8a)
	}
}

func (cpu *CPU) sraA() { cpu.sra(&cpu.a) }

func (cpu *CPU) sraB() { cpu.sra(&cpu.b) }

func (cpu *CPU) sraC() { cpu.sra(&cpu.c) }

func (cpu *CPU) sraD() { cpu.sra(&cpu.d) }

func (cpu *CPU) sraE() { cpu.sra(&cpu.e) }

func (cpu *CPU) sraH() { cpu.sra(&cpu.h) }

func (cpu *CPU) sraL() { cpu.sra(&cpu.l) }

func (cpu *CPU) sra(r8 *uint8) {
	cf := (*r8 & 0x01) > 0
	bit7 := (*r8 & 0x80) > 0
	*r8 >>= 1
	if bit7 {
		*r8 |= 0x80
	}
	// [Z 0 0 C]
	cpu.setZf(*r8 == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(cf)
}

func (cpu *CPU) srlM(mapper *memory.Mapper) func() {
	return func() {
		cpu.srl(&cpu.m8a)
		mapper.Write(cpu.hl(), cpu.m8a)
	}
}

func (cpu *CPU) srlA() { cpu.srl(&cpu.a) }

func (cpu *CPU) srlB() { cpu.srl(&cpu.b) }

func (cpu *CPU) srlC() { cpu.srl(&cpu.c) }

func (cpu *CPU) srlD() { cpu.srl(&cpu.d) }

func (cpu *CPU) srlE() { cpu.srl(&cpu.e) }

func (cpu *CPU) srlH() { cpu.srl(&cpu.h) }

func (cpu *CPU) srlL() { cpu.srl(&cpu.l) }

func (cpu *CPU) srl(r8 *uint8) {
	cf := (*r8 & 0x01) > 0
	*r8 >>= 1
	// [Z 0 0 C]
	cpu.setZf(*r8 == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(cf)
}

func (cpu *CPU) swapM(mapper *memory.Mapper) func() {
	return func() {
		cpu.swap(&cpu.m8a)
		mapper.Write(cpu.hl(), cpu.m8a)
	}
}

func (cpu *CPU) swapA() { cpu.swap(&cpu.a) }

func (cpu *CPU) swapB() { cpu.swap(&cpu.b) }

func (cpu *CPU) swapC() { cpu.swap(&cpu.c) }

func (cpu *CPU) swapD() { cpu.swap(&cpu.d) }

func (cpu *CPU) swapE() { cpu.swap(&cpu.e) }

func (cpu *CPU) swapH() { cpu.swap(&cpu.h) }

func (cpu *CPU) swapL() { cpu.swap(&cpu.l) }

func (cpu *CPU) swap(r8 *uint8) {
	u8 := *r8
	*r8 = u8<<4 | u8>>4
	// [Z 0 0 0]
	cpu.setZf(*r8 == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(false)
}

func (cpu *CPU) scf() {
	// [- 0 0 1]
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(true)
}

func (cpu *CPU) stop() {
	cpu.stopped = true
}

func (cpu *CPU) sbcM(mapper *memory.Mapper) func() {
	return func() {
		cpu.sbc(mapper.Read(cpu.hl()))
	}
}

func (cpu *CPU) sbcA() { cpu.sbc(cpu.a) }

func (cpu *CPU) sbcB() { cpu.sbc(cpu.b) }

func (cpu *CPU) sbcC() { cpu.sbc(cpu.c) }

func (cpu *CPU) sbcD() { cpu.sbc(cpu.d) }

func (cpu *CPU) sbcE() { cpu.sbc(cpu.e) }

func (cpu *CPU) sbcH() { cpu.sbc(cpu.h) }

func (cpu *CPU) sbcL() { cpu.sbc(cpu.l) }

func (cpu *CPU) sbcU() { cpu.sbc(cpu.u8a) }

func (cpu *CPU) sbc(u8 uint8) {
	a := cpu.a
	cpu.a -= u8
	hf := hc8Sub(a, u8)
	cf := c8Sub(a, u8)
	if cpu.cf() {
		cpu.a--
		hf = hf || cpu.a&0x0f == 0x0f
		cf = cf || cpu.a == 0xff
	}
	// [Z 1 H C]
	cpu.setZf(cpu.a == 0)
	cpu.setNf(true)
	cpu.setHf(hf)
	cpu.setCf(cf)

}

func (cpu *CPU) subM(mapper *memory.Mapper) func() {
	return func() {
		cpu.sub(mapper.Read(cpu.hl()))
	}
}

func (cpu *CPU) subA() { cpu.sub(cpu.a) }

func (cpu *CPU) subB() { cpu.sub(cpu.b) }

func (cpu *CPU) subC() { cpu.sub(cpu.c) }

func (cpu *CPU) subD() { cpu.sub(cpu.d) }

func (cpu *CPU) subE() { cpu.sub(cpu.e) }

func (cpu *CPU) subH() { cpu.sub(cpu.h) }

func (cpu *CPU) subL() { cpu.sub(cpu.l) }

func (cpu *CPU) subU() { cpu.sub(cpu.u8a) }

func (cpu *CPU) sub(u8 uint8) {
	a := cpu.a
	cpu.a -= u8
	// [Z 1 H C]
	cpu.setZf(cpu.a == 0)
	cpu.setNf(true)
	cpu.setHf(hc8Sub(a, u8))
	cpu.setCf(c8Sub(a, u8))
}

func (cpu *CPU) xorM(mapper *memory.Mapper) func() {
	return func() {
		cpu.xor(mapper.Read(cpu.hl()))
	}
}

func (cpu *CPU) xorA() { cpu.xor(cpu.a) }

func (cpu *CPU) xorB() { cpu.xor(cpu.b) }

func (cpu *CPU) xorC() { cpu.xor(cpu.c) }

func (cpu *CPU) xorD() { cpu.xor(cpu.d) }

func (cpu *CPU) xorE() { cpu.xor(cpu.e) }

func (cpu *CPU) xorH() { cpu.xor(cpu.h) }

func (cpu *CPU) xorL() { cpu.xor(cpu.l) }

func (cpu *CPU) xorU() { cpu.xor(cpu.u8a) }

func (cpu *CPU) xor(u8 uint8) {
	cpu.a = cpu.a ^ u8
	// [Z 0 0 0]
	cpu.setZf(cpu.a == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(false)
}
