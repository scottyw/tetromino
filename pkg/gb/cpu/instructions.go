package cpu

import (
	"fmt"

	"github.com/scottyw/tetromino/pkg/gb/mem"
)

func (cpu *CPU) adcA() { cpu.adc(cpu.a) }

func (cpu *CPU) adcB() { cpu.adc(cpu.b) }

func (cpu *CPU) adcC() { cpu.adc(cpu.c) }

func (cpu *CPU) adcD() { cpu.adc(cpu.d) }

func (cpu *CPU) adcE() { cpu.adc(cpu.e) }

func (cpu *CPU) adcH() { cpu.adc(cpu.h) }

func (cpu *CPU) adcL() { cpu.adc(cpu.l) }

func (cpu *CPU) adcU8() { cpu.adc(cpu.u8a) }

func (cpu *CPU) adcM8() { cpu.adc(cpu.m8a) }

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

func (cpu *CPU) addA() { cpu.add(cpu.a) }

func (cpu *CPU) addB() { cpu.add(cpu.b) }

func (cpu *CPU) addC() { cpu.add(cpu.c) }

func (cpu *CPU) addD() { cpu.add(cpu.d) }

func (cpu *CPU) addE() { cpu.add(cpu.e) }

func (cpu *CPU) addH() { cpu.add(cpu.h) }

func (cpu *CPU) addL() { cpu.add(cpu.l) }

func (cpu *CPU) addU() { cpu.add(cpu.u8a) }

func (cpu *CPU) addM() { cpu.add(cpu.m8a) }

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

func (cpu *CPU) andA() { cpu.and(cpu.a) }

func (cpu *CPU) andB() { cpu.and(cpu.b) }

func (cpu *CPU) andC() { cpu.and(cpu.c) }

func (cpu *CPU) andD() { cpu.and(cpu.d) }

func (cpu *CPU) andE() { cpu.and(cpu.e) }

func (cpu *CPU) andH() { cpu.and(cpu.h) }

func (cpu *CPU) andL() { cpu.and(cpu.l) }

func (cpu *CPU) andU() { cpu.and(cpu.u8a) }

func (cpu *CPU) andM() { cpu.and(cpu.m8a) }

func (cpu *CPU) and(u8 uint8) {
	cpu.a &= u8
	// [Z 0 1 0]
	cpu.setZf(cpu.a == 0)
	cpu.setNf(false)
	cpu.setHf(true)
	cpu.setCf(false)
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

func (cpu *CPU) call(kind string, a16 uint16, mem *mem.Memory) {
	switch kind {
	case "":
		cpu.sp--
		mem.Write(cpu.sp, byte(cpu.pc>>8))
		cpu.sp--
		mem.Write(cpu.sp, byte(cpu.pc&0xff))
		cpu.pc = a16
	case "NZ":
		if !cpu.zf() {
			cpu.call("", a16, mem)
		} else {
			// cpu.altTicks = true
		}
	case "Z":
		if cpu.zf() {
			cpu.call("", a16, mem)
		} else {
			// cpu.altTicks = true
		}
	case "NC":
		if !cpu.cf() {
			cpu.call("", a16, mem)
		} else {
			// cpu.altTicks = true
		}
	case "C":
		if cpu.cf() {
			cpu.call("", a16, mem)
		} else {
			// cpu.altTicks = true
		}
	default:
		panic(fmt.Sprintf("No implementation for call: %v %v", kind, a16))
	}
}

func (cpu *CPU) ccf() {
	// [- 0 0 C]
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(!cpu.cf())
}

func (cpu *CPU) cp(u8 uint8) {
	// [Z 1 H C]
	cpu.setZf(cpu.a == u8)
	cpu.setNf(true)
	cpu.setHf(hc8Sub(cpu.a, u8))
	cpu.setCf(c8Sub(cpu.a, u8))
}

func (cpu *CPU) cpAddr(a16 uint16, mem *mem.Memory) {
	cpu.cp(mem.Read(a16))
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

func (cpu *CPU) dec(r8 *uint8) {
	old := *r8
	*r8--
	// [Z 1 H -]
	cpu.setZf(*r8 == 0)
	cpu.setNf(true)
	cpu.setHf(hc8Sub(old, 1))
}

func (cpu *CPU) dec16(msb, lsb *uint8) {
	new := uint16(*msb)<<8 + uint16(*lsb) - 1
	*msb = uint8(new >> 8)
	*lsb = uint8(new)
}

func (cpu *CPU) decSP() {
	cpu.sp--
}

func (cpu *CPU) decAddr(a16 uint16, mem *mem.Memory) {
	value := mem.Read(a16)
	cpu.dec(&value)
	mem.Write(a16, value)
}

func (cpu *CPU) di() {
	cpu.ime = false
}

func (cpu *CPU) ei() {
	cpu.ime = true
}

func (cpu *CPU) halt() {
	cpu.halted = true
	// FIXME halt bug needs to be implemented
}

func (cpu *CPU) inc(r8 *uint8) {
	old := *r8
	*r8++
	// [Z 0 H -]
	cpu.setZf(*r8 == 0)
	cpu.setNf(false)
	cpu.setHf(hc8(old, 1))
}

func (cpu *CPU) inc16(msb, lsb *uint8) {
	new := uint16(*msb)<<8 + uint16(*lsb) + 1
	*msb = uint8(new >> 8)
	*lsb = uint8(new)
}

func (cpu *CPU) incSP() {
	cpu.sp++
}

func (cpu *CPU) incAddr(a16 uint16, mem *mem.Memory) {
	value := mem.Read(a16)
	cpu.inc(&value)
	mem.Write(a16, value)
}

func (cpu *CPU) jp(kind string, a16 uint16) {
	switch kind {
	case "":
		cpu.pc = a16
	case "NZ":
		if !cpu.zf() {
			cpu.jp("", a16)
		} else {
			// cpu.altTicks = true
		}
	case "Z":
		if cpu.zf() {
			cpu.jp("", a16)
		} else {
			// cpu.altTicks = true
		}
	case "NC":
		if !cpu.cf() {
			cpu.jp("", a16)
		} else {
			// cpu.altTicks = true
		}
	case "C":
		if cpu.cf() {
			cpu.jp("", a16)
		} else {
			// cpu.altTicks = true
		}
	default:
		panic(fmt.Sprintf("No implementation for jp: %v %v", kind, a16))
	}
}

func (cpu *CPU) jr(kind string, i8 int8) {
	address := uint16(int16(cpu.pc) + int16(i8))
	cpu.jp(kind, address)
}

func (cpu *CPU) ld(r8 *uint8, u8 uint8) {
	*r8 = u8
}

func (cpu *CPU) ld16(msb, lsb *uint8, u16 uint16) {
	*msb = uint8(u16 >> 8)
	*lsb = uint8(u16)
}

func (cpu *CPU) ldR8A16(r8 *uint8, a16 uint16, mem *mem.Memory) {
	*r8 = mem.Read(a16)
}

func (cpu *CPU) ldA16U8(a16 uint16, u8 uint8, mem *mem.Memory) {
	mem.Write(a16, u8)
}

func (cpu *CPU) ldSP(u16 uint16) {
	cpu.sp = u16
}

func (cpu *CPU) ldSPHL() {
	cpu.sp = cpu.hl()
}

func (cpu *CPU) ldA16SP(a16 uint16, mem *mem.Memory) {
	mem.Write(a16, uint8(cpu.sp))
	mem.Write(a16+1, uint8(cpu.sp>>8))
}

func (cpu *CPU) ldHLSP(i8 int8) {
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

func (cpu *CPU) lddAA16(mem *mem.Memory) {
	cpu.ldR8A16(&cpu.a, cpu.hl(), mem)
	cpu.dec16(&cpu.h, &cpu.l)
}

func (cpu *CPU) lddA16A(mem *mem.Memory) {
	cpu.ldA16U8(cpu.hl(), cpu.a, mem)
	cpu.dec16(&cpu.h, &cpu.l)
}

func (cpu *CPU) ldiAA16(mem *mem.Memory) {
	cpu.ldR8A16(&cpu.a, cpu.hl(), mem)
	cpu.inc16(&cpu.h, &cpu.l)
}

func (cpu *CPU) ldiA16A(mem *mem.Memory) {
	cpu.ldA16U8(cpu.hl(), cpu.a, mem)
	cpu.inc16(&cpu.h, &cpu.l)
}

func (cpu *CPU) or(u8 uint8) {
	cpu.a |= u8
	// [Z 0 0 0]
	cpu.setZf(cpu.a == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(false)
}

func (cpu *CPU) orAddr(a16 uint16, mem *mem.Memory) {
	cpu.or(mem.Read(a16))
}

func (cpu *CPU) pop(mem *mem.Memory, r8 *uint8) func() {
	return func() {
		*r8 = mem.Read(cpu.sp)
		cpu.sp++
	}
}

func (cpu *CPU) popF(mem *mem.Memory) func() {
	return func() {
		// Lower nibble is always zero no matter what data was written
		cpu.f = mem.Read(cpu.sp) & 0xf0
		cpu.sp++
	}
}

func (cpu *CPU) push(mem *mem.Memory, r8 *uint8) func() {
	return func() {
		cpu.sp--
		mem.Write(cpu.sp, *r8)
	}
}

func (cpu *CPU) res(pos uint8, r8 *uint8) func() {
	return func() {
		*r8 &^= bits[pos]
	}
}

func (cpu *CPU) ret(kind string, mem *mem.Memory) {
	switch kind {
	case "":
		lsb := mem.Read(cpu.sp)
		cpu.sp++
		msb := mem.Read(cpu.sp)
		cpu.sp++
		retAddr := uint16(msb)<<8 | uint16(lsb)
		cpu.pc = retAddr
	case "NZ":
		if !cpu.zf() {
			cpu.ret("", mem)
		} else {
			// cpu.altTicks = true
		}
	case "Z":
		if cpu.zf() {
			cpu.ret("", mem)
		} else {
			// cpu.altTicks = true
		}
	case "NC":
		if !cpu.cf() {
			cpu.ret("", mem)
		} else {
			// cpu.altTicks = true
		}
	case "C":
		if cpu.cf() {
			cpu.ret("", mem)
		} else {
			// cpu.altTicks = true
		}
	default:
		panic(fmt.Sprintf("No implementation for ret"))
	}
}

func (cpu *CPU) reti(mem *mem.Memory) {
	cpu.ret("", mem)
	cpu.ei()
}

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

func (cpu *CPU) rlAddr(a16 uint16, mem *mem.Memory) {
	value := mem.Read(a16)
	cpu.rl(&value)
	mem.Write(a16, value)
}

func (cpu *CPU) rla() {
	cpu.rl(&cpu.a)
	cpu.f &^= zFlag
	// [0 0 0 C]
	cpu.setZf(false)
}

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

func (cpu *CPU) rlcAddr(a16 uint16, mem *mem.Memory) {
	value := mem.Read(a16)
	cpu.rlc(&value)
	mem.Write(a16, value)
}

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

func (cpu *CPU) rrAddr(a16 uint16, mem *mem.Memory) {
	value := mem.Read(a16)
	cpu.rr(&value)
	mem.Write(a16, value)
}

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

func (cpu *CPU) rrcAddr(a16 uint16, mem *mem.Memory) {
	value := mem.Read(a16)
	cpu.rrc(&value)
	mem.Write(a16, value)
}

func (cpu *CPU) rst(a16 uint16, mem *mem.Memory) {
	cpu.call("", a16, mem)
}

func (cpu *CPU) set(pos uint8, r8 *uint8) func() {
	return func() {
		*r8 |= bits[pos]
	}
}

func (cpu *CPU) sla(r8 *uint8) {
	cf := (*r8 & 0x80) > 0
	*r8 <<= 1
	// [Z 0 0 C]
	cpu.setZf(*r8 == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(cf)
}

func (cpu *CPU) slaAddr(a16 uint16, mem *mem.Memory) {
	value := mem.Read(a16)
	cpu.sla(&value)
	mem.Write(a16, value)
}

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

func (cpu *CPU) sraAddr(a16 uint16, mem *mem.Memory) {
	value := mem.Read(a16)
	cpu.sra(&value)
	mem.Write(a16, value)
}

func (cpu *CPU) srl(r8 *uint8) {
	cf := (*r8 & 0x01) > 0
	*r8 >>= 1
	// [Z 0 0 C]
	cpu.setZf(*r8 == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(cf)
}

func (cpu *CPU) srlAddr(a16 uint16, mem *mem.Memory) {
	value := mem.Read(a16)
	cpu.srl(&value)
	mem.Write(a16, value)
}

func (cpu *CPU) swap(r8 *uint8) {
	u8 := *r8
	*r8 = u8<<4 | u8>>4
	// [Z 0 0 0]
	cpu.setZf(*r8 == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(false)
}

func (cpu *CPU) swapAddr(a16 uint16, mem *mem.Memory) {
	value := mem.Read(a16)
	cpu.swap(&value)
	mem.Write(a16, value)
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

func (cpu *CPU) sbcAddr(a16 uint16, mem *mem.Memory) {
	cpu.sbc(mem.Read(a16))
}

func (cpu *CPU) sub(u8 uint8) {
	a := cpu.a
	cpu.a -= u8
	// [Z 1 H C]
	cpu.setZf(cpu.a == 0)
	cpu.setNf(true)
	cpu.setHf(hc8Sub(a, u8))
	cpu.setCf(c8Sub(a, u8))
}

func (cpu *CPU) subAddr(a16 uint16, mem *mem.Memory) {
	cpu.sub(mem.Read(a16))
}

func (cpu *CPU) xor(u8 uint8) {
	cpu.a = cpu.a ^ u8
	// [Z 0 0 0]
	cpu.setZf(cpu.a == 0)
	cpu.setNf(false)
	cpu.setHf(false)
	cpu.setCf(false)
}

func (cpu *CPU) xorAddr(a16 uint16, mem *mem.Memory) {
	cpu.xor(mem.Read(a16))
}
