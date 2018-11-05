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
	m.updateMBC1()
}

func (m *mbc) updateMBC1() {

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
