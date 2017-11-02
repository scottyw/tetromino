package cpu

import (
	"fmt"

	"github.com/scottyw/goomba/mem"
)

const (
	bit0 uint8 = 1 << iota
	bit1 uint8 = 1 << iota
	bit2 uint8 = 1 << iota
	bit3 uint8 = 1 << iota
	bit4 uint8 = 1 << iota
	bit5 uint8 = 1 << iota
	bit6 uint8 = 1 << iota
	bit7 uint8 = 1 << iota

	zFlag = bit7
	nFlag = bit6
	hFlag = bit5
	cFlag = bit4
)

var bits = [8]uint8{bit0, bit1, bit2, bit3, bit4, bit5, bit6, bit7}

var flags = [4]uint8{zFlag, nFlag, hFlag, cFlag}

var opcodes [256]opcodeMetadata

var prefixedOpcodes [256]opcodeMetadata

var cycles int

// CPU stores the internal CPU state
type CPU struct {
	// Registers
	a uint8
	b uint8
	c uint8
	d uint8
	e uint8
	f uint8
	h uint8
	l uint8

	// Flags
	zf bool
	nf bool
	hf bool
	cf bool

	// State
	sp  uint16
	pc  uint16
	ime bool
}

// NewCPU returns a CPU initialized as a Gameboy does on start
func NewCPU() CPU {
	return CPU{
		ime: true,
		a:   0x01,
		f:   0xb0,
		b:   0x00,
		c:   0x13,
		d:   0x00,
		e:   0xd8,
		h:   0x01,
		l:   0x4d,
		sp:  0xfffe,
		pc:  0x0100}
}

func (cpu CPU) String() string {
	return fmt.Sprintf("{ime:%v a:%02x b:%02x c:%02x d:%02x e:%02x f:%02x h:%02x l:%02x sp:%04x pc:%04x}",
		cpu.ime, cpu.a, cpu.b, cpu.c, cpu.d, cpu.e, cpu.f, cpu.h, cpu.l, cpu.sp, cpu.pc)
}

func (cpu *CPU) get8(mem mem.Memory, name string) uint8 {
	switch name {
	case "A":
		return cpu.a
	case "B":
		return cpu.b
	case "C":
		return cpu.c
	case "D":
		return cpu.d
	case "E":
		return cpu.e
	case "H":
		return cpu.h
	case "L":
		return cpu.l
	case "(HL)":
		return mem.Read(cpu.get16(mem, "HL"))
	case "(a16)":
		return mem.Read(cpu.get16(mem, "MOOO"))
	default:
		panic(fmt.Sprintf("get8: %v", name))
	}
}

func (cpu *CPU) get16(mem mem.Memory, name string) uint16 {
	switch name {
	case "BC":
		return uint16(cpu.b)<<8 + uint16(cpu.c)
	case "DE":
		return uint16(cpu.d)<<8 + uint16(cpu.e)
	case "HL":
		return uint16(cpu.h)<<8 + uint16(cpu.l)
	default:
		panic(fmt.Sprintf("get16: %v", name))
	}
}

func (cpu *CPU) set8(mem mem.Memory, name string, val uint8) {
	switch name {
	case "A":
		cpu.a = val
	case "B":
		cpu.b = val
	case "C":
		cpu.c = val
	case "D":
		cpu.d = val
	case "E":
		cpu.e = val
	case "H":
		cpu.h = val
	case "L":
		cpu.l = val
	case "(HL)":
		mem.Write(cpu.get16(mem, "HL"), val)
	default:
		panic(fmt.Sprintf("set8: %v", name))
	}
}

func (cpu *CPU) set16(mem mem.Memory, name string, val uint16) {
	switch name {
	case "BC":
		cpu.b = uint8(val >> 8)
		cpu.c = uint8(val)
	case "DE":
		cpu.d = uint8(val >> 8)
		cpu.e = uint8(val)
	case "HL":
		cpu.h = uint8(val >> 8)
		cpu.l = uint8(val)
	case "SP":
		cpu.sp = val
	default:
		panic(fmt.Sprintf("set16: %v", name))
	}
}

func (cpu *CPU) checkInterrupts(mem mem.Memory) {
	if cpu.ime {
		interrupts := mem.Read(0xffff) & mem.Read(0xff0f)
		if interrupts > 0 {
			cpu.ime = false
			mem.Write(cpu.sp, uint8(cpu.pc))
			cpu.sp--
			mem.Write(cpu.sp, uint8(cpu.pc>>8))
			cpu.sp--
			switch {
			case interrupts&0x01 > 0:
				cpu.pc = 0x0040                     // V-Blank
				mem.Write(0xff0f, interrupts&^0x01) // Reset IF
			case interrupts&0x02 > 0:
				cpu.pc = 0x0048                     // LCDC status
				mem.Write(0xff0f, interrupts&^0x02) // Reset IF
			case interrupts&0x04 > 0:
				cpu.pc = 0x0050                     // Timer Overflow
				mem.Write(0xff0f, interrupts&^0x04) // Reset IF
			case interrupts&0x08 > 0:
				cpu.pc = 0x0058                     // Serial Transfer
				mem.Write(0xff0f, interrupts&^0x08) // Reset IF
			case interrupts&0x10 > 0:
				cpu.pc = 0x0060                     // Hi-Lo of P10-P13
				mem.Write(0xff0f, interrupts&^0x10) // Reset IF
			}
		}
	}
}

func z(new uint8) bool {
	return new == 0
}

func h(old, new uint8) bool {
	return old&0xf > new&0xf
}

func c(old, new uint8) bool {
	return old > new
}

func setBitByPos(b uint8, pos uint8) {
	b |= bits[pos]
}

func setBitByPattern(b uint8, pattern uint8) {
	b |= pattern
}

func resetBitByPos(b uint8, pos uint8) {
	b &^= bits[pos]
}

func resetBitByPattern(b uint8, pattern uint8) {
	b &^= pattern
}

// Return true if the flag is set and false if not
func (cpu *CPU) isFlagSet(flag uint8) bool {
	return cpu.f&flag > 0
}

func (cpu *CPU) setFlag(flag uint8) {
	cpu.f |= flag
}

func (cpu *CPU) resetFlag(flagBit uint8) {
	cpu.f &^= flagBit
}

func (cpu *CPU) execute(mem mem.Memory) {
	defer mem.MemoryDump()
	cpu.checkInterrupts(mem)
	instruction := mem.Read(cpu.pc)
	opcode := opcodes[instruction]
	switch opcode.Length {
	case 1:
		fmt.Printf("0x%02x : %v\n", cpu.pc, opcode)
		cpu.pc++
		cpu.dispatch(instruction)
	case 2:
		u8 := mem.Read(cpu.pc + 1)
		fmt.Printf("0x%02x : %v u8=0x%02x\n", cpu.pc, opcode, u8)
		cpu.pc += 2
		cpu.dispatch8(instruction, u8)
	case 3:
		u16 := uint16(mem.Read(cpu.pc+1)) | uint16(mem.Read(cpu.pc+2))<<8
		fmt.Printf("0x%02x : %v u8=0x%04x\n", cpu.pc, opcode, u16)
		cpu.pc += 3
		cpu.dispatch16(instruction, u16)
	}
	// FIXME - Most instructions have a single cycle count - handle the conditional ones later.
	cycles = opcode.Cycles[0]
}

// Tick runs the CPU for one machine cycle i.e. 4 clock cycles
func (cpu *CPU) Tick(mem mem.Memory) {
	if cycles == 0 {
		cpu.execute(mem)
	}
	cycles--
}
