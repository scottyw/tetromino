package cpu

import (
	"testing"

	"github.com/scottyw/tetromino/pkg/gb/mem"
)

func TestSanity(t *testing.T) {
	// Check all unprefixed instructions are fundamentally sound
	for instruction := 0; instruction < 256; instruction++ {
		if instruction == 0xcb || instructionMetadata[instruction].Mnemonic == "" {
			continue
		}
		hwr := mem.NewHardwareRegisters(nil)
		cpu := NewCPU(hwr, false, false, false)
		cpu.validateFlags = true
		rom := make([]byte, 0x0200)
		memory := mem.NewMemory(hwr, rom)
		rom[0x100] = uint8(instruction)
		cpu.execute(memory)
	}
	// Check all prefixed instructions are fundamentally sound
	for instruction := 0; instruction < 256; instruction++ {
		hwr := mem.NewHardwareRegisters(nil)
		cpu := NewCPU(hwr, false, false, false)
		cpu.validateFlags = true
		rom := make([]byte, 0x0200)
		memory := mem.NewMemory(hwr, rom)
		rom[0x100] = 0xcb
		rom[0x101] = uint8(instruction)
		cpu.execute(memory)
	}
}
