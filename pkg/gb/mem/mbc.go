package mem

import (
	"fmt"
)

type mbc interface {
	read(uint16) uint8
	write(uint16, uint8)
}

type mbc1 struct {
	rom        []byte
	ram        [4][0x2000]byte
	romBank    uint8
	ramBank    uint8
	mode       uint8
	ramEnabled bool
}

func (m *mbc1) read(addr uint16) uint8 {
	switch {
	case addr < 0x4000:
		return m.rom[addr]
	case addr < 0x8000:
		var bank uint8
		if m.mode == 0 {
			bank = m.ramBank<<5 | m.romBank
		} else {
			bank = m.romBank
		}
		offset := (uint(bank)-1)*0x4000 + uint(addr)
		return m.rom[offset]
	case addr < 0xa000:
		panic(fmt.Sprintf("MBC1 has no mapping for address 0x%04x", addr))
	case addr < 0xc000:
		// FIXME support reading from RAM
		return 0
	default:
		panic(fmt.Sprintf("MBC1 has no mapping for address 0x%04x", addr))
	}
}

func (m *mbc1) write(addr uint16, value uint8) {
	switch {
	case addr < 0x2000:
		m.ramEnabled = value&0x0f == 0x0a
	case addr < 0x4000:
		m.romBank = value & 0x1f
		if m.romBank == 0 {
			m.romBank = 1
		}
	case addr < 0x6000:
		m.ramBank = value & 0x03
	case addr < 0x8000:
		m.mode = value & 0x01
		if m.mode == 0 {
			m.ramBank = 0
		}
	case addr < 0xa000:
		panic(fmt.Sprintf("MBC1 has no mapping for address 0x%04x", addr))
	case addr < 0xc000:
		// FIXME support writing to RAM
	default:
		panic(fmt.Sprintf("MBC1 has no mapping for address 0x%04x", addr))
	}
}

func newMBC(rom []byte) mbc {
	return &mbc1{
		rom:     rom,
		ram:     [4][0x2000]byte{[0x2000]byte{}, [0x2000]byte{}, [0x2000]byte{}, [0x2000]byte{}},
		romBank: 1,
	}
}
