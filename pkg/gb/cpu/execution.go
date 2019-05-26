package cpu

import (
	"fmt"
)

func (d *Dispatch) handleInterrupt() func() {
	cpu := d.cpu
	memory := d.memory
	return func() {
		if cpu.ime {
			interrupts := memory.IE & memory.IF & 0x1f
			cpu.ime = false
			switch {
			case interrupts&bit0 > 0:
				// 0040 Vertical Blank Interrupt Start Address
				cpu.rst(0x0040, memory)
				memory.IF &^= bit0
			case interrupts&bit1 > 0:
				// 0048 LCDC Status Interrupt Start Address
				cpu.rst(0x0048, memory)
				memory.IF &^= bit1
			case interrupts&bit2 > 0:
				// 0050 Timer OverflowInterrupt Start Address
				cpu.rst(0x0050, memory)
				memory.IF &^= bit2
			case interrupts&bit3 > 0:
				// 0058 Serial Transfer Completion Interrupt Start Address
				cpu.rst(0x0058, memory)
				memory.IF &^= bit3
			case interrupts&bit4 > 0:
				// 0060 High-to-Low of P10-P13 Interrupt Start Address
				cpu.rst(0x0060, memory)
				memory.IF &^= bit4
			}
		}
	}
}

func (d *Dispatch) checkInterrupts() *[]func() {
	cpu := d.cpu
	memory := d.memory
	var length int
	interrupts := memory.IE & memory.IF & 0x1f
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
	steps := make([]func(), length)
	for i := range steps {
		steps[i] = nop
	}
	steps[length-1] = d.handleInterrupt()
	return &steps
}

// Every instruction is implemented as a list of steps that take one machine cycle each
func (d *Dispatch) peek() *[]func() {
	cpu := d.cpu
	memory := d.memory
	instructionByte := memory.Read(cpu.pc)
	// Mooneye uses the 0x40 instruction as a magic breakpoint
	// to indicate that a test rom has completeed
	if instructionByte == 0x40 {
		d.Mooneye = true
	}
	md := instructionMetadata[instructionByte]
	if md.Addr == "" {
		panic(fmt.Sprintf("Unknown instruction opcode: %v", md))
	}
	if instructionByte == 0xcb {
		instructionByte = memory.Read(cpu.pc + 1)
		md = prefixedInstructionMetadata[instructionByte]
	}
	pc := cpu.pc
	var steps *[]func()
	var value string
	if md.Prefixed {
		cpu.pc += 2
		steps = &d.prefix[md.Dispatch]
	} else {
		// Reset any context from previous instructions
		cpu.u8a = 0
		cpu.u8b = 0
		cpu.m8a = 0
		cpu.m8b = 0
		if cpu.debugCPU {
			// Peek at the instruction arguments for debug purposes
			switch md.Length {
			case 2:
				u8 := memory.Read(cpu.pc + 1)
				value = fmt.Sprintf("%02x", u8)
			case 3:
				u16 := uint16(memory.Read(cpu.pc+1)) | uint16(memory.Read(cpu.pc+2))<<8
				value = fmt.Sprintf("%04x", u16)
			}
		}
		cpu.pc++
		steps = &d.normal[md.Dispatch]
	}
	if md.AltMachineCycles != 0 {
		d.altStepIndex = md.AltMachineCycles
	} else {
		d.altStepIndex = 0
	}
	if cpu.debugCPU {
		fmt.Printf("0x%04x: [%02x] %-12s | %-4s | a:%02x b:%02x c:%02x d:%02x e:%02x f:%02x h:%02x l:%02x sp:%04x\n",
			pc, md.Dispatch, fmt.Sprintf("%s %s %s", md.Mnemonic, md.Operand1, md.Operand2), value, cpu.a, cpu.b, cpu.c, cpu.d, cpu.e, cpu.f, cpu.h, cpu.l, cpu.sp)
	}
	return steps
}

// ExecuteMachineCycle runs the CPU for one machine cycle
func (d *Dispatch) ExecuteMachineCycle() {
	cpu := d.cpu
	if d.stepIndex == len(*d.steps) ||
		(cpu.altTicks && d.stepIndex == d.altStepIndex) {
		cpu.altTicks = false
		var steps *[]func()
		if !d.handlingInterrupt {
			steps = d.checkInterrupts()
			if cpu.halted || cpu.stopped {
				return
			}
			if steps != nil {
				d.handlingInterrupt = true
			}
		} else {
			d.handlingInterrupt = false
		}
		if steps == nil {
			steps = d.peek()
		}
		d.stepIndex = 0
		d.steps = steps
	}
	step := (*d.steps)[d.stepIndex]
	step()
	d.stepIndex++
}
