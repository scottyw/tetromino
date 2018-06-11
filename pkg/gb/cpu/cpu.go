package cpu

import (
	"fmt"

	"github.com/scottyw/tetromino/pkg/gb/mem"
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

	// State
	sp      uint16
	pc      uint16
	ime     bool
	halted  bool
	stopped bool

	// Hardware
	hwr   *mem.HardwareRegisters
	ticks int

	// Debug
	debugCPU         bool
	debugFlowControl bool
	debugJumps       bool
}

// NewCPU returns a CPU initialized as a Gameboy does on start
func NewCPU(hwr *mem.HardwareRegisters, debugCPU, debugFlowControl, debugJumps bool) *CPU {
	return &CPU{
		hwr:              hwr,
		debugCPU:         debugCPU,
		debugFlowControl: debugFlowControl,
		debugJumps:       debugJumps,
		ime:              true,
		a:                0x01,
		f:                0xb0,
		b:                0x00,
		c:                0x13,
		d:                0x00,
		e:                0xd8,
		h:                0x01,
		l:                0x4d,
		sp:               0xfffe,
		pc:               0x0100}
}

func (cpu CPU) String() string {
	return fmt.Sprintf("{ime:%v a:%02x b:%02x c:%02x d:%02x e:%02x f:%02x h:%02x l:%02x sp:%04x pc:%04x}",
		cpu.ime, cpu.a, cpu.b, cpu.c, cpu.d, cpu.e, cpu.f, cpu.h, cpu.l, cpu.sp, cpu.pc)
}

// Start the CPU again on button press
func (cpu *CPU) Start() {
	cpu.stopped = false
}

func (cpu *CPU) bc() uint16 {
	return uint16(cpu.b)<<8 + uint16(cpu.c)
}

func (cpu *CPU) de() uint16 {
	return uint16(cpu.d)<<8 + uint16(cpu.e)
}

func (cpu *CPU) af() uint16 {
	return uint16(cpu.a)<<8 + uint16(cpu.f)
}

func (cpu *CPU) hl() uint16 {
	return uint16(cpu.h)<<8 + uint16(cpu.l)
}

func (cpu *CPU) zf() bool {
	return cpu.f&zFlag > 0
}

func (cpu *CPU) nf() bool {
	return cpu.f&nFlag > 0
}

func (cpu *CPU) hf() bool {
	return cpu.f&hFlag > 0
}

func (cpu *CPU) cf() bool {
	return cpu.f&cFlag > 0
}

func (cpu *CPU) setZf(value bool) {
	if value {
		cpu.f |= zFlag
	} else {
		cpu.f &^= zFlag
	}
}

func (cpu *CPU) setNf(value bool) {
	if value {
		cpu.f |= nFlag
	} else {
		cpu.f &^= nFlag
	}
}

func (cpu *CPU) setHf(value bool) {
	if value {
		cpu.f |= hFlag
	} else {
		cpu.f &^= hFlag
	}
}

func (cpu *CPU) setCf(value bool) {
	if value {
		cpu.f |= cFlag
	} else {
		cpu.f &^= cFlag
	}
}

func hc8(a, b uint8) bool {
	return a&0x0f+b&0x0f > 0x0f
}
func c8(a, b uint8) bool {
	return int(a)+int(b) > 0xff
}

func hc16(a, b uint16) bool {
	return a&0x0fff+b&0x0fff > 0x0fff
}

func c16(a, b uint16) bool {
	return int(a)+int(b) > 0xffff
}

func hc8Sub(a, b uint8) bool {
	return int(a)&0x0f-int(b)&0x0f < 0
}

func c8Sub(a, b uint8) bool {
	return int(a)-int(b) < 0
}

func hc16Sub(a, b uint16) bool {
	return int(a)&0x0fff-int(b)&0x0fff < 0
}

func c16Sub(a, b uint16) bool {
	return int(a)-int(b) < 0
}

func (cpu *CPU) checkInterrupts(memory *mem.Memory) {
	interrupts := cpu.hwr.IE & cpu.hwr.IF
	if interrupts > 0 {
		cpu.halted = false
		if cpu.ime {
			cpu.ime = false
			switch {
			case interrupts&bit0 > 0:
				// 0040 Vertical Blank Interrupt Start Address
				if cpu.debugFlowControl {
					fmt.Printf("==== V-Blank interrupt ...\n")
				}
				cpu.rst(0x0040, memory)
				cpu.hwr.IF &^= bit0
			case interrupts&bit1 > 0:
				// 0048 LCDC Status Interrupt Start Address
				if cpu.debugFlowControl {
					fmt.Printf("==== LCDC Status interrupt ...\n")
				}
				cpu.rst(0x0048, memory)
				cpu.hwr.IF &^= bit1
			case interrupts&bit2 > 0:
				// 0050 Timer OverflowInterrupt Start Address
				if cpu.debugFlowControl {
					fmt.Printf("==== Timer Overflow interrupt ...\n")
				}
				cpu.rst(0x0050, memory)
				cpu.hwr.IF &^= bit2
			case interrupts&bit3 > 0:
				// 0058 Serial Transfer Completion Interrupt Start Address
				if cpu.debugFlowControl {
					fmt.Printf("==== Serial Transfer interrupt ...\n")
				}
				cpu.rst(0x0058, memory)
				cpu.hwr.IF &^= bit3
			case interrupts&bit4 > 0:
				// 0060 High-to-Low of P10-P13 Interrupt Start Address
				if cpu.debugFlowControl {
					fmt.Printf("==== High-to-Low Pin interrupt ...\n")
				}
				cpu.rst(0x0060, memory)
				cpu.hwr.IF &^= bit4
			}
		}
	}
}

func flagMetadata(i uint, flags []string) string {
	if len(flags) == 0 {
		return "-"
	}
	return flags[i]
}

func validateFlag(label string, i uint, f1, f2 uint8, im metadata) {
	bit := uint8(0x80) >> i
	switch flagMetadata(i, im.Flags) {
	case "-":
		if f1&bit != f2&bit {
			panic(fmt.Sprintf("%s flag invalid! Should not change: before=0x%02x after=0x%02x metadata=%v", label, f1, f2, im))
		}
	case "0":
		if f2&bit != 0 {
			panic(fmt.Sprintf("%s flag invalid! Should be reset: flags=0x%02x metadata=%v", label, f2, im))
		}
	case "1":
		if f2&bit == 0 {
			panic(fmt.Sprintf("%s flag invalid! Should be set: flags=0x%02x metadata=%v", label, f2, im))
		}
	}
}

func validateFlags(f1, f2 uint8, im metadata) {
	validateFlag("Z", 0, f1, f2, im)
	validateFlag("N", 1, f1, f2, im)
	validateFlag("H", 2, f1, f2, im)
	validateFlag("C", 3, f1, f2, im)
}

func (cpu *CPU) execute(mem *mem.Memory) int {
	// f := cpu.f // FIXME flag validation
	pc := cpu.pc
	instruction := mem.Read(cpu.pc)
	im := instructionMetadata[instruction]
	if im.Addr == "" {
		panic(fmt.Sprintf("Unknown instruction opcode: %v", instruction))
	}
	if instruction == 0xcb {
		instruction = mem.Read(cpu.pc + 1)
		im = prefixedInstructionMetadata[instruction]
		cpu.pc += 2
		cpu.dispatchPrefixedInstruction(mem, instruction)
		if cpu.debugCPU {
			fmt.Printf("0x%04x: [%02x] %-12s |      | a:%02x b:%02x c:%02x d:%02x e:%02x f:%02x h:%02x l:%02x sp:%04x\n",
				pc, instruction, fmt.Sprintf("%s %s %s", im.Mnemonic, im.Operand1, im.Operand2), cpu.a, cpu.b, cpu.c, cpu.d, cpu.e, cpu.f, cpu.h, cpu.l, cpu.sp)
		}
	} else {
		switch im.Length {
		case 1:
			cpu.pc++
			cpu.dispatchOneByteInstruction(mem, instruction)
			if cpu.debugCPU {
				fmt.Printf("0x%04x: [%02x] %-12s |      | a:%02x b:%02x c:%02x d:%02x e:%02x f:%02x h:%02x l:%02x sp:%04x\n",
					pc, instruction, fmt.Sprintf("%s %s %s", im.Mnemonic, im.Operand1, im.Operand2), cpu.a, cpu.b, cpu.c, cpu.d, cpu.e, cpu.f, cpu.h, cpu.l, cpu.sp)
			}
		case 2:
			u8 := mem.Read(cpu.pc + 1)
			cpu.pc += 2
			cpu.dispatchTwoByteInstruction(mem, instruction, u8)
			if cpu.debugCPU {
				fmt.Printf("0x%04x: [%02x] %-12s | %02x   | a:%02x b:%02x c:%02x d:%02x e:%02x f:%02x h:%02x l:%02x sp:%04x\n",
					pc, instruction, fmt.Sprintf("%s %s %s", im.Mnemonic, im.Operand1, im.Operand2), u8, cpu.a, cpu.b, cpu.c, cpu.d, cpu.e, cpu.f, cpu.h, cpu.l, cpu.sp)
			}
		case 3:
			u16 := uint16(mem.Read(cpu.pc+1)) | uint16(mem.Read(cpu.pc+2))<<8
			cpu.pc += 3
			cpu.dispatchThreeByteInstruction(mem, instruction, u16)
			if cpu.debugCPU {
				fmt.Printf("0x%04x: [%02x] %-12s | %04x | a:%02x b:%02x c:%02x d:%02x e:%02x f:%02x h:%02x l:%02x sp:%04x\n",
					pc, instruction, fmt.Sprintf("%s %s %s", im.Mnemonic, im.Operand1, im.Operand2), u16, cpu.a, cpu.b, cpu.c, cpu.d, cpu.e, cpu.f, cpu.h, cpu.l, cpu.sp)
			}
		}
	}
	// FIXME temporary flag validation
	// validateFlags(f, cpu.f, im)
	// FIXME - Most instructions have a single cycle count - handle the conditional ones later.
	return im.Cycles[0]
}

// Tick runs the CPU for one machine cycle i.e. 4 clock cycles
func (cpu *CPU) Tick(mem *mem.Memory) {
	if cpu.ticks == 0 {
		cpu.checkInterrupts(mem)
		if !cpu.halted && !cpu.stopped {
			cpu.ticks = cpu.execute(mem)
		}
	} else {
		cpu.ticks--
	}
}
