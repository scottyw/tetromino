package mem

import (
	"fmt"
	"io/ioutil"
)

// Memory allows read and write access to memory
type Memory struct {
	hwr               *HardwareRegisters
	mbc               mbc
	VideoRAM          [0x2000]byte
	internalRAM       [0x2000]byte
	OAM               [0xa0]byte
	zeroPage          [0x8f]byte
	WriteNotification WriteNotification
}

// WriteNotification provides a mechanism to notify other subsystems about memory writes
type WriteNotification interface {
	WriteToVideoRAM(addr uint16)
}

// NewMemoryFromFile loads the specified ROM file and calls NewMemory
func NewMemoryFromFile(hwr *HardwareRegisters, romFilename string) *Memory {
	var rom []byte
	if romFilename == "" {
		panic(fmt.Sprintf("No ROM file specified"))
	}
	rom, err := ioutil.ReadFile(romFilename)
	if err != nil {
		panic(fmt.Sprintf("Failed to read the ROM file at \"%s\" (%v)", romFilename, err))
	}
	return NewMemory(hwr, rom)
}

// NewMemory creates the memory struct and initializes it with ROM contents and default values
func NewMemory(hwr *HardwareRegisters, rom []byte) *Memory {
	return &Memory{
		hwr: hwr,
		mbc: newMBC(rom),
	}
}

// Read a byte from the chosen memory location
func (mem *Memory) Read(addr uint16) byte {
	switch {
	case addr < 0x8000:
		return mem.mbc.read(addr)
	case addr < 0xa000:
		return mem.VideoRAM[addr-0x8000]
	case addr < 0xc000:
		// FIXME
		// panic(fmt.Sprintf("Read on cartridge RAM is not implemented: 0x%04x", addr))
		return 0
	case addr < 0xe000:
		return mem.internalRAM[addr-0xc000]
	case addr < 0xfe00:
		return mem.internalRAM[addr-0xe000]
	case addr < 0xfea0:
		return mem.OAM[addr-0xfe00]
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
		if mem.WriteNotification != nil {
			mem.WriteNotification.WriteToVideoRAM(addr)
		}
		mem.VideoRAM[addr-0x8000] = value
	case addr < 0xc000:
		// FIXME maybe write to file
		// panic(fmt.Sprintf("Write on cartridge RAM is not implemented: 0x%04x", addr))
	case addr < 0xe000:
		mem.internalRAM[addr-0xc000] = value
	case addr < 0xfe00:
		mem.internalRAM[addr-0xe000] = value
	case addr < 0xfea0:
		mem.OAM[addr-0xfe00] = value
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
