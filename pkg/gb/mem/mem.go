package mem

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/scottyw/tetromino/pkg/gb/audio"
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
)

// Memory allows read and write access to memory
type Memory struct {

	// HW registers
	IE   byte
	IF   byte
	LCDC byte
	LY   byte
	LYC  byte
	SCX  byte
	SCY  byte
	STAT byte
	WX   byte
	WY   byte
	JOYP byte
	SB   byte
	SC   byte
	BGP  byte
	OBP0 byte
	OBP1 byte

	// Implementation
	mbc               *mbc
	VideoRAM          [0x2000]byte
	internalRAM       [0x2000]byte
	OAM               [0xa0]byte
	zeroPage          [0x8f]byte
	WriteNotification WriteNotification
	biosDisabled      bool
	oamRunning        bool
	oamCycle          uint16
	oamBaseAddr       uint16
	oamRead           uint8
	DirectionInput    uint8 // JOYP
	ButtonInput       uint8 // JOYP
	timer             *timer.Timer
	audio             *audio.Audio
	sbWriter          io.Writer
}

// WriteNotification provides a mechanism to notify other subsystems about memory writes
type WriteNotification interface {
	WriteToVideoRAM(addr uint16)
}

// NewMemory creates the memory struct and initializes it with ROM contents and default values
func NewMemory(rom []byte, sbWriter io.Writer, timer *timer.Timer, audio *audio.Audio) *Memory {
	if sbWriter == nil {
		sbWriter = ioutil.Discard
	}
	return &Memory{

		// HW register defaults
		IE:   0x00,
		IF:   0x01,
		LCDC: 0x91,
		LY:   0x00,
		LYC:  0x00,
		SCX:  0x00,
		SCY:  0x00,
		STAT: 0x00,
		WX:   0x00,
		WY:   0x00,
		JOYP: 0x0f,
		SB:   0x00,
		SC:   0x7e,
		BGP:  0xfc,
		OBP0: 0xff,
		OBP1: 0xff,

		// Implementation
		mbc:            newMBC(rom),
		DirectionInput: 0x0f,
		ButtonInput:    0x0f,
		timer:          timer,
		audio:          audio,
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
	case addr < 0x0100:
		if m.biosDisabled {
			return m.mbc.read(addr)
		}
		return bios[addr]
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
		return m.STAT | 0x80 // First bit is always high
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
	case addr == BGP:
		return m.BGP
	case addr == OBP0:
		return m.OBP0
	case addr == OBP1:
		return m.OBP1
	case addr == NR10:
		return m.audio.ReadNR10()
	case addr == NR11:
		return m.audio.ReadNR11()
	case addr == NR12:
		return m.audio.ReadNR12()
	case addr == NR13:
		return m.audio.ReadNR13()
	case addr == NR14:
		return m.audio.ReadNR14()
	case addr == NR21:
		return m.audio.ReadNR21()
	case addr == NR22:
		return m.audio.ReadNR22()
	case addr == NR23:
		return m.audio.ReadNR23()
	case addr == NR24:
		return m.audio.ReadNR24()
	case addr == NR30:
		return m.audio.ReadNR30()
	case addr == NR31:
		return m.audio.ReadNR31()
	case addr == NR32:
		return m.audio.ReadNR32()
	case addr == NR33:
		return m.audio.ReadNR33()
	case addr == NR34:
		return m.audio.ReadNR34()
	case addr == NR41:
		return m.audio.ReadNR41()
	case addr == NR42:
		return m.audio.ReadNR42()
	case addr == NR43:
		return m.audio.ReadNR43()
	case addr == NR44:
		return m.audio.ReadNR44()
	case addr == NR50:
		return m.audio.ReadNR50()
	case addr == NR51:
		return m.audio.ReadNR51()
	case addr == NR52:
		return m.audio.ReadNR52()
	case addr == JOYP:
		return m.readJOYP() | 0xc0 // First 2 bits are always high
	case addr == SB:
		return m.SB
	case addr == SC:
		return m.SC
	case addr == DIV:
		return m.timer.DIV()
	case addr == TIMA:
		return m.timer.TIMA()
	case addr == TMA:
		return m.timer.TMA()
	case addr == TAC:
		return m.timer.TAC() | 0xf8 // First 5 bits are always high
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
		m.audio.WriteNR10(value)
	case addr == NR11:
		m.audio.WriteNR11(value)
	case addr == NR12:
		m.audio.WriteNR12(value)
	case addr == NR13:
		m.audio.WriteNR13(value)
	case addr == NR14:
		m.audio.WriteNR14(value)
	case addr == NR21:
		m.audio.WriteNR21(value)
	case addr == NR22:
		m.audio.WriteNR22(value)
	case addr == NR23:
		m.audio.WriteNR23(value)
	case addr == NR24:
		m.audio.WriteNR24(value)
	case addr == NR30:
		m.audio.WriteNR30(value)
	case addr == NR31:
		m.audio.WriteNR31(value)
	case addr == NR32:
		m.audio.WriteNR32(value)
	case addr == NR33:
		m.audio.WriteNR33(value)
	case addr == NR34:
		m.audio.WriteNR34(value)
	case addr == NR41:
		m.audio.WriteNR41(value)
	case addr == NR42:
		m.audio.WriteNR42(value)
	case addr == NR43:
		m.audio.WriteNR43(value)
	case addr == NR44:
		m.audio.WriteNR44(value)
	case addr == NR50:
		m.audio.WriteNR50(value)
	case addr == NR51:
		m.audio.WriteNR51(value)
	case addr == NR52:
		m.audio.WriteNR52(value)
	case addr == JOYP:
		m.JOYP = value
	case addr == SB:
		_, err := m.sbWriter.Write([]byte{value})
		if err != nil {
			panic(fmt.Sprintf("Write to SB failed: %v", err))
		}
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

// DisableBIOS disables access to the BIOS once it has run at start time
func (m *Memory) DisableBIOS() {
	m.biosDisabled = true
}
