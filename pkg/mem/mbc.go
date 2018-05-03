package mem

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/scottyw/tetromino/pkg/options"
)

type mbc interface {
	read(uint16) uint8
	write(uint16, uint8)
}

type mbc1 struct {
	rom        []byte
	bank       uint16
	romBank    uint8
	ramBank    uint8
	romMode    bool
	ramEnabled bool
}

func (m *mbc1) accessRom(addr uint16) uint8 {
	if int(addr) > len(m.rom) {
		if len(m.rom) == 0 {
			fmt.Println("No ROM filename specified ...")
			os.Exit(1)
		}
		panic(fmt.Sprintf("Attempt to access address 0x%04x but MBC1 ROM only has length 0x%04x", addr, len(m.rom)))
	}
	return m.rom[addr]
}

func (m *mbc1) read(addr uint16) uint8 {
	switch {
	case addr < 0x4000:
		return m.accessRom(addr)
	case addr < 0x8000:
		return m.accessRom((m.bank-1)*0x4000 + addr)
	case addr < 0xa000:
		panic(fmt.Sprintf("MBC1 has no mapping for address 0x%04x", addr))
	case addr < 0xc000:
		panic(fmt.Sprintf("MBC1 RAM Bank not implemented! Address: 0x%04x", addr))
	default:
		panic(fmt.Sprintf("MBC1 has no mapping for address 0x%04x", addr))
	}
}

func (m *mbc1) updateBank() {
	if m.romMode {
		m.bank = uint16(m.ramBank<<5 | m.romBank)
	} else {
		m.bank = uint16(m.romBank)
	}
}

func (m *mbc1) write(addr uint16, value uint8) {
	switch {
	case addr < 0x2000:
		//		m.ramEnabled = value&0x0a > 0
	case addr < 0x4000:
		m.romBank = value & 0x1f
		if m.romBank == 0 {
			m.romBank++
		}
	case addr < 0x6000:
		m.ramBank = value & 0x03
	case addr < 0x8000:
		m.romMode = value&0x01 == 0
	default:
		panic(fmt.Sprintf("MBC1 has no mapping for address 0x%04x", addr))
	}
	m.updateBank()
}

func newMBC() mbc {
	var rom []byte
	if *options.RomFilename != "" {
		var err error
		rom, err = ioutil.ReadFile(*options.RomFilename)
		if err != nil {
			panic(fmt.Sprintf("Failed to read the ROM file at \"%s\" (%v)", *options.RomFilename, err))
		}
	}
	return &mbc1{
		rom:  rom,
		bank: 1,
	}
}
