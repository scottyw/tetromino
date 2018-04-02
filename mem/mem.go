package mem

import (
	"fmt"
)

// Memory allows read and write access to memory
type Memory struct {
	hwr         *HardwareRegisters
	mbc         mbc
	videoRAM    [0x2000]byte
	internalRAM [0x2000]byte
	oam         [0xa0]byte
	zeroPage    [0x8f]byte
}

// NewMemory creates the memory and initializes it with ROM contents and default values
func NewMemory(hwr *HardwareRegisters) *Memory {
	return &Memory{
		hwr: hwr,
		mbc: newMBC(),
	}
}

// Read a byte from the chosen memory location
func (mem *Memory) Read(addr uint16) byte {
	switch {
	case addr < 0x8000:
		return mem.mbc.read(addr)
	case addr < 0xa000:
		return mem.videoRAM[addr-0x8000]
	case addr < 0xc000:
		panic(fmt.Sprintf("Read on cartridge RAM is not implemented: 0x%04x", addr))
	case addr < 0xe000:
		return mem.internalRAM[addr-0xc000]
	case addr < 0xfe00:
		return mem.internalRAM[addr-0xe000]
	case addr < 0xfea0:
		return mem.oam[addr-0xfe00]
	case addr < 0xff00:
		return 0 // Unusable region
	case addr < 0xff80:
		return mem.readHardwareRegisters(addr)
	case addr < 0xffff:
		return mem.zeroPage[addr-0xff80]
	case addr == 0xffff:
		return mem.hwr.IE
	default:
		panic(fmt.Sprintf("Read failed: 0x%04x", addr))
	}
}

// Write a byte to the chosen memory location
func (mem *Memory) Write(addr uint16, value byte) {
	switch {
	case addr < 0x8000:
		mem.mbc.write(addr, value)
	case addr < 0xa000:
		mem.videoRAM[addr-0x8000] = value
	case addr < 0xc000:
		panic(fmt.Sprintf("Read on cartridge RAM is not implemented: 0x%04x", addr))
	case addr < 0xe000:
		mem.internalRAM[addr-0xc000] = value
	case addr < 0xfe00:
		mem.internalRAM[addr-0xe000] = value
	case addr < 0xfea0:
		mem.oam[addr-0xfe00] = value
	case addr < 0xff00:
		// Unusable region
	case addr < 0xff80:
		mem.writeHardwareRegisters(addr, value)
	case addr < 0xffff:
		mem.zeroPage[addr-0xff80] = value
	case addr == IE:
		mem.hwr.IE = value
	default:
		panic(fmt.Sprintf("Write failed: 0x%04x", addr))
	}
}
