package oam

// OAM captures the current state of sprite memory
type OAM struct {
	oam         [0xa0]byte
	dmaRunning  bool
	dmaCycle    uint16
	dmaBaseAddr uint16
	dmaRead     uint8
}

// New OAM
func New() *OAM {
	return &OAM{}
}

func (m *OAM) ReadOAM(addr uint16) uint8 {
	if m.dmaRunning {
		return 0xff
	}
	return m.oam[addr-0xfe00]
}

func (m *OAM) WriteOAM(addr uint16, value uint8) {
	m.oam[addr-0xfe00] = value
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
