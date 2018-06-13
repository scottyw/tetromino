package cpu

import (
	"fmt"

	"github.com/scottyw/tetromino/pkg/gb/mem"
)

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

func (cpu *CPU) adcAddr(a16 uint16, mem *mem.Memory) {
	cpu.adc(mem.Read(a16))
}

func (cpu *CPU) add(u8 uint8) {
	a := cpu.a
	cpu.a += u8
	// [Z 0 H C]
	cpu.setZf(cpu.a == 0)
	cpu.setNf(false)
	cpu.setHf(hc8(a, u8))
	cpu.setCf(c8(a, u8))
}

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

func (cpu *CPU) addSP(i8 int8) {
	sp := cpu.sp
	delta := int16(i8)
	cpu.sp = uint16(int16(cpu.sp) + delta)
	var hf, cf bool
	if delta >= 0 {
		hf = hc16(sp, uint16(delta))
		cf = c16(sp, uint16(delta))
	} else {
		hf = hc16Sub(sp, uint16(-delta))
		cf = c16Sub(sp, uint16(-delta))
	}
	// [0 0 H C]
	cpu.setZf(false)
	cpu.setNf(false)
	cpu.setHf(hf)
	cpu.setCf(cf)
}

func (cpu *CPU) addAddr(a16 uint16, mem *mem.Memory) {
	cpu.add(mem.Read(a16))
}

func (cpu *CPU) and(u8 uint8) {
	cpu.a &= u8
	// [Z 0 1 0]
	cpu.setZf(cpu.a == 0)
	cpu.setNf(false)
	cpu.setHf(true)
	cpu.setCf(false)
}

func (cpu *CPU) andAddr(a16 uint16, mem *mem.Memory) {
	cpu.and(mem.Read(a16))
}

func (cpu *CPU) bit(pos uint8, u8 uint8) {
	zero := u8&bits[pos] == 0
	// [Z 0 1 -]
	cpu.setZf(zero)
	cpu.setNf(false)
	cpu.setHf(true)
}

func (cpu *CPU) bitAddr(pos uint8, a16 uint16, mem *mem.Memory) {
	cpu.bit(pos, mem.Read(a16))
}

func (cpu *CPU) call(kind string, a16 uint16, mem *mem.Memory) {
	switch kind {
	case "":
		if cpu.debugFlowControl {
			fmt.Printf("==== CALL %04x --> %04x\n", cpu.pc, a16)
		}
		cpu.sp--
		mem.Write(cpu.sp, byte(cpu.pc>>8))
		cpu.sp--
		mem.Write(cpu.sp, byte(cpu.pc&0xff))
		cpu.pc = a16
	case "NZ":
		if !cpu.zf() {
			if cpu.debugFlowControl {
				fmt.Printf("==== CALL NZ ...\n")
			}
			cpu.call("", a16, mem)
		}
	case "Z":
		if cpu.zf() {
			if cpu.debugFlowControl {
				fmt.Printf("==== CALL Z ...\n")
			}
			cpu.call("", a16, mem)
		}
	case "NC":
		if !cpu.cf() {
			if cpu.debugFlowControl {
				fmt.Printf("==== CALL NC ...\n")
			}
			cpu.call("", a16, mem)
		}
	case "C":
		if cpu.cf() {
			if cpu.debugFlowControl {
				fmt.Printf("==== CALL C ...\n")
			}
			cpu.call("", a16, mem)
		}
	default:
		panic(fmt.Sprintf("No implementation for call: %v %v", kind, a16))
	}
}

func (cpu *CPU) ccf() {
	// [- 0 0 *]
	cpu.setZf(cpu.a == 0)
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
	cf := cpu.cf()
	if !cpu.nf() {
		if cpu.cf() || cpu.a > 0x99 {
			cpu.a += 0x60
			cf = true
		}
		if cpu.hf() || (cpu.a&0x0f) > 0x09 {
			cpu.a += 0x6
		}
	} else {
		if cpu.cf() {
			cpu.a -= 0x60
		}
		if cpu.hf() {
			cpu.a -= 0x6
		}
	}
	// [Z - 0 C]
	cpu.setZf(cpu.a == 0)
	cpu.setHf(false)
	cpu.setCf(cf)
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
		if cpu.debugJumps {
			fmt.Printf("==== JP %04x --> %04x\n", cpu.pc, a16)
		}
		cpu.pc = a16
	case "NZ":
		if !cpu.zf() {
			if cpu.debugJumps {
				fmt.Printf("==== JP NZ\n")
			}
			cpu.jp("", a16)
		}
	case "Z":
		if cpu.zf() {
			if cpu.debugJumps {
				fmt.Printf("==== JP Z\n")
			}
			cpu.jp("", a16)
		}
	case "NC":
		if !cpu.cf() {
			if cpu.debugJumps {
				fmt.Printf("==== JP NC\n")
			}
			cpu.jp("", a16)
		}
	case "C":
		if cpu.cf() {
			if cpu.debugJumps {
				fmt.Printf("==== JP C\n")
			}
			cpu.jp("", a16)
		}
	default:
		panic(fmt.Sprintf("No implementation for jp: %v %v", kind, a16))
	}
}

func (cpu *CPU) jr(kind string, i8 int8) {
	address := uint16(int16(cpu.pc) + int16(i8))
	if cpu.debugJumps {
		fmt.Printf("==== JR %02x\n", i8)
	}
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

func (cpu *CPU) ldhAA8(u8 uint8, mem *mem.Memory) {
	address := uint16(0xff00 + uint16(u8))
	cpu.a = mem.Read(address)
}

func (cpu *CPU) ldhA8A(u8 uint8, mem *mem.Memory) {
	address := uint16(0xff00 + uint16(u8))
	mem.Write(address, cpu.a)
}

func (cpu *CPU) ldSP(u16 uint16) {
	cpu.sp = u16
}

func (cpu *CPU) ldSPHL() {
	cpu.sp = cpu.hl()
}

func (cpu *CPU) ldA16SP(a16 uint16, mem *mem.Memory) {
	mem.Write(a16, uint8(cpu.sp>>8))
	mem.Write(a16+1, uint8(cpu.sp|0x0f))
}

func (cpu *CPU) ldHLSP(i8 int8) {
	old := cpu.hl()
	new := uint16(int16(cpu.sp) + int16(i8))
	cpu.h = uint8(new >> 8)
	cpu.l = uint8(new)
	// [0 0 H C]
	cpu.setZf(false)
	cpu.setNf(false)
	cpu.setHf(hc16(old, new))
	cpu.setCf(c16(old, new))
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

func (cpu *CPU) nop() {
	// Do nothing
	return
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

func (cpu *CPU) pop(msb, lsb *uint8, mem *mem.Memory) {
	*lsb = mem.Read(cpu.sp)
	cpu.sp++
	*msb = mem.Read(cpu.sp)
	cpu.sp++
}

func (cpu *CPU) popAF(mem *mem.Memory) {
	cpu.pop(&cpu.a, &cpu.f, mem)
	// Lower nibble is always zero no matter what data was written
	cpu.f &= 0xf0
}

func (cpu *CPU) push(msb, lsb uint8, mem *mem.Memory) {
	cpu.sp--
	mem.Write(cpu.sp, msb)
	cpu.sp--
	mem.Write(cpu.sp, lsb)
}

func (cpu *CPU) res(pos uint8, r8 *uint8) {
	*r8 &^= bits[pos]
}

func (cpu *CPU) resAddr(pos uint8, a16 uint16, mem *mem.Memory) {
	value := mem.Read(a16)
	cpu.res(pos, &value)
	mem.Write(a16, value)
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
		if cpu.debugFlowControl {
			fmt.Printf("==== RET %04x --> %04x\n", retAddr, cpu.pc)
		}
	case "NZ":
		if !cpu.zf() {
			if cpu.debugFlowControl {
				fmt.Printf("==== RET NZ ...\n")
			}
			cpu.ret("", mem)
		}
	case "Z":
		if cpu.zf() {
			if cpu.debugFlowControl {
				fmt.Printf("==== RET Z ...\n")
			}
			cpu.ret("", mem)
		}
	case "NC":
		if !cpu.cf() {
			if cpu.debugFlowControl {
				fmt.Printf("==== RET NC ...\n")
			}
			cpu.ret("", mem)
		}
	case "C":
		if cpu.cf() {
			if cpu.debugFlowControl {
				fmt.Printf("==== RET C ...\n")
			}
			cpu.ret("", mem)
		}
	default:
		panic(fmt.Sprintf("No implementation for ret"))
	}
}

func (cpu *CPU) reti(mem *mem.Memory) {
	if cpu.debugFlowControl {
		fmt.Printf("==== RETI ...\n")
	}
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
	if cpu.debugFlowControl {
		fmt.Printf("==== RST %04x ...\n", a16)
	}
	cpu.call("", a16, mem)
}

func (cpu *CPU) set(pos uint8, r8 *uint8) {
	*r8 |= bits[pos]
}

func (cpu *CPU) setAddr(pos uint8, a16 uint16, mem *mem.Memory) {
	value := mem.Read(a16)
	cpu.set(pos, &value)
	mem.Write(a16, value)
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
	cpu.setZf(cpu.a == 0)
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
