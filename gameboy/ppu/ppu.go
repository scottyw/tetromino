package ppu

import (
	"image"

	"github.com/scottyw/tetromino/gameboy/interrupts"
)

type PPU struct {
	interrupts *interrupts.Interrupts

	lcdc uint8
	ly   uint8
	lyc  uint8
	scx  uint8
	scy  uint8
	stat uint8
	wx   uint8
	wy   uint8
	bgp  uint8
	obp0 uint8
	obp1 uint8

	videoRAM [0x2000]byte
	oam      *[0xa0]byte
	tick     int
	debug    bool

	// LCD
	tileCache      [384]*[8][8]uint8
	previousBg     [32][32]uint16
	previousWindow [32][32]uint16
	bg             [256][256]uint8
	window         [256][256]uint8
	sprites        [144][160]uint8
	frame          *image.RGBA
}

func New(oam *[0xa0]byte, interrupts *interrupts.Interrupts, debug bool) *PPU {
	ppu := &PPU{
		oam:        oam,
		interrupts: interrupts,
		lcdc:       0x91,
		ly:         0x00,
		lyc:        0x00,
		scx:        0x00,
		scy:        0x00,
		stat:       0x00,
		wx:         0x00,
		wy:         0x00,
		bgp:        0xfc,
		obp0:       0xff,
		obp1:       0xff,
		frame:      image.NewRGBA(image.Rect(0, 0, 256, 256)),
		debug:      debug,
	}
	return ppu
}

func (ppu *PPU) ReadVideoRAM(addr uint16) uint8 {
	return ppu.videoRAM[addr-0x8000]
}

func (ppu *PPU) WriteVideoRAM(addr uint16, value uint8) {
	if addr < 0x9800 {
		tileNumber := (addr - 0x8000) / 16
		ppu.tileCache[tileNumber] = nil
	}
	ppu.videoRAM[addr-0x8000] = value
}
