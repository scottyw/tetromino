package lcd

import (
	"fmt"

	"github.com/scottyw/tetromino/pkg/gb/mem"
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

// LCD represents the LCD display of the Gameboy
type LCD struct {
	hwr      *mem.HardwareRegisters
	videoRAM *[0x2000]byte
	oam      *[0xa0]byte
	data     [256 * 144]uint8
	debug    bool
}

// NewLCD returns the configured LCD
func NewLCD(hwr *mem.HardwareRegisters, memory *mem.Memory, debug bool) *LCD {
	return &LCD{
		hwr:      hwr,
		videoRAM: memory.VideoRAM(),
		oam:      memory.OAM(),
		debug:    debug,
	}
}

// Tick runs the LCD driver for one machine cycle i.e. 4 clock cycles
func (lcd *LCD) Tick(cycle int, foo bool) {
	lcd.hwr.LY = uint8(cycle / 114)
	x := cycle % 114

	// Set mode on STAT register
	switch {
	case x == 0 && lcd.hwr.LY == 144:
		// V-Blank period starts
		lcd.hwr.STAT = (lcd.hwr.STAT & 0xfc) | 0x01
		// V-Blank interrupt
		lcd.hwr.IF |= 0x01
		// Is LCD STAT interrupt enabled?
		if lcd.hwr.STAT&0x10 > 0 {
			lcd.hwr.IF |= 0x02
		}
	case x == 0 && lcd.hwr.LY < 144:
		// OAM period starts
		lcd.hwr.STAT = (lcd.hwr.STAT & 0xfc) | 0x02
		// Is LCD STAT interrupt enabled?
		if lcd.hwr.STAT&0x20 > 0 {
			lcd.hwr.IF |= 0x02
		}
	case x == 20 && lcd.hwr.LY < 144:
		// LCD data transfer period starts
		lcd.hwr.STAT = (lcd.hwr.STAT & 0xfc) | 0x03
	case x == 63 && lcd.hwr.LY < 144:
		// H-Blank period starts
		lcd.hwr.STAT = (lcd.hwr.STAT & 0xfc)
		// Is LCD STAT interrupt enabled?
		if lcd.hwr.STAT&0x08 > 0 {
			lcd.hwr.IF |= 0x02
		}
		// Render LCD line
		lcd.updateLcdLine()
	}

	// Check coincidence flag
	if x == 0 {
		if lcd.hwr.LY == lcd.hwr.LYC {
			lcd.hwr.STAT |= 0x04
			// Is LCD STAT interrupt enabled?
			if lcd.hwr.STAT&0x40 > 0 {
				lcd.hwr.IF |= 0x02
			}
		} else {
			lcd.hwr.STAT &^= 0x04
		}
	}
}

// FrameData returns the frame data as a 256x144 array of bytes where each element is a colour value between 0 and 3
func (lcd *LCD) FrameData() [256 * 144]uint8 {
	return lcd.data
}

func lowTileAbsoluteAddress(tileNumber uint8) uint16 {
	return 0x8000 + uint16(tileNumber)*16
}

func highTileAbsoluteAddress(tileNumber int8) uint16 {
	return uint16(0x9000 + int(tileNumber)*16)
}

func (lcd *LCD) readVideoRAM(memoryAddr uint16) uint8 {
	return uint8(lcd.videoRAM[memoryAddr&0x7fff])
}

func (lcd *LCD) pixel(memoryAddr uint16, bit uint8) uint8 {
	var a, b, pixel uint8
	a = uint8(lcd.readVideoRAM(memoryAddr))
	b = uint8(lcd.readVideoRAM(memoryAddr + 1))
	switch bit {
	case 0:
		pixel = (a&bit0)>>7 | (b&bit0)>>6
	case 1:
		pixel = (a&bit1)>>6 | (b&bit1)>>5
	case 2:
		pixel = (a&bit2)>>5 | (b&bit2)>>4
	case 3:
		pixel = (a&bit3)>>4 | (b&bit3)>>3
	case 4:
		pixel = (a&bit4)>>3 | (b&bit4)>>2
	case 5:
		pixel = (a&bit5)>>2 | (b&bit5)>>1
	case 6:
		pixel = (a&bit6)>>1 | (b & bit6)
	case 7:
		pixel = (a & bit7) | (b&bit7)<<1
	default:
		panic(fmt.Sprintf("Bad bit in pixel(): %v", bit))
	}
	return pixel
}

// Returns the memory address of the tile
func (lcd *LCD) tileDataAddr(highTileMap, lowTileData bool, tileX, tileY uint8) uint16 {
	var tileNumberAddr, tileIndex uint16
	tileIndex = uint16(tileY)*32 + uint16(tileX)
	if highTileMap {
		tileNumberAddr = 0x9c00 + tileIndex
	} else {
		tileNumberAddr = 0x9800 + tileIndex
	}
	tileNumber := lcd.readVideoRAM(tileNumberAddr)
	if lowTileData {
		return lowTileAbsoluteAddress(tileNumber)
	}
	return highTileAbsoluteAddress(int8(tileNumber))
}

func (lcd *LCD) findSprites(lcdY uint8) []uint8 {
	var spriteAddrs []uint8
	for spriteAddr := uint8(0x00); spriteAddr < 0x9f; spriteAddr += 4 {
		startY := lcd.oam[spriteAddr]
		if startY == 0 || startY > 160 {
			continue
		}
		startX := lcd.oam[spriteAddr+1]
		if startX == 0 || startX > 168 {
			continue
		}
		if lcdY >= startY-16 &&
			lcdY < startY-8 {
			spriteAddrs = append(spriteAddrs, spriteAddr)
		}
	}
	return spriteAddrs
}

func (lcd *LCD) deriveSpritePixel(lcdX, lcdY, scx, scy uint8, spriteAddrs []uint8) (uint8, bool) {
	if lcd.debug {
		lcdX -= scx // Overflows deliberately
		lcdY -= scy // Overflows deliberately
	}
	// FIXME search for sprites on the current Y line outside this function
	largeSpriteSize := largeSpriteSize(lcd.hwr)
	if largeSpriteSize {
		panic(fmt.Sprintf("Large sprites are not supported"))
	}
	for _, spriteAddr := range spriteAddrs {
		startY := lcd.oam[spriteAddr]
		startX := lcd.oam[spriteAddr+1]
		tileNumber := lcd.oam[spriteAddr+2]
		// attributes := lcd.oam[spriteAddr+ 3]
		if lcdX >= startX-8 &&
			lcdX < startX {
			var tileOffsetX, tileOffsetY uint8
			var tileAddr, memoryAddr uint16
			tileAddr = lowTileAbsoluteAddress(tileNumber)
			tileOffsetX = lcdX - startX + 8
			tileOffsetY = lcdY - startY + 16
			memoryAddr = tileAddr + uint16(tileOffsetY)*2
			pixel := lcd.pixel(memoryAddr, tileOffsetX)
			if pixel > 0 {
				return pixel, true
			}
		}
	}
	return 0, false
}

func (lcd *LCD) deriveTilePixel(highTileMap, lowTileData bool, lcdX, lcdY uint8) uint8 {
	var tileX, tileY, tileOffsetX, tileOffsetY uint8
	var tileAddr, memoryAddr uint16
	tileX = lcdX / 8
	tileY = lcdY / 8
	tileAddr = lcd.tileDataAddr(highTileMap, lowTileData, tileX, tileY)
	tileOffsetX = lcdX % 8
	tileOffsetY = lcdY % 8
	memoryAddr = tileAddr + uint16(tileOffsetY)*2
	return lcd.pixel(memoryAddr, tileOffsetX)
}

func (lcd *LCD) deriveWindowPixel(highTileMap, lowTileData bool, wx, wy, lcdX, lcdY uint8) (uint8, bool) {
	if windowDisplayEnable(lcd.hwr) &&
		wx >= 0 &&
		wx <= 166 &&
		wy >= 0 &&
		wy <= 143 &&
		lcdX >= wx &&
		lcdY >= wy+7 {
		return lcd.deriveTilePixel(highTileMap, lowTileData, lcdX, lcdY), true
	}
	return 0, false
}

func (lcd *LCD) deriveBackgroundPixel(highTileMap, lowTileData bool, scx, scy, lcdX, lcdY uint8) uint8 {
	if !bgDisplay(lcd.hwr) {
		return 0
	}
	var pixel uint8
	if lcd.debug {
		// Use SCX/SCY to colour pixels to illusrate the offset
		if lcdX-scx >= 160 || lcdY-scy >= 144 {
			pixel += 0x10
		}
	} else {
		// Use SCX/SCY to shift the visible pixels
		lcdX += scx // Overflows deliberately
		lcdY += scy // Overflows deliberately
	}
	pixel += lcd.deriveTilePixel(highTileMap, lowTileData, lcdX, lcdY)
	return pixel
}

func (lcd *LCD) derivePixel(highBgTileMap, highWindowTileMap, lowTileData bool, scx, scy, wx, wy, lcdX, lcdY uint8, spriteAddrs []uint8) uint8 {
	// if !lcdDisplayEnable(memory) {
	// 	return 0
	// }
	if pixel, found := lcd.deriveSpritePixel(lcdX, lcdY, scx, scy, spriteAddrs); found {
		return pixel + 0x30 // Colour offset
	}
	if pixel, found := lcd.deriveWindowPixel(highWindowTileMap, lowTileData, wx, wy, lcdX, lcdY); found {
		return pixel + 0x20 // Colour offset
	}
	return lcd.deriveBackgroundPixel(highBgTileMap, lowTileData, scx, scy, lcdX, lcdY)
}

func (lcd *LCD) updateLcdLine() {
	lcdY := lcd.hwr.LY
	scx := lcd.hwr.SCX
	scy := lcd.hwr.SCY
	wx := lcd.hwr.WX
	wy := lcd.hwr.WY
	highBgTileMap := highBgTileMapDisplaySelect(lcd.hwr)
	highWindowTileMap := highWindowTileMapDisplaySelect(lcd.hwr)
	lowTileData := lowTileDataSelect(lcd.hwr)
	spriteAddrs := lcd.findSprites(lcdY)
	var width int
	if lcd.debug {
		width = 256
	} else {
		width = 160
	}
	for lcdX := 0; lcdX < width; lcdX++ {
		index := uint16(lcdY)*256 + uint16(lcdX)
		lcd.data[index] = lcd.derivePixel(highBgTileMap, highWindowTileMap, lowTileData, scx, scy, wx, wy, uint8(lcdX), lcdY, spriteAddrs)
	}
}
