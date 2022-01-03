package memory

import (
	"fmt"
)

type mbc interface {
	Read(addr uint16) uint8
	Write(addr uint16, value uint8)
	DumpRAM() []byte
}

type none struct {
	rom []byte
}

func (n *none) Read(addr uint16) uint8 {
	return n.rom[addr]
}

func (n *none) Write(addr uint16, value uint8) {
	// Do nothing
}

func (n *none) DumpRAM() []byte {
	return []byte{}
}

func newMBC(romImage []byte) mbc {

	if romImage == nil || len(romImage) < 0x0148 {
		return nil
	}
	romSize := romImage[0x0148]
	rom := prepareROM(romSize, romImage)
	cartType := romImage[0x0147]
	ramSize := romImage[0x0149]
	ram := prepareRAM(cartType, ramSize)

	switch cartType {
	case 0x00:
		// 00 - ROM ONLY
		return &none{rom: romImage}
	case 0x01:
		// 01 - ROM + MBC1
		return newMBC1(rom, ram)
	case 0x02:
		// 02 - ROM + MBC1 + RAM
		return newMBC1(rom, ram)
	case 0x03:
		// 03 - ROM + MBC1 + RAM + BATT
		return newMBC1(rom, ram)
	case 0x05:
		// 05 - ROM + MBC2
	case 0x06:
		// 06 - ROM + MBC2 + BATT
	case 0x08:
		// 08 - ROM + RAM
	case 0x09:
		// 09 - ROM + RAM + BATT + MMM01
	case 0x0c:
		// 0C - ROM + MMM01 + SRAM
	case 0x0d:
		// 0D - ROM + MMM01 + SRAM + BATT
	case 0x0f:
		// 0f - ROM + MBC3 + TIMER + BATT
		// return newMBC3(rom, ram)
	case 0x10:
		// 10 - ROM + MBC3 + RAM + TIMER + BATT
		// return newMBC3(rom, ram)
	case 0x11:
		// 11 - ROM + MBC3
		// return newMBC3(rom, ram)
	case 0x12:
		// 12 - ROM + MBC3 + RAM
		// return newMBC3(rom, ram)
	case 0x13:
		// 13 - ROM + MBC3 + RAM + BATT
		// return newMBC3(rom, ram)
	case 0x19:
		// 19 - ROM + MBC5
	case 0x1a:
		// 1A - ROM + MBC5 + RAM
	case 0x1b:
		// 1B - ROM + MBC5 + RAM + BATT
	case 0x1c:
		// 1C - ROM + MBC5 + RUMBLE
	case 0x1d:
		// 1D - ROM + MBC5 + RUMBLE + SRAM
	case 0x1e:
		// 1E - ROM + MBC5 + RUMBLE + SRAM + BATT
	case 0x20:
		// 20 - ROM + MBC6 + RAM + BATT
	case 0x22:
		// 22 - ROM + MBC7 + RAM + BATT + ACCELEROMETER
	case 0xfc:
		// FC - POCKET CAMERA
	case 0xfd:
		// FD - Bandai TAMA5
	case 0xfe:
		// FE - Hudson on HuC-1
	case 0xff:
		// FF - Hudson on HuC-1 + RAM + BATTERY
	}
	panic(fmt.Sprintf("mbc does not support cart type 0x%02x", cartType))
}

func prepareROM(romSize uint8, rom []byte) [][0x4000]byte {
	if len(rom)%0x4000 != 0 {
		panic(fmt.Sprintf("ROM size must be a multiple of 32KB. Current size: 0x%02x", len(rom)))
	}
	pageCount := len(rom) / 0x4000
	if pageCount != (0x02 << romSize) {
		panic(fmt.Sprintf("Actual ROM size must match reported size: Actual=0x%04x Reported=0x%04x", pageCount, (0x02 << romSize)))
	}
	pages := make([][0x4000]byte, pageCount)
	for i := 0; i < pageCount; i++ {
		copy(pages[i][:], rom[i*0x4000:])
	}
	return pages
}

func prepareRAM(cartType, ramSize uint8) [][0x2000]byte {
	var ram [][0x2000]byte
	if cartType == 0x05 || cartType == 0x06 {
		ram = make([][0x2000]byte, 1)
	} else {
		switch ramSize {
		case 0x01:
			ram = make([][0x2000]byte, 1)
		case 0x02:
			ram = make([][0x2000]byte, 1)
		case 0x03:
			ram = make([][0x2000]byte, 4)
		case 0x04:
			ram = make([][0x2000]byte, 16)
		case 0x05:
			ram = make([][0x2000]byte, 8)
		default:
			// This cart has no RAM but capture writes anyway to allow test ROM validation
			ram = make([][0x2000]byte, 1)
		}
	}
	// Initialize it to 0xff
	for i := 0; i < len(ram); i++ {
		for j := 0; j < 0x2000; j++ {
			ram[i][j] = 0xff
		}
	}
	return ram
}
