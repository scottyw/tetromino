package cpu

import (
	"fmt"

	"github.com/scottyw/tetromino/pkg/gb/mem"
)

func (cpu *CPU) nopStep(*mem.Memory) {
	// Do nothing
}

func (cpu *CPU) interruptStep(memory *mem.Memory) {
	if cpu.ime {
		interrupts := cpu.hwr.IE & cpu.hwr.IF & 0x1f
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
