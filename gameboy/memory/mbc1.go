package memory

import (
	"fmt"
)

type mbc1 struct {
	// ROM and RAM data and mask read from the cart
	rom [][0x4000]byte
	ram [][0x2000]byte

	// Record of what as written between 0x0000 and 0x8000
	ramEnabled bool
	bank1      uint8
	bank2      uint8
	mode1      bool

	// Precomputed values for ROM and RAM banks determined at write time to optimize for read time
	romBank0 uint8
	romBank1 uint8
	ramBank  uint8
}

func newMBC1(rom [][0x4000]byte, ram [][0x2000]byte) mbc {
	mbc := &mbc1{
		rom:   rom,
		ram:   ram,
		bank1: 1,
		bank2: 0,
	}
	mbc.updateBanks()
	return mbc
}

func (m *mbc1) Read(addr uint16) uint8 {
	switch {
	case addr < 0x4000:
		return m.rom[m.romBank0][addr]
	case addr < 0x8000:
		offset := addr - 0x4000
		return m.rom[m.romBank1][offset]
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

func (m *mbc1) Write(addr uint16, value uint8) {
	switch {
	case addr < 0x2000:
		m.ramEnabled = value&0x0f == 0x0a
		m.updateBanks()
	case addr < 0x4000:
		m.bank1 = value & 0x1f
		if m.bank1 == 0 {
			m.bank1++
		}
		m.updateBanks()
	case addr < 0x6000:
		m.bank2 = value & 0x03
		m.updateBanks()
	case addr < 0x8000:
		m.mode1 = value&0x01 != 0
		m.updateBanks()
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
}

func (m *mbc1) updateBanks() {

	// Update ROM bank 0
	if m.mode1 {
		m.romBank0 = (m.bank2 << 5) % uint8(len(m.rom))
	} else {
		m.romBank0 = 0
	}

	// Update ROM bank 1
	m.romBank1 = (m.bank1 | m.bank2<<5) % uint8(len(m.rom))

	// Update RAM bank
	if m.ramEnabled {
		if m.mode1 {
			m.ramBank = m.bank2 % uint8(len(m.ram))
		} else {
			m.ramBank = 0
		}
	}

}

func (m *mbc1) DumpRAM() []byte {
	var dump []byte
	for _, r := range m.ram {
		dump = append(dump, r[:]...)
	}
	return dump
}
