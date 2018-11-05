package mem

import (
	"fmt"
)

type mbc struct {
	// ROM and RAM data and mask read from the cart
	rom [][0x4000]byte
	ram [][0x2000]byte

	// Record of what as written between 0x0000 and 0x8000
	enabledRegion uint8
	romRegion     uint8
	ramRegion     uint8
	modeRegion    uint8

	// Selected ROM and RAM banks computed from written values on an MBC-type basis
	ramEnabled bool
	romBank0   int
	romBankX   int
	ramBank    int
	update     func(*mbc)
}

func newMBC(rom []byte) *mbc {
	// FIXME this is really for unit tests - need to be more graceful about this and missing ROM files in general
	if rom == nil || len(rom) < 0x0148 {
		return nil
	}
	cartType := rom[0x0147]
	romSize := rom[0x0148]
	ramSize := rom[0x0149]
	return &mbc{
		rom:      splitROMIntoPages(romSize, rom),
		ram:      createRAM(cartType, ramSize),
		romBank0: 0,
		romBankX: 1,
		update:   chooseUpdateFunc(cartType),
	}
}

func splitROMIntoPages(romSize uint8, rom []byte) [][0x4000]byte {
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

func createRAM(cartType, ramSize uint8) [][0x2000]byte {
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
			return nil
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

func chooseUpdateFunc(cartType uint8) func(*mbc) {
	switch cartType {
	case 0x00:
		// 00 - ROM ONLY
		return func(_ *mbc) {}
	case 0x01:
		// 01 - ROM + MBC1
		return updateMBC1
	case 0x02:
		// 02 - ROM + MBC1 + RAM
		return updateMBC1
	case 0x03:
		// 03 - ROM + MBC1 + RAM + BATT
		return updateMBC1
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
		return updateMBC3
	case 0x10:
		// 10 - ROM + MBC3 + RAM + TIMER + BATT
		return updateMBC3
	case 0x11:
		// 11 - ROM + MBC3
		return updateMBC3
	case 0x12:
		// 12 - ROM + MBC3 + RAM
		return updateMBC3
	case 0x13:
		// 13 - ROM + MBC3 + RAM + BATT
		return updateMBC3
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

func (m *mbc) read(addr uint16) uint8 {
	switch {
	case addr < 0x4000:
		return m.rom[m.romBank0][addr]
	case addr < 0x8000:
		offset := addr - 0x4000
		return m.rom[m.romBankX][offset]
	case addr < 0xa000:
		panic(fmt.Sprintf("mbc has no read mapping for address 0x%04x", addr))
	case addr < 0xc000:
		if m.ramEnabled {
			offset := addr - 0xa000
			return m.ram[m.ramBank][offset]
		}
		return 0xff
	default:
		panic(fmt.Sprintf("mbc has no read mapping for address 0x%04x", addr))
	}
}

func (m *mbc) write(addr uint16, value uint8) {
	switch {
	case addr < 0x2000:
		m.enabledRegion = value
	case addr < 0x4000:
		m.romRegion = value
	case addr < 0x6000:
		m.ramRegion = value
	case addr < 0x8000:
		m.modeRegion = value
	case addr < 0xa000:
		panic(fmt.Sprintf("mbc has no write mapping for address 0x%04x", addr))
	case addr < 0xc000:
		offset := addr - 0xa000
		if m.ramEnabled {
			m.ram[m.ramBank][offset] = value
		}
	default:
		panic(fmt.Sprintf("mbc has no write mapping for address 0x%04x", addr))
	}
	m.update(m)
}

func updateMBC1(m *mbc) {

	// Check if RAM is enabled
	m.ramEnabled = m.enabledRegion&0x0f == 0x0a

	// Check ROM bank 0
	if m.modeRegion&0x01 == 0 {
		m.romBank0 = 0
	} else {
		m.romBank0 = int((m.ramRegion & 0x03) << 5)
		m.romBank0 = m.romBank0 % len(m.rom)
	}

	// Check ROM bank 1
	m.romBankX = int((m.romRegion & 0x1f))
	if m.romBankX == 0 {
		m.romBankX = 1
	}
	m.romBankX |= int((m.ramRegion & 0x03 << 5))
	m.romBankX = m.romBankX % len(m.rom)

	// Check RAM bank
	if m.ramEnabled {
		if m.modeRegion&0x01 == 0 {
			m.ramBank = 0
		} else {
			m.ramBank = int(m.ramRegion & 0x03)
			m.ramBank = m.ramBank % len(m.ram)
		}
	}

}

func updateMBC3(m *mbc) {

	// Check if RAM is enabled
	m.ramEnabled = m.enabledRegion&0x0f == 0x0a

	// Check ROM bank 1
	m.romBankX = int((m.romRegion & 0x7f))
	if m.romBankX == 0 {
		m.romBankX = 1
	}
	m.romBankX = m.romBankX % len(m.rom)

	// Check RAM bank
	if m.ramEnabled {
		m.ramBank = int(m.ramRegion & 0x07)
		m.ramBank = m.ramBank % len(m.ram)
	}

}
