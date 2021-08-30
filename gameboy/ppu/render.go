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
	// if ppu.spritesEnabled {
	// 	for sprite := range ppu.spriteOverlaps {
	// 		spriteAddr := 0xfe00 + uint16(sprite*4) + 1
	// 		spriteX := ppu.oam.ReadOAM(spriteAddr)
	// 		if x >= spriteX && x < spriteX+8 {
	// 			spriteY := ppu.oam.ReadOAM(spriteAddr + 2)
	// 			tileNumber := ppu.oam.ReadOAM(spriteAddr + 2)
	// 			attributes := ppu.oam.ReadOAM(spriteAddr + 3)
	// 			ppu.drawSpritePixel(x, y, spriteX, spriteY, tileNumber, attributes)
	// 			return
	// 		}
	// 	}
	// }

	// Use the background
	if ppu.bgEnabled {
		ppu.drawBackgroundPixel(x, y)
	}

}

func (ppu *PPU) drawSpritePixel(x, y, spriteX, spriteY, tileNumber, attributes uint8) {

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
		offsetAddr = 0x9c00
	} else {
		offsetAddr = 0x9800
	}

	var tileNumber uint16
	tileAddr := 32*uint16(tileY) + uint16(tileX)
	tileByte := ppu.ReadVideoRAM(offsetAddr + tileAddr)
	if ppu.lowTileData {
		tileNumber = uint16(tileByte)
	} else {
		tileNumber = uint16(256 + int(int8(tileByte)))
	}

	startAddr := uint16(0x8000 + (tileNumber * 16))
	a := ppu.ReadVideoRAM(startAddr + uint16(tileOffsetY*2))
	b := ppu.ReadVideoRAM(startAddr + uint16(tileOffsetY*2) + 1)
	aset := a&patterns[7-tileOffsetX] > 0
	bset := b&patterns[7-tileOffsetX] > 0

	var pixel uint8
	switch {
	case !aset && !bset:
		pixel = ppu.bgpColour[0]
	case !aset && bset:
		pixel = ppu.bgpColour[1]
	case aset && !bset:
		pixel = ppu.bgpColour[2]
	case aset && bset:
		pixel = ppu.bgpColour[3]
	}

	ppu.frame.SetRGBA(int(x), int(y), grey[pixel])

}
