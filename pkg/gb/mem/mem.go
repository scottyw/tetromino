package mem

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/scottyw/tetromino/pkg/gb/timer"
)

const (

	//
	// Register constants
	//

	LCDC = 0xFF40
	STAT = 0xFF41
	SCY  = 0xFF42
	SCX  = 0xFF43
	LY   = 0xFF44
	LYC  = 0xFF45
	WY   = 0xFF4A
	WX   = 0xFF4B
	DMA  = 0xFF46
	BGP  = 0xFF47
	OBP0 = 0xFF48
	OBP1 = 0xFF49
	NR10 = 0xFF10
	NR11 = 0xFF11
	NR12 = 0xFF12
	NR13 = 0xFF13
	NR14 = 0xFF14
	NR21 = 0xFF16
	NR22 = 0xFF17
	NR23 = 0xFF18
	NR24 = 0xFF19
	NR30 = 0xFF1A
	NR31 = 0xFF1B
	NR32 = 0xFF1C
	NR33 = 0xFF1D
	NR34 = 0xFF1E
	NR41 = 0xFF20
	NR42 = 0xFF21
	NR43 = 0xFF22
	NR44 = 0xFF23
	NR50 = 0xFF24
	NR51 = 0xFF25
	NR52 = 0xFF26
	JOYP = 0xFF00
	SB   = 0xFF01
	SC   = 0xFF02
	DIV  = 0xFF04
	TIMA = 0xFF05
	TMA  = 0xFF06
	TAC  = 0xFF07
	IF   = 0xFF0F
	IE   = 0xFFFF

	//
	// Defaults
	//

	JOYPDEFAULT = 0x0f
	NR10DEFAULT = 0x80
	NR11DEFAULT = 0xbf
	NR12DEFAULT = 0xf3
	NR13DEFAULT = 0xbf
	NR21DEFAULT = 0x3f
	NR24DEFAULT = 0xbf
	NR30DEFAULT = 0x7f
	NR31DEFAULT = 0xff
	NR32DEFAULT = 0x9f
	NR34DEFAULT = 0xbf
	NR41DEFAULT = 0xff
	NR44DEFAULT = 0xbf
	NR50DEFAULT = 0x77
	NR51DEFAULT = 0xf3
	NR52DEFAULT = 0xf1
	LCDCDEFAULT = 0x91
	BGPDEFAULT  = 0xfc
	OBP0DEFAULT = 0xff
	OBP1DEFAULT = 0xff
)

// Memory allows read and write access to memory
type Memory struct {
	mbc               *mbc
	VideoRAM          [0x2000]byte
	internalRAM       [0x2000]byte
	OAM               [0xa0]byte
	zeroPage          [0x8f]byte
	WriteNotification WriteNotification
	oamRunning        bool
	oamCycle          uint16
	oamBaseAddr       uint16
	oamRead           uint8
	IE                byte
	IF                byte
	LCDC              byte
	LY                byte
	LYC               byte
	SCX               byte
	SCY               byte
	STAT              byte
	WX                byte
	WY                byte
	JOYP              byte
	DirectionInput    uint8 // JOYP
	ButtonInput       uint8 // JOYP
	timer             *timer.Timer
	sbWriter          io.Writer
}

// WriteNotification provides a mechanism to notify other subsystems about memory writes
type WriteNotification interface {
	WriteToVideoRAM(addr uint16)
}

// NewMemory creates the memory struct and initializes it with ROM contents and default values
func NewMemory(rom []byte, sbWriter io.Writer, timer *timer.Timer) *Memory {
	if sbWriter == nil {
		sbWriter = ioutil.Discard
	}
	return &Memory{
		mbc:            newMBC(rom),
		LCDC:           LCDCDEFAULT,
		DirectionInput: JOYPDEFAULT,
		ButtonInput:    JOYPDEFAULT,
		timer:          timer,
		sbWriter:       sbWriter,
	}
}

// ExecuteMachineCycle updates the OAM after a machine cycle
func (m *Memory) ExecuteMachineCycle() {
	if m.oamRunning {
		if m.oamCycle == 0 {
			// Setup
		} else if m.oamCycle == 1 {
			m.oamRead = m.Read(m.oamBaseAddr)
		} else if m.oamCycle == 161 {
			m.OAM[159] = m.oamRead
			m.oamRunning = false
		} else {
			m.OAM[m.oamCycle-2] = m.oamRead
			m.oamRead = m.Read(m.oamBaseAddr + m.oamCycle - 1)
		}
		m.oamCycle++
	}
}

func (m *Memory) startOAM(value uint8) {
	m.oamRunning = true
	m.oamCycle = 0
	m.oamBaseAddr = uint16(value) << 8
	if m.oamBaseAddr >= 0xe000 {
		m.oamBaseAddr -= 0x2000
	}
}

func (m *Memory) readJOYP() uint8 {
	// Bit 5 - P15 Select Button Keys      (0=Select)
	// Bit 4 - P14 Select Direction Keys   (0=Select)
	if m.JOYP&0x10 == 0 {
		return m.JOYP&0xf0 | m.DirectionInput&0x0f
	}
	if m.JOYP&0x20 == 0 {
		return m.JOYP&0xf0 | m.ButtonInput&0x0f
	}
	return m.JOYP | 0x0f
}

// Read a byte from the chosen memory location
func (m *Memory) Read(addr uint16) byte {
	switch {
	case addr < 0x8000:
		return m.mbc.read(addr)
	case addr < 0xa000:
		return m.VideoRAM[addr-0x8000]
	case addr < 0xc000:
		return m.mbc.read(addr)
	case addr < 0xe000:
		return m.internalRAM[addr-0xc000]
	case addr < 0xfe00:
		return m.internalRAM[addr-0xe000]
	case addr < 0xfea0:
		if m.oamRunning {
			return 0xff
		}
		return m.OAM[addr-0xfe00]
	case addr < 0xff00:
		return 0 // Unusable region
	case addr == DMA:
		return uint8(m.oamBaseAddr >> 8)
	case addr == IE:
		return m.IE
	case addr == IF:
		return m.IF | 0xe0 // Top 3 bits are always high
	case addr == LCDC:
		return m.LCDC
	case addr == STAT:
		return m.STAT
	case addr == SCY:
		return m.SCY
	case addr == SCX:
		return m.SCX
	case addr == LY:
		return m.LY
	case addr == LYC:
		return m.LYC
	case addr == WY:
		return m.WY
	case addr == WX:
		return m.WX
		// case BGP:
		// case OBP0:
		// case OBP1:
		// case NR10:
		// case NR11:
		// case NR12:
		// case NR13:
		// case NR14:
		// case NR21:
		// case NR22:
		// case NR23:
		// case NR24:
		// case NR30:
		// case NR31:
		// case NR32:
		// case NR33:
		// case NR34:
		// case NR41:
		// case NR42:
		// case NR43:
		// case NR44:
		// case NR50:
		// case NR51:
		// case NR52:
	case addr == JOYP:
		return m.readJOYP()
		// case SB:
		// case SC:
	case addr == DIV:
		return m.timer.DIV()
	case addr == TIMA:
		return m.timer.TIMA()
	case addr == TMA:
		return m.timer.TMA()
	case addr == TAC:
		return m.timer.TAC()
	case addr < 0xff80:
		return 0xff // Default of a non-hardware register is read
	case addr < 0xffff:
		return m.zeroPage[addr-0xff80]
	case addr == 0xffff:
		return m.IE
	default:
		panic(fmt.Sprintf("Read failed: 0x%04x", addr))
	}
}

// Write a byte to the chosen memory location
func (m *Memory) Write(addr uint16, value byte) {
	switch {
	case addr < 0x8000:
		m.mbc.write(addr, value)
	case addr < 0xa000:
		if m.WriteNotification != nil {
			m.WriteNotification.WriteToVideoRAM(addr)
		}
		m.VideoRAM[addr-0x8000] = value
	case addr < 0xc000:
		m.mbc.write(addr, value)
	case addr < 0xe000:
		m.internalRAM[addr-0xc000] = value
	case addr < 0xfe00:
		m.internalRAM[addr-0xe000] = value
	case addr < 0xfea0:
		m.OAM[addr-0xfe00] = value
	case addr < 0xff00:
		// Unusable region
	case addr == DMA:
		m.startOAM(value)
	case addr == IE:
		m.IE = value
	case addr == IF:
		m.IF = value
	case addr == LCDC:
		m.LCDC = value
	case addr == STAT:
		m.STAT = value
	case addr == SCY:
		m.SCY = value
	case addr == SCX:
		m.SCX = value
	case addr == LY:
		m.LY = value
	case addr == LYC:
		m.LYC = value
	case addr == WY:
		m.WY = value
	case addr == WX:
		m.WX = value
	case addr == BGP:
		// FIXME palette support
	case addr == OBP0:
		// FIXME sprite palette support
	case addr == OBP1:
		// FIXME sprite palette support
	case addr == NR10:
		// FIXME sound support
	case addr == NR11:
		// FIXME sound support
	case addr == NR12:
		// FIXME sound support
	case addr == NR13:
		// FIXME sound support
	case addr == NR14:
		// FIXME sound support
	case addr == NR21:
		// FIXME sound support
	case addr == NR22:
		// FIXME sound support
	case addr == NR23:
		// FIXME sound support
	case addr == NR24:
		// FIXME sound support
	case addr == NR30:
		// FIXME sound support
	case addr == NR31:
		// FIXME sound support
	case addr == NR32:
		// FIXME sound support
	case addr == NR33:
		// FIXME sound support
	case addr == NR34:
		// FIXME sound support
	case addr == NR41:
		// FIXME sound support
	case addr == NR42:
		// FIXME sound support
	case addr == NR43:
		// FIXME sound support
	case addr == NR44:
		// FIXME sound support
	case addr == NR50:
		// FIXME sound support
	case addr == NR51:
		// FIXME sound support
	case addr == NR52:
		// FIXME sound support
	case addr == JOYP:
		m.JOYP = value
	case addr == SB:
		m.sbWriter.Write([]byte{value})
	case addr == SC:
		// FIXME serial bus support
	case addr == DIV:
		m.timer.Reset()
	case addr == TIMA:
		m.timer.WriteTIMA(value)
	case addr == TMA:
		m.timer.WriteTMA(value)
	case addr == TAC:
		m.timer.WriteTAC(value)
	case addr < 0xff80:
	// Do nothing if a non-hardware register is written
	case addr < 0xffff:
		m.zeroPage[addr-0xff80] = value
	case addr == IE:
		m.IE = value
	default:
		panic(fmt.Sprintf("Write failed: 0x%04x", addr))
	}
}

// CartRAM returns the contents of cartridge RAM, which is useful for verifing test results
func (m *Memory) CartRAM() [][0x2000]byte {
	return m.mbc.ram
}
