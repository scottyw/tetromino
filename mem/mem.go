package mem

import (
	"io/ioutil"
	"os"
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
	IE   = 0xFFFF
	IF   = 0xFF0F
)

// Memory allows read and write access to memory
type Memory struct {
	mem                  []byte
	storesDirectionInput bool
	DirectionInput       *uint8
	ButtonInput          *uint8
}

// NewMemory creates the memory and initializes it with ROM contents and default values
func NewMemory() *Memory {
	rom, err := ioutil.ReadFile(os.Getenv("ROM_FILENAME"))
	if err != nil {
		panic(err)
	}
	mem := make([]byte, 65536)
	copy(mem, rom)
	mem[0xff00] = 0x0f
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
	direction := uint8(0x0f)
	button := uint8(0x0f)
	return &Memory{
		mem:            mem,
		DirectionInput: &direction,
		ButtonInput:    &button,
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

func (mem *Memory) joypRead() uint8 {
	if mem.storesDirectionInput {
		return *mem.DirectionInput
	}
	return *mem.ButtonInput
}

func (mem *Memory) joypWrite(value uint8) {
	// Bit 5 - P15 Select Button Keys      (0=Select)
	// Bit 4 - P14 Select Direction Keys   (0=Select)
	if value&0x20 == 0 {
		mem.storesDirectionInput = false
	} else if value&0x10 == 0 {
		mem.storesDirectionInput = true
	}
}

func (mem *Memory) dma(addrPrefix uint8) {
	srcBaseAddr := uint16(addrPrefix) << 8
	for i := uint16(0x00); i < 0x0100; i++ {
		mem.mem[0xfe00+i] = mem.mem[srcBaseAddr+i]
	}
}

// Read a byte from the chosen memory location
func (mem *Memory) Read(addr uint16) byte {
	switch addr {
	case JOYP:
		return mem.joypRead()
	case DMA:
		return 0
	}
	return mem.mem[addr]
}

// Write a byte to the chosen memory location
func (mem *Memory) Write(addr uint16, value byte) {
	switch addr {
	case JOYP:
		mem.joypWrite(value)
	case DMA:
		mem.dma(value)
	default:
		mem.mem[addr] = value
	}
}

// ReadRegion of memory
func (mem *Memory) ReadRegion(startAddr, length uint16) []byte {
	return mem.mem[startAddr : startAddr+length]
}

// GenerateCrashReport writes the contents of the whole address space to file
func (mem *Memory) GenerateCrashReport() {
	if r := recover(); r != nil {
		ioutil.WriteFile("memory.bin", mem.mem, 0644)
		panic(r)
	}
}

// RaiseInterrupt updates the IF register with the supplied pattern
func (mem *Memory) RaiseInterrupt(pattern uint8) {
	mem.mem[IF] |= 0x01
}

// ResetInterrupt updates the IF register with the supplied pattern
func (mem *Memory) ResetInterrupt(pattern uint8) {
	mem.mem[IF] &^= 0x01
}
