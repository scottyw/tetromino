package ppu

import (
	"image/color"
)

var (

	// Tile bits are counted left-to-right so bit 0 of the tile is bit 7 of the byte
	patterns = []uint8{
		0b10000000,
		0b01000000,
		0b00100000,
		0b00010000,
		0b00001000,
		0b00000100,
		0b00000010,
		0b00000001,
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

func (ppu *PPU) renderPixel(x, y uint8) {

	// Find the first overlapping spite if there is one
	if ppu.spritesEnabled {
		for sprite, overlaps := range ppu.spriteOverlaps {
			if !overlaps {
				continue
			}
			spriteAddr := 0xfe00 + uint16(sprite*4)
			spriteX := ppu.oam.ReadOAM(spriteAddr + 1)
			if x+8 >= spriteX && x < spriteX {
				spriteY := ppu.oam.ReadOAM(spriteAddr)
				tileNumber := ppu.oam.ReadOAM(spriteAddr + 2)
				attributes := ppu.oam.ReadOAM(spriteAddr + 3)
				if ppu.renderSpritePixel(x, y, spriteX, spriteY, int(tileNumber), attributes) {
					return
				}
			}
		}
	}

	// Use the background
	if ppu.bgEnabled {
		ppu.renderBackgroundPixel(x, y)
	}

}

func (ppu *PPU) renderSpritePixel(x, y, spriteX, spriteY uint8, tileNumber int, attributes uint8) bool {
	tileOffsetX := (x - spriteX) % 8
	tileOffsetY := (y - spriteY) % 8
	behind := attributes&0x80 > 0
	flipY := attributes&0x40 > 0
	flipX := attributes&0x20 > 0
	usePalette1 := attributes&0x08 > 0
	if flipX {
		tileOffsetX = 7 - tileOffsetX
	}
	if flipY {
		tileOffsetY = 7 - tileOffsetY
	}
	if behind {
		panic("behind")
	}
	pixel := ppu.readTilePixel(tileNumber, tileOffsetX, tileOffsetY)
	if pixel > 0 {
		var colour color.RGBA
		if usePalette1 {
			colour = blue[ppu.obp1Colour[pixel]]
		} else {
			colour = blue[ppu.obp0Colour[pixel]]
		}
		ppu.frame.SetRGBA(int(x), int(y), colour)
		return true
	}
	return false
}

func (ppu *PPU) renderBackgroundPixel(x, y uint8) {

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
	aset := a&patterns[tileOffsetX] > 0
	bset := b&patterns[tileOffsetX] > 0
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
