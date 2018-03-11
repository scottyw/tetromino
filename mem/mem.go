package mem

import (
	"io/ioutil"
	"os"
)

// Memory allows read and write access to memory
type Memory struct {
	mem  []byte
	LCDC *byte // FF40 - LCDC - LCD Control (R/W)
	STAT *byte // FF41 - STAT - LCDC Status   (R/W)
	SCY  *byte // FF42 - SCY - Scroll Y   (R/W)
	SCX  *byte // FF43 - SCX - Scroll X   (R/W)
	LY   *byte // FF44 - LY - LCDC Y-Coordinate (R)
	LYC  *byte // FF45 - LYC - LY Compare  (R/W)
	WY   *byte // FF4A - WY - Window Y Position (R/W)
	WX   *byte // FF4B - WX - Window X Position minus 7 (R/W)
	BGP  *byte // FF47 - BGP - BG Palette Data  (R/W) - Non CGB Mode Only
	OBP0 *byte // FF48 - OBP0 - Object Palette 0 Data (R/W) - Non CGB Mode Only
	OBP1 *byte // FF49 - OBP1 - Object Palette 1 Data (R/W) - Non CGB Mode Only
	DMA  *byte // FF46 - DMA - DMA Transfer and Start Address (W)
	NR10 *byte // FF10 - NR10 - Channel 1 Sweep register (R/W)
	NR11 *byte // FF11 - NR11 - Channel 1 Sound length/Wave pattern duty (R/W)
	NR12 *byte // FF12 - NR12 - Channel 1 Volume Envelope (R/W)
	NR13 *byte // FF13 - NR13 - Channel 1 Frequency lo (Write Only)
	NR14 *byte // FF14 - NR14 - Channel 1 Frequency hi (R/W)
	NR21 *byte // FF16 - NR21 - Channel 2 Sound Length/Wave Pattern Duty (R/W)
	NR22 *byte // FF17 - NR22 - Channel 2 Volume Envelope (R/W)
	NR23 *byte // FF18 - NR23 - Channel 2 Frequency lo data (W)
	NR24 *byte // FF19 - NR24 - Channel 2 Frequency hi data (R/W)
	NR30 *byte // FF1A - NR30 - Channel 3 Sound on/off (R/W)
	NR31 *byte // FF1B - NR31 - Channel 3 Sound Length
	NR32 *byte // FF1C - NR32 - Channel 3 Select output level (R/W)
	NR33 *byte // FF1D - NR33 - Channel 3 Frequency's lower data (W)
	NR34 *byte // FF1E - NR34 - Channel 3 Frequency's higher data (R/W)
	NR41 *byte // FF20 - NR41 - Channel 4 Sound Length (R/W)
	NR42 *byte // FF21 - NR42 - Channel 4 Volume Envelope (R/W)
	NR43 *byte // FF22 - NR43 - Channel 4 Polynomial Counter (R/W)
	NR44 *byte // FF23 - NR44 - Channel 4 Counter/consecutive; Inital (R/W)
	NR50 *byte // FF24 - NR50 - Channel control / ON-OFF / Volume (R/W)
	NR51 *byte // FF25 - NR51 - Selection of Sound output terminal (R/W)
	NR52 *byte // FF26 - NR52 - Sound on/off
	JOYP *byte // FF00 - P1/JOYP - Joypad (R/W)
	SB   *byte // FF01 - SB - Serial transfer data (R/W)
	SC   *byte // FF02 - SC - Serial Transfer Control  (R/W)
	DIV  *byte // FF04 - DIV - Divider Register (R/W)
	TIMA *byte // FF05 - TIMA - Timer counter (R/W)
	TMA  *byte // FF06 - TMA - Timer Modulo (R/W)
	TAC  *byte // FF07 - TAC - Timer Control (R/W)
	IE   *byte // FFFF - IE - Interrupt Enable (R/W)
	IF   *byte // FF0F - IF - Interrupt Flag (R/W)
}

// NewMemory creates the memory and initializes it with ROM contents and default values
func NewMemory() Memory {
	rom, err := ioutil.ReadFile(os.Getenv("ROM_FILENAME"))
	if err != nil {
		panic(err)
	}
	mem := make([]byte, 65536)
	copy(mem, rom)
	mem[0xff05] = 0x00
	mem[0xff06] = 0x00
	mem[0xff07] = 0x00
	mem[0xff10] = 0x80
	mem[0xff11] = 0xbf
	mem[0xff12] = 0xf3
	mem[0xff14] = 0xbf
	mem[0xff16] = 0x3f
	mem[0xff17] = 0x00
	mem[0xff19] = 0xbf
	mem[0xff1a] = 0x7f
	mem[0xff1b] = 0xff
	mem[0xff1c] = 0x9f
	mem[0xff1e] = 0xbf
	mem[0xff20] = 0xff
	mem[0xff21] = 0x00
	mem[0xff22] = 0x00
	mem[0xff23] = 0xbf
	mem[0xff24] = 0x77
	mem[0xff25] = 0xf3
	mem[0xff26] = 0xf1
	mem[0xff40] = 0x91
	mem[0xff42] = 0x00
	mem[0xff43] = 0x00
	mem[0xff45] = 0x00
	mem[0xff47] = 0xfc
	mem[0xff48] = 0xff
	mem[0xff49] = 0xff
	mem[0xff4a] = 0x00
	mem[0xff4b] = 0x00
	mem[0xffff] = 0x00
	return Memory{
		mem:  mem,
		LCDC: &mem[0xFF40], // FF40] - LCDC - LCD Control (R/W)
		STAT: &mem[0xFF41], // FF41 - STAT - LCDC Status   (R/W)
		SCY:  &mem[0xFF42], // FF42 - SCY - Scroll Y   (R/W)
		SCX:  &mem[0xFF43], // FF43 - SCX - Scroll X   (R/W)
		LY:   &mem[0xFF44], // FF44 - LY - LCDC Y-Coordinate (R)
		LYC:  &mem[0xFF45], // FF45 - LYC - LY Compare  (R/W)
		WY:   &mem[0xFF4A], // FF4A - WY - Window Y Position (R/W)
		WX:   &mem[0xFF4B], // FF4B - WX - Window X Position minus 7 (R/W)
		BGP:  &mem[0xFF47], // FF47 - BGP - BG Palette Data  (R/W) - Non CGB Mode Only
		OBP0: &mem[0xFF48], // FF48 - OBP0 - Object Palette 0 Data (R/W) - Non CGB Mode Only
		OBP1: &mem[0xFF49], // FF49 - OBP1 - Object Palette 1 Data (R/W) - Non CGB Mode Only
		DMA:  &mem[0xFF46], // FF46 - DMA - DMA Transfer and Start Address (W)
		NR10: &mem[0xFF10], // FF10 - NR10 - Channel 1 Sweep register (R/W)
		NR11: &mem[0xFF11], // FF11 - NR11 - Channel 1 Sound length/Wave pattern duty (R/W)
		NR12: &mem[0xFF12], // FF12 - NR12 - Channel 1 Volume Envelope (R/W)
		NR13: &mem[0xFF13], // FF13 - NR13 - Channel 1 Frequency lo (Write Only)
		NR14: &mem[0xFF14], // FF14 - NR14 - Channel 1 Frequency hi (R/W)
		NR21: &mem[0xFF16], // FF16 - NR21 - Channel 2 Sound Length/Wave Pattern Duty (R/W)
		NR22: &mem[0xFF17], // FF17 - NR22 - Channel 2 Volume Envelope (R/W)
		NR23: &mem[0xFF18], // FF18 - NR23 - Channel 2 Frequency lo data (W)
		NR24: &mem[0xFF19], // FF19 - NR24 - Channel 2 Frequency hi data (R/W)
		NR30: &mem[0xFF1A], // FF1A - NR30 - Channel 3 Sound on/off (R/W)
		NR31: &mem[0xFF1B], // FF1B - NR31 - Channel 3 Sound Length
		NR32: &mem[0xFF1C], // FF1C - NR32 - Channel 3 Select output level (R/W)
		NR33: &mem[0xFF1D], // FF1D - NR33 - Channel 3 Frequency's lower data (W)
		NR34: &mem[0xFF1E], // FF1E - NR34 - Channel 3 Frequency's higher data (R/W)
		NR41: &mem[0xFF20], // FF20 - NR41 - Channel 4 Sound Length (R/W)
		NR42: &mem[0xFF21], // FF21 - NR42 - Channel 4 Volume Envelope (R/W)
		NR43: &mem[0xFF22], // FF22 - NR43 - Channel 4 Polynomial Counter (R/W)
		NR44: &mem[0xFF23], // FF23 - NR44 - Channel 4 Counter/consecutive; Inital (R/W)
		NR50: &mem[0xFF24], // FF24 - NR50 - Channel control / ON-OFF / Volume (R/W)
		NR51: &mem[0xFF25], // FF25 - NR51 - Selection of Sound output terminal (R/W)
		NR52: &mem[0xFF26], // FF26 - NR52 - Sound on/off
		JOYP: &mem[0xFF00], // FF00 - P1/JOYP - Joypad (R/W)
		SB:   &mem[0xFF01], // FF01 - SB - Serial transfer data (R/W)
		SC:   &mem[0xFF02], // FF02 - SC - Serial Transfer Control  (R/W)
		DIV:  &mem[0xFF04], // FF04 - DIV - Divider Register (R/W)
		TIMA: &mem[0xFF05], // FF05 - TIMA - Timer counter (R/W)
		TMA:  &mem[0xFF06], // FF06 - TMA - Timer Modulo (R/W)
		TAC:  &mem[0xFF07], // FF07 - TAC - Timer Control (R/W)
		IE:   &mem[0xFFFF], // FFFF - IE - Interrupt Enable (R/W)
		IF:   &mem[0xFF0F], // FF0F - IF - Interrupt Flag (R/W)
	}
}

// Debug function
func region(addr uint16) string {
	switch {
	case addr == 0xffff:
		return "IE Register"
	case addr >= 0xff80:
		return "RAM (FF80)"
	case addr >= 0xff4c:
		return "EMPTY"
	case addr >= 0xff00:
		return "IO"
	case addr >= 0xfea0:
		return "EMPTY"
	case addr >= 0xfe00:
		return "Sprite Attribute Mem"
	case addr >= 0xe000:
		return "RAM (echo)"
	case addr >= 0xc000:
		return "RAM"
	case addr >= 0xa000:
		return "Switchable RAM"
	case addr >= 0x8000:
		return "Video RAM"
	case addr >= 0x4000:
		return "Switchable ROM"
	default:
		return "ROM"
	}
}

func (mem Memory) dma(addrPrefix uint8) {
	srcBaseAddr := uint16(addrPrefix) << 8
	for i := uint16(0x00); i < 0x0100; i++ {
		mem.mem[0xfe00+i] = mem.mem[srcBaseAddr+i]
	}
}

// Read a byte from the chosen memory location
func (mem Memory) Read(addr uint16) byte {
	// if addr >= 0x8000 {
	// 	fmt.Printf("DEBUG: Read %s - 0x%04x\n", region(addr), addr)
	// }
	return mem.mem[addr]
}

// Write a byte to the chosen memory location
func (mem Memory) Write(addr uint16, value byte) {
	switch addr {
	case 0xff46:
		mem.dma(value)
	default:
		mem.mem[addr] = value
	}
}

// ReadRegion of memory
func (mem Memory) ReadRegion(startAddr, length uint16) []byte {
	return mem.mem[startAddr : startAddr+length]
}

// GenerateCrashReport writes the contents of the whole address space to file
func (mem Memory) GenerateCrashReport() {
	if r := recover(); r != nil {
		ioutil.WriteFile("memory.bin", mem.mem, 0644)
		panic(r)
	}
}
