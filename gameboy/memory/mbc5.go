package memory

type mbc5 struct {
	// ROM and RAM data and mask read from the cart
	rom [][0x4000]byte
	ram [][0x2000]byte

	// Internal state
	ramEnabled bool
	romBank    uint16
	ramBank    uint8
}

func newMBC5(rom [][0x4000]byte, ram [][0x2000]byte) mbc {
	mbc := &mbc5{
		rom:     rom,
		ram:     ram,
		romBank: 1,
	}
	return mbc
}

func (m *mbc5) Read(addr uint16) uint8 {
	switch {
	case addr < 0x4000:
		return m.rom[0][addr]
	case addr < 0x8000:
		offset := addr - 0x4000
		return m.rom[m.romBank][offset]
	case addr < 0xa000:
		return 0xff
	case addr < 0xc000:
		if m.ramEnabled {
			offset := addr - 0xa000
			return m.ram[m.ramBank][offset]
		}
		return 0xff
	default:
		return 0xff
	}
}

func (m *mbc5) Write(addr uint16, value uint8) {
	switch {
	case addr < 0x2000:
		m.ramEnabled = value&0x0f == 0x0a
	case addr < 0x3000:
		m.romBank = m.romBank&0xff00 + uint16(value)
		m.romBank %= uint16(len(m.rom))
	case addr < 0x4000:
		m.romBank = uint16(value)<<8 + m.romBank&0x00ff
		m.romBank %= uint16(len(m.rom))
	case addr < 0x6000:
		m.ramBank = value & 0x0f
		m.ramBank %= uint8(len(m.ram))
	default:
		// Ignore
	}
}

func (m *mbc5) DumpRAM() []byte {
	var dump []byte
	for _, r := range m.ram {
		dump = append(dump, r[:]...)
	}
	return dump
}
