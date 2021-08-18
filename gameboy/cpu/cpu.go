package cpu

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

var instructionMetadata [256]*metadata

var prefixedInstructionMetadata [256]*metadata

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
	haltbug bool
	stopped bool

	// Context
	u8a uint8 // 8-bit instruction argument
	u8b uint8 // Additional 8-bit instruction argument
	m8a uint8 // Cached copy of an 8-bit memory value for instructions that operate over it
	m8b uint8 // Additional cached copy of an 8-bit memory value for instructions that operate over it

	// Debug
	debugCPU bool
}

// NewCPU returns a CPU initialized as a Gameboy does on start
func NewCPU(debugCPU bool) *CPU {
	return &CPU{
		debugCPU: debugCPU,
		ime:      true,
		a:        0x01,
		f:        0xb0,
		b:        0x00,
		c:        0x13,
		d:        0x00,
		e:        0xd8,
		h:        0x01,
		l:        0x4d,
		sp:       0xfffe,
		pc:       0x0100,
	}
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
