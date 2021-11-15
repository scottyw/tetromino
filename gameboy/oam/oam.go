package oam

// OAM captures the current state of sprite memory
type OAM struct {
	oam           [0xa0]byte
	dmaRunning    bool
	dmaCycle      uint16
	dmaBaseAddr   uint16
	dmaRead       uint8
	corrupt       bool
	ppuLastAccess uint16
	read          bool
	write         bool
	doubleWrite   bool
}

// New OAM
func New() *OAM {
	return &OAM{}
}

func (m *OAM) TriggerWriteCorruption(u16 uint16) {
	if !m.corrupt || u16 < 0xfe00 || u16 > 0xfeff {
		return
	}
	if m.write {
		m.doubleWrite = true
	} else {
		m.write = true
	}
}

func (m *OAM) Corrupt() {

	if !m.read && !m.write {
		return
	}

	if m.read && m.write {
		m.readWriteCorruption()
	}

	if m.read {
		m.readCorruption()
	} else {

		if m.doubleWrite {
			m.doubleWriteCorruption()
		}

		if m.write {
			m.writeCorruption()
		}

	}

	m.read = false
	m.write = false
	m.doubleWrite = false

}

func (m *OAM) writeCorruption() {

	// A write corruption corrupts the currently access row in the following manner,
	// as long as it's not the first row (containing the first two sprites):
	// - The first word in the row is replaced with this bitwise expression:
	// 	 ((a ^ c) & (b ^ c)) ^ c, where:
	//   a is the original value of that word,
	//   b is the first word in the preceding row, and
	//   c is the third word in the preceding row.
	// - The last three words are copied from the last three words in the preceding row.

	rowStart := ((m.ppuLastAccess - 0xfe00) / 8) * 8
	if rowStart == 0 {
		return
	}
	a := uint16(m.oam[rowStart])<<8 + uint16(m.oam[rowStart+1])
	b := uint16(m.oam[rowStart-8])<<8 + uint16(m.oam[rowStart-7])
	c := uint16(m.oam[rowStart-4])<<8 + uint16(m.oam[rowStart-3])
	a = ((a ^ c) & (b ^ c)) ^ c
	m.oam[rowStart] = uint8(a >> 8)
	m.oam[rowStart+1] = uint8(a & 0xff)
	copy(m.oam[rowStart+2:rowStart+8], m.oam[rowStart-6:rowStart])

}

func (m *OAM) readCorruption() {

	// A read corruption works similarly to a write corruption, except the bitwise expression is b | (a & c).

	rowStart := ((m.ppuLastAccess - 0xfe00) / 8) * 8
	if rowStart == 0 {
		return
	}
	a := uint16(m.oam[rowStart])<<8 + uint16(m.oam[rowStart+1])
	b := uint16(m.oam[rowStart-8])<<8 + uint16(m.oam[rowStart-7])
	c := uint16(m.oam[rowStart-4])<<8 + uint16(m.oam[rowStart-3])
	a = b | (a & c)
	m.oam[rowStart] = uint8(a >> 8)
	m.oam[rowStart+1] = uint8(a & 0xff)
	copy(m.oam[rowStart+2:rowStart+8], m.oam[rowStart-6:rowStart])

}

func (m *OAM) readWriteCorruption() {

	// This corruption will not happen if the accessed row is one of the first four, as well as if it's the last row:
	// The first word in the row preceding the currently accessed row is replaced with the following bitwise expression:
	// (b & (a | c | d)) | (a & c & d)
	// where
	// a is the first word two rows before the currently accessed row,
	// b is the first word in the preceding row (the word being corrupted),
	// c is the first word in the currently accessed row, and
	// d is the third word in the preceding row.
	// The contents of the preceding row is copied (after the corruption of the first word in it) both
	// to the currently accessed row and
	// to two rows before the currently accessed row
	// Regardless of wether the previous corruption occurred or not, a normal read corruption is then applied.

	row := (m.ppuLastAccess - 0xfe00) / 8
	if row < 5 || row == 19 {
		return
	}
	rowStart := row * 8
	a := uint16(m.oam[rowStart-16])<<8 + uint16(m.oam[rowStart-15])
	b := uint16(m.oam[rowStart-8])<<8 + uint16(m.oam[rowStart-7])
	c := uint16(m.oam[rowStart])<<8 + uint16(m.oam[rowStart+1])
	d := uint16(m.oam[rowStart-4])<<8 + uint16(m.oam[rowStart-3])
	b = (b & (a | c | d)) | (a & c & d)
	m.oam[rowStart-8] = uint8(b >> 8)
	m.oam[rowStart-7] = uint8(b & 0xff)
	copy(m.oam[rowStart:rowStart+8], m.oam[rowStart-8:rowStart])
	copy(m.oam[rowStart-16:rowStart-8], m.oam[rowStart-8:rowStart])

}

func (m *OAM) doubleWriteCorruption() {

	// Just like a normal write corruption but affects the previous row

	rowStart := (((m.ppuLastAccess - 0xfe00) / 8) - 1) * 8
	if rowStart < 1 {
		return
	}
	a := uint16(m.oam[rowStart])<<8 + uint16(m.oam[rowStart+1])
	b := uint16(m.oam[rowStart-8])<<8 + uint16(m.oam[rowStart-7])
	c := uint16(m.oam[rowStart-4])<<8 + uint16(m.oam[rowStart-3])
	a = ((a ^ c) & (b ^ c)) ^ c
	m.oam[rowStart] = uint8(a >> 8)
	m.oam[rowStart+1] = uint8(a & 0xff)
	copy(m.oam[rowStart+2:rowStart+8], m.oam[rowStart-6:rowStart])

}

func (m *OAM) EnterMode2() {
	m.corrupt = true
}

func (m *OAM) ExitMode2() {
	m.corrupt = false
}

func (m *OAM) Read(addr uint16) uint8 {
	if m.dmaRunning {
		return 0xff
	}
	if m.corrupt {
		m.read = true
	}
	if addr >= 0xfea0 {
		return 0
	}

	if addr >= 0xfea0 {
		return 0
	}
	return m.oam[addr-0xfe00]
}

func (m *OAM) Write(addr uint16, value uint8) {
	if m.corrupt {
		if m.write {
			m.doubleWrite = true
		} else {
			m.write = true
		}
	}
	if addr < 0xfea0 {
		m.oam[addr-0xfe00] = value
	}
}

func (m *OAM) PPURead(addr uint16) uint8 {
	m.ppuLastAccess = addr
	if m.dmaRunning {
		return 0xff
	}
	return m.oam[addr-0xfe00]
}

func (m *OAM) TickDMA(read func(uint16) uint8) {
	if m.dmaRunning {
		if m.dmaCycle == 0 {
			// Setup
		} else if m.dmaCycle == 1 {
			m.dmaRead = read(m.dmaBaseAddr)
		} else if m.dmaCycle == 161 {
			m.oam[159] = m.dmaRead
			m.dmaRunning = false
		} else {
			m.oam[m.dmaCycle-2] = m.dmaRead
			m.dmaRead = read(m.dmaBaseAddr + m.dmaCycle - 1)
		}
		m.dmaCycle++
	}
}

func (m *OAM) startDMA(value uint8) {
	m.dmaRunning = true
	m.dmaCycle = 0
	m.dmaBaseAddr = uint16(value) << 8
	if m.dmaBaseAddr >= 0xe000 {
		m.dmaBaseAddr -= 0x2000
	}
}

// FF46 - DMA - DMA Transfer and Start Address (W)
// Writing to this register launches a DMA transfer from ROM or RAM to OAM memory
// (sprite attribute table). The written value specifies the transfer source
// address divided by 100h, ie. source & destination are:
// Source:      XX00-XX9F   ;XX in range from 00-F1h
// Destination: FE00-FE9F
// It takes 160 microseconds until the transfer has completed (80 microseconds in
// CGB Double Speed Mode), during this time the CPU can access only HRAM (memory
// at FF80-FFFE). For this reason, the programmer must copy a short procedure
// into HRAM, and use this procedure to start the transfer from inside HRAM, and
// wait until the transfer has finished:
//  ld  (0FF46h),a ;start DMA transfer, a=start address/100h
//  ld  a,28h      ;delay...
// wait:           ;total 5x40 cycles, approx 200ms
//  dec a          ;1 cycle
//  jr  nz,wait    ;4 cycles
// Most programs are executing this procedure from inside of their VBlank
// procedure, but it is possible to execute it during display redraw also,
// allowing to display more than 40 sprites on the screen (ie. for example 40
// sprites in upper half, and other 40 sprites in lower half of the screen).

func (m *OAM) WriteDMA(value uint8) {
	// fmt.Printf("> DMA - 0x%02x\n", value)
	m.startDMA(value)
}

func (m *OAM) ReadDMA() uint8 {
	dma := uint8(m.dmaBaseAddr >> 8)
	// fmt.Printf("< DMA - 0x%02x\n", dma)
	return dma
}
