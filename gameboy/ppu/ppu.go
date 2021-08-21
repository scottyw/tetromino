package ppu

import (
	"image"

	"github.com/scottyw/tetromino/gameboy/interrupts"
	"github.com/scottyw/tetromino/gameboy/oam"
)

type PPU struct {

	//LCDC
	enabled           bool
	highWindowTileMap bool
	windowEnabled     bool
	lowTileData       bool
	highBgTileMap     bool
	spritesLarge      bool
	spritesEnabled    bool
	bgEnabled         bool

	// STAT
	coincidenceInterrupt bool
	oamInterrupt         bool
	vblankInterrupt      bool
	hlankInterrupt       bool
	coincidence          bool
	mode                 uint8

	// BGP
	bgpColour0 uint8
	bgpColour1 uint8
	bgpColour2 uint8
	bgpColour3 uint8

	// OBP0
	obp0Colour1 uint8
	obp0Colour2 uint8
	obp0Colour3 uint8

	// OBP1
	obp1Colour1 uint8
	obp1Colour2 uint8
	obp1Colour3 uint8

	// Single-value registers
	ly  uint8
	lyc uint8
	scx uint8
	scy uint8
	wx  uint8
	wy  uint8

	// Internal state
	interrupts *interrupts.Interrupts
	oam        *oam.OAM
	videoRAM   [0x2000]byte
	tick       int
	debug      bool

	// Internal LCD state
	tileCache      [384]*[8][8]uint8
	previousBg     [32][32]uint16
	previousWindow [32][32]uint16
	bg             [256][256]uint8
	window         [256][256]uint8
	sprites        [144][160]uint8
	frame          *image.RGBA
}

func New(interrupts *interrupts.Interrupts, oam *oam.OAM, debug bool) *PPU {
	ppu := &PPU{
		interrupts: interrupts,
		oam:        oam,
		frame:      image.NewRGBA(image.Rect(0, 0, 256, 256)),
		debug:      debug,
	}
	ppu.WriteLCDC(0x91)
	ppu.WriteLY(0x00)
	ppu.WriteLYC(0x00)
	ppu.WriteSCX(0x00)
	ppu.WriteSCY(0x00)
	ppu.WriteSTAT(0x00)
	ppu.WriteWX(0x00)
	ppu.WriteWY(0x00)
	ppu.WriteBGP(0xFC)
	ppu.WriteOBP0(0xFF)
	ppu.WriteOBP1(0xFF)
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
