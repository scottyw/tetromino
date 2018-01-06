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

var instructionMetadata [256]metadata

var prefixedInstructionMetadata [256]metadata

var cycles int

type register16 interface {
	Get() uint16
	Set(uint16)
}

// CPU stores the internal CPU state
type CPU struct {
	// 8-bit registers
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
	return fmt.Sprintf("{ime:%v a:%02x b:%02x c:%02x d:%02x e:%02x f:%02x h:%02x l:%02x sp:%04x pc:%04x zf:%v nf:%v hf:%v cf:%v}",
		cpu.ime, cpu.a, cpu.b, cpu.c, cpu.d, cpu.e, cpu.f, cpu.h, cpu.l, cpu.sp, cpu.pc, cpu.zf, cpu.nf, cpu.hf, cpu.cf)
}

func (cpu *CPU) bc() register16 {
	return newRegister16(&cpu.b, &cpu.c)
}

func (cpu *CPU) de() register16 {
	return newRegister16(&cpu.d, &cpu.e)
}

func (cpu *CPU) af() register16 {
	return newRegister16(&cpu.a, &cpu.f)
}

func (cpu *CPU) hl() register16 {
	return newRegister16(&cpu.h, &cpu.l)
}

func z(new uint8) bool {
	return new == 0
}

func h(old, new uint8) bool {
	return old&0xf > new&0xf
}

func h16(old, new uint16) bool {
	return h(uint8(old>>4), uint8(new>>4))
}

func c(old, new uint8) bool {
	return old > new
}

func c16(old, new uint16) bool {
	return c(uint8(old>>4), uint8(new>>4))
}

func (cpu *CPU) flags(zf, nf, hf, cf bool) {
	cpu.zf = zf
	cpu.nf = nf
	cpu.hf = hf
	cpu.cf = cf
}

func (cpu *CPU) isFlagSet(flag uint8) bool {
	return cpu.f&uint8(flag) > 0
}

func (cpu *CPU) checkInterrupts(mem mem.Memory) {
	if cpu.ime {
		interrupts := *mem.Read(0xffff) & *mem.Read(0xff0f)
		if interrupts > 0 {
			cpu.ime = false
			*mem.Read(cpu.sp) = uint8(cpu.pc)
			cpu.sp--
			*mem.Read(cpu.sp) = uint8(cpu.pc >> 8)
			cpu.sp--
			switch {
			case interrupts&0x01 > 0:
				cpu.pc = 0x0040                        // V-Blank
				*mem.Read(0xff0f) = interrupts &^ 0x01 // Reset IF
			case interrupts&0x02 > 0:
				cpu.pc = 0x0048                        // LCDC status
				*mem.Read(0xff0f) = interrupts &^ 0x02 // Reset IF
			case interrupts&0x04 > 0:
				cpu.pc = 0x0050                        // Timer Overflow
				*mem.Read(0xff0f) = interrupts &^ 0x04 // Reset IF
			case interrupts&0x08 > 0:
				cpu.pc = 0x0058                        // Serial Transfer
				*mem.Read(0xff0f) = interrupts &^ 0x08 // Reset IF
			case interrupts&0x10 > 0:
				cpu.pc = 0x0060                        // Hi-Lo of P10-P13
				*mem.Read(0xff0f) = interrupts &^ 0x10 // Reset IF
			}
		}
	}
}

func (cpu *CPU) execute(mem mem.Memory) {
	defer mem.GenerateCrashReport()
	cpu.checkInterrupts(mem)
	instruction := *mem.Read(cpu.pc)
	im := instructionMetadata[instruction]
	if instruction == 0xcb {
		instruction := *mem.Read(cpu.pc + 1)
		fmt.Printf("0xcb%02x : %v\n", cpu.pc, im)
		cpu.pc += 2
		cpu.dispatchPrefixedInstruction(mem, instruction)
	} else {
		switch im.Length {
		case 1:
			fmt.Printf("0x%02x : %v\n", cpu.pc, im)
			cpu.pc++
			cpu.dispatchOneByteInstruction(mem, instruction)
		case 2:
			u8 := *mem.Read(cpu.pc + 1)
			fmt.Printf("0x%02x : %v u8=0x%02x\n", cpu.pc, im, u8)
			cpu.pc += 2
			cpu.dispatchTwoByteInstruction(mem, instruction, u8)
		case 3:
			u16 := uint16(*mem.Read(cpu.pc + 1)) | uint16(*mem.Read(cpu.pc + 2))<<8
			fmt.Printf("0x%02x : %v u8=0x%04x\n", cpu.pc, im, u16)
			cpu.pc += 3
			cpu.dispatchThreeByteInstruction(mem, instruction, u16)
		}
	}
	// FIXME - Most instructions have a single cycle count - handle the conditional ones later.
	cycles = im.Cycles[0]
}

// Tick runs the CPU for one machine cycle i.e. 4 clock cycles
func (cpu *CPU) Tick(mem mem.Memory) {
	if cycles == 0 {
		cpu.execute(mem)
	}
	cycles--
}
