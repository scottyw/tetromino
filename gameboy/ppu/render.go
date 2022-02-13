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

	// red = []color.RGBA{
	// 	{0xff, 0xaa, 0xaa, 0xff},
	// 	{0xdd, 0x77, 0x77, 0xff},
	// 	{0xaa, 0x33, 0x33, 0xff},
	// 	{0x55, 0x00, 0x00, 0xff},
	// }

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

	colours := grey

	var spritePixel uint8
	var spriteBehindBackground bool
	var useSpritePalette1 bool

	// Does this pixel intersect with a sprite?
	if ppu.spritesEnabled {

		for sprite, overlaps := range ppu.spriteOverlaps {
			if !overlaps {
				continue
			}
			spriteAddr := 0xfe00 + uint16(sprite*4)
			spriteX := ppu.oam.PPURead(spriteAddr + 1)
			if x+8 >= spriteX && x < spriteX {
				spriteY := ppu.oam.PPURead(spriteAddr)
				tileNumber := ppu.oam.PPURead(spriteAddr + 2)
				attributes := ppu.oam.PPURead(spriteAddr + 3)
				tileOffsetX := (x - spriteX) % 8
				tileOffsetY := (y - spriteY) % 8
				spriteBehindBackground = attributes&0x80 > 0
				flipY := attributes&0x40 > 0
				flipX := attributes&0x20 > 0
				useSpritePalette1 = attributes&0x08 > 0
				if flipX {
					tileOffsetX = 7 - tileOffsetX
				}
				if flipY {
					tileOffsetY = 7 - tileOffsetY
				}
				spritePixel = ppu.readTilePixel(int(tileNumber), tileOffsetX, tileOffsetY)
				if spritePixel > 0 {
					break
				}
			}
		}

		// Did we intersect with a foreground sprite?
		if spritePixel > 0 && !spriteBehindBackground {
			if ppu.debug {
				colours = blue
			}
			if useSpritePalette1 {
				ppu.frame.SetRGBA(int(x), int(y), colours[ppu.obp1Colour[spritePixel]])
			} else {
				ppu.frame.SetRGBA(int(x), int(y), colours[ppu.obp0Colour[spritePixel]])
			}
			return
		}

	}

	var pixel uint8

	// Does this pixel intersect with the window?
	var intersectedWindow bool
	if ppu.windowEnabled {
		wx := ppu.ReadWX()
		wy := ppu.ReadWY()
		if wx <= 166 && wy <= 143 && x >= (wx-7) && y >= wy {
			pixel = ppu.findWindowPixel(x-(wx-7), y-wy)
			intersectedWindow = true
		}
	}

	// If we didn't find a window pixel then find where this pixel intersects the background
	if !intersectedWindow && ppu.bgEnabled {
		pixel = ppu.findBackgroundPixel(x, y)
	}

	// Does this pixel intersect with a background sprite?
	if pixel == 0 && spritePixel != 0 && spriteBehindBackground {
		if ppu.debug {
			colours = blue
		}
		if useSpritePalette1 {
			ppu.frame.SetRGBA(int(x), int(y), colours[ppu.obp1Colour[spritePixel]])
		} else {
			ppu.frame.SetRGBA(int(x), int(y), colours[ppu.obp0Colour[spritePixel]])
		}
	} else {
		if intersectedWindow {
			if ppu.debug {
				colours = green
			}
			ppu.frame.SetRGBA(int(x), int(y), colours[ppu.bgpColour[pixel]])
		} else {
			ppu.frame.SetRGBA(int(x), int(y), grey[ppu.bgpColour[pixel]])
		}
	}

}

func (ppu *PPU) findWindowPixel(x, y uint8) uint8 {

	tileX := x / 8
	tileY := y / 8
	tileOffsetX := x % 8
	tileOffsetY := y % 8

	var offsetAddr uint16
	if ppu.highWindowTileMap {
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

	return ppu.readTilePixel(tileNumber, tileOffsetX, tileOffsetY)

}

func (ppu *PPU) findBackgroundPixel(x, y uint8) uint8 {

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

	return ppu.readTilePixel(tileNumber, tileOffsetX, tileOffsetY)

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
