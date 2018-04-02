package mem

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
	//
	// [0x00] = 0x0f
	// [0x10] = 0x80
	// [0x11] = 0xbf
	// [0x12] = 0xf3
	// [0x14] = 0xbf
	// [0x16] = 0x3f
	// [0x19] = 0xbf
	// [0x1a] = 0x7f
	// [0x1b] = 0xff
	// [0x1c] = 0x9f
	// [0x1e] = 0xbf
	// [0x20] = 0xff
	// [0x23] = 0xbf
	// [0x24] = 0x77
	// [0x25] = 0xf3
	// [0x26] = 0xf1
	// [0x40] = 0x91
	// [0x47] = 0xfc
	// [0x48] = 0xff
	// [0x49] = 0xff

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

	// JOYP
	storingDirectionInput bool
	DirectionInput        uint8
	ButtonInput           uint8
}

// NewHardwareRegisters creates a new representation of the hardware registers
func NewHardwareRegisters() *HardwareRegisters {
	return &HardwareRegisters{
		DirectionInput: 0x0f,
		ButtonInput:    0x0f,
	}
}

func (mem *Memory) readHardwareRegisters(addr uint16) uint8 {
	switch addr {
	case DMA:
		return 0
	case IE:
		return mem.hwr.IE
	case IF:
		return mem.hwr.IF
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
		return mem.hwr.joypRead()
	// case SB:
	// case SC:
	// case DIV:
	// case TIMA:
	// case TMA:
	// case TAC:
	default:
		// fmt.Printf("No read support for hardware register 0x%04x\n", addr)
	}
	return 0 // Temporary default value
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
		mem.hwr.joypWrite(value)
	// case SB:
	// case SC:
	// case DIV:
	// case TIMA:
	// case TMA:
	// case TAC:
	default:
		// fmt.Printf("No write support for hardware register 0x%04x\n", addr)
	}
}

func (r *HardwareRegisters) joypRead() uint8 {
	if r.storingDirectionInput {
		return r.DirectionInput
	}
	return r.ButtonInput
}

func (r *HardwareRegisters) joypWrite(value uint8) {
	// Bit 5 - P15 Select Button Keys      (0=Select)
	// Bit 4 - P14 Select Direction Keys   (0=Select)
	if value&0x20 == 0 {
		r.storingDirectionInput = false
	} else if value&0x10 == 0 {
		r.storingDirectionInput = true
	}
}

func (mem *Memory) dma(addrPrefix uint8) {
	srcBaseAddr := uint16(addrPrefix) << 8
	for i := uint16(0x00); i < 0x0a0; i++ {
		mem.oam[i] = mem.Read(srcBaseAddr + i)
	}
}
