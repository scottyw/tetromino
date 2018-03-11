package lcd

import (
	"fmt"

	"github.com/scottyw/tetromino/mem"
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
	data [23040]uint8
}

// NewLCD returns the configured LCD
func NewLCD() *LCD {
	return &LCD{}
}

// Tick runs the LCD driver for one machine cycle i.e. 4 clock cycles
func (lcd *LCD) Tick(memory mem.Memory, cycle int) {
	ly := uint8(cycle / 114)
	lyRemainder := cycle % 114
	var stat uint8
	// Set mode on stat register
	switch {
	case ly == 144:
		// V-Blank period starts
		stat = 1
		*memory.IF |= 0x01
	case lyRemainder == 0:
		// OAM period starts
		stat = 2
	case lyRemainder == 20:
		// LCD data transfer period starts
		stat = 3
	case lyRemainder == 63:
		// H-Blank period starts
		stat = 0
		if ly < 144 {
			lcd.updateLcdLine(memory, ly)
		}
	default:
		stat = *memory.STAT
	}
	// Set coincidence flag and coincidence interrupt on stat register
	if ly == *memory.LYC {
		stat |= 0x44
	} else {
		stat &^= 0x44
	}
	// Set interrupts on stat register
	switch {
	case ly == 144:
		stat |= 0x10
	case lyRemainder == 0:
		stat |= 0x20
	case lyRemainder == 63:
		stat |= 0x08
	}
	*memory.STAT = stat
	*memory.LY = ly
}

// FrameData returns the frame data as a 160x144 array of bytes where each element is a colour value between 0 and 3
func (lcd *LCD) FrameData() [23040]uint8 {
	return lcd.data
}

func lowTileAbsoluteAddress(tileNumber uint8) uint16 {
	return 0x8000 + uint16(tileNumber)*16
}

func highTileAbsoluteAddress(tileNumber int8) uint16 {
	return uint16(0x9000 + int(tileNumber)*16)
}

// Returns 16 bytes representing one 8x8 tile
func tileData(mem mem.Memory, tile uint16) []byte {
	var tileAddr uint16
	if highBgTileMapDisplaySelect(mem) {
		tileAddr = 0x9c00 + tile
	} else {
		tileAddr = 0x9800 + tile
	}
	tileNumber := mem.Read(tileAddr)
	if lowTileDataSelect(mem) {
		return mem.ReadRegion(lowTileAbsoluteAddress(tileNumber), 16)
	}
	return mem.ReadRegion(highTileAbsoluteAddress(int8(tileNumber)), 16)
}

func pixel(mem mem.Memory, memoryAddr uint16, bit uint8) uint8 {
	var a, b, pixel uint8
	a = uint8(mem.Read(memoryAddr))
	b = uint8(mem.Read(memoryAddr + 1))
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
func tileDataAddr(mem mem.Memory, tileX, tileY uint8, highTileMap, lowTileData bool) uint16 {
	var tileNumberAddr, tileIndex uint16
	tileIndex = uint16(tileY)*32 + uint16(tileX)
	if highTileMap {
		tileNumberAddr = 0x9c00 + tileIndex
	} else {
		tileNumberAddr = 0x9800 + tileIndex
	}
	tileNumber := mem.Read(tileNumberAddr)
	if lowTileData {
		return lowTileAbsoluteAddress(tileNumber)
	}
	return highTileAbsoluteAddress(int8(tileNumber))
}

func deriveSpritePixel(mem mem.Memory, lcdX, lcdY uint8) (uint8, bool) {
	// FIXME search for sprites on the current Y line outside this function
	largeSpriteSize := largeSpriteSize(mem)
	var spriteAddr uint16
	for spriteAddr = 0xfe00; spriteAddr < 0xfe9f; spriteAddr += 4 {
		startY := mem.Read(spriteAddr)
		if startY == 0 || startY > 160 {
			continue
		}
		startX := mem.Read(spriteAddr + 1)
		if startX == 0 || startX > 168 {
			continue
		}
		tileNumber := mem.Read(spriteAddr + 2)
		// attributes := mem.Read(spriteAddr + 3)
		endY := startY + 8
		endX := startX + 8
		if largeSpriteSize {
			endY += 8
			panic(fmt.Sprintf("Large sprites are not supported"))
		}
		if lcdX >= startX &&
			lcdX < endX &&
			lcdY >= startY &&
			lcdY < endY {
			var tileOffsetX, tileOffsetY uint8
			var tileAddr, memoryAddr uint16
			tileAddr = lowTileAbsoluteAddress(tileNumber)
			tileOffsetX = lcdX - startX
			tileOffsetY = lcdY - startY
			memoryAddr = tileAddr + uint16(tileOffsetY)*2
			return pixel(mem, memoryAddr, tileOffsetX), true
		}
	}
	return 0, false
}

func deriveTilePixel(mem mem.Memory, lcdX, lcdY uint8, highTileMap, lowTileData bool) uint8 {
	var tileX, tileY, tileOffsetX, tileOffsetY uint8
	var tileAddr, memoryAddr uint16
	tileX = lcdX / 8
	tileY = lcdY / 8
	tileAddr = tileDataAddr(mem, tileX, tileY, highTileMap, lowTileData)
	tileOffsetX = lcdX % 8
	tileOffsetY = lcdY % 8
	memoryAddr = tileAddr + uint16(tileOffsetY)*2
	return pixel(mem, memoryAddr, tileOffsetX)
}

func deriveWindowPixel(mem mem.Memory, lcdX, lcdY uint8) (uint8, bool) {
	if windowDisplayEnable(mem) &&
		*mem.WX >= 0 &&
		*mem.WX <= 166 &&
		*mem.WY >= 0 &&
		*mem.WY <= 143 &&
		lcdX >= *mem.WX &&
		lcdY >= *mem.WY+7 {
		highTileMap := highWindowTileMapDisplaySelect(mem)
		lowTileData := lowTileDataSelect(mem)
		return deriveTilePixel(mem, lcdX, lcdY, highTileMap, lowTileData), true
	}
	return 0, false
}

func deriveBackgroundPixel(mem mem.Memory, lcdX, lcdY uint8) uint8 {
	if !bgDisplay(mem) {
		return 0
	}
	highTileMap := highBgTileMapDisplaySelect(mem)
	lowTileData := lowTileDataSelect(mem)
	lcdX += *mem.SCX // Overflows deliberately
	lcdY += *mem.SCY // Overflows deliberately
	return deriveTilePixel(mem, lcdX, lcdY, highTileMap, lowTileData)
}

func derivePixel(mem mem.Memory, lcdX, lcdY uint8) uint8 {
	// if !lcdDisplayEnable(mem) {
	// 	return 0
	// }
	if pixel, found := deriveSpritePixel(mem, lcdX, lcdY); found {
		return pixel + 0x30 // Colour offset
	}
	if pixel, found := deriveWindowPixel(mem, lcdX, lcdY); found {
		return pixel + 0x20 // Colour offset
	}
	return deriveBackgroundPixel(mem, lcdX, lcdY) + 0x10 // Colour offset
}

func (lcd *LCD) updateLcdLine(mem mem.Memory, lcdY uint8) {
	var lcdX uint8
	var index uint16
	// FIXME search for all sprites on this Y line here
	for lcdX = 0; lcdX < 160; lcdX++ {
		index = uint16(lcdY)*160 + uint16(lcdX)
		lcd.data[index] = derivePixel(mem, lcdX, lcdY)
	}
}
