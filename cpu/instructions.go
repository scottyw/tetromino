package cpu

import (
	"fmt"

	"github.com/scottyw/goomba/mem"
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
	hl := cpu.hl()
	old := hl.Get()
	hl.Set(old + u16)
	cpu.flags(cpu.zf, false, h16(old, hl.Get()), c16(old, hl.Get())) // [- 0 H C]
}

func (cpu *CPU) addSP(u8 uint8) {
	panic(fmt.Sprintf("Missing implementation for addSP: %v", u8))
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

func (cpu *CPU) call(kind string, u16 uint16) {
	panic(fmt.Sprintf("Missing implementation for call: %v %v", kind, u16))
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
	panic(fmt.Sprintf("Missing implementation for cpl"))
}

func (cpu *CPU) daa() {
	panic(fmt.Sprintf("Missing implementation for daa"))
}

func (cpu *CPU) dec(r8 *uint8) {
	old := *r8
	*r8--
	cpu.flags(z(*r8), true, h(*r8, old), cpu.cf) //	[Z 1 H -]
}

func (cpu *CPU) dec16(r16 register16) {
	r16.Set(r16.Get() - 1)
}

func (cpu *CPU) decSP() {
	panic(fmt.Sprintf("Missing implementation for decSP"))
}

func (cpu *CPU) decAddr(a16 uint16, mem mem.Memory) {
	panic(fmt.Sprintf("Missing implementation for decAddr: %v", a16))
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

func (cpu *CPU) inc16(r16 register16) {
	r16.Set(r16.Get() + 1)
}

func (cpu *CPU) incSP() {
	panic(fmt.Sprintf("Missing implementation for incSP"))
}

func (cpu *CPU) incAddr(a16 uint16, mem mem.Memory) {
	panic(fmt.Sprintf("Missing implementation for incAddr: %v", a16))
}

func (cpu *CPU) jp(kind string, u16 uint16) {
	switch kind {
	case "":
		cpu.pc = u16
	case "NZ":
		if !cpu.zf {
			cpu.pc = u16
		}
	case "Z":
		if cpu.zf {
			cpu.pc = u16
		}
	case "NC":
		if !cpu.cf {
			cpu.pc = u16
		}
	case "C":
		if cpu.cf {
			cpu.pc = u16
		}
	default:
		panic(fmt.Sprintf("Missing implementation for jp: %v %v", kind, u16))
	}
}

func (cpu *CPU) jr(kind string, u8 uint8) {
	cpu.jp(kind, cpu.pc+uint16(u8))
}

func (cpu *CPU) ld(r8 *uint8, u8 uint8) {
	*r8 = u8
}

func (cpu *CPU) ld16(r16 register16, u16 uint16) {
	r16.Set(u16)
}

func (cpu *CPU) ldFromAddr(r8 *uint8, a16 uint16, mem mem.Memory) {
	*r8 = *mem.Read(a16)
}

func (cpu *CPU) ldToAddr(a16 uint16, u8 uint8, mem mem.Memory) {
	*mem.Read(a16) = u8
}

func (cpu *CPU) ldhFromAddr(a8 uint8, mem mem.Memory) {
	panic(fmt.Sprintf("Missing implementation for ldhFromAddr: %v", a8))
}

func (cpu *CPU) ldhToAddr(a8 uint8, mem mem.Memory) {
	panic(fmt.Sprintf("Missing implementation for ldhToAddr: %v", a8))
}

func (cpu *CPU) ldAFromAddrC() {
	panic(fmt.Sprintf("Missing implementation for ldAFromAddrC"))
}

func (cpu *CPU) ldAToAddrC() {
	panic(fmt.Sprintf("Missing implementation for ldAToAddrC"))
}

func (cpu *CPU) ldSP(u16 uint16) {
	panic(fmt.Sprintf("Missing implementation for ldSP: %v", u16))
}

func (cpu *CPU) ldHLToSP() {
	panic(fmt.Sprintf("Missing implementation for ldHLToSP"))
}

func (cpu *CPU) ldSPToAddr(a16 uint16, mem mem.Memory) {
	panic(fmt.Sprintf("Missing implementation for ldSPToAddr: %v", a16))
}

func (cpu *CPU) ldSPToHL(u8 uint8) {
	panic(fmt.Sprintf("Missing implementation for ldSPToHL: %v", u8))
}

func (cpu *CPU) lddFromAddr(mem mem.Memory) {
	hl := cpu.hl()
	cpu.ldFromAddr(&cpu.a, hl.Get(), mem)
	cpu.dec16(hl)
}

func (cpu *CPU) lddToAddr(mem mem.Memory) {
	hl := cpu.hl()
	cpu.ldToAddr(hl.Get(), cpu.a, mem)
	cpu.dec16(hl)
}

func (cpu *CPU) ldiFromAddr(mem mem.Memory) {
	panic(fmt.Sprintf("Missing implementation for ldiFromAddr"))
}

func (cpu *CPU) ldiToAddr(mem mem.Memory) {
	panic(fmt.Sprintf("Missing implementation for ldiToAddr"))
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

func (cpu *CPU) pop(r16 register16) {
	panic(fmt.Sprintf("Missing implementation for pop:   %v", r16))
}

func (cpu *CPU) push(r16 register16) {
	panic(fmt.Sprintf("Missing implementation for push:   %v", r16))
}

func (cpu *CPU) res(pos uint8, r8 *uint8) {
	*r8 &^= bits[pos]
}

func (cpu *CPU) resAddr(pos uint8, a16 uint16, mem mem.Memory) {
	cpu.res(pos, mem.Read(a16))
}

func (cpu *CPU) ret(kind string) {
	panic(fmt.Sprintf("Missing implementation for ret:   %v", kind))
}

func (cpu *CPU) reti() {
	panic(fmt.Sprintf("Missing implementation for reti"))
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
	panic(fmt.Sprintf("Missing implementation for rlAddr: %v", a16))
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
	panic(fmt.Sprintf("Missing implementation for rlcAddr: %v", a16))
}

func (cpu *CPU) rr(r8 *uint8) {

}

func (cpu *CPU) rra() {
	cpu.rr(&cpu.a)
}

func (cpu *CPU) rrAddr(a16 uint16, mem mem.Memory) {
	panic(fmt.Sprintf("Missing implementation for rrAddr: %v", a16))
}

func (cpu *CPU) rrc(r8 *uint8) {
	panic(fmt.Sprintf("Missing implementation for rrc: %v", r8))
}

func (cpu *CPU) rrca() {
	panic(fmt.Sprintf("Missing implementation for rrca"))
}

func (cpu *CPU) rrcAddr(a16 uint16, mem mem.Memory) {
	panic(fmt.Sprintf("Missing implementation for rrcAddr: %v", a16))
}

func (cpu *CPU) rst(u8 uint8) {
	panic(fmt.Sprintf("Missing implementation for rst: %v", u8))
}

func (cpu *CPU) set(pos uint8, r8 *uint8) {
	*r8 |= bits[pos]
}

func (cpu *CPU) setAddr(pos uint8, a16 uint16, mem mem.Memory) {
	cpu.set(pos, mem.Read(a16))
}

func (cpu *CPU) sla(r8 *uint8) {
	panic(fmt.Sprintf("Missing implementation for sla: %v", r8))
}

func (cpu *CPU) slaAddr(a16 uint16, mem mem.Memory) {
	panic(fmt.Sprintf("Missing implementation for slaAddr: %v", a16))
}

func (cpu *CPU) sra(r8 *uint8) {
	panic(fmt.Sprintf("Missing implementation for sra: %v", r8))
}

func (cpu *CPU) sraAddr(a16 uint16, mem mem.Memory) {
	panic(fmt.Sprintf("Missing implementation for sraAddr: %v", a16))
}

func (cpu *CPU) srl(r8 *uint8) {
	panic(fmt.Sprintf("Missing implementation for srl: %v", r8))
}

func (cpu *CPU) srlAddr(a16 uint16, mem mem.Memory) {
	panic(fmt.Sprintf("Missing implementation for srlAddr: %v", a16))
}

func (cpu *CPU) swap(r8 *uint8) {
	panic(fmt.Sprintf("Missing implementation for swap: %v", r8))
}

func (cpu *CPU) swapAddr(a16 uint16, mem mem.Memory) {
	panic(fmt.Sprintf("Missing implementation for swapAddr: %v", a16))
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
	panic(fmt.Sprintf("Missing implementation for sub: %v", u8))
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
