package cpu

import (
	"fmt"

	"github.com/scottyw/goomba/mem"
)

func (cpu *CPU) add(value uint8) {
	a := cpu.a
	cpu.a += value
	cpu.zf = z(cpu.a)
	cpu.hf = h(a, cpu.a)
	cpu.cf = c(a, cpu.a)
}

func (cpu *CPU) jp(kind string, u16 uint16) {
	switch kind {
	case "":
		cpu.pc = u16
	default:
		panic(fmt.Sprintf("Missing implementation for jp: %v %v", kind, u16))
	}
}

func (cpu *CPU) nop() {
	// Do nothing
	return
}

////////////////
////////////////
// OLD
////////////////
////////////////

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

func (cpu *CPU) and(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	switch operand1 {
	case "d8":
		cpu.a &= u8
	default:
		cpu.a &= cpu.get8(mem, operand1)
	}
	return map[string]bool{"Z": cpu.a == 0}
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

func (cpu *CPU) di(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	cpu.ime = false
	return
}

func (cpu *CPU) ei(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	cpu.ime = true
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

func (cpu *CPU) xor(mem mem.Memory, operand1, operand2 string, u8 uint8, u16 uint16) (flags map[string]bool) {
	cpu.a = cpu.a ^ cpu.get8(mem, operand1)
	return map[string]bool{"Z": cpu.a == 0}
}
