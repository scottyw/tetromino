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
	data     [160 * 144]uint8
}

// NewLCD returns the configured LCD
func NewLCD(hwr *mem.HardwareRegisters, memory *mem.Memory) *LCD {
	return &LCD{
		hwr:      hwr,
		videoRAM: memory.VideoRAM(),
		oam:      memory.OAM(),
	}
}

// Tick runs the LCD driver for one machine cycle i.e. 4 clock cycles
func (lcd *LCD) Tick(cycle int) {
	lcd.hwr.LY = uint8(cycle / 114)
	lyRemainder := cycle % 114
	// Set mode on stat register
	switch {
	case lcd.hwr.LY == 144:
		// V-Blank period starts
		lcd.hwr.STAT = 1
		lcd.hwr.IF |= 0x01
	case lyRemainder == 0:
		// OAM period starts
		lcd.hwr.STAT = 2
	case lyRemainder == 20:
		// LCD data transfer period starts
		lcd.hwr.STAT = 3
	case lyRemainder == 63:
		// H-Blank period starts
		lcd.hwr.STAT = 0
		if lcd.hwr.LY < 144 {
			lcd.updateLcdLine()
		}
	}
	// Set coincidence flag and coincidence interrupt on stat register
	if lcd.hwr.LY == lcd.hwr.LYC {
		lcd.hwr.STAT |= 0x44
	} else {
		lcd.hwr.STAT &^= 0x44
	}
	// Set interrupts on stat register
	switch {
	case lcd.hwr.LY == 144:
		lcd.hwr.STAT |= 0x10
	case lyRemainder == 0:
		lcd.hwr.STAT |= 0x20
	case lyRemainder == 63:
		lcd.hwr.STAT |= 0x08
	}
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

func (lcd *LCD) deriveSpritePixel(lcdX, lcdY uint8, spriteAddrs []uint8) (uint8, bool) {
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
	lcdX += scx // Overflows deliberately
	lcdY += scy // Overflows deliberately
	return lcd.deriveTilePixel(highTileMap, lowTileData, lcdX, lcdY)
}

func (lcd *LCD) derivePixel(highBgTileMap, highWindowTileMap, lowTileData bool, scx, scy, wx, wy, lcdX, lcdY uint8, spriteAddrs []uint8) uint8 {
	// if !lcdDisplayEnable(memory) {
	// 	return 0
	// }
	if pixel, found := lcd.deriveSpritePixel(lcdX, lcdY, spriteAddrs); found {
		return pixel + 0x30 // Colour offset
	}
	if pixel, found := lcd.deriveWindowPixel(highWindowTileMap, lowTileData, wx, wy, lcdX, lcdY); found {
		return pixel + 0x20 // Colour offset
	}
	pixel := lcd.deriveBackgroundPixel(highBgTileMap, lowTileData, scx, scy, lcdX, lcdY)
	return pixel + 0x10 // Colour offset
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
	for lcdX := uint8(0); lcdX < 160; lcdX++ {
		index := uint16(lcdY)*160 + uint16(lcdX)
		lcd.data[index] = lcd.derivePixel(highBgTileMap, highWindowTileMap, lowTileData, scx, scy, wx, wy, lcdX, lcdY, spriteAddrs)
	}
}
