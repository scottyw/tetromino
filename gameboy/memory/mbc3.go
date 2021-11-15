package memory

import (
	"fmt"
)

type mbc3 struct {
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

func newMBC3(rom [][0x4000]byte, ram [][0x2000]byte) mbc {
	return &mbc3{
		rom:      rom,
		ram:      ram,
		romBank0: 0,
		romBankX: 1,
	}
}

func (m *mbc3) Read(addr uint16) uint8 {
	switch {
	case addr < 0x4000:
		return m.rom[m.romBank0][addr]
	case addr < 0x8000:
		offset := addr - 0x4000
		return m.rom[m.romBankX][offset]
	case addr < 0xa000:
		panic(fmt.Sprintf("mbc1 has no read mapping for address 0x%04x", addr))
	case addr < 0xc000:
		if m.ramEnabled {
			offset := addr - 0xa000
			return m.ram[m.ramBank][offset]
		}
		return 0xff
	default:
		panic(fmt.Sprintf("mbc1 has no read mapping for address 0x%04x", addr))
	}
}

func (m *mbc3) Write(addr uint16, value uint8) {

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
		panic(fmt.Sprintf("mbc1 has no write mapping for address 0x%04x", addr))
	case addr < 0xc000:
		offset := addr - 0xa000
		if m.ramEnabled {
			m.ram[m.ramBank][offset] = value
		}
	default:
		panic(fmt.Sprintf("mbc1 has no write mapping for address 0x%04x", addr))
	}

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

func (m *mbc3) DumpRAM() []byte {
	var dump []byte
	for _, r := range m.ram {
		dump = append(dump, r[:]...)
	}
	return dump
}
