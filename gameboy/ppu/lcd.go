package ppu

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

const (
	bit7 = 1 << iota
	bit6 = 1 << iota
	bit5 = 1 << iota
	bit4 = 1 << iota
	bit3 = 1 << iota
	bit2 = 1 << iota
	bit1 = 1 << iota
	bit0 = 1 << iota
)

var (
	gray = []color.RGBA{
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

// EndMachineCycle updates the LCD driver after each machine cycle i.e. 4 clock cycles
func (ppu *PPU) EndMachineCycle() {

	// is the lcd enabled?
	if ppu.lcdc&0x80 == 0 {
		ppu.ly = 0
		ppu.tick = 0
		return
	}

	// where are we on the lcd?
	ppu.ly = uint8(ppu.tick / 114)
	x := ppu.tick % 114
	ppu.tick++
	if ppu.tick >= 17556 {
		ppu.tick = 0
	}

	// set mode on stat register
	switch {
	case x == 0 && ppu.ly == 144:
		// v-blank period starts
		ppu.stat = (ppu.stat & 0xfc) | 0x01

		// v-blank interrupt
		ppu.interrupts.IF |= 0x01
		// is lcd stat interrupt enabled?
		if ppu.stat&0x10 > 0 {
			ppu.interrupts.IF |= 0x02
		}

	case x == 0 && ppu.ly < 144:
		// oam period starts
		ppu.stat = (ppu.stat & 0xfc) | 0x02

		// is lcd stat interrupt enabled?
		if ppu.stat&0x20 > 0 {
			ppu.interrupts.IF |= 0x02
		}

	case x == 20 && ppu.ly < 144:
		// lcd data transfer period starts
		ppu.stat = (ppu.stat & 0xfc) | 0x03
	case x == 63 && ppu.ly < 144:
		// h-blank period starts
		ppu.stat = (ppu.stat & 0xfc)

		// is lcd stat interrupt enabled?
		if ppu.stat&0x08 > 0 {
			ppu.interrupts.IF |= 0x02
		}

		// render lcd line
		ppu.updateLcdLine(ppu.ly)
	}

	// check coincidence flag
	if x == 0 {
		if ppu.ly == ppu.lyc {
			ppu.stat |= 0x04

			// is lcd stat interrupt enabled?
			if ppu.stat&0x40 > 0 {
				ppu.interrupts.IF |= 0x02
			}

		} else {
			ppu.stat &^= 0x04
		}
	}
}

func (ppu *PPU) readTile(tileNumber uint16) (*[8][8]uint8, bool) {
	tile := ppu.tileCache[tileNumber]
	if tile != nil {
		return tile, true
	}
	tile = &[8][8]uint8{}
	startAddr := uint16(0x8000 + (tileNumber * 16))
	for y := uint16(0); y < 8; y++ {
		a := ppu.ReadVideoRAM(startAddr + y*2)
		b := ppu.ReadVideoRAM(startAddr + y*2 + 1)
		for x := uint16(0); x < 8; x++ {
			tile[y][0] = (a&bit0)>>7 | (b&bit0)>>6
			tile[y][1] = (a&bit1)>>6 | (b&bit1)>>5
			tile[y][2] = (a&bit2)>>5 | (b&bit2)>>4
			tile[y][3] = (a&bit3)>>4 | (b&bit3)>>3
			tile[y][4] = (a&bit4)>>3 | (b&bit4)>>2
			tile[y][5] = (a&bit5)>>2 | (b&bit5)>>1
			tile[y][6] = (a&bit6)>>1 | (b & bit6)
			tile[y][7] = (a & bit7) | (b&bit7)<<1
		}
	}
	ppu.tileCache[tileNumber] = tile
	return tile, false
}

func (ppu *PPU) updateTiles(lcdY uint8, offsetAddr uint16, layer *[256][256]uint8, previousTiles *[32][32]uint16) {
	lowTileData := ppu.lowTileDataSelect()
	tileY := lcdY / 8
	for tileX := 0; tileX < 32; tileX++ {
		var tileNumber uint16
		tileAddr := 32*uint16(tileY) + uint16(tileX)
		tileByte := ppu.ReadVideoRAM(offsetAddr + tileAddr)
		if !lowTileData {
			tileNumber = uint16(256 + int(int8(tileByte)))
		} else {
			tileNumber = uint16(tileByte)
		}
		tile, cacheHit := ppu.readTile(tileNumber)
		if !cacheHit || tileNumber != previousTiles[tileY][tileX] {
			lcdX := uint8(tileX * 8)
			for y := uint8(0); y < 8; y++ {
				for x := uint8(0); x < 8; x++ {
					layer[lcdY+y][lcdX+x] = tile[y][x]
				}
			}
		}
		previousTiles[tileY][tileX] = tileNumber
	}
}

func (ppu *PPU) updateBG(lcdY, scy uint8) {
	if !ppu.bgDisplayEnable() {
		return
	}
	var offsetAddr uint16
	if ppu.highBgTileMapDisplaySelect() {
		offsetAddr = 0x9c00
	} else {
		offsetAddr = 0x9800
	}
	ppu.updateTiles(lcdY+scy, offsetAddr, &ppu.bg, &ppu.previousBg)
}

func (ppu *PPU) updateWindow(lcdY uint8) {
	if !ppu.windowDisplayEnable() {
		return
	}
	var offsetAddr uint16
	if ppu.highWindowTileMapDisplaySelect() {
		offsetAddr = 0x9c00
	} else {
		offsetAddr = 0x9800
	}
	ppu.updateTiles(lcdY, offsetAddr, &ppu.window, &ppu.previousWindow)
}

func (ppu *PPU) readSpriteTile(tileNumber uint16, att uint8) *[8][8]uint8 {
	tile, _ := ppu.readTile(uint16(tileNumber))
	if spriteXFlip(att) {
		new := [8][8]uint8{}
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				new[y][x] = tile[y][7-x]
			}
		}
		tile = &new
	}
	if spriteYFlip(att) {
		new := [8][8]uint8{}
		for y := 0; y < 8; y++ {
			new[y] = tile[7-y]
		}
		tile = &new
	}
	return tile
}

func (ppu *PPU) updateSprites(lcdY uint8) {
	if lcdY >= 144 {
		return
	}
	if ppu.largeSprites() {
		panic(fmt.Sprintf("Large sprites are not supported"))
	}
	ppu.sprites[lcdY] = [160]uint8{}
	for sprite := 0; sprite < 40; sprite++ {
		spriteAddr := uint16(sprite * 4)
		startY := ppu.oam[spriteAddr]
		if startY == 0 || lcdY < startY-16 || lcdY >= startY-8 {
			continue
		}
		startX := ppu.oam[spriteAddr+1]
		if startX == 0 || startX >= 168 {
			continue
		}
		tileNumber := ppu.oam[spriteAddr+2]
		attributes := ppu.oam[spriteAddr+3]
		tile := ppu.readSpriteTile(uint16(tileNumber), attributes)
		spriteX := startX - 8
		for tileX := uint8(0); tileX < 8; tileX++ {
			lcdX := spriteX + tileX
			if lcdX < 160 {
				ppu.sprites[lcdY][lcdX] = tile[lcdY-startY+16][tileX]
			}
		}
	}
}

func (ppu *PPU) renderPixel(x, y, scx, scy, wx, wy uint8, debug bool) color.RGBA {
	if ppu.spriteDisplayEnable() {
		if x < 160 && y < 144 {
			pixel := ppu.sprites[y][x]
			if pixel > 0 {
				if debug {
					return blue[pixel]
				}
				return gray[pixel]
			}
		}
	}

	// Make tiles visible
	// if (x+scx)%8 == 0 && (y+scy)%2 == 0 ||
	// 	(x+scx)%2 == 0 && (y+scy)%8 == 0 {
	// 	return color.RGBA{0xff, 0, 0, 0xff}
	// }

	if ppu.windowDisplayEnable() {
		// Use WX/WY to shift the visible pixels
		if x >= wx && y >= wy {
			pixel := ppu.window[y-wy][x-wx]
			if debug {
				return green[pixel]
			}
			return gray[pixel]
		}
	}
	if ppu.bgDisplayEnable() {
		// Use SCX/SCY to shift the visible pixels
		pixel := ppu.bg[y+scy][x+scx]
		if debug && (x >= 160 || y >= 144) {
			return red[pixel]
		}

		return gray[pixel]
	}
	return gray[0]
}

func (ppu *PPU) renderLine(y, scy uint8) {
	scx := ppu.ReadSCX()
	wx := ppu.ReadWX()
	wy := ppu.ReadWY()
	if ppu.debug {
		for x := 0; x < 256; x++ {
			pixel := ppu.renderPixel(uint8(x)-scx, y-scy, scx, scy, wx, wy, true)
			ppu.frame.SetRGBA(x, int(y), pixel)
		}
	} else {
		for x := 0; x < 160; x++ {
			pixel := ppu.renderPixel(uint8(x), y, scx, scy, wx, wy, false)
			ppu.frame.SetRGBA(x, int(y), pixel)
		}
	}
}

func (ppu *PPU) updateLcdLine(y uint8) {
	scy := ppu.ReadSCY()
	ppu.updateBG(y, scy)
	ppu.updateWindow(y)
	ppu.updateSprites(y)
	ppu.renderLine(y, scy)
}

// FrameEnd writes any remaining VRAM lines to the GUI for debugging
func (ppu *PPU) Frame() *image.RGBA {
	// if ppu.debug {
	// 	for y := ppu.memory.SCY + 144; y < ppu.memory.SCY || y >= ppu.memory.SCY+144; y++ {
	// 		ppu.updateLcdLine(y)
	// 	}
	// }
	return ppu.frame
}

// Screenshot writes a screenshot to file
func (ppu *PPU) Screenshot(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	err = png.Encode(f, ppu.frame)
	if err != nil {
		fmt.Println(err)
	}
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 7 - LCD Display Enable             (0=Off, 1=On)
func (ppu *PPU) lcdDisplayEnable() bool {
	return ppu.lcdc&0x80 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 6 - Window Tile Map Display Select (0=9800-9BFF, 1=9C00-9FFF)
func (ppu *PPU) highWindowTileMapDisplaySelect() bool {
	return ppu.lcdc&0x40 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 5 - Window Display Enable          (0=Off, 1=On)
func (ppu *PPU) windowDisplayEnable() bool {
	return ppu.lcdc&0x20 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 4 - BG & Window Tile Data Select   (0=8800-97FF, 1=8000-8FFF)
func (ppu *PPU) lowTileDataSelect() bool {
	return ppu.lcdc&0x10 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 3 - BG Tile Map Display Select     (0=9800-9BFF, 1=9C00-9FFF)
func (ppu *PPU) highBgTileMapDisplaySelect() bool {
	return ppu.lcdc&0x08 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 2 - OBJ (Sprite) Size              (0=8x8, 1=8x16)
func (ppu *PPU) largeSprites() bool {
	return ppu.lcdc&0x04 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 1 - OBJ (Sprite) Display Enable    (0=Off, 1=On)
func (ppu *PPU) spriteDisplayEnable() bool {
	return ppu.lcdc&0x02 > 0
}

// FF40 - LCDC - LCD Control (R/W)
// Bit 0 - BG Display (for CGB see below) (0=Off, 1=On)
func (ppu *PPU) bgDisplayEnable() bool {
	return ppu.lcdc&0x01 > 0
}

// Byte3 - Attributes/Flags:
//   Bit7   OBJ-to-BG Priority (0=OBJ Above BG, 1=OBJ Behind BG color 1-3)
//          (Used for both BG and Window. BG color 0 is always behind OBJ)
//   Bit6   Y flip          (0=Normal, 1=Vertically mirrored)
//   Bit5   X flip          (0=Normal, 1=Horizontally mirrored)
//   Bit4   Palette number  **Non CGB Mode Only** (0=OBP0, 1=OBP1)
//   Bit3   Tile VRAM-Bank  **CGB Mode Only**     (0=Bank 0, 1=Bank 1)
//   Bit2-0 Palette number  **CGB Mode Only**     (OBP0-7)

func spriteBehindBg(att uint8) bool {
	return att&0x80 > 0
}

func spriteYFlip(att uint8) bool {
	return att&0x40 > 0
}

func spriteXFlip(att uint8) bool {
	return att&0x20 > 0
}

func spritePalatteSet(att uint8) bool {
	return att&0x08 > 0
}
