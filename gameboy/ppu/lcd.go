package ppu

// func (ppu *PPU) readTile(tileNumber uint16) (*[8][8]uint8, bool) {
// 	tile := ppu.tileCache[tileNumber]
// 	if tile != nil {
// 		return tile, true
// 	}
// 	tile = &[8][8]uint8{}
// 	startAddr := uint16(0x8000 + (tileNumber * 16))
// 	for y := uint16(0); y < 8; y++ {
// 		a := ppu.ReadVideoRAM(startAddr + y*2)
// 		b := ppu.ReadVideoRAM(startAddr + y*2 + 1)
// 		for x := uint16(0); x < 8; x++ {
// 			tile[y][0] = (a&bit0)>>7 | (b&bit0)>>6
// 			tile[y][1] = (a&bit1)>>6 | (b&bit1)>>5
// 			tile[y][2] = (a&bit2)>>5 | (b&bit2)>>4
// 			tile[y][3] = (a&bit3)>>4 | (b&bit3)>>3
// 			tile[y][4] = (a&bit4)>>3 | (b&bit4)>>2
// 			tile[y][5] = (a&bit5)>>2 | (b&bit5)>>1
// 			tile[y][6] = (a&bit6)>>1 | (b & bit6)
// 			tile[y][7] = (a & bit7) | (b&bit7)<<1
// 		}
// 	}
// 	ppu.tileCache[tileNumber] = tile
// 	return tile, false
// }

// func (ppu *PPU) updateTiles(lcdY uint8, offsetAddr uint16, layer *[256][256]uint8, previousTiles *[32][32]uint16) {

// }

// func (ppu *PPU) updateWindow(lcdY uint8) {
// 	if !ppu.windowEnabled {
// 		return
// 	}
// 	var offsetAddr uint16
// 	if ppu.highWindowTileMap {
// 		offsetAddr = 0x9c00
// 	} else {
// 		offsetAddr = 0x9800
// 	}
// 	ppu.updateTiles(lcdY, offsetAddr, &ppu.window, &ppu.previousWindow)
// }

// func (ppu *PPU) readSpriteTile(tileNumber uint16, att uint8) *[8][8]uint8 {
// 	tile, _ := ppu.readTile(uint16(tileNumber))
// 	if spriteXFlip(att) {
// 		new := [8][8]uint8{}
// 		for y := 0; y < 8; y++ {
// 			for x := 0; x < 8; x++ {
// 				new[y][x] = tile[y][7-x]
// 			}
// 		}
// 		tile = &new
// 	}
// 	if spriteYFlip(att) {
// 		new := [8][8]uint8{}
// 		for y := 0; y < 8; y++ {
// 			new[y] = tile[7-y]
// 		}
// 		tile = &new
// 	}
// 	return tile
// }

// func (ppu *PPU) updateSprites(lcdY uint8) {
// 	if lcdY >= 144 {
// 		return
// 	}
// 	if ppu.spritesLarge {
// 		panic(fmt.Sprintf("Large sprites are not supported"))
// 	}
// 	ppu.sprites[lcdY] = [160]uint8{}
// 	for sprite := 0; sprite < 40; sprite++ {
// 		spriteAddr := 0xfe00 + uint16(sprite*4)
// 		startY := ppu.oam.ReadOAM(spriteAddr)
// 		if startY == 0 || lcdY < startY-16 || lcdY >= startY-8 {
// 			continue
// 		}
// 		startX := ppu.oam.ReadOAM(spriteAddr + 1)
// 		if startX == 0 || startX >= 168 {
// 			continue
// 		}
// 		tileNumber := ppu.oam.ReadOAM(spriteAddr + 2)
// 		attributes := ppu.oam.ReadOAM(spriteAddr + 3)
// 		tile := ppu.readSpriteTile(uint16(tileNumber), attributes)
// 		spriteX := startX - 8
// 		for tileX := uint8(0); tileX < 8; tileX++ {
// 			lcdX := spriteX + tileX
// 			if lcdX < 160 {
// 				ppu.sprites[lcdY][lcdX] = tile[lcdY-startY+16][tileX]
// 			}
// 		}
// 	}
// }

// func (ppu *PPU) renderPixel(x, y, scx, scy, wx, wy uint8, debug bool) color.RGBA {
// 	if ppu.spritesEnabled {
// 		if x < 160 && y < 144 {
// 			pixel := ppu.sprites[y][x]
// 			if pixel > 0 {
// 				if debug {
// 					return blue[pixel]
// 				}
// 				return gray[pixel]
// 			}
// 		}
// 	}

// 	// Make tiles visible
// 	// if (x+scx)%8 == 0 && (y+scy)%2 == 0 ||
// 	// 	(x+scx)%2 == 0 && (y+scy)%8 == 0 {
// 	// 	return color.RGBA{0xff, 0, 0, 0xff}
// 	// }

// 	if ppu.windowEnabled {
// 		// Use WX/WY to shift the visible pixels
// 		if x >= wx && y >= wy {
// 			pixel := ppu.window[y-wy][x-wx]
// 			if debug {
// 				return green[pixel]
// 			}
// 			return gray[pixel]
// 		}
// 	}
// 	if ppu.bgEnabled {
// 		// Use SCX/SCY to shift the visible pixels
// 		pixel := ppu.bg[y+scy][x+scx]
// 		if debug && (x >= 160 || y >= 144) {
// 			return red[pixel]
// 		}

// 		return gray[pixel]
// 	}
// 	return gray[0]
// }

// func (ppu *PPU) renderLine(y, scy uint8) {
// 	scx := ppu.ReadSCX()
// 	wx := ppu.ReadWX()
// 	wy := ppu.ReadWY()
// 	if ppu.debug {
// 		for x := 0; x < 256; x++ {
// 			pixel := ppu.renderPixel(uint8(x)-scx, y-scy, scx, scy, wx, wy, true)
// 			ppu.frame.SetRGBA(x, int(y), pixel)
// 		}
// 	} else {
// 		for x := 0; x < 160; x++ {
// 			pixel := ppu.renderPixel(uint8(x), y, scx, scy, wx, wy, false)
// 			ppu.frame.SetRGBA(x, int(y), pixel)
// 		}
// 	}
// }

// func (ppu *PPU) updateLcdLine(y uint8) {
// 	scy := ppu.ReadSCY()
// 	ppu.updateBG(y, scy)
// 	ppu.updateWindow(y)
// 	ppu.updateSprites(y)
// 	ppu.renderLine(y, scy)
// }

// // Screenshot writes a screenshot to file
// func (ppu *PPU) Screenshot(filename string) {
// 	f, err := os.Create(filename)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer f.Close()
// 	err = png.Encode(f, ppu.frame)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

// // Byte3 - Attributes/Flags:
// //   Bit7   OBJ-to-BG Priority (0=OBJ Above BG, 1=OBJ Behind BG color 1-3)
// //          (Used for both BG and Window. BG color 0 is always behind OBJ)
// //   Bit6   Y flip          (0=Normal, 1=Vertically mirrored)
// //   Bit5   X flip          (0=Normal, 1=Horizontally mirrored)
// //   Bit4   Palette number  **Non CGB Mode Only** (0=OBP0, 1=OBP1)
// //   Bit3   Tile VRAM-Bank  **CGB Mode Only**     (0=Bank 0, 1=Bank 1)
// //   Bit2-0 Palette number  **CGB Mode Only**     (OBP0-7)

// func spriteBehindBg(att uint8) bool {
// 	return att&0x80 > 0
// }

// func spriteYFlip(att uint8) bool {
// 	return att&0x40 > 0
// }

// func spriteXFlip(att uint8) bool {
// 	return att&0x20 > 0
// }

// func spritePalatteSet(att uint8) bool {
// 	return att&0x08 > 0
// }
