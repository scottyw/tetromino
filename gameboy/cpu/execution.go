package cpu

import "fmt"

func (cpu *CPU) handleInterrupt() {
	if cpu.interrupts.Enabled() {
		cpu.interrupts.Disable()
		switch {
		case cpu.interrupts.VblankPending():
			// 0040 Vertical Blank Interrupt Start Address
			cpu.rst(0x0040)()
			cpu.interrupts.ResetVblank()
		case cpu.interrupts.StatPending():
			// 0048 LCDC Status Interrupt Start Address
			cpu.rst(0x0048)()
			cpu.interrupts.ResetStat()
		case cpu.interrupts.TimerPending():
			// 0050 Timer OverflowInterrupt Start Address
			cpu.rst(0x0050)()
			cpu.interrupts.ResetTimer()
		case cpu.interrupts.SerialPending():
			// 0058 Serial Transfer Completion Interrupt Start Address
			cpu.rst(0x0058)()
			cpu.interrupts.ResetSerial()
		case cpu.interrupts.JoypadPending():
			// 0060 High-to-Low of P10-P13 Interrupt Start Address
			cpu.rst(0x0060)()
			cpu.interrupts.ResetJoypad()
		}

		// Now push the PC
		cpu.push(cpu.mapper, &cpu.m8b)()
		cpu.push(cpu.mapper, &cpu.m8a)()
	}
}

func (cpu *CPU) checkInterrupts() int {
	var length int
	if cpu.interrupts.Pending() {
		if cpu.halted {
			cpu.halted = false
			length++
		}
		if cpu.interrupts.Enabled() {
			length += 5
		}
	}
	return length
}

func (cpu *CPU) next() {

	if cpu.halted || cpu.stopped {
		return
	}

	// if cpu.interruptCycles == 0 {
	// 	cpu.interruptCycles = cpu.checkInterrupts()
	// }

	// if cpu.interruptCycles > 0 {
	// 	return
	// }

	mapper := cpu.mapper
	cpu.currentInstruction = mapper.Read(cpu.pc)

	// FIXME
	// every possible instruction value should be dispatchable
	// invalid instructions should panic

	if cpu.currentInstruction == 0xcb {
		cpu.pc++
		cpu.currentInstruction = mapper.Read(cpu.pc)
		cpu.currentCycle = 0
		cpu.currentEnded = cpu.ended
		cpu.currentMetadata = prefixedInstructionMetadata[cpu.currentInstruction]
		cpu.currentSubinstructions = prefix[cpu.currentInstruction]
	} else {
		cpu.currentCycle = 0
		cpu.currentEnded = cpu.ended
		cpu.currentMetadata = instructionMetadata[cpu.currentInstruction]
		cpu.currentSubinstructions = normal[cpu.currentInstruction]
	}

	// Reset any context from previous instructions
	cpu.u8a = 0
	cpu.u8b = 0
	cpu.m8a = 0
	cpu.m8b = 0

	// Finally increment PC
	if !cpu.haltbug {
		cpu.pc++
	} else {
		cpu.haltbug = false
	}

	if cpu.debugCPU {
		metadata := cpu.currentMetadata
		var pc uint16
		if cpu.currentMetadata.Prefixed {
			pc = cpu.pc - 2
		} else {
			pc = cpu.pc - 1
		}
		var operandValue string
		switch metadata.Length {
		case 2:
			u8 := mapper.Read(cpu.pc)
			operandValue = fmt.Sprintf("%02x", u8)
		case 3:
			u16 := uint16(mapper.Read(cpu.pc)) | uint16(mapper.Read(cpu.pc+1))<<8
			operandValue = fmt.Sprintf("%04x", u16)
		}
		fmt.Printf("0x%04x: [%02x] %-12s | %-4s | a:%02x b:%02x c:%02x d:%02x e:%02x f:%02x h:%02x l:%02x sp:%04x\n",
			pc, metadata.Addr, fmt.Sprintf("%s %s %s", metadata.Mnemonic, metadata.Operand1, metadata.Operand2), operandValue, cpu.a, cpu.b, cpu.c, cpu.d, cpu.e, cpu.f, cpu.h, cpu.l, cpu.sp)
	}

}

// ExecuteMachineCycle runs the CPU for one machine cycle
func (cpu *CPU) ExecuteMachineCycle() {
	if cpu.currentEnded() {
		cpu.next()
	}
	cpu.currentSubinstructions[cpu.currentCycle]()
	cpu.oam.Corrupt()
	cpu.currentCycle++
}
