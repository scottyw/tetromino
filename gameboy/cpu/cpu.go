package cpu

import (
	"github.com/scottyw/tetromino/gameboy/interrupts"
	"github.com/scottyw/tetromino/gameboy/memory"
	"github.com/scottyw/tetromino/gameboy/oam"
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
	sp                     uint16
	pc                     uint16
	halted                 bool
	haltbug                bool
	stopped                bool
	interrupts             *interrupts.Interrupts
	oam                    *oam.OAM
	mapper                 *memory.Mapper
	currentInstruction     uint8
	currentSubinstructions []func()
	currentCycle           int
	currentEnded           func() bool
	currentMetadata        *metadata

	// Context
	u8a uint8 // 8-bit instruction argument
	u8b uint8 // Additional 8-bit instruction argument
	m8a uint8 // Cached copy of an 8-bit memory value for instructions that operate over it
	m8b uint8 // Additional cached copy of an 8-bit memory value for instructions that operate over it

	// Debug
	debugCPU               bool
	mooneyeDebugBreakpoint bool
}

// NewCPU returns a CPU initialized as a Gameboy does on start
func New(interrupts *interrupts.Interrupts, oam *oam.OAM, debugCPU bool, mapper *memory.Mapper) *CPU {
	return &CPU{
		interrupts: interrupts,
		oam:        oam,
		mapper:     mapper,
		debugCPU:   debugCPU,
		a:          0x01,
		f:          0xb0,
		b:          0x00,
		c:          0x13,
		d:          0x00,
		e:          0xd8,
		h:          0x01,
		l:          0x4d,
		sp:         0xfffe,
		pc:         0x0100,
	}
}

// Restart the CPU again on button press
func (cpu *CPU) Restart() {
	cpu.stopped = false
}

func (cpu *CPU) bc() uint16 {
	return uint16(cpu.b)<<8 + uint16(cpu.c)
}

func (cpu *CPU) de() uint16 {
	return uint16(cpu.d)<<8 + uint16(cpu.e)
}

func (cpu *CPU) hl() uint16 {
	return uint16(cpu.h)<<8 + uint16(cpu.l)
}

func (cpu *CPU) u16() uint16 {
	return uint16(cpu.u8b)<<8 + uint16(cpu.u8a)
}

// CheckMooneye checks if the magic Mooneye breakpoint was hit and if so returns
// the register values needed to see whether the test passed
func (cpu *CPU) CheckMooneye() []uint8 {
	if cpu.mooneyeDebugBreakpoint {
		return []uint8{cpu.a, cpu.b, cpu.c, cpu.d, cpu.e, cpu.h, cpu.l}
	}
	return nil
}

func (cpu *CPU) mooneye() {
	// Mooneye uses this instruction (0x40) as a magic breakpoint
	// to indicate that a test rom has completed
	cpu.mooneyeDebugBreakpoint = true
}

// Read an 8-bit instruction argument
func (cpu *CPU) readParamA() {
	cpu.u8a = cpu.mapper.Read(cpu.pc)
	cpu.pc++
}

// Read an additonal 8-bit instruction argument
func (cpu *CPU) readParamB() {
	cpu.u8b = cpu.mapper.Read(cpu.pc)
	cpu.pc++
}
