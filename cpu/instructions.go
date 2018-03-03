package cpu

import (
	"fmt"

	"github.com/scottyw/tetromino/mem"
	"github.com/scottyw/tetromino/options"
)

func (cpu *CPU) adc(u8 uint8) {
	a := cpu.a
	cpu.a += u8
	hf := h(a, cpu.a)
	cf := c(a, cpu.a)
	if cpu.cf {
		a = cpu.a
		cpu.a++
		hf = hf || h(a, cpu.a)
		cf = cf || c(a, cpu.a)
	}
	cpu.flags(z(cpu.a), false, hf, cf) // [Z 0 H C]
}

func (cpu *CPU) adcAddr(a16 uint16, mem mem.Memory) {
	cpu.adc(*mem.Read(a16))
}

func (cpu *CPU) add(u8 uint8) {
	a := cpu.a
	cpu.a += u8
	cpu.flags(z(cpu.a), false, h(a, cpu.a), c(a, cpu.a)) // [Z 0 H C]
}

func (cpu *CPU) addHL(u16 uint16) {
	old := cpu.hl()
	new := old + u16
	cpu.h = uint8(new >> 8)
	cpu.l = uint8(new)
	cpu.flags(cpu.zf, false, h16(old, new), c16(old, new)) // [- 0 H C]
}

func (cpu *CPU) addSP(i8 int8) {
	panic(fmt.Sprintf("Missing implementation for addSP: %v", i8))
}

func (cpu *CPU) addAddr(a16 uint16, mem mem.Memory) {
	panic(fmt.Sprintf("Missing implementation for addAddr: %v", a16))
}

func (cpu *CPU) and(u8 uint8) {
	cpu.a &= u8
	cpu.flags(z(cpu.a), false, true, false) // [Z 0 1 0]
}

func (cpu *CPU) andAddr(a16 uint16, mem mem.Memory) {
	cpu.and(*mem.Read(a16))
}

func (cpu *CPU) bit(pos uint8, u8 uint8) {
	zero := u8&bits[pos] == 0
	cpu.flags(zero, false, true, cpu.cf) // [Z 0 1 -]
}

func (cpu *CPU) bitAddr(pos uint8, a16 uint16, mem mem.Memory) {
	cpu.bit(pos, *mem.Read(a16))
}

func (cpu *CPU) call(kind string, a16 uint16, mem mem.Memory) {
	switch kind {
	case "":
		if options.DebugFlowControl() {
			fmt.Printf("==== CALL %04x --> %04x\n", cpu.pc, a16)
		}
		*mem.Read(cpu.sp) = byte(cpu.pc & 0xff)
		cpu.sp--
		*mem.Read(cpu.sp) = byte(cpu.pc >> 8)
		cpu.sp--
		cpu.pc = a16
	case "NZ":
		if !cpu.zf {
			if options.DebugFlowControl() {
				fmt.Printf("==== CALL NZ ...\n")
			}
			cpu.call("", a16, mem)
		}
	case "Z":
		if cpu.zf {
			if options.DebugFlowControl() {
				fmt.Printf("==== CALL Z ...\n")
			}
			cpu.call("", a16, mem)
		}
	case "NC":
		if !cpu.cf {
			if options.DebugFlowControl() {
				fmt.Printf("==== CALL NC ...\n")
			}
			cpu.call("", a16, mem)
		}
	case "C":
		if cpu.cf {
			if options.DebugFlowControl() {
				fmt.Printf("==== CALL C ...\n")
			}
			cpu.call("", a16, mem)
		}
	default:
		panic(fmt.Sprintf("Missing implementation for call: %v %v", kind, a16))
	}
}

func (cpu *CPU) ccf() {
	panic(fmt.Sprintf("Missing implementation for ccf"))
}

func (cpu *CPU) cp(u8 uint8) {
	cpu.flags(cpu.a == u8, true, h(u8, cpu.a), c(u8, cpu.a)) // [Z 1 H C]
}

func (cpu *CPU) cpAddr(a16 uint16, mem mem.Memory) {
	cpu.cp(*mem.Read(a16))
}

func (cpu *CPU) cpl() {
	cpu.a = ^cpu.a
	cpu.flags(cpu.zf, true, true, cpu.cf)
}

func (cpu *CPU) daa() {
	panic(fmt.Sprintf("Missing implementation for daa"))
}

func (cpu *CPU) dec(r8 *uint8) {
	old := *r8
	*r8--
	cpu.flags(z(*r8), true, h(*r8, old), cpu.cf) //	[Z 1 H -]
}

func (cpu *CPU) dec16(msb, lsb *uint8) {
	old := uint16(*msb)<<8 + uint16(*lsb)
	new := old - 1
	*msb = uint8(new >> 8)
	*lsb = uint8(new)
	cpu.flags(z16(new), true, h16(new, old), cpu.cf) //	[Z 1 H -]
}

func (cpu *CPU) decSP() {
	panic(fmt.Sprintf("Missing implementation for decSP"))
}

func (cpu *CPU) decAddr(a16 uint16, mem mem.Memory) {
	cpu.dec(mem.Read(a16))
}

func (cpu *CPU) di() {
	cpu.ime = false
}

func (cpu *CPU) ei() {
	cpu.ime = true
}

func (cpu *CPU) halt() {
	panic(fmt.Sprintf("Missing implementation for halt"))
}

func (cpu *CPU) inc(r8 *uint8) {
	old := *r8
	*r8++
	cpu.flags(z(*r8), true, h(old, *r8), cpu.cf) // [Z 0 H -]
}

func (cpu *CPU) inc16(msb, lsb *uint8) {
	old := uint16(*msb)<<8 + uint16(*lsb)
	new := old + 1
	*msb = uint8(new >> 8)
	*lsb = uint8(new)
	cpu.flags(z16(new), true, h16(old, new), cpu.cf) //	[Z 1 H -]
}

func (cpu *CPU) incSP() {
	panic(fmt.Sprintf("Missing implementation for incSP"))
}

func (cpu *CPU) incAddr(a16 uint16, mem mem.Memory) {
	cpu.inc(mem.Read(a16))
}

func (cpu *CPU) jp(kind string, a16 uint16) {
	switch kind {
	case "":
		if options.DebugJumps() {
			fmt.Printf("==== JP %04x --> %04x\n", cpu.pc, a16)
		}
		cpu.pc = a16
	case "NZ":
		if !cpu.zf {
			if options.DebugJumps() {
				fmt.Printf("==== JP NZ\n")
			}
			cpu.jp("", a16)
		}
	case "Z":
		if cpu.zf {
			if options.DebugJumps() {
				fmt.Printf("==== JP Z\n")
			}
			cpu.jp("", a16)
		}
	case "NC":
		if !cpu.cf {
			if options.DebugJumps() {
				fmt.Printf("==== JP NC\n")
			}
			cpu.jp("", a16)
		}
	case "C":
		if cpu.cf {
			if options.DebugJumps() {
				fmt.Printf("==== JP C\n")
			}
			cpu.jp("", a16)
		}
	default:
		panic(fmt.Sprintf("Missing implementation for jp: %v %v", kind, a16))
	}
}

func (cpu *CPU) jr(kind string, i8 int8) {
	address := uint16(int16(cpu.pc) + int16(i8))
	if options.DebugJumps() {
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

func (cpu *CPU) ldFromAddr(r8 *uint8, a16 uint16, mem mem.Memory) {
	*r8 = *mem.Read(a16)
}

func (cpu *CPU) ldToAddr(a16 uint16, u8 uint8, mem mem.Memory) {
	*mem.Read(a16) = u8
}

func (cpu *CPU) ldhFromAddr(u8 uint8, mem mem.Memory) {
	address := uint16(0xff00 + uint16(u8))
	cpu.a = *mem.Read(address)
}

func (cpu *CPU) ldhToAddr(u8 uint8, mem mem.Memory) {
	address := uint16(0xff00 + uint16(u8))
	*mem.Read(address) = cpu.a
}

func (cpu *CPU) ldAFromAddrC(mem mem.Memory) {
	cpu.ldhFromAddr(cpu.c, mem)
}

func (cpu *CPU) ldAToAddrC(mem mem.Memory) {
	cpu.ldhToAddr(cpu.c, mem)
}

func (cpu *CPU) ldSP(u16 uint16) {
	cpu.sp = u16
}

func (cpu *CPU) ldHLToSP() {
	panic(fmt.Sprintf("Missing implementation for ldHLToSP"))
}

func (cpu *CPU) ldSPToAddr(a16 uint16, mem mem.Memory) {
	*mem.Read(a16) = uint8(cpu.sp >> 8)
	*mem.Read(a16 + 1) = uint8(cpu.sp | 0x0f)
}

func (cpu *CPU) ldSPToHL(i8 int8) {
	panic(fmt.Sprintf("Missing implementation for ldSPToHL: %v", i8))
}

func (cpu *CPU) lddFromAddr(mem mem.Memory) {
	cpu.ldFromAddr(&cpu.a, cpu.hl(), mem)
	cpu.dec16(&cpu.h, &cpu.l)
}

func (cpu *CPU) lddToAddr(mem mem.Memory) {
	cpu.ldToAddr(cpu.hl(), cpu.a, mem)
	cpu.dec16(&cpu.h, &cpu.l)
}

func (cpu *CPU) ldiFromAddr(mem mem.Memory) {
	cpu.ldFromAddr(&cpu.a, cpu.hl(), mem)
	cpu.inc16(&cpu.h, &cpu.l)
}

func (cpu *CPU) ldiToAddr(mem mem.Memory) {
	cpu.ldToAddr(cpu.hl(), cpu.a, mem)
	cpu.inc16(&cpu.h, &cpu.l)
}

func (cpu *CPU) nop() {
	// Do nothing
	return
}

func (cpu *CPU) or(u8 uint8) {
	cpu.a |= u8
	cpu.flags(z(cpu.a), false, false, false) // [Z 0 0 0]
}

func (cpu *CPU) orAddr(a16 uint16, mem mem.Memory) {
	cpu.or(*mem.Read(a16))
}

func (cpu *CPU) pop(msb, lsb *uint8, mem mem.Memory) {
	cpu.sp++
	*msb = *mem.Read(cpu.sp)
	cpu.sp++
	*lsb = *mem.Read(cpu.sp)
}

func (cpu *CPU) popAF(mem mem.Memory) {
	cpu.pop(&cpu.a, &cpu.f, mem)
	cpu.zf = cpu.f&8 > 1
	cpu.nf = cpu.f&4 > 1
	cpu.hf = cpu.f&2 > 1
	cpu.cf = cpu.f&1 > 1
}

func (cpu *CPU) push(msb, lsb uint8, mem mem.Memory) {
	*mem.Read(cpu.sp) = lsb
	cpu.sp--
	*mem.Read(cpu.sp) = msb
	cpu.sp--
}

func (cpu *CPU) res(pos uint8, r8 *uint8) {
	*r8 &^= bits[pos]
}

func (cpu *CPU) resAddr(pos uint8, a16 uint16, mem mem.Memory) {
	cpu.res(pos, mem.Read(a16))
}

func (cpu *CPU) ret(kind string, mem mem.Memory) {
	switch kind {
	case "":
		cpu.sp++
		msb := *mem.Read(cpu.sp)
		cpu.sp++
		lsb := *mem.Read(cpu.sp)
		retAddr := uint16(msb)<<8 | uint16(lsb)
		cpu.pc = retAddr
		if options.DebugFlowControl() {
			fmt.Printf("==== RET %04x --> %04x\n", retAddr, cpu.pc)
		}
	case "NZ":
		if !cpu.zf {
			if options.DebugFlowControl() {
				fmt.Printf("==== RET NZ ...\n")
			}
			cpu.ret("", mem)
		}
	case "Z":
		if cpu.zf {
			if options.DebugFlowControl() {
				fmt.Printf("==== RET Z ...\n")
			}
			cpu.ret("", mem)
		}
	case "NC":
		if !cpu.cf {
			if options.DebugFlowControl() {
				fmt.Printf("==== RET NC ...\n")
			}
			cpu.ret("", mem)
		}
	case "C":
		if cpu.cf {
			if options.DebugFlowControl() {
				fmt.Printf("==== RET C ...\n")
			}
			cpu.ret("", mem)
		}
	default:
		panic(fmt.Sprintf("Missing implementation for ret"))
	}
}

func (cpu *CPU) reti(mem mem.Memory) {
	if options.DebugFlowControl() {
		fmt.Printf("==== RETI ...\n")
	}
	cpu.ret("", mem)
	cpu.ei()
}

func (cpu *CPU) rl(r8 *uint8) {
	cf := (*r8 & 0x80) > 0
	*r8 <<= 1
	if cpu.cf {
		*r8 |= 0x01
	}
	cpu.flags(z(*r8), false, false, cf) //  [Z 0 0 C]
}

func (cpu *CPU) rlAddr(a16 uint16, mem mem.Memory) {
	cpu.rl(mem.Read(a16))
}

func (cpu *CPU) rla() {
	cpu.rl(&cpu.a)
	cpu.zf = false
}

func (cpu *CPU) rlc(r8 *uint8) {
	cf := (*r8 & 0x80) > 0
	*r8 <<= 1
	if cf {
		*r8 |= 0x01
	}
	cpu.flags(z(*r8), false, false, cf) //  [Z 0 0 C]
}

func (cpu *CPU) rlca() {
	cpu.rlc(&cpu.a)
	cpu.zf = false
}

func (cpu *CPU) rlcAddr(a16 uint16, mem mem.Memory) {
	cpu.rlc(mem.Read(a16))
}

func (cpu *CPU) rr(r8 *uint8) {
	cf := (*r8 & 0x01) > 0
	*r8 >>= 1
	if cpu.cf {
		*r8 |= 0x80
	}
	cpu.flags(z(*r8), false, false, cf) //  [Z 0 0 C]
}

func (cpu *CPU) rra() {
	cpu.rr(&cpu.a)
}

func (cpu *CPU) rrAddr(a16 uint16, mem mem.Memory) {
	cpu.rr(mem.Read(a16))
}

func (cpu *CPU) rrc(r8 *uint8) {
	cf := (*r8 & 0x01) > 0
	*r8 >>= 1
	if cf {
		*r8 |= 0x80
	}
	cpu.flags(z(*r8), false, false, cf) //  [Z 0 0 C]
}

func (cpu *CPU) rrca() {
	cpu.rrc(&cpu.a)
}

func (cpu *CPU) rrcAddr(a16 uint16, mem mem.Memory) {
	cpu.rrc(mem.Read(a16))
}

func (cpu *CPU) rst(a16 uint16, mem mem.Memory) {
	if options.DebugFlowControl() {
		fmt.Printf("==== RST %04x ...\n", a16)
	}
	cpu.call("", a16, mem)
}

func (cpu *CPU) set(pos uint8, r8 *uint8) {
	*r8 |= bits[pos]
}

func (cpu *CPU) setAddr(pos uint8, a16 uint16, mem mem.Memory) {
	cpu.set(pos, mem.Read(a16))
}

func (cpu *CPU) sla(r8 *uint8) {
	cf := (*r8 & 0x80) > 0
	*r8 <<= 1
	cpu.flags(z(*r8), false, false, cf) //  [Z 0 0 C]
}

func (cpu *CPU) slaAddr(a16 uint16, mem mem.Memory) {
	cpu.sla(mem.Read(a16))
}

func (cpu *CPU) sra(r8 *uint8) {
	cf := (*r8 & 0x01) > 0
	bit7 := (*r8 & 0x80) > 0
	*r8 >>= 1
	if bit7 {
		*r8 |= 0x80
	}
	cpu.flags(z(*r8), false, false, cf) //  [Z 0 0 C]
}

func (cpu *CPU) sraAddr(a16 uint16, mem mem.Memory) {
	cpu.sra(mem.Read(a16))
}

func (cpu *CPU) srl(r8 *uint8) {
	cf := (*r8 & 0x01) > 0
	*r8 >>= 1
	cpu.flags(z(*r8), false, false, cf) //  [Z 0 0 C]
}

func (cpu *CPU) srlAddr(a16 uint16, mem mem.Memory) {
	cpu.srl(mem.Read(a16))
}

func (cpu *CPU) swap(r8 *uint8) {
	u8 := *r8
	*r8 = u8<<4 | u8>>4
	cpu.flags(z(*r8), false, false, false)
}

func (cpu *CPU) swapAddr(a16 uint16, mem mem.Memory) {
	cpu.swap(mem.Read(a16))
}

func (cpu *CPU) scf() {
	panic(fmt.Sprintf("Missing implementation for scf"))
}

func (cpu *CPU) stop() {
	panic(fmt.Sprintf("Missing implementation for stop"))
}

func (cpu *CPU) sbc(u8 uint8) {
	panic(fmt.Sprintf("Missing implementation for sbc: %v", u8))
}

func (cpu *CPU) sbcAddr(a16 uint16, mem mem.Memory) {
	panic(fmt.Sprintf("Missing implementation for sbcAddr: %v", a16))
}

func (cpu *CPU) sub(u8 uint8) {
	a := cpu.a
	cpu.a -= u8
	cpu.flags(z(cpu.a), false, h(cpu.a, a), c(cpu.a, a)) // [Z 0 H C]
}

func (cpu *CPU) subAddr(a16 uint16, mem mem.Memory) {
	panic(fmt.Sprintf("Missing implementation for subAddr: %v", a16))
}

func (cpu *CPU) xor(u8 uint8) {
	cpu.a = cpu.a ^ u8
	cpu.flags(z(cpu.a), false, false, false) // [Z 0 0 0]
}

func (cpu *CPU) xorAddr(a16 uint16, mem mem.Memory) {
	panic(fmt.Sprintf("Missing implementation for xorAddr: %v", a16))
}
