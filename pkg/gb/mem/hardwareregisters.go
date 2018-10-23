package mem

import (
	"io"
	"io/ioutil"

	"github.com/scottyw/tetromino/pkg/gb/timer"
)

// Register constants
const (
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

// HardwareRegisters represents hardware registers between 0xff00 and 0xff7f
type HardwareRegisters struct {
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

	// JOYP
	DirectionInput uint8
	ButtonInput    uint8

	// Timer
	timer *timer.Timer

	// Misc
	sbWriter io.Writer
}

// NewHardwareRegisters creates a new representation of the hardware registers
func NewHardwareRegisters(timer *timer.Timer, sbWriter io.Writer) *HardwareRegisters {
	if sbWriter == nil {
		sbWriter = ioutil.Discard
	}
	return &HardwareRegisters{
		LCDC:           LCDCDEFAULT,
		DirectionInput: JOYPDEFAULT,
		ButtonInput:    JOYPDEFAULT,
		timer:          timer,
		sbWriter:       sbWriter,
	}
}

func (mem *Memory) readHardwareRegisters(addr uint16) uint8 {
	switch addr {
	case DMA:
		return 0
	case IE:
		return mem.hwr.IE
	case IF:
		return mem.hwr.IF | 0xe0 // Top 3 bits are always high
	case LCDC:
		return mem.hwr.LCDC
	case STAT:
		return mem.hwr.STAT
	case SCY:
		return mem.hwr.SCY
	case SCX:
		return mem.hwr.SCX
	case LY:
		return mem.hwr.LY
	case LYC:
		return mem.hwr.LYC
	case WY:
		return mem.hwr.WY
	case WX:
		return mem.hwr.WX
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
	case JOYP:
		return mem.hwr.readJOYP()
	// case SB:
	// case SC:
	case DIV:
		return mem.hwr.timer.DIV()
	case TIMA:
		return mem.hwr.timer.TIMA()
	case TMA:
		return mem.hwr.timer.TMA()
	case TAC:
		return mem.hwr.timer.TAC()
	default:
		return 0xff
	}
}

func (mem *Memory) writeHardwareRegisters(addr uint16, value uint8) {
	switch addr {
	case DMA:
		mem.dma(value)
	case IE:
		mem.hwr.IE = value
	case IF:
		mem.hwr.IF = value
	case LCDC:
		mem.hwr.LCDC = value
	case STAT:
		mem.hwr.STAT = value
	case SCY:
		mem.hwr.SCY = value
	case SCX:
		mem.hwr.SCX = value
	case LY:
		mem.hwr.LY = value
	case LYC:
		mem.hwr.LYC = value
	case WY:
		mem.hwr.WY = value
	case WX:
		mem.hwr.WX = value
	case BGP:
		// FIXME palette support
	case OBP0:
		// FIXME sprite palette support
	case OBP1:
		// FIXME sprite palette support
	case NR10:
		// FIXME sound support
	case NR11:
		// FIXME sound support
	case NR12:
		// FIXME sound support
	case NR13:
		// FIXME sound support
	case NR14:
		// FIXME sound support
	case NR21:
		// FIXME sound support
	case NR22:
		// FIXME sound support
	case NR23:
		// FIXME sound support
	case NR24:
		// FIXME sound support
	case NR30:
		// FIXME sound support
	case NR31:
		// FIXME sound support
	case NR32:
		// FIXME sound support
	case NR33:
		// FIXME sound support
	case NR34:
		// FIXME sound support
	case NR41:
		// FIXME sound support
	case NR42:
		// FIXME sound support
	case NR43:
		// FIXME sound support
	case NR44:
		// FIXME sound support
	case NR50:
		// FIXME sound support
	case NR51:
		// FIXME sound support
	case NR52:
		// FIXME sound support
	case JOYP:
		mem.hwr.JOYP = value
	case SB:
		mem.hwr.sbWriter.Write([]byte{value})
	case SC:
		// FIXME serial bus support
	case DIV:
		mem.hwr.timer.Reset()
	case TIMA:
		mem.hwr.timer.WriteTIMA(value)
	case TMA:
		mem.hwr.timer.WriteTMA(value)
	case TAC:
		mem.hwr.timer.WriteTAC(value)
	default:
		// Do nothing
	}
}

func (r *HardwareRegisters) readJOYP() uint8 {

	// Bit 5 - P15 Select Button Keys      (0=Select)
	// Bit 4 - P14 Select Direction Keys   (0=Select)

	if r.JOYP&0x10 == 0 {
		return r.JOYP&0xf0 | r.DirectionInput&0x0f
	}

	if r.JOYP&0x20 == 0 {
		return r.JOYP&0xf0 | r.ButtonInput&0x0f
	}

	return r.JOYP | 0x0f
}

func (mem *Memory) dma(addrPrefix uint8) {
	srcBaseAddr := uint16(addrPrefix) << 8
	for i := uint16(0x00); i < 0x0a0; i++ {
		mem.OAM[i] = mem.Read(srcBaseAddr + i)
	}
}
