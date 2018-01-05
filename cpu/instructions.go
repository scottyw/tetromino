package cpu

import (
	"fmt"

	"github.com/scottyw/goomba/mem"
)

func (cpu *CPU) adc(u8 uint8) {
	// cpu.flags(z(cpu.a), false, h(a, cpu.a), c(a, cpu.a)) // [Z 0 H C]
}

func (cpu *CPU) adcAddr(a16 uint16) {
	panic(fmt.Sprintf("Missing implementation for adcAddr: %v", a16))
}

func (cpu *CPU) add(u8 uint8) {
	a := cpu.a
	cpu.a += u8
	cpu.flags(z(cpu.a), false, h(a, cpu.a), c(a, cpu.a)) // [Z 0 H C]
}

func (cpu *CPU) addHL(u16 uint16) {
	panic(fmt.Sprintf("Missing implementation for addHL: %v", u16))
}

func (cpu *CPU) addSP(u8 uint8) {
	panic(fmt.Sprintf("Missing implementation for addSP: %v", u8))
}

func (cpu *CPU) addAddr(a16 uint16) {
	panic(fmt.Sprintf("Missing implementation for addAddr: %v", a16))
}

func (cpu *CPU) and(u8 uint8) {
	panic(fmt.Sprintf("Missing implementation for and: %v", u8))
}

func (cpu *CPU) andAddr(a16 uint16) {
	panic(fmt.Sprintf("Missing implementation for andAddr: %v", a16))
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
	panic(fmt.Sprintf("Missing implementation for cp: %v", u8))
}

func (cpu *CPU) cpAddr(a16 uint16) {
	panic(fmt.Sprintf("Missing implementation for cp: %v", a16))
}

func (cpu *CPU) cpl() {
	panic(fmt.Sprintf("Missing implementation for cpl"))
}

func (cpu *CPU) daa() {
	panic(fmt.Sprintf("Missing implementation for daa"))
}

func (cpu *CPU) dec(r8 *uint8) {
	panic(fmt.Sprintf("Missing implementation for dec: %v", r8))
}

func (cpu *CPU) dec16(r16 register16) {
	panic(fmt.Sprintf("Missing implementation for decAddr: %v", r16))
}

func (cpu *CPU) decSP() {
	panic(fmt.Sprintf("Missing implementation for decSP"))
}

func (cpu *CPU) decAddr(a16 uint16) {
	panic(fmt.Sprintf("Missing implementation for decAddr: %v", a16))
}

func (cpu *CPU) di() {
	panic(fmt.Sprintf("Missing implementation for di"))
}

func (cpu *CPU) ei() {
	panic(fmt.Sprintf("Missing implementation for ei"))
}

func (cpu *CPU) halt() {
	panic(fmt.Sprintf("Missing implementation for halt"))
}

func (cpu *CPU) inc(r8 *uint8) {
	panic(fmt.Sprintf("Missing implementation for inc8: %v", r8))
}

func (cpu *CPU) inc16(r16 register16) {
	panic(fmt.Sprintf("Missing implementation for incAddr: %v", r16))
}

func (cpu *CPU) incSP() {
	panic(fmt.Sprintf("Missing implementation for incSP"))
}

func (cpu *CPU) incAddr(a16 uint16) {
	panic(fmt.Sprintf("Missing implementation for incAddr: %v", a16))
}

func (cpu *CPU) jp(kind string, u16 uint16) {
	switch kind {
	case "":
		cpu.pc = u16
	default:
		panic(fmt.Sprintf("Missing implementation for jp: %v %v", kind, u16))
	}
}

func (cpu *CPU) jr(kind string, u8 uint8) {
	panic(fmt.Sprintf("Missing implementation for jr: %v %v", kind, u8))
}

func (cpu *CPU) ld(r8 *uint8, u8 uint8) {
	panic(fmt.Sprintf("Missing implementation for ld: %v %v", r8, u8))
}

func (cpu *CPU) ld16(r16 register16, u16 uint16) {
	panic(fmt.Sprintf("Missing implementation for ld16: %v %v", r16, u16))
}

func (cpu *CPU) ldFromAddr(r8 *uint8, a16 uint16) {
	panic(fmt.Sprintf("Missing implementation for ldFromAddr: %v %v", r8, a16))
}

func (cpu *CPU) ldToAddr(a16 uint16, u8 uint8) {
	panic(fmt.Sprintf("Missing implementation for ldToAddr: %v %v", a16, u8))
}

func (cpu *CPU) ldhFromAddr(a8 uint8) {
	panic(fmt.Sprintf("Missing implementation for ldhFromAddr: %v", a8))
}

func (cpu *CPU) ldhToAddr(a8 uint8) {
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

func (cpu *CPU) ldSPToAddr(a16 uint16) {
	panic(fmt.Sprintf("Missing implementation for ldSPToAddr: %v", a16))
}

func (cpu *CPU) ldSPToHL(u8 uint8) {
	panic(fmt.Sprintf("Missing implementation for ldSPToHL: %v", u8))
}

func (cpu *CPU) lddFromAddr() {
	panic(fmt.Sprintf("Missing implementation for lddFromAddr"))
}

func (cpu *CPU) lddToAddr() {
	panic(fmt.Sprintf("Missing implementation for lddToAddr"))
}

func (cpu *CPU) ldiFromAddr() {
	panic(fmt.Sprintf("Missing implementation for ldiFromAddr"))
}

func (cpu *CPU) ldiToAddr() {
	panic(fmt.Sprintf("Missing implementation for ldiToAddr"))
}

func (cpu *CPU) nop() {
	// Do nothing
	return
}

func (cpu *CPU) or(u8 uint8) {
	panic(fmt.Sprintf("Missing implementation for or: %v", u8))
}

func (cpu *CPU) orAddr(a16 uint16) {
	panic(fmt.Sprintf("Missing implementation for orAddr: %v", a16))
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
	panic(fmt.Sprintf("Missing implementation for rl: %v", r8))
}

func (cpu *CPU) rlAddr(a16 uint16) {
	panic(fmt.Sprintf("Missing implementation for rlAddr: %v", a16))
}

func (cpu *CPU) rla() {
	panic(fmt.Sprintf("Missing implementation for rla"))
}

func (cpu *CPU) rlc(r8 *uint8) {
	panic(fmt.Sprintf("Missing implementation for rlc: %v", r8))
}

func (cpu *CPU) rlca() {
	panic(fmt.Sprintf("Missing implementation for rlca"))
}

func (cpu *CPU) rlcAddr(a16 uint16) {
	panic(fmt.Sprintf("Missing implementation for rlcAddr: %v", a16))
}

func (cpu *CPU) rr(r8 *uint8) {
	panic(fmt.Sprintf("Missing implementation for rr: %v", r8))
}

func (cpu *CPU) rra() {
	panic(fmt.Sprintf("Missing implementation for rra"))
}

func (cpu *CPU) rrAddr(a16 uint16) {
	panic(fmt.Sprintf("Missing implementation for rrAddr: %v", a16))
}

func (cpu *CPU) rrc(r8 *uint8) {
	panic(fmt.Sprintf("Missing implementation for rrc: %v", r8))
}

func (cpu *CPU) rrca() {
	panic(fmt.Sprintf("Missing implementation for rrca"))
}

func (cpu *CPU) rrcAddr(a16 uint16) {
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

func (cpu *CPU) slaAddr(a16 uint16) {
	panic(fmt.Sprintf("Missing implementation for slaAddr: %v", a16))
}

func (cpu *CPU) sra(r8 *uint8) {
	panic(fmt.Sprintf("Missing implementation for sra: %v", r8))
}

func (cpu *CPU) sraAddr(a16 uint16) {
	panic(fmt.Sprintf("Missing implementation for sraAddr: %v", a16))
}

func (cpu *CPU) srl(r8 *uint8) {
	panic(fmt.Sprintf("Missing implementation for srl: %v", r8))
}

func (cpu *CPU) srlAddr(a16 uint16) {
	panic(fmt.Sprintf("Missing implementation for srlAddr: %v", a16))
}

func (cpu *CPU) swap(r8 *uint8) {
	panic(fmt.Sprintf("Missing implementation for swap: %v", r8))
}

func (cpu *CPU) swapAddr(a16 uint16) {
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

func (cpu *CPU) sbcAddr(a16 uint16) {
	panic(fmt.Sprintf("Missing implementation for sbcAddr: %v", a16))
}

func (cpu *CPU) sub(u8 uint8) {
	panic(fmt.Sprintf("Missing implementation for sub: %v", u8))
}

func (cpu *CPU) subAddr(a16 uint16) {
	panic(fmt.Sprintf("Missing implementation for subAddr: %v", a16))
}

func (cpu *CPU) xor(u8 uint8) {
	cpu.a = cpu.a ^ u8
	cpu.flags(z(cpu.a), false, false, false) // [Z 0 0 0]
}

func (cpu *CPU) xorAddr(a16 uint16) {
	panic(fmt.Sprintf("Missing implementation for xorAddr: %v", a16))
}

////////////////
////////////////
// OLD
////////////////
////////////////

// func (cpu *CPU) adc(mem mem.Memory, operand1, operand2 string, u8 uint8, r16 register16) (flags map[string]bool) {
// 	switch operand1 {
// 	case "A":
// 		a := cpu.a
// 		if operand2 == "d8" {
// 			cpu.a += u8
// 		} else {
// 			cpu.a += cpu.get8(mem, operand2)
// 		}
// 		halfCarry := h(a, cpu.a)
// 		carry := c(a, cpu.a)
// 		if cpu.isFlagSet(cFlag) {
// 			a = cpu.a
// 			cpu.a++
// 			halfCarry = halfCarry || h(a, cpu.a)
// 			carry = carry || c(a, cpu.a)
// 		}
// 		return map[string]bool{
// 			"Z": z(cpu.a),
// 			"H": halfCarry,
// 			"C": carry}
// 	default:
// 		panic(fmt.Sprintf("Missing implementation for adc: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
// 	}
// }

// func (cpu *CPU) and(mem mem.Memory, operand1, operand2 string, u8 uint8, r16 register16) (flags map[string]bool) {
// 	switch operand1 {
// 	case "d8":
// 		cpu.a &= u8
// 	default:
// 		cpu.a &= cpu.get8(mem, operand1)
// 	}
// 	return map[string]bool{"Z": cpu.a == 0}
// }

// func (cpu *CPU) cp(mem mem.Memory, operand1, operand2 string, u8 uint8, r16 register16) (flags map[string]bool) {
// 	var val uint8
// 	switch operand1 {
// 	case "d8":
// 		val = u8
// 	case "(HL)":
// 		val = mem.Read(cpu.get16(mem, "HL"))
// 	default:
// 		val = cpu.get8(mem, operand1)
// 	}
// 	return map[string]bool{
// 		"Z": val == cpu.a,
// 		"H": (val & 0xf) > (cpu.a & 0xf),
// 		"C": val > cpu.a}
// }

// func (cpu *CPU) cpl(mem mem.Memory, operand1, operand2 string, u8 uint8, r16 register16) (flags map[string]bool) {
// 	panic(fmt.Sprintf("Missing implementation for cpl: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
// }

// func (cpu *CPU) daa(mem mem.Memory, operand1, operand2 string, u8 uint8, r16 register16) (flags map[string]bool) {
// 	panic(fmt.Sprintf("Missing implementation for daa: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
// }

// func (cpu *CPU) dec(mem mem.Memory, operand1, operand2 string, u8 uint8, r16 register16) (flags map[string]bool) {
// 	val := cpu.get8(mem, operand1)
// 	val--
// 	cpu.set8(mem, operand1, val)
// 	return map[string]bool{
// 		"Z": val == 0,
// 		"H": (val & 0xf) == 0xf}
// }

// func (cpu *CPU) di(mem mem.Memory, operand1, operand2 string, u8 uint8, r16 register16) (flags map[string]bool) {
// 	cpu.ime = false
// 	return
// }

// func (cpu *CPU) ei(mem mem.Memory, operand1, operand2 string, u8 uint8, r16 register16) (flags map[string]bool) {
// 	cpu.ime = true
// 	return
// }

// func (cpu *CPU) jr(mem mem.Memory, operand1, operand2 string, u8 uint8, r16 register16) (flags map[string]bool) {
// 	switch operand1 {
// 	case "Z":
// 		if cpu.isFlagSet(zFlag) {
// 			cpu.pc = register16(int32(cpu.pc) + int32(int8(u8)))
// 		}
// 	case "NZ":
// 		if !cpu.isFlagSet(zFlag) {
// 			cpu.pc = register16(int32(cpu.pc) + int32(int8(u8)))
// 		}
// 	default:
// 		panic(fmt.Sprintf("Missing implementation for jr: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
// 	}
// 	return
// }

// func (cpu *CPU) ld(mem mem.Memory, operand1, operand2 string, u8 uint8, r16 register16) (flags map[string]bool) {
// 	if operand1 == "(a16)" {
// 		mem.Write(u16, cpu.get8(mem, operand2))
// 	} else if operand1 == "(C)" {
// 		c := cpu.get8(mem, "C")
// 		mem.Write(0xff00|register16(c), cpu.get8(mem, operand2))
// 	} else if operand2 == "(C)" {
// 		c := cpu.get8(mem, "C")
// 		cpu.set8(mem, operand1, mem.Read(0xff00|register16(c)))
// 	} else if operand1 == "(HL+)" {
// 		hl := cpu.get16(mem, "HL")
// 		mem.Write(hl, cpu.get8(mem, operand2))
// 		hl++
// 		cpu.set16(mem, "HL", hl)
// 	} else if operand1 == "(HL-)" {
// 		hl := cpu.get16(mem, "HL")
// 		mem.Write(hl, cpu.get8(mem, operand2))
// 		hl--
// 		cpu.set16(mem, "HL", hl)
// 	} else if operand2 == "(HL+)" {
// 		hl := cpu.get16(mem, "HL")
// 		cpu.set8(mem, operand1, mem.Read(hl))
// 		hl++
// 		cpu.set16(mem, "HL", hl)
// 	} else if operand2 == "(HL-)" {
// 		hl := cpu.get16(mem, "HL")
// 		cpu.set8(mem, operand1, mem.Read(hl))
// 		hl--
// 		cpu.set16(mem, "HL", hl)
// 	} else if len(operand2) == 1 {
// 		cpu.set8(mem, operand1, cpu.get8(mem, operand2))
// 	} else if operand2 == "d8" {
// 		cpu.set8(mem, operand1, u8)
// 	} else if operand2 == "d16" {
// 		cpu.set16(mem, operand1, u16)
// 	} else {
// 		cpu.set16(mem, operand1, cpu.get16(mem, operand2))
// 	}
// 	return
// }

// func (cpu *CPU) ldh(mem mem.Memory, operand1, operand2 string, u8 uint8, r16 register16) (flags map[string]bool) {
// 	if operand1 == "(a8)" && operand2 == "A" {
// 		mem.Write(0xff00+register16(int8(u8)), cpu.a)
// 	} else if operand1 == "A" && operand2 == "(a8)" {
// 		cpu.a = mem.Read(0xff00 + register16(int8(u8)))
// 	} else {
// 		panic(fmt.Sprintf("Missing implementation for ldh: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
// 	}
// 	return
// }
