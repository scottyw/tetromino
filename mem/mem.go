package mem

import (
	"io/ioutil"
	"os"
)

// Memory allows read and write access to memory
type Memory interface {
	Read(uint16) *byte
	GenerateCrashReport()
}

type memory struct {
	mem []byte
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
	return memory{mem: mem}
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

// Read a byte from the chosen memory location
func (mem memory) Read(addr uint16) *byte {
	// if addr >= 0x8000 {
	// 	fmt.Printf("DEBUG: Read %s - 0x%04x\n", region(addr), addr)
	// }
	return &mem.mem[addr]
}

// GenerateCrashReport writes the contents of the whole address space to file
func (mem memory) GenerateCrashReport() {
	if r := recover(); r != nil {
		ioutil.WriteFile("memory.bin", mem.mem, 0644)
		// drawWindow()
		panic(r)
	}
}
