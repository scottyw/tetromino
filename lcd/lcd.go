package lcd

import (
	"github.com/scottyw/tetromino/mem"
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
		*memory.Read(mem.IF) |= 0x01
	case ly > 144:
		// V-Blank period
		stat = 1
	case lyRemainder < 20:
		// OAM period starts
		stat = 2
	case lyRemainder < 63:
		// LCD data transfer period starts
		stat = 3
	case lyRemainder < 114:
		// H-Blank period starts
		stat = 0
		updateLcdLine(ly)
	default:
		panic("LCD driver error setting mode")
	}
	// Set coincidence flag and coincidence interrupt on stat register
	if ly == uint8(*memory.Read(0xff45)) {
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
	*memory.Read(mem.STAT) = stat
	*memory.Read(mem.LY) = ly
}

// FrameData returns the frame data as a 160x144 array of bytes where each element is a colour value between 0 and 3
func (lcd *LCD) FrameData(mem mem.Memory) [23040]uint8 {
	lcd.drawTiles(mem)
	// if windowDisplayEnable(mem) {
	// 	lcd.drawTiles(mem, highWindowTileMapDisplaySelect)
	// }
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
	tileNumber := *mem.Read(tileAddr)
	if lowTileDataSelect(mem) {
		return mem.ReadRegion(lowTileAbsoluteAddress(tileNumber), 16)
	}
	return mem.ReadRegion(highTileAbsoluteAddress(int8(tileNumber)), 16)
}

func (lcd *LCD) drawTiles2(mem mem.Memory) {
	var x, y, row, col uint16
	var pixel uint8
	for y = 0; y < 32; y++ {
		for x = 0; x < 32; x++ {
			tileData := tileData(mem, y*32+x)
			for row = 0; row < 8; row++ {
				a := tileData[row*2]
				b := tileData[row*2+1]
				for col = 0; col < 8; col++ {
					pixel = (a>>uint(7-col))&1 | ((b>>uint(7-col))&1)<<1
					index := (((y * 8) + row) * 256) + ((x * 8) + col)
					lcd.data[index] = pixel
				}
			}
		}
	}
}

func updateLcdLine(ly uint8) {
	// FIXME
}

func pixel(addr uint16, bit uint8) uint8 {
	return bit % 4
}

func (lcd *LCD) drawTiles(mem mem.Memory) {
	var lcdX, lcdY, scrX, scrY, vramX, vramY, tileX, tileY, tileOffsetX, tileOffsetY uint8
	var index, addr uint16
	for lcdY = 0; lcdY < 144; lcdY++ {
		for lcdX = 0; lcdX < 160; lcdX++ {
			index = uint16(lcdY)*160 + uint16(lcdX)
			scrX = 0            // FIXME scroll register
			scrY = 0            // FIXME scroll register
			vramX = lcdX + scrX // Overflows deliberately
			vramY = lcdY + scrY // Overflows deliberately
			tileX = vramX / 32
			tileY = vramY / 32
			tileOffsetX = vramX % 32
			tileOffsetY = vramY % 32
			addr = uint16(tileY)*512 + uint16(tileX)*16 + uint16(tileOffsetY)*2
			lcd.data[index] = pixel(addr, tileOffsetX)
		}
	}
}
