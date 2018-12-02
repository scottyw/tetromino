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
	sp       uint16
	pc       uint16
	ime      bool
	halted   bool
	stopped  bool
	altTicks bool

	// Hardware
	hwr               *mem.HardwareRegisters
	steps             *[]func(*CPU, *mem.Memory)
	stepIndex         int
	altStepIndex      int
	handlingInterrupt bool

	// Debug
	debugCPU         bool
	debugFlowControl bool
	debugJumps       bool
	validateFlags    bool
	Mooneye          bool
}

// NewCPU returns a CPU initialized as a Gameboy does on start
func NewCPU(hwr *mem.HardwareRegisters, debugCPU, debugFlowControl, debugJumps bool) *CPU {
	initialSteps := []func(*CPU, *mem.Memory){}
	return &CPU{
		hwr:              hwr,
		steps:            &initialSteps,
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

// Start the CPU again on button press
func (cpu *CPU) Start() {
	cpu.stopped = false
}

// A returns the value of register a
func (cpu *CPU) A() uint8 {
	return cpu.a
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

func (cpu *CPU) checkInterrupts(memory *mem.Memory) *[]func(*CPU, *mem.Memory) {
	var length int
	interrupts := cpu.hwr.IE & cpu.hwr.IF & 0x1f
	if interrupts > 0 {
		if cpu.halted {
			cpu.halted = false
			length++
		}
		if cpu.ime {
			length += 5
		}
	}
	if length == 0 {
		return nil
	}
	steps := make([]func(*CPU, *mem.Memory), length)
	for i := range steps {
		steps[i] = (*CPU).nopStep
	}
	steps[length-1] = (*CPU).interruptStep
	return &steps
}

// Every instruction is implemented as a list of steps that take one machine cycle each
func (cpu *CPU) peek(m *mem.Memory) *[]func(*CPU, *mem.Memory) {
	instructionByte := m.Read(cpu.pc)
	md := instructionMetadata[instructionByte]
	if md.Addr == "" {
		panic(fmt.Sprintf("Unknown instruction opcode: %v", md))
	}
	if instructionByte == 0xcb {
		instructionByte = m.Read(cpu.pc + 1)
		md = prefixedInstructionMetadata[instructionByte]
	}
	// The step lists should be statically defined some place
	// For now we manufacture a list of the right length and fill it with nops
	// The last step is a single monolithic "do everything" step
	// This allows gradual migration of instructions to the new architecture
	steps := make([]func(*CPU, *mem.Memory), md.MachineCycles)
	for i := range steps {
		steps[i] = (*CPU).nopStep
	}
	if md.AltMachineCycles == 0 {
		steps[md.MachineCycles-1] = monolithStep(md)
		cpu.altStepIndex = 0
	} else {
		steps[md.AltMachineCycles-1] = monolithStep(md)
		cpu.altStepIndex = md.AltMachineCycles
	}
	return &steps
}

// FIXME used to debug timing but should be removed
var ticks uint32
var lastTicks uint32

// ExecuteMachineCycle runs the CPU for one machine cycle
func (cpu *CPU) ExecuteMachineCycle(m *mem.Memory) {
	if cpu.stepIndex == len(*cpu.steps) ||
		(cpu.altTicks && cpu.stepIndex == cpu.altStepIndex) {
		lastTicks = ticks
		cpu.altTicks = false
		var steps *[]func(*CPU, *mem.Memory)
		if !cpu.handlingInterrupt {
			steps = cpu.checkInterrupts(m)
			if cpu.halted || cpu.stopped {
				return
			}
			if steps != nil {
				cpu.handlingInterrupt = true
			}
		} else {
			cpu.handlingInterrupt = false
		}
		if steps == nil {
			steps = cpu.peek(m)
		}
		cpu.stepIndex = 0
		cpu.steps = steps
	}
	step := (*cpu.steps)[cpu.stepIndex]
	step(cpu, m)
	cpu.stepIndex++
	ticks += 4
}
