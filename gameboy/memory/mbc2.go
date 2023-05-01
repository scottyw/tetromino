package memory

type mbc2 struct {
	// ROM and RAM data and mask read from the cart
	rom [][0x4000]byte
	ram []byte

	// Internal state
	ramEnabled bool
	romBank    uint8
}

func newMBC2(rom [][0x4000]byte) mbc {
	mbc := &mbc2{
		rom:     rom,
		ram:     make([]byte, 512),
		romBank: 1,
	}
	return mbc
}

func (m *mbc2) Read(addr uint16) uint8 {
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
			offset %= 0x0200
			return m.ram[offset]
		}
		return 0xff
	default:
		return 0xff
	}
}

func (m *mbc2) Write(addr uint16, value uint8) {
	switch {
	case addr < 0x4000:
		if addr&0x0100 == 0 {
			m.ramEnabled = value&0x0f == 0x0a
		} else {
			m.romBank = value & 0x0f
			if m.romBank == 0 {
				m.romBank = 0x01
			}
			m.romBank %= uint8(len(m.rom))
		}
	case addr < 0xa000:
		// Ignore
	case addr < 0xc000:
		if m.ramEnabled {
			offset := addr - 0xa000
			offset %= 0x0200
			m.ram[offset] = value | 0xf0
		}
	default:
		// Ignore
	}
}

func (m *mbc2) DumpRAM() []byte {
	return m.ram
}
