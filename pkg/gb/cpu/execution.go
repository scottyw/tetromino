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
				cpu.rst(0x0040)()
				memory.IF &^= bit0
			case interrupts&bit1 > 0:
				// 0048 LCDC Status Interrupt Start Address
				cpu.rst(0x0048)()
				memory.IF &^= bit1
			case interrupts&bit2 > 0:
				// 0050 Timer OverflowInterrupt Start Address
				cpu.rst(0x0050)()
				memory.IF &^= bit2
			case interrupts&bit3 > 0:
				// 0058 Serial Transfer Completion Interrupt Start Address
				cpu.rst(0x0058)()
				memory.IF &^= bit3
			case interrupts&bit4 > 0:
				// 0060 High-to-Low of P10-P13 Interrupt Start Address
				cpu.rst(0x0060)()
				memory.IF &^= bit4
			}

			// Now push the PC
			cpu.push(memory, &cpu.m8b)()
			cpu.push(memory, &cpu.m8a)()
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

func (d *Dispatch) useAltMachineCycles(instruction uint8) bool {
	cpu := d.cpu
	switch instruction {

	// CALL 3 [24 12] [] 0xc4 NZ a16 196 false 6 3 <nil>
	// JP 3 [16 12] [] 0xc2 NZ a16 194 false 4 3 <nil>
	// JR 2 [12 8] [] 0x20 NZ r8 32 false 3 2 <nil>
	// RET 1 [20 8] [] 0xc0 NZ  192 false 5 2 <nil>
	case 0xc4, 0xc2, 0x20, 0xc0:
		return cpu.zf()

	// CALL 3 [24 12] [] 0xcc Z a16 204 false 6 3 <nil>
	// JP 3 [16 12] [] 0xca Z a16 202 false 4 3 <nil>
	// JR 2 [12 8] [] 0x28 Z r8 40 false 3 2 <nil>
	// RET 1 [20 8] [] 0xc8 Z  200 false 5 2 <nil>
	case 0xcc, 0xca, 0x28, 0xc8:
		return !cpu.zf()

	// CALL 3 [24 12] [] 0xd4 NC a16 212 false 6 3 <nil>
	// JP 3 [16 12] [] 0xd2 NC a16 210 false 4 3 <nil>
	// JR 2 [12 8] [] 0x30 NC r8 48 false 3 2 <nil>
	// RET 1 [20 8] [] 0xd0 NC  208 false 5 2 <nil>
	case 0xd4, 0xd2, 0x30, 0xd0:
		return cpu.cf()

	// CALL 3 [24 12] [] 0xdc C a16 220 false 6 3 <nil>
	// JP 3 [16 12] [] 0xda C a16 218 false 4 3 <nil>
	// JR 2 [12 8] [] 0x38 C r8 56 false 3 2 <nil>
	// RET 1 [20 8] [] 0xd8 C  216 false 5 2 <nil>
	case 0xdc, 0xda, 0x38, 0xd8:
		return !cpu.cf()

	}
	return false
}

// Every instruction is implemented as a list of steps that take one machine cycle each
func (d *Dispatch) peek() *[]func() {
	cpu := d.cpu
	memory := d.memory
	if cpu.pc >= 0x0100 {
		memory.DisableBIOS()
	}
	instructionByte := memory.Read(cpu.pc)
	// Mooneye uses the 0x40 instruction as a magic breakpoint
	// to indicate that a test rom has completed
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
	var steps []func()
	var value string
	if md.Prefixed {
		steps = d.prefix[md.Dispatch]
		if !cpu.haltbug {
			cpu.pc += 2
		} else {
			cpu.haltbug = false
		}
	} else {
		// Peek at the instruction arguments for debug purposes
		if cpu.debugCPU {
			switch md.Length {
			case 2:
				u8 := memory.Read(cpu.pc + 1)
				value = fmt.Sprintf("%02x", u8)
			case 3:
				u16 := uint16(memory.Read(cpu.pc+1)) | uint16(memory.Read(cpu.pc+2))<<8
				value = fmt.Sprintf("%04x", u16)
			}
		}

		// Reset any context from previous instructions
		cpu.u8a = 0
		cpu.u8b = 0
		cpu.m8a = 0
		cpu.m8b = 0

		// Get the steps associated with this instruction
		steps = d.normal[md.Dispatch]

		// Check for instructions that need to use the shorter alt machine cycle count
		if d.useAltMachineCycles(md.Dispatch) {
			steps = steps[:md.AltMachineCycles]
		}

		// Finally increment PC
		if !cpu.haltbug {
			cpu.pc++
		} else {
			cpu.haltbug = false
		}
	}
	if cpu.debugCPU {
		fmt.Printf("0x%04x: [%02x] %-12s | %-4s | a:%02x b:%02x c:%02x d:%02x e:%02x f:%02x h:%02x l:%02x sp:%04x\n",
			pc, md.Dispatch, fmt.Sprintf("%s %s %s", md.Mnemonic, md.Operand1, md.Operand2), value, cpu.a, cpu.b, cpu.c, cpu.d, cpu.e, cpu.f, cpu.h, cpu.l, cpu.sp)
	}
	return &steps
}

// ExecuteMachineCycle runs the CPU for one machine cycle
func (d *Dispatch) ExecuteMachineCycle() {
	cpu := d.cpu
	if d.stepIndex == len(*d.steps) {
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
