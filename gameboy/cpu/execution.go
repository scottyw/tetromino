package cpu

import "fmt"

func (cpu *CPU) handleInterrupt() func() {
	return func() {
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

func (cpu *CPU) useAltMachineCycles(instruction uint8) bool {

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

func (cpu *CPU) next() {

	if cpu.halted || cpu.stopped {
		return
	}

	if cpu.interruptCycles == 0 {
		cpu.interruptCycles = cpu.checkInterrupts()
	}

	if cpu.interruptCycles > 0 {
		fmt.Println("interrupt")
		return
	}

	mapper := cpu.mapper
	cpu.cycle = 0
	cpu.instruction = mapper.Read(cpu.pc)

	// FIXME
	// every possible instruction value should be dispatchable
	// invalid instructions should panic

	cpu.prefixed = cpu.instruction == 0xcb
	if cpu.prefixed {
		cpu.pc++
		cpu.instruction = mapper.Read(cpu.pc)
	}

	// Reset any context from previous instructions
	cpu.u8a = 0
	cpu.u8b = 0
	cpu.m8a = 0
	cpu.m8b = 0

	// Check for instructions that need to use the shorter alt machine cycle count
	// FIXME
	// if cpu.useAltMachineCycles(md.Dispatch) {
	// 	steps = steps[:md.AltMachineCycles]
	// }

	// Finally increment PC
	if !cpu.haltbug {
		cpu.pc++
	} else {
		cpu.haltbug = false
	}

	if cpu.debugCPU {
		var metadata *metadata
		var pc uint16
		if cpu.prefixed {
			metadata = prefixedInstructionMetadata[cpu.instruction]
			pc = cpu.pc - 2
		} else {
			metadata = instructionMetadata[cpu.instruction]
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
			pc, cpu.instruction, fmt.Sprintf("%s %s %s", metadata.Mnemonic, metadata.Operand1, metadata.Operand2), operandValue, cpu.a, cpu.b, cpu.c, cpu.d, cpu.e, cpu.f, cpu.h, cpu.l, cpu.sp)
	}

}

// ExecuteMachineCycle runs the CPU for one machine cycle
func (cpu *CPU) ExecuteMachineCycle() {
	if cpu.prefixed {
		if cpu.cycle == len(prefix[cpu.instruction]) {
			cpu.next()
		}
	} else {
		if cpu.cycle == len(normal[cpu.instruction]) {
			cpu.next()
		}
	}
	if cpu.interruptCycles == 1 {
		cpu.handleInterrupt()
	}
	if cpu.interruptCycles > 0 {
		cpu.interruptCycles--
	} else {
		if cpu.prefixed {
			prefix[cpu.instruction][cpu.cycle]()
		} else {
			normal[cpu.instruction][cpu.cycle]()
		}
	}
	cpu.oam.Corrupt()
	cpu.cycle++
}
