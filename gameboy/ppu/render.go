package ppu

import (
	"image/color"
)

var (
	patterns = []uint8{
		0b00000001,
		0b00000010,
		0b00000100,
		0b00001000,
		0b00010000,
		0b00100000,
		0b01000000,
		0b10000000,
	}

	grey = []color.RGBA{
		{0xff, 0xff, 0xff, 0xff},
		{0xaa, 0xaa, 0xaa, 0xff},
		{0x77, 0x77, 0x77, 0xff},
		{0x33, 0x33, 0x33, 0xff},
	}

	red = []color.RGBA{
		{0xff, 0xaa, 0xaa, 0xff},
		{0xdd, 0x77, 0x77, 0xff},
		{0xaa, 0x33, 0x33, 0xff},
		{0x55, 0x00, 0x00, 0xff},
	}

	green = []color.RGBA{
		{0xaa, 0xff, 0xaa, 0xff},
		{0x77, 0xdd, 0x77, 0xff},
		{0x33, 0xaa, 0x33, 0xff},
		{0x00, 0x55, 0x00, 0xff},
	}

	blue = []color.RGBA{
		{0xaa, 0xaa, 0xff, 0xff},
		{0x77, 0x77, 0xdd, 0xff},
		{0x33, 0x33, 0xaa, 0xff},
		{0x00, 0x00, 0x55, 0xff},
	}
)

func (ppu *PPU) drawPixel(x, y uint8) {

	// Find the first overlapping spite if there is one
	if ppu.spritesEnabled {
		for sprite, overlaps := range ppu.spriteOverlaps {
			if !overlaps {
				continue
			}
			spriteAddr := 0xfe00 + uint16(sprite*4)
			spriteX := ppu.oam.ReadOAM(spriteAddr + 1)
			if x >= spriteX-8 && x < spriteX {
				spriteY := ppu.oam.ReadOAM(spriteAddr)
				tileNumber := ppu.oam.ReadOAM(spriteAddr + 2)
				attributes := ppu.oam.ReadOAM(spriteAddr + 3)
				ppu.drawSpritePixel(x, y, spriteX, spriteY, int(tileNumber), attributes)
				return
			}
		}
	}

	// Use the background
	if ppu.bgEnabled {
		ppu.drawBackgroundPixel(x, y)
	}

}

func (ppu *PPU) drawSpritePixel(x, y, spriteX, spriteY uint8, tileNumber int, attributes uint8) {
	tileOffsetX := (x - spriteX) % 8
	tileOffsetY := (y - spriteY) % 8
	pixel := ppu.readTilePixel(tileNumber, tileOffsetX, tileOffsetY)
	ppu.frame.SetRGBA(int(x), int(y), grey[ppu.bgpColour[pixel]])
}

func (ppu *PPU) drawBackgroundPixel(x, y uint8) {

	scx := ppu.ReadSCX()
	scy := ppu.ReadSCY()
	tileX := (x + scx) / 8
	tileY := (y + scy) / 8
	tileOffsetX := (x + scx) % 8
	tileOffsetY := (y + scy) % 8

	var offsetAddr uint16
	if ppu.highBgTileMap {
		offsetAddr = 0x9c00 - 0x8000
	} else {
		offsetAddr = 0x9800 - 0x8000
	}

	var tileNumber int
	tileAddr := 32*uint16(tileY) + uint16(tileX)
	tileByte := ppu.videoRAM[offsetAddr+tileAddr]
	if ppu.lowTileData {
		tileNumber = int(tileByte)
	} else {
		tileNumber = 256 + int(int8(tileByte))
	}

	pixel := ppu.readTilePixel(tileNumber, tileOffsetX, tileOffsetY)

	ppu.frame.SetRGBA(int(x), int(y), grey[ppu.bgpColour[pixel]])

}

func (ppu *PPU) readTilePixel(tileNumber int, tileOffsetX, tileOffsetY uint8) uint8 {
	startAddr := tileNumber * 16
	a := ppu.videoRAM[startAddr+int(tileOffsetY*2)]
	b := ppu.videoRAM[startAddr+int(tileOffsetY*2)+1]
	aset := a&patterns[7-tileOffsetX] > 0
	bset := b&patterns[7-tileOffsetX] > 0
	switch {
	case !aset && !bset:
		return 0
	case !aset && bset:
		return 1
	case aset && !bset:
		return 2
	case aset && bset:
		return 3
	default:
		panic("error reading tile")
	}
}
