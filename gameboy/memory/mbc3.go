package memory

import (
	"fmt"
)

type mbc3 struct {
	// ROM and RAM data and mask read from the cart
	rom [][0x4000]byte
	ram [][0x2000]byte

	// Internal state
	ramEnabled bool
	romBank    uint8
	ramBank    uint8
	rtc        *rtc
}

func newMBC3(rom [][0x4000]byte, ram [][0x2000]byte, rtc *rtc) mbc {
	mbc := &mbc3{
		rom:     rom,
		ram:     ram,
		romBank: 1,
		rtc:     rtc,
	}
	return mbc
}

func (m *mbc3) Read(addr uint16) uint8 {
	switch {
	case addr < 0x4000:
		return m.rom[0][addr]
	case addr < 0x8000:
		offset := addr - 0x4000
		return m.rom[m.romBank][offset]
	case addr < 0xa000:
		panic(fmt.Sprintf("mbc3 has no read mapping for address 0x%04x", addr))
	case addr < 0xc000:
		if m.ramEnabled {
			if m.ramBank >= 0x08 {
				return m.rtc.read(m.ramBank)
			}
			offset := addr - 0xa000
			return m.ram[m.ramBank][offset]
		}
		return 0xff
	default:
		panic(fmt.Sprintf("mbc3 has no read mapping for address 0x%04x", addr))
	}
}

func (m *mbc3) Write(addr uint16, value uint8) {
	switch {
	case addr < 0x2000:
		m.ramEnabled = value&0x0f == 0x0a
	case addr < 0x4000:
		m.romBank = value & 0x7f
	case addr < 0x6000:
		m.ramBank = value & 0x0f
	case addr < 0x8000:
		if value&0x01 == 0 {
			m.rtc.latchLow()
		} else {
			m.rtc.latchHigh()
		}
	case addr < 0xa000:
		panic(fmt.Sprintf("mbc3 has no write mapping for address 0x%04x", addr))
	case addr < 0xc000:
		offset := addr - 0xa000
		if m.ramEnabled {
			if m.ramBank >= 0x08 {
				m.rtc.write(m.ramBank, value)
				return
			}
			m.ram[m.ramBank][offset] = value
		}
	default:
		panic(fmt.Sprintf("mbc3 has no write mapping for address 0x%04x", addr))
	}
}

func (m *mbc3) DumpRAM() []byte {
	var dump []byte
	for _, r := range m.ram {
		dump = append(dump, r[:]...)
	}
	return dump
}
