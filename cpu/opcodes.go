package cpu

import (
	"fmt"

	"github.com/scottyw/goomba/mem"
)

var mapping = map[string]func(*CPU, mem.Memory, string, string, uint8, uint16) map[string]bool{
	"ADC":    (*CPU).adc,
	"ADD":    (*CPU).add,
	"ADD16":  (*CPU).add16,
	"AND":    (*CPU).and,
	"BIT":    (*CPU).bit,
	"CALL":   (*CPU).call,
	"CCF":    (*CPU).ccf,
	"CP":     (*CPU).cp,
	"CPL":    (*CPU).cpl,
	"DAA":    (*CPU).daa,
	"DEC":    (*CPU).dec,
	"DEC16":  (*CPU).dec16,
	"DI":     (*CPU).di,
	"EI":     (*CPU).ei,
	"HALT":   (*CPU).halt,
	"INC":    (*CPU).inc,
	"INC16":  (*CPU).inc16,
	"JR":     (*CPU).jr,
	"JP":     (*CPU).jp,
	"LD":     (*CPU).ld,
	"LD16":   (*CPU).ld16,
	"LDH":    (*CPU).ldh,
	"NOP":    (*CPU).nop,
	"OR":     (*CPU).or,
	"POP":    (*CPU).pop,
	"POP16":  (*CPU).pop16,
	"PREFIX": (*CPU).prefix,
	"PUSH":   (*CPU).push,
	"PUSH16": (*CPU).push16,
	"RES":    (*CPU).res,
	"RET":    (*CPU).ret,
	"RETI":   (*CPU).reti,
	"RL":     (*CPU).rl,
	"RLA":    (*CPU).rla,
	"RLC":    (*CPU).rlc,
	"RLCA":   (*CPU).rlca,
	"RR":     (*CPU).rr,
	"RRA":    (*CPU).rra,
	"RRC":    (*CPU).rrc,
	"RRCA":   (*CPU).rrca,
	"RST":    (*CPU).rst,
	"SBC":    (*CPU).sbc,
	"SCF":    (*CPU).scf,
	"SET":    (*CPU).set,
	"SLA":    (*CPU).sla,
	"SRA":    (*CPU).sra,
	"SRL":    (*CPU).srl,
	"STOP":   (*CPU).stop,
	"SUB":    (*CPU).sub,
	"SWAP":   (*CPU).swap,
	"XOR":    (*CPU).xor}

func z(new uint8) bool {
	return new == 0
}

func h(old, new uint8) bool {
	return old&0xf > new&0xf
}

func c(old, new uint8) bool {
	return old > new
}

func (cpu *CPU) adc(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	switch operand1 {
	case "A":
		a := cpu.a
		if operand2 == "d8" {
			cpu.a += u8
		} else {
			cpu.a += cpu.get8(mem, operand2)
		}
		halfCarry := h(a, cpu.a)
		carry := c(a, cpu.a)
		if cpu.isFlagSet(cFlag) {
			a = cpu.a
			cpu.a++
			halfCarry = halfCarry || h(a, cpu.a)
			carry = carry || c(a, cpu.a)
		}
		return map[string]bool{
			"Z": z(cpu.a),
			"H": halfCarry,
			"C": carry}
	default:
		panic(fmt.Sprintf("Missing implementation for adc: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
	}
}

func (cpu *CPU) add(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	switch operand1 {
	case "A":
		a := cpu.a
		if operand2 == "d8" {
			cpu.a += u8
		} else {
			cpu.a += cpu.get8(mem, operand2)
		}
		return map[string]bool{
			"Z": z(cpu.a),
			"H": h(a, cpu.a),
			"C": c(a, cpu.a)}
	default:
		panic(fmt.Sprintf("Missing implementation for add: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
	}
}

func (cpu *CPU) add16(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for add16: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) and(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	switch operand1 {
	case "d8":
		cpu.a &= u8
	default:
		cpu.a &= cpu.get8(mem, operand1)
	}
	return map[string]bool{"Z": cpu.a == 0}
}

func (cpu *CPU) bit(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for bit: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) call(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for call: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) ccf(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for ccf: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) cp(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	var val uint8
	switch operand1 {
	case "d8":
		val = u8
	case "(HL)":
		val = mem.Read(cpu.get16(mem, "HL"))
	default:
		val = cpu.get8(mem, operand1)
	}
	return map[string]bool{
		"Z": val == cpu.a,
		"H": (val & 0xf) > (cpu.a & 0xf),
		"C": val > cpu.a}
}

func (cpu *CPU) cpl(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for cpl: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) daa(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for daa: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) dec(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	val := cpu.get8(mem, operand1)
	val--
	cpu.set8(mem, operand1, val)
	return map[string]bool{
		"Z": val == 0,
		"H": (val & 0xf) == 0xf}
}

func (cpu *CPU) dec16(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for dec16: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) di(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	cpu.ime = false
	return
}

func (cpu *CPU) ei(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	cpu.ime = true
	return
}

func (cpu *CPU) halt(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for halt: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) inc(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for inc: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) inc16(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for inc16: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) jp(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	switch operand1 {
	case "a16":
		cpu.pc = u16
	default:
		panic(fmt.Sprintf("Missing implementation for jp: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
	}
	return
}

func (cpu *CPU) jr(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	switch operand1 {
	case "Z":
		if cpu.isFlagSet(zFlag) {
			cpu.pc = uint16(int32(cpu.pc) + int32(int8(u8)))
		}
	case "NZ":
		if !cpu.isFlagSet(zFlag) {
			cpu.pc = uint16(int32(cpu.pc) + int32(int8(u8)))
		}
	default:
		panic(fmt.Sprintf("Missing implementation for jr: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
	}
	return
}

func (cpu *CPU) ld(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	if operand1 == "(a16)" {
		mem.Write(u16, cpu.get8(mem, operand2))
	} else if operand1 == "(C)" {
		c := cpu.get8(mem, "C")
		mem.Write(0xff00|uint16(c), cpu.get8(mem, operand2))
	} else if operand2 == "(C)" {
		c := cpu.get8(mem, "C")
		cpu.set8(mem, operand1, mem.Read(0xff00|uint16(c)))
	} else if operand1 == "(HL+)" {
		hl := cpu.get16(mem, "HL")
		mem.Write(hl, cpu.get8(mem, operand2))
		hl++
		cpu.set16(mem, "HL", hl)
	} else if operand1 == "(HL-)" {
		hl := cpu.get16(mem, "HL")
		mem.Write(hl, cpu.get8(mem, operand2))
		hl--
		cpu.set16(mem, "HL", hl)
	} else if operand2 == "(HL+)" {
		hl := cpu.get16(mem, "HL")
		cpu.set8(mem, operand1, mem.Read(hl))
		hl++
		cpu.set16(mem, "HL", hl)
	} else if operand2 == "(HL-)" {
		hl := cpu.get16(mem, "HL")
		cpu.set8(mem, operand1, mem.Read(hl))
		hl--
		cpu.set16(mem, "HL", hl)
	} else if len(operand2) == 1 {
		cpu.set8(mem, operand1, cpu.get8(mem, operand2))
	} else if operand2 == "d8" {
		cpu.set8(mem, operand1, u8)
	} else if operand2 == "d16" {
		cpu.set16(mem, operand1, u16)
	} else {
		cpu.set16(mem, operand1, cpu.get16(mem, operand2))
	}
	return
}

func (cpu *CPU) ld16(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for ld16: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) ldh(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	if operand1 == "(a8)" && operand2 == "A" {
		mem.Write(0xff00+uint16(int8(u8)), cpu.a)
	} else if operand1 == "A" && operand2 == "(a8)" {
		cpu.a = mem.Read(0xff00 + uint16(int8(u8)))
	} else {
		panic(fmt.Sprintf("Missing implementation for ldh: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
	}
	return
}

func (cpu *CPU) nop(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	// Do nothing
	return
}

func (cpu *CPU) or(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for or: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) pop(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for pop: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) pop16(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for pop16: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) prefix(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for prefix: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) push(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for push: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) push16(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for push16: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) res(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for res: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) ret(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for ret: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) reti(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for reti: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) rl(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for rl: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) rla(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for rla: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) rlc(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for rlc: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) rlca(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for rlca: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) rr(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for rr: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) rra(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for rra: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) rrc(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for rrc: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) rrca(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for rrca: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) rst(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for rst: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) sbc(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for sbc: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) scf(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for scf: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) set(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for set: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) sla(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for sla: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) sra(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for sra: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) srl(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for srl: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) stop(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for stop: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) sub(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for sub: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) swap(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	panic(fmt.Sprintf("Missing implementation for swap: op1=%v op2=%v u8=%v u16=%v", operand1, operand2, u8, u16))
}

func (cpu *CPU) xor(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	cpu.a = cpu.a ^ cpu.get8(mem, operand1)
	return map[string]bool{"Z": cpu.a == 0}
}
